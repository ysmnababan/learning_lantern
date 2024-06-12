package main

import (
	"learning_lantern/config"
	"learning_lantern/router"
	"os"

	// _ "learning-lantern/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	// echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a Restful API Gym App
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host secure-shore-13090-19e4736be1d4.herokuapp.com
// @BasePath /
func main() {
	db := config.Connect()

	e := echo.New()
	router.SetupRouter(e, db)

	err := godotenv.Load()
	if err != nil {
		panic("unable to open .env file:")
	}

	// e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

}
