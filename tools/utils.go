package tools

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/xhit/go-str2duration/v2"
	"gorm.io/gorm"
	"os"
	"time"
)

// GetDurationEnv : Get durations time from environment
func GetDurationEnv(envVar string, defaultValue time.Duration) time.Duration {
	durationStr := os.Getenv(envVar)
	if durationStr == "" {
		log.Errorf("Error parsing %s, using default value: %v", envVar, defaultValue)
		return defaultValue
	}

	duration, err := str2duration.ParseDuration(durationStr)
	if err != nil {
		log.Errorf("Error parsing %s, using default value: %v", envVar, defaultValue)
		return defaultValue
	}
	return duration
}

// GetDBFromContext : Getting DB Connection from echo context
func GetDBFromContext(c echo.Context) *gorm.DB {
	db, ok := c.Get("db").(*gorm.DB)
	if !ok {
		log.Error("cant extract database from context")
		return nil
	}
	return db
}
