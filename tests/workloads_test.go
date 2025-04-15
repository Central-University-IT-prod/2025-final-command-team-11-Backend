package tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// Тест для получения нагрузки на рабочее место с проверкой данных
func TestGetWorkloadForEntity(t *testing.T) {
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
	entityID := floorData["entities"].([]map[string]interface{})[0]["id"].(string)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteFloor(e, floorID) // Удаляем этаж
	})

	// Данные для создания бронирования
	bookingData := map[string]interface{}{
		"entity_id": entityID, // ID сущности
		"time_from": 0,
		"time_to":   2700,
	}

	// Создаем бронирование
	booking := createBooking(e, token, bookingData)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, booking["id"].(string))
	})

	t.Run("Get Workload for Entity - Success", func(t *testing.T) {
		resp := e.GET("/booking/workloads/{entityId}", entityID).
			WithQuery("timeFrom", 0).
			WithQuery("timeTo", 9000).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array().
			NotEmpty()

		// Проверяем, что в указанное время бронирования рабочее место занято
		item1 := resp.Value(0).Object()
		item1.Value("time").Number().IsEqual(bookingData["time_from"].(int))
		item1.Value("is_free").Boolean().IsFalse()
	})

	t.Run("Get Workload for Entity - Invalid Time Range", func(t *testing.T) {
		e.GET("/booking/workloads/{entityId}", entityID).
			WithQuery("timeFrom", 18000). // timeFrom > timeTo
			WithQuery("timeTo", 17100).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("Get Workload for Entity - Unauthorized", func(t *testing.T) {
		e.GET("/booking/workloads/{entityId}", entityID).
			WithQuery("timeFrom", 27000).
			WithQuery("timeTo", 270000).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

// Тест для получения нагрузки на этаж с проверкой данных
func TestGetWorkloadForFloor(t *testing.T) {
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
	entityID := floorData["entities"].([]map[string]interface{})[0]["id"].(string)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteFloor(e, floorID) // Удаляем этаж
	})

	// Данные для создания бронирования
	bookingData := map[string]interface{}{
		"entity_id": entityID, // ID сущности
		"time_from": 0,
		"time_to":   2700,
	}

	// Создаем бронирование
	booking := createBooking(e, token, bookingData)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, booking["id"].(string))
	})

	t.Run("Get Workload for Floor - Success", func(t *testing.T) {
		resp := e.GET("/booking/workloads/floors/{floorId}", floorID).
			WithQuery("timeFrom", 0).
			WithQuery("timeTo", 9000).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array().
			NotEmpty()

		// Проверяем, что возвращаемые данные соответствуют ожидаемым
		entity := resp.Value(0).Object().Value("entity").Object()
		entity.Value("id").String().IsEqual(entityID)
		entity.Value("type").String().IsEqual("ROOM")
		entity.Value("title").String().IsEqual("Test Room")
		entity.Value("floor_id").String().IsEqual(floorID)
		entity.Value("capacity").Number().IsEqual(10)

		// Проверяем, что в указанное время бронирования рабочее место занято
		resp.Value(0).Object().
			Value("is_free").Boolean().IsFalse()
	})

	t.Run("Get Workload for Floor - Invalid Time Range", func(t *testing.T) {
		e.GET("/booking/workloads/floors/{floorId}", floorID).
			WithQuery("timeFrom", 27000). // timeFrom > timeTo
			WithQuery("timeTo", 26100).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("Get Workload for Floor - Unauthorized", func(t *testing.T) {
		e.GET("/booking/workloads/floors/{floorId}", floorID).
			WithQuery("timeFrom", 54000).
			WithQuery("timeTo", 54900).
			Expect().
			Status(http.StatusUnauthorized)
	})
}
