package middlewares

import (
	"ToDoProject/internal/models"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"net/http"
	"os"
)

var (
	jwtKey        = os.Getenv("JWT_KEY")
	jwtRefreshKey = os.Getenv("JWT_REFRESH_KEY")
)

// ContextJWT : Add data from JWT token to context
func ContextJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenHeader := c.Request().Header.Get("Authorization")
			if tokenHeader == "" {
				c.Set("jwt_status", map[string]interface{}{
					"status":  http.StatusUnauthorized,
					"message": "Authorization is required",
				})
				return next(c)
			}

			claims, err := parseJWTClaims(tokenHeader)
			if err != nil {
				log.Error(err)
				c.Set("jwt_status", map[string]interface{}{
					"status":  http.StatusUnauthorized,
					"message": "Invalid token",
				})
				return next(c)
			}

			db := c.Get("db").(*gorm.DB)
			var userTokenModel models.UserToken

			db.First(&userTokenModel, claims.ID)
			if userTokenModel.IsDisabled {
				c.Set("jwt_status", map[string]interface{}{
					"status":  http.StatusUnauthorized,
					"message": "token expired",
				})
				return next(c)
			}

			var user models.User
			result := db.First(&user, claims.UserID)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.Set("jwt_status", map[string]interface{}{
					"status":  http.StatusUnauthorized,
					"message": "User not found",
				})
				return next(c)
			}
			c.Set("jwt_claims", claims)
			c.Set("db_user", &user)

			return next(c)
		}
	}
}
