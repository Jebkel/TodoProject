package tools

import (
	"ToDoProject/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

// GetDurationEnv : Get durations time from environment
func GetDurationEnv(envVar string, defaultValue time.Duration) time.Duration {
	durationStr := os.Getenv(envVar)
	if durationStr == "" {
		log.Errorf("Error parsing %s, using default value: %v", envVar, defaultValue)
		return defaultValue
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Errorf("Error parsing %s, using default value: %v", envVar, defaultValue)
		return defaultValue
	}

	return duration
}

// CreateJwtClaims : Creating a JwtCustomClaims
func CreateJwtClaims(userID uint64, tokenType models.TokenType, tokenID uint64, expiresIn time.Duration) *models.JwtCustomClaims {
	return &models.JwtCustomClaims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strconv.FormatUint(uint64(tokenID), 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		},
	}
}

// SignJwt : Signing JWT
func SignJwt(claims *models.JwtCustomClaims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

// GetDBFromContext : Getting DB Connection from echo context
func GetDBFromContext(c echo.Context) *gorm.DB {
	db, _ := c.Get("db").(*gorm.DB)
	return db
}

// GetJWTFromContext : Getting JWTCustomClaims from echo context
func GetJWTFromContext(c echo.Context) *models.JwtCustomClaims {
	jwtClaims, _ := c.Get("jwt_claims").(*models.JwtCustomClaims)
	return jwtClaims
}

// GetUserModelFromContext : Getting User Model from echo context
func GetUserModelFromContext(c echo.Context) *models.User {
	user, _ := c.Get("db_user").(*models.User)
	return user
}
