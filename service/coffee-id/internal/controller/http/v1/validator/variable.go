package validator

import (
	types "github.com/nikitaSstepanov/coffee-id/internal/entity/type"
	e "github.com/nikitaSstepanov/tools/error"
)

const (
	Password Arg = iota
	Birthday
	Fields
)

type Arg int

var (
	validFields = []types.Field{
		types.ID, types.EMAIL,
		types.NAME, types.VERIFIED,
		types.BIRTHDAY, types.ROLES,
	}
)

var (
	lenErr   = e.New("Bad string length", e.BadInput)
	uuidErr  = e.New("Id must be uuid", e.BadInput)
	emailErr = e.New("Invlaid email", e.BadInput)
)

type email struct {
	Value string `validate:"email"`
}

type uuid struct {
	Value string `validate:"uuid"`
}
