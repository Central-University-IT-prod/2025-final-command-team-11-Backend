package entity

import (
	"time"

	"github.com/nikitaSstepanov/tools/client/pg"
)

type Order struct {
	Id        string
	BookingId string
	Completed bool
	Thing     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderBooking struct {
	Order   *Order
	Booking *BookingEntity
}

func (o *Order) Scan(r pg.Row) error {
	return r.Scan(
		&o.Id,
		&o.BookingId,
		&o.Completed,
		&o.Thing,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
}
