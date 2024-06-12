package repository

import (
	"learning_lantern/helper"
	"learning_lantern/models"
)

type BookRepo interface {
	GetAllBooks() ([]models.Book, error)
}

func (r *Repo) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	res := r.DB.Find(&books)
	if res.Error != nil {
		return nil, helper.ErrQuery
	}

	return books, nil
}
