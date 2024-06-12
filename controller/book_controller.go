package controller

import (
	"learning_lantern/helper"
	"learning_lantern/models"
	"learning_lantern/repository"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type BookController struct {
	repository.BookRepo
}

func (s *BookController) ListAllBooks(c echo.Context) error {
	books, err := s.BookRepo.GetAllBooks()
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Get All Books", "Books": books})
}

func (s *BookController) ListAvailableBooks(c echo.Context) error {
	books, err := s.BookRepo.GetAllAvailableBooks()
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Get All Available Books", "Books": books})
}

func (s *BookController) ListRentedBook(c echo.Context) error {
	return nil
}

func (s *BookController) AddNewBook(c echo.Context) error {
	cred := helper.GetCredential(c)
	if cred.Role != "admin" {
		return helper.ParseError(helper.ErrMustAdmin, c)
	}
	var GetB models.BookRequest
	err := c.Bind(&GetB)
	if err != nil {
		return helper.ParseError(helper.ErrBindJSON, c)
	}

	//validate books
	if GetB.BookName == "" || GetB.Stock < 0 || GetB.RentalCost < 0 || GetB.Author == "" {
		return helper.ParseError(helper.ErrParam, c)
	}

	var book models.Book
	copier.Copy(&book, GetB)
	err = s.BookRepo.CreateBook(&book)
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"Message": "New Book Created", "Books": book})
}

func (s *BookController) EditBook(c echo.Context) error {
	return nil
}

func (s *BookController) DeleteBook(c echo.Context) error {
	return nil
}
