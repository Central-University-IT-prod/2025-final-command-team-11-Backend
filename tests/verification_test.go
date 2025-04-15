package tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// Тест для проверки данных верификации
func TestCheckVerification(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем пользователей для тестирования
	_, token := createUser(e)
	verifiedUserID := getUserId(e, token)

	_, token = createUser(e)
	notVerifiedUserID := getUserId(e, token)

	// Устанавливаем данные верификации
	e.POST("/admin/verification/{id}/set", verifiedUserID).
		WithHeader("Authorization", "Bearer "+adminToken).
		WithMultipart().
		WithFileBytes("passport", "passport.jpg", []byte("test data")).
		Expect().
		Status(http.StatusNoContent)

	t.Run("Check Verification - Success", func(t *testing.T) {
		e.GET("/admin/verification/{id}/check", verifiedUserID).
			WithHeader("Authorization", "Bearer "+adminToken).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("verified").
			Value("verified").
			IsEqual(true)

		e.GET("/admin/verification/{id}/check", notVerifiedUserID).
			WithHeader("Authorization", "Bearer "+adminToken).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("verified").
			Value("verified").
			IsEqual(false)
	})

	t.Run("Check Verification - Not Found", func(t *testing.T) {
		nonExistentUserID := "00000000-0000-0000-0000-000000000000"
		e.GET("/admin/verification/{id}/check", nonExistentUserID).
			WithHeader("Authorization", "Bearer "+adminToken).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Check Verification - Unauthorized", func(t *testing.T) {
		e.GET("/admin/verification/{id}/check", verifiedUserID).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

// Тест для установки данных верификации
func TestSetVerification(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Создаем пользователя для тестирования
	_, token := createUser(e)
	userID := getUserId(e, token)

	t.Run("Set Verification - Success", func(t *testing.T) {
		// Устанавливаем данные верификации
		e.POST("/admin/verification/{id}/set", userID).
			WithHeader("Authorization", "Bearer "+adminToken).
			WithMultipart().
			WithFileBytes("passport", "passport.jpg", []byte("test data")).
			Expect().
			Status(http.StatusNoContent)
	})

	t.Run("Set Verification - Invalid Data", func(t *testing.T) {
		// Пытаемся установить данные верификации без файла
		e.POST("/admin/verification/{id}/set", userID).
			WithHeader("Authorization", "Bearer "+adminToken).
			WithMultipart().
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("Set Verification - Not Found", func(t *testing.T) {
		nonExistentUserID := "00000000-0000-0000-0000-000000000000"
		e.POST("/admin/verification/{id}/set", nonExistentUserID).
			WithHeader("Authorization", "Bearer "+adminToken).
			WithMultipart().
			WithFileBytes("passport", "passport.jpg", []byte("test data")).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Set Verification - Unauthorized", func(t *testing.T) {
		e.POST("/admin/verification/{id}/set", userID).
			WithMultipart().
			WithFileBytes("passport", "passport.jpg", []byte("test data")).
			Expect().
			Status(http.StatusUnauthorized)
	})
}
