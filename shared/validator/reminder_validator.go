package validator

import (
	"github.com/MapMinder/mapminder_backend/feature/reminder/domain"
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

func ValidateStatus(fl validator.FieldLevel) bool {
	param := fl.Field().String()
	if param == "" {
		return true
	}
	return domain.IsValidStatus(param)
}
