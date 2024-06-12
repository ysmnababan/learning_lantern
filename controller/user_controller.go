package controller

import (
	"learning_lantern/helper"
	"learning_lantern/models"
	"learning_lantern/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	repository.UserRepo
}

// Login godoc
// @Summary Login as user
// @Description login as user and generate token
// @Tags User
// @Accept  json
// @Produce  json
// @Param student body models.UserRequest true "Login using email and password"
// @Success 200 {object} map[string]interface{} "message : string, token: string"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/users/login [post]
func (s *UserController) Login(c echo.Context) error {
	var GetU models.User
	err := c.Bind(&GetU)
	if err != nil {
		return helper.ParseError(helper.ErrBindJSON, c)
	}

	//validate user
	if GetU.Email == "" || GetU.Password == "" {
		return helper.ParseError(helper.ErrParam, c)
	}

	tokenString, err := s.UserRepo.Login(GetU)
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"message": "login success",
			"token":   tokenString,
		})
}

// Register godoc
// @Summary Register as user
// @Description Register as user and return user data
// @Tags User
// @Accept  json
// @Produce  json
// @Param student body models.UserRegister true "Register new user"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/users/register [post]
func (s *UserController) Register(c echo.Context) error {
	var GetU models.User
	err := c.Bind(&GetU)
	if err != nil {
		return helper.ParseError(helper.ErrBindJSON, c)
	}

	//validate user
	// if GetU.Email == "" || GetU.FullName == "" || GetU.Password == "" || GetU.Weight <= 0 || GetU.Height <= 0 {
	if GetU.Email == "" {
		return helper.ParseError(helper.ErrParam, c)
	}

	respU, err := s.UserRepo.Register(GetU)
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "New User Created", "User": respU})
}

// GetUserInfo godoc
// @Summary Get info about a user
// @Description must be authenticated user and return user detail data using third party API
// @Tags User
// @Accept  json
// @Produce  json
// @Param   Auth  header  string  true  "Authentication token"  default()
// @Success 200 {object} models.UserResponseDetail
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/users [get]
// func (s *UserController) GetUserInfo(c echo.Context) error {
// 	ctx_id := c.Get("id")
// 	id := uint(ctx_id.(float64))
// 	respU, err := s.UserRepo.GetInfo(id)
// 	if err != nil {
// 		return helper.ParseError(err, c)
// 	}
// 	return c.JSON(http.StatusOK, respU)
// }
