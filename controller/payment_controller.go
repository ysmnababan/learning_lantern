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

// GetTotalRevenue godoc
// @Summary Get total revenue [ONLY FOR ADMIN]
// @Description Get total revenue from all returned books
// @Tags History
// @Accept  json
// @Produce  json
// @Param   Authorization  header  string  true  "Authentication token"  default()
// @Success 200 {object} map[string]interface{} "message : string, Revenue(USD): float64"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/history/revenue [get]
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
