package middlewares

import (
	"net/http"

	"REDACTED/team-11/backend/booking/pkg/logger"

	"go.uber.org/zap"
)

func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.FromCtx(r.Context()).Error("recovered from error", zap.Error(err.(error)))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
