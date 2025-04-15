package middleware

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	resp "REDACTED/team-11/backend/admin/internal/controller/response"
	types "REDACTED/team-11/backend/admin/internal/entity/type"
)

func (m *Middleware) CheckAccess(roles ...types.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		if header == "" {
			resp.AbortErrMsg(ctx, foundErr)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) < 2 {
			resp.AbortErrMsg(ctx, bearerErr)
			return
		}

		bearer := parts[0]
		token := parts[1]

		if bearer != bearerType {
			resp.AbortErrMsg(ctx, bearerErr)
			return
		}

		claims, err := m.auth.ValidateToken(token, false)
		if err != nil {
			resp.AbortErrMsg(ctx, unauthErr.WithErr(err))
			return
		}

		if len(roles) != 0 {
			access := false

			if slices.Contains(roles, claims.Role) {
				access = true
			}

			if !access {
				resp.AbortErrMsg(ctx, forbiddenErr)
				return
			}
		}

		ctx.Set("userId", claims.Id)

		ctx.Next()
	}
}
