package middlewares

import (
	"ToDoProject/internal/mail"
	"github.com/labstack/echo/v4"
)

func ContextMail(mailer *mail.Mailer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("mailer", mailer)
			return next(c)
		}
	}
}
