package passwordRecovery

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
	"time"
)

func (RouterPassRecovery) SendPasswordResetCode(c echo.Context) error {
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

	if result := db.Where("username = ?", body.Username).First(&user); result.Error != nil || result.RowsAffected == 0 {
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

func (RouterPassRecovery) ResetPassword(c echo.Context) error {
	type RequestBody struct {
		ResetCode   string `json:"reset_code" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,gte=8,lte=255"`
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

	language := c.Get("lang").(string)
	i18n := c.Get("i18n").(*goeasyi18n.I18n)

	db := tools.GetDBFromContext(c)
	var passwordResetModel models.PasswordResetCode

	if result := db.Preload("User").Find(&passwordResetModel, "token = ?", body.ResetCode); result.Error != nil || result.RowsAffected == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": map[string]string{
				"reset_code": i18n.T(language, "incorrect_code", goeasyi18n.Options{}),
			},
		})
	}

	check, err := passwordResetModel.User.ValidatePassword(body.NewPassword)
	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	if check {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": map[string]string{
				"new_password": i18n.T(language, "password_already_using", goeasyi18n.Options{}),
			},
		})
	}

	passwordResetModel.User.Password = body.NewPassword
	if err = passwordResetModel.User.HashPassword(); err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	db.Save(&passwordResetModel.User)
	db.Delete(&passwordResetModel)
	return c.NoContent(http.StatusNoContent)
}
