package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/ads/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
)

func TestAdUseCase_GetAllPlaces(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	expectedAds := []domain.GetAllAdsResponse{
		{UUID: "1234", CityID: 1, AuthorUUID: "user123"},
	}
	mockRepo.MockGetAllPlaces = func(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error) {
		return expectedAds, nil
	}

	ctx := context.Background()
	filter := domain.AdFilter{Location: "New York"}
	ads, err := useCase.GetAllPlaces(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, expectedAds, ads)
}

func TestAdUseCase_GetOnePlace(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	adID := "ad123"
	expectedAd := domain.GetAllAdsResponse{UUID: adID, CityID: 2, AuthorUUID: "user567"}
	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return expectedAd, nil
	}

	ctx := context.Background()
	ad, err := useCase.GetOnePlace(ctx, adID)

	assert.NoError(t, err)
	assert.Equal(t, expectedAd, ad)
}

func TestAdUseCase_CreatePlace(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	newAd := domain.Ad{}
	fileHeaders := []*multipart.FileHeader{nil, nil}
	createRequest := domain.CreateAdRequest{
		CityName: "Los Angeles", Address: "123 Main St", Description: "Nice place", RoomsNumber: 2,
	}

	mockRepo.MockCreatePlace = func(ctx context.Context, ad *domain.Ad, req domain.CreateAdRequest) error {
		return nil
	}

	mockMinioService.UploadFileFunc = func(fileHeader *multipart.FileHeader, destinationPath string) (string, error) {
		return "uploadedPath", nil
	}

	mockRepo.MockSaveImages = func(ctx context.Context, adUUID string, imagePaths []string) error {
		return nil
	}

	ctx := context.Background()
	err := useCase.CreatePlace(ctx, &newAd, fileHeaders, createRequest)

	assert.NoError(t, err)
}

func TestAdUseCase_UpdatePlace(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	adID := "ad123"
	userID := "user456"
	existingAd := domain.Ad{}
	updateRequest := domain.UpdateAdRequest{
		CityName: "New City", Address: "456 New St", Description: "Updated description", RoomsNumber: 3,
	}
	fileHeaders := []*multipart.FileHeader{nil, nil}

	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return domain.GetAllAdsResponse{UUID: adID}, nil
	}

	mockRepo.MockUpdatePlace = func(ctx context.Context, ad *domain.Ad, aID, uID string, req domain.UpdateAdRequest) error {
		return nil
	}

	mockMinioService.UploadFileFunc = func(fileHeader *multipart.FileHeader, destinationPath string) (string, error) {
		return "uploadedPath", nil
	}

	mockRepo.MockSaveImages = func(ctx context.Context, adUUID string, imagePaths []string) error {
		return nil
	}

	ctx := context.Background()
	err := useCase.UpdatePlace(ctx, &existingAd, adID, userID, fileHeaders, updateRequest)

	assert.NoError(t, err)
}

func TestAdUseCase_DeletePlace(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	adID := "ad123"
	userID := "user456"
	imagePaths := []string{"images/path1", "images/path2"}

	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return domain.GetAllAdsResponse{UUID: adID}, nil
	}

	mockRepo.MockGetAdImages = func(ctx context.Context, aID string) ([]string, error) {
		return imagePaths, nil
	}

	mockMinioService.DeleteFileFunc = func(filePath string) error {
		return nil
	}

	mockRepo.MockDeletePlace = func(ctx context.Context, aID, uID string) error {
		return nil
	}

	ctx := context.Background()
	err := useCase.DeletePlace(ctx, adID, userID)

	assert.NoError(t, err)
}

func TestAdUseCase_GetPlacesPerCity(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	city := "New York"
	expectedPlaces := []domain.GetAllAdsResponse{
		{UUID: "1234", CityID: 1, AuthorUUID: "user123"},
	}

	ctx := context.Background()

	mockRepo.MockGetPlacesPerCity = func(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
		return expectedPlaces, nil
	}

	places, err := useCase.GetPlacesPerCity(ctx, city)

	assert.NoError(t, err)
	assert.Equal(t, expectedPlaces, places)

	mockRepo.MockGetPlacesPerCity = func(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
		return []domain.GetAllAdsResponse{}, nil
	}

	places, err = useCase.GetPlacesPerCity(ctx, city)

	assert.NoError(t, err)
	assert.Equal(t, []domain.GetAllAdsResponse{}, places)

	mockRepo.MockGetPlacesPerCity = func(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
		return nil, errors.New("database error")
	}

	places, err = useCase.GetPlacesPerCity(ctx, city)

	assert.Error(t, err)
	assert.Nil(t, places)
	if err != nil {
		assert.Equal(t, "database error", err.Error())
	}
}

func TestAdUseCase_GetUserPlaces(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	userID := "user123"
	expectedPlaces := []domain.GetAllAdsResponse{
		{UUID: "ad123", CityID: 2, AuthorUUID: userID},
	}

	// Успешный случай - объявления пользователя найдены
	mockRepo.MockGetUserPlaces = func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
		return expectedPlaces, nil
	}

	ctx := context.Background()
	places, err := useCase.GetUserPlaces(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedPlaces, places)

	// Случай, когда у пользователя нет объявлений
	mockRepo.MockGetUserPlaces = func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
		return []domain.GetAllAdsResponse{}, nil
	}

	places, err = useCase.GetUserPlaces(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(places))

	// Случай, когда произошла ошибка при запросе
	mockRepo.MockGetUserPlaces = func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
		return nil, errors.New("database error")
	}

	places, err = useCase.GetUserPlaces(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, places)
	assert.Equal(t, "database error", err.Error())
}

func TestAdUseCase_GetAllPlaces_Error(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	mockRepo.MockGetAllPlaces = func(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error) {
		return nil, errors.New("database error")
	}

	ctx := context.Background()
	filter := domain.AdFilter{Location: "New York"}
	ads, err := useCase.GetAllPlaces(ctx, filter)

	assert.Error(t, err)
	assert.Nil(t, ads)
	assert.Equal(t, "database error", err.Error())
}

func TestAdUseCase_GetOnePlace_Error(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return domain.GetAllAdsResponse{}, errors.New("ad not found")
	}

	ctx := context.Background()
	adID := "invalid_ad_id"
	ad, err := useCase.GetOnePlace(ctx, adID)

	assert.Error(t, err)
	assert.Equal(t, "ad not found", err.Error())
	assert.Equal(t, domain.GetAllAdsResponse{}, ad)
}

func TestAdUseCase_CreatePlace_ErrorOnCreate(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	newAd := domain.Ad{}
	createRequest := domain.CreateAdRequest{
		CityName: "Los Angeles", Address: "123 Main St", Description: "Nice place", RoomsNumber: 2,
	}

	mockRepo.MockCreatePlace = func(ctx context.Context, ad *domain.Ad, req domain.CreateAdRequest) error {
		return errors.New("creation failed")
	}

	ctx := context.Background()
	err := useCase.CreatePlace(ctx, &newAd, nil, createRequest)

	assert.Error(t, err)
	assert.Equal(t, "creation failed", err.Error())
}

func TestAdUseCase_CreatePlace_ErrorOnUpload(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	newAd := domain.Ad{}
	fileHeaders := []*multipart.FileHeader{&multipart.FileHeader{}, &multipart.FileHeader{}}
	createRequest := domain.CreateAdRequest{
		CityName: "Los Angeles", Address: "123 Main St", Description: "Nice place", RoomsNumber: 2,
	}

	mockRepo.MockCreatePlace = func(ctx context.Context, ad *domain.Ad, req domain.CreateAdRequest) error {
		return nil
	}

	mockMinioService.UploadFileFunc = func(fileHeader *multipart.FileHeader, destinationPath string) (string, error) {
		return "", errors.New("upload failed")
	}

	ctx := context.Background()
	err := useCase.CreatePlace(ctx, &newAd, fileHeaders, createRequest)

	assert.Error(t, err)
	assert.Equal(t, "upload failed", err.Error())
}

func TestAdUseCase_UpdatePlace_ErrorOnGet(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	adID := "invalid_ad_id"
	userID := "user456"
	updateRequest := domain.UpdateAdRequest{
		CityName: "New City", Address: "456 New St", Description: "Updated description", RoomsNumber: 3,
	}

	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return domain.GetAllAdsResponse{}, errors.New("ad not found")
	}

	ctx := context.Background()
	err := useCase.UpdatePlace(ctx, nil, adID, userID, nil, updateRequest)

	assert.Error(t, err)
	assert.Equal(t, "ad not found", err.Error())
}

func TestAdUseCase_DeletePlace_ErrorOnGet(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	adID := "invalid_ad_id"
	userID := "user456"

	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return domain.GetAllAdsResponse{}, errors.New("ad not found")
	}

	ctx := context.Background()
	err := useCase.DeletePlace(ctx, adID, userID)

	assert.Error(t, err)
	assert.Equal(t, "ad not found", err.Error())
}

func TestDeleteAdImage(t *testing.T) {
	ctx := context.Background()
	adId := "ad-uuid"
	imageId := 123
	userId := "user-uuid"
	imageURL := "/images/image.jpg"

	// Тест успешного удаления изображения
	t.Run("success", func(t *testing.T) {
		adRepoMock := &mocks.MockAdRepository{
			MockDeleteAdImage: func(ctx context.Context, adId string, imageId int, userId string) (string, error) {
				return imageURL, nil
			},
		}
		minioServiceMock := &mocks.MockMinioService{
			DeleteFileFunc: func(imageURL string) error {
				return nil
			},
		}

		adUseCase := NewAdUseCase(adRepoMock, minioServiceMock)

		// Вызываем функцию
		err := adUseCase.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.NoError(t, err)
	})

	// Тест ошибки при удалении изображения в репозитории
	t.Run("repository delete error", func(t *testing.T) {
		expectedErr := errors.New("repository delete error")

		adRepoMock := &mocks.MockAdRepository{
			MockDeleteAdImage: func(ctx context.Context, adId string, imageId int, userId string) (string, error) {
				return "", expectedErr
			},
		}
		minioServiceMock := &mocks.MockMinioService{
			DeleteFileFunc: func(imageURL string) error {
				return nil
			},
		}

		adUseCase := NewAdUseCase(adRepoMock, minioServiceMock)

		// Вызываем функцию
		err := adUseCase.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	// Тест ошибки при удалении файла в MinIO
	t.Run("minio delete error", func(t *testing.T) {
		adRepoMock := &mocks.MockAdRepository{
			MockDeleteAdImage: func(ctx context.Context, adId string, imageId int, userId string) (string, error) {
				return imageURL, nil
			},
		}
		minioErr := errors.New("failed to delete file from MinIO")
		minioServiceMock := &mocks.MockMinioService{
			DeleteFileFunc: func(imageURL string) error {
				return minioErr
			},
		}

		adUseCase := NewAdUseCase(adRepoMock, minioServiceMock)

		// Вызываем функцию
		err := adUseCase.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат: основная операция успешна, но возникает ошибка в логировании MinIO
		assert.NoError(t, err)
	})
}
