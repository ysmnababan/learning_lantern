package controller

import (
	"fmt"
	"learning_lantern/helper"
	"learning_lantern/models"
	"learning_lantern/repository"
	"net/http"
	"strconv"

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
	cred := helper.GetCredential(c)
	if cred.Role != "user" {
		return helper.ParseError(helper.ErrOnlyUser, c)
	}

	resp, err := s.RentRepo.GetRentedBooks(uint(cred.UserID))
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "History of Book Rented", "Books": resp})
}

func (s *RentController) MyRentedBooks(c echo.Context) error {
	cred := helper.GetCredential(c)
	if cred.Role != "user" {
		return helper.ParseError(helper.ErrOnlyUser, c)
	}

	resp, err := s.RentRepo.GetStillRentingBooks(cred.UserID)
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Book that is still renting", "Books": resp})
}

func (s *RentController) DetailRentedBook(c echo.Context) error {
	cred := helper.GetCredential(c)
	if cred.Role != "user" {
		return helper.ParseError(helper.ErrOnlyUser, c)
	}

	// get rent id
	rent_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.ParseError(helper.ErrInvalidId, c)
	}

	resp, err := s.RentRepo.GetStillRentingBookByID(cred.UserID, uint(rent_id))
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Book that is still renting", "Books": resp})
}

func (s *RentController) ReturnBook(c echo.Context) error {
	cred := helper.GetCredential(c)
	if cred.Role != "user" {
		return helper.ParseError(helper.ErrOnlyUser, c)
	}

	// get rent id
	rent_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.ParseError(helper.ErrInvalidId, c)
	}

	// get request data
	var GetR models.RentPayment
	err = c.Bind(&GetR)
	if err != nil {
		return helper.ParseError(helper.ErrBindJSON, c)
	}

	//validate payment requirement
	if GetR.PaymentMethod == "" || !(GetR.PaymentMethod == "cash" || GetR.PaymentMethod == "VA") {
		return helper.ParseError(helper.ErrParam, c)
	}

	if GetR.PaymentMethod == "VA" && !(GetR.BankCode == "BCA" ||
		GetR.BankCode == "BNI" ||
		GetR.BankCode == "BRI" ||
		GetR.BankCode == "BJB" ||
		GetR.BankCode == "BSI" ||
		GetR.BankCode == "BNC" ||
		GetR.BankCode == "CIMB" ||
		GetR.BankCode == "DBS" ||
		GetR.BankCode == "MANDIRI" ||
		GetR.BankCode == "PERMATA" ||
		GetR.BankCode == "SAHABAT_SAMPOERNA") {
		return helper.ParseError(helper.ErrParam, c)
	}

	resp, err := s.RentRepo.ReturnBook(cred.UserID, uint(rent_id), &GetR)
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Book returned successfully", "Detail": resp})

}

func (s *RentController) CobaVA(c echo.Context) error {
	fmt.Println("here")

	resp, err := s.RentRepo.VAPayment(12, 1, 100.0, "BNI")
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Book returned successfully", "Detail": resp})
}
