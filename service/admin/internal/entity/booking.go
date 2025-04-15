package entity

import (
	"time"

	"github.com/nikitaSstepanov/tools/client/pg"
)

type Booking struct {
	Id        string
	EntityId  string
	UserId    string
	TimeFrom  time.Time
	TimeTo    time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Guest struct {
	UserId    string
	BookingId string
	CreatedAt time.Time
}

func (g *Guest) Scan(r pg.Row) error {
	return r.Scan(
		&g.UserId,
		&g.BookingId,
		&g.CreatedAt,
	)
}

func (b *Booking) Scan(r pg.Row) error {
	return r.Scan(
		&b.Id,
		&b.EntityId,
		&b.UserId,
		&b.TimeFrom,
		&b.TimeTo,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
}
