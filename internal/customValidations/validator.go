package customValidations

import "github.com/go-playground/validator/v10"

func Init() *validator.Validate {
	v := validator.New()
	err := v.RegisterValidation("datetime", customDateTime)
	if err != nil {
		panic(err)
	}
	return v
}
