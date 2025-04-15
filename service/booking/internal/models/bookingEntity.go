package models

import (
	"time"

	"github.com/google/uuid"
)

type BookingEntityType string

var (
	BookingEntityTypeRoom      BookingEntityType = "ROOM"
	BookingEntityTypeOpenSpace BookingEntityType = "OPEN_SPACE"
)

type BookingEntity struct {
	Id        uuid.UUID         `db:"id"`
	Type      BookingEntityType `db:"type"`
	Title     string            `db:"title"`
	X         int               `db:"x"`
	Y         int               `db:"y"`
	FloorId   uuid.UUID         `db:"floor_id"`
	Width     int               `db:"width"`
	Height    int               `db:"height"`
	Capacity  int               `db:"capacity"`
	CreatedAt time.Time         `db:"created_at"`
	UpdatedAt time.Time         `db:"updated_at"`
}
