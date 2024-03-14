package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type (
	Field       = validator.FieldLevel
	ValidateMsg = map[string]string
)

func msg(validate validator.ValidationErrors) ValidateMsg {
	var msg ValidateMsg

	for _, i := range validate {
		// low casting first letter
		key := strings.ToLower(i.Field()[:1]) + i.Field()[1:]

		if msg == nil {
			msg = ValidateMsg{}
		}

		// push msg
		switch i.Tag() {
		case "required":
			msg[key] = "This field cannot empty"
		case "min":
			msg[key] =
				fmt.Sprintf(
					"This field need value minimum %s",
					i.Param())
		case "max":
			msg[key] =
				fmt.Sprintf(
					"This field need value maximum %s",
					i.Param())
		case "email":
			msg[key] = "This field need valid email"
		case "unique":
			msg[key] = "Email already exists"
		case "eqfield":
			msg[key] = "Password not same"
		case "imageValidate":
			msg[key] = "This fill with file ext .png, .jpg, .img and size file must be less then 2mb"
		case "uniqueProduct":
			msg[key] = "This product already exits"
		}
	}

	return msg
}

func registerValidate(v *validator.Validate) {
	v.RegisterValidation("imageValidate", imageValidate)
	v.RegisterValidation("unique", unique)
	v.RegisterValidation("uniqueProduct", uniqueProduct)
}

func IsValid(value interface{}) ValidateMsg {
	validate := validator.New()
	registerValidate(validate)

	if err := validate.Struct(value); err != nil {
		return msg(err.(validator.ValidationErrors))
	}
	return nil
}
