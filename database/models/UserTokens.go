package models

import "gorm.io/gorm"

type TokenType string

const (
	JWTAccess     TokenType = "jwt_access"
	JWTRefresh    TokenType = "jwt_refresh"
	PasswordReset TokenType = "password_reset"
)

type UserToken struct {
	gorm.Model
	ID         uint64    `gorm:"primary_key"`
	TokenType  TokenType `gorm:"token_type"`
	IsDisabled bool      `gorm:"is_disabled,default:0"`
	User       *User     `gorm:"foreignKey:id"`
}
