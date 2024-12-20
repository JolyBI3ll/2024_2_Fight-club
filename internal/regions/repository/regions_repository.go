package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type RegionRepository struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) domain.RegionRepository {
	return &RegionRepository{
		db: db,
	}
}

func (r *RegionRepository) GetVisitedRegions(ctx context.Context, userId string) ([]domain.VisitedRegions, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetVisitedRegions called", zap.String("request_id", requestID), zap.String("userID", userId))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetVisitedRegions", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetVisitedRegions", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetVisitedRegions").Observe(duration)
	}()
	var user domain.User
	var regions []domain.VisitedRegions
	if err := r.db.Where("uuid = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("User not found", zap.String("request_id", requestID), zap.String("userID", userId))
			return nil, errors.New("user not found")
		}
		logger.DBLogger.Error("Error fetching user by ID", zap.String("request_id", requestID), zap.String("userID", userId), zap.Error(err))
		return nil, errors.New("error fetching user by ID")
	}

	if err := r.db.Model(&domain.VisitedRegions{}).
		Where("\"userId\"", userId).
		Order("\"startVisitDate\" ASC").
		Find(&regions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("No regions found", zap.String("request_id", requestID), zap.String("userID", userId))
			return nil, errors.New("no regions found")
		}
		logger.DBLogger.Error("Error fetching regions", zap.String("request_id", requestID), zap.String("userID", userId), zap.Error(err))
		return nil, errors.New("error fetching regions")
	}

	logger.DBLogger.Info("Successfully fetched user by ID", zap.String("request_id", requestID), zap.String("userID", userId))
	return regions, nil
}
