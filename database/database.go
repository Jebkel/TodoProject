package database

import (
	"ToDoProject/database/models"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

// Connect : Database connect
func Connect() (*gorm.DB, *sql.DB) {
	db, err := gorm.Open(sqlite.Open("./database/data.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	models.Migrate(db)

	return db, sqlDB
}
