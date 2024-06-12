package helper

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// check authorization
		// do the auth here
		tokenString := ctx.Request().Header.Get("Authorization")
		response := map[string]interface{}{}
		if tokenString == "" {
			Logging(ctx).Warning("unable to get the token")
			response["message"] = "unauthorized"
			return ctx.JSON(http.StatusUnauthorized, response)
		}

		_ = godotenv.Load()
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(os.Getenv("KEY")), nil
		})

		if err != nil || !token.Valid {
			Logging(ctx).Warning("token invalid")
			response["message"] = "unauthorized"
			return ctx.JSON(http.StatusUnauthorized, response)
		}

		// change token -> struct
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx.Set("email", claims["email"])
			ctx.Set("id", claims["id"])
			ctx.Set("role", claims["role"])
		} else {
			Logging(ctx).Warning("invalid claims")
			response["message"] = "invalid claims"
			return ctx.JSON(http.StatusUnauthorized, response)
		}

		return next(ctx)
	}
}
