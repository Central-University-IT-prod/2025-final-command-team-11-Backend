package tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// Тест для получения списка этажей
func TestGetFloors(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем несколько этажей для тестирования
	floorData1 := map[string]interface{}{
		"id":   "550e8400-e29b-41d4-a716-446655440000",
		"name": "Test Floor 1",
		"entities": []map[string]interface{}{
			{
				"id":       "550e8400-e29b-41d4-a716-446655440001",
				"type":     "ROOM",
				"title":    "Test Room 1",
				"x":        10,
				"y":        20,
				"width":    100,
				"height":   100,
				"capacity": 10,
				"floor_id": "550e8400-e29b-41d4-a716-446655440000",
			},
		},
	}

	floorData2 := map[string]interface{}{
		"id":   "550e8400-e29b-41d4-a716-446655440002",
		"name": "Test Floor 2",
		"entities": []map[string]interface{}{
			{
				"id":       "550e8400-e29b-41d4-a716-446655440003",
				"type":     "ROOM",
				"title":    "Test Room 2",
				"x":        30,
				"y":        40,
				"width":    150,
				"height":   150,
				"capacity": 20,
				"floor_id": "550e8400-e29b-41d4-a716-446655440002",
			},
		},
	}

	// Создаем этажи
	createFloor(e, floorData1)
	createFloor(e, floorData2)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteFloor(e, floorData1["id"].(string))
		deleteFloor(e, floorData2["id"].(string))
	})

	t.Run("Get Floors - Success", func(t *testing.T) {
		resp := e.GET("/admin/layout/floors").
			WithHeader("Authorization", "Bearer "+adminToken).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array()

		resp.Length().IsEqual(2)

		resp.Value(0).Object().
			HasValue("id", floorData1["id"].(string)).
			HasValue("name", "Test Floor 1")

		resp.Value(1).Object().
			HasValue("id", floorData2["id"].(string)).
			HasValue("name", "Test Floor 2")
	})
}

// Тест для получения сущностей для этажа
func TestGetEntitiesForFloor(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем этаж для тестирования
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

	createFloor(e, floorData)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteFloor(e, floorData["id"].(string))
	})

	t.Run("Get Entities for Floor - Success", func(t *testing.T) {
		e.GET("/admin/layout/floors/{id}", floorData["id"].(string)).
			WithHeader("Authorization", "Bearer "+adminToken).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array().
			NotEmpty()
	})

	t.Run("Get Entities for Floor - Not Found", func(t *testing.T) {
		nonExistentFloorID := "00000000-0000-0000-0000-000000000000"
		e.GET("/admin/layout/floors/{id}", nonExistentFloorID).
			WithHeader("Authorization", "Bearer "+adminToken).
			Expect().
			Status(http.StatusNotFound)
	})

}

// Тест для получения сущности по ID
func TestGetEntityByID(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем этаж и сущность для тестирования
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

	createFloor(e, floorData)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteFloor(e, floorData["id"].(string))
	})

	entityID := "550e8400-e29b-41d4-a716-446655440001"

	t.Run("Get Entity by ID - Success", func(t *testing.T) {
		e.GET("/admin/layout/entities/{id}", entityID).
			WithHeader("Authorization", "Bearer "+adminToken).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			HasValue("id", entityID).
			HasValue("type", "ROOM").
			HasValue("title", "Test Room").
			HasValue("x", 10).
			HasValue("y", 20).
			HasValue("width", 100).
			HasValue("height", 100).
			HasValue("capacity", 10)
	})

	t.Run("Get Entity by ID - Not Found", func(t *testing.T) {
		nonExistentEntityID := "00000000-0000-0000-0000-000000000000"
		e.GET("/admin/layout/entities/{id}", nonExistentEntityID).
			WithHeader("Authorization", "Bearer "+adminToken).
			Expect().
			Status(http.StatusNotFound)
	})

}

// Тест для сохранения макета этажа
func TestSaveFloorLayout(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

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

	t.Run("Save Floor Layout - Success", func(t *testing.T) {
		createFloor(e, floorData)

		// Очистка данных после теста
		t.Cleanup(func() {
			deleteFloor(e, floorData["id"].(string))
		})

		e.POST("/admin/layout/floors").
			WithHeader("Authorization", "Bearer "+adminToken).
			WithJSON(floorData).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("Save Floor Layout - Invalid Data", func(t *testing.T) {
		invalidFloorData := map[string]interface{}{
			"id":   "invalid-id",
			"name": "Test Floor",
			"entities": []map[string]interface{}{
				{
					"id":       "invalid-id",
					"type":     "INVALID_TYPE",
					"title":    "Test Room",
					"x":        10,
					"y":        20,
					"width":    100,
					"height":   100,
					"capacity": 10,
					"floor_id": "invalid-id",
				},
			},
		}

		e.POST("/admin/layout/floors").
			WithHeader("Authorization", "Bearer "+adminToken).
			WithJSON(invalidFloorData).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("Save Floor Layout - Unauthorized", func(t *testing.T) {
		e.POST("/admin/layout/floors").
			WithJSON(floorData).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

// Тест для удаления этажа
func TestDeleteFloor(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

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

	t.Run("Delete Floor - Success", func(t *testing.T) {
		createFloor(e, floorData)

		t.Cleanup(func() {
			e.DELETE("/admin/layout/floors/{id}", floorData["id"].(string)).
				WithHeader("Authorization", "Bearer "+adminToken).
				Expect().
				Status(http.StatusNotFound)
		})

		e.DELETE("/admin/layout/floors/{id}", floorData["id"].(string)).
			WithHeader("Authorization", "Bearer "+adminToken).
			Expect().
			Status(http.StatusNoContent)

		e.GET("/admin/layout/floors/{id}", floorData["id"].(string)).
			WithHeader("Authorization", "Bearer "+adminToken).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Delete Floor - Not Found", func(t *testing.T) {
		nonExistentFloorID := "00000000-0000-0000-0000-000000000000"
		e.DELETE("/admin/layout/floors/{id}", nonExistentFloorID).
			WithHeader("Authorization", "Bearer "+adminToken).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Delete Floor - Unauthorized", func(t *testing.T) {
		createFloor(e, floorData)

		// Очистка данных после теста
		t.Cleanup(func() {
			deleteFloor(e, floorData["id"].(string))
		})

		e.DELETE("/admin/layout/floors/{id}", floorData["id"].(string)).
			Expect().
			Status(http.StatusUnauthorized)
	})
}
