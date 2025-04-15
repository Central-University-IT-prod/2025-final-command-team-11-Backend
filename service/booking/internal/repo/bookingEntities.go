package repo

import (
	"context"

	"github.com/google/uuid"
	"REDACTED/team-11/backend/booking/internal/models"
)

type BookingEntitiesRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (models.BookingEntity, error)
	GetForFloor(ctx context.Context, floorId uuid.UUID) ([]models.BookingEntity, error)
}
