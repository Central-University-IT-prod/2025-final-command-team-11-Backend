package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
)

func TestUserRegistration(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Успешная регистрация
	t.Run("Successful Registration", func(t *testing.T) {
		email := generateUniqueEmail() // Генерация уникального email
		userData := JSON{
			"email":    email,
			"name":     "Test User",
			"password": "ValidPass123!",
		}

		e.POST("/id/account/new").
			WithJSON(userData).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsKey("token")
	})

	// Регистрация с уже существующим email
	t.Run("Registration with Existing Email", func(t *testing.T) {
		email := generateUniqueEmail() // Генерация уникального email
		userData := JSON{
			"email":    email,
			"name":     "Test User",
			"password": "ValidPass123!",
		}

		// Сначала регистрируем пользователя
		e.POST("/id/account/new").
			WithJSON(userData).
			Expect().
			Status(http.StatusCreated)

		// Пытаемся зарегистрироваться с тем же email
		e.POST("/id/account/new").
			WithJSON(userData).
			Expect().
			Status(http.StatusConflict).
			JSON().
			Object().
			ContainsKey("error")
	})

	// Регистрация с некорректными данными
	t.Run("Registration with Invalid Data", func(t *testing.T) {
		userData := JSON{
			"email":    "invalid-email",
			"name":     "T",
			"password": "short",
		}

		e.POST("/id/account/new").
			WithJSON(userData).
			Expect().
			Status(http.StatusBadRequest).
			JSON().
			Object().
			ContainsKey("error")
	})

	// Регистрация с невалидным паролем
	t.Run("Registration with Invalid Password", func(t *testing.T) {
		email := generateUniqueEmail() // Генерация уникального email
		userData := JSON{
			"email":    email,
			"name":     "Another User",
			"password": "invalid", // Пароль не соответствует требованиям
		}

		e.POST("/id/account/new").
			WithJSON(userData).
			Expect().
			Status(http.StatusBadRequest).
			JSON().
			Object().
			ContainsKey("error")
	})
}

func TestUserLogin(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Успешный вход
	t.Run("Successful Login", func(t *testing.T) {
		email := generateUniqueEmail() // Генерация уникального email
		password := "ValidPass123!"

		// Сначала регистрируем пользователя
		userData := JSON{
			"email":    email,
			"name":     "Test User",
			"password": password,
		}

		e.POST("/id/account/new").
			WithJSON(userData).
			Expect().
			Status(http.StatusCreated)

		// Выполняем вход
		loginData := JSON{
			"email":    email,
			"password": password,
		}

		e.POST("/id/auth/login").
			WithJSON(loginData).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("token")
	})

	// Вход с некорректным email
	t.Run("Login with Invalid Email", func(t *testing.T) {
		loginData := JSON{
			"email":    "nonexistent@example.com",
			"password": "ValidPass123!",
		}

		e.POST("/id/auth/login").
			WithJSON(loginData).
			Expect().
			Status(http.StatusNotFound).
			JSON().
			Object().
			ContainsKey("error")
	})

	// Вход с некорректным паролем
	t.Run("Login with Invalid Password", func(t *testing.T) {
		email := generateUniqueEmail() // Генерация уникального email
		password := "ValidPass123!"

		// Сначала регистрируем пользователя
		userData := JSON{
			"email":    email,
			"name":     "Test User",
			"password": password,
		}

		e.POST("/id/account/new").
			WithJSON(userData).
			Expect().
			Status(http.StatusCreated)

		// Пытаемся войти с неправильным паролем
		loginData := JSON{
			"email":    email,
			"password": "WrongPass123!", // Неверный пароль
		}

		e.POST("/id/auth/login").
			WithJSON(loginData).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().
			Object().
			ContainsKey("error")
	})

	// Вход с невалидным паролем (не соответствует требованиям)
	t.Run("Login with Invalid Password Format", func(t *testing.T) {
		email := generateUniqueEmail() // Генерация уникального email
		password := "ValidPass123!"

		// Сначала регистрируем пользователя
		userData := JSON{
			"email":    email,
			"name":     "Test User",
			"password": password,
		}

		e.POST("/id/account/new").
			WithJSON(userData).
			Expect().
			Status(http.StatusCreated)

		// Пытаемся войти с невалидным паролем
		loginData := JSON{
			"email":    email,
			"password": "invalid", // Пароль не соответствует требованиям
		}

		e.POST("/id/auth/login").
			WithJSON(loginData).
			Expect().
			Status(http.StatusBadRequest).
			JSON().
			Object().
			ContainsKey("error")
	})
}

func TestUserInfo(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Успешное получение информации о пользователе
	t.Run("Get User Info", func(t *testing.T) {
		email := generateUniqueEmail()
		password := "ValidPass123!"

		// Сначала регистрируем пользователя
		userData := JSON{
			"email":    email,
			"name":     "Test User",
			"password": password,
		}

		e.POST("/id/account/new").
			WithJSON(userData).
			Expect().
			Status(http.StatusCreated)

		// Выполняем вход, чтобы получить токен
		loginData := JSON{
			"email":    email,
			"password": password,
		}

		token := e.POST("/id/auth/login").
			WithJSON(loginData).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Value("token").String().Raw()

		// Используем токен для получения информации о пользователе
		e.GET("/id/account/").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			HasValue("email", email).
			HasValue("name", "Test User")
	})

	// Получение информации о несуществующем пользователе
	t.Run("Get Non-Existent User Info", func(t *testing.T) {
		e.GET("/id/account/{userId}", uuid.New().String()).
			Expect().
			Status(http.StatusNotFound).
			JSON().
			Object().
			ContainsKey("error")
	})
}

func TestUserUpdate(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	// Успешное обновление информации о пользователе
	t.Run("Successful Update", func(t *testing.T) {
		email := generateUniqueEmail()
		password := "ValidPass123!"

		// Сначала регистрируем пользователя
		userData := JSON{
			"email":    email,
			"name":     "Test User",
			"password": password,
		}

		e.POST("/id/account/new").
			WithJSON(userData).
			Expect().
			Status(http.StatusCreated)

		// Выполняем вход, чтобы получить токен
		loginData := JSON{
			"email":    email,
			"password": password,
		}

		token := e.POST("/id/auth/login").
			WithJSON(loginData).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Value("token").String().Raw()

		// Обновляем информацию о пользователе
		newEmail := "updated-" + uuid.New().String() + "@example.com"
		newName := "Updated User"
		updateData := JSON{
			"email":       newEmail,
			"name":        newName,
			"oldPassword": password,
			"password":    "NewValidPass123!",
		}

		e.PATCH("/id/account/edit").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(updateData).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			HasValue("email", newEmail).
			HasValue("name", newName)
	})

	// Обновление с некорректными данными
	t.Run("Update with Invalid Data", func(t *testing.T) {
		email := generateUniqueEmail() // Генерация уникального email
		password := "ValidPass123!"

		// Сначала регистрируем пользователя
		userData := JSON{
			"email":    email,
			"name":     "Test User",
			"password": password,
		}

		e.POST("/id/account/new").
			WithJSON(userData).
			Expect().
			Status(http.StatusCreated)

		// Выполняем вход, чтобы получить токен
		loginData := JSON{
			"email":    email,
			"password": password,
		}

		token := e.POST("/id/auth/login").
			WithJSON(loginData).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Value("token").String().Raw()

		// Обновляем информацию с некорректными данными
		updateData := JSON{
			"email":    "invalid-email",
			"name":     "T",
			"password": "short",
		}

		e.PATCH("/id/account/edit").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(updateData).
			Expect().
			Status(http.StatusBadRequest).
			JSON().
			Object().
			ContainsKey("error")
	})

	// Обновление с уже существующим email
	t.Run("Update with Existing Email", func(t *testing.T) {
		email1 := generateUniqueEmail() // Генерация уникального email для первого пользователя
		email2 := generateUniqueEmail() // Генерация уникального email для второго пользователя
		password := "ValidPass123!"

		// Регистрируем первого пользователя
		userData1 := JSON{
			"email":    email1,
			"name":     "User 1",
			"password": password,
		}

		e.POST("/id/account/new").
			WithJSON(userData1).
			Expect().
			Status(http.StatusCreated)

		// Регистрируем второго пользователя
		userData2 := JSON{
			"email":    email2,
			"name":     "User 2",
			"password": password,
		}

		e.POST("/id/account/new").
			WithJSON(userData2).
			Expect().
			Status(http.StatusCreated)

		// Выполняем вход вторым пользователем, чтобы получить токен
		loginData := JSON{
			"email":    email2,
			"password": password,
		}

		token := e.POST("/id/auth/login").
			WithJSON(loginData).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Value("token").String().Raw()

		// Пытаемся обновить email второго пользователя на email первого пользователя
		updateData := JSON{
			"email":    email1, // Уже существующий email
			"name":     "Updated User",
			"password": "NewValidPass123!",
		}

		e.PATCH("/id/account/edit").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(updateData).
			Expect().
			Status(http.StatusConflict).
			JSON().
			Object().
			ContainsKey("error")
	})

	// Обновление с невалидным паролем
	t.Run("Update with Invalid Password", func(t *testing.T) {
		email := generateUniqueEmail() // Генерация уникального email
		password := "ValidPass123!"

		// Сначала регистрируем пользователя
		userData := JSON{
			"email":    email,
			"name":     "Test User",
			"password": password,
		}

		e.POST("/id/account/new").
			WithJSON(userData).
			Expect().
			Status(http.StatusCreated)

		// Выполняем вход, чтобы получить токен
		loginData := JSON{
			"email":    email,
			"password": password,
		}

		token := e.POST("/id/auth/login").
			WithJSON(loginData).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Value("token").String().Raw()

		// Обновляем информацию с невалидным паролем
		updateData := JSON{
			"email":    "updated-" + uuid.New().String() + "@example.com", // Новый уникальный email
			"name":     "Updated User",
			"password": "invalid", // Пароль не соответствует требованиям
		}

		e.PATCH("/id/account/edit").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(updateData).
			Expect().
			Status(http.StatusBadRequest).
			JSON().
			Object().
			ContainsKey("error")
	})
}

func TestLogoutAccount(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	_, token := createUser(e)

	tests := []struct {
		TestName          string
		Token             string
		Status            int
		WithoutAuthHeader bool
		IsError           bool
	}{
		{
			TestName: "Success",
			Token:    fmt.Sprintf("Bearer %s", token),
			Status:   http.StatusOK,
		},
		{
			TestName: "Invalid token",
			Token:    fmt.Sprintf("Bearer %s", "REDACTED"),
			Status:   http.StatusUnauthorized,
			IsError:  true,
		},
		{
			TestName:          "No Authorization header",
			Status:            http.StatusUnauthorized,
			WithoutAuthHeader: true,
			IsError:           true,
		},
		// TODO: Add more test cases
	}

	for _, tc := range tests {
		t.Run(tc.TestName, func(t *testing.T) {
			if tc.IsError {
				if tc.WithoutAuthHeader {
					e.POST("/id/auth/logout").Expect().Status(tc.Status).JSON().Object().ContainsKey("error")
				} else {
					e.POST("/id/auth/logout").WithHeader("Authorization", tc.Token).Expect().Status(tc.Status).JSON().Object().
						ContainsKey("error")
				}
			} else {
				obj := e.POST("/id/auth/logout").WithHeader("Authorization", tc.Token).Expect().Status(tc.Status).JSON().Object()
				obj.ContainsKey("message")
			}
		})
	}
}

func TestRefresh(t *testing.T) {
	e := httpexpect.Default(t, baseURL)

	t.Run("Without cookie", func(t *testing.T) {
		e.GET("/id/auth/refresh").Expect().Status(http.StatusBadRequest).
			JSON().Object().ContainsKey("error")
	})

	user := genRandUser()

	val := e.POST("/id/account/new").WithJSON(user).Expect().Status(http.StatusCreated).Cookie("refreshToken").Value()
	token := val.Raw()

	test := []struct {
		TestName      string
		RefreshToken  string
		Status        int
		WithoutCookie bool
		IsError       bool
	}{
		{
			TestName:     "Success",
			Status:       http.StatusOK,
			RefreshToken: token,
		},
		{
			TestName:     "Invalid token",
			Status:       http.StatusUnauthorized,
			RefreshToken: "REDACTED",
			IsError:      true,
		},
		{
			TestName: "Empty refresh token",
			Status:   http.StatusUnauthorized,
			IsError:  true,
		},
	}

	for _, tc := range test {
		t.Run(tc.TestName, func(t *testing.T) {
			if tc.IsError {
				e.GET("/id/auth/refresh").WithCookie("refreshToken", tc.RefreshToken).Expect().Status(tc.Status).
					JSON().Object().ContainsKey("error")
			} else {
				e.GET("/id/auth/refresh").WithCookie("refreshToken", tc.RefreshToken).Expect().Status(tc.Status).
					JSON().Object().ContainsKey("token")
			}
		})
	}

}
