package models

import "gorm.io/gorm"

type TokenType string

const (
	JWTAccess  TokenType = "jwt_access"  // Используется, для получения доступа к защищённым данным
	JWTRefresh TokenType = "jwt_refresh" // Используется, для обновление JWTAccess, доступа к данным не имеет
)

type UserToken struct {
	gorm.Model
	ID         uint64    `gorm:"primary_key"`
	TokenType  TokenType `gorm:"token_type"`
	IsDisabled bool      `gorm:"is_disabled,default:0"`
	UserID     uint64
	User       *User `gorm:"foreignKey:UserID"`
}
