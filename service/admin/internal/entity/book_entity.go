package entity

import (
	"time"

	"github.com/nikitaSstepanov/tools/client/pg"
	types "REDACTED/team-11/backend/admin/internal/entity/type"
)

type FloorEntity struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BookingEntity struct {
	Id        string
	Type      types.BookingType
	Title     string
	X         int
	Y         int
	FloorId   string
	Width     int
	Height    int
	Capacity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *BookingEntity) Scan(r pg.Row) error {
	return r.Scan(
		&b.Id,
		&b.Type,
		&b.Title,
		&b.X,
		&b.Y,
		&b.FloorId,
		&b.Width,
		&b.Height,
		&b.Capacity,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
}

func (f *FloorEntity) Scan(r pg.Row) error {
	return r.Scan(
		&f.Id,
		&f.Name,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
}
