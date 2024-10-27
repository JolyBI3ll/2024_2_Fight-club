package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/ads/mocks"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdHandler_GetAllPlaces(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := &mocks.MockAdUseCase{}
	mockSession := &mocks.MockServiceSession{}

	handler := NewAdHandler(mockUseCase, mockSession)

	filter := domain.AdFilter{
		Location:    "Test City",
		Rating:      "5",
		NewThisWeek: "true",
		HostGender:  "any",
		GuestCount:  "2",
	}

	expectedAds := []domain.Ad{
		{
			ID:           "1",
			LocationMain: "Test City",
			Position:     []float64{40.7128, -74.0060},
			Distance:     10.5,
		},
	}

	mockUseCase.MockGetAllPlaces = func(ctx context.Context, f domain.AdFilter) ([]domain.Ad, error) {
		assert.Equal(t, filter, f, "Filter should match")
		return expectedAds, nil
	}

	req, err := http.NewRequest("GET", "/api/ads/?location=Test+City&rating=5&new=true&gender=any&guests=2", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.GetAllPlaces(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string][]domain.Ad
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, expectedAds, response["places"], "Returned places should match expected")

	mockUseCase.MockGetAllPlaces = func(ctx context.Context, f domain.AdFilter) ([]domain.Ad, error) {
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := &mocks.MockAdUseCase{}
	mockSession := &mocks.MockServiceSession{}

	handler := NewAdHandler(mockUseCase, mockSession)

	adId := "1"
	expectedAd := domain.Ad{
		ID:           adId,
		LocationMain: "Test City",
		Position:     []float64{40.7128, -74.0060},
		Distance:     10.5,
	}

	mockUseCase.MockGetOnePlace = func(ctx context.Context, id string) (domain.Ad, error) {
		assert.Equal(t, adId, id, "Ad ID should match")
		return expectedAd, nil
	}

	req, err := http.NewRequest("GET", "/api/ads/"+adId, nil)
	assert.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{"adId": adId})

	rr := httptest.NewRecorder()

	handler.GetOnePlace(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string]domain.Ad
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, expectedAd, response["place"], "Returned place should match expected")

	mockUseCase.MockGetOnePlace = func(ctx context.Context, id string) (domain.Ad, error) {
		return domain.Ad{}, fmt.Errorf("place not found")
	}
	rr = httptest.NewRecorder()
	handler.GetOnePlace(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status 500")
}

func TestAdHandler_CreatePlace(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := &mocks.MockAdUseCase{}
	mockSession := &mocks.MockServiceSession{}

	handler := NewAdHandler(mockUseCase, mockSession)

	newAd := domain.Ad{
		ID:              "1",
		LocationMain:    "Test City",
		LocationStreet:  "Test Street",
		Position:        []float64{40.7128, -74.0060},
		Images:          []string{"image1.jpg", "image2.jpg"},
		AuthorUUID:      "user123",
		PublicationDate: "2024-04-27",
		AvailableDates:  []string{"2024-05-01", "2024-05-02"},
		Distance:        10.5,
	}

	mockUseCase.MockCreatePlace = func(ctx context.Context, ad *domain.Ad, files []*multipart.FileHeader) error {
		assert.Equal(t, &newAd, ad, "Ad should match expected")
		return nil
	}

	mockSession.MockGetUserID = func(ctx context.Context, r *http.Request, w http.ResponseWriter) (string, error) {
		return "user123", nil
	}

	metaData, err := json.Marshal(newAd)
	assert.NoError(t, err)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("metadata", string(metaData))
	writer.Close()

	req, err := http.NewRequest("POST", "/api/createAd", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler.CreatePlace(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string]domain.Ad
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, newAd, response["place"], "Returned place should match expected")

	mockUseCase.MockCreatePlace = func(ctx context.Context, ad *domain.Ad, files []*multipart.FileHeader) error {
		return fmt.Errorf("could not create place")
	}
	rr = httptest.NewRecorder()
	handler.CreatePlace(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status 500")
}

func TestAdHandler_UpdatePlace(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := &mocks.MockAdUseCase{}
	mockSession := &mocks.MockServiceSession{}

	handler := NewAdHandler(mockUseCase, mockSession)

	adId := "1"
	updatedAd := domain.Ad{
		ID:              adId,
		LocationMain:    "Updated City",
		LocationStreet:  "Updated Street",
		Position:        []float64{41.7128, -73.0060},
		Images:          []string{"updated_image1.jpg"},
		AuthorUUID:      "user123",
		PublicationDate: "2024-04-28",
		AvailableDates:  []string{"2024-05-03", "2024-05-04"},
		Distance:        12.0,
	}

	mockUseCase.MockUpdatePlace = func(ctx context.Context, ad *domain.Ad, id string, userId string, files []*multipart.FileHeader) error {
		assert.Equal(t, adId, id, "Ad ID should match")
		assert.Equal(t, "user123", userId, "User ID should match")
		assert.Equal(t, &updatedAd, ad, "Ad should match expected")
		return nil
	}

	mockSession.MockGetUserID = func(ctx context.Context, r *http.Request, w http.ResponseWriter) (string, error) {
		return "user123", nil
	}

	metaData, err := json.Marshal(updatedAd)
	assert.NoError(t, err)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("metadata", string(metaData))
	writer.Close()

	req, err := http.NewRequest("PUT", "/api/ads/"+adId, body)
	assert.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"adId": adId})
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler.UpdatePlace(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Update successfully", response["response"], "Response message should match expected")

	mockUseCase.MockUpdatePlace = func(ctx context.Context, ad *domain.Ad, id string, userId string, files []*multipart.FileHeader) error {
		return fmt.Errorf("could not update place")
	}
	rr = httptest.NewRecorder()
	handler.UpdatePlace(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status 500")
}

func TestAdHandler_DeletePlace(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := &mocks.MockAdUseCase{}
	mockSession := &mocks.MockServiceSession{}

	handler := NewAdHandler(mockUseCase, mockSession)

	adId := "1"

	mockUseCase.MockDeletePlace = func(ctx context.Context, id string, userId string) error {
		assert.Equal(t, adId, id, "Ad ID should match")
		assert.Equal(t, "user123", userId, "User ID should match")
		return nil
	}

	mockSession.MockGetUserID = func(ctx context.Context, r *http.Request, w http.ResponseWriter) (string, error) {
		return "user123", nil
	}

	req, err := http.NewRequest("DELETE", "/api/ads/"+adId, nil)
	assert.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"adId": adId})

	rr := httptest.NewRecorder()

	handler.DeletePlace(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	mockUseCase.MockDeletePlace = func(ctx context.Context, id string, userId string) error {
		return fmt.Errorf("could not delete place")
	}
	rr = httptest.NewRecorder()
	handler.DeletePlace(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status 500")
}

func TestAdHandler_GetPlacesPerCity(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := &mocks.MockAdUseCase{}
	mockSession := &mocks.MockServiceSession{}

	handler := NewAdHandler(mockUseCase, mockSession)

	city := "Test City"
	expectedAds := []domain.Ad{
		{
			ID:           "1",
			LocationMain: city,
			Position:     []float64{40.7128, -74.0060},
			Distance:     10.5,
		},
		{
			ID:           "2",
			LocationMain: city,
			Position:     []float64{40.7138, -74.0070},
			Distance:     12.5,
		},
	}

	mockUseCase.MockGetPlacesPerCity = func(ctx context.Context, c string) ([]domain.Ad, error) {
		assert.Equal(t, city, c, "City should match")
		return expectedAds, nil
	}

	req, err := http.NewRequest("GET", "/api/ads/city/"+city, nil)
	assert.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{"city": city})

	rr := httptest.NewRecorder()

	handler.GetPlacesPerCity(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	var response map[string][]domain.Ad
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, expectedAds, response["places"], "Returned places should match expected")

	mockUseCase.MockGetPlacesPerCity = func(ctx context.Context, c string) ([]domain.Ad, error) {
		return []domain.Ad{}, fmt.Errorf("could not get places")
	}
	rr = httptest.NewRecorder()
	handler.GetPlacesPerCity(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status 500")
}

func TestAdHandler_GetAllPlaces_Error(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := &mocks.MockAdUseCase{}
	mockSession := &mocks.MockServiceSession{}

	handler := NewAdHandler(mockUseCase, mockSession)

	_ = domain.AdFilter{
		Location:    "Test City",
		Rating:      "5",
		NewThisWeek: "true",
		HostGender:  "any",
		GuestCount:  "2",
	}

	mockUseCase.MockGetAllPlaces = func(ctx context.Context, f domain.AdFilter) ([]domain.Ad, error) {
		return nil, errors.New("database error")
	}

	req, err := http.NewRequest("GET", "/api/ads/?location=Test+City&rating=5&new=true&gender=any&guests=2", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.GetAllPlaces(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error")

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "database error", response["error"], "Error message should match expected")
}
