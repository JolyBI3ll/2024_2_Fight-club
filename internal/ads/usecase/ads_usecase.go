package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"errors"
)

type AdUseCase interface {
	GetAllPlaces() ([]domain.Ad, error)
	GetOnePlace(adId string) (domain.Ad, error)
	CreatePlace(place *domain.Ad) error
	UpdatePlace(place *domain.Ad, adId string, userId string) error
	DeletePlace(adId string, userId string) error
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

func (uc *adUseCase) GetOnePlace(adId string) (domain.Ad, error) {
	ad, err := uc.adRepository.GetPlaceById(adId)
	if err != nil {
		return ad, errors.New("ad not found")
	}
	return ad, nil
}

func (uc *adUseCase) CreatePlace(place *domain.Ad) error {
	err := uc.adRepository.CreatePlace(place)
	if err != nil {
		return err
	}
	return nil
}

func (uc *adUseCase) UpdatePlace(place *domain.Ad, adId string, userId string) error {
	err := uc.adRepository.UpdatePlace(place, adId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (uc *adUseCase) DeletePlace(adId string, userId string) error {
	err := uc.adRepository.DeletePlace(adId, userId)
	if err != nil {
		return err
	}
	return nil
}