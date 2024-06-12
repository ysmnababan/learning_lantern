package models

import "time"

type Book struct {
	BookID      uint    `json:"book_id" gorm:"primaryKey;autoIncrement"`
	BookName    string  `json:"book_name" gorm:"type:varchar(255);not null"`
	Stock       int     `json:"stock" gorm:"type:int;not null;default:0;check:stock>=0"`
	RentalCost  float64 `json:"rental_cost" gorm:"type:decimal(10,2);default:0;check:rental_cost>=0"`
	Category    string  `json:"category" gorm:"type:varchar(255)"`
	Description string  `json:"description" gorm:"type:text"`
	Author      string  `json:"author" gorm:"type:varchar(255);not null"`
	Publisher   string  `json:"publisher" gorm:"type:varchar(255)"`
}

type BookRequest struct {
	BookName    string  `json:"book_name"`
	Stock       int     `json:"stock"`
	RentalCost  float64 `json:"rental_cost"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Author      string  `json:"author"`
	Publisher   string  `json:"publisher"`
}

type BookAvailable struct {
	BookID      uint    `json:"book_id" gorm:"primaryKey;autoIncrement"`
	BookName    string  `json:"book_name" gorm:"type:varchar(255);not null"`
	RentalCost  float64 `json:"rental_cost" gorm:"type:decimal(10,2);default:0;check:rental_cost>=0"`
	Category    string  `json:"category" gorm:"type:varchar(255)"`
	Description string  `json:"description" gorm:"type:text"`
	Author      string  `json:"author" gorm:"type:varchar(255);not null"`
	Publisher   string  `json:"publisher" gorm:"type:varchar(255)"`
}

type BookUnavailable struct {
	BookID     uint      `json:"book_id" gorm:"primaryKey;autoIncrement"`
	BookName   string    `json:"book_name" gorm:"type:varchar(255);not null"`
	RentAt     time.Time `json:"rent_at" gorm:"not null"`
	Deadline time.Time `json:"deadline"`
	RentBy     string    `json:"rent_by"`
}
