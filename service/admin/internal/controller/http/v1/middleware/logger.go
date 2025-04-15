package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nikitaSstepanov/tools/ctx"
	"github.com/nikitaSstepanov/tools/sl"
	ct "REDACTED/team-11/backend/admin/pkg/utils/controller"
)

func (m *Middleware) InitLogger(c ctx.Context) gin.HandlerFunc {
	log := c.Logger()

	log.Info("logger middleware enabled.")

	return func(c *gin.Context) {
		ctx := ctx.New(log)

		c.Set(ct.CtxKey, ctx)
		
		req := c.Request

		c.Next()
		entry := log.With(
			sl.StringAttr("method", req.Method),
			sl.StringAttr("path", req.URL.Path),
			sl.StringAttr("remote_addr", req.RemoteAddr),
			sl.StringAttr("user_agent", req.UserAgent()),
		)

		t1 := time.Now()
		defer func() {
			entry.Info("request completed",
				sl.IntAttr("status", c.Writer.Status()),
				sl.StringAttr("duration", time.Since(t1).String()),
			)
		}()
	}
}
