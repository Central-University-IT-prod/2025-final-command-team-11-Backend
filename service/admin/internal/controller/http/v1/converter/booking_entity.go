package converter

import (
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/dto"
	"REDACTED/team-11/backend/admin/internal/entity"
)

func DtoFloor(floor *entity.FloorEntity) *dto.FloorEntity {
	return &dto.FloorEntity{
		Id:        floor.Id,
		Name:      floor.Name,
		CreatedAt: floor.CreatedAt,
		UpdatedAt: floor.UpdatedAt,
	}
}

func DtoEntity(entity *entity.BookingEntity) *dto.BookingEntity {
	return &dto.BookingEntity{
		Id:        entity.Id,
		Type:      entity.Type,
		Title:     entity.Title,
		X:         entity.X,
		Y:         entity.Y,
		Width:     entity.Width,
		Height:    entity.Height,
		Capacity:  entity.Capacity,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
