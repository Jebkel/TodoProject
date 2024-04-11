package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type PasswordResetCode struct {
	ID        uint64 `gorm:"primary_key"`
	UserID    uint64 `gorm:"index"`
	Token     string `gorm:"uniqueIndex"`
	User      *User  `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *PasswordResetCode) HashToken() {
	hash := sha256.New()
	hash.Write([]byte(p.Token))
	p.Token = hex.EncodeToString(hash.Sum(nil))
}
