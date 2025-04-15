package dto

import (
	"time"

	types "REDACTED/team-11/backend/admin/internal/entity/type"
)

type FloorEntity struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BookingEntity struct {
	Id        string            `json:"id"`
	Type      types.BookingType `json:"type"`
	Title     string            `json:"title"`
	X         int               `json:"x"`
	Y         int               `json:"y"`
	Width     int               `json:"width"`
	Height    int               `json:"height"`
	Capacity  int               `json:"capacity"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type UpsertFloor struct {
	Id       string   `json:"id" validate:"uuid"`
	Name     string   `json:"name"`
	Entities []Entity `json:"entities"`
}

type Entity struct {
	Id       string            `json:"id"       validate:"uuid"`
	FloorId  string            `json:"floor_id" validate:"uuid"`
	Type     types.BookingType `json:"type"     validate:"booking"`
	Title    string            `json:"title"`
	X        int               `json:"x"`
	Y        int               `json:"y"`
	Width    int               `json:"width"`
	Height   int               `json:"height"`
	Capacity int               `json:"capacity"`
}
