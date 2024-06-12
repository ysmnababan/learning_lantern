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
}
