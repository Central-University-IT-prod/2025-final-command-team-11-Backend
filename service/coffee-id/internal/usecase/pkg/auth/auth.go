package auth

import (
	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/utils/coder"
)

type Auth struct {
	user      UserStorage
	oauth     OAuthStorage
	yndxOAuth YandexOAuth
	coder     *coder.Coder
	Jwt       *Jwt
}

func New(store *Storages, uc *UseCases) *Auth {
	return &Auth{
		user:      store.User,
		yndxOAuth: uc.YandexOAuth,
		oauth:     store.OAuth,
		coder:     uc.Coder,
		Jwt:       uc.Jwt,
	}
}

func (a *Auth) Login(c ctx.Context, user *entity.User) (*entity.Tokens, *entity.User, e.Error) {
	candidate, err := a.user.GetByEmail(c, user.Email)
	if err != nil {
		return nil, nil, err
	}

	if candidate.Id == "" {
		return nil, nil, badDataErr
	}

	if candidate.OAuth != "" {
		return nil, nil, notFoundErr
	}

	if err := a.coder.CompareHash(candidate.Password, user.Password); err != nil {
		return nil, nil, badDataErr.WithErr(err)
	}

	tokens, err := a.getTokens(candidate)
	if err != nil {
		return nil, nil, err
	}

	return tokens, candidate, nil
}

func (a *Auth) Refresh(ctx ctx.Context, refresh string) (*entity.Tokens, e.Error) {
	claims, err := a.Jwt.ValidateToken(refresh, true)
	if err != nil {
		return nil, err
	}

	user, err := a.user.GetById(ctx, claims.Id)
	if err != nil {
		return nil, err
	}

	return a.getTokens(user)
}

func (a *Auth) ValidateToken(jwtString string, isRefresh bool) (*Claims, e.Error) {
	return a.Jwt.ValidateToken(jwtString, isRefresh)
}

func (a *Auth) getTokens(user *entity.User) (*entity.Tokens, e.Error) {
	access, err := a.Jwt.GenerateToken(user, false)
	if err != nil {
		return nil, err
	}

	refresh, err := a.Jwt.GenerateToken(user, true)
	if err != nil {
		return nil, err
	}

	return &entity.Tokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}
