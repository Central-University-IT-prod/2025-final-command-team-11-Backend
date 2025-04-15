package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id        uuid.UUID `db:"id"`
	BookingId uuid.UUID `db:"booking_id"`
	Completed bool      `db:"completed"`
	Thing     string    `db:"thing"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
