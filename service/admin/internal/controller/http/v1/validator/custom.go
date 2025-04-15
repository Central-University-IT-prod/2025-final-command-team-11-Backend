package validator

import (
	"slices"

	"github.com/go-playground/validator/v10"
	types "REDACTED/team-11/backend/admin/internal/entity/type"
)

func setupArgs(validate *validator.Validate, args []Arg) error {
	for i := range len(args) {
		switch args[i] {
		case Booking:
			err := validate.RegisterValidation("booking", validateBookingType)
			if err != nil {
				return err
			}
		default:
			return nil
		}
	}

	return nil
}

func validateBookingType(fl validator.FieldLevel) bool {
	booking := fl.Field().Interface().(types.BookingType)

	return slices.Contains(validBooking, booking)
}
