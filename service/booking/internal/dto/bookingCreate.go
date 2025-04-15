package dto

import (
	"time"

	"github.com/google/uuid"
)

type BookingCreateDto struct {
	EntityId uuid.UUID
	UserId   uuid.UUID
	TimeFrom time.Time
	TimeTo   time.Time
}
