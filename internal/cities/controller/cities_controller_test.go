package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/cities/mocks"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetCitiesSuccess(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	mockUseCase := &mocks.MockCitiesUseCase{
		MockGetCities: func(ctx context.Context) ([]domain.City, error) {
			return []domain.City{
				{ID: 1, Title: "Moscow", EnTitle: "moscow", Description: "A large city in Russia."},
			}, nil
		},
	}

	handler := NewCityHandler(mockUseCase)

	req, err := http.NewRequest("GET", "/api/cities", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetCities(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %v, got %v", http.StatusOK, status)
	}

	expectedCities := []domain.City{
		{ID: 1, Title: "Moscow", EnTitle: "moscow", Description: "A large city in Russia."},
	}
	var responseData map[string][]domain.City
	if err := json.Unmarshal(rr.Body.Bytes(), &responseData); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !reflect.DeepEqual(responseData["cities"], expectedCities) {
		t.Errorf("expected body %v, got %v", expectedCities, responseData["cities"])
	}
}

func TestGetCitiesFailure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	mockUseCase := &mocks.MockCitiesUseCase{
		MockGetCities: func(ctx context.Context) ([]domain.City, error) {
			return nil, errors.New("failed to retrieve cities")
		},
	}

	handler := NewCityHandler(mockUseCase)

	req, err := http.NewRequest("GET", "/api/cities", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetCities(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected status code %v, got %v", http.StatusInternalServerError, status)
	}

	expectedError := map[string]string{"error": "failed to retrieve cities"}
	var actualError map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &actualError); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if !reflect.DeepEqual(actualError, expectedError) {
		t.Errorf("expected body %v, got %v", expectedError, actualError)
	}
}
