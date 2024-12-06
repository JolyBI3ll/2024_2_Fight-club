package middleware_test

import (
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestIDMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(middleware.RequestIDKey)
		assert.NotEmpty(t, requestID, "Request ID должно быть установлено")

		ctxRequestID := middleware.GetRequestID(r.Context())
		assert.Equal(t, requestID, ctxRequestID, "Request ID в контексте и заголовке должны совпадать")
	})

	middlewareHandler := middleware.RequestIDMiddleware(nextHandler)
	req := httptest.NewRequest("GET", "/api/ads", nil)
	rec := httptest.NewRecorder()

	middlewareHandler.ServeHTTP(rec, req)

	requestID := rec.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID, "Заголовок X-Request-ID должен быть установлен")
}

func TestGetRequestID(t *testing.T) {
	requestID := uuid.New().String()
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, requestID)

	extractedRequestID := middleware.GetRequestID(ctx)
	assert.Equal(t, requestID, extractedRequestID, "Должен извлечь корректный Request ID из контекста")

	emptyCtx := context.Background()
	extractedRequestID = middleware.GetRequestID(emptyCtx)
	assert.Empty(t, extractedRequestID, "Должен вернуть пустую строку, если Request ID не найден")
}

func TestJwtToken_Create(t *testing.T) {
	secret := "mysecretkey"
	jwtService, err := middleware.NewJwtToken(secret)
	assert.NoError(t, err)

	// Создаем сессию с необходимыми значениями
	session := "session123"

	// Задаем время истечения токена (через 1 час)
	expirationTime := time.Now().Add(1 * time.Hour).Unix()

	tokenString, err := jwtService.Create(session, expirationTime)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	//Парсим токен для проверки claims
	claims := &middleware.JwtCsrfClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, jwtService.ParseSecretGetter)
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, "session123", claims.SessionID)
	assert.Equal(t, expirationTime, claims.ExpiresAt)
}

func TestJwtToken_Validate(t *testing.T) {
	secret := "mysecretkey"
	jwtService, err := middleware.NewJwtToken(secret)
	assert.NoError(t, err)

	// Создаем сессию с необходимыми значениями
	session := "session123"
	expirationTime := time.Now().Add(1 * time.Hour).Unix()
	tokenString, err := jwtService.Create(session, expirationTime)
	assert.NoError(t, err)

	// Успешная валидация
	claims, err := jwtService.Validate(tokenString, session)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "session123", claims.SessionID)

	// Валидация с неправильной подписью
	invalidService, _ := middleware.NewJwtToken("wrongsecret")
	_, err = invalidService.Validate(tokenString, session)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token parse error")

	// Валидация истёкшего токена
	expiredTime := time.Now().Add(-1 * time.Hour).Unix()
	expiredToken, err := jwtService.Create(session, expiredTime)

	_, err = jwtService.Validate(expiredToken, session)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token parse error")
	// Валидация токена с неправильными claims
	// Создаем токен с неверными claims
	wrongClaims := middleware.JwtCsrfClaims{
		SessionID: "wrong_session",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			IssuedAt:  time.Now().Unix(),
		},
	}
	wrongToken := jwt.NewWithClaims(jwt.SigningMethodHS256, wrongClaims)
	wrongTokenString, err := wrongToken.SignedString([]byte(secret))
	assert.NoError(t, err)

	validatedClaims, err := jwtService.Validate(wrongTokenString, session)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token invalid")
	assert.Nil(t, validatedClaims)

	//Валидация токена с неправильным session_id
	wrongSession := "wrong_session123"

	_, err = invalidService.Validate(tokenString, wrongSession)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token parse error")
}

func TestJwtToken_ParseSecretGetter(t *testing.T) {
	secret := "mysecretkey"
	jwtService, err := middleware.NewJwtToken(secret)
	assert.NoError(t, err)

	// Создаем токен с правильным методом подписи
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &middleware.JwtCsrfClaims{
		SessionID: "s123",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	// Тест правильного метода подписи
	parsedSecret, err := jwtService.ParseSecretGetter(token)
	assert.NoError(t, err)
	assert.Equal(t, []byte(secret), parsedSecret)

	// Создаем токен с неправильным методом подписи
	token.Method = jwt.SigningMethodRS256
	parsedSecret, err = jwtService.ParseSecretGetter(token)
	assert.Error(t, err)
	assert.Nil(t, parsedSecret)
	assert.Contains(t, err.Error(), "bad sign method")
}
