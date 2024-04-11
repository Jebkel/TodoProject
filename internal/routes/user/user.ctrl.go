package user

import (
	"ToDoProject/internal/mail"
	"ToDoProject/internal/mail/structures"
	"ToDoProject/internal/models"
	"ToDoProject/tools"
	"fmt"
	"github.com/eduardolat/goeasyi18n"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type ApiError struct {
	Param   string
	Message string
}

var (
	MAIL_PORT = os.Getenv("MAIL_PORT")
	MAIL_HOST = os.Getenv("MAIL_HOST")
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

func (RouterUser) SendPasswordResetCode(c echo.Context) error {
	type RequestBody struct {
		Username string `json:"username" validate:"required"`
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

	db := tools.GetDBFromContext(c)

	var user models.User

	if db.Where("username = ?", body.Username).First(&user).Error != nil {
		// Отправляем, что всё окей, что бы нельзя было перебирать пользователей
		return c.NoContent(http.StatusOK)
	}

	// Генерируем 8-значный код для подтверждения
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const numbers = "0123456789"
	var resultCode string
	for i := 0; i < 8; i++ {
		resultCode += string(numbers[r.Intn(len(numbers))])
	}

	// Удаление прошлых кодов востановления для пользователя
	db.Where("user_id = ?", user.ID).Delete(&models.PasswordResetCode{})

	// Создание записи в бд
	passwordResetModel := models.PasswordResetCode{
		User:  &user,
		Token: resultCode,
	}
	passwordResetModel.HashToken()
	db.Create(&passwordResetModel)

	// Отправка кода на mail
	mailer := c.Get("mailer").(*mail.Mailer)
	i18n := c.Get("i18n").(*goeasyi18n.I18n)
	language := c.Get("lang").(string)
	mailer.QueueEmail("test@gmail.com", "Востановление пароля", structures.MessagesData{
		PreHeader: i18n.T(language, fmt.Sprintf("password_reeset"), goeasyi18n.Options{}),
		Messages: []string{
			i18n.T(language, "password_reset_msg_1", goeasyi18n.Options{}),
			i18n.T(language, "password_reset_msg_2", goeasyi18n.Options{
				Data: map[string]string{
					"Code": resultCode,
				},
			}),
			i18n.T(language, "password_reset_msg_3", goeasyi18n.Options{}),
		},
	})
	return c.NoContent(http.StatusNoContent)
}
