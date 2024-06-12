package helper

import "github.com/labstack/echo/v4"

type Credential struct {
	UserID uint
	Email  string
	Role   string
}

func GetCredential(c echo.Context) Credential {
	ctx_id := c.Get("id").(float64)
	ctx_email := c.Get("email").(string)
	ctx_role := c.Get("role").(string)
	var cred Credential
	cred.UserID = uint(ctx_id)
	cred.Email = ctx_email
	cred.Role = ctx_role

	return cred
}
