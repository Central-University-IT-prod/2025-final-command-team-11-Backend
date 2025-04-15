package validator

import (
	"slices"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	types "github.com/nikitaSstepanov/coffee-id/internal/entity/type"
)

func setupArgs(validate *validator.Validate, args []Arg) error {
	for i := range len(args) {
		switch args[i] {
		case Password:
			err := validate.RegisterValidation("password", validatePassword)
			if err != nil {
				return err
			}
		case Birthday:
			err := validate.RegisterValidation("age", validateAge)
			if err != nil {
				return err
			}
		case Fields:
			err := validate.RegisterValidation("fields", validateFields)
			if err != nil {
				return err
			}
		default:
			return nil
		}
	}

	return nil
}

func validatePassword(fl validator.FieldLevel) bool {
	pass := fl.Field().String()

	var number, upper, lower, special bool

	for _, s := range pass {
		switch {
		case unicode.IsNumber(s):
			number = true
		case unicode.IsUpper(s):
			upper = true
		case unicode.IsLower(s):
			lower = true
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			special = true
		default:
			return false
		}
	}

	return number && upper && lower && special && len(pass) >= 8
}

func validateAge(fl validator.FieldLevel) bool {
	birthday := fl.Field().Interface().(types.Date)
	age := getAge(time.Time(birthday))

	return age >= 8
}

func validateFields(fl validator.FieldLevel) bool {
	fields := fl.Field().Interface().([]types.Field)

	for _, field := range fields {
		if !slices.Contains(validFields, field) {
			return false
		}
	}

	return true
}

func getAge(birthdate time.Time) int {
	now := time.Now()

	years := now.Year() - birthdate.Year()
	months := int(now.Month()) - int(birthdate.Month())
	days := now.Day() - birthdate.Day()

	if days < 0 {
		months--
	}

	if months < 0 {
		years--
	}

	return years
}
