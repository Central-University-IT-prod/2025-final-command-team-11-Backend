package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"REDACTED/team-11/backend/booking/internal/dto"
	"REDACTED/team-11/backend/booking/internal/models"
	"REDACTED/team-11/backend/booking/internal/repo"
)

type OrdersService struct {
	ordersRepo   repo.OrdersRepo
	bookingsRepo repo.BookingsRepo
}

func NewOrdersService(
	ordersRepo repo.OrdersRepo,
	bookingsRepo repo.BookingsRepo,
) *OrdersService {
	return &OrdersService{
		ordersRepo:   ordersRepo,
		bookingsRepo: bookingsRepo,
	}
}

func (os *OrdersService) Create(ctx context.Context, input dto.OrderCreateDto, userId uuid.UUID) (models.Order, error) {
	op := "service.OrdersService.Create"

	booking, err := os.bookingsRepo.GetById(ctx, input.BookingId)
	if err != nil {
		if errors.Is(err, models.ErrBookingNotFound) {
			return models.Order{}, models.ErrBookingNotFound
		}

		return models.Order{}, fmt.Errorf("%s: bookingsRepo.GetById: %w", op, err)
	}

	if booking.UserId != userId {
		return models.Order{}, models.ErrNoAccessToBooking
	}

	created, err := os.ordersRepo.Create(ctx, input)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: ordersRepo.Create: %w", op, err)
	}

	return created, nil
}

func (os *OrdersService) GetForBooking(ctx context.Context, bookingId uuid.UUID, token models.Token) ([]models.Order, error) {
	op := "service.OrdersService.GetForBooking"

	booking, err := os.bookingsRepo.GetById(ctx, bookingId)
	if err != nil {
		if errors.Is(err, models.ErrBookingNotFound) {
			return nil, models.ErrBookingNotFound
		}

		return nil, fmt.Errorf("%s: bookingsRepo.GetById: %w", op, err)
	}

	if booking.UserId != token.UserId && token.Role != models.RoleAdmin && token.Role != models.RoleSuperAdmin {
		return nil, models.ErrNoAccessToBooking
	}

	orders, err := os.ordersRepo.GetForBooking(ctx, bookingId)
	if err != nil {
		return nil, fmt.Errorf("%s: ordersRepo.GetForBooking: %w", op, err)
	}

	return orders, nil
}

func (os *OrdersService) Delete(ctx context.Context, bookingId, orderId uuid.UUID, token models.Token) error {
	op := "service.OrdersService.Delete"

	order, err := os.ordersRepo.GetById(ctx, orderId)
	if err != nil {
		if errors.Is(err, models.ErrOrderNotFound) {
			return models.ErrOrderNotFound
		}

		return fmt.Errorf("%s: ordersRepo.GetById: %w", op, err)
	}

	if order.BookingId != bookingId {
		return models.ErrOrderNotFound
	}

	booking, err := os.bookingsRepo.GetById(ctx, order.BookingId)
	if err != nil {
		if errors.Is(err, models.ErrBookingNotFound) {
			return models.ErrBookingNotFound
		}

		return fmt.Errorf("%s: bookingsRepo.GetById: %w", op, err)
	}

	if booking.UserId != token.UserId && token.Role != models.RoleAdmin && token.Role != models.RoleSuperAdmin {
		return models.ErrNoAccessToBooking
	}

	err = os.ordersRepo.Delete(ctx, orderId)
	if err != nil {
		return fmt.Errorf("%s: ordersRepo.Delete: %w", op, err)
	}

	return nil
}
