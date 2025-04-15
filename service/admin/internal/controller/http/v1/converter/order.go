package converter

import (
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/dto"
	"REDACTED/team-11/backend/admin/internal/entity"
)

func DtoOrder(order *entity.OrderBooking) *dto.Order {
	return &dto.Order{
		Id:        order.Order.Id,
		BookingId: order.Order.BookingId,
		Completed: order.Order.Completed,
		Thing:     order.Order.Thing,
		CreatedAt: order.Order.CreatedAt,
		UpdatedAt: order.Order.UpdatedAt,
		BookName:  order.Booking.Title,
	}
}
