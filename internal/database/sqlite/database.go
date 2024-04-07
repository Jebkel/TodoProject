package sqlite

import (
	"ToDoProject/internal/models"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
)

var (
	dbFile = os.Getenv("DB_FILE")
)

// Connect : Database connect
func Connect() (*gorm.DB, *sql.DB) {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
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
