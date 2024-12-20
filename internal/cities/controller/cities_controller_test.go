package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/utils"
	"2024_2_FIGHT-CLUB/microservices/city_service/controller/gen"
	"2024_2_FIGHT-CLUB/microservices/city_service/mocks"
	"errors"
	"github.com/gorilla/mux"
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

func TestCityHandler_GetCities_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	// Моки
	mockClient := new(mocks.MockGrpcClient)
	mockUtils := &utils.MockUtils{}

	mockResponse := &gen.GetCitiesResponse{
		Cities: []*gen.City{
			{
				Id:          1,
				Title:       "Москва",
				Entitle:     "Moscow",
				Description: "Some Desc",
				Image:       "test_image",
			},
			{
				Id:          2,
				Title:       "Волгоград",
				Entitle:     "Volgograd",
				Description: "Some Desc",
				Image:       "test_image",
			},
		},
	}

	mockPayload := domain.AllCitiesResponse{
		Cities: []*domain.City{
			{
				ID:          1,
				Title:       "Москва",
				EnTitle:     "Moscow",
				Description: "Some Desc",
				Image:       "test_image",
			},
			{
				ID:          2,
				Title:       "Волгоград",
				EnTitle:     "Volgograd",
				Description: "Some Desc",
				Image:       "test_image",
			},
		},
	}

	// Настройка моков
	mockClient.On("GetCities", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)
	mockUtils.On("ConvertAllCitiesProtoToGo", mockResponse).Return(mockPayload, nil)

	cityHandler := CityHandler{
		client: mockClient,
		utils:  mockUtils,
	}

	req := httptest.NewRequest(http.MethodGet, "/cities", nil)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	w := httptest.NewRecorder()

	// Вызов метода
	cityHandler.GetCities(w, req)

	// Проверка
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	body, _ := io.ReadAll(result.Body)
	expectedResponse := `{"cities":[{"description":"Some Desc", "enTitle":"Moscow", "id":1, "image":"test_image", "title":"Москва"},{"description":"Some Desc", "enTitle":"Volgograd", "id":2, "image":"test_image", "title":"Волгоград"}]}`
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.JSONEq(t, expectedResponse, string(body))

	mockClient.AssertExpectations(t)
	mockUtils.AssertExpectations(t)
}

func TestCityHandler_GetCities_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	mockUtils := &utils.MockUtils{}

	grpcErr := status.Error(codes.Internal, "gRPC error")

	// Настройка мока
	mockClient.On("GetCities", mock.Anything, &gen.GetCitiesRequest{}, mock.Anything).Return(&gen.GetCitiesResponse{}, grpcErr)

	cityHandler := CityHandler{
		client: mockClient,
		utils:  mockUtils,
	}

	req := httptest.NewRequest(http.MethodGet, "/cities", nil)
	w := httptest.NewRecorder()

	// Вызов метода
	cityHandler.GetCities(w, req)

	// Проверка
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

func TestCityHandler_GetCities_ConversionError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockClient := new(mocks.MockGrpcClient)
	mockUtils := &utils.MockUtils{}

	mockResponse := &gen.GetCitiesResponse{
		Cities: []*gen.City{
			{
				Id:          1,
				Title:       "Москва",
				Entitle:     "Moscow",
				Description: "Some Desc",
				Image:       "test_image",
			},
			{
				Id:          2,
				Title:       "Волгоград",
				Entitle:     "Volgograd",
				Description: "Some Desc",
				Image:       "test_image",
			},
		},
	}

	conversionErr := errors.New("conversion error")

	// Настройка моков
	mockClient.On("GetCities", mock.Anything, &gen.GetCitiesRequest{}, mock.Anything).Return(mockResponse, nil)
	mockUtils.On("ConvertAllCitiesProtoToGo", mockResponse).Return(nil, conversionErr)

	cityHandler := CityHandler{
		client: mockClient,
		utils:  mockUtils,
	}

	req := httptest.NewRequest(http.MethodGet, "/cities", nil)
	w := httptest.NewRecorder()

	// Вызов метода
	cityHandler.GetCities(w, req)

	// Проверка
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

	mockClient.AssertExpectations(t)
	mockUtils.AssertExpectations(t)
}

func TestCityHandler_GetOneCity_Success(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	// Моки
	mockClient := new(mocks.MockGrpcClient)
	mockUtils := new(utils.MockUtils)

	mockCityResponse := &gen.GetOneCityResponse{
		City: &gen.City{
			Id:          1,
			Title:       "Москва",
			Entitle:     "Moscow",
			Description: "Some Desc",
			Image:       "test_image",
		},
	}

	mockClient.On("GetOneCity", mock.Anything, mock.Anything, mock.Anything).Return(mockCityResponse, nil)
	mockUtils.On("ConvertOneCityProtoToGo", mockCityResponse.City).Return(domain.City{
		ID:          1,
		Title:       "Москва",
		EnTitle:     "Moscow",
		Description: "Some Desc",
		Image:       "test_image",
	}, nil)

	// Инициализация обработчика
	cityHandler := CityHandler{
		client: mockClient,
		utils:  mockUtils,
	}

	// HTTP-запрос
	req := httptest.NewRequest(http.MethodGet, "/city/test_city", nil)
	req = mux.SetURLVars(req, map[string]string{"city": "test_city"})
	w := httptest.NewRecorder()

	// Вызов метода
	cityHandler.GetOneCity(w, req)

	// Проверка
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	mockClient.AssertExpectations(t)
	mockUtils.AssertExpectations(t)
}

func TestCityHandler_GetOneCity_GrpcError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	// Моки
	mockClient := new(mocks.MockGrpcClient)
	grpcErr := status.Error(codes.Internal, "gRPC error")

	mockClient.On("GetOneCity", mock.Anything, mock.Anything, mock.Anything).Return(&gen.GetOneCityResponse{}, grpcErr)

	// Инициализация обработчика
	cityHandler := CityHandler{
		client: mockClient,
	}

	req := httptest.NewRequest(http.MethodGet, "/city/test_city", nil)
	req = mux.SetURLVars(req, map[string]string{"city": "test_city"})
	w := httptest.NewRecorder()

	// Вызов метода
	cityHandler.GetOneCity(w, req)

	// Проверка
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

func TestCityHandler_GetOneCity_ConversionError(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	// Моки
	mockClient := new(mocks.MockGrpcClient)
	mockUtils := new(utils.MockUtils)

	mockCityResponse := &gen.GetOneCityResponse{
		City: &gen.City{
			Id:          1,
			Title:       "Москва",
			Entitle:     "Moscow",
			Description: "Some Desc",
			Image:       "test_image",
		},
	}

	mockClient.On("GetOneCity", mock.Anything, mock.Anything, mock.Anything).Return(mockCityResponse, nil)
	mockUtils.On("ConvertOneCityProtoToGo", mockCityResponse.City).Return(nil, errors.New("conversion error"))

	// Инициализация обработчика
	cityHandler := CityHandler{
		client: mockClient,
		utils:  mockUtils,
	}

	req := httptest.NewRequest(http.MethodGet, "/city/test_city", nil)
	req = mux.SetURLVars(req, map[string]string{"city": "test_city"})
	w := httptest.NewRecorder()

	// Вызов метода
	cityHandler.GetOneCity(w, req)

	// Проверка
	result := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(result.Body)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	mockClient.AssertExpectations(t)
	mockUtils.AssertExpectations(t)
}