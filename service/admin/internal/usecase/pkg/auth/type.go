package auth

import (
	e "github.com/nikitaSstepanov/tools/error"
)

var (
	unauthErr = e.New("Token is invalid", e.Unauthorize)
)
