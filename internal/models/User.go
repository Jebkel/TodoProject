package models

import (
	"ToDoProject/tools"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

type User struct {
	gorm.Model
	ID uint64 `json:"id" gorm:"primary_key"`

	Username    string `json:"username" gorm:"uniqueIndex"`
	Password    string `json:"-"`
	DisplayName string `json:"display_name"`
	Language    string `json:"language" gorm:"default:'en'"`

	UserTokens []UserToken `json:"-"`
}

type JwtCustomClaims struct {
	UserID        uint64    `json:"user_id"`
	TokenType     TokenType `json:"token_type"`
	LinkedTokenID uint64    `json:"linked_token_id,omitempty"`
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
	log.Info(u.Password)
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
	jwtAccessClaims := createJwtClaims(u.ID, JWTAccess, jwtAccessModel.ID, accessDuration, jwtRefreshModel.ID)
	jwtRefreshClaims := createJwtClaims(u.ID, JWTRefresh, jwtRefreshModel.ID, refreshDuration, jwtAccessModel.ID)

	// Подпись и получение JWT токенов
	accessTokenSigned, err = signJwt(jwtAccessClaims, jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshTokenSigned, err = signJwt(jwtRefreshClaims, jwtRefreshKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenSigned, refreshTokenSigned, nil
}

// CreateJwtClaims : Creating a JwtCustomClaims
func createJwtClaims(userID uint64, tokenType TokenType, tokenID uint64, expiresIn time.Duration, linkedTokenId uint64) *JwtCustomClaims {
	return &JwtCustomClaims{
		UserID:        userID,
		TokenType:     tokenType,
		LinkedTokenID: linkedTokenId,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strconv.FormatUint(tokenID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		},
	}
}

// SignJwt : Signing JWT
func signJwt(claims *JwtCustomClaims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}
