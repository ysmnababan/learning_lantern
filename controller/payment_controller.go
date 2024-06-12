package controller

import (
	"learning_lantern/helper"
	"learning_lantern/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PaymentController struct {
	repository.PaymentRepo
}

func (s *PaymentController) GetTotalRevenue(c echo.Context) error {
	cred := helper.GetCredential(c)
	if cred.Role != "admin" {
		return helper.ParseError(helper.ErrOnlyUser, c)
	}

	resp, err := s.PaymentRepo.GetRevenue()
	if err != nil {
		return helper.ParseError(err, c)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"Message": "Total Revenue", "Revenue": resp})
}
