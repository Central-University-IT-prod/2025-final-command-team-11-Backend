package account

import (
	types "github.com/nikitaSstepanov/coffee-id/internal/entity/type"
	e "github.com/nikitaSstepanov/tools/error"
)

var (
	conflictErr = e.New("User with this email already exist", e.Conflict)
	badPassErr  = e.New("Incorrect password", e.Forbidden)

	defaultRole = types.USER
)
