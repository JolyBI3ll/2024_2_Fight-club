package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/utils"
	"2024_2_FIGHT-CLUB/microservices/auth_service/controller/gen"
	"2024_2_FIGHT-CLUB/microservices/auth_service/mocks"
	"bytes"
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// GenerateImage создает изображение с указанным форматом и размером
func GenerateImage(format string, width, height int) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Заливаем белым цветом
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, white)
		}
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err := jpeg.Encode(&buf, img, nil)
		if err != nil {
			return nil, err
		}
	case "png":
		err := png.Encode(&buf, img)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported format")
	}

	return buf.Bytes(), nil
}

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
	utilsMock := &utils.MockUtils{}
	mockResponse := &gen.UserResponse{
		SessionId: "session123",
		Jwttoken:  "token123",
		User: &gen.User{
			Id:       "test_user_id",
			Username: "test_user_name",
			Email:    "test@example.com",
		},
	}
	userSession := "session123"
	mockGrpcClient.On("RegisterUser", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)
	utilsMock.On("ConvertAuthResponseProtoToGo", mock.Anything, userSession).Return(domain.AuthResponse{
		SessionId: userSession,
		User: domain.AuthData{
			Id:       mockResponse.User.Id,
			Email:    mockResponse.User.Email,
			Username: mockResponse.User.Username,
		},
	}, nil)
	// Инициализация обработчика
	authHandler := AuthHandler{
		client: mockGrpcClient,
		utils:  utilsMock,
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

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
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
	mockGrpcClient.On("RegisterUser", mock.Anything, mock.Anything, mock.Anything).Return(&gen.UserResponse{}, grpcErr)

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
	mockResponse := &gen.UserResponse{
		SessionId: "session123",
		Jwttoken:  "token123",
		User: &gen.User{
			Id:       "test_user_id",
			Username: "test_user_name",
			Email:    "test@example.com",
		},
	}
	userSession := "session123"
	utilsMock := &utils.MockUtils{}
	mockGrpcClient := new(mocks.MockGrpcClient)
	authHandler := AuthHandler{
		client: mockGrpcClient,
		utils:  utilsMock,
	}
	mockGrpcClient.On("RegisterUser", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)
	// Переопределяем конвертацию ответа, чтобы она возвращала ошибку
	utilsMock.On("ConvertAuthResponseProtoToGo", mockResponse, userSession).
		Return(nil, errors.New("conversion error"))

	user := domain.User{
		Username: "test_user_name",
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

func TestAuthHandler_LoginUser_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	// Моки
	mockGrpcClient := new(mocks.MockGrpcClient)
	utilsMock := &utils.MockUtils{}

	mockResponse := &gen.UserResponse{
		SessionId: "session123",
		Jwttoken:  "token123",
		User: &gen.User{
			Id:       "test_user_id",
			Username: "test_user_name",
			Email:    "test@example.com",
		},
	}

	// Настройка моков
	mockGrpcClient.On("LoginUser", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)
	utilsMock.On("ConvertAuthResponseProtoToGo", mockResponse, "session123").Return(domain.AuthResponse{
		SessionId: "session123",
		User: domain.AuthData{
			Id:       "test_user_id",
			Email:    "test@example.com",
			Username: "test_user_name",
		},
	}, nil)

	// Инициализация обработчика
	authHandler := AuthHandler{
		client: mockGrpcClient,
		utils:  utilsMock,
	}

	// Тело запроса
	user := domain.User{
		Username: "test_user",
		Password: "password123",
	}
	body, _ := easyjson.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-IP", "127.0.0.1")
	w := httptest.NewRecorder()

	// Вызов метода
	authHandler.LoginUser(w, req)

	// Проверка
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)
}

func TestAuthHandler_LoginUser_DecodeError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	authHandler := AuthHandler{}

	// Некорректное тело запроса
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte("invalid_json")))
	w := httptest.NewRecorder()

	authHandler.LoginUser(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestAuthHandler_LoginUser_ExistingCsrfToken(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	authHandler := AuthHandler{}

	user := domain.User{
		Username: "test_user",
		Password: "password123",
	}
	body, _ := easyjson.Marshal(user)

	// Запрос с существующим CSRF куки
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "already_exists"})
	w := httptest.NewRecorder()

	authHandler.LoginUser(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestAuthHandler_LoginUser_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	grpcErr := status.Error(codes.Internal, "gRPC Error")

	// Настройка мока
	mockGrpcClient.On("LoginUser", mock.Anything, mock.Anything, mock.Anything).Return(&gen.UserResponse{}, grpcErr)

	authHandler := AuthHandler{
		client: mockGrpcClient,
	}

	user := domain.User{
		Username: "test_user",
		Password: "password123",
	}
	body, _ := easyjson.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.LoginUser(w, req)

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

func TestAuthHandler_LoginUser_ConversionError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	utilsMock := &utils.MockUtils{}

	mockResponse := &gen.UserResponse{
		SessionId: "session123",
		Jwttoken:  "token123",
		User: &gen.User{
			Id:       "test_user_id",
			Username: "test_user",
			Email:    "test@example.com",
		},
	}

	// Настройка моков
	mockGrpcClient.On("LoginUser", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)
	utilsMock.On("ConvertAuthResponseProtoToGo", mockResponse, "session123").
		Return(domain.AuthResponse{}, errors.New("conversion error"))

	authHandler := AuthHandler{
		client: mockGrpcClient,
		utils:  utilsMock,
	}

	user := domain.User{
		Username: "test_user",
		Password: "password123",
	}
	body, _ := easyjson.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.LoginUser(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)
}

func TestAuthHandler_LogoutUser_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() { _ = logger.SyncLoggers() }()

	mockGrpcClient := new(mocks.MockGrpcClient)
	mockSession := new(mocks.MockServiceSession)
	mockResponse := &gen.LogoutUserResponse{}

	// Настройка моков
	mockSession.MockLogoutSession = func(ctx context.Context, sessionID string) error {
		return nil
	}
	mockGrpcClient.On("LogoutUser", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)

	authHandler := AuthHandler{
		client:         mockGrpcClient,
		sessionService: mockSession,
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	// Установка куков в запрос
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test_session_id",
	})
	req.AddCookie(&http.Cookie{
		Name:  "csrf_token",
		Value: "test_csrf_token",
	})

	// Запись ответа
	w := httptest.NewRecorder()

	authHandler.LogoutUser(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)

	cookies := result.Cookies()
	require.Len(t, cookies, 2)
	assert.Equal(t, "session_id", cookies[0].Name)
	assert.Equal(t, "csrf_token", cookies[1].Name)
	assert.Equal(t, "", cookies[0].Value)
	assert.Equal(t, "", cookies[1].Value)
}

func TestAuthHandler_LogoutUser_SessionIDError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() { _ = logger.SyncLoggers() }()

	mockSession := new(mocks.MockServiceSession)

	authHandler := AuthHandler{
		sessionService: mockSession,
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	w := httptest.NewRecorder()

	authHandler.LogoutUser(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
}

func TestAuthHandler_LogoutUser_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() { _ = logger.SyncLoggers() }()

	mockGrpcClient := new(mocks.MockGrpcClient)
	mockSession := new(mocks.MockServiceSession)
	grpcErr := status.Error(codes.Internal, "gRPC Logout Error")

	mockSession.MockLogoutSession = func(ctx context.Context, sessionID string) error {
		return errors.New("session error")
	}
	mockGrpcClient.On("LogoutUser", mock.Anything, mock.Anything, mock.Anything).Return(&gen.LogoutUserResponse{}, grpcErr)

	authHandler := AuthHandler{
		client:         mockGrpcClient,
		sessionService: mockSession,
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	// Установка куков в запрос
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test_session_id",
	})
	req.AddCookie(&http.Cookie{
		Name:  "csrf_token",
		Value: "test_csrf_token",
	})
	w := httptest.NewRecorder()

	authHandler.LogoutUser(w, req)

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

func TestAuthHandler_PutUser_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() { _ = logger.SyncLoggers() }()

	// Моки
	mockGrpcClient := new(mocks.MockGrpcClient)
	mockSession := new(mocks.MockServiceSession)

	mockGrpcClient.On("PutUser", mock.Anything, mock.Anything, mock.Anything).Return(&gen.UpdateResponse{}, nil)

	authHandler := AuthHandler{
		client:         mockGrpcClient,
		sessionService: mockSession,
	}

	// Тело запроса
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("metadata", `{
		"UUID": "123",
		"Username": "test_user",
		"Password": "password123",
		"Email": "test@example.com",
		"Name": "Test User",
		"Score": 100,
		"Sex": "M",
		"GuestCount": 2,
		"Birthdate": "2000-01-01T00:00:00Z",
		"IsHost": true
	}`)

	part, _ := writer.CreateFormFile("avatar", "avatar.png")
	_, err := part.Write([]byte("fake_image_data"))
	if err != nil {
		return
	}

	err = writer.Close()
	if err != nil {
		return
	}

	req := httptest.NewRequest(http.MethodPut, "/user", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-CSRF-Token", "csrf_token")
	req.Header.Set("X-Real-IP", "127.0.0.1")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test_session_id",
	})
	w := httptest.NewRecorder()

	authHandler.PutUser(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
}

func TestAuthHandler_PutUser_SessionIDError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() { _ = logger.SyncLoggers() }()

	mockSession := new(mocks.MockServiceSession)

	authHandler := AuthHandler{
		sessionService: mockSession,
	}

	req := httptest.NewRequest(http.MethodPut, "/user", nil)
	w := httptest.NewRecorder()

	authHandler.PutUser(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
}

func TestAuthHandler_PutUser_ParseMultipartFormError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() { _ = logger.SyncLoggers() }()

	mockSession := new(mocks.MockServiceSession)

	authHandler := AuthHandler{
		sessionService: mockSession,
	}

	req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer([]byte("invalid_data")))
	req.Header.Set("Content-Type", "multipart/form-data")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test_session_id",
	})
	w := httptest.NewRecorder()

	authHandler.PutUser(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestAuthHandler_PutUser_MetadataParseError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() { _ = logger.SyncLoggers() }()

	mockSession := new(mocks.MockServiceSession)

	authHandler := AuthHandler{
		sessionService: mockSession,
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("metadata", "invalid_json")
	err := writer.Close()
	if err != nil {
		return
	}

	req := httptest.NewRequest(http.MethodPut, "/user", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test_session_id",
	})
	w := httptest.NewRecorder()

	authHandler.PutUser(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestAuthHandler_PutUser_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() { _ = logger.SyncLoggers() }()

	mockGrpcClient := new(mocks.MockGrpcClient)
	mockSession := new(mocks.MockServiceSession)
	grpcErr := status.Error(codes.Internal, "gRPC Error")
	mockGrpcClient.On("PutUser", mock.Anything, mock.Anything, mock.Anything).Return(&gen.UpdateResponse{}, grpcErr)

	authHandler := AuthHandler{
		client:         mockGrpcClient,
		sessionService: mockSession,
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("metadata", `{"Username": "test_user"}`)
	err := writer.Close()
	if err != nil {
		return
	}

	req := httptest.NewRequest(http.MethodPut, "/user", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test_session_id",
	})
	w := httptest.NewRecorder()

	authHandler.PutUser(w, req)

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

func TestAuthHandler_GetUserById_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	utilsMock := &utils.MockUtils{}

	mockResponse := &gen.GetUserByIdResponse{
		User: &gen.MetadataOneUser{
			Uuid:       "test_uuid",
			Username:   "test_username",
			Email:      "test@example.com",
			Name:       "test_name",
			Score:      2,
			Avatar:     "avatar.png",
			Sex:        "M",
			GuestCount: 5,
			Birthdate:  nil,
			IsHost:     false,
		},
	}

	// Настройка моков
	mockGrpcClient.On("GetUserById", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)
	utilsMock.On("ConvertUserResponseProtoToGo", mockResponse.User).Return(domain.UserDataResponse{
		Uuid:       "test_uuid",
		Username:   "test_username",
		Email:      "test@example.com",
		Name:       "test_name",
		Score:      2,
		Avatar:     "avatar.png",
		Sex:        "M",
		GuestCount: 5,
		Birthdate:  time.Time{},
		IsHost:     false,
	}, nil)

	authHandler := AuthHandler{
		client: mockGrpcClient,
		utils:  utilsMock,
	}

	req := httptest.NewRequest(http.MethodGet, "/user/user123", nil)
	req = mux.SetURLVars(req, map[string]string{"userId": "user123"})
	req.Header.Set("X-Real-IP", "127.0.0.1")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test_session_id",
	})
	w := httptest.NewRecorder()

	authHandler.GetUserById(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)
}

func TestAuthHandler_GetUserById_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)

	grpcErr := status.Error(codes.Internal, "gRPC error")
	mockGrpcClient.On("GetUserById", mock.Anything, mock.Anything, mock.Anything).Return(&gen.GetUserByIdResponse{}, grpcErr)

	authHandler := AuthHandler{
		client: mockGrpcClient,
	}

	req := httptest.NewRequest(http.MethodGet, "/user/user123", nil)
	req = mux.SetURLVars(req, map[string]string{"userId": "user123"})
	w := httptest.NewRecorder()

	authHandler.GetUserById(w, req)

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

func TestAuthHandler_GetUserById_ConversionError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	utilsMock := &utils.MockUtils{}

	mockResponse := &gen.GetUserByIdResponse{
		User: &gen.MetadataOneUser{
			Uuid:       "test_uuid",
			Username:   "test_username",
			Email:      "test@example.com",
			Name:       "test_name",
			Score:      2,
			Avatar:     "avatar.png",
			Sex:        "M",
			GuestCount: 5,
			Birthdate:  nil,
			IsHost:     false,
		},
	}

	// Моки
	mockGrpcClient.On("GetUserById", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)
	utilsMock.On("ConvertUserResponseProtoToGo", mockResponse.User).
		Return(domain.UserDataResponse{}, errors.New("conversion error"))

	authHandler := AuthHandler{
		client: mockGrpcClient,
		utils:  utilsMock,
	}

	req := httptest.NewRequest(http.MethodGet, "/user/user123", nil)
	req = mux.SetURLVars(req, map[string]string{"userId": "user123"})
	w := httptest.NewRecorder()

	authHandler.GetUserById(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)
}

func TestAuthHandler_GetAllUsers_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	// Моки
	mockGrpcClient := new(mocks.MockGrpcClient)
	utilsMock := &utils.MockUtils{}

	mockUsers := &gen.AllUsersResponse{
		Users: []*gen.MetadataOneUser{
			{
				Uuid:       "test_uuid1",
				Username:   "test_username1",
				Email:      "test1@example.com",
				Name:       "test_name1",
				Score:      2,
				Avatar:     "avatar1.png",
				Sex:        "M",
				GuestCount: 5,
				Birthdate:  nil,
				IsHost:     false,
			},
			{
				Uuid:       "test_uuid2",
				Username:   "test_username2",
				Email:      "test2@example.com",
				Name:       "test_name2",
				Score:      3,
				Avatar:     "avatar2.png",
				Sex:        "F",
				GuestCount: 2,
				Birthdate:  nil,
				IsHost:     false,
			},
		},
	}

	convertedUsers := []*domain.UserDataResponse{
		{
			Uuid:       "test_uuid1",
			Username:   "test_username1",
			Email:      "test1@example.com",
			Name:       "test_name1",
			Score:      2,
			Avatar:     "avatar1.png",
			Sex:        "M",
			GuestCount: 5,
			Birthdate:  time.Time{},
			IsHost:     false,
		},
		{
			Uuid:       "test_uuid2",
			Username:   "test_username2",
			Email:      "test2@example.com",
			Name:       "test_name2",
			Score:      3,
			Avatar:     "avatar2.png",
			Sex:        "F",
			GuestCount: 2,
			Birthdate:  time.Time{},
			IsHost:     false,
		},
	}

	mockGrpcClient.On("GetAllUsers", mock.Anything, mock.Anything, mock.Anything).Return(mockUsers, nil)
	utilsMock.On("ConvertUsersProtoToGo", mockUsers).Return(convertedUsers, nil)

	authHandler := AuthHandler{
		client: mockGrpcClient,
		utils:  utilsMock,
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	w := httptest.NewRecorder()

	// Вызов метода
	authHandler.GetAllUsers(w, req)

	// Проверка
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)

	body, err := io.ReadAll(result.Body)
	require.NoError(t, err)

	var response domain.GetAllUsersResponse
	err = easyjson.Unmarshal(body, &response)
	require.NoError(t, err)
	assert.Equal(t, convertedUsers, response.Users)
}

func TestAuthHandler_GetAllUsers_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	grpcErr := status.Error(codes.Internal, "gRPC Error")

	mockGrpcClient.On("GetAllUsers", mock.Anything, mock.Anything, mock.Anything).Return(&gen.AllUsersResponse{}, grpcErr)

	authHandler := AuthHandler{
		client: mockGrpcClient,
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	// Вызов метода
	authHandler.GetAllUsers(w, req)

	// Проверка
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

func TestAuthHandler_GetAllUsers_ConversionError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	utilsMock := &utils.MockUtils{}

	mockUsers := &gen.AllUsersResponse{
		Users: []*gen.MetadataOneUser{
			{
				Uuid:       "test_uuid1",
				Username:   "test_username1",
				Email:      "test1@example.com",
				Name:       "test_name1",
				Score:      2,
				Avatar:     "avatar1.png",
				Sex:        "M",
				GuestCount: 5,
				Birthdate:  nil,
				IsHost:     false,
			},
			{
				Uuid:       "test_uuid2",
				Username:   "test_username2",
				Email:      "test2@example.com",
				Name:       "test_name2",
				Score:      3,
				Avatar:     "avatar2.png",
				Sex:        "F",
				GuestCount: 2,
				Birthdate:  nil,
				IsHost:     false,
			},
		},
	}

	mockGrpcClient.On("GetAllUsers", mock.Anything, mock.Anything, mock.Anything).Return(mockUsers, nil)
	utilsMock.On("ConvertUsersProtoToGo", mockUsers).Return(nil, errors.New("conversion error"))

	authHandler := AuthHandler{
		client: mockGrpcClient,
		utils:  utilsMock,
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	// Вызов метода
	authHandler.GetAllUsers(w, req)

	// Проверка
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)
}

func TestAuthHandler_GetSessionData_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	// Моки
	mockGrpcClient := new(mocks.MockGrpcClient)
	utilsMock := &utils.MockUtils{}

	sessionID := "test_session_id"
	mockResponse := &gen.SessionDataResponse{
		Id:     "test_id",
		Avatar: "test_avatar",
	}

	convertedResponse := domain.SessionData{
		Id:     "test_id",
		Avatar: "test_avatar",
	}

	// Настройка моков
	mockGrpcClient.On("GetSessionData", mock.Anything, &gen.GetSessionDataRequest{SessionId: sessionID}, mock.Anything).Return(mockResponse, nil)
	utilsMock.On("ConvertSessionDataProtoToGo", mockResponse).Return(convertedResponse, nil)
	sessionMock := &mocks.MockServiceSession{}

	authHandler := AuthHandler{
		client:         mockGrpcClient,
		utils:          utilsMock,
		sessionService: sessionMock,
	}

	// Запрос
	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test_session_id",
	})
	w := httptest.NewRecorder()

	// Вызов метода
	authHandler.GetSessionData(w, req)

	// Проверка
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)

	var response domain.SessionData
	err := easyjson.UnmarshalFromReader(result.Body, &response)
	require.NoError(t, err)
	assert.Equal(t, convertedResponse, response)

	mockGrpcClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)
}

func TestAuthHandler_GetSessionData_SessionError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	sessionMock := &mocks.MockServiceSession{}

	authHandler := AuthHandler{
		sessionService: sessionMock,
	}

	// Запрос
	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	w := httptest.NewRecorder()

	// Вызов метода
	authHandler.GetSessionData(w, req)

	// Проверка
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
}

func TestAuthHandler_GetSessionData_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	sessionMock := &mocks.MockServiceSession{}
	sessionID := "test_session_id"
	grpcErr := status.Error(codes.Internal, "gRPC error")

	// Настройка моков
	mockGrpcClient.On("GetSessionData", mock.Anything, &gen.GetSessionDataRequest{SessionId: sessionID}, mock.Anything).Return(&gen.SessionDataResponse{}, grpcErr)

	authHandler := AuthHandler{
		client:         mockGrpcClient,
		sessionService: sessionMock,
	}

	// Запрос
	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test_session_id",
	})
	w := httptest.NewRecorder()

	// Вызов метода
	authHandler.GetSessionData(w, req)

	// Проверка
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

func TestAuthHandler_GetSessionData_ConversionError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	utilsMock := &utils.MockUtils{}
	sessionMock := &mocks.MockServiceSession{}

	sessionID := "test_session_id"
	mockResponse := &gen.SessionDataResponse{
		Id:     "test_id",
		Avatar: "test_avatar",
	}

	// Настройка моков
	mockGrpcClient.On("GetSessionData", mock.Anything, &gen.GetSessionDataRequest{SessionId: sessionID}, mock.Anything).Return(mockResponse, nil)
	utilsMock.On("ConvertSessionDataProtoToGo", mockResponse).Return(nil, errors.New("conversion error"))

	authHandler := AuthHandler{
		client:         mockGrpcClient,
		utils:          utilsMock,
		sessionService: sessionMock,
	}

	// Запрос
	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test_session_id",
	})
	w := httptest.NewRecorder()

	// Вызов метода
	authHandler.GetSessionData(w, req)

	// Проверка
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

	mockGrpcClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)
}

func TestAuthHandler_RefreshCsrfToken_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	// Моки
	mockGrpcClient := new(mocks.MockGrpcClient)
	mockResponse := &gen.RefreshCsrfTokenResponse{
		CsrfToken: "new_csrf_token",
	}

	// Настройка мока
	mockGrpcClient.On("RefreshCsrfToken", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)

	authHandler := AuthHandler{
		client: mockGrpcClient,
	}

	// Запрос с валидным sessionID
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh-csrf", nil)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test_session_id",
	})
	w := httptest.NewRecorder()

	// Вызов метода
	authHandler.RefreshCsrfToken(w, req)

	// Проверка
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)

	// Проверяем CSRF куки
	cookies := result.Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, "csrf_token", cookies[0].Name)
	assert.Equal(t, "new_csrf_token", cookies[0].Value)

	mockGrpcClient.AssertExpectations(t)
}

func TestAuthHandler_RefreshCsrfToken_NoSessionID(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	authHandler := AuthHandler{}

	// Запрос без sessionID
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh-csrf", nil)
	w := httptest.NewRecorder()

	// Вызов метода
	authHandler.RefreshCsrfToken(w, req)

	// Проверка
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
}

func TestAuthHandler_RefreshCsrfToken_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	grpcErr := status.Error(codes.Internal, "gRPC error")

	// Настройка мока
	mockGrpcClient.On("RefreshCsrfToken", mock.Anything, mock.Anything, mock.Anything).Return(&gen.RefreshCsrfTokenResponse{}, grpcErr)

	authHandler := AuthHandler{
		client: mockGrpcClient,
	}

	// Запрос с валидным sessionID
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh-csrf", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test_session_id",
	})
	w := httptest.NewRecorder()

	// Вызов метода
	authHandler.RefreshCsrfToken(w, req)

	// Проверка
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
