package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/auth/mocks"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	mockAuthUseCase := &mocks.MockAuthUseCase{}
	mockSessionService := &mocks.MockServiceSession{}

	handler := NewAuthHandler(mockAuthUseCase, mockSessionService)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	user := domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Name:     "testuser",
	}
	metadata, _ := json.Marshal(user)
	_ = writer.WriteField("metadata", string(metadata))
	avatarContent := []byte("fake image content")
	part, _ := writer.CreateFormFile("avatar", "avatar.jpg")
	part.Write(avatarContent)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/auth/register", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	mockAuthUseCase.MockRegisterUser = func(ctx context.Context, user *domain.User, avatar *multipart.FileHeader) error {
		user.UUID = "test-uuid"
		user.Username = "testuser"
		user.Email = "test@example.com"
		user.Password = "password123"
		user.Name = "testuser"
		return nil
	}

	mockSessionService.MockCreateSession = func(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (string, error) {
		return "session-id-123", nil
	}

	handler.RegisterUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "session-id-123", response["session_id"])
	userData := response["user"].(map[string]interface{})
	assert.Equal(t, "test-uuid", userData["id"])
	assert.Equal(t, "testuser", userData["username"])
	assert.Equal(t, "test@example.com", userData["email"])

	mockAuthUseCase.MockRegisterUser = func(ctx context.Context, user *domain.User, avatar *multipart.FileHeader) error {
		return fmt.Errorf("register error")
	}

	mockSessionService.MockCreateSession = func(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (string, error) {
		return "", fmt.Errorf("error creating session")
	}
	w = httptest.NewRecorder()
	handler.RegisterUser(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status 500")
}

func TestLoginUser(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	mockAuthUseCase := &mocks.MockAuthUseCase{}
	mockSessionService := &mocks.MockServiceSession{}
	handler := NewAuthHandler(mockAuthUseCase, mockSessionService)

	user := domain.User{
		Username: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Настройка моков
	mockAuthUseCase.MockLoginUser = func(ctx context.Context, user *domain.User) (*domain.User, error) {
		if user.Username == "testuser" && user.Password == "password123" {
			return &domain.User{
				UUID:     "test-uuid",
				Username: "testuser",
				Email:    "test@example.com",
			}, nil
		}
		return nil, fmt.Errorf("invalid credentials")
	}

	mockSessionService.MockCreateSession = func(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (string, error) {
		return "session-id-123", nil
	}

	handler.LoginUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "session-id-123", response["session_id"])
	userData := response["user"].(map[string]interface{})
	assert.Equal(t, "test-uuid", userData["id"])
	assert.Equal(t, "testuser", userData["username"])
	assert.Equal(t, "test@example.com", userData["email"])

	mockAuthUseCase.MockLoginUser = func(ctx context.Context, user *domain.User) (*domain.User, error) {
		return nil, fmt.Errorf("invalid credentials")
	}

	mockSessionService.MockCreateSession = func(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (string, error) {
		return "", fmt.Errorf("error creating session")
	}
	w = httptest.NewRecorder()
	handler.RegisterUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status 400")
}

func TestLogoutUser(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	mockAuthUseCase := &mocks.MockAuthUseCase{}
	mockSessionService := &mocks.MockServiceSession{}
	handler := NewAuthHandler(mockAuthUseCase, mockSessionService)

	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	mockSessionService.MockLogoutSession = func(ctx context.Context, r *http.Request, w http.ResponseWriter) error {
		return nil
	}

	w = httptest.NewRecorder()

	handler.LogoutUser(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	mockSessionService.MockLogoutSession = func(ctx context.Context, r *http.Request, w http.ResponseWriter) error {
		return fmt.Errorf("logout error") // Возвращаем ошибку
	}

	w = httptest.NewRecorder()

	handler.LogoutUser(w, req)

	res = w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

}

func TestPutUser(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	mockAuthUseCase := &mocks.MockAuthUseCase{}
	mockSessionService := &mocks.MockServiceSession{}
	handler := NewAuthHandler(mockAuthUseCase, mockSessionService)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	user := domain.User{
		UUID:     "test-uuid",
		Username: "updateduser",
		Email:    "updated@example.com",
		Password: "newpassword123",
	}

	metadata, _ := json.Marshal(user)
	_ = writer.WriteField("metadata", string(metadata))
	avatarContent := []byte("fake image content")
	part, _ := writer.CreateFormFile("avatar", "avatar.jpg")
	part.Write(avatarContent)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/putUser", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	mockAuthUseCase.MockPutUser = func(ctx context.Context, creds *domain.User, userID string, avatar *multipart.FileHeader) error {
		return nil
	}

	mockSessionService.MockGetUserID = func(ctx context.Context, r *http.Request, w http.ResponseWriter) (string, error) {
		return "user123", nil
	}

	handler.PutUser(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	mockAuthUseCase.MockPutUser = func(ctx context.Context, creds *domain.User, userID string, avatar *multipart.FileHeader) error {
		return fmt.Errorf("update error") // Ошибка обновления
	}
	w = httptest.NewRecorder()

	handler.PutUser(w, req)

	res = w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestGetUserById(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	mockAuthUseCase := &mocks.MockAuthUseCase{}
	mockSessionService := &mocks.MockServiceSession{}
	handler := NewAuthHandler(mockAuthUseCase, mockSessionService)

	req := httptest.NewRequest("GET", "/user/test-uuid", nil)
	w := httptest.NewRecorder()

	mockSessionService.MockGetUserID = func(ctx context.Context, r *http.Request, w http.ResponseWriter) (string, error) {
		return "user123", nil
	}

	mockAuthUseCase.MockGetUserById = func(ctx context.Context, id string) (*domain.User, error) {
		return &domain.User{UUID: "test-uuid", Username: "testuser", Email: "test@example.com"}, nil
	}

	handler.GetUserById(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "test-uuid", response["UUID"])
	assert.Equal(t, "testuser", response["Username"])
	assert.Equal(t, "test@example.com", response["Email"])

	mockAuthUseCase.MockGetUserById = func(ctx context.Context, id string) (*domain.User, error) {
		return nil, fmt.Errorf("user not found")
	}
	w = httptest.NewRecorder()
	handler.GetUserById(w, req)

	res = w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGetAllUsers(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	mockAuthUseCase := &mocks.MockAuthUseCase{}
	mockSessionService := &mocks.MockServiceSession{}
	handler := NewAuthHandler(mockAuthUseCase, mockSessionService)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	mockAuthUseCase.MockGetAllUser = func(ctx context.Context) ([]domain.User, error) {
		return []domain.User{
			{UUID: "test-uuid-1", Username: "user1", Email: "user1@example.com"},
			{UUID: "test-uuid-2", Username: "user2", Email: "user2@example.com"},
		}, nil
	}

	handler.GetAllUsers(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)

	usersData := response["users"].([]interface{})
	assert.Len(t, usersData, 2)

	mockAuthUseCase.MockGetAllUser = func(ctx context.Context) ([]domain.User, error) {
		return nil, fmt.Errorf("error fetching users") // Ошибка при извлечении пользователей
	}
	w = httptest.NewRecorder()
	handler.GetAllUsers(w, req)

	res = w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestGetSessionData(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	mockAuthUseCase := &mocks.MockAuthUseCase{}
	mockSessionService := &mocks.MockServiceSession{}
	handler := NewAuthHandler(mockAuthUseCase, mockSessionService)

	req := httptest.NewRequest("GET", "/session", nil)
	w := httptest.NewRecorder()

	mockSessionService.MockGetSessionData = func(ctx context.Context, r *http.Request) (*map[string]interface{}, error) {
		return &map[string]interface{}{
			"id":     "test-uuid",
			"avatar": "images/avatar.jpg",
		}, nil
	}

	handler.GetSessionData(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "test-uuid", response["id"])
	assert.Equal(t, "images/avatar.jpg", response["avatar"])

	mockSessionService.MockGetSessionData = func(ctx context.Context, r *http.Request) (*map[string]interface{}, error) {
		return nil, fmt.Errorf("session error")
	}
	w = httptest.NewRecorder()
	handler.GetSessionData(w, req)

	res = w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
