package booking_entity

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type EntityStorage interface {
	UpdateEntity(c ctx.Context, ent *entity.BookingEntity) e.Error
	GetFloor(c ctx.Context, id string) (*entity.FloorEntity, e.Error)
	GetEntities(c ctx.Context, id string) ([]*entity.BookingEntity, e.Error)
	GetFloors(c ctx.Context) ([]*entity.FloorEntity, e.Error)
	CreateEntity(c ctx.Context, entity *entity.BookingEntity) e.Error
	DeleteEntity(c ctx.Context, id string) e.Error
	CreateFloor(c ctx.Context, floor *entity.FloorEntity) e.Error
	UpdateFloor(c ctx.Context, floor *entity.FloorEntity) e.Error
	DeleteFloor(c ctx.Context, id string) e.Error
	GetEntity(c ctx.Context, id string) (*entity.BookingEntity, e.Error)
}
