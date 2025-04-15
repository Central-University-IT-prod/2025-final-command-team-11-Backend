package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"REDACTED/team-11/backend/booking/internal/dto"
	"REDACTED/team-11/backend/booking/internal/models"
)

type BookingsRepo interface {
	Create(ctx context.Context, input dto.BookingCreateDto) (models.Booking, error)
	GetById(ctx context.Context, id uuid.UUID) (models.Booking, error)
	ListForUser(ctx context.Context, userId uuid.UUID) ([]models.Booking, error)
	Update(ctx context.Context, input dto.BookingUpdateDto) (models.Booking, error)
	Delete(ctx context.Context, id uuid.UUID) error

	ListAll(ctx context.Context) ([]models.Booking, error)

	ListIntersectedForUser(ctx context.Context, userId uuid.UUID, timeFrom, timeTo time.Time) ([]models.Booking, error)
	ListIntersectedForEntity(ctx context.Context, entityId uuid.UUID, timeFrom, timeTo time.Time) ([]models.Booking, error)
}
