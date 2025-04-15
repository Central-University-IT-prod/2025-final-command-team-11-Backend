package validator

import (
	e "github.com/nikitaSstepanov/tools/error"
	types "REDACTED/team-11/backend/admin/internal/entity/type"
)

const (
	Booking Arg = iota
)

type Arg int

var (
	lenErr   = e.New("Bad string length", e.BadInput)
	uuidErr  = e.New("Id must be uuid", e.BadInput)
	emailErr = e.New("Invlaid email", e.BadInput)
)

var (
	validBooking = []types.BookingType{
		types.OPENSPACE,
		types.ROOM,
	}
)

type email struct {
	Value string `validate:"email"`
}

type uuid struct {
	Value string `validate:"uuid"`
}
