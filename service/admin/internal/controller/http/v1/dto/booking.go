package dto

import "time"

type BookingAccess struct {
	Status    string `json:"status"`
	BookingId string `json:"booking_id"`
}

type GuestId struct {
	UserId string `json:"email" validate:"email"`
}

type Guest struct {
	UserId    string    `json:"email"`
	BookingId string    `json:"booking_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Stats struct {
	Count int `json:"count"`
}
