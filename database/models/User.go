package models

import (
	"ToDoProject/utils"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"os"
	"strconv"
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
	argonP        = &utils.ArgonParams{
		Memory:      64 * 1024,
		Iterations:  2,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
)

// HashPassword : Hash Password
func (u *User) HashPassword() error {
	encodedHash, err := utils.GenerateFromPassword(u.Password, argonP)
	if err != nil {
		return err
	}
	u.Password = encodedHash
	return nil
}

func (u *User) ValidatePassword(password string) (bool, error) {
	match, err := utils.ComparePasswordAndHash(password, u.Password)
	if err != nil {
		return false, err
	}
	return match, nil
}

// GenerateJwt : Generate JWT
func (u *User) GenerateJwt(db *gorm.DB) (accessTokenSigned string, refreshTokenSigned string, error error) {

	jwtAccessModel := UserToken{
		TokenType: JWTAccess,
		User:      u,
	}
	jwtRefreshModel := UserToken{
		TokenType: JWTRefresh,
		User:      u,
	}
	db.Create(&jwtAccessModel)
	db.Create(&jwtRefreshModel)

	jwtAccessClaims := &JwtCustomClaims{
		UserID:         u.ID,
		TokenType:      JWTAccess,
		RefreshTokenID: jwtRefreshModel.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strconv.FormatUint(jwtAccessModel.ID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	jwtRefreshClaims := &JwtCustomClaims{
		UserID:        u.ID,
		TokenType:     JWTRefresh,
		AccessTokenID: jwtAccessModel.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strconv.FormatUint(jwtRefreshModel.ID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtAccessClaims)

	accessTokenSigned, error = jwtAccessToken.SignedString([]byte(jwtKey))

	jwtRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtRefreshClaims)
	refreshTokenSigned, error = jwtRefreshToken.SignedString([]byte(jwtRefreshKey))
	return accessTokenSigned, refreshTokenSigned, error
}
