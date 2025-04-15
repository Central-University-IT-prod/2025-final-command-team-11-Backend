package middleware

import "github.com/nikitaSstepanov/coffee-id/internal/usecase/pkg/auth"

type Middleware struct {
	auth *auth.Auth
}

func New(auth *auth.Auth) *Middleware {
	return &Middleware{
		auth,
	}
}
