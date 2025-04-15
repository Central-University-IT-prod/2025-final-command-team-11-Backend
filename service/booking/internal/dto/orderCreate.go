package dto

import "github.com/google/uuid"

type OrderCreateDto struct {
	BookingId uuid.UUID
	Thing     string
}
