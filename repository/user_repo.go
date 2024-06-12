package repository

import (
	"fmt"
	"learning_lantern/helper"
	"learning_lantern/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo interface {
	Login(u models.User) (string, error)
	Register(u models.User) (models.UserResponse, error)
}

type Repo struct {
	DB *gorm.DB
}

func (r *Repo) isUserExist(email string) (bool, error) {
	var user models.User
	res := r.DB.Where("email= ?", email).First(&user)
	// user exist
	if res.Error == nil {
		return true, nil
	}

	// error query
	if res.Error != gorm.ErrRecordNotFound {
		return false, helper.ErrQuery
	}

	return false, nil

}

func generateToken(u models.User) (string, error) {
	// create the payload
	payload := jwt.MapClaims{
		"id":    u.UserID,
		"email": u.Email,
		"role":  u.Role,
		"exp":   time.Now().Add(time.Hour * 48).Unix(),
	}

	// define the method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// get token string
	_ = godotenv.Load()
	tokenString, err := token.SignedString([]byte(os.Getenv("KEY")))
	if err != nil {
		return "", fmt.Errorf("error when creating token: %v", err)
	}

	return tokenString, nil
}

func (r *Repo) Login(u models.User) (string, error) {
	isExist, err := r.isUserExist(u.Email)
	if err != nil {
		return "", err
	}

	if !isExist {
		return "", helper.ErrCredential
	}

	// var u models.User
	var getU models.User
	r.DB.Where("email=?", u.Email).First(&getU)

	//check pwd
	err = bcrypt.CompareHashAndPassword([]byte(getU.Password), []byte(u.Password))
	if err != nil {
		return "", helper.ErrCredential
	}

	// generate token
	return generateToken(getU)

}

func (r *Repo) Register(u models.User) (models.UserResponse, error) {
	var newU models.UserResponse

	isExist, err := r.isUserExist(u.Email)
	if err != nil {
		return models.UserResponse{}, err
	}

	if isExist {
		return models.UserResponse{}, helper.ErrUserExists
	}

	hashedpwd, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedpwd)
	res := r.DB.Create(&u)
	if res.Error != nil {
		return models.UserResponse{}, helper.ErrQuery
	}

	// return response
	newU.UserID = u.UserID
	newU.Email = u.Email
	// newU.FullName = u.FullName
	// newU.Weight = u.Weight
	// newU.Height = u.Height

	return newU, nil
}

// func (r *Repo) GetInfo(id uint) (models.UserResponseDetail, error) {
// 	var newU models.UserResponseDetail
// 	var u models.User
// 	res := r.DB.First(&u, id)
// 	if res.Error != nil {
// 		return models.UserResponseDetail{}, helper.ErrQuery
// 	}

// 	// return response
// 	newU.UserID = u.UserID
// 	newU.Email = u.Email
// 	newU.FullName = u.FullName
// 	newU.Weight = u.Weight
// 	newU.Height = u.Height

// 	bmi, err := getBMIIndex(u.Weight, u.Height)
// 	if err != nil {
// 		return models.UserResponseDetail{}, err
// 	}

// 	newU.BMI = bmi
// 	return newU, nil
// }
