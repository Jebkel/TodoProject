package auth

import (
	"ToDoProject/internal/models"
	"ToDoProject/tools"
	"github.com/eduardolat/goeasyi18n"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

// Register : Register RouterAuth
func (RouterAuth) Register(c echo.Context) error {
	// Структура для тела запроса
	type RequestBody struct {
		Username    string `json:"username" validate:"required,gte=6,lte=32"`
		Password    string `json:"password" validate:"required,gte=8,lte=255"`
		DisplayName string `json:"display_name" validate:"required,gte=6,lte=32"`
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

	// Получение подключения к базе данных
	db := tools.GetDBFromContext(c)

	// Проверка наличия пользователя с таким же именем
	var user models.User

	if db.Where("username = ?", body.Username).First(&user).Error == nil {
		return c.NoContent(http.StatusConflict)
	}

	// Создание нового пользователя
	user = models.User{
		Username:    body.Username,
		Password:    body.Password,
		DisplayName: body.DisplayName,
	}

	// Хеширование пароля пользователя
	if err := user.HashPassword(); err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// Сохранение пользователя в базе данных
	if err := db.Create(&user).Error; err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// Генерация JWT токенов
	accessToken, refreshToken, err := user.GenerateJwt(db)
	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// Возврат ответа с JWT токенами и информацией о пользователе
	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  &accessToken,
		"refresh_token": &refreshToken,
		"user":          &user,
	})
}

// Login : Login RouterAuth
func (RouterAuth) Login(c echo.Context) error {
	// Структура для тела запроса
	type RequestBody struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
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

	// Получение подключения к базе данных
	db := tools.GetDBFromContext(c)

	// Проверка наличия пользователя с таким же именем
	var user models.User

	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// Проверка пароля
	check, err := user.ValidatePassword(body.Password)
	if err != nil || !check {
		i18n := c.Get("i18n").(*goeasyi18n.I18n)
		language := c.Get("lang").(string)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": i18n.T(language, "invalid_login_or_password", goeasyi18n.Options{}),
		})
	}

	// Генерация JWT токенов
	accessToken, refreshToken, err := user.GenerateJwt(db)
	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// Возврат ответа с JWT токенами и информацией о пользователе
	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  &accessToken,
		"refresh_token": &refreshToken,
		"user":          &user,
	})
}

func (RouterAuth) Logout(c echo.Context) error {
	// Получение коннекта с бд из контекста echo
	db := tools.GetDBFromContext(c)

	// Получение JWT токена из контекста и обновление их в бд, как не работающие
	jwtClaims, _ := c.Get("jwt_claims").(*models.JwtCustomClaims)
	db.Model(&models.UserToken{}).Where("id = ?", jwtClaims.LinkedTokenID).Update("is_disabled",
		true)
	db.Model(&models.UserToken{}).Where(" = ?", jwtClaims.ID).Update("is_disabled", true)

	return c.NoContent(http.StatusNoContent)
}
