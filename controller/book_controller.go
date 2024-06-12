package controller

import (
	"learning_lantern/repository"

	"github.com/labstack/echo/v4"
)

type BookController struct {
	repository.BookRepo
}

func (s *BookController) ListAllBooks(c echo.Context) error {
	return nil
}

func (s *BookController) ListAvailableBooks(c echo.Context) error {
	return nil
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
