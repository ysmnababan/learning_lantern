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
			"message": "Login success",
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

	if GetU.Email == "" || GetU.Password == "" || GetU.Username == "" {
		return helper.ParseError(helper.ErrParam, c)
	}

	// validate role
	if GetU.Role == "" {
		GetU.Role = "user"
	} else {
		if GetU.Role != "user" && GetU.Role != "admin" {
			return helper.ParseError(helper.ErrParam, c)
		}
	}

	respU, err := s.UserRepo.Register(GetU)
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "New User Created", "User": respU})
}

// GetUserInfo godoc
// @Summary Get info about a user
// @Description must be authenticated user and return user detail data
// @Tags User
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Success 200 {object} models.UserDetailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/user [get]
func (s *UserController) GetUserInfo(c echo.Context) error {
	cred := helper.GetCredential(c)
	resp, err := s.UserRepo.GetInfo(cred.UserID)
	if err != nil {
		return helper.ParseError(err, c)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Get User Info", "User": resp})
}

// GetAllUser godoc
// @Summary Get info about a user ONLY FOR ADMIN
// @Description must be authenticated user and return all user data
// @Tags User
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Success 200 {array} models.UserDetailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/users [get]
func (s *UserController) GetAllUser(c echo.Context) error {
	cred := helper.GetCredential(c)
	if cred.Role != "admin" {
		return helper.ParseError(helper.ErrMustAdmin, c)
	}

	resp, err := s.UserRepo.GetAllUser(cred.UserID)
	if err != nil {
		return helper.ParseError(err, c)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Get All User", "User": resp})
}

// UpdateUser godoc
// @Summary Update user information
// @Description must be authenticated user and update detail info of a user
// @Tags User
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Param student body models.UserUpdateRequest true "Data to be updated"
// @Success 200 {object} models.UserDetailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/user [put]
func (s *UserController) UpdateUser(c echo.Context) error {
	cred := helper.GetCredential(c)
	var GetU models.UserUpdateRequest
	err := c.Bind(&GetU)
	if err != nil {
		return helper.ParseError(helper.ErrBindJSON, c)
	}

	//validate user
	if GetU.Username == "" {
		return helper.ParseError(helper.ErrParam, c)
	}

	respU, err := s.UserRepo.Update(cred.UserID, GetU)
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "User Data Updated", "User": respU})
}

// TopUpDeposit godoc
// @Summary Update user information
// @Description must be authenticated user and update detail info of a user
// @Tags User
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Param student body models.TopUpReq true "Data to be updated"
// @Success 200 {object} map[string]interface{} "message : string, Amount of Deposit: float64"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/user/topup [put]
func (s *UserController) TopUpDeposit(c echo.Context) error {
	cred := helper.GetCredential(c)
	if cred.Role != "user" {
		return helper.ParseError(helper.ErrOnlyUser, c)
	}
	var GetD models.TopUpReq
	err := c.Bind(&GetD)
	if err != nil {
		return helper.ParseError(helper.ErrBindJSON, c)
	}

	//validate deposit
	if GetD.Deposit <= 0 {
		return helper.ParseError(helper.ErrParam, c)
	}

	respU, err := s.UserRepo.TopUp(cred.UserID, GetD.Deposit)
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Top Up success", "Amount of Deposit": respU})
}
