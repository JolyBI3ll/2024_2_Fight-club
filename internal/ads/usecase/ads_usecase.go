package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/images"
	ntype "2024_2_FIGHT-CLUB/internal/service/type"

	"context"
	"errors"
	"fmt"
	"mime/multipart"
)

type AdUseCase interface {
	GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.Ad, error)
	GetOnePlace(ctx context.Context, adId string) (domain.Ad, error)
	CreatePlace(ctx context.Context, place *domain.Ad, fileHeader []*multipart.FileHeader) error
	UpdatePlace(ctx context.Context, place *domain.Ad, adId string, userId string, fileHeader []*multipart.FileHeader) error
	DeletePlace(ctx context.Context, adId string, userId string) error
	GetPlacesPerCity(ctx context.Context, city string) ([]domain.Ad, error)
}

type adUseCase struct {
	adRepository domain.AdRepository
	minioService *images.MinioService
}

func NewAdUseCase(adRepository domain.AdRepository, minioService *images.MinioService) AdUseCase {
	return &adUseCase{
		adRepository: adRepository,
		minioService: minioService,
	}
}

func (uc *adUseCase) GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.Ad, error) {
	ads, err := uc.adRepository.GetAllPlaces(ctx, filter)
	if err != nil {
		return nil, err
	}
	return ads, nil
}

func (uc *adUseCase) GetOnePlace(ctx context.Context, adId string) (domain.Ad, error) {
	ad, err := uc.adRepository.GetPlaceById(ctx, adId)
	if err != nil {
		return ad, errors.New("ad not found")
	}
	return ad, nil
}

func (uc *adUseCase) CreatePlace(ctx context.Context, place *domain.Ad, fileHeaders []*multipart.FileHeader) error {
	err := uc.adRepository.CreatePlace(ctx, place)
	if err != nil {
		return err
	}
	var uploadedPaths ntype.StringArray

	for _, fileHeader := range fileHeaders {
		if fileHeader != nil {
			filePath := fmt.Sprintf("ads/%s/%s", place.ID, fileHeader.Filename)

			uploadedPath, err := uc.minioService.UploadFile(fileHeader, filePath)
			if err != nil {
				for _, path := range uploadedPaths {
					_ = uc.minioService.DeleteFile(path)
				}
				return err
			}
			uploadedPaths = append(uploadedPaths, "http://localhost:9000/images/"+uploadedPath)
		}
	}

	place.Images = uploadedPaths

	err = uc.adRepository.SavePlace(ctx, place)
	if err != nil {
		return err
	}
	return nil
}

func (uc *adUseCase) UpdatePlace(ctx context.Context, place *domain.Ad, adId string, userId string, fileHeaders []*multipart.FileHeader) error {
	existingPlace, err := uc.adRepository.GetPlaceById(ctx, adId)
	if err != nil {
		return err
	}

	oldImages := existingPlace.Images

	var newUploadedPaths ntype.StringArray

	for _, fileHeader := range fileHeaders {
		if fileHeader != nil {
			filePath := fmt.Sprintf("ads/%s/%s", adId, fileHeader.Filename)

			uploadedPath, err := uc.minioService.UploadFile(fileHeader, filePath)
			if err != nil {
				for _, path := range newUploadedPaths {
					_ = uc.minioService.DeleteFile(path)
				}
				return err
			}

			newUploadedPaths = append(newUploadedPaths, "http://localhost:9000/images/"+uploadedPath)
		}
	}

	place.Images = append(oldImages, newUploadedPaths...)

	err = uc.adRepository.UpdatePlace(ctx, place, adId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (uc *adUseCase) DeletePlace(ctx context.Context, adId string, userId string) error {
	place, err := uc.adRepository.GetPlaceById(ctx, adId)
	if err != nil {
		return err
	}

	for _, imagePath := range place.Images {
		_ = uc.minioService.DeleteFile(imagePath)
	}

	err = uc.adRepository.DeletePlace(ctx, adId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (uc *adUseCase) GetPlacesPerCity(ctx context.Context, city string) ([]domain.Ad, error) {
	places, err := uc.adRepository.GetPlacesPerCity(ctx, city)
	if err != nil || len(places) == 0 {
		return nil, errors.New("ad not found")
	}
	return places, nil
}
