package models

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

/*
Структуры для валидации запросов с json

В доках echo:
https://echo.labstack.com/docs/request#validate-data
*/

type PlaygroundValidator struct {
	Validator *validator.Validate
}

func (pv *PlaygroundValidator) Validate(i interface{}) error {
	if err := pv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	return nil
}
