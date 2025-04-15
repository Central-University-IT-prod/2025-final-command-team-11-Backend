package account

import (
	"github.com/nikitaSstepanov/coffee-id/internal/controller/http/v1/dto"
	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
)

type AccountUseCase interface {
	GetList(c ctx.Context, page, size int, token string) ([]*dto.AP, int, e.Error)
	Get(ctx ctx.Context, userId string) (*entity.User, e.Error)
	Create(ctx ctx.Context, user *entity.User) (*entity.Tokens, e.Error)
	Update(ctx ctx.Context, user *entity.User, pass string) (*entity.User, e.Error)
	AddRole(ctx ctx.Context, user *entity.User) e.Error
	Delete(ctx ctx.Context, user *entity.User) e.Error
	GetByEmail(c ctx.Context, email string) (*entity.User, e.Error)
}
