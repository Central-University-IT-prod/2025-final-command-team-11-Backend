package middleware

import "REDACTED/team-11/backend/admin/internal/usecase/pkg/auth"

type Middleware struct {
	auth *auth.Auth
}

func New(auth *auth.Auth) *Middleware {
	return &Middleware{
		auth,
	}
}
