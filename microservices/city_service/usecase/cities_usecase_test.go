package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/microservices/city_service/mocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

	cityUsecase := NewCityUseCase(mockRepo)

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

	cityUsecase := NewCityUseCase(mockRepo)

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

func TestCityUseCase_GetOneCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Контекст для тестов.
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Инициализируем ожидаемое значение для успешного выполнения.
		expectedCity := domain.City{
			EnTitle: "New York",
			Title:   "Нью-Йорк",
		}

		// Создаем мок репозитория.
		mockRepo := &mocks.MockCitiesRepository{
			MockGetCityByEnName: func(ctx context.Context, cityEnName string) (domain.City, error) {
				if cityEnName == "New York" {
					return expectedCity, nil
				}
				return domain.City{}, errors.New("city not found")
			},
		}

		// Инициализируем CityUseCase с нашим моком.
		cityUC := NewCityUseCase(mockRepo)

		// Вызываем тестируемый метод.
		city, err := cityUC.GetOneCity(ctx, "New York")

		// Проверяем, что не возникло ошибок и город совпадает с ожидаемым.
		assert.NoError(t, err)
		assert.Equal(t, expectedCity, city)
	})

	t.Run("City Not Found", func(t *testing.T) {
		// Создаем мок, возвращающий ошибку, если город не найден.
		mockRepo := &mocks.MockCitiesRepository{
			MockGetCityByEnName: func(ctx context.Context, cityEnName string) (domain.City, error) {
				return domain.City{}, errors.New("city not found")
			},
		}

		cityUC := NewCityUseCase(mockRepo)

		// Проверка на случай, если город не найден.
		city, err := cityUC.GetOneCity(ctx, "Unknown City")

		assert.Error(t, err)
		assert.Equal(t, "city not found", err.Error())
		assert.Equal(t, domain.City{}, city)
	})

	t.Run("Repository Error", func(t *testing.T) {
		// Мок репозитория, который всегда возвращает ошибку.
		mockRepo := &mocks.MockCitiesRepository{
			MockGetCityByEnName: func(ctx context.Context, cityEnName string) (domain.City, error) {
				return domain.City{}, errors.New("repository error")
			},
		}

		cityUC := NewCityUseCase(mockRepo)

		// Проверка на случай, если репозиторий вернул ошибку.
		city, err := cityUC.GetOneCity(ctx, "Any City")

		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())
		assert.Equal(t, domain.City{}, city)
	})
}
