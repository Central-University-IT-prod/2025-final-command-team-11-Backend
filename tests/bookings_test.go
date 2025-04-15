package tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// Тест для создания бронирования
func TestCreateBooking(t *testing.T) {
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

	// Данные для создания бронирования
	bookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001", // ID сущности
		"time_from": 1672502400,                             // Пример времени начала (Unix timestamp)
		"time_to":   1672506000,                             // Пример времени окончания (Unix timestamp)
	}

	t.Run("Create Booking - Success", func(t *testing.T) {
		booking := createBooking(e, token, bookingData)

		// Проверяем, что возвращенные данные соответствуют ожидаемым
		e.GET("/booking/bookings/{bookingId}", booking["id"].(string)).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("id").
			ContainsKey("user").
			ContainsKey("entity").
			HasValue("time_from", bookingData["time_from"]).
			HasValue("time_to", bookingData["time_to"])

		// Очистка данных после теста
		t.Cleanup(func() {
			deleteBooking(e, token, booking["id"].(string))
		})
	})

	t.Run("Create Booking - Invalid Data", func(t *testing.T) {
		invalidBookingData := map[string]interface{}{
			"entity_id": "invalid-id", // Некорректный ID рабочего места
			"time_from": 1672502400,
			"time_to":   1672506000,
		}

		e.POST("/booking/bookings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(invalidBookingData).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("Create Booking - Unauthorized", func(t *testing.T) {
		e.POST("/booking/bookings").
			WithJSON(bookingData).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

// Тест для получения списка всех бронирований (только для админа)
func TestListAllBookings(t *testing.T) {
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

	t.Run("List All Bookings - Success", func(t *testing.T) {
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

	t.Run("List All Bookings - Unauthorized", func(t *testing.T) {
		e.GET("/booking/bookings").
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("List All Bookings - Insufficient Permissions", func(t *testing.T) {
		// Создаем обычного пользователя
		_, token := createUser(e)

		e.GET("/booking/bookings").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusForbidden)
	})
}

// Тест для получения списка моих бронирований
func TestListMyBookings(t *testing.T) {
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

	// Создаем бронирование для пользователя
	bookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}

	booking := createBooking(e, token, bookingData)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, booking["id"].(string))
	})

	t.Run("List My Bookings - Success", func(t *testing.T) {
		e.GET("/booking/bookings/my").
			WithHeader("Authorization", "Bearer "+token).
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

	t.Run("List My Bookings - Unauthorized", func(t *testing.T) {
		e.GET("/booking/bookings/my").
			Expect().
			Status(http.StatusUnauthorized)
	})
}

// Тест для получения бронирования по ID
func TestGetBookingById(t *testing.T) {
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

	// Создаем бронирование для пользователя
	bookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}

	booking := createBooking(e, token, bookingData)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, booking["id"].(string))
	})

	t.Run("Get Booking by ID - Success", func(t *testing.T) {
		e.GET("/booking/bookings/{bookingId}", booking["id"].(string)).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("id").
			ContainsKey("user").
			ContainsKey("entity")

	})

	t.Run("Get Booking by ID - Not Found", func(t *testing.T) {
		nonExistentBookingID := "00000000-0000-0000-0000-000000000000"
		e.GET("/booking/bookings/{bookingId}", nonExistentBookingID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Get Booking by ID - Unauthorized", func(t *testing.T) {
		e.GET("/booking/bookings/{bookingId}", booking["id"].(string)).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

// Тест для обновления бронирования по ID
func TestUpdateBooking(t *testing.T) {
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

	// Создаем бронирование для пользователя
	bookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}

	booking := createBooking(e, token, bookingData)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, booking["id"].(string))
	})

	t.Run("Update Booking - Success", func(t *testing.T) {
		updateData := map[string]interface{}{
			"time_from": 1672506000, // Новое время начала
			"time_to":   1672509600, // Новое время окончания
		}

		e.PATCH("/booking/bookings/{bookingId}", booking["id"].(string)).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(updateData).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("id").
			ContainsKey("user_id").
			ContainsKey("entity_id")

	})

	t.Run("Update Booking - Invalid Data", func(t *testing.T) {
		invalidUpdateData := map[string]interface{}{
			"time_from": "invalid-time", // Некорректное время
			"time_to":   1672509600,
		}

		e.PATCH("/booking/bookings/{bookingId}", booking["id"].(string)).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(invalidUpdateData).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("Update Booking - Unauthorized", func(t *testing.T) {
		updateData := map[string]interface{}{
			"time_from": 1672506000,
			"time_to":   1672509600,
		}

		e.PATCH("/booking/bookings/{bookingId}", booking["id"].(string)).
			WithJSON(updateData).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

// Тест для удаления бронирования по ID
func TestDeleteBooking(t *testing.T) {
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

	// Создаем бронирование для пользователя
	bookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}

	booking := createBooking(e, token, bookingData)

	t.Run("Delete Booking - Success", func(t *testing.T) {
		e.DELETE("/booking/bookings/{bookingId}", booking["id"].(string)).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNoContent)
	})

	t.Run("Delete Booking - Not Found", func(t *testing.T) {
		nonExistentBookingID := "00000000-0000-0000-0000-000000000000"
		e.DELETE("/booking/bookings/{bookingId}", nonExistentBookingID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Delete Booking - Unauthorized", func(t *testing.T) {
		e.DELETE("/booking/bookings/{bookingId}", booking["id"].(string)).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

func TestCreateBooking_InvalidTimeRange(t *testing.T) {
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
		deleteFloor(e, floorID)
	})

	// Данные для создания бронирования с некорректным временем
	invalidBookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672506000, // time_from >= time_to
		"time_to":   1672502400,
	}

	// Пытаемся создать бронирование с некорректным временем
	e.POST("/booking/bookings").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(invalidBookingData).
		Expect().
		Status(http.StatusBadRequest)
}

func TestCreateBooking_OverlappingBooking(t *testing.T) {
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
		deleteFloor(e, floorID)
	})

	// Создаем первое бронирование
	bookingData1 := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}
	booking := createBooking(e, token, bookingData1)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, booking["id"].(string))
	})

	// Пытаемся создать второе бронирование, которое пересекается с первым
	overlappingBookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672503300, // Пересекается с первым бронированием
		"time_to":   1672505100,
	}

	// Пытаемся создать пересекающееся бронирование
	e.POST("/booking/bookings").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(overlappingBookingData).
		Expect().
		Status(http.StatusConflict)
}

func TestCreateBooking_NoAvailableSeats(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем пользователя для тестирования
	_, token := createUser(e)

	// Данные для создания этажа с сущностью типа OPEN_SPACE и маленькой вместимостью
	floorData := map[string]interface{}{
		"id":   "550e8400-e29b-41d4-a716-446655440000",
		"name": "Test Floor",
		"entities": []map[string]interface{}{
			{
				"id":       "550e8400-e29b-41d4-a716-446655440001",
				"type":     "OPEN_SPACE",
				"title":    "Test Open Space",
				"x":        10,
				"y":        20,
				"width":    100,
				"height":   100,
				"capacity": 1, // Очень маленькая вместимость
				"floor_id": "550e8400-e29b-41d4-a716-446655440000",
			},
		},
	}

	// Создаем этаж с сущностью
	createFloor(e, floorData)
	floorID := floorData["id"].(string)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteFloor(e, floorID)
	})

	// Создаем первое бронирование для другого пользователя
	_, otherToken := createUser(e)
	bookingData1 := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}
	booking := createBooking(e, otherToken, bookingData1)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, otherToken, booking["id"].(string))
	})

	// Пытаемся создать бронирование для текущего пользователя, когда все места заняты
	bookingData2 := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}

	e.POST("/booking/bookings").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(bookingData2).
		Expect().
		Status(http.StatusForbidden)
}

func TestUpdateBooking_InvalidTimeRange(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем пользователя для тестирования
	_, token := createUser(e)

	// Данные для создания этажа с сущностью
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
		deleteFloor(e, floorID)
	})

	// Создаем бронирование
	bookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}
	booking := createBooking(e, token, bookingData)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, booking["id"].(string))
	})

	// Пытаемся обновить бронирование с некорректным временем
	invalidUpdateData := map[string]interface{}{
		"time_from": 1672506000, // time_from >= time_to
		"time_to":   1672502400,
	}

	e.PATCH("/booking/bookings/{bookingId}", booking["id"].(string)).
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(invalidUpdateData).
		Expect().
		Status(http.StatusBadRequest)
}

func TestUpdateBooking_OverlappingBooking(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем пользователя для тестирования
	_, token := createUser(e)

	// Данные для создания этажа с сущностью
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
		deleteFloor(e, floorID)
	})

	// Создаем первое бронирование
	bookingData1 := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}
	booking1 := createBooking(e, token, bookingData1)

	// Создаем второе бронирование
	bookingData2 := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672506000,
		"time_to":   1672509600,
	}
	booking2 := createBooking(e, token, bookingData2)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, booking1["id"].(string))
		deleteBooking(e, token, booking2["id"].(string))
	})

	// Пытаемся обновить второе бронирование, чтобы оно пересекалось с первым
	invalidUpdateData := map[string]interface{}{
		"time_from": 1672503300, // Пересекается с первым бронированием
		"time_to":   1672508700,
	}

	e.PATCH("/booking/bookings/{bookingId}", booking2["id"].(string)).
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(invalidUpdateData).
		Expect().
		Status(http.StatusConflict)
}
