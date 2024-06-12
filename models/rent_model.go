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
