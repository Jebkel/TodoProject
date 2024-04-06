package main

import (
	"ToDoProject/app/middlewares"
	"ToDoProject/database"
	"ToDoProject/routes"
	"ToDoProject/utils"
	"database/sql"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
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
	db, sqlDB := database.Connect()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic(err)
		}
	}(sqlDB)

	e := echo.New()

	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	e.Use(middlewares.ContextDB(db))

	routes.Routes(e.Group(""))

	e.Use()

	e.Logger.Fatal(e.Start(":8000"))
}
