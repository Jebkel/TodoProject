package middlewares

import (
	"ToDoProject/internal/models"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
)

var (
	jwtKey        = os.Getenv("JWT_KEY")
	jwtRefreshKey = os.Getenv("JWT_REFRESH_KEY")
)

// Authorized : Check Auth
func Authorized(tokenType models.TokenType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenHeader := c.Request().Header.Get("Authorization")
			if tokenHeader == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			claims, err := parseJWTClaims(tokenHeader, tokenType)
			if err != nil || claims.TokenType != tokenType {
				log.Error(err)
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
			}

			db := c.Get("db").(*gorm.DB)
			var userTokenModel models.UserToken

			db.First(&userTokenModel, claims.ID)
			if userTokenModel.IsDisabled {
				return c.NoContent(http.StatusUnauthorized)
			}

			var user models.User
			result := db.First(&user, claims.UserID)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "User not found"})
			}
			c.Set("jwt_claims", claims)
			c.Set("db_user", &user)
			return next(c)
		}
	}
}

func parseJWTClaims(tokenHeader string, tokenType models.TokenType) (*models.JwtCustomClaims, error) {
	tokenString := strings.Replace(tokenHeader, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка на метод, используемый для подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// определения ключа подписи в зависимости от типа токена
		switch tokenType {
		case models.JWTAccess:
			return []byte(jwtKey), nil
		case models.JWTRefresh:
			return []byte(jwtRefreshKey), nil
		default:
			return nil, fmt.Errorf("unexcepted token type: %v", tokenType)
		}
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*models.JwtCustomClaims)
	if !ok || claims == nil {
		log.Error("failed to parse JWT claims")
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
