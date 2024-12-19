package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"go.uber.org/zap"
	"regexp"
)

type RegionUsecase interface {
	GetVisitedRegions(ctx context.Context, userId string) ([]domain.VisitedRegions, error)
}

type regionUsecase struct {
	repository domain.RegionRepository
}

func NewRegionUsecase(repository domain.RegionRepository) RegionUsecase {
	return &regionUsecase{
		repository: repository,
	}
}

func (r *regionUsecase) GetVisitedRegions(ctx context.Context, userId string) ([]domain.VisitedRegions, error) {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(userId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return nil, errors.New("input contains invalid characters")
	}

	if len(userId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return nil, errors.New("input exceeds character limit")
	}

	reviews, err := r.repository.GetVisitedRegions(ctx, userId)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
