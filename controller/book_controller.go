package controller

import (
	"fmt"
	"learning_lantern/helper"
	"learning_lantern/models"
	"learning_lantern/repository"
	"net/http"
	"reflect"
	"strconv"

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
	cred := helper.GetCredential(c)
	if cred.Role != "admin" {
		return helper.ParseError(helper.ErrMustAdmin, c)
	}
	// get book id
	book_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helper.ParseError(helper.ErrInvalidId, c)
	}

	var GetB models.BookRequest
	err = c.Bind(&GetB)
	if err != nil {
		return helper.ParseError(helper.ErrBindJSON, c)
	}

	// PrintFieldsAndValues(GetB)
	//validate books
	// if GetB.BookName == "" || GetB.Stock < 0 || GetB.RentalCost < 0 || GetB.Author == "" {
	// 	return helper.ParseError(helper.ErrParam, c)
	// }

	// var book models.Book
	// // book.BookName = "wow"
	// helper.CopyNonEmptyFields(GetB, &book)
	// log.Println(book)

	// // copier.Copy(&book, GetB)
	// book.BookID = uint(book_id)
	resp, err := s.BookRepo.UpdateBook(uint(book_id), &GetB)
	if err != nil {
		return helper.ParseError(err, c)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"Message": "Book Updated", "Books": resp})
}
func PrintFieldsAndValues(obj interface{}) {
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	// Ensure the object is a struct
	if typ.Kind() != reflect.Struct {
		fmt.Println("The provided object is not a struct")
		return
	}

	// Iterate over the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Get the zero value for the field's type
		zeroValue := reflect.Zero(field.Type).Interface()
		currentValue := value.Interface()

		// Check if the current value is equal to the zero value
		isEmpty := reflect.DeepEqual(currentValue, zeroValue)

		if isEmpty {
			fmt.Printf("%s: (empty)\n", field.Name)
		} else {
			fmt.Printf("%s: %v\n", field.Name, currentValue)
		}
	}
}
func (s *BookController) DeleteBook(c echo.Context) error {
	return nil
}
