package resp

import (
	"github.com/gin-gonic/gin"
	ct "github.com/nikitaSstepanov/coffee-id/pkg/utils/controller"
	e "github.com/nikitaSstepanov/tools/error"
)

type Message struct {
	Message string `json:"message"`
}

func NewMessage(msg string) *Message {
	return &Message{
		Message: msg,
	}
}

// JsonError use only for doc and represent e.JsonError
type JsonError struct {
	Error string `json:"error"`
}

func AbortErrMsg(c *gin.Context, err e.Error) {
	ctx := ct.GetCtx(c)
	log := ctx.Logger()

	if err.GetCode() == e.Internal {
		log.Error("Something going wrong...", err.SlErr())
	} else {
		log.Info("Invalid input data")
	}

	c.AbortWithStatusJSON(
		err.ToHttpCode(),
		err.ToJson(),
	)
}
