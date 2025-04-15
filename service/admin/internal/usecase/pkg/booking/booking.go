package booking

import (
	"fmt"

	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type Booking struct {
	booking BookingStorage
}

func New(booking BookingStorage) *Booking {
	return &Booking{
		booking: booking,
	}
}

func (b *Booking) CheckAccess(c ctx.Context, id string) (*entity.Booking, e.Error) {
	booking, err := b.booking.GetNearest(c, id)
	if err != nil && err.GetCode() != e.NotFound {
		return nil, err
	}
	fmt.Println(err)
	if err == nil {
		return booking, nil
	}

	booking, err = b.booking.GetNearestGuest(c, id)
	if err != nil && err.GetCode() != e.NotFound {
		return nil, err
	}

	if err == nil {
		return booking, nil
	}

	return nil, e.New("User hasn`t access.", e.NotFound)
}

func (b *Booking) Stats(c ctx.Context, filter string) (int, e.Error) {
	return b.booking.GetStats(c, filter)
}
