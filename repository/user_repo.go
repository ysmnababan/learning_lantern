package repository

import (
	"fmt"
	"learning_lantern/helper"
	"learning_lantern/models"
	"log"
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
	GetInfo(user_id uint) (models.UserDetailResponse, error)
	GetAllUser(user_id uint) ([]models.UserDetailResponse, error)
	Update(user_id uint, u models.UserUpdateRequest) (models.UserDetailResponse, error)
	TopUp(user_id uint, amount float64) (float64, error)
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
	token, err := generateToken(getU)
	if err != nil {
		return "", err
	}

	// update token and login time in db
	getU.JwtToken = token
	getU.LastLoginDate = time.Now().Truncate(time.Second) //'2023-06-10 12:34:56'
	res := r.DB.Save(&getU)
	if res.Error != nil {
		return "", helper.ErrQuery
	}

	return token, nil
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

	// create userdetail in userdetail table using emtpy data
	var ud models.UserDetail
	ud.UserID = u.UserID
	res = r.DB.Create(&ud)
	if res.Error != nil {
		return models.UserResponse{}, helper.ErrQuery
	}

	// return response
	newU.UserID = u.UserID
	newU.Email = u.Email
	newU.Username = u.Username
	newU.Deposit = u.Deposit
	return newU, nil
}

func (r *Repo) GetInfo(user_id uint) (models.UserDetailResponse, error) {
	var respU models.UserDetailResponse

	var u models.User
	res := r.DB.First(&u, user_id)
	if res.Error != nil {
		return models.UserDetailResponse{}, helper.ErrQuery
	}
	var ud models.UserDetail
	res = r.DB.Where("user_id=?", user_id).First(&ud)
	if res.Error != nil {
		return models.UserDetailResponse{}, helper.ErrQuery
	}

	// return response
	respU.UserID = u.UserID
	respU.Username = u.Username
	respU.Email = u.Email
	respU.Deposit = u.Deposit
	respU.Fname = ud.Fname
	respU.Lname = ud.Lname
	respU.Address = ud.Address
	respU.Age = ud.Age
	respU.PhoneNumber = ud.PhoneNumber

	return respU, nil
}

func (r *Repo) GetAllUser(user_id uint) ([]models.UserDetailResponse, error) {
	var users []models.User
	res := r.DB.Where("role = ?", "user").Find(&users)
	if res.Error != nil {
		return nil, helper.ErrQuery
	}
	var alluser []models.UserDetailResponse
	for _, user := range users {
		u, err := r.GetInfo(uint(user.UserID))
		if err != nil {
			return nil, err
		}
		alluser = append(alluser, u)
	}

	return alluser, nil
}

func (r *Repo) Update(user_id uint, u models.UserUpdateRequest) (models.UserDetailResponse, error) {
	var user models.User
	res := r.DB.First(&user, user_id)
	if res.Error != nil {
		return models.UserDetailResponse{}, helper.ErrQuery
	}

	// update username
	user.Username = u.Username
	res = r.DB.Save(&user)
	if res.Error != nil {
		return models.UserDetailResponse{}, helper.ErrQuery
	}

	// update detail of user
	var updateU models.UserDetail
	res = r.DB.Where("user_id=?", user_id).First(&updateU)
	if res.Error != nil {
		return models.UserDetailResponse{}, helper.ErrQuery
	}
	updateU.Fname = u.Fname
	updateU.Lname = u.Lname
	updateU.Address = u.Address
	updateU.Age = u.Age
	updateU.PhoneNumber = u.PhoneNumber
	res = r.DB.Save(&updateU)
	if res.Error != nil {
		return models.UserDetailResponse{}, helper.ErrQuery
	}

	return r.GetInfo(user_id)
}

func (r *Repo) TopUp(user_id uint, amount float64) (float64, error) {
	var user models.User
	res := r.DB.First(&user, user_id)
	if res.Error != nil {
		return 0, helper.ErrQuery
	}
	log.Println("HERE")
	// update user deposit
	user.Deposit = user.Deposit + amount
	res = r.DB.Save(&user)
	if res.Error != nil {
		return 0, helper.ErrQuery
	}
	return user.Deposit, nil
}
