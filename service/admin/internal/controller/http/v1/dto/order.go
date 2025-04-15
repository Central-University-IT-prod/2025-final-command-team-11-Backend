package dto

import "time"

type Order struct {
	Id        string    `json:"id"`
	BookingId string    `json:"booking_id"`
	Completed bool      `json:"completed"`
	Thing     string    `json:"thing"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	BookName  string    `json:"booking_title"`
}

type Orders struct {
	Values []*Order `json:"orders"`
	Count  int     `json:"count"`
}
