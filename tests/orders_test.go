package tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// Тест для создания заказа
func TestCreateOrder(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем пользователя для тестирования
	_, token := createUser(e)

	// Создаем этаж и сущность для бронирования
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

	// Создаем бронирование для пользователя
	bookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}

	booking := createBooking(e, token, bookingData)
	bookingID := booking["id"].(string)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, bookingID)
	})

	// Данные для создания заказа
	orderData := map[string]interface{}{
		"thing": "laptop",
	}

	t.Run("Create Order - Success", func(t *testing.T) {
		order := createOrder(e, token, bookingID, orderData)

		// Проверяем, что заказ успешно создан
		e.GET("/booking/bookings/{bookingId}/orders", bookingID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array().
			NotEmpty().
			Value(0).
			Object().
			ContainsKey("id").
			ContainsKey("booking_id").
			ContainsKey("thing").
			HasValue("thing", "laptop")

		// Очистка данных после теста
		t.Cleanup(func() {
			deleteOrder(e, token, bookingID, order["id"].(string))
		})
	})

	t.Run("Create Order - Invalid Data", func(t *testing.T) {
		invalidOrderData := map[string]interface{}{
			"thing": "invalid-thing", // Некорректный тип заказа
		}

		e.POST("/booking/bookings/{bookingId}/orders", bookingID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(invalidOrderData).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("Create Order - Unauthorized", func(t *testing.T) {
		e.POST("/booking/bookings/{bookingId}/orders", bookingID).
			WithJSON(orderData).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

// Тест для получения списка заказов
func TestListOrders(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем пользователя для тестирования
	_, token := createUser(e)

	// Создаем этаж и сущность для бронирования
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

	// Создаем бронирование для пользователя
	bookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}

	booking := createBooking(e, token, bookingData)
	bookingID := booking["id"].(string)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, bookingID)
	})

	// Создаем заказы для бронирования
	orderData1 := map[string]interface{}{
		"thing": "laptop",
	}

	orderData2 := map[string]interface{}{
		"thing": "coffee",
	}

	order1 := createOrder(e, token, bookingID, orderData1)
	order2 := createOrder(e, token, bookingID, orderData2)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteOrder(e, token, bookingID, order1["id"].(string))
		deleteOrder(e, token, bookingID, order2["id"].(string))
	})

	t.Run("List Orders - Success", func(t *testing.T) {
		resp := e.GET("/booking/bookings/{bookingId}/orders", bookingID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array()

		resp.Length().IsEqual(2)

		resp.Value(0).Object().
			ContainsKey("id").
			ContainsKey("booking_id").
			ContainsKey("thing")

		resp.Value(1).Object().
			ContainsKey("id").
			ContainsKey("booking_id").
			ContainsKey("thing")
	})

	t.Run("List Orders - Unauthorized", func(t *testing.T) {
		e.GET("/booking/bookings/{bookingId}/orders", bookingID).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

// Тест для удаления заказа
func TestDeleteOrder(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем пользователя для тестирования
	_, token := createUser(e)

	// Создаем этаж и сущность для бронирования
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

	// Создаем бронирование для пользователя
	bookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}

	booking := createBooking(e, token, bookingData)
	bookingID := booking["id"].(string)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, bookingID)
	})

	// Создаем заказ для бронирования
	orderData := map[string]interface{}{
		"thing": "laptop",
	}

	order := createOrder(e, token, bookingID, orderData)
	orderID := order["id"].(string)

	t.Run("Delete Order - Success", func(t *testing.T) {
		e.DELETE("/booking/bookings/{bookingId}/orders/{orderId}", bookingID, orderID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNoContent)
	})

	t.Run("Delete Order - Not Found", func(t *testing.T) {
		nonExistentOrderID := "00000000-0000-0000-0000-000000000000"
		e.DELETE("/booking/bookings/{bookingId}/orders/{orderId}", bookingID, nonExistentOrderID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Delete Order - Unauthorized", func(t *testing.T) {
		e.DELETE("/booking/bookings/{bookingId}/orders/{orderId}", bookingID, orderID).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

// Тест для получения списка всех заказов (для админа)
func TestListAllOrdersAdmin(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем пользователя для тестирования
	_, token := createUser(e)

	// Создаем этаж и сущность для бронирования
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

	// Создаем бронирование для пользователя
	bookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}

	booking := createBooking(e, token, bookingData)
	bookingID := booking["id"].(string)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, bookingID)
	})

	// Создаем заказы для бронирования
	orderData1 := map[string]interface{}{
		"thing": "laptop",
	}

	orderData2 := map[string]interface{}{
		"thing": "coffee",
	}

	order1 := createOrder(e, token, bookingID, orderData1)
	order2 := createOrder(e, token, bookingID, orderData2)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteOrder(e, token, bookingID, order1["id"].(string))
		deleteOrder(e, token, bookingID, order2["id"].(string))
	})

	t.Run("List All Orders - Success", func(t *testing.T) {

		// Проверяем, что ответ содержит список заказов
		e.GET("/admin/orders").
			WithHeader("Authorization", "Bearer "+adminToken).
			WithQuery("page", 1).
			WithQuery("size", 10).
			WithQuery("completed", "").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("count").
			ContainsKey("orders")
	})

	t.Run("List All Orders - Unauthorized", func(t *testing.T) {
		e.GET("/admin/orders").
			WithQuery("page", 1).
			WithQuery("size", 10).
			WithQuery("completed", "").
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("List All Orders - Insufficient Permissions", func(t *testing.T) {
		e.GET("/admin/orders").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("page", 1).
			WithQuery("size", 10).
			WithQuery("completed", "").
			Expect().
			Status(http.StatusForbidden)
	})
}

// Тест для получения статистики по заказам (для админа)
func TestGetOrderStatsAdmin(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	t.Run("Get Order Stats - Success", func(t *testing.T) {

		// Проверяем, что ответ содержит статистику
		e.GET("/admin/orders/stats").
			WithHeader("Authorization", "Bearer "+adminToken).
			WithQuery("filter", "day").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("count")
	})

	t.Run("Get Order Stats - Unauthorized", func(t *testing.T) {
		e.GET("/admin/orders/stats").
			WithQuery("filter", "day").
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("Get Order Stats - Insufficient Permissions", func(t *testing.T) {
		_, token := createUser(e)

		e.GET("/admin/orders/stats").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("filter", "day").
			Expect().
			Status(http.StatusForbidden)
	})
}

// Тест для установки статуса заказа как выполненного (для админа)
func TestSetOrderCompletedAdmin(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем пользователя для тестирования
	_, token := createUser(e)

	// Создаем этаж и сущность для бронирования
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

	// Создаем бронирование для пользователя
	bookingData := map[string]interface{}{
		"entity_id": "550e8400-e29b-41d4-a716-446655440001",
		"time_from": 1672502400,
		"time_to":   1672506000,
	}

	booking := createBooking(e, token, bookingData)
	bookingID := booking["id"].(string)

	// Очистка данных после теста
	t.Cleanup(func() {
		deleteBooking(e, token, bookingID)
	})

	// Создаем заказ для бронирования
	orderData := map[string]interface{}{
		"thing": "laptop",
	}

	order := createOrder(e, token, bookingID, orderData)
	orderID := order["id"].(string)

	t.Run("Set Order Completed - Success", func(t *testing.T) {
		setOrderCompleted(e, adminToken, orderID)

		// Проверяем, что статус заказа изменился
		e.GET("/booking/bookings/{bookingId}/orders", bookingID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array().
			Value(0).
			Object().
			HasValue("completed", true)
	})

	t.Run("Set Order Completed - Unauthorized", func(t *testing.T) {
		e.POST("/admin/orders/{id}", orderID).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("Set Order Completed - Insufficient Permissions", func(t *testing.T) {
		e.POST("/admin/orders/{id}", orderID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusForbidden)
	})
}
