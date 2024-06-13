package main

import (
	"learning_lantern/config"
	"learning_lantern/router"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"

	_ "learning_lantern/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is a Restful API Learning Lantern Library
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
// @host thawing-tor-78922-34d29536655f.herokuapp.com
func main() {
	db := config.Connect()

	e := echo.New()
	router.SetupRouter(e, db)

	err := godotenv.Load()
	if err != nil {
		panic("unable to open .env file:")
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

}
