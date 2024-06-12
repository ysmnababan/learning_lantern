package controller

import (
	"learning_lantern/helper"
	"learning_lantern/repository"
	"net/http"

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
	return nil
}

func (s *BookController) EditBook(c echo.Context) error {
	return nil
}

func (s *BookController) DeleteBook(c echo.Context) error {
	return nil
}
