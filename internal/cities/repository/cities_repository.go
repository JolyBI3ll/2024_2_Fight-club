package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
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
		return nil, err
	}

	logger.DBLogger.Info("Successfully fetched all cities", zap.String("request_id", requestID), zap.Int("count", len(cities)))
	return cities, nil
}
