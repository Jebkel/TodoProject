package user

import (
	"ToDoProject/internal/models"
	"ToDoProject/tools"
	"github.com/eduardolat/goeasyi18n"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func (RouterUser) Me(c echo.Context) error {
	user := c.Get("db_user").(*models.User)

	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
}

func (RouterUser) RefreshJWTToken(c echo.Context) error {
	// Получение данных из echo контекста
	db := tools.GetDBFromContext(c)
	jwtClaims, _ := c.Get("jwt_claims").(*models.JwtCustomClaims)
	user, _ := c.Get("db_user").(*models.User)
	// Обновление JWT токенов в бд, как не рабочие
	db.Model(&models.UserToken{}).Where("id = ?", jwtClaims.LinkedTokenID).Update("is_disabled",
		true)
	db.Model(&models.UserToken{}).Where("id = ?", jwtClaims.ID).Update("is_disabled", true)

	// Генерация новых JWT токенов
	accessToken, refreshToken, err := user.GenerateJwt(db)
	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// Возврат ответа с JWT токенами и информацией о пользователе
	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func (RouterUser) UpdatePassword(c echo.Context) error {
	// Структура для тела запроса
	type RequestBody struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required"`
	}

	var body RequestBody

	// Привязка данных запроса к структуре
	if err := c.Bind(&body); err != nil {
		return err
	}

	// Валидация данных запроса
	if err := c.Validate(&body); err != nil {
		return err
	}

	// Получение модели пользователя из контекста echo
	user := c.Get("db_user").(*models.User)

	// Получение коннекта с бд из контекста echo
	db := tools.GetDBFromContext(c)

	checkPassword, err := user.ValidatePassword(body.OldPassword)
	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	if !checkPassword {
		i18n := c.Get("i18n").(*goeasyi18n.I18n)
		language := c.Get("lang").(string)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": i18n.T(language, "bad_old_password", goeasyi18n.Options{}),
		})
	}
	user.Password = body.NewPassword
	if user.HashPassword() != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	db.Save(&user)
	return c.NoContent(http.StatusNoContent)
}
