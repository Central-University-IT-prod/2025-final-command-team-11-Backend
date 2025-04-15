package auth

import (
	"time"

	e "github.com/nikitaSstepanov/tools/error"
)

const (
	refreshExpires = 72 * time.Hour
	accessExpires  = 72 * time.Hour
)

var (
	badDataErr  = e.New("Incorrect email or password", e.Unauthorize)
	unauthErr   = e.New("Token is invalid", e.Unauthorize)
	notFoundErr = e.New("This Coffee ID user wasn`t found.", e.NotFound)
)
