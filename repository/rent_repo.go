package repository

import (
	"learning_lantern/helper"
	"learning_lantern/models"
	"time"
)

type RentRepo interface {
	CreateNewRent(user_id uint, req *models.RentRequest) (models.RentResponse, error)
}

func (r *Repo) CreateNewRent(user_id uint, req *models.RentRequest) (models.RentResponse, error) {

	// check book existence and stock
	isExist, err := r.isBookExist(req.BookID)
	if err != nil {
		return models.RentResponse{}, err
	}
	if !isExist {
		return models.RentResponse{}, helper.ErrNoData
	}

	isReady, err := r.isStockReady(req.BookID)
	if err != nil {
		return models.RentResponse{}, err
	}
	if !isReady {
		return models.RentResponse{}, helper.ErrNoStock
	}

	// create the data to upload
	var rent models.Rent
	rent.BookID = req.BookID
	rent.UserID = user_id
	rent.RentStatus = "pending"
	if req.RentAt == "" {
		rent.RentAt = time.Now().Truncate(time.Second)
	} else {
		layoutFormat := "2006-01-02 15:04:05"
		rentAt, err := time.Parse(layoutFormat, req.RentAt)
		if err != nil {
			return models.RentResponse{}, err
		}
		rent.RentAt = rentAt.Truncate(time.Second)
	}
	deadline := rent.RentAt
	rent.Deadline = deadline.AddDate(0, 0, req.Days)

	res := r.DB.Create(&rent)
	if res.Error != nil {
		return models.RentResponse{}, helper.ErrQuery
	}

	// reduce book's stock by one
	var book models.Book
	r.DB.First(&book, req.BookID)
	book.Stock = book.Stock - 1
	r.DB.Save(&book)

	var resp models.RentResponse
	helper.CopyNonEmptyFields(rent, &resp)

	return resp, nil
}
