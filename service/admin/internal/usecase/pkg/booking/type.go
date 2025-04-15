package booking

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type BookingStorage interface {
	GetNearest(c ctx.Context, id string) (*entity.Booking, e.Error)
	GetNearestGuest(c ctx.Context, id string) (*entity.Booking, e.Error)
	GetStats(c ctx.Context, filter string) (int, e.Error)
}
