package main

import (
	"ToDoProject/internal/database"
	"ToDoProject/internal/middlewares"
	"ToDoProject/internal/routes"
	"ToDoProject/tools"
	"database/sql"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	log.Info("Initializing environment")
}

func main() {
	db, sqlDB := database.ConnectDB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic(err)
		}
	}(sqlDB)

	e := echo.New()

	e.Validator = &tools.CustomValidator{Validator: validator.New()}

	e.Use(middlewares.ContextDB(db))

	routes.Routes(e.Group(""))

	e.Use()

	e.Logger.Fatal(e.Start(":8000"))
}
