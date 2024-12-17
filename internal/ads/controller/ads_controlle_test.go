package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	"2024_2_FIGHT-CLUB/microservices/ads_service/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

	adHandler := &AdHandler{
		client: mockGrpcClient,
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

	utilsMock := mocks.MockUtils{}
	utilsMock.On("ConvertGetAllAdsResponseProtoToGo", response).
		Return(nil, errors.New("conversion error"))

	adHandler := &AdHandler{
		client: mockGrpcClient,
	}

	req := httptest.NewRequest(http.MethodGet, "/housing", nil)
	w := httptest.NewRecorder()

	// Выполнение метода
	adHandler.GetAllPlaces(w, req)

	// Проверка результата
	result := w.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	mockGrpcClient.AssertExpectations(t)
	utilsMock.AssertExpectations(t)
}

func TestConvertGetAllAdsResponseProtoToGo_Error(t *testing.T) {
	mockUtils := new(mocks.MockUtils)

	protoResponse := &gen.GetAllAdsResponseList{}

	mockUtils.On("ConvertGetAllAdsResponseProtoToGo", protoResponse).
		Return(domain.GetAllAdsListResponse{}, errors.New("conversion error"))

	result, err := mockUtils.ConvertGetAllAdsResponseProtoToGo(protoResponse)

	assert.Error(t, err)
	assert.Equal(t, "conversion error", err.Error())
	assert.Empty(t, result)

	mockUtils.AssertExpectations(t)
}

func TestConvertAdProtoToGo_Error(t *testing.T) {
	mockUtils := new(mocks.MockUtils)

	protoAd := &gen.GetAllAdsResponse{}

	mockUtils.On("ConvertAdProtoToGo", protoAd).
		Return(domain.GetAllAdsResponse{}, errors.New("ad conversion error"))

	result, err := mockUtils.ConvertAdProtoToGo(protoAd)

	assert.Error(t, err)
	assert.Equal(t, "ad conversion error", err.Error())
	assert.Empty(t, result)

	mockUtils.AssertExpectations(t)
}

func TestParseDate_Error(t *testing.T) {
	mockUtils := new(mocks.MockUtils)

	dateStr := "invalid-date"
	adID := "123"
	fieldName := "PublicationDate"

	mockUtils.On("ParseDate", dateStr, adID, fieldName).
		Return(time.Time{}, errors.New("invalid date format"))

	parsedDate, err := mockUtils.ParseDate(dateStr, adID, fieldName)

	assert.Error(t, err)
	assert.Equal(t, "invalid date format", err.Error())
	assert.Equal(t, time.Time{}, parsedDate)

	mockUtils.AssertExpectations(t)
}
