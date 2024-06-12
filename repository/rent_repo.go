package repository

import (
	"learning_lantern/helper"
	"learning_lantern/models"
	"time"
)

type RentRepo interface {
	CreateNewRent(user_id uint, req *models.RentRequest) (models.RentResponse, error)
	GetRentedBooks(user_id uint) ([]models.RentHistory, error)
	GetStillRentingBooks(user_id uint) ([]models.RentedResponse, error)
}

func (r *Repo) CreateNewRent(user_id uint, req *models.RentRequest) (models.RentResponse, error) {

	// check book existence and stock
	isExist, err := r.isBookExist(req.BookID)
	if err != nil {
		return models.RentResponse{}, err
	}
	if !isExist {
		return models.RentResponse{}, helper.ErrNoData
	}

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
