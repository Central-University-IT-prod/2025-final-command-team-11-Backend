package handlers

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"REDACTED/team-11/backend/booking/internal/dto"
	"REDACTED/team-11/backend/booking/internal/models"
	"REDACTED/team-11/backend/booking/internal/transport/http/v1/security"
	"REDACTED/team-11/backend/booking/pkg/logger"
	api "REDACTED/team-11/backend/booking/pkg/ogen"
	"go.uber.org/zap"
)

type OrdersUsecase interface {
	Create(ctx context.Context, input dto.OrderCreateDto, userId uuid.UUID) (models.Order, error)
	GetForBooking(ctx context.Context, bookingId uuid.UUID, token models.Token) ([]models.Order, error)
	Delete(ctx context.Context, bookingId, orderId uuid.UUID, token models.Token) error
}

type OrdersHandler struct {
	usecase OrdersUsecase
}

func NewOrdersHandler(
	usecase OrdersUsecase,
) *OrdersHandler {
	return &OrdersHandler{
		usecase: usecase,
	}
}

// CreateOrder implements createOrder operation.
//
// Create order.
//
// POST /bookings/{bookingId}/orders
func (oh *OrdersHandler) CreateOrder(ctx context.Context, req *api.OrderCreate, params api.CreateOrderParams) (api.CreateOrderRes, error) {
	token := security.TokenFromCtx(ctx)

	created, err := oh.usecase.Create(ctx, dto.OrderCreateDto{
		BookingId: params.BookingId,
		Thing:     string(req.GetThing()),
	}, token.UserId)
	if err != nil {
		if errors.Is(err, models.ErrBookingNotFound) {
			return &api.Response404{
				Resource: api.NewOptResponse404Resource(api.Response404ResourceBooking),
			}, nil
		}
		if errors.Is(err, models.ErrNoAccessToBooking) {
			return &api.Response404{
				Resource: api.NewOptResponse404Resource(api.Response404ResourceBooking),
			}, nil
		}

		logger.FromCtx(ctx).Error("create order", zap.Error(err))
		return nil, err
	}

	res := convertOrder(created)
	return &res, nil
}

// DeleteOrders implements deleteOrders operation.
//
// Delete order.
//
// DELETE /bookings/{bookingId}/orders/{orderId}
func (oh *OrdersHandler) DeleteOrders(ctx context.Context, params api.DeleteOrdersParams) (api.DeleteOrdersRes, error) {
	token := security.TokenFromCtx(ctx)

	err := oh.usecase.Delete(ctx, params.BookingId, params.OrderId, token)
	if err != nil {
		if errors.Is(err, models.ErrOrderNotFound) {
			return &api.Response404{
				Resource: api.NewOptResponse404Resource(api.Response404ResourceOrder),
			}, nil
		}
		if errors.Is(err, models.ErrNoAccessToBooking) {
			return &api.Response404{
				Resource: api.NewOptResponse404Resource(api.Response404ResourceOrder),
			}, nil
		}

		logger.FromCtx(ctx).Error("delete order", zap.Error(err))
		return nil, err
	}

	return &api.DeleteOrdersNoContent{}, nil
}

// ListOrders implements listOrders operation.
//
// List orders.
//
// GET /bookings/{bookingId}/orders
func (oh *OrdersHandler) ListOrders(ctx context.Context, params api.ListOrdersParams) (api.ListOrdersRes, error) {
	token := security.TokenFromCtx(ctx)

	orders, err := oh.usecase.GetForBooking(ctx, params.BookingId, token)
	if err != nil {
		if errors.Is(err, models.ErrBookingNotFound) {
			return &api.Response404{
				Resource: api.NewOptResponse404Resource(api.Response404ResourceBooking),
			}, nil
		}
		if errors.Is(err, models.ErrNoAccessToBooking) {
			return &api.Response404{
				Resource: api.NewOptResponse404Resource(api.Response404ResourceBooking),
			}, nil
		}

		logger.FromCtx(ctx).Error("list orders for booking", zap.Error(err))
		return nil, err
	}

	res := api.ListOrdersOKApplicationJSON(make(api.ListOrdersOKApplicationJSON, 0, len(orders)))
	for _, order := range orders {
		res = append(res, convertOrder(order))
	}

	return &res, nil
}

func convertOrder(order models.Order) api.Order {
	return api.Order{
		ID:        order.Id,
		BookingID: order.BookingId,
		Completed: order.Completed,
		Thing:     api.OrderThingEnum(order.Thing),
		CreatedAt: api.Time(order.CreatedAt.Unix()),
		UpdatedAt: api.Time(order.UpdatedAt.Unix()),
	}
}
