package router

import (
	"learning_lantern/controller"
	"learning_lantern/helper"
	"learning_lantern/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)



func SetupRouter(e *echo.Echo, db *gorm.DB) {
	// using logger for each api
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			helper.Logging(c).Info("Calling an API")
			return next(c)
		}
	})
	e.Use(middleware.Recover())

	repo := &repository.Repo{DB: db}

	uc := &controller.UserController{UserRepo: repo}
	bc := &controller.BookController{BookRepo: repo}
	rc := &controller.RentController{RentRepo: repo}
	rp := &controller.PaymentController{PaymentRepo: repo}

	// no need authorization
	e.POST("/api/users/register", uc.Register)
	e.POST("/api/users/login", uc.Login)

	service := e.Group("/api")

	// authentification middleware
	service.Use(helper.Auth)
	{
		// for user models
		service.GET("/user", uc.GetUserInfo)
		service.GET("/users", uc.GetAllUser)
		service.PUT("/user", uc.UpdateUser)
		service.PUT("/user/topup", uc.TopUpDeposit)

		// for book models
		service.GET("/books", bc.ListAllBooks)
		service.GET("/books/available", bc.ListAvailableBooks)
		service.GET("/books/unavailable", bc.ListOfUnavailableBooks)
		service.POST("/book", bc.AddNewBook)
		service.PUT("/book/:id", bc.EditBook)
		service.DELETE("/book/:id", bc.DeleteBook)

		// for user rent models
		service.POST("/rent", rc.RentBook)
		service.GET("/rents", rc.MyRentedBooks)
		service.GET("/rent/:id", rc.DetailRentedBook)
		service.POST("/rent/return_cash/:id", rc.ReturnBookCash)
		service.POST("/rent/return_va/:id", rc.ReturnBookVA)

		// for history
		service.GET("/history/rent", rc.MyRentHistory)
		service.GET("/history/revenue", rp.GetTotalRevenue)

		// test 3rd party api
		service.GET("/cobava", rc.CobaVA)
	}
}
