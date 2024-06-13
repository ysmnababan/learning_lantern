package models

import "time"

type Payment struct {
	PaymentID     uint      `json:"payment_id" gorm:"primaryKey;autoIncrement"`
	RentID        uint      `json:"rent_id" gorm:"unique;not null;"`
	PaymentDate   time.Time `json:"payment_date" gorm:"type:timestamp;not null"`
	PaymentAmount float64   `json:"payment_amount" gorm:"type:decimal(10,2);not null"`
	PaymentMethod string    `json:"payment_method" gorm:"type:varchar(50);not null"`
}

type RentPayment struct {
	PaymentMethod string `json:"payment_method"`
	BankCode      string `json:"bank_code"`
}

type PaymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type VAResponse struct {
	ID             string    `json:"id"`
	ExternalID     string    `json:"external_id"`
	OwnerID        string    `json:"owner_id"`
	BankCode       string    `json:"bank_code"`
	MerchantCode   string    `json:"merchant_code"`
	AccountNumber  string    `json:"account_number"`
	Name           string    `json:"name"`
	IsSingleUse    bool      `json:"is_single_use"`
	IsClosed       bool      `json:"is_closed"`
	ExpirationDate time.Time `json:"expiration_date"`
	Status         string    `json:"status"`
	Currency       string    `json:"currency"`
	Country        string    `json:"country"`
}
