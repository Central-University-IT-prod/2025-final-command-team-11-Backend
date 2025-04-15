package account

import (
	resp "github.com/nikitaSstepanov/coffee-id/internal/controller/response"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/httper"
)

const (
	okStatus        = httper.StatusOK
	createdStatus   = httper.StatusCreated
	noContentStatus = httper.StatusNoContent
)

var (
	badReqErr = e.New("Incorrect data.", e.BadInput)
)

var (
	okMsg = resp.NewMessage("Ok.")
)
