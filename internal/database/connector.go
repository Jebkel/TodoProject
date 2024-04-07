package database

import (
	"ToDoProject/internal/database/sqlite"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"os"
)

var (
	dbDriver = os.Getenv("DB_DRIVER")
)

func ConnectDB() (*gorm.DB, *sql.DB) {
	switch dbDriver {
	case "sqlite":
		return sqlite.Connect()
	}
	panic(fmt.Sprintf("database driver '%s' not found", dbDriver))
}
