package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	et "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(f reflect.StructField) (name string) {
		name_ := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name_ == "-" {
			return
		}
		return name_
	})
	var en = en.New()
	var ut = ut.New(en, en)
	trans, _ = ut.GetTranslator("en")
	_ = et.RegisterDefaultTranslations(validate, trans)
}

type ValidationError map[string]string

func translateError(err error, tr ...ut.Translator) (ve ValidationError) {
	if err == nil {
		return
	}
	tr0 := trans
	if len(tr) > 0 {
		tr0 = tr[0]
	}
	validatorErrs := err.(validator.ValidationErrors)
	ve = ValidationError{}
	for _, e := range validatorErrs {
		ve[e.Field()] = e.Translate(tr0)
	}
	return
}

func ValidateStruct(s interface{}) (ve ValidationError) {
	err := validate.Struct(s)
	if err == nil {
		return
	}
	ve = translateError(err)
	return
}
