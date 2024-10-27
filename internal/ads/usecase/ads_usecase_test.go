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

func TestGetAllPlaces(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{
		GetAllPlacesFunc: func(ctx context.Context, filter domain.AdFilter) ([]domain.Ad, error) {
			return []domain.Ad{{ID: "1", LocationMain: "Test Location Main"}}, nil
		},
	}

	useCase := NewAdUseCase(mockRepo, nil)
	ads, err := useCase.GetAllPlaces(context.Background(), domain.AdFilter{})

	assert.NoError(t, err)
	assert.NotEmpty(t, ads)
	assert.Equal(t, "Test Location Main", ads[0].LocationMain)
}

func TestGetOnePlace(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{
		GetPlaceByIdFunc: func(ctx context.Context, adId string) (domain.Ad, error) {
			if adId == "1" {
				return domain.Ad{ID: "1", LocationMain: "Test Location Main"}, nil
			}
			return domain.Ad{}, errors.New("ad not found")
		},
	}

	useCase := NewAdUseCase(mockRepo, nil)
	ad, err := useCase.GetOnePlace(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, "Test Location Main", ad.LocationMain)

	_, err = useCase.GetOnePlace(context.Background(), "999")
	assert.Error(t, err)
	assert.Equal(t, "ad not found", err.Error())
}

func TestCreatePlace(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{
		CreatePlaceFunc: func(ctx context.Context, place *domain.Ad) error {
			return nil
		},
		SavePlaceFunc: func(ctx context.Context, place *domain.Ad) error {
			return nil
		},
	}

	mockMinio := &mocks.MockMinioService{
		UploadFileFunc: func(file *multipart.FileHeader, path string) (string, error) {
			return "uploaded/path", nil
		},
		DeleteFileFunc: func(path string) error {
			return nil
		},
	}

	useCase := NewAdUseCase(mockRepo, mockMinio)
	err := useCase.CreatePlace(context.Background(), &domain.Ad{ID: "1"}, []*multipart.FileHeader{})

	assert.NoError(t, err)
}

func TestUpdatePlace(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{
		GetPlaceByIdFunc: func(ctx context.Context, adId string) (domain.Ad, error) {
			return domain.Ad{ID: adId, Images: []string{"old/image/path"}}, nil
		},
		UpdatePlaceFunc: func(ctx context.Context, place *domain.Ad, adId, userId string) error {
			return nil
		},
	}

	mockMinio := &mocks.MockMinioService{
		UploadFileFunc: func(file *multipart.FileHeader, path string) (string, error) {
			return "new/image/path", nil
		},
		DeleteFileFunc: func(path string) error {
			return nil
		},
	}

	useCase := NewAdUseCase(mockRepo, mockMinio)
	err := useCase.UpdatePlace(context.Background(), &domain.Ad{ID: "1"}, "1", "1", []*multipart.FileHeader{})

	assert.NoError(t, err)
}

func TestDeletePlace(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{
		GetPlaceByIdFunc: func(ctx context.Context, adId string) (domain.Ad, error) {
			return domain.Ad{ID: adId, Images: []string{"old/image/path"}}, nil
		},
		DeletePlaceFunc: func(ctx context.Context, adId, userId string) error {
			return nil
		},
	}

	mockMinio := &mocks.MockMinioService{
		DeleteFileFunc: func(path string) error {
			return nil
		},
	}

	useCase := NewAdUseCase(mockRepo, mockMinio)
	err := useCase.DeletePlace(context.Background(), "1", "1")

	assert.NoError(t, err)
}
