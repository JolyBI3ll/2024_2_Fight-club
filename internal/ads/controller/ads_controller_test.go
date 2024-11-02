package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/ads/mocks"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAdHandler_GetAllPlaces(t *testing.T) {
	// Инициализация логгера
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	// Создание мока use case и других зависимостей
	mockUseCase := &mocks.MockAdUseCase{}
	mockSession := &mocks.MockServiceSession{}
	mockJwtToken := &mocks.MockJwtTokenService{}
	handler := NewAdHandler(mockUseCase, mockSession, mockJwtToken)

	// Определение фильтра и ожидаемых результатов
	filter := domain.AdFilter{
		Location:    "Test City",
		Rating:      "5",
		NewThisWeek: "true",
		HostGender:  "any",
		GuestCount:  "2",
	}

	expectedAds := []domain.GetAllAdsResponse{
		{
			UUID:        "1",
			Cityname:    "Test City",
			Address:     "test address",
			RoomsNumber: 10,
		},
	}

	mockUseCase.MockGetAllPlaces = func(ctx context.Context, f domain.AdFilter) ([]domain.GetAllAdsResponse, error) {
		assert.Equal(t, filter, f, "Filter should match")
		return expectedAds, nil
	}

	req, err := http.NewRequest("GET", "/api/ads/?location=Test+City&rating=5&new=true&gender=any&guests=2", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.GetAllPlaces(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	places, ok := response["places"].([]interface{})
	assert.True(t, ok, "Expected 'places' field in response")

	assert.Equal(t, len(expectedAds), len(places), "Returned places count should match expected")

	mockUseCase.MockGetAllPlaces = func(ctx context.Context, f domain.AdFilter) ([]domain.GetAllAdsResponse, error) {
		return nil, fmt.Errorf("server error")
	}

	rr = httptest.NewRecorder()
	handler.GetAllPlaces(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status 500")
}

func TestAdHandler_GetOnePlace(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	mockUseCase := &mocks.MockAdUseCase{}
	mockSession := &mocks.MockServiceSession{}
	mockJwtToken := &mocks.MockJwtTokenService{}
	handler := NewAdHandler(mockUseCase, mockSession, mockJwtToken)

	adId := "1"
	expectedAd := domain.GetAllAdsResponse{
		UUID:        "1",
		Cityname:    "Test City",
		Address:     "test address",
		RoomsNumber: 10,
	}

	mockUseCase.MockGetOnePlace = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		assert.Equal(t, adId, id, "Ad ID should match")
		return expectedAd, nil
	}

	req, err := http.NewRequest("GET", "/api/ads/"+adId, nil)
	assert.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{"adId": adId})

	rr := httptest.NewRecorder()

	handler.GetOnePlace(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string]domain.GetAllAdsResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, expectedAd, response["place"], "Returned place should match expected")

	mockUseCase.MockGetOnePlace = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return domain.GetAllAdsResponse{}, fmt.Errorf("place not found")
	}

	rr = httptest.NewRecorder()
	handler.GetOnePlace(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status 500")
}

func TestAdHandler_CreatePlace(t *testing.T) {
	// Инициализация логгера
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	// Ручное создание мок объектов
	mockUseCase := &mocks.MockAdUseCase{}
	mockSession := &mocks.MockServiceSession{}
	mockJwtToken := &mocks.MockJwtTokenService{}
	handler := NewAdHandler(mockUseCase, mockSession, mockJwtToken)

	newAdRequest := domain.CreateAdRequest{
		CityName:    "Test City",
		Address:     "Test Street",
		Description: "Test Description",
		RoomsNumber: 3,
	}

	testTime := time.Now()

	expectedAd := domain.Ad{
		UUID:            "1",
		CityID:          123,
		AuthorUUID:      "user123",
		Address:         "Test Street",
		PublicationDate: testTime,
		Description:     "Test Description",
		RoomsNumber:     3,
	}

	// Мокаем метод CreatePlace
	mockUseCase.MockCreatePlace = func(ctx context.Context, ad *domain.Ad, files []*multipart.FileHeader, newPlace domain.CreateAdRequest) error {
		ad.UUID = expectedAd.UUID
		ad.CityID = expectedAd.CityID
		ad.AuthorUUID = expectedAd.AuthorUUID
		ad.Address = expectedAd.Address
		ad.PublicationDate = expectedAd.PublicationDate
		ad.Description = expectedAd.Description
		ad.RoomsNumber = expectedAd.RoomsNumber
		return nil
	}

	// Мокаем другие компоненты
	mockSession.MockGetUserID = func(ctx context.Context, r *http.Request) (string, error) {
		return "user123", nil
	}

	mockJwtToken.MockValidate = func(tokenString string) (*middleware.JwtCsrfClaims, error) {
		if tokenString == "valid_token" {
			return &middleware.JwtCsrfClaims{}, nil
		}
		return nil, fmt.Errorf("invalid token")
	}

	metaData, err := json.Marshal(newAdRequest)
	assert.NoError(t, err)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("metadata", string(metaData))
	writer.Close()

	req, err := http.NewRequest("POST", "/api/ads", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-CSRF-Token", "bearer valid_token")

	rr := httptest.NewRecorder()

	handler.CreatePlace(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string]domain.Ad
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Используем сравнение без учёта времени
	if responseAd, ok := response["place"]; ok {
		assert.Equal(t, expectedAd.UUID, responseAd.UUID)
		assert.Equal(t, expectedAd.CityID, responseAd.CityID)
		assert.Equal(t, expectedAd.AuthorUUID, responseAd.AuthorUUID)
		assert.Equal(t, expectedAd.Address, responseAd.Address)
		assert.Equal(t, expectedAd.Description, responseAd.Description)
		assert.Equal(t, expectedAd.RoomsNumber, responseAd.RoomsNumber)
	} else {
		t.FailNow()
	}

	// Тест на случай ошибки в создании
	mockUseCase.MockCreatePlace = func(ctx context.Context, ad *domain.Ad, files []*multipart.FileHeader, newPlace domain.CreateAdRequest) error {
		return fmt.Errorf("could not create place")
	}

	rr = httptest.NewRecorder()
	handler.CreatePlace(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status 500")
}

func TestAdHandler_UpdatePlace(t *testing.T) {
	// Инициализация логгеров
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	// Ручное создание мок-объектов для теста
	mockUseCase := &mocks.MockAdUseCase{}
	mockSession := &mocks.MockServiceSession{}
	mockJwtToken := &mocks.MockJwtTokenService{}

	handler := NewAdHandler(mockUseCase, mockSession, mockJwtToken)

	adId := "1"
	updatedAdRequest := domain.UpdateAdRequest{
		CityName:    "Updated City",
		Address:     "Updated Street",
		Description: "Updated Description",
		RoomsNumber: 3,
	}

	mockUseCase.MockUpdatePlace = func(ctx context.Context, ad *domain.Ad, id string, userID string, files []*multipart.FileHeader, updatedPlace domain.UpdateAdRequest) error {
		assert.Equal(t, adId, id, "Ad ID должен совпадать")
		assert.Equal(t, "user123", userID, "User ID должен совпадать")
		assert.Equal(t, updatedAdRequest, updatedPlace, "Обновленный ад должен совпадать с ожидаемым")
		return nil
	}

	mockSession.MockGetUserID = func(ctx context.Context, r *http.Request) (string, error) {
		return "user123", nil
	}

	mockJwtToken.MockValidate = func(tokenString string) (*middleware.JwtCsrfClaims, error) {
		if tokenString == "valid_token" {
			return &middleware.JwtCsrfClaims{}, nil
		}
		return nil, fmt.Errorf("invalid token")
	}

	metaData, err := json.Marshal(updatedAdRequest)
	assert.NoError(t, err)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("metadata", string(metaData))
	writer.Close()

	req, err := http.NewRequest("PUT", "/api/ads/"+adId, body)
	assert.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"adId": adId})
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-CSRF-Token", "Bearer valid_token")

	rr := httptest.NewRecorder()

	handler.UpdatePlace(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Update successfully", response["response"], "Response message должен совпадать с ожидаемым")

	// Проверка на случай ошибки в обновлении
	mockUseCase.MockUpdatePlace = func(ctx context.Context, ad *domain.Ad, id string, userID string, files []*multipart.FileHeader, updatedPlace domain.UpdateAdRequest) error {
		return fmt.Errorf("could not update place")
	}

	rr = httptest.NewRecorder()
	handler.UpdatePlace(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status 500")
}

func TestAdHandler_DeletePlace(t *testing.T) {
	// Инициализация логгеров
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	// Ручное создание мок-объектов
	mockUseCase := &mocks.MockAdUseCase{}
	mockSession := &mocks.MockServiceSession{}
	mockJwtToken := &mocks.MockJwtTokenService{}

	handler := NewAdHandler(mockUseCase, mockSession, mockJwtToken)

	adId := "1"

	mockUseCase.MockDeletePlace = func(ctx context.Context, id string, userID string) error {
		assert.Equal(t, adId, id, "Ad ID должен совпадать")
		assert.Equal(t, "user123", userID, "User ID должен совпадать")
		return nil
	}

	mockSession.MockGetUserID = func(ctx context.Context, r *http.Request) (string, error) {
		return "user123", nil
	}

	mockJwtToken.MockValidate = func(tokenString string) (*middleware.JwtCsrfClaims, error) {
		if tokenString == "valid_token" {
			return &middleware.JwtCsrfClaims{}, nil
		}
		return nil, fmt.Errorf("invalid token")
	}

	req, err := http.NewRequest("DELETE", "/api/ads/"+adId, nil)
	assert.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"adId": adId})
	req.Header.Set("X-CSRF-Token", "Bearer valid_token")

	rr := httptest.NewRecorder()

	handler.DeletePlace(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Delete successfully", response["response"], "Response message должен совпадать с ожидаемым")

	// Проверка на случай ошибки в удалении
	mockUseCase.MockDeletePlace = func(ctx context.Context, id string, userID string) error {
		return fmt.Errorf("could not delete place")
	}

	rr = httptest.NewRecorder()
	handler.DeletePlace(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status 500")
}

func TestAdHandler_GetPlacesPerCity(t *testing.T) {
	// Инициализация логгеров
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	// Ручное создание мок-объектов
	mockUseCase := &mocks.MockAdUseCase{}

	handler := NewAdHandler(mockUseCase, nil, nil)

	cityName := "Test City"
	expectedAds := []domain.GetAllAdsResponse{
		{UUID: "1", Address: "Test Street", RoomsNumber: 10},
		{UUID: "2", Address: "Test Street2", RoomsNumber: 12},
	}
	// Успешный сценарий
	mockUseCase.MockGetPlacesPerCity = func(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
		assert.Equal(t, cityName, city, "City name должен совпадать")
		return expectedAds, nil
	}

	req, err := http.NewRequest("GET", "/api/ads/city/"+cityName, nil)
	assert.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"city": cityName})

	rr := httptest.NewRecorder()
	handler.GetPlacesPerCity(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string][]domain.GetAllAdsResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedAds, response["places"], "Returned places should match expected")

	// Ошибочный сценарий
	mockUseCase.MockGetPlacesPerCity = func(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
		return nil, fmt.Errorf("could not get places")
	}

	rr = httptest.NewRecorder()
	handler.GetPlacesPerCity(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status 500")
}

func TestAdHandler_GetUserPlaces(t *testing.T) {
	// Инициализация логгеров
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	// Ручное создание мок-объектов
	mockUseCase := &mocks.MockAdUseCase{}

	handler := NewAdHandler(mockUseCase, nil, nil)

	expectedUserId := "user123"
	expectedAds := []domain.GetAllAdsResponse{
		{UUID: "1", Address: "Test Street", RoomsNumber: 10},
		{UUID: "2", Address: "Test Street2", RoomsNumber: 12},
	}

	// Успешный сценарий
	mockUseCase.MockGetUserPlaces = func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
		assert.Equal(t, expectedUserId, userId, "User ID должен совпадать")
		return expectedAds, nil
	}

	req, err := http.NewRequest("GET", "/api/users/"+expectedUserId+"/ads", nil)
	assert.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"userId": expectedUserId})

	rr := httptest.NewRecorder()
	handler.GetUserPlaces(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string][]domain.GetAllAdsResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedAds, response["places"], "Returned places should match expected")

	// Ошибочный сценарий
	mockUseCase.MockGetUserPlaces = func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
		return nil, fmt.Errorf("could not get places")
	}

	rr = httptest.NewRecorder()
	handler.GetUserPlaces(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status 500")
}
