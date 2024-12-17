package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/utils"
	"2024_2_FIGHT-CLUB/microservices/auth_service/controller/gen"
	"2024_2_FIGHT-CLUB/microservices/auth_service/mocks"
	"bytes"
	"errors"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тест для успешного выполнения RegisterUser
func TestAuthHandler_RegisterUser_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	// Мок gRPC клиента
	mockGrpcClient := new(mocks.MockGrpcClient)
	mockResponse := &gen.UserResponse{
		SessionId: "session123",
		Jwttoken:  "token123",
		User: &gen.User{
			Id:       "test_user_id",
			Username: "test_user_name",
			Email:    "test@example.com",
		},
	}
	mockGrpcClient.On("RegisterUser", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)

	// Инициализация обработчика
	authHandler := AuthHandler{
		client: mockGrpcClient,
	}

	// Создание тела запроса
	user := domain.User{
		Username: "test_user_name",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}
	body, _ := easyjson.Marshal(user)

	// Создание HTTP запроса
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-IP", "127.0.0.1")

	w := httptest.NewRecorder()

	// Вызов обработчика
	authHandler.RegisterUser(w, req)

	// Проверка результата
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusCreated, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
}

// Тест на ошибку при декодировании JSON
func TestAuthHandler_RegisterUser_DecodeError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	authHandler := AuthHandler{}

	// Некорректное тело запроса
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer([]byte("{invalid_json}")))
	w := httptest.NewRecorder()

	// Вызов обработчика
	authHandler.RegisterUser(w, req)

	// Проверка результата
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

// Тест на ошибку от gRPC клиента
func TestAuthHandler_RegisterUser_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	grpcErr := status.Error(codes.Internal, "gRPC Error")
	mockGrpcClient := new(mocks.MockGrpcClient)
	mockGrpcClient.On("RegisterUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, grpcErr)

	authHandler := AuthHandler{
		client: mockGrpcClient,
	}

	user := domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}
	body, _ := easyjson.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	authHandler.RegisterUser(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
}

// Тест на ошибку при конверсии ответа
func TestAuthHandler_RegisterUser_ConvertError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	mockResponse := &gen.UserResponse{
		SessionId: "session123",
		Jwttoken:  "token123",
		User: &gen.User{
			Id:       "test_user_id",
			Username: "test_user_name",
			Email:    "test@example.com",
		},
	}
	mockGrpcClient.On("RegisterUser", mock.Anything, mock.Anything).Return(mockResponse, nil)

	// Переопределяем конвертацию ответа, чтобы она возвращала ошибку
	utils.ConvertAuthResponseProtoToGo = func(protoResponse *gen.RegisterUserResponse, sessionID string) (interface{}, error) {
		return nil, errors.New("conversion error")
	}

	authHandler := AuthHandler{
		client: mockGrpcClient,
	}

	user := domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}
	body, _ := easyjson.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	authHandler.RegisterUser(w, req)

	result := w.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
}
