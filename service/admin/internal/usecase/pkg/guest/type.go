package guest

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type IdUseCase interface {
	GetUser(c ctx.Context, email string) (*entity.User, e.Error)
	GetUserById(c ctx.Context, id string) (*entity.User, e.Error)
}

type GuestStorage interface {
	Create(c ctx.Context, guest *entity.Guest) e.Error
	Get(c ctx.Context, id string) ([]*entity.Guest, e.Error)
	GetById(c ctx.Context, bookId, userId string) (*entity.Guest, e.Error)
	Delete(c ctx.Context, bookId string, userId string) e.Error
}

type EntityStorage interface {
	GetEntity(c ctx.Context, id string) (*entity.BookingEntity, e.Error)
}

type BookingStorage interface {
	GetById(c ctx.Context, id string) (*entity.Booking, e.Error)
}
