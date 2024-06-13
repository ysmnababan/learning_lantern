package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"learning_lantern/models"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func randGenerator() int {
	// Generate a 6-digit random integer
	min := 100000
	max := 999999
	return rand.Intn(max-min+1) + min
}

func CreateVirtualAccount(full_name string, user_id, rent_id uint, bank_code string) (models.VAResponse, error) {
	err := godotenv.Load()
	if err != nil {
		return models.VAResponse{}, fmt.Errorf("error opening .env")
	}
	apiKey := os.Getenv("API_KEY")
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

	return va, err
}

func SimulatePayment(va *models.VAResponse, amount float64) (models.PaymentResponse, error) {
	err := godotenv.Load()
	if err != nil {
		return models.PaymentResponse{}, fmt.Errorf("error opening .env")
	}
	apiKey := os.Getenv("API_KEY")
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

	// unmarshall body
	pr := models.PaymentResponse{}
	if err := json.Unmarshal(body, &pr); err != nil {
		return models.PaymentResponse{}, err
	}
	return pr, nil
}
