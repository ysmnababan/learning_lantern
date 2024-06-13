package controller

import (
	"learning_lantern/helper"
	"learning_lantern/models"
	"learning_lantern/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

type BookController struct {
	repository.BookRepo
}

// ListAllBooks godoc
// @Summary Get all books
// @Description Get all books
// @Tags Books
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Success 200 {array} models.Book
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/books [get]
func (s *BookController) ListAllBooks(c echo.Context) error {
	books, err := s.BookRepo.GetAllBooks()
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Get All Books", "Books": books})
}

// ListAvailableBooks godoc
// @Summary Get all books
// @Description Get all books that can be rented, stock >0
// @Tags Books
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Success 200 {array} models.BookAvailable
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/books/available [get]
func (s *BookController) ListAvailableBooks(c echo.Context) error {
	books, err := s.BookRepo.GetAllAvailableBooks()
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Get All Available Books", "Books": books})
}

// AddNewBook godoc
// @Summary Add new book to library [ONLY FOR ADMIN]
// @Description Add new book to library
// @Tags Books
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Param student body models.BookRequest true "Books to be inserted"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/book [post]
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

// EditBook godoc
// @Summary Edit book to library [ONLY FOR ADMIN]
// @Description Edit book to library but you can insert body data that you need to update only
// @Tags Books
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Param id path string true "Book ID"
// @Param student body models.BookRequest true "Books to be inserted"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/book/{id} [put]
func (s *BookController) EditBook(c echo.Context) error {
	cred := helper.GetCredential(c)
	if cred.Role != "admin" {
		return helper.ParseError(helper.ErrMustAdmin, c)
	}
	// get book id
	book_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.ParseError(helper.ErrInvalidId, c)
	}
	log.Println(book_id)
	var GetB models.BookRequest
	err = c.Bind(&GetB)
	if err != nil {
		return helper.ParseError(helper.ErrBindJSON, c)
	}

	resp, err := s.BookRepo.UpdateBook(uint(book_id), &GetB)
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Book Updated", "Books": resp})
}

// DeleteBook godoc
// @Summary Delete book [ONLY FOR ADMIN]
// @Description Delete book from library database
// @Tags Books
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/book/{id} [delete]
func (s *BookController) DeleteBook(c echo.Context) error {
	cred := helper.GetCredential(c)
	if cred.Role != "admin" {
		return helper.ParseError(helper.ErrMustAdmin, c)
	}
	// get book id
	book_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.ParseError(helper.ErrInvalidId, c)
	}
	resp, err := s.BookRepo.DeleteBook(uint(book_id))
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Book Deleted", "Books": resp})
}

// ListOfUnavailableBooks godoc
// @Summary Get all unavailable books [ONLY FOR USER]
// @Description Get all books that can be rented but now is out of stock because another user still using it
// @Tags Books
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Success 200 {array} models.BookUnavailable
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/books/unavailable [get]
func (s *BookController) ListOfUnavailableBooks(c echo.Context) error {
	cred := helper.GetCredential(c)
	if cred.Role != "user" {
		return helper.ParseError(helper.ErrMustAdmin, c)
	}

	books, err := s.BookRepo.GetUnavailableBooks()
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Get All Unavailable Books", "Books": books})

}
