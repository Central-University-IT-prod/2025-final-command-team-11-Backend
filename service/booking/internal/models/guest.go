package models

import (
	"time"

	"github.com/google/uuid"
)

type Guest struct {
	UserId    uuid.UUID `db:"user_id"`
	BookingId uuid.UUID `db:"booking_id"`
	CreatedAt time.Time `db:"created_at"`
}
