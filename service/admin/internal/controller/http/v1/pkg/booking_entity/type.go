package booking_entity

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type EntityuseCase interface {
	Save(c ctx.Context, bookings []*entity.BookingEntity, floorEntity *entity.FloorEntity) e.Error
	GetEntities(c ctx.Context, id string) ([]*entity.BookingEntity, e.Error)
	GetFloors(c ctx.Context) ([]*entity.FloorEntity, e.Error)
	DeleteFloor(c ctx.Context, id string) e.Error
	GetEntity(c ctx.Context, id string) (*entity.BookingEntity, e.Error)
}
