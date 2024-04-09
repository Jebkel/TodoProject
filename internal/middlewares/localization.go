package middlewares

import (
	"ToDoProject/internal/models"
	"ToDoProject/tools"
	"github.com/eduardolat/goeasyi18n"
	"github.com/labstack/echo/v4"
)

func Localization(i18n *goeasyi18n.I18n) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			dbUser := c.Get("db_user")
			language := "en"

			if user, ok := dbUser.(*models.User); ok {

				language = user.Language
			} else if acceptLanguageHeader := c.Request().Header.Get("Accept-Language"); acceptLanguageHeader != "" {
				language = tools.ParseAcceptLanguage(acceptLanguageHeader)[0].Lang
			}
			c.Set("i18n", i18n)
			c.Set("lang", language)
			return next(c)
		}
	}
}
