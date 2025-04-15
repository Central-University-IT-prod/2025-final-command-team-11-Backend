package account

import (
	"github.com/nikitaSstepanov/coffee-id/internal/entity"
	"github.com/nikitaSstepanov/coffee-id/internal/usecase/admin"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/utils/coder"
)

type UseCases struct {
	Jwt   JwtUseCase
	Mail  MailUseCase
	Coder *coder.Coder
	Admin
}

type Storages struct {
	User UserStorage
	Code CodeStorage
}

type JwtUseCase interface {
	GenerateToken(user *entity.User, isRefresh bool) (string, e.Error)
}

type MailUseCase interface {
	SendActivation(c ctx.Context, to string, code string) e.Error
}

type UserStorage interface {
	Get(c ctx.Context, limit, offset int) ([]*entity.User, e.Error)
	Getc(c ctx.Context) ([]*entity.User, e.Error)
	GetById(c ctx.Context, id string) (*entity.User, e.Error)
	GetByEmail(c ctx.Context, email string) (*entity.User, e.Error)
	Create(c ctx.Context, user *entity.User) e.Error
	Update(ctx ctx.Context, user *entity.User) (*entity.User, e.Error)
	AddRole(c ctx.Context, user *entity.User) e.Error
	Verify(c ctx.Context, user *entity.User) e.Error
	Delete(c ctx.Context, user *entity.User) e.Error
}

type CodeStorage interface {
	Get(c ctx.Context, userId uint64) (*entity.ActivationCode, e.Error)
	Set(c ctx.Context, code *entity.ActivationCode) e.Error
	Del(c ctx.Context, userId uint64) e.Error
}

type ClientStorage interface {
	GetById(c ctx.Context, id uint64) (*entity.Client, e.Error)
}

type Admin interface {
	GetUser(c ctx.Context, email string, tokem string) (*admin.User, e.Error)
}
