package repository

import (
	"learning_lantern/helper"
	"learning_lantern/models"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type BookRepo interface {
	GetAllBooks() ([]models.Book, error)
	GetAllAvailableBooks() ([]models.BookAvailable, error)
	CreateBook(b *models.Book) error
}

func (r *Repo) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	res := r.DB.Find(&books)
	if res.Error != nil {
		return nil, helper.ErrQuery
	}

	return books, nil
}

func (r *Repo) GetAllAvailableBooks() ([]models.BookAvailable, error) {
	var books []models.Book
	res := r.DB.Where("stock>0").Find(&books)
	if res.Error != nil {
		return nil, helper.ErrQuery
	}

	var bAvail []models.BookAvailable
	for _, book := range books {
		var b models.BookAvailable
		copier.Copy(&b, &book)
		bAvail = append(bAvail, b)
	}
	return bAvail, nil
}

func (r *Repo) isBookUnique(bookname, author string) (bool, error) {
	// var book models.Book
	res := r.DB.Where("book_name = ? AND author = ?", bookname, author).First(&models.Book{})
	// combibation of book and author exist
	if res.Error == nil {
		return false, nil
	}

	// error query
	if res.Error != gorm.ErrRecordNotFound {
		return false, helper.ErrQuery
	}

	return true, nil
}

func (r *Repo) CreateBook(b *models.Book) error {
	isUnique, err := r.isBookUnique(b.BookName, b.Author)
	if err != nil {
		return helper.ErrQuery
	}

	if !isUnique {
		return helper.ErrAuthorBookUQ
	}
	res := r.DB.Create(b)
	if res.Error != nil {
		return helper.ErrQuery
	}
	return nil
}
