package booking

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type BookingUseCase interface {
	CheckAccess(c ctx.Context, id string) (*entity.Booking, e.Error)
	Stats(c ctx.Context, filter string) (int, e.Error)
}
