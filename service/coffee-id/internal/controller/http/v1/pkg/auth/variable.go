package auth

import (
	resp "github.com/nikitaSstepanov/coffee-id/internal/controller/response"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/httper"
)

const (
	okStatus        = httper.StatusOK
)

var (
	badReqErr = e.New("Incorrect data.", e.BadInput)
)

var (
	logoutMsg = resp.NewMessage("Logout success.")
)
