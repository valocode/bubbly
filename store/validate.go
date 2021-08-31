package store

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(validateTagNameFunc)
	return v
}

func validateTagNameFunc(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	// If the struct field is another struct embedded, then check if there
	// is an alias tag attached so that we can get a nicer name than the
	// horrible struct names, e.g. MyTypeModelCreate `alias:"type"`
	if name == "" && fld.Anonymous {
		return strings.SplitN(fld.Tag.Get("alias"), ",", 2)[0]
	}

	if name == "-" {
		return ""
	}

	return name
}
