package repository

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"learning_lantern/helper"
	"learning_lantern/models"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"time"

	"gorm.io/gorm"
)

type RentRepo interface {
	CreateNewRent(user_id uint, req *models.RentRequest) (models.RentResponse, error)
	GetRentedBooks(user_id uint) ([]models.RentHistory, error)
	GetStillRentingBooks(user_id uint) ([]models.RentedResponse, error)
	GetStillRentingBookByID(user_id uint, book_id uint) (models.RentedResponse, error)
	ReturnBook(user_id, rent_id uint, p *models.RentPayment) (models.ReturnBook, error)
	VAPayment(user_id, rent_id uint, amount float64, bank_code string) (models.PaymentResponse, error)
}

func (r *Repo) CreateNewRent(user_id uint, req *models.RentRequest) (models.RentResponse, error) {

	// check book existence
	isExist, err := r.isBookExist(req.BookID)
	if err != nil {
		return models.RentResponse{}, err
	}
	if !isExist {
		return models.RentResponse{}, helper.ErrNoData
	}

	// check book stock
	isReady, err := r.isStockReady(req.BookID)
	if err != nil {
		return models.RentResponse{}, err
	}
	if !isReady {
		return models.RentResponse{}, helper.ErrNoStock
	}

	// create the data to upload
	var rent models.Rent
	rent.BookID = req.BookID
	rent.UserID = user_id
	rent.RentStatus = "pending"
	if req.RentAt == "" {
		rent.RentAt = time.Now().Truncate(time.Second)
	} else {
		layoutFormat := "2006-01-02 15:04:05"
		rentAt, err := time.Parse(layoutFormat, req.RentAt)
		if err != nil {
			return models.RentResponse{}, err
		}
		rent.RentAt = rentAt.Truncate(time.Second)
	}
	deadline := rent.RentAt
	rent.Deadline = deadline.AddDate(0, 0, req.Days)

	res := r.DB.Create(&rent)
	if res.Error != nil {
		return models.RentResponse{}, helper.ErrQuery
	}

	// reduce book's stock by one
	var book models.Book
	r.DB.First(&book, req.BookID)
	book.Stock = book.Stock - 1
	r.DB.Save(&book)

	var resp models.RentResponse
	helper.CopyNonEmptyFields(rent, &resp)

	return resp, nil
}

func (r *Repo) GetRentedBooks(user_id uint) ([]models.RentHistory, error) {
	query := `
		SELECT r.rent_id, r.book_id, r.total_price, r.rent_at, r.returned_at, DATE_PART('day', r.returned_at::TIMESTAMP - r.rent_at::TIMESTAMP) AS days_rented from users as u
		join rents as r on r.user_id = u.user_id
		where u.user_id = ? and r.rent_status = 'returned';`
	var history []models.RentHistory
	res := r.DB.Raw(query, user_id).Find(&history)
	if res.Error != nil {
		return nil, helper.ErrQuery
	}

	return history, nil
}

func (r *Repo) GetStillRentingBooks(user_id uint) ([]models.RentedResponse, error) {
	var rentingList []models.Rent
	res := r.DB.Where("rent_status = 'pending' AND user_id = ?", user_id).Find(&rentingList)
	if res.Error != nil {
		return nil, helper.ErrQuery
	}
	var responList []models.RentedResponse
	for _, rent := range rentingList {
		var resp models.RentedResponse
		helper.CopyNonEmptyFields(rent, &resp)
		responList = append(responList, resp)
	}

	return responList, nil
}

func (r *Repo) GetStillRentingBookByID(user_id uint, rent_id uint) (models.RentedResponse, error) {
	var rent models.Rent
	res := r.DB.Where("rent_status = 'pending' AND user_id = ?", user_id).First(&rent, rent_id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return models.RentedResponse{}, helper.ErrNoData
		}
		return models.RentedResponse{}, helper.ErrQuery
	}

	var resp models.RentedResponse
	helper.CopyNonEmptyFields(rent, &resp)

	return resp, nil
}

func (r *Repo) ReturnBook(user_id, rent_id uint, p *models.RentPayment) (models.ReturnBook, error) {
	var rent models.Rent
	res := r.DB.Where("rent_status = 'pending' AND user_id = ?", user_id).First(&rent, rent_id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return models.ReturnBook{}, helper.ErrNoData
		}
		return models.ReturnBook{}, helper.ErrQuery
	}

	// start transaction
	rb := models.ReturnBook{}

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// count how may days
		now := time.Now()
		returned_at := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		deadline := time.Date(rent.Deadline.Year(), rent.Deadline.Month(), rent.Deadline.Day(), 0, 0, 0, 0, time.UTC)
		rent_at := time.Date(rent.RentAt.Year(), rent.RentAt.Month(), rent.RentAt.Day(), 0, 0, 0, 0, time.UTC)

		returnPlan := int(deadline.Sub(rent_at).Hours() / 24)
		days_rented := int(returned_at.Sub(rent_at).Hours() / 24)
		if days_rented == 0 {
			days_rented++
		}

		helper.Logging(nil).Info("return plan: ", returnPlan)
		helper.Logging(nil).Info("days rented: ", days_rented)

		fineDay := days_rented - returnPlan

		// get renting price for a day
		// renting price is calculated based on deadline and fine (10% from rent price each day)
		b := models.Book{}
		r.DB.First(&b, rent.BookID)

		var fine float64 = b.RentalCost * float64(fineDay) * 1.1 // this 1.1 is fine rate
		if fine < 0 {
			// if return date not exceeded deadline, fine is 0
			fine = 0
		}
		helper.Logging(nil).Info("fine: ", fine)

		total_price := float64(returnPlan)*b.RentalCost + fine
		helper.Logging(nil).Info("total price: ", total_price)

		// get user balance
		u := models.User{UserID: user_id}
		r.DB.First(&u)
		deposit := u.Deposit
		if deposit < total_price {
			return helper.ErrUnsufficientBalance
		}

		// if user balance is enough, pay the price
		// if payment is cash , substract form user deposit
		// if payment is VA, use 3rd party API to pay
		if p.PaymentMethod == "cash" {
			u.Deposit = u.Deposit - total_price
			r.DB.Save(&u)
		} else if p.PaymentMethod == "VA" {
			paymentRespons, err := r.VAPayment(user_id, rent_id, total_price, p.BankCode)
			if err != nil {
				return err
			}
			if paymentRespons.Status != "COMPLETED" {
				return fmt.Errorf(paymentRespons.Message)
			}
		}

		// update rent status
		rent.RentStatus = "returned"
		rent.TotalPrice = total_price
		rent.ReturnedAt = time.Now().Truncate(time.Hour)
		r.DB.Save(rent)

		// increase book stock
		b.Stock = b.Stock + 1
		r.DB.Save(&b)

		// create new payment
		var payment models.Payment
		payment.RentID = rent.RentID
		payment.PaymentDate = rent.ReturnedAt
		payment.PaymentAmount = total_price
		payment.PaymentMethod = p.PaymentMethod
		res := r.DB.Create(&payment)
		if res.Error != nil {
			return res.Error
		}

		// update data for response
		rb.BookID = b.BookID
		rb.RentID = rent.RentID
		rb.RentAt = rent.RentAt
		rb.TotalPrice = total_price
		rb.ReturnedAt = returned_at
		rb.DaysRented = days_rented
		rb.PaymentMethod = p.PaymentMethod
		return nil
	})

	if err != nil {
		return models.ReturnBook{}, err
	}
	return rb, nil
}

func (r *Repo) VAPayment(user_id, rent_id uint, amount float64, bank_code string) (models.PaymentResponse, error) {
	full_name := ""
	res := r.DB.Model(&models.User{}).Select("CONCAT (ud.fname, ' ',ud.lname) as full_name").
		Joins("LEFT JOIN user_details as ud ON ud.user_id = users.user_id").Where("users.user_id = ?", user_id).Find(&full_name)
	if res.Error != nil {
		return models.PaymentResponse{}, res.Error
	}
	log.Println(full_name)
	va, err := CreateVirtualAccount(full_name, user_id, rent_id, bank_code)
	if err != nil {
		return models.PaymentResponse{}, err
	}
	log.Println(va)

	paymentRespon, err := SimulatePayment(&va, amount)
	if err != nil {
		return models.PaymentResponse{}, err
	}
	log.Println(paymentRespon)
	return paymentRespon, nil
}

func randGenerator() int {
	// Generate a 6-digit random integer
	min := 100000
	max := 999999
	return rand.Intn(max-min+1) + min
}

func CreateVirtualAccount(full_name string, user_id, rent_id uint, bank_code string) (models.VAResponse, error) {
	apiKey := "xnd_development_tT8Uulp13l497VsTfUzzdz7Jnf49qmJiDSRkcEbiYthiXtI9eiaYunUVdXNe"
	url := "https://api.xendit.co/callback_virtual_accounts"

	type bodyReq struct {
		ExternalID     string `json:"external_id"`
		BankCode       string `json:"bank_code"`
		Name           string `json:"name"`
		IsClosed       bool   `json:"is_closed"`
		ExpirationDate string `json:"expiration_date"`
		Country        string `json:"country"`
		IsSingleUse    bool   `json:"is_single_use"`
		Currency       string `json:"currency"`
	}

	var bq bodyReq
	bq.ExternalID = fmt.Sprintf("VA_fixed-%v%v%v", user_id, rent_id, randGenerator())
	bq.BankCode = bank_code
	bq.Name = full_name
	bq.IsClosed = false
	bq.ExpirationDate = time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	bq.Country = "ID"
	bq.IsSingleUse = true
	bq.Currency = "IDR"

	// marshall body
	bodyReqJSON, err := json.Marshal(bq)
	if err != nil {
		return models.VAResponse{}, err
	}
	log.Println(bq)

	// create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyReqJSON))
	if err != nil {
		return models.VAResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	// Encode API key to Base64
	auth := fmt.Sprintf("%s:", apiKey)
	b64 := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+b64)

	// req.Header.Set("Authorization", apiKey)

	// req.SetBasicAuth(apiKey, "")

	// // send the request
	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	return models.VAResponse{}, err
	// }
	// defer resp.Body.Close()
	// // open the request
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return models.VAResponse{}, err
	// }

	// // unmarshall body
	// va := models.VAResponse{}
	// if err := json.Unmarshal(body, &va); err != nil {
	// 	return models.VAResponse{}, err
	// }

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.VAResponse{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.VAResponse{}, err
	}

	// For debugging, log the request and response
	reqDump, _ := httputil.DumpRequest(req, true)
	log.Printf("REQUEST:\n%s", string(reqDump))
	log.Printf("RESPONSE:\n%s", body)

	// unmarshall body
	va := models.VAResponse{}
	if err := json.Unmarshal(body, &va); err != nil {
		return models.VAResponse{}, err
	}

	// log.Println("creating va success", va)
	return va, err
}

func SimulatePayment(va *models.VAResponse, amount float64) (models.PaymentResponse, error) {
	apiKey := "xnd_development_tT8Uulp13l497VsTfUzzdz7Jnf49qmJiDSRkcEbiYthiXtI9eiaYunUVdXNe"
	url := fmt.Sprintf("https://api.xendit.co/callback_virtual_accounts/external_id=%s/simulate_payment", va.ExternalID)
	log.Println(url)
	type bReq struct {
		Amount float64 `json:"amount"`
	}

	bq := bReq{Amount: amount}
	// marshall body
	bodyReqJSON, err := json.Marshal(&bq)
	if err != nil {
		return models.PaymentResponse{}, err
	}
	log.Println(bodyReqJSON)

	// create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyReqJSON))
	if err != nil {
		return models.PaymentResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	// Encode API key to Base64
	auth := fmt.Sprintf("%s:", apiKey)
	b64 := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+b64)

	// req.Header.Add("Authorization", "xnd_development_cR3ylqJJEUrcEesGiQFO4pPxPalg4yj2lYEHFUcrKiP0v6Wo5lNlL9ugbSCGHVOZ")

	// // send the request
	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	return models.PaymentResponse{}, err
	// }
	// defer resp.Body.Close()

	// // open the request
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return models.PaymentResponse{}, err
	// }

	// // unmarshall body
	// pr := models.PaymentResponse{}
	// if err := json.Unmarshal(body, &pr); err != nil {
	// 	return models.PaymentResponse{}, err
	// }

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.PaymentResponse{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.PaymentResponse{}, err
	}

	// For debugging, log the request and response
	reqDump, _ := httputil.DumpRequest(req, true)
	log.Printf("REQUEST:\n%s", string(reqDump))
	log.Printf("RESPONSE:\n%s", body)

	// log.Println("creating va success", va)
	// return models.PaymentResponse{}, err

	// unmarshall body
	pr := models.PaymentResponse{}
	if err := json.Unmarshal(body, &pr); err != nil {
		return models.PaymentResponse{}, err
	}
	return pr, nil
}
