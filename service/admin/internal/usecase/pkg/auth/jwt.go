package auth

import (
	"github.com/golang-jwt/jwt/v5"
	e "github.com/nikitaSstepanov/tools/error"
)

type Jwt struct {
	Audience  []string
	issuer    string
	accessKey string
}

func NewJwt(options *JwtOptions) *Jwt {
	return &Jwt{
		Audience:  options.Audience,
		issuer:    options.Issuer,
		accessKey: options.AccessKey,
	}
}

func (j *Jwt) ValidateToken(jwtString string, isRefresh bool) (*Claims, e.Error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(j.accessKey), nil
	}

	token, err := jwt.ParseWithClaims(jwtString, &Claims{}, keyFunc)
	if err != nil {
		return nil, unauthErr.WithErr(err)
	}

	if !token.Valid {
		return nil, unauthErr.WithErr(err)
	}

	return token.Claims.(*Claims), nil
}
