package tests

import (
	"net/http"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
)

const baseURL = "http://localhost:8080/api/v1"

var adminToken = os.Getenv("ADMIN_TOKEN") // Получаем токен администратора из переменной окружения

type JSON map[string]any

// Генерация уникального email с использованием UUID
func generateUniqueEmail() string {
	return "testuser-" + uuid.New().String() + "@example.com"
}

func genRandUser() JSON {
	return JSON{
		"email":    generateUniqueEmail(),
		"name":     gofakeit.Name(),
		"password": "ValidPass123!",
	}
}

func createUser(e *httpexpect.Expect) (JSON, string) {
	user := genRandUser()

	val := e.POST("/id/account/new").WithJSON(user).Expect().Status(http.StatusCreated).JSON().Object().Value("token")
	token := val.String().Raw()
	return user, token
}

func getUserId(e *httpexpect.Expect, token string) string {
	return e.GET("/id/account/").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("id").String().Raw()
}

// Вспомогательная функция для создания этажа
func createFloor(e *httpexpect.Expect, floorData map[string]interface{}) {
	e.POST("/admin/layout/floors").
		WithHeader("Authorization", "Bearer "+adminToken).
		WithJSON(floorData).
		Expect().
		Status(http.StatusOK)
}

// Вспомогательная функция для удаления этажа
func deleteFloor(e *httpexpect.Expect, floorID string) {
	e.DELETE("/admin/layout/floors/{id}", floorID).
		WithHeader("Authorization", "Bearer "+adminToken).
		Expect().
		Status(http.StatusNoContent)
}

// Вспомогательная функция для создания бронирования
func createBooking(e *httpexpect.Expect, token string, bookingData map[string]interface{}) map[string]interface{} {
	var resp map[string]interface{}

	e.POST("/booking/bookings").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(bookingData).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Decode(&resp)

	return resp
}

// Вспомогательная функция для удаления бронирования
func deleteBooking(e *httpexpect.Expect, token string, bookingID string) {
	e.DELETE("/booking/bookings/{bookingId}", bookingID).
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusNoContent)
}

// Вспомогательная функция для создания заказа
func createOrder(e *httpexpect.Expect, token string, bookingID string, orderData map[string]interface{}) map[string]interface{} {
	var resp map[string]interface{}

	e.POST("/booking/bookings/{bookingId}/orders", bookingID).
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(orderData).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Decode(&resp)

	return resp
}

// Вспомогательная функция для удаления заказа
func deleteOrder(e *httpexpect.Expect, token string, bookingID string, orderID string) {
	e.DELETE("/booking/bookings/{bookingId}/orders/{orderId}", bookingID, orderID).
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusNoContent)
}

// Вспомогательная функция для установки статуса заказа как выполненного (для админа)
func setOrderCompleted(e *httpexpect.Expect, token string, orderID string) {
	e.POST("/admin/orders/{id}", orderID).
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusNoContent)
}
