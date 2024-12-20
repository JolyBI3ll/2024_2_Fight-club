package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/utils"
	"2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	"2024_2_FIGHT-CLUB/microservices/ads_service/mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdHandler_GetAllPlaces_Success(t *testing.T) {
	// Инициализация логгера
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	testResponse := &gen.GetAllAdsResponseList{}
	mockGrpcClient.On("GetAllPlaces", mock.Anything, mock.Anything, mock.Anything).Return(testResponse, nil)
	response := &gen.GetAllAdsResponseList{}
	utilsMock := &utils.MockUtils{}
	utilsMock.On("ConvertGetAllAdsResponseProtoToGo", response).
		Return(nil, nil)

	adHandler := &AdHandler{
		client: mockGrpcClient,
		utils:  utilsMock,
	}

	req := httptest.NewRequest(http.MethodGet, "/housing?location=test", nil)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	w := httptest.NewRecorder()

	adHandler.GetAllPlaces(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
}

func TestAdHandler_GetAllPlaces_Error(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	grpcErr := status.Error(codes.Internal, "Simulated gRPC Error")

	mockGrpcClient := new(mocks.MockGrpcClient)
	mockGrpcClient.On("GetAllPlaces", mock.Anything, mock.Anything, mock.Anything).
		Return((*gen.GetAllAdsResponseList)(nil), grpcErr)

	adHandler := &AdHandler{
		client: mockGrpcClient,
	}

	req := httptest.NewRequest(http.MethodGet, "/housing", nil)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	w := httptest.NewRecorder()

	adHandler.GetAllPlaces(w, req)

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

func TestAdHandler_GetAllPlaces_ConvertError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	response := &gen.GetAllAdsResponseList{}
	mockGrpcClient := new(mocks.MockGrpcClient)
	mockGrpcClient.On("GetAllPlaces", mock.Anything, mock.Anything, mock.Anything).
		Return(response, nil)

	utilsMock := &utils.MockUtils{}
	utilsMock.On("ConvertGetAllAdsResponseProtoToGo", response).
		Return(nil, errors.New("conversion error"))

	adHandler := &AdHandler{
		client: mockGrpcClient,
		utils:  utilsMock,
	}

	req := httptest.NewRequest(http.MethodGet, "/housing", nil)
	w := httptest.NewRecorder()

	// Выполнение метода
	adHandler.GetAllPlaces(w, req)

	// Проверка результата
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

func TestAdHandler_GetOnePlace_Success(t *testing.T) {
	// Инициализация логгера
	require.NoError(t, logger.InitLoggers())
	defer func() {
		require.NoError(t, logger.SyncLoggers())
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	mockSessionService := &mocks.MockServiceSession{}
	mockUtils := new(utils.MockUtils)

	testResponse := &gen.GetAllAdsResponse{}
	mockGrpcClient.On("GetOnePlace", mock.Anything, mock.Anything, mock.Anything).Return(testResponse, nil)

	convertedResponse := &domain.GetAllAdsResponse{}
	mockUtils.On("ConvertAdProtoToGo", testResponse).Return(convertedResponse, nil)

	mockSessionService.MockGetUserID = func(ctx context.Context, sessionID string) (string, error) {
		return "123", nil
	}

	adHandler := &AdHandler{
		client:         mockGrpcClient,
		sessionService: mockSessionService,
		utils:          mockUtils,
	}

	req := httptest.NewRequest(http.MethodGet, "/housing/123", nil)
	req = mux.SetURLVars(req, map[string]string{"adId": "123"})
	req.Header.Set("X-Real-IP", "127.0.0.1")
	w := httptest.NewRecorder()

	adHandler.GetOnePlace(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
	mockUtils.AssertExpectations(t)
}

func TestAdHandler_GetOnePlace_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	mockSessionService := &mocks.MockServiceSession{}

	grpcErr := status.Error(codes.Internal, "Simulated gRPC Error")
	mockGrpcClient.On("GetOnePlace", mock.Anything, mock.Anything, mock.Anything).Return((*gen.GetAllAdsResponse)(nil), grpcErr)

	mockSessionService.MockGetUserID = func(ctx context.Context, sessionID string) (string, error) {
		return "123", nil
	}

	adHandler := &AdHandler{
		client:         mockGrpcClient,
		sessionService: mockSessionService,
	}

	req := httptest.NewRequest(http.MethodGet, "/housing/123", nil)
	req = mux.SetURLVars(req, map[string]string{"adId": "123"})
	req.Header.Set("X-Real-IP", "127.0.0.1")
	w := httptest.NewRecorder()

	adHandler.GetOnePlace(w, req)

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

func TestAdHandler_GetOnePlace_ConvertError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockGrpcClient := new(mocks.MockGrpcClient)
	mockSessionService := &mocks.MockServiceSession{}
	mockUtils := new(utils.MockUtils)

	testResponse := &gen.GetAllAdsResponse{}
	mockGrpcClient.On("GetOnePlace", mock.Anything, mock.Anything, mock.Anything).Return(testResponse, nil)

	mockUtils.On("ConvertAdProtoToGo", testResponse).Return(nil, errors.New("conversion error"))

	mockSessionService.MockGetUserID = func(ctx context.Context, sessionID string) (string, error) {
		return "123", nil
	}

	adHandler := &AdHandler{
		client:         mockGrpcClient,
		sessionService: mockSessionService,
		utils:          mockUtils,
	}

	req := httptest.NewRequest(http.MethodGet, "/housing/123", nil)
	req = mux.SetURLVars(req, map[string]string{"adId": "123"})
	req.Header.Set("X-Real-IP", "127.0.0.1")
	w := httptest.NewRecorder()

	adHandler.GetOnePlace(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
	mockUtils.AssertExpectations(t)
}

func TestAdHandler_CreatePlace_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{
		client: mockClient,
	}

	// Создаем multipart запрос с метаданными и файлами
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Добавляем JSON metadata
	metadata := `{
		"CityName": "TestCity",
		"Description": "Test Description",
		"Address": "Test Address",
		"RoomsNumber": 2,
		"SquareMeters": 60,
		"Floor": 3,
		"BuildingType": "Apartment",
		"HasBalcony": true
	}`
	_ = writer.WriteField("metadata", metadata)

	part, _ := writer.CreateFormFile("images", "image1.jpg")
	_, err := part.Write([]byte("mock image data"))
	if err != nil {
		return
	}

	err = writer.Close()
	if err != nil {
		return
	}

	req := httptest.NewRequest("POST", "/housing", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-CSRF-Token", "test-token")
	req.Header.Set("X-Real-IP", "127.0.0.1")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})

	w := httptest.NewRecorder()

	mockClient.On("CreatePlace", mock.Anything, mock.Anything, mock.Anything).Return(&gen.Ad{}, nil)

	handler.CreatePlace(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	bodyResponse, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(bodyResponse), "Successfully created ad")

	mockClient.AssertExpectations(t)
}

func TestAdHandler_CreatePlace_FailedToReadFile(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{
		client: mockClient,
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_, err := writer.CreateFormFile("images", "image1.jpg")
	if err != nil {
		return
	}

	err = writer.Close()
	if err != nil {
		return
	}

	req := httptest.NewRequest("POST", "/ads/create", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})

	w := httptest.NewRecorder()

	// Выполнение метода
	handler.CreatePlace(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestAdHandler_CreatePlace_CreatePlaceFailure(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{
		client: mockClient,
	}

	// Создаем multipart запрос с метаданными и файлами
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("metadata", `{"CityName":"TestCity"}`)
	part, _ := writer.CreateFormFile("images", "image1.jpg")
	_, err := part.Write([]byte("mock image data"))
	if err != nil {
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}

	req := httptest.NewRequest("POST", "/ads/create", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})

	w := httptest.NewRecorder()
	grpcErr := status.Error(codes.Internal, "Simulated gRPC Error")
	mockClient.On("CreatePlace", mock.Anything, mock.Anything, mock.Anything).Return(&gen.Ad{}, grpcErr)

	handler.CreatePlace(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockClient.AssertExpectations(t)
}

func TestAdHandler_CreatePlace_NoCookie(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{
		client: mockClient,
	}

	// Создаем multipart запрос с метаданными и файлами
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("metadata", `{"CityName":"TestCity"}`)
	part, _ := writer.CreateFormFile("images", "image1.jpg")
	_, err := part.Write([]byte("mock image data"))
	if err != nil {
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}

	req := httptest.NewRequest("POST", "/ads/create", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	handler.CreatePlace(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockClient.AssertExpectations(t)
}

func TestAdHandler_UpdatePlace_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	metadata := `{
		"CityName": "TestCity",
		"Description": "Updated Description",
		"Address": "Test Address",
		"RoomsNumber": 3,
		"SquareMeters": 70,
		"Floor": 4,
		"BuildingType": "Apartment",
		"HasBalcony": false
	}`

	_ = writer.WriteField("metadata", metadata)

	part, _ := writer.CreateFormFile("images", "image1.jpg")
	_, err := part.Write([]byte("mock image data"))
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "/housing/123", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-CSRF-Token", "test-token")
	req.Header.Set("X-Real-IP", "127.0.0.1")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})

	w := httptest.NewRecorder()

	mockClient.On("UpdatePlace", mock.Anything, mock.Anything, mock.Anything).Return(&gen.AdResponse{}, nil)

	handler.UpdatePlace(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	bodyResponse, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(bodyResponse), "Successfully updated ad")

	mockClient.AssertExpectations(t)
}

func TestAdHandler_UpdatePlace_ParseMultipartError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	handler := &AdHandler{client: new(mocks.MockGrpcClient)}

	req := httptest.NewRequest("PUT", "/housing/123", bytes.NewBufferString("invalid body"))
	req.Header.Set("Content-Type", "multipart/form-data")

	w := httptest.NewRecorder()

	handler.UpdatePlace(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAdHandler_UpdatePlace_MetadataDecodeError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	handler := &AdHandler{client: new(mocks.MockGrpcClient)}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("metadata", "invalid json")

	_ = writer.Close()

	req := httptest.NewRequest("PUT", "/housing/123", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})

	w := httptest.NewRecorder()

	handler.UpdatePlace(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAdHandler_UpdatePlace_SessionIDError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	handler := &AdHandler{client: new(mocks.MockGrpcClient)}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("metadata", "{}")
	_ = writer.Close()

	req := httptest.NewRequest("PUT", "/housing/123", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	handler.UpdatePlace(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestAdHandler_UpdatePlace_GRPCError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	metadata := `{"CityName": "TestCity"}`
	_ = writer.WriteField("metadata", metadata)
	_ = writer.Close()

	req := httptest.NewRequest("PUT", "/housing/123", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})

	w := httptest.NewRecorder()
	grpcErr := status.Error(codes.Internal, "Simulated gRPC Error")
	mockClient.On("UpdatePlace", mock.Anything, mock.Anything, mock.Anything).Return(&gen.AdResponse{}, grpcErr)

	handler.UpdatePlace(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestAdHandler_DeletePlace_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("DELETE", "/housing/123", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})

	w := httptest.NewRecorder()

	mockClient.On("DeletePlace", mock.Anything, mock.Anything, mock.Anything).Return(&gen.DeleteResponse{}, nil)

	handler.DeletePlace(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, w.Body.String(), "Successfully deleted place")

	mockClient.AssertExpectations(t)
}

func TestAdHandler_DeletePlace_FailedToGetSessionID(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("DELETE", "/housing/123", nil)
	req.Header.Set("X-CSRF-Token", "test-token")

	w := httptest.NewRecorder()

	handler.DeletePlace(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.Contains(t, w.Body.String(), "failed to get session id from request cookie")

	mockClient.AssertExpectations(t)
}

func TestAdHandler_DeletePlace_ClientError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("DELETE", "/housing/123", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})

	w := httptest.NewRecorder()
	grpcErr := status.Error(codes.Internal, "failed to delete place")
	mockClient.On("DeletePlace", mock.Anything, mock.Anything, mock.Anything).Return(&gen.DeleteResponse{}, grpcErr)

	handler.DeletePlace(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.Contains(t, w.Body.String(), "failed to delete place")

	mockClient.AssertExpectations(t)
}

func TestAdHandler_GetPlacesPerCity_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	utilsMock := &utils.MockUtils{}
	handler := &AdHandler{client: mockClient, utils: utilsMock}

	req := httptest.NewRequest("GET", "/housing/city/TestCity", nil)
	req.Header.Set("X-CSRF-Token", "test-token")

	w := httptest.NewRecorder()

	mockClient.On("GetPlacesPerCity", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.GetAllAdsResponseList{}, nil)

	utilsMock.On("ConvertGetAllAdsResponseProtoToGo", mock.Anything).
		Return([]domain.GetAllAdsListResponse{}, nil)

	handler.GetPlacesPerCity(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)

	mockClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)
}

func TestAdHandler_GetPlacesPerCity_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("GET", "/housing/city/TestCity", nil)
	req.Header.Set("X-CSRF-Token", "test-token")

	w := httptest.NewRecorder()
	grpcError := status.Error(codes.Internal, "failed to get place per city")
	mockClient.On("GetPlacesPerCity", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.GetAllAdsResponseList{}, grpcError)

	handler.GetPlacesPerCity(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	mockClient.AssertExpectations(t)
}

func TestAdHandler_GetPlacesPerCity_ConvertError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	utilsMock := &utils.MockUtils{}
	handler := &AdHandler{client: mockClient, utils: utilsMock}

	req := httptest.NewRequest("GET", "/housing/city/TestCity", nil)
	req.Header.Set("X-CSRF-Token", "test-token")

	w := httptest.NewRecorder()

	mockClient.On("GetPlacesPerCity", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.GetAllAdsResponseList{}, nil)

	utilsMock.On("ConvertGetAllAdsResponseProtoToGo", mock.Anything).
		Return([]domain.GetAllAdsListResponse{}, errors.New("conversion error"))

	handler.GetPlacesPerCity(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

	mockClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)
}

func TestAdHandler_GetUserPlaces_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	mockUtils := &utils.MockUtils{}
	handler := &AdHandler{client: mockClient, utils: mockUtils}

	req := httptest.NewRequest("GET", "/housing/user/testUserId", nil)
	req.Header.Set("X-CSRF-Token", "test-token")

	w := httptest.NewRecorder()

	mockClient.On("GetUserPlaces", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.GetAllAdsResponseList{}, nil)

	mockUtils.On("ConvertGetAllAdsResponseProtoToGo", mock.Anything).
		Return([]domain.GetAllAdsListResponse{}, nil)

	handler.GetUserPlaces(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	mockClient.AssertExpectations(t)
	mockUtils.AssertExpectations(t)
}

func TestAdHandler_GetUserPlaces_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("GET", "/housing/user/testUserId", nil)
	req.Header.Set("X-CSRF-Token", "test-token")

	w := httptest.NewRecorder()
	grpcError := status.Error(codes.Internal, "failed to get place per city")
	mockClient.On("GetUserPlaces", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.GetAllAdsResponseList{}, grpcError)

	handler.GetUserPlaces(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

	mockClient.AssertExpectations(t)
}

func TestAdHandler_GetUserPlaces_ConvertError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	utilsMock := &utils.MockUtils{}
	handler := &AdHandler{client: mockClient, utils: utilsMock}

	req := httptest.NewRequest("GET", "/housing/user/testUserId", nil)
	req.Header.Set("X-CSRF-Token", "test-token")

	w := httptest.NewRecorder()

	mockClient.On("GetUserPlaces", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.GetAllAdsResponseList{}, nil)

	utilsMock.On("ConvertGetAllAdsResponseProtoToGo", mock.Anything).
		Return([]domain.GetAllAdsListResponse{}, errors.New("conversion error"))

	handler.GetUserPlaces(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

	mockClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)
}

func TestAdHandler_DeleteAdImage_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)

	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("DELETE", "/housing/{adId}/images/{imageId}", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})

	w := httptest.NewRecorder()

	mockClient.On("DeleteAdImage", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.DeleteResponse{}, nil)

	handler.DeleteAdImage(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, w.Body.String(), "\"message\":\"Successfully deleted ad image\"")

	mockClient.AssertExpectations(t)
}

func TestAdHandler_DeleteAdImage_FailedToGetSessionID(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	handler := &AdHandler{}

	req := httptest.NewRequest("DELETE", "/housing/{adId}/images/{imageId}", nil)
	req.Header.Set("X-CSRF-Token", "test-token")

	w := httptest.NewRecorder()

	handler.DeleteAdImage(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.Contains(t, w.Body.String(), "failed to get session id from request cookie")
}

func TestAdHandler_DeleteAdImage_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)

	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("DELETE", "/housing/{adId}/images/{imageId}", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})

	w := httptest.NewRecorder()
	grpcError := status.Error(codes.Internal, "failed to get place per city")
	mockClient.On("DeleteAdImage", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.DeleteResponse{}, grpcError)

	handler.DeleteAdImage(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	require.Contains(t, w.Body.String(), "\"error\":\"failed to get place per city\"")

	mockClient.AssertExpectations(t)
}

func TestAdHandler_AddToFavorites_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("POST", "/housing/{adId}/like", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})

	w := httptest.NewRecorder()

	mockClient.On("AddToFavorites", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.AdResponse{}, nil)

	handler.AddToFavorites(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, w.Body.String(), "\"message\":\"Successfully added to favorites\"")

	mockClient.AssertExpectations(t)
}

func TestAdHandler_AddToFavorites_FailedToGetSessionID(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	handler := &AdHandler{}

	req := httptest.NewRequest("POST", "/housing/{adId}/like", nil)
	req.Header.Set("X-CSRF-Token", "test-token")

	w := httptest.NewRecorder()

	handler.AddToFavorites(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.Contains(t, w.Body.String(), "failed to get session id from request cookie")

}

func TestAdHandler_AddToFavorites_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("POST", "/housing/{adId}/like", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})

	w := httptest.NewRecorder()
	grpcError := status.Error(codes.Internal, "failed to add to favorites")
	mockClient.On("AddToFavorites", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.AdResponse{}, grpcError)

	handler.AddToFavorites(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	require.Contains(t, w.Body.String(), "\"error\":\"failed to add to favorites\"")

	mockClient.AssertExpectations(t)
}

func TestAdHandler_DeleteFromFavorites_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("DELETE", "/housing/{adId}/dislike", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})

	w := httptest.NewRecorder()

	mockClient.On("DeleteFromFavorites", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.AdResponse{}, nil)

	handler.DeleteFromFavorites(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, w.Body.String(), "\"message\":\"Successfully deleted from favorites\"")

	mockClient.AssertExpectations(t)
}

func TestAdHandler_DeleteFromFavorites_FailedToGetSessionID(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	handler := &AdHandler{}

	req := httptest.NewRequest("POST", "/housing/{adId}/dislike", nil)
	req.Header.Set("X-CSRF-Token", "test-token")

	w := httptest.NewRecorder()

	handler.DeleteFromFavorites(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.Contains(t, w.Body.String(), "failed to get session id from request cookie")

}

func TestAdHandler_DeleteFromFavorites_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("POST", "/housing/{adId}/dislike", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})

	w := httptest.NewRecorder()
	grpcError := status.Error(codes.Internal, "failed to delete ad from favorites")
	mockClient.On("DeleteFromFavorites", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.AdResponse{}, grpcError)

	handler.DeleteFromFavorites(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	require.Contains(t, w.Body.String(), "\"error\":\"failed to delete ad from favorites\"")

	mockClient.AssertExpectations(t)
}

func TestAdHandler_GetUserFavorites_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	mockUtils := &utils.MockUtils{}
	handler := &AdHandler{client: mockClient, utils: mockUtils}

	req := httptest.NewRequest("GET", "/users/{userId}/favorites", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})

	w := httptest.NewRecorder()

	mockClient.On("GetUserFavorites", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.GetAllAdsResponseList{}, nil)

	mockUtils.On("ConvertGetAllAdsResponseProtoToGo", mock.Anything).
		Return([]domain.GetAllAdsListResponse{}, nil)

	handler.GetUserFavorites(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	mockClient.AssertExpectations(t)
	mockUtils.AssertExpectations(t)
}

func TestAdHandler_GetUserFavorites_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("GET", "/users/{userId}/favorites", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})

	w := httptest.NewRecorder()
	grpcError := status.Error(codes.Internal, "failed to user favorites")
	mockClient.On("GetUserFavorites", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.GetAllAdsResponseList{}, grpcError)

	handler.GetUserFavorites(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

	mockClient.AssertExpectations(t)
}

func TestAdHandler_GetUserFavorites_ConvertError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	utilsMock := &utils.MockUtils{}
	handler := &AdHandler{client: mockClient, utils: utilsMock}

	req := httptest.NewRequest("GET", "/users/{userId}/favorites", nil)
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})

	w := httptest.NewRecorder()

	mockClient.On("GetUserFavorites", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.GetAllAdsResponseList{}, nil)

	utilsMock.On("ConvertGetAllAdsResponseProtoToGo", mock.Anything).
		Return([]domain.GetAllAdsListResponse{}, errors.New("conversion error"))

	handler.GetUserFavorites(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

	mockClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)
}

func TestAdHandler_GetUserFavorites_FailedToGetSessionID(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	handler := &AdHandler{}

	req := httptest.NewRequest("GET", "/users/{userId}/favorites", nil)
	req.Header.Set("X-CSRF-Token", "test-token")

	w := httptest.NewRecorder()

	handler.GetUserFavorites(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.Contains(t, w.Body.String(), "failed to get session id from request cookie")
}

func TestAdHandler_UpdatePriorityWithPayment_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	card := domain.PaymentInfo{
		DonationAmount: "100",
	}
	cardData, _ := json.Marshal(card)

	req := httptest.NewRequest("POST", "/housing/{adId}/payment", bytes.NewReader(cardData))
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})
	w := httptest.NewRecorder()

	mockClient.On("UpdatePriority", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.AdResponse{}, nil)

	handler.UpdatePriorityWithPayment(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, w.Body.String(), "\"message\":\"Successfully update ad priority\"")

	mockClient.AssertExpectations(t)
}

func TestAdHandler_UpdatePriorityWithPayment_DecodeError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	req := httptest.NewRequest("POST", "/housing/{adId}/payment", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})
	w := httptest.NewRecorder()

	handler.UpdatePriorityWithPayment(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockClient.AssertExpectations(t)
}

func TestAdHandler_UpdatePriorityWithPayment_FailedToGetSessionID(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}
	card := domain.PaymentInfo{
		DonationAmount: "100",
	}
	cardData, _ := json.Marshal(card)

	req := httptest.NewRequest("POST", "/housing/{adId}/payment", bytes.NewReader(cardData))
	req.Header.Set("X-CSRF-Token", "test-token")
	w := httptest.NewRecorder()

	handler.UpdatePriorityWithPayment(w, req)

	resp := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.Contains(t, w.Body.String(), "failed to get session id from request cookie")
}

func TestAdHandler_UpdatePriorityWithPayment_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	handler := &AdHandler{client: mockClient}

	card := domain.PaymentInfo{
		DonationAmount: "100",
	}
	cardData, _ := json.Marshal(card)

	req := httptest.NewRequest("POST", "/housing/{adId}/payment", bytes.NewReader(cardData))
	req.Header.Set("X-CSRF-Token", "test-token")
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session-id",
	})
	w := httptest.NewRecorder()
	grpcError := status.Error(codes.Internal, "failed to update ad priority")
	mockClient.On("UpdatePriority", mock.Anything, mock.Anything, mock.Anything).
		Return(&gen.AdResponse{}, grpcError)

	handler.UpdatePriorityWithPayment(w, req)

	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

	mockClient.AssertExpectations(t)
}