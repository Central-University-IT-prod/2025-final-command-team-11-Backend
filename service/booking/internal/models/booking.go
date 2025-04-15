package models

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	Id        uuid.UUID `db:"id"`
	EntityId  uuid.UUID `db:"entity_id"`
	UserId    uuid.UUID `db:"user_id"`
	TimeFrom  time.Time `db:"time_from"`
	TimeTo    time.Time `db:"time_to"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type BookingInfo struct {
	Booking
	User   User
	Orders []Order
	Entity BookingEntity
}
