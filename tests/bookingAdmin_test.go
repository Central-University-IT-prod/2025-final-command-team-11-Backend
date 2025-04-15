package tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// Тест для получения списка всех бронирований (только для админа)
func TestListAllBookingsAdmin(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем пользователя для тестирования
	_, token := createUser(e)

	// Данные для создания этажа с сущностью (bookingEntity)
	floorData := map[string]interface{}{
		"id":   "550e8400-e29b-41d4-a716-446655440000",
		"name": "Test Floor",
		"entities": []map[string]interface{}{
			{
				"id":       "550e8400-e29b-41d4-a716-446655440001",
				"type":     "ROOM",
				"title":    "Test Room",
				"x":        10,
				"y":        20,
				"width":    100,
				"height":   100,
				"capacity": 10,
				"floor_id": "550e8400-e29b-41d4-a716-446655440000",
			},
		},
	}

	// Создаем этаж с сущностью
	createFloor(e, floorData)
	floorID := floorData["id"].(string)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteFloor(e, floorID) // Удаляем этаж
	})

	// Данные для создания бронирований
	bookingData1 := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001", // ID сущности
		"time_from": 1672502400,                             // Пример времени начала (Unix timestamp)
		"time_to":   1672506000,                             // Пример времени окончания (Unix timestamp)
	}

	bookingData2 := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001", // ID сущности
		"time_from": 1672506000,                             // Пример времени начала (Unix timestamp)
		"time_to":   1672509600,                             // Пример времени окончания (Unix timestamp)
	}

	// Создаем бронирования
	booking1 := createBooking(e, token, bookingData1)
	booking2 := createBooking(e, token, bookingData2)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, booking1["id"].(string))
		deleteBooking(e, token, booking2["id"].(string))
	})

	t.Run("List All Bookings Admin - Success", func(t *testing.T) {
		e.GET("/booking/bookings").
			WithHeader("Authorization", "Bearer "+adminToken).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array().
			NotEmpty().
			Value(0).
			Object().
			ContainsKey("id").
			ContainsKey("user").
			ContainsKey("entity")
	})

	t.Run("List All Bookings Admin - Unauthorized", func(t *testing.T) {
		e.GET("/booking/bookings").
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("List All Bookings Admin - Insufficient Permissions", func(t *testing.T) {
		// Создаем обычного пользователя
		_, userToken := createUser(e)

		e.GET("/booking/bookings").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(http.StatusForbidden)
	})
}

// Тест для получения статистики по бронированиям (только для админа)
func TestGetBookingStatsAdmin(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	t.Run("Get Booking Stats Admin - Success", func(t *testing.T) {
		e.GET("/admin/booking/stats").
			WithHeader("Authorization", "Bearer "+adminToken).
			WithQuery("filter", "day"). // Пример фильтра (день, неделя, месяц)
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("count") // Проверяем, что возвращается статистика
	})

	t.Run("Get Booking Stats Admin - Unauthorized", func(t *testing.T) {
		e.GET("/admin/booking/stats").
			WithQuery("filter", "day").
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("Get Booking Stats Admin - Insufficient Permissions", func(t *testing.T) {
		// Создаем обычного пользователя
		_, userToken := createUser(e)

		e.GET("/admin/booking/stats").
			WithHeader("Authorization", "Bearer "+userToken).
			WithQuery("filter", "day").
			Expect().
			Status(http.StatusForbidden)
	})
}
