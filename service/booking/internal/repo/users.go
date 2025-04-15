package repo

import (
	"context"

	"github.com/google/uuid"
	"REDACTED/team-11/backend/booking/internal/models"
)

type UsersRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (models.User, error)
}
