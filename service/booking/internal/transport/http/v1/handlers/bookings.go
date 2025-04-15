package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"REDACTED/team-11/backend/booking/internal/dto"
	"REDACTED/team-11/backend/booking/internal/models"
	"REDACTED/team-11/backend/booking/internal/transport/http/v1/security"
	"REDACTED/team-11/backend/booking/pkg/logger"
	api "REDACTED/team-11/backend/booking/pkg/ogen"
	"go.uber.org/zap"
)

var (
	bookingIntervalMinutes int64 = 15
)

type BookingUsecase interface {
	Create(ctx context.Context, input dto.BookingCreateDto) (models.Booking, error)
	GetById(ctx context.Context, bookingId uuid.UUID, token models.Token) (models.BookingInfo, error)
	ListAll(ctx context.Context, token models.Token) ([]models.BookingInfo, error)
	ListForUser(ctx context.Context, userId uuid.UUID) ([]models.BookingInfo, error)
	Update(ctx context.Context, input dto.BookingUpdateDto, token models.Token) (models.Booking, error)
	Delete(ctx context.Context, bookingId uuid.UUID, token models.Token) error
}

type BookingsHandler struct {
	usecase BookingUsecase
}

func NewBookingsHandler(
	usecase BookingUsecase,
) *BookingsHandler {
	return &BookingsHandler{
		usecase: usecase,
	}
}

// CreateBooking implements createBooking operation.
//
// Create booking.
//
// POST /bookings
func (bh *BookingsHandler) CreateBooking(ctx context.Context, req *api.BookingCreate) (api.CreateBookingRes, error) {
	token := security.TokenFromCtx(ctx)

	if req.GetTimeFrom() >= req.GetTimeTo() {
		return &api.Response400{
			Message: api.NewOptString("time_from must be before time_to"),
		}, nil
	}

	if int64(req.GetTimeFrom())%(bookingIntervalMinutes*60) != 0 {
		return &api.Response400{
			Message: api.NewOptString("time_from must be multiple of 15 minutes"),
		}, nil
	}

	if int64(req.GetTimeTo())%(bookingIntervalMinutes*60) != 0 {
		return &api.Response400{
			Message: api.NewOptString("time_to must be multiple of 15 minutes"),
		}, nil
	}

	timeFrom := time.Unix(int64(req.GetTimeFrom()), 0).UTC()
	timeTo := time.Unix(int64(req.GetTimeTo()), 0).UTC()

	booking, err := bh.usecase.Create(ctx, dto.BookingCreateDto{
		EntityId: req.GetEntityID(),
		UserId:   token.UserId,
		TimeFrom: timeFrom,
		TimeTo:   timeTo,
	})

	if err != nil {
		if errors.Is(err, models.ErrBookingEntityNotFound) {
			return &api.Response404{
				Resource: api.NewOptResponse404Resource(api.Response404ResourceBookingEntity),
			}, nil
		}
		if errors.Is(err, models.ErrAlreadyHaveBooking) {
			return &api.CreateBookingConflict{}, nil
		}
		if errors.Is(err, models.ErrNoFreePlaces) {
			return &api.CreateBookingForbidden{}, nil
		}

		logger.FromCtx(ctx).Error("create booking", zap.Error(err))
		return nil, err
	}

	res := convertBooking(booking)
	return &res, nil

}

func (bh *BookingsHandler) CreateBookingForAdmin(ctx context.Context, req *api.BookingCreate, params api.CreateBookingForAdminParams) (api.CreateBookingForAdminRes, error) {
	token := security.TokenFromCtx(ctx)

	if token.Role != "ADMIN" {
		return &api.CreateBookingForAdminForbidden{}, nil
	}

	if req.GetTimeFrom() >= req.GetTimeTo() {
		return &api.Response400{
			Message: api.NewOptString("time_from must be before time_to"),
		}, nil
	}

	if int64(req.GetTimeFrom())%(bookingIntervalMinutes*60) != 0 {
		return &api.Response400{
			Message: api.NewOptString("time_from must be multiple of 15 minutes"),
		}, nil
	}

	if int64(req.GetTimeTo())%(bookingIntervalMinutes*60) != 0 {
		return &api.Response400{
			Message: api.NewOptString("time_to must be multiple of 15 minutes"),
		}, nil
	}

	timeFrom := time.Unix(int64(req.GetTimeFrom()), 0).UTC()
	timeTo := time.Unix(int64(req.GetTimeTo()), 0).UTC()

	booking, err := bh.usecase.Create(ctx, dto.BookingCreateDto{
		EntityId: req.GetEntityID(),
		UserId:   params.UserId,
		TimeFrom: timeFrom,
		TimeTo:   timeTo,
	})

	if err != nil {
		if errors.Is(err, models.ErrBookingEntityNotFound) {
			return &api.Response404{
				Resource: api.NewOptResponse404Resource(api.Response404ResourceBookingEntity),
			}, nil
		}
		if errors.Is(err, models.ErrAlreadyHaveBooking) {
			return &api.CreateBookingForAdminConflict{}, nil
		}
		if errors.Is(err, models.ErrNoFreePlaces) {
			return &api.CreateBookingForAdminForbidden{}, nil
		}

		logger.FromCtx(ctx).Error("create booking", zap.Error(err))
		return nil, err
	}

	res := convertBooking(booking)
	return &res, nil
}

// DeleteBooking implements deleteBooking operation.
//
// Delete booking by ID.
//
// DELETE /bookings/{bookingId}
func (bh *BookingsHandler) DeleteBooking(ctx context.Context, params api.DeleteBookingParams) (api.DeleteBookingRes, error) {
	token := security.TokenFromCtx(ctx)

	err := bh.usecase.Delete(ctx, params.BookingId, token)
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

		logger.FromCtx(ctx).Error("delete booking", zap.Error(err))
		return nil, err
	}

	return &api.DeleteBookingNoContent{}, nil
}

// GetBookingById implements getBookingById operation.
//
// Get booking by ID.
//
// GET /bookings/{bookingId}
func (bh *BookingsHandler) GetBookingById(ctx context.Context, params api.GetBookingByIdParams) (api.GetBookingByIdRes, error) {
	token := security.TokenFromCtx(ctx)

	bookingInfo, err := bh.usecase.GetById(ctx, params.BookingId, token)
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

		logger.FromCtx(ctx).Error("get booking by id", zap.Error(err))
		return nil, err
	}

	res := convertBookingInfo(bookingInfo)
	return &res, nil
}

// ListMyBookings implements listMyBookings operation.
//
// Get list of my bookings.
//
// GET /bookings/my
func (bh *BookingsHandler) ListMyBookings(ctx context.Context) (api.ListMyBookingsRes, error) {
	token := security.TokenFromCtx(ctx)

	bookings, err := bh.usecase.ListForUser(ctx, token.UserId)
	if err != nil {
		logger.FromCtx(ctx).Error("list my bookings", zap.Error(err))
		return nil, err
	}

	res := api.ListMyBookingsOKApplicationJSON(make([]api.BookingInfo, 0, len(bookings)))
	for _, booking := range bookings {
		res = append(res, convertBookingInfo(booking))
	}

	return &res, nil
}

// ListAllBookings implements listAllBookings operation.
//
// Возвращает список всех бронирований.
//
// GET /bookings
func (bh *BookingsHandler) ListAllBookings(ctx context.Context) (api.ListAllBookingsRes, error) {
	token := security.TokenFromCtx(ctx)

	bookings, err := bh.usecase.ListAll(ctx, token)
	if err != nil {
		if errors.Is(err, models.ErrNoRights) {
			return &api.ListAllBookingsForbidden{}, nil
		}

		logger.FromCtx(ctx).Error("list all books", zap.Error(err))
		return nil, err
	}

	res := api.ListAllBookingsOKApplicationJSON(make([]api.BookingInfo, 0, len(bookings)))
	for _, booking := range bookings {
		res = append(res, convertBookingInfo(booking))
	}

	return &res, nil
}

// UpdateBooking implements updateBooking operation.
//
// Update booking by ID.
//
// PUT /bookings/{bookingId}
func (bh *BookingsHandler) UpdateBooking(ctx context.Context, req *api.BookingUpdate, params api.UpdateBookingParams) (api.UpdateBookingRes, error) {
	token := security.TokenFromCtx(ctx)

	if req.GetTimeFrom().IsSet() && int64(req.GetTimeFrom().Value)%(bookingIntervalMinutes*60) != 0 {
		return &api.Response400{
			Message: api.NewOptString("time_from must be multiple of 15 minutes"),
		}, nil
	}

	if req.GetTimeTo().IsSet() && int64(req.GetTimeTo().Value)%(bookingIntervalMinutes*60) != 0 {
		return &api.Response400{
			Message: api.NewOptString("time_to must be multiple of 15 minutes"),
		}, nil
	}

	if req.GetTimeFrom().IsSet() && req.GetTimeTo().IsSet() && req.GetTimeFrom().Value >= req.GetTimeTo().Value {
		return &api.Response400{
			Message: api.NewOptString("time_from must be before time_to"),
		}, nil
	}

	var (
		timeFrom *time.Time
		timeTo   *time.Time
	)

	if req.GetTimeFrom().IsSet() {
		timeFrom = pointer(time.Unix(int64(req.GetTimeFrom().Value), 0).UTC())
	}
	if req.GetTimeTo().IsSet() {
		timeTo = pointer(time.Unix(int64(req.GetTimeTo().Value), 0).UTC())
	}

	updated, err := bh.usecase.Update(ctx, dto.BookingUpdateDto{
		BookingId: params.BookingId,
		TimeFrom:  timeFrom,
		TimeTo:    timeTo,
	}, token)
	if err != nil {
		if errors.Is(err, models.ErrInvalidBookingTime) {
			return &api.Response400{
				Message: api.NewOptString("time_from must be before time_to"),
			}, nil
		}
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
		if errors.Is(err, models.ErrAlreadyHaveBooking) {
			return &api.UpdateBookingConflict{}, nil
		}
		if errors.Is(err, models.ErrNoFreePlaces) {
			return &api.UpdateBookingForbidden{}, nil
		}

		logger.FromCtx(ctx).Error("update booking", zap.Error(err))
		return nil, err
	}

	res := convertBooking(updated)
	return &res, nil
}

func convertBooking(booking models.Booking) api.Booking {
	return api.Booking{
		ID:        booking.Id,
		EntityID:  booking.EntityId,
		UserID:    booking.UserId,
		TimeFrom:  api.Time(booking.TimeFrom.Unix()),
		TimeTo:    api.Time(booking.TimeTo.Unix()),
		CreatedAt: api.Time(booking.CreatedAt.Unix()),
		UpdatedAt: api.Time(booking.UpdatedAt.Unix()),
	}
}

func convertBookingInfo(bookingInfo models.BookingInfo) api.BookingInfo {
	orders := make([]api.Order, 0, len(bookingInfo.Orders))
	for _, order := range bookingInfo.Orders {
		orders = append(orders, convertOrder(order))
	}

	return api.BookingInfo{
		ID:     bookingInfo.Id,
		Entity: convertBookingEntity(bookingInfo.Entity),
		User: api.User{
			ID:    bookingInfo.User.Id,
			Email: bookingInfo.User.Email,
			Name:  bookingInfo.User.Name,
		},
		Orders:    orders,
		TimeFrom:  api.Time(bookingInfo.TimeFrom.Unix()),
		TimeTo:    api.Time(bookingInfo.TimeTo.Unix()),
		CreatedAt: api.Time(bookingInfo.CreatedAt.Unix()),
		UpdatedAt: api.Time(bookingInfo.UpdatedAt.Unix()),
	}
}

func convertBookingEntity(entity models.BookingEntity) api.BookingEntity {
	return api.BookingEntity{
		ID:        entity.Id,
		Type:      api.BookingEntityType(entity.Type),
		Title:     entity.Title,
		X:         entity.X,
		Y:         entity.Y,
		FloorID:   entity.FloorId,
		Width:     entity.Width,
		Height:    entity.Height,
		Capacity:  entity.Capacity,
		CreatedAt: api.Time(entity.CreatedAt.UTC().Unix()),
		UpdatedAt: api.Time(entity.UpdatedAt.UTC().Unix()),
	}
}

func pointer[T any](v T) *T {
	return &v
}
