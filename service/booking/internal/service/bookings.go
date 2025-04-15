package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"REDACTED/team-11/backend/booking/internal/dto"
	"REDACTED/team-11/backend/booking/internal/models"
	"REDACTED/team-11/backend/booking/internal/repo"
)

type BookingsService struct {
	bookingsRepo        repo.BookingsRepo
	bookingEntitiesRepo repo.BookingEntitiesRepo
	ordersRepo          repo.OrdersRepo
	workloadsService    WorkloadsService
	usersRepo           repo.UsersRepo
}

func NewBookingsService(
	bookingsRepo repo.BookingsRepo,
	bookingEntitiesRepo repo.BookingEntitiesRepo,
	ordersRepo repo.OrdersRepo,
	workloadsService WorkloadsService,
	usersRepo repo.UsersRepo,
) *BookingsService {
	return &BookingsService{
		bookingsRepo:        bookingsRepo,
		bookingEntitiesRepo: bookingEntitiesRepo,
		ordersRepo:          ordersRepo,
		workloadsService:    workloadsService,
		usersRepo:           usersRepo,
	}
}

func (bs *BookingsService) Create(ctx context.Context, input dto.BookingCreateDto) (models.Booking, error) {
	op := "service.BookingsService.Create"

	intersected, err := bs.bookingsRepo.ListIntersectedForUser(ctx, input.UserId, input.TimeFrom, input.TimeTo)
	if err != nil {
		return models.Booking{}, fmt.Errorf("%s: bookingsRepo.ListIntersectedForUser: %w", op, err)
	}

	if len(intersected) != 0 {
		return models.Booking{}, models.ErrAlreadyHaveBooking
	}

	workload, err := bs.workloadsService.Get(ctx, input.EntityId, input.TimeFrom, input.TimeTo)
	if err != nil {
		if errors.Is(err, models.ErrBookingEntityNotFound) {
			return models.Booking{}, models.ErrBookingEntityNotFound
		}

		return models.Booking{}, fmt.Errorf("%s: workloadsService.Get: %w", op, err)
	}

	free := true
	for _, snapshot := range workload {
		if !snapshot.IsFree {
			free = false
			break
		}
	}

	if !free {
		return models.Booking{}, models.ErrNoFreePlaces
	}

	res, err := bs.bookingsRepo.Create(ctx, input)
	if err != nil {
		return models.Booking{}, fmt.Errorf("%s: bookingsRepo.Create: %w", op, err)
	}

	return res, nil
}

func (bs *BookingsService) GetById(ctx context.Context, bookingId uuid.UUID, token models.Token) (models.BookingInfo, error) {
	op := "service.BookingsService.GetById"

	booking, err := bs.bookingsRepo.GetById(ctx, bookingId)
	if err != nil {
		if errors.Is(err, models.ErrBookingNotFound) {
			return models.BookingInfo{}, models.ErrBookingNotFound
		}

		return models.BookingInfo{}, fmt.Errorf("%s: bookingsRepo.GetById: %w", op, err)
	}

	if booking.UserId != token.UserId && token.Role != models.RoleAdmin && token.Role != models.RoleSuperAdmin {
		return models.BookingInfo{}, models.ErrNoAccessToBooking
	}

	orders, err := bs.ordersRepo.GetForBooking(ctx, bookingId)
	if err != nil {
		return models.BookingInfo{}, fmt.Errorf("%s: ordersRepo.GetForBooking: %w", op, err)
	}

	user, err := bs.usersRepo.GetById(ctx, booking.UserId)
	if err != nil {
		return models.BookingInfo{}, fmt.Errorf("%s: usersRepo.GetById: %w", op, err)
	}

	entity, err := bs.bookingEntitiesRepo.GetById(ctx, booking.EntityId)
	if err != nil {
		return models.BookingInfo{}, fmt.Errorf("%s: bookingEntitiesRepo.GetById: %w", op, err)
	}

	bookingInfo := models.BookingInfo{
		Booking: booking,
		Entity:  entity,
		Orders:  orders,
		User:    user,
	}

	return bookingInfo, nil
}

func (bs *BookingsService) ListForUser(ctx context.Context, userId uuid.UUID) ([]models.BookingInfo, error) {
	op := "service.BookingsService.ListForUser"

	bookings, err := bs.bookingsRepo.ListForUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: bookingsRepo.ListForUser: %w", op, err)
	}

	bookingsInfos := make([]models.BookingInfo, 0, len(bookings))

	for _, booking := range bookings {
		user, err := bs.usersRepo.GetById(ctx, booking.UserId)
		if err != nil {
			return nil, fmt.Errorf("%s: usersRepo.GetById: %w", op, err)
		}

		orders, err := bs.ordersRepo.GetForBooking(ctx, booking.Id)
		if err != nil {
			return nil, fmt.Errorf("%s: ordersRepo.GetForBooking: %w", op, err)
		}

		entity, err := bs.bookingEntitiesRepo.GetById(ctx, booking.EntityId)
		if err != nil {
			return nil, fmt.Errorf("%s: bookingEntitiesRepo.GetById: %w", op, err)
		}

		bookingsInfos = append(bookingsInfos, models.BookingInfo{
			Booking: booking,
			Entity:  entity,
			User:    user,
			Orders:  orders,
		})
	}

	return bookingsInfos, nil
}

func (bs *BookingsService) ListAll(ctx context.Context, token models.Token) ([]models.BookingInfo, error) {
	op := "service.BookingService.ListAll"

	if token.Role != models.RoleAdmin && token.Role != models.RoleSuperAdmin {
		return nil, models.ErrNoRights
	}

	bookings, err := bs.bookingsRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: bookingsRepo.ListForUser: %w", op, err)
	}

	bookingsInfos := make([]models.BookingInfo, 0, len(bookings))

	for _, booking := range bookings {
		user, err := bs.usersRepo.GetById(ctx, booking.UserId)
		if err != nil {
			return nil, fmt.Errorf("%s: usersRepo.GetById: %w", op, err)
		}

		orders, err := bs.ordersRepo.GetForBooking(ctx, booking.Id)
		if err != nil {
			return nil, fmt.Errorf("%s: ordersRepo.GetForBooking: %w", op, err)
		}

		entity, err := bs.bookingEntitiesRepo.GetById(ctx, booking.EntityId)
		if err != nil {
			return nil, fmt.Errorf("%s: bookingEntitiesRepo.GetById: %w", op, err)
		}

		bookingsInfos = append(bookingsInfos, models.BookingInfo{
			Booking: booking,
			Entity:  entity,
			User:    user,
			Orders:  orders,
		})
	}

	return bookingsInfos, nil
}

func (bs *BookingsService) Update(ctx context.Context, input dto.BookingUpdateDto, token models.Token) (models.Booking, error) {
	op := "service.BookingsService.Update"

	booking, err := bs.bookingsRepo.GetById(ctx, input.BookingId)
	if err != nil {
		if errors.Is(err, models.ErrBookingNotFound) {
			return models.Booking{}, models.ErrBookingNotFound
		}

		return models.Booking{}, fmt.Errorf("%s: bookingsRepo.GetById: %w", op, err)
	}

	if booking.UserId != token.UserId && token.Role != models.RoleAdmin && token.Role != models.RoleSuperAdmin {
		return models.Booking{}, models.ErrNoAccessToBooking
	}

	var (
		resTimeFrom time.Time = booking.TimeFrom
		resTimeTo   time.Time = booking.TimeTo
	)

	if input.TimeFrom != nil {
		resTimeFrom = *input.TimeFrom
	}
	if input.TimeTo != nil {
		resTimeTo = *input.TimeTo
	}

	if !resTimeFrom.Before(resTimeTo) {
		return models.Booking{}, models.ErrInvalidBookingTime
	}

	intersectedMayBeWithSame, err := bs.bookingsRepo.ListIntersectedForUser(ctx, token.UserId, resTimeFrom, resTimeTo)
	if err != nil {
		return models.Booking{}, fmt.Errorf("%s: bookingsRepo.ListIntersectedForUser: %w", op, err)
	}

	intersected := []models.Booking{}
	for _, booking := range intersectedMayBeWithSame {
		if booking.Id != input.BookingId {
			intersected = append(intersected, booking)
		}
	}

	if len(intersected) != 0 {
		return models.Booking{}, models.ErrAlreadyHaveBooking
	}

	workload, err := bs.workloadsService.Get(ctx, booking.EntityId, resTimeFrom, resTimeTo)
	if err != nil {
		return models.Booking{}, fmt.Errorf("%s: workloadsService.Get: %w", op, err)
	}

	free := true
	for _, snapshot := range workload {
		if booking.TimeFrom.UTC().Unix() <= snapshot.Time.UTC().Unix() &&
			snapshot.Time.UTC().Unix() <= booking.TimeTo.UTC().Unix() {
			continue
		}
		if !snapshot.IsFree {
			free = false
			break
		}
	}

	if !free {
		return models.Booking{}, models.ErrNoFreePlaces
	}

	updated, err := bs.bookingsRepo.Update(ctx, input)
	if err != nil {
		return models.Booking{}, fmt.Errorf("%s: bookingsRepo.Update: %w", op, err)
	}

	return updated, nil
}

func (bs *BookingsService) Delete(ctx context.Context, bookingId uuid.UUID, token models.Token) error {
	op := "service.BookingsService.Delete"

	booking, err := bs.bookingsRepo.GetById(ctx, bookingId)
	if err != nil {
		if errors.Is(err, models.ErrBookingNotFound) {
			return models.ErrBookingNotFound
		}

		return fmt.Errorf("%s: bookingsRepo.GetById: %w", op, err)
	}

	if booking.UserId != token.UserId && token.Role != models.RoleAdmin && token.Role != models.RoleSuperAdmin {
		return models.ErrNoAccessToBooking
	}

	err = bs.bookingsRepo.Delete(ctx, bookingId)
	if err != nil {
		return fmt.Errorf("%s: bookingsRepo.Delete: %w", op, err)
	}

	return nil
}
