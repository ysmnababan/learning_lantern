package repository

import "learning_lantern/models"

type RentRepo interface {
	CreateNewRent(user_id uint, book_id uint) (models.Rent, error)
}

func (r *Repo) CreateNewRent(user_id uint, book_id uint) (models.Rent, error) {

	return models.Rent{}, nil
}
