package order

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type OrderUseCase interface {
	ChangeStatus(c ctx.Context, id string, status bool) e.Error
	GetOrders(c ctx.Context, page, size int, filter string) ([]*entity.OrderBooking, int, e.Error)
	Stats(c ctx.Context, filter string) (int, e.Error)
}
