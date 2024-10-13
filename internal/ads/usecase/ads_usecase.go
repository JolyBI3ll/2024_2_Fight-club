package usecase

import "2024_2_FIGHT-CLUB/domain"

type AdUseCase interface {
	GetAllPlaces() ([]domain.Ad, error)
}

type adUseCase struct {
	adRepository domain.AdRepository
}

func NewAdUseCase(adRepository domain.AdRepository) AdUseCase {
	return &adUseCase{
		adRepository: adRepository,
	}
}

func (uc *adUseCase) GetAllPlaces() ([]domain.Ad, error) {
	ads, err := uc.adRepository.GetAllPlaces()
	if err != nil {
		return nil, err
	}
	return ads, nil
}
