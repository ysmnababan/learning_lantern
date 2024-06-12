package models

import "time"

type Rent struct {
	RentID     uint      `json:"rent_id" gorm:"primaryKey;autoIncrement"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	BookID     uint      `json:"book_id" gorm:"not null"`
	TotalPrice float64   `json:"total_price" gorm:"type:decimal(10,2)"`
	RentStatus string    `json:"rent_status" gorm:"not null"`
	RentAt     time.Time `json:"rent_at" gorm:"not null"`
	Deadline   time.Time `json:"deadline"`
	ReturnedAt time.Time `json:"returned_at"`
}

type RentRequest struct {
	BookID uint   `json:"book_id"`
	RentAt string `json:"rent_at"`
	Days   int    `json:"days"`
}

type RentResponse struct {
	RentID     uint      `json:"rent_id" gorm:"primaryKey;autoIncrement"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	BookID     uint      `json:"book_id" gorm:"not null"`
	RentStatus string    `json:"rent_status" gorm:"not null"`
	RentAt     time.Time `json:"rent_at" gorm:"not null"`
	Deadline   time.Time `json:"deadline"`
}

type RentedResponse struct {
	RentID     uint      `json:"rent_id" gorm:"primaryKey;autoIncrement"`
	BookID     uint      `json:"book_id" gorm:"not null"`
	RentStatus string    `json:"rent_status" gorm:"not null"`
	RentAt     time.Time `json:"rent_at" gorm:"not null"`
	Deadline   time.Time `json:"deadline"`
}

type RentHistory struct {
	RentID     uint      `json:"rent_id"`
	BookID     uint      `json:"book_id"`
	RentAt     time.Time `json:"rent_at" gorm:"not null"`
	TotalPrice float64   `json:"total_price" gorm:"type:decimal(10,2)"`
	ReturnedAt time.Time `json:"returned_at"`
	DaysRented int       `json:"days_rented"`
}
