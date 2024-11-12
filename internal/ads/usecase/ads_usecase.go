package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/images"
	ntype "2024_2_FIGHT-CLUB/internal/service/type"
	"log"

	"context"
	"errors"
	"mime/multipart"
)

type AdUseCase interface {
	GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error)
	GetOnePlace(ctx context.Context, adId string, isAuthorized bool) (domain.GetAllAdsResponse, error)
	CreatePlace(ctx context.Context, place *domain.Ad, fileHeader []*multipart.FileHeader, newPlace domain.CreateAdRequest) error
	UpdatePlace(ctx context.Context, place *domain.Ad, adId string, userId string, fileHeader []*multipart.FileHeader, updatedPlace domain.UpdateAdRequest) error
	DeletePlace(ctx context.Context, adId string, userId string) error
	GetPlacesPerCity(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error)
	GetUserPlaces(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error)
	DeleteAdImage(ctx context.Context, adId string, imageId int, userId string) error
}

type adUseCase struct {
	adRepository domain.AdRepository
	minioService images.MinioServiceInterface
}

func NewAdUseCase(adRepository domain.AdRepository, minioService images.MinioServiceInterface) AdUseCase {
	return &adUseCase{
		adRepository: adRepository,
		minioService: minioService,
	}
}

func (uc *adUseCase) GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error) {
	ads, err := uc.adRepository.GetAllPlaces(ctx, filter)
	if err != nil {
		return nil, err
	}
	return ads, nil
}

func (uc *adUseCase) GetOnePlace(ctx context.Context, adId string, isAuthorized bool) (domain.GetAllAdsResponse, error) {
	ad, err := uc.adRepository.GetPlaceById(ctx, adId)
	if err != nil {
		return ad, errors.New("ad not found")
	}

	if isAuthorized {
		ad, err = uc.adRepository.UpdateViewsCount(ctx, ad)
		if err != nil {
			return ad, err
		}
	}

	return ad, nil
}

func (uc *adUseCase) CreatePlace(ctx context.Context, place *domain.Ad, fileHeaders []*multipart.FileHeader, newPlace domain.CreateAdRequest) error {
	place.Description = newPlace.Description
	place.Address = newPlace.Address
	place.RoomsNumber = newPlace.RoomsNumber
	err := uc.adRepository.CreatePlace(ctx, place, newPlace)
	if err != nil {
		return err
	}
	var uploadedPaths ntype.StringArray

	for _, fileHeader := range fileHeaders {
		if fileHeader != nil {
			uploadedPath, err := uc.minioService.UploadFile(fileHeader, "ads/"+place.UUID)
			if err != nil {
				for _, path := range uploadedPaths {
					_ = uc.minioService.DeleteFile(path)
				}
				return err
			}
			uploadedPaths = append(uploadedPaths, "/images/"+uploadedPath)
		}
	}

	err = uc.adRepository.SaveImages(ctx, place.UUID, uploadedPaths)
	if err != nil {
		return err
	}

	return nil
}

func (uc *adUseCase) UpdatePlace(ctx context.Context, place *domain.Ad, adId string, userId string, fileHeaders []*multipart.FileHeader, updatedPlace domain.UpdateAdRequest) error {
	_, err := uc.adRepository.GetPlaceById(ctx, adId)
	if err != nil {
		return err
	}
	place.Description = updatedPlace.Description
	place.Address = updatedPlace.Address
	place.RoomsNumber = updatedPlace.RoomsNumber
	var newUploadedPaths ntype.StringArray

	for _, fileHeader := range fileHeaders {
		if fileHeader != nil {
			uploadedPath, err := uc.minioService.UploadFile(fileHeader, "ads/"+adId)
			if err != nil {
				for _, path := range newUploadedPaths {
					_ = uc.minioService.DeleteFile(path)
				}
				return err
			}
			newUploadedPaths = append(newUploadedPaths, "/images/"+uploadedPath)
		}
	}

	err = uc.adRepository.UpdatePlace(ctx, place, adId, userId, updatedPlace)
	if err != nil {
		return err
	}

	err = uc.adRepository.SaveImages(ctx, adId, newUploadedPaths)
	if err != nil {
		return err
	}
	return nil
}

func (uc *adUseCase) DeletePlace(ctx context.Context, adId string, userId string) error {
	_, err := uc.adRepository.GetPlaceById(ctx, adId)
	if err != nil {
		return err
	}
	imagesPath, err := uc.adRepository.GetAdImages(ctx, adId)
	if err != nil {
		return err
	}
	for _, imagePath := range imagesPath {
		_ = uc.minioService.DeleteFile(imagePath)
	}

	err = uc.adRepository.DeletePlace(ctx, adId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (uc *adUseCase) GetPlacesPerCity(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
	places, err := uc.adRepository.GetPlacesPerCity(ctx, city)
	if err != nil {
		return nil, err
	}
	return places, nil
}

func (uc *adUseCase) GetUserPlaces(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
	places, err := uc.adRepository.GetUserPlaces(ctx, userId)
	if err != nil {
		return nil, err
	}
	return places, nil
}

func (uc *adUseCase) DeleteAdImage(ctx context.Context, adId string, imageId int, userId string) error {
	imageURL, err := uc.adRepository.DeleteAdImage(ctx, adId, imageId, userId)
	if err != nil {
		return err
	}

	if err := uc.minioService.DeleteFile(imageURL); err != nil {
		log.Printf("Warning: failed to delete file from MinIO: %v", err)
	}

	return nil
}
