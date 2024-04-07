package models

import (
	"ToDoProject/tools"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"os"
	"time"
)

type User struct {
	gorm.Model
	ID uint64 `gorm:"primary_key"`

	Username    string
	Password    string
	DisplayName string

	UserTokens []UserToken
}

type JwtCustomClaims struct {
	UserID         uint64    `json:"user_id"`
	TokenType      TokenType `json:"token_type"`
	AccessTokenID  uint64    `json:"access_token_id,omitempty"`
	RefreshTokenID uint64    `json:"refresh_token_id,omitempty"`
	jwt.RegisteredClaims
}

var (
	jwtKey        = os.Getenv("JWT_KEY")
	jwtRefreshKey = os.Getenv("JWT_REFRESH_KEY")
	argonP        = &tools.ArgonParams{
		Memory:      64 * 1024,
		Iterations:  2,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
)

// HashPassword : Hash Password
func (u *User) HashPassword() error {
	encodedHash, err := tools.GenerateFromPassword(u.Password, argonP)
	if err != nil {
		return err
	}
	u.Password = encodedHash
	return nil
}

func (u *User) ValidatePassword(password string) (bool, error) {
	match, err := tools.ComparePasswordAndHash(password, u.Password)
	if err != nil {
		return false, err
	}
	return match, nil
}

// GenerateJwt : Generate JWT
func (u *User) GenerateJwt(db *gorm.DB) (accessTokenSigned string, refreshTokenSigned string, err error) {
	// Создание моделей токенов
	jwtAccessModel := UserToken{TokenType: JWTAccess, User: u}
	jwtRefreshModel := UserToken{TokenType: JWTRefresh, User: u}
	db.Create(&jwtAccessModel)
	db.Create(&jwtRefreshModel)

	// Получение длительности жизни токенов из переменных окружения
	accessDuration := tools.GetDurationEnv("JWT_LIFETIME", time.Hour*3)
	refreshDuration := tools.GetDurationEnv("JWT_REFRESH_LIFETIME", time.Hour*24*7)

	// Создание утверждений для JWT токенов
	jwtAccessClaims := tools.CreateJwtClaims(u.ID, JWTAccess, jwtAccessModel.ID, accessDuration)
	jwtRefreshClaims := tools.CreateJwtClaims(u.ID, JWTRefresh, jwtRefreshModel.ID, refreshDuration)

	// Подпись и получение JWT токенов
	accessTokenSigned, err = tools.SignJwt(jwtAccessClaims, jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshTokenSigned, err = tools.SignJwt(jwtRefreshClaims, jwtRefreshKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenSigned, refreshTokenSigned, nil
}
