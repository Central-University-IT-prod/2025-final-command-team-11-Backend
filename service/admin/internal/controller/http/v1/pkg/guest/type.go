package guest

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type GuestUseCase interface {
	Create(c ctx.Context, bookId, email, owner string) e.Error 
	Get(c ctx.Context, bookingId string, userId string) ([]*entity.Guest, e.Error)
	Delete(c ctx.Context, bookId, email, owner string) e.Error 
}
