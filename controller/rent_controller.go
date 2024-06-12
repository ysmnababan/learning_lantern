package controller

import (
	"learning_lantern/helper"
	"learning_lantern/models"
	"learning_lantern/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RentController struct {
	repository.RentRepo
}

func (s *RentController) RentBook(c echo.Context) error {
	cred := helper.GetCredential(c)
	if cred.Role != "user" {
		return helper.ParseError(helper.ErrOnlyUser, c)
	}

	// get request data
	var GetR models.RentRequest
	err := c.Bind(&GetR)
	if err != nil {
		return helper.ParseError(helper.ErrBindJSON, c)
	}

	//validate rent requirement
	if GetR.Days <= 0 || GetR.Days > 30 || GetR.BookID <= 0 {
		return helper.ParseError(helper.ErrParam, c)
	}

	resp, err := s.RentRepo.CreateNewRent(cred.UserID, &GetR)
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Book Rent", "Books": resp})
}

func (s *RentController) MyRentHistory(c echo.Context) error {

	return nil
}

func (s *RentController) MyRentedBooks(c echo.Context) error {

	return nil
}

func (s *RentController) DetailRentedBook(c echo.Context) error {

	return nil
}

func (s *RentController) ReturnBook(c echo.Context) error {

	return nil
}
