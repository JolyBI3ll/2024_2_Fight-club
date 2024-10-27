package middleware_test

import (
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	// Создаем поддельный обработчик для проверки установки ID запроса
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовок запроса
		requestID := r.Context().Value(middleware.RequestIDKey)
		assert.NotEmpty(t, requestID, "Request ID должно быть установлено")

		// Извлекаем Request ID из контекста
		ctxRequestID := middleware.GetRequestID(r.Context())
		assert.Equal(t, requestID, ctxRequestID, "Request ID в контексте и заголовке должны совпадать")
	})

	// Создаем тестовый сервер с использованием middleware
	middlewareHandler := middleware.RequestIDMiddleware(nextHandler)
	req := httptest.NewRequest("GET", "/api/ads", nil)
	rec := httptest.NewRecorder()

	// Выполняем запрос с использованием middleware
	middlewareHandler.ServeHTTP(rec, req)

	// Проверяем заголовок X-Request-ID
	requestID := rec.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID, "Заголовок X-Request-ID должен быть установлен")
}

func TestGetRequestID(t *testing.T) {
	// Создаем новый UUID для проверки
	requestID := uuid.New().String()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, requestID)

	// Проверка правильности извлечения ID
	extractedRequestID := middleware.GetRequestID(ctx)
	assert.Equal(t, requestID, extractedRequestID, "Должен извлечь корректный Request ID из контекста")

	// Проверка случая, когда ID отсутствует
	emptyCtx := context.Background()
	extractedRequestID = middleware.GetRequestID(emptyCtx)
	assert.Empty(t, extractedRequestID, "Должен вернуть пустую строку, если Request ID не найден")
}
