package mocks

import (
	"2024_2_FIGHT-CLUB/domain"
	"context"
)

type MockCitiesRepository struct {
	MockGetCities func(ctx context.Context) ([]domain.City, error)
}

func (m *MockCitiesRepository) GetCities(ctx context.Context) ([]domain.City, error) {
	return m.MockGetCities(ctx)
}

type MockCitiesUseCase struct {
	MockGetCities func(ctx context.Context) ([]domain.City, error)
}

func (m *MockCitiesUseCase) GetCities(ctx context.Context) ([]domain.City, error) {
	return m.MockGetCities(ctx)
}
