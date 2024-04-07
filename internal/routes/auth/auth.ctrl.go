package auth

import (
	models2 "ToDoProject/internal/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

// Register : Register RouterAuth
func (RouterAuth) Register(c echo.Context) error {
	type RequestBody struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`

		DisplayName string `json:"display_name" validate:"required"`
	}

	var body RequestBody

	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	db, _ := c.Get("db").(*gorm.DB)

	if err := db.Where("username = ?", body.Username).First(&models2.User{}).Error; err == nil {
		fmt.Println(err)
		return c.NoContent(http.StatusConflict)
	}

	user := models2.User{
		Username: body.Username,
		Password: body.Password,

		DisplayName: body.DisplayName,
	}

	err := user.HashPassword()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	db.Create(&user)

	accessToken, refreshToken, err := user.GenerateJwt(db)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

// Login : Login RouterAuth
func (RouterAuth) Login(c echo.Context) error {
	type RequestBody struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var body RequestBody
	if err := c.Bind(&body); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := c.Validate(&body); err != nil {
		return err
	}

	db, _ := c.Get("db").(*gorm.DB)

	var user models2.User

	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	check, err := user.ValidatePassword(body.Password)
	if err != nil || !check {
		return c.NoContent(http.StatusBadRequest)
	}
	accessToken, refreshToken, err := user.GenerateJwt(db)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func (RouterAuth) Logout(c echo.Context) error {
	db, _ := c.Get("db").(*gorm.DB)
	jwtClaims, _ := c.Get("jwt_claims").(*models2.JwtCustomClaims)
	db.Model(&models2.UserToken{}).Where("id = ?", jwtClaims.RefreshTokenID).Update("is_disabled",
		true)
	db.Model(&models2.UserToken{}).Where(" = ?", jwtClaims.ID).Update("is_disabled", true)

	return c.NoContent(http.StatusNoContent)
}

func (RouterAuth) Me(c echo.Context) error {
	user := c.Get("db_user").(*models2.User)

	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
}

func (RouterAuth) RefreshJWTToken(c echo.Context) error {
	db, _ := c.Get("db").(*gorm.DB)
	jwtClaims, _ := c.Get("jwt_claims").(*models2.JwtCustomClaims)
	user := c.Get("db_user").(*models2.User)

	db.Model(&models2.UserToken{}).Where("id = ?", jwtClaims.AccessTokenID).Update("is_disabled",
		true)
	db.Model(&models2.UserToken{}).Where("id = ?", jwtClaims.ID).Update("is_disabled", true)

	accessToken, refreshToken, err := user.GenerateJwt(db)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})

}
