package auth

import (
	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
)

type AuthUseCase interface {
	Login(c ctx.Context, user *entity.User) (*entity.Tokens, *entity.User, e.Error)
	Refresh(ctx ctx.Context, refresh string) (*entity.Tokens, e.Error)
}
