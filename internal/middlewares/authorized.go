package middlewares

import (
	"ToDoProject/internal/models"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
)

// Authorized : Check Auth
func Authorized(tokenType models.TokenType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			jwtStatus, ok := c.Get("jwt_status").(map[string]interface{})
			if ok {
				token := c.Get("jwt_claims").(*models.JwtCustomClaims)
				if token.TokenType != tokenType {
					return c.JSON(http.StatusUnauthorized, echo.Map{"error": "bad token type"})
				}
				if jwtStatus["message"] == nil {
					return c.NoContent(jwtStatus["status"].(int))
				}
				return c.JSON(jwtStatus["status"].(int), echo.Map{"error": jwtStatus["message"]})
			}

			return next(c)
		}
	}
}

func parseJWTClaims(tokenHeader string) (*models.JwtCustomClaims, error) {
	tokenString := strings.Replace(tokenHeader, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка на метод, используемый для подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// определения ключа подписи в зависимости от типа токена
		switch token.Claims.(*models.JwtCustomClaims).TokenType {
		case models.JWTAccess:
			return []byte(jwtKey), nil
		case models.JWTRefresh:
			return []byte(jwtRefreshKey), nil
		default:
			return nil, fmt.Errorf("unexcepted token type")
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
