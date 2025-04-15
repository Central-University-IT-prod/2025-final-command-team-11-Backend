package http

import api "REDACTED/team-11/backend/booking/pkg/ogen"

type Handler struct {
	api.BookingsHandler
	api.OrdersHandler
	api.WorkloadsHandler
}

func NewHandler(
	bookingsHandler api.BookingsHandler,
	ordersHandler api.OrdersHandler,
	workloadsHandler api.WorkloadsHandler,
) api.Handler {
	return &Handler{
		BookingsHandler:  bookingsHandler,
		OrdersHandler:    ordersHandler,
		WorkloadsHandler: workloadsHandler,
	}
}
