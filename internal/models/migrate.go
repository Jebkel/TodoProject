package models

import "gorm.io/gorm"

// Migrate : migrate models
func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		// В Go не принято вызывать panic, но в данном случае от этой миграции зависит то, будет ли функционировать
		// наше приложение, по этому вызов panic для меня рационален в этом месте
		panic(err.Error())
	}
}
