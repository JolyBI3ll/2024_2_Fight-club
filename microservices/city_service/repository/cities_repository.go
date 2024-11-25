package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type cityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) domain.CityRepository {
	return &cityRepository{
		db: db,
	}
}

func (c cityRepository) GetCities(ctx context.Context) ([]domain.City, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetCities called", zap.String("request_id", requestID))

	var cities []domain.City
	if err := c.db.Find(&cities).Error; err != nil {
		logger.DBLogger.Error("Error fetching all cities", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("error fetching all cities")
	}

	logger.DBLogger.Info("Successfully fetched all cities", zap.String("request_id", requestID), zap.Int("count", len(cities)))
	return cities, nil
}

func (c cityRepository) GetCityByEnName(ctx context.Context, cityEnName string) (domain.City, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetCityByEnName called", zap.String("request_id", requestID))

	var city domain.City
	if err := c.db.First(&city, "\"enTitle\" = ?", cityEnName).Error; err != nil {
		logger.DBLogger.Error("Error fetching city", zap.String("request_id", requestID), zap.Error(err))
		return domain.City{}, errors.New("error fetching city")
	}

	logger.DBLogger.Info("Successfully fetched all cities", zap.String("request_id", requestID))
	return city, nil
}
