package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
	validate.RegisterValidation("latitude", ValidateLatitude)
	validate.RegisterValidation("longitude", ValidateLongitude)
	validate.RegisterValidation("status", ValidateStatus)
}

func ValidateStruct(s any) error {
	return validate.Struct(s)
}
