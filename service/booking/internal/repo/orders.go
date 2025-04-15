package repo

import (
	"context"

	"github.com/google/uuid"
	"REDACTED/team-11/backend/booking/internal/dto"
	"REDACTED/team-11/backend/booking/internal/models"
)

type OrdersRepo interface {
	Create(ctx context.Context, input dto.OrderCreateDto) (models.Order, error)
	GetById(ctx context.Context, id uuid.UUID) (models.Order, error)
	GetForBooking(ctx context.Context, bookignId uuid.UUID) ([]models.Order, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
