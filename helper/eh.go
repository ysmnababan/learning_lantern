package helper

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrNoData        = errors.New("no data in result set")
	ErrNoUser        = errors.New("no user exist")
	ErrQuery         = errors.New("query execution failed")
	ErrScan          = errors.New("row scanning failed")
	ErrInvalidId     = errors.New("invalid id")
	ErrUserExists    = errors.New("user already exist")
	ErrRowsAffected  = errors.New("unable to get affected row")
	ErrNoAffectedRow = errors.New("rows affected is 0")
	ErrLastInsertId  = errors.New("unable to get last insert id")
	ErrNoUpdate      = errors.New("data already exists")
	ErrBindJSON      = errors.New("unable to bind json")
	ErrParam         = errors.New("error or missing parameter")
	ErrCredential    = errors.New("password or email doesn't match")
	ErrGeneratedPwd  = errors.New("error generating password hash")
)

func ParseError(err error, ctx echo.Context) error {
	Logging(ctx).Error(err)
	status := http.StatusOK
	message := ""
	switch {
	case errors.Is(err, ErrQuery):
		fallthrough
	case errors.Is(err, ErrGeneratedPwd):
		fallthrough
	case errors.Is(err, ErrScan):
		fallthrough
	case errors.Is(err, ErrRowsAffected):
		fallthrough
	case errors.Is(err, ErrLastInsertId):
		fallthrough
	case errors.Is(err, ErrNoAffectedRow):
		status = http.StatusInternalServerError
		message = "Internal Server Error"
	case errors.Is(err, ErrNoUser):
		status = http.StatusNotFound
		message = "No User found"
	case errors.Is(err, ErrNoData):
		status = http.StatusNotFound
		message = "No data found"
	case errors.Is(err, ErrParam):
		status = http.StatusBadRequest
		message = "error or missing param"
	case errors.Is(err, ErrBindJSON):
		status = http.StatusBadRequest
		message = "Bad request"
	case errors.Is(err, ErrInvalidId):
		status = http.StatusBadRequest
		message = "Invalid ID"
	case errors.Is(err, ErrCredential):
		status = http.StatusBadRequest
		message = "email or password missmatch"
	case errors.Is(err, ErrUserExists):
		status = http.StatusBadRequest
		message = "User Already Exists"
	case errors.Is(err, ErrNoUpdate):
		status = http.StatusBadRequest
		message = "Data is the same"
	default:
		status = http.StatusInternalServerError
		message = "Unknown error:" + err.Error()
	}

	return ctx.JSON(status, map[string]interface{}{"message": message})
}
