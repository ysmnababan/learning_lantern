package repository

import (
	"learning_lantern/helper"
	"learning_lantern/models"

	"github.com/jinzhu/copier"
)

type BookRepo interface {
	GetAllBooks() ([]models.Book, error)
	GetAllAvailableBooks() ([]models.BookAvailable, error)
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
