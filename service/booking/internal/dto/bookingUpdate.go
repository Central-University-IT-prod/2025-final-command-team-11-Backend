package dto

import (
	"time"

	"github.com/google/uuid"
)

type BookingUpdateDto struct {
	BookingId uuid.UUID
	TimeFrom  *time.Time
	TimeTo    *time.Time
}
