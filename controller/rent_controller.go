package controller

import (
	"learning_lantern/repository"

	"github.com/labstack/echo/v4"
)

type RentController struct {
	repository.RentRepo
}

func (s *RentController) RentBook(c echo.Context) error {

	return nil
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
