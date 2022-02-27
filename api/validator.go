package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/hamdysherif/simplebank/util"
)

var validCurrency validator.Func = func(field validator.FieldLevel) bool {
	if currency, ok := field.Field().Interface().(string); ok {
		return util.ValidCurrency(currency)
	}

	return false
}
