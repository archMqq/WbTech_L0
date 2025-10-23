package validation

import (
	"log"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
)

type OrderValidator struct {
	validator *validator.Validate
}

func NewValidator() *OrderValidator {
	v := validator.New()

	if err := v.RegisterValidation("real_date", validateRealDate); err != nil {
		log.Printf("error adding 'real_date' validation: %s", err)
	}

	return &OrderValidator{v}
}

func (ov *OrderValidator) Validate(i interface{}) error {
	return ov.validator.Struct(i)
}

func validateRealDate(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Type() != reflect.TypeOf(time.Time{}) {
		return false
	}

	date := field.Interface().(time.Time)
	return date.Before(time.Now())
}
