package auth

import e "github.com/nikitaSstepanov/tools/error"

type Auth struct {
	jwt *Jwt
}

func New(opts *JwtOptions) *Auth {
	return &Auth{
		jwt: NewJwt(opts),
	}
}

func (a *Auth) ValidateToken(token string, isRefresh bool) (*Claims, e.Error) {
	return a.jwt.ValidateToken(token, isRefresh)
}
