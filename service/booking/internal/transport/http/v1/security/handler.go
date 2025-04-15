package security

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"REDACTED/team-11/backend/booking/internal/models"
	api "REDACTED/team-11/backend/booking/pkg/ogen"
)

type tokenCtxKey struct{}

type SecurityHandler struct {
	secret string
}

func NewSecurityHandler(secret string) *SecurityHandler {
	return &SecurityHandler{
		secret: secret,
	}
}

func (sh *SecurityHandler) HandleBearerAuth(ctx context.Context, operationName api.OperationName, t api.BearerAuth) (context.Context, error) {
	jwtToken, err := jwt.Parse(t.GetToken(), func(t *jwt.Token) (interface{}, error) {
		return []byte(sh.secret), nil
	})
	if err != nil {
		return ctx, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return ctx, models.ErrInvalidToken
	}

	token, err := models.TokenFromCliams(map[string]any(claims))
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, tokenCtxKey{}, token)

	return ctx, nil
}

func TokenFromCtx(ctx context.Context) models.Token {
	return ctx.Value(tokenCtxKey{}).(models.Token)
}
