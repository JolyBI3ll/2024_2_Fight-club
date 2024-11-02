package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/cities/mocks"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"context"
	"errors"
	"log"
	"reflect"
	"testing"
)

func TestGetCitiesSuccess(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	mockRepo := &mocks.MockCitiesRepository{
		MockGetCities: func(ctx context.Context) ([]domain.City, error) {
			return []domain.City{
				{ID: 1, Title: "Moscow", EnTitle: "moscow", Description: "A large city in Russia."},
			}, nil
		},
	}

	cityUsecase := NewCityUSeCase(mockRepo)

	ctx := context.TODO()
	cities, err := cityUsecase.GetCities(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedCities := []domain.City{
		{ID: 1, Title: "Moscow", EnTitle: "moscow", Description: "A large city in Russia."},
	}
	if !reflect.DeepEqual(cities, expectedCities) {
		t.Errorf("expected %v, got %v", expectedCities, cities)
	}
}

func TestGetCitiesFailure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	mockRepo := &mocks.MockCitiesRepository{
		MockGetCities: func(ctx context.Context) ([]domain.City, error) {
			return nil, errors.New("failed to retrieve cities")
		},
	}

	cityUsecase := NewCityUSeCase(mockRepo)

	ctx := context.TODO()
	_, err := cityUsecase.GetCities(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}

	expectedErrorMsg := "failed to retrieve cities"
	if err.Error() != expectedErrorMsg {
		t.Errorf("expected error message %v, got %v", expectedErrorMsg, err.Error())
	}
}
