package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"context"
)

type CityUseCase interface {
	GetCities(ctx context.Context) ([]domain.City, error)
}

type cityUseCase struct {
	cityRepository domain.CityRepository
}

func NewCityUSeCase(cityRepository domain.CityRepository) CityUseCase {
	return &cityUseCase{
		cityRepository: cityRepository,
	}
}

func (c *cityUseCase) GetCities(ctx context.Context) ([]domain.City, error) {
	cities, err := c.cityRepository.GetCities(ctx)
	if err != nil {
		return nil, err
	}
	return cities, nil
}
