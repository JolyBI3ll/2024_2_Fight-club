package controller

//
//import (
//	"2024_2_FIGHT-CLUB/domain"
//	"2024_2_FIGHT-CLUB/internal/auth/mocks"
//	"2024_2_FIGHT-CLUB/internal/service/logger"
//	"2024_2_FIGHT-CLUB/internal/service/middleware"
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/gorilla/mux"
//	"github.com/gorilla/sessions"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"io"
//	"log"
//	"mime/multipart"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//	"time"
//)
//
//func TestRegisterUser(t *testing.T) {
//	if err := logger.InitLoggers(); err != nil {
//		log.Fatalf("Failed to initialize loggers: %v", err)
//	}
//	defer logger.SyncLoggers()
//
//	mockAuthUseCase := &mocks.MockAuthUseCase{}
//	mockSessionService := &mocks.MockServiceSession{}
//	mockJwtTokenService := &mocks.MockJwtTokenService{}
//
//	handler := &AuthHandler{
//		authUseCase:    mockAuthUseCase,
//		sessionService: mockSessionService,
//		jwtToken:       mockJwtTokenService,
//	}
//
//	user := domain.User{
//		Username: "testuser",
//		Email:    "test@example.com",
//		Password: "password123",
//		Name:     "testuser",
//	}
//
//	body, _ := json.Marshal(user)
//	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
//	req.Header.Set("Content-Type", "application/json")
//	w := httptest.NewRecorder()
//
//	// Mock функции
//	mockAuthUseCase.MockRegisterUser = func(ctx context.Context, user *domain.User) error {
//		user.UUID = "test-uuid"
//		return nil
//	}
//
//	mockSession := &sessions.Session{
//		Values: map[interface{}]interface{}{
//			"session_id": "session-id-123",
//		},
//	}
//
//	mockSessionService.MockCreateSession = func(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (*sessions.Session, error) {
//		return mockSession, nil
//	}
//
//	mockJwtTokenService.MockCreate = func(s *sessions.Session, tokenExpTime int64) (string, error) {
//		return "fake-jwt-token", nil
//	}
//
//	// Тест успешного запроса
//	handler.RegisterUser(w, req)
//	res := w.Result()
//	defer res.Body.Close()
//
//	assert.Equal(t, http.StatusCreated, res.StatusCode)
//
//	// Проверка наличия csrf_token в cookies
//	cookies := res.Cookies()
//	foundCsrf := false
//	for _, cookie := range cookies {
//		if cookie.Name == "csrf_token" && cookie.Value == "fake-jwt-token" {
//			foundCsrf = true
//			break
//		}
//	}
//	assert.True(t, foundCsrf, "Expected csrf_token cookie")
//
//	var response map[string]interface{}
//	err := json.NewDecoder(res.Body).Decode(&response)
//	assert.NoError(t, err)
//	assert.Equal(t, "session-id-123", response["session_id"])
//
//	userData := response["user"].(map[string]interface{})
//	assert.Equal(t, "test-uuid", userData["id"])
//	assert.Equal(t, "testuser", userData["username"])
//	assert.Equal(t, "test@example.com", userData["email"])
//
//	// Тест обработки ошибки при регистрации пользователя
//	mockAuthUseCase.MockRegisterUser = func(ctx context.Context, user *domain.User) error {
//		return fmt.Errorf("register error")
//	}
//	w = httptest.NewRecorder()
//	handler.RegisterUser(w, req)
//	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status 500")
//
//	// Тест обработки ошибки при создании сессии
//	mockAuthUseCase.MockRegisterUser = func(ctx context.Context, user *domain.User) error {
//		return nil
//	}
//	mockSessionService.MockCreateSession = func(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (*sessions.Session, error) {
//		return nil, fmt.Errorf("error creating session")
//	}
//	w = httptest.NewRecorder()
//	handler.RegisterUser(w, req)
//	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status 500")
//
//	// Тест обработки ошибки при создании JWT-токена
//	mockSessionService.MockCreateSession = func(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (*sessions.Session, error) {
//		return mockSession, nil
//	}
//	mockJwtTokenService.MockCreate = func(s *sessions.Session, tokenExpTime int64) (string, error) {
//		return "", fmt.Errorf("error creating JWT token")
//	}
//	w = httptest.NewRecorder()
//	handler.RegisterUser(w, req)
//	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status 500")
//}
//
//func TestLoginUser(t *testing.T) {
//	if err := logger.InitLoggers(); err != nil {
//		log.Fatalf("Failed to initialize loggers: %v", err)
//	}
//	defer logger.SyncLoggers()
//	mockAuthUseCase := &mocks.MockAuthUseCase{}
//	mockSessionService := &mocks.MockServiceSession{}
//	mockJwtToken := &mocks.MockJwtTokenService{}
//	handler := NewAuthHandler(mockAuthUseCase, mockSessionService, mockJwtToken)
//
//	user := domain.User{
//		Username: "testuser",
//		Password: "password123",
//	}
//	body, _ := json.Marshal(user)
//	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
//	req.Header.Set("Content-Type", "application/json")
//	w := httptest.NewRecorder()
//
//	// Настройка моков для успешного входа
//	mockAuthUseCase.MockLoginUser = func(ctx context.Context, user *domain.User) (*domain.User, error) {
//		if user.Username == "testuser" && user.Password == "password123" {
//			return &domain.User{
//				UUID:     "test-uuid",
//				Username: "testuser",
//				Email:    "test@example.com",
//			}, nil
//		}
//		return nil, fmt.Errorf("invalid credentials")
//	}
//
//	mockSessionService.MockCreateSession = func(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (*sessions.Session, error) {
//		session := &sessions.Session{
//			Values: map[interface{}]interface{}{
//				"session_id": "session-id-123",
//			},
//		}
//		return session, nil
//	}
//
//	mockJwtToken.MockCreate = func(session *sessions.Session, exp int64) (string, error) {
//		return "fake-jwt-token", nil
//	}
//
//	// Вызов обработчика
//	handler.LoginUser(w, req)
//
//	res := w.Result()
//	defer res.Body.Close()
//
//	// Проверка кода статуса
//	assert.Equal(t, http.StatusOK, res.StatusCode)
//
//	// Проверка на наличие и значение csrf_token cookie
//	cookies := res.Cookies()
//	var csrfTokenCookie *http.Cookie
//	for _, cookie := range cookies {
//		if cookie.Name == "csrf_token" {
//			csrfTokenCookie = cookie
//			break
//		}
//	}
//	require.NotNil(t, csrfTokenCookie, "Expected csrf_token cookie")
//	assert.Equal(t, "fake-jwt-token", csrfTokenCookie.Value)
//
//	// Проверка структуры ответа
//	var response map[string]interface{}
//	err := json.NewDecoder(res.Body).Decode(&response)
//	assert.NoError(t, err)
//
//	assert.Equal(t, "session-id-123", response["session_id"])
//
//	// Проверка данных пользователя в ответе
//	userData := response["user"].(map[string]interface{})
//	assert.Equal(t, "test-uuid", userData["id"])
//	assert.Equal(t, "testuser", userData["username"])
//	assert.Equal(t, "test@example.com", userData["email"])
//
//	// Тестирование на случай неверных учетных данных
//	mockAuthUseCase.MockLoginUser = func(ctx context.Context, user *domain.User) (*domain.User, error) {
//		return nil, fmt.Errorf("invalid credentials")
//	}
//
//	w = httptest.NewRecorder()
//	handler.LoginUser(w, req)
//	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status 500")
//}
//
//func TestLogoutUser(t *testing.T) {
//	if err := logger.InitLoggers(); err != nil {
//		log.Fatalf("Failed to initialize loggers: %v", err)
//	}
//	defer logger.SyncLoggers()
//
//	mockAuthUseCase := &mocks.MockAuthUseCase{}
//	mockSessionService := &mocks.MockServiceSession{}
//	mockJwtToken := &mocks.MockJwtTokenService{}
//	handler := NewAuthHandler(mockAuthUseCase, mockSessionService, mockJwtToken)
//
//	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
//	req.Header.Set("X-CSRF-Token", "Bearer valid-token")
//	w := httptest.NewRecorder()
//
//	// Настройка моков для успешного выхода
//	mockJwtToken.MockValidate = func(tokenString string) (*middleware.JwtCsrfClaims, error) {
//		if tokenString == "valid-token" {
//			return &middleware.JwtCsrfClaims{}, nil
//		}
//		return nil, fmt.Errorf("invalid token")
//	}
//
//	mockSessionService.MockLogoutSession = func(ctx context.Context, r *http.Request, w http.ResponseWriter) error {
//		return nil
//	}
//	handler.LogoutUser(w, req)
//	res := w.Result()
//	defer res.Body.Close()
//	assert.Equal(t, http.StatusOK, res.StatusCode)
//
//	var logoutResponse map[string]string
//	err := json.NewDecoder(res.Body).Decode(&logoutResponse)
//	assert.NoError(t, err)
//	assert.Equal(t, "Logout successfully", logoutResponse["response"])
//
//	cookies := res.Cookies()
//	var csrfTokenCookie *http.Cookie
//	for _, cookie := range cookies {
//		if cookie.Name == "csrf_token" {
//			csrfTokenCookie = cookie
//			break
//		}
//	}
//	require.NotNil(t, csrfTokenCookie, "Expected csrf_token cookie")
//	assert.Empty(t, csrfTokenCookie.Value)
//	assert.True(t, csrfTokenCookie.Expires.Before(time.Now()), "Expected csrf_token cookie to be expired")
//
//	req.Header.Set("X-CSRF-Token", "Bearer invalid-token")
//	w = httptest.NewRecorder()
//
//	handler.LogoutUser(w, req)
//
//	res = w.Result()
//	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
//
//	req.Header.Set("X-CSRF-Token", "Bearer valid-token")
//	mockSessionService.MockLogoutSession = func(ctx context.Context, r *http.Request, w http.ResponseWriter) error {
//		return fmt.Errorf("logout error")
//	}
//	w = httptest.NewRecorder()
//
//	handler.LogoutUser(w, req)
//
//	res = w.Result()
//	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
//}
//
//func TestPutUser(t *testing.T) {
//	// Инициализация логгера
//	if err := logger.InitLoggers(); err != nil {
//		log.Fatalf("Failed to initialize loggers: %v", err)
//	}
//	defer logger.SyncLoggers()
//
//	mockAuthUseCase := &mocks.MockAuthUseCase{}
//	mockSessionService := &mocks.MockServiceSession{}
//	mockJwtToken := &mocks.MockJwtTokenService{}
//	handler := NewAuthHandler(mockAuthUseCase, mockSessionService, mockJwtToken)
//
//	var buf bytes.Buffer
//	writer := multipart.NewWriter(&buf)
//
//	user := domain.User{
//		Username: "updateduser",
//		Email:    "updated@example.com",
//		Sex:      "M",
//		Name:     "Robert Baron",
//	}
//
//	metadata, _ := json.Marshal(user)
//	_ = writer.WriteField("metadata", string(metadata))
//	writer.Close()
//
//	req := httptest.NewRequest("PUT", "/api/users/", &buf)
//	req.Header.Set("Content-Type", writer.FormDataContentType())
//	req.Header.Set("X-CSRF-Token", "Bearer valid-token")
//	w := httptest.NewRecorder()
//
//	mockJwtToken.MockValidate = func(tokenString string) (*middleware.JwtCsrfClaims, error) {
//		if tokenString == "valid-token" {
//			return &middleware.JwtCsrfClaims{}, nil
//		}
//		return nil, fmt.Errorf("invalid token")
//	}
//
//	mockSessionService.MockGetUserID = func(ctx context.Context, r *http.Request) (string, error) {
//		return "user123", nil
//	}
//
//	mockAuthUseCase.MockPutUser = func(ctx context.Context, creds *domain.User, userID string, avatar *multipart.FileHeader) error {
//		return nil
//	}
//
//	handler.PutUser(w, req)
//
//	res := w.Result()
//	assert.Equal(t, http.StatusOK, res.StatusCode)
//
//	mockAuthUseCase.MockPutUser = func(ctx context.Context, creds *domain.User, userID string, avatar *multipart.FileHeader) error {
//		return fmt.Errorf("update error") // Ошибка обновления
//	}
//
//	w = httptest.NewRecorder()
//	handler.PutUser(w, req)
//	res = w.Result()
//	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
//}
//
//func TestGetUserById(t *testing.T) {
//	if err := logger.InitLoggers(); err != nil {
//		t.Fatalf("Failed to initialize loggers: %v", err)
//	}
//	defer logger.SyncLoggers()
//
//	mockAuthUseCase := &mocks.MockAuthUseCase{}
//
//	handler := NewAuthHandler(mockAuthUseCase, nil, nil) // вторым параметром передаем nil, так как MockServiceSession больше не нужен
//
//	router := mux.NewRouter()
//	router.HandleFunc("/api/users/{userId}", handler.GetUserById)
//
//	router.Use(middleware.RequestIDMiddleware)
//
//	// Тестовый пользователь (успешный случай)
//	t.Run("Successful GetUserById", func(t *testing.T) {
//		// Устанавливаем поведение мока для успешного вызова
//		mockAuthUseCase.MockGetUserById = func(ctx context.Context, id string) (*domain.User, error) {
//			return &domain.User{UUID: "test-uuid", Username: "testuser", Email: "test@example.com"}, nil
//		}
//
//		req := httptest.NewRequest("GET", "/api/users/test-uuid", nil)
//		w := httptest.NewRecorder()
//
//		router.ServeHTTP(w, req)
//
//		// Проверяем статус код
//		res := w.Result()
//		assert.Equal(t, http.StatusOK, res.StatusCode)
//
//		// Проверяем заголовок Content-Type
//		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
//
//		// Декодируем ответ
//		var response domain.User
//		err := json.NewDecoder(res.Body).Decode(&response)
//		assert.NoError(t, err)
//
//		// Проверяем содержимое ответа
//		assert.Equal(t, "test-uuid", response.UUID)
//		assert.Equal(t, "testuser", response.Username)
//		assert.Equal(t, "test@example.com", response.Email)
//	})
//
//	// Тестовый пользователь (случай ошибки - пользователь не найден)
//	t.Run("User Not Found", func(t *testing.T) {
//		// Устанавливаем поведение мока для ошибки
//		mockAuthUseCase.MockGetUserById = func(ctx context.Context, id string) (*domain.User, error) {
//			return nil, fmt.Errorf("user not found")
//		}
//
//		// Создаем новый тестовый запрос
//		req := httptest.NewRequest("GET", "/api/users/non-existent-uuid", nil)
//		w := httptest.NewRecorder()
//
//		// Выполняем запрос через маршрутизатор
//		router.ServeHTTP(w, req)
//
//		// Проверяем статус код
//		res := w.Result()
//		assert.Equal(t, http.StatusNotFound, res.StatusCode)
//
//		// Проверяем содержание ответа (можно уточнить в зависимости от реализации h.handleError)
//		var response map[string]string
//		err := json.NewDecoder(res.Body).Decode(&response)
//		assert.NoError(t, err)
//		assert.Contains(t, response["error"], "user not found")
//	})
//}
//
//func TestGetAllUsers(t *testing.T) {
//	if err := logger.InitLoggers(); err != nil {
//		log.Fatalf("Failed to initialize loggers: %v", err)
//	}
//	defer logger.SyncLoggers()
//
//	mockAuthUseCase := &mocks.MockAuthUseCase{}
//	mockSessionService := &mocks.MockServiceSession{}
//	mockJwtToken := &mocks.MockJwtTokenService{}
//	handler := NewAuthHandler(mockAuthUseCase, mockSessionService, mockJwtToken)
//
//	req := httptest.NewRequest("GET", "/users", nil)
//	w := httptest.NewRecorder()
//
//	mockAuthUseCase.MockGetAllUser = func(ctx context.Context) ([]domain.User, error) {
//		return []domain.User{
//			{UUID: "test-uuid-1", Username: "user1", Email: "user1@example.com"},
//			{UUID: "test-uuid-2", Username: "user2", Email: "user2@example.com"},
//		}, nil
//	}
//
//	handler.GetAllUsers(w, req)
//
//	res := w.Result()
//	assert.Equal(t, http.StatusOK, res.StatusCode)
//
//	var response map[string]interface{}
//	err := json.NewDecoder(res.Body).Decode(&response)
//	assert.NoError(t, err)
//
//	usersData := response["users"].([]interface{})
//	assert.Len(t, usersData, 2)
//
//	mockAuthUseCase.MockGetAllUser = func(ctx context.Context) ([]domain.User, error) {
//		return nil, fmt.Errorf("error fetching users") // Ошибка при извлечении пользователей
//	}
//	w = httptest.NewRecorder()
//	handler.GetAllUsers(w, req)
//
//	res = w.Result()
//	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
//}
//
//func TestGetSessionData(t *testing.T) {
//	if err := logger.InitLoggers(); err != nil {
//		log.Fatalf("Failed to initialize loggers: %v", err)
//	}
//	defer logger.SyncLoggers()
//
//	mockAuthUseCase := &mocks.MockAuthUseCase{}
//	mockSessionService := &mocks.MockServiceSession{}
//	mockJwtToken := &mocks.MockJwtTokenService{}
//	handler := NewAuthHandler(mockAuthUseCase, mockSessionService, mockJwtToken)
//
//	req := httptest.NewRequest("GET", "/session", nil)
//	w := httptest.NewRecorder()
//
//	mockSessionService.MockGetSessionData = func(ctx context.Context, r *http.Request) (*map[string]interface{}, error) {
//		return &map[string]interface{}{
//			"id":     "test-uuid",
//			"avatar": "images/avatar.jpg",
//		}, nil
//	}
//
//	handler.GetSessionData(w, req)
//
//	res := w.Result()
//	assert.Equal(t, http.StatusOK, res.StatusCode)
//
//	var response map[string]interface{}
//	err := json.NewDecoder(res.Body).Decode(&response)
//	assert.NoError(t, err)
//
//	assert.Equal(t, "test-uuid", response["id"])
//	assert.Equal(t, "images/avatar.jpg", response["avatar"])
//
//	mockSessionService.MockGetSessionData = func(ctx context.Context, r *http.Request) (*map[string]interface{}, error) {
//		return nil, fmt.Errorf("session error")
//	}
//	w = httptest.NewRecorder()
//	handler.GetSessionData(w, req)
//
//	res = w.Result()
//	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
//}
//
//func TestRefreshCsrfToken(t *testing.T) {
//	if err := logger.InitLoggers(); err != nil {
//		t.Fatalf("Failed to initialize loggers: %v", err)
//	}
//	defer logger.SyncLoggers()
//
//	mockAuthUseCase := &mocks.MockAuthUseCase{}
//	mockSessService := &mocks.MockServiceSession{}
//	mockJWT := &mocks.MockJwtTokenService{}
//
//	handler := NewAuthHandler(mockAuthUseCase, mockSessService, mockJWT)
//
//	router := mux.NewRouter()
//	router.HandleFunc("/api/csrf/refresh", handler.RefreshCsrfToken)
//	router.Use(middleware.RequestIDMiddleware) // Ваш собственный middleware для установки request_id
//
//	testRequestID := "test-request-id"
//
//	router.Use(func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			ctx := context.WithValue(r.Context(), "request_id", testRequestID)
//			next.ServeHTTP(w, r.WithContext(ctx))
//		})
//	})
//
//	newCsrfToken := "new-csrf-token"
//
//	mockSession := &sessions.Session{
//		Values: map[interface{}]interface{}{
//			"session_id": "session-id-123",
//		},
//	}
//
//	t.Run("Successful RefreshCsrfToken", func(t *testing.T) {
//		mockSessService.MockGetSession = func(ctx context.Context, r *http.Request) (*sessions.Session, error) {
//			return mockSession, nil
//		}
//
//		mockJWT.MockCreate = func(s *sessions.Session, tokenExpTime int64) (string, error) {
//			return newCsrfToken, nil
//		}
//
//		req := httptest.NewRequest("POST", "/api/csrf/refresh", nil)
//		w := httptest.NewRecorder()
//
//		router.ServeHTTP(w, req)
//
//		res := w.Result()
//		assert.Equal(t, http.StatusOK, res.StatusCode)
//
//		cookies := res.Cookies()
//		var csrfCookie *http.Cookie
//		for _, cookie := range cookies {
//			if cookie.Name == "csrf_token" {
//				csrfCookie = cookie
//				break
//			}
//		}
//		assert.NotNil(t, csrfCookie)
//		assert.Equal(t, newCsrfToken, csrfCookie.Value)
//		assert.Equal(t, "/", csrfCookie.Path)
//		assert.Equal(t, http.SameSiteStrictMode, csrfCookie.SameSite)
//
//		var response map[string]string
//		err := json.NewDecoder(res.Body).Decode(&response)
//		assert.NoError(t, err)
//		assert.Equal(t, newCsrfToken, response["csrf_token"])
//	})
//
//	t.Run("Failed to GetSession - Unauthorized", func(t *testing.T) {
//		mockSessService.MockGetSession = func(ctx context.Context, r *http.Request) (*sessions.Session, error) {
//			return nil, fmt.Errorf("session not found")
//		}
//		req := httptest.NewRequest("POST", "/api/csrf/refresh", nil)
//		w := httptest.NewRecorder()
//
//		router.ServeHTTP(w, req)
//
//		res := w.Result()
//		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
//
//		// Используем io.ReadAll для чтения тела ответа
//		bodyBytes, err := io.ReadAll(res.Body)
//		assert.NoError(t, err)
//		assert.Equal(t, "Unauthorized\n", string(bodyBytes))
//	})
//
//	t.Run("Failed to Create CSRF Token - Internal Server Error", func(t *testing.T) {
//		mockSessService.MockGetSession = func(ctx context.Context, r *http.Request) (*sessions.Session, error) {
//			return mockSession, nil
//		}
//		mockJWT.MockCreate = func(s *sessions.Session, tokenExpTime int64) (string, error) {
//			return "", fmt.Errorf("failed to create CSRF token")
//		}
//
//		req := httptest.NewRequest("POST", "/api/csrf/refresh", nil)
//		w := httptest.NewRecorder()
//
//		router.ServeHTTP(w, req)
//
//		res := w.Result()
//		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
//
//		// Проверяем тело ответа
//		bodyBytes, err := io.ReadAll(res.Body)
//		assert.NoError(t, err)
//		assert.Equal(t, "Failed to create CSRF token\n", string(bodyBytes))
//	})
//}
