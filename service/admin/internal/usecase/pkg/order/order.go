package order

import (
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"REDACTED/team-11/backend/admin/internal/entity"
)

type Order struct {
	order OrderStorage
}

func New(order OrderStorage) *Order {
	return &Order{
		order: order,
	}
}

func (o *Order) ChangeStatus(c ctx.Context, id string, status bool) e.Error {
	order, err := o.order.GetById(c, id)
	if err != nil {
		return err
	}

	order.Completed = status

	return o.order.Update(c, order)
}

func (o *Order) GetOrders(c ctx.Context, page, size int, filter string) ([]*entity.OrderBooking, int, e.Error) {
	all, err := o.order.Get(c)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.OrderBooking, 0)

	if filter == "" {
		result = all
	} else {
		for _, order := range all {
			if (order.Order.Completed && filter == "true") || (!order.Order.Completed && filter == "false") {
				result = append(result, order)
			}
		}
	}

	if page*size >= len(result) {
		return nil, len(result), nil
	}

	right := min(page*size+size, len(result))

	return result[page*size : right], len(result), nil
}

func (o *Order) Stats(c ctx.Context, filter string) (int, e.Error) {
	return o.order.GetStats(c, filter)
}
