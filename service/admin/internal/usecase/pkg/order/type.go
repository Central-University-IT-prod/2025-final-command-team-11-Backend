package order

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type OrderStorage interface {
	GetById(c ctx.Context, id string) (*entity.Order, e.Error)
	Get(c ctx.Context) ([]*entity.OrderBooking, e.Error)
	Update(c ctx.Context, order *entity.Order) e.Error
	GetStats(c ctx.Context, filter string) (int, e.Error)
}
