package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var GuestUserId = "guest_user_id"

func CreateSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := c.Cookie(GuestUserId)
		if err != nil {
			cookie := new(http.Cookie)
			cookie.Name = GuestUserId
			cookie.Value = uuid.New().String()
			c.SetCookie(cookie)
		}
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
