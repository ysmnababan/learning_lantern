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

// RentBook godoc
// @Summary Rent a book [ONLY FOR USER]
// @Description Rent a book in a library, 1-30 days
// @Tags Rents
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Param student body models.RentRequest true "Books to rent"
// @Success 200 {object} models.RentResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/rent [post]
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

// MyRentHistory godoc
// @Summary History of rented book [ONLY FOR USER]
// @Description List of books those rented by user (already returned)
// @Tags History
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Success 200 {array} models.RentHistory
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/history/rent [get]
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

// MyRentedBooks godoc
// @Summary List of rented book [ONLY FOR USER]
// @Description List of books those rented by user (not returned yet)
// @Tags Rents
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Success 200 {array} models.RentedResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/rents [get]
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

// DetailRentedBook godoc
// @Summary Detail of rented book [ONLY FOR USER]
// @Description Detail of a book that rented by user (not returned yet)
// @Tags Rents
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Param id path string true "Rent ID"
// @Success 200 {objcect} models.RentedResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/rent/{id} [get]
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

// ReturnBookCash godoc
// @Summary Return book by rent id [ONLY FOR USER]
// @Description Return book by rent id and payment using cash
// @Tags Rents
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Param id path string true "Rent ID"
// @Success 200 {objcect} models.ReturnBook
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/rent/return_cash/{id} [post]
func (s *RentController) ReturnBookCash(c echo.Context) error {
	cred := helper.GetCredential(c)
	if cred.Role != "user" {
		return helper.ParseError(helper.ErrOnlyUser, c)
	}

	// get rent id
	rent_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.ParseError(helper.ErrInvalidId, c)
	}

	resp, err := s.RentRepo.ReturnBookCash(cred.UserID, uint(rent_id))
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Book returned successfully", "Detail": resp})

}

// ReturnBookVA godoc
// @Summary Return book by rent id [ONLY FOR USER]
// @Description Return book by rent id and payment using virtual account
// @Tags Rents
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Param student body models.RentPayment true "Payment method (VA) and Bank Code, BRI, BNI, MANDIRI"
// @Param id path string true "Rent ID"
// @Success 200 {objcect} models.ReturnBookVA
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/rent/return_va/{id} [post]
func (s *RentController) ReturnBookVA(c echo.Context) error {
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
	if !(GetR.PaymentMethod == "" || GetR.PaymentMethod == "VA") {
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

	resp, err := s.RentRepo.ReturnBookVA(cred.UserID, uint(rent_id), &GetR)
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Book returned successfully", "Detail": resp})

}

func (s *RentController) CobaVA(c echo.Context) error {
	fmt.Println("here")

	_, resp, err := s.RentRepo.VAPayment(12, 1, 100.0, "BNI")
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Book returned successfully", "Detail": resp})
}
