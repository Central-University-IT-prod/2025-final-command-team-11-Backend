package converter

import (
	"REDACTED/team-11/backend/admin/internal/controller/http/v1/dto"
	"REDACTED/team-11/backend/admin/internal/entity"
)

func DtoAccess(status string, booking *entity.Booking) map[string]interface{} {
	data := make(map[string]interface{})

	data["status"] = status

	if status != "NOT_READY" {
		data["booking_id"] = booking.Id
	}

	return data
}

func DtoGuest(guest *entity.Guest) *dto.Guest {
	return &dto.Guest{
		UserId:    guest.UserId,
		BookingId: guest.BookingId,
		CreatedAt: guest.CreatedAt,
	}
}
