package customValidations

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"time"
)

func customDateTime(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	layout := "02.01.2006 15:04"

	// Парсинг происходит в формате UTC
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Error(err)
		return false
	}

	// Время сервера переводим в UTC для стандартизации
	return t.After(time.Now().UTC())
}
