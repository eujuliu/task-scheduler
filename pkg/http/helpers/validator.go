package helpers

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var Datatime validator.Func = func(fl validator.FieldLevel) bool {
	_, ok := fl.Field().Interface().(time.Time)

	return ok
}

var UTCDateTime validator.Func = func(fl validator.FieldLevel) bool {
	t, ok := fl.Field().Interface().(time.Time)

	return ok && t.Location() == time.UTC
}
