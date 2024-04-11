package main

import (
	"ToDoProject/internal/database"
	"ToDoProject/internal/mail"
	"ToDoProject/internal/middlewares"
	"ToDoProject/internal/routes"
	"ToDoProject/locales"
	"ToDoProject/tools"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"os"
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

	i18n := locales.Init()

	e.HTTPErrorHandler = middlewares.NewHttpErrorHandler(tools.NewErrorStatusCodeMaps()).Handler
	e.Validator = &tools.CustomValidator{Validator: validator.New()}

	mailer := &mail.Mailer{}
	mailer.Init()
	go mailer.StartHandling()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(middlewares.ContextMail(mailer))
	e.Use(middlewares.ContextDB(db))
	e.Use(middlewares.ContextJWT())
	e.Use(middlewares.Localization(i18n))

	routes.Routes(e.Group(""))

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))))
}
