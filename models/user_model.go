package models

import "time"

type User struct {
	UserID        int       `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username      string    `json:"username" gorm:"type:varchar(255);not null"`
	Email         string    `json:"email" gorm:"type:varchar(255);unique;not null"`
	Password      string    `json:"password" gorm:"type:varchar(255);not null"`
	Role          string    `json:"role" gorm:"type:varchar(50);not null"`
	Deposit       float64   `json:"deposit" gorm:"type:decimal(10,2);check:deposit>=0"`
	LastLoginDate time.Time `json:"last_login_date" gorm:"type:timestamp"`
	JwtToken      string    `json:"jwt_token" gorm:"type:text"`
}

type UserResponse struct {
	UserID   int     `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username string  `json:"username" gorm:"type:varchar(255);not null"`
	Email    string  `json:"email" gorm:"type:varchar(255);unique;not null"`
	Deposit  float64 `json:"deposit" gorm:"type:decimal(10,2);check:deposit>=0"`
}

type UserLoginResponse struct {
	UserID        int       `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username      string    `json:"username" gorm:"type:varchar(255);not null"`
	Email         string    `json:"email" gorm:"type:varchar(255);unique;not null"`
	Deposit       float64   `json:"deposit" gorm:"type:decimal(10,2);check:deposit>=0"`
	LastLoginDate time.Time `json:"last_login_date" gorm:"type:timestamp"`
	JwtToken      string    `json:"jwt_token" gorm:"type:text"`
}

type UserDetail struct {
	UserDetailID int    `json:"user_detail_id" gorm:"primaryKey;autoIncrement"`
	UserID       int    `json:"user_id" gorm:"unique;not null"`
	Fname        string `json:"fname" gorm:"type:varchar(255)"`
	Lname        string `json:"lname" gorm:"type:varchar(255)"`
	Address      string `json:"address" gorm:"type:text"`
	Age          int    `json:"age" gorm:"check:age>0"`
	PhoneNumber  string `json:"phone_number" gorm:"type:varchar(20)"`
}

type UserDetailResponse struct {
	UserID      int     `json:"user_id" gorm:"unique;not null"`
	Username    string  `json:"username" gorm:"type:varchar(255);not null"`
	Email       string  `json:"email" gorm:"type:varchar(255);unique;not null"`
	Deposit     float64 `json:"deposit" gorm:"type:decimal(10,2);check:deposit>=0"`
	Fname       string  `json:"fname" gorm:"type:varchar(255)"`
	Lname       string  `json:"lname" gorm:"type:varchar(255)"`
	Address     string  `json:"address" gorm:"type:text"`
	Age         int     `json:"age" gorm:"check:age>0"`
	PhoneNumber string  `json:"phone_number" gorm:"type:varchar(20)"`
}
