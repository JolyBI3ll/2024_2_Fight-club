package mocks

import (
	"2024_2_FIGHT-CLUB/domain"
	"context"
)

type MockCitiesRepository struct {
	MockGetCities       func(ctx context.Context) ([]domain.City, error)
	MockGetCityByEnName func(ctx context.Context, cityEnName string) (domain.City, error)
}

func (m *MockCitiesRepository) GetCities(ctx context.Context) ([]domain.City, error) {
	return m.MockGetCities(ctx)
}

func (m *MockCitiesRepository) GetCityByEnName(ctx context.Context, cityEnName string) (domain.City, error) {
	return m.MockGetCityByEnName(ctx, cityEnName)
}

type MockCitiesUseCase struct {
	MockGetCities  func(ctx context.Context) ([]domain.City, error)
	MockGetOneCity func(ctx context.Context, cityEnName string) (domain.City, error)
}

func (m *MockCitiesUseCase) GetCities(ctx context.Context) ([]domain.City, error) {
	return m.MockGetCities(ctx)
}

func (m *MockCitiesUseCase) GetOneCity(ctx context.Context, cityEnName string) (domain.City, error) {
	return m.MockGetOneCity(ctx, cityEnName)
}
