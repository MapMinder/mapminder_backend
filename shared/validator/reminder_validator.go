package validator

import (
	"github.com/go-playground/validator/v10"
)

func ValidateLatitude(fl validator.FieldLevel) bool {
	param := fl.Field().Float()
	return param >= -90 && param <= 90
}

func ValidateLongitude(fl validator.FieldLevel) bool {
	param := fl.Field().Float()
	return param >= -180 && param <= 180
}
