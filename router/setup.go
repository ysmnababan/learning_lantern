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
	// wc := &controller.WorkoutController{WorkoutRepo: repo}
	// lc := &controller.LogController{LogRepo: repo}
	// ec := &controller.ExerciseController{ER: repo}

	// no need authorization
	e.POST("/api/users/register", uc.Register)
	e.POST("/api/users/login", uc.Login)

	service := e.Group("/api")

	// authentification middleware
	service.Use(helper.Auth)
	{
		// for user
		service.GET("/user", uc.GetUserInfo)
		service.GET("/users", uc.GetAllUser)
		service.PUT("/user", uc.UpdateUser)
		service.PUT("/user/topup", uc.TopUpDeposit)

		// // for workou
		// service.GET("/workouts", wc.GetAllWorkout)
		// service.GET("/workouts/:id", wc.GetDetailWorkout)
		// service.POST("/workouts", wc.CreateWorkout)
		// service.PUT("/workouts/:id", wc.UpdateWorkout)
		// service.DELETE("/workouts/:id", wc.DeleteWorkout)

		// // for exercise
		// service.POST("/exercises", ec.CreateExercise)
		// service.DELETE("/exercises/:id", ec.DeleteExercise)

		// // for log
		// service.POST("/logs", lc.CreateLog)
		// service.GET("/logs", lc.GetLog)
	}
}
