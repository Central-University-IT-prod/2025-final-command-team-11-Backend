package booking_entity

import (
	"slices"

	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type BookingEntity struct {
	booking EntityStorage
}

func New(booking EntityStorage) *BookingEntity {
	return &BookingEntity{
		booking: booking,
	}
}

func (b *BookingEntity) GetEntity(c ctx.Context, id string) (*entity.BookingEntity, e.Error) {
	return b.booking.GetEntity(c, id)
}

func (b *BookingEntity) GetFloors(c ctx.Context) ([]*entity.FloorEntity, e.Error) {
	return b.booking.GetFloors(c)
}

func (b *BookingEntity) GetEntities(c ctx.Context, id string) ([]*entity.BookingEntity, e.Error) {
	_, err := b.booking.GetFloor(c, id)
	if err != nil {
		return nil, err
	}

	return b.booking.GetEntities(c, id)
}

func (b *BookingEntity) Save(c ctx.Context, bookings []*entity.BookingEntity, floorEntity *entity.FloorEntity) e.Error {
	floor, err := b.booking.GetFloor(c, floorEntity.Id)
	if err != nil && err.GetCode() != e.NotFound {
		return err
	}

	if err != nil {
		floor := &entity.FloorEntity{
			Id:        floorEntity.Id,
			Name:      floorEntity.Name,
			CreatedAt: floorEntity.CreatedAt,
			UpdatedAt: floorEntity.UpdatedAt,
		}

		err := b.booking.CreateFloor(c, floor)
		if err != nil {
			return err
		}
	} else {
		floor.UpdatedAt = floorEntity.UpdatedAt
		floor.Name = floorEntity.Name

		err := b.booking.UpdateFloor(c, floor)
		if err != nil {
			return err
		}
	}

	entities, err := b.booking.GetEntities(c, floorEntity.Id)
	if err != nil && err.GetCode() != e.NotFound {
		return err
	}

	id1 := make([]string, 0)

	for _, entit := range entities {
		id1 = append(id1, entit.Id)
	}

	id2 := make([]string, 0)

	for _, booking := range bookings {
		id2 = append(id2, booking.Id)
	}

	toDel := make([]string, 0)

	for _, i := range id1 {
		if !slices.Contains(id2, i) {
			toDel = append(toDel, i)
		}
	}

	for _, to := range toDel {
		if err := b.booking.DeleteEntity(c, to); err != nil {
			return err
		}
	}

	for _, u := range bookings {
		_, err := b.booking.GetEntity(c, u.Id)
		if err != nil && err.GetCode() != e.NotFound {
			return err
		}

		if err != nil {
			err := b.booking.CreateEntity(c, u)
			if err != nil {
				return err
			}
		} else {
			err := b.booking.UpdateEntity(c, u)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *BookingEntity) DeleteFloor(c ctx.Context, id string) e.Error {
	_, err := b.booking.GetFloor(c, id)
	if err != nil {
		return err
	}

	return b.booking.DeleteFloor(c, id)
}
