package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) domain.ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

func (r *ReviewRepository) CreateReview(ctx context.Context, review *domain.Review) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("CreateReview called", zap.String("HostId", review.HostID), zap.String("request_id", requestID))

	if err := r.db.Where("uuid = ?", review.HostID).First(&domain.User{}).Error; err != nil {
		logger.DBLogger.Error("Error finding host", zap.String("userId", review.HostID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("Error finding host")
	}

	if err := r.db.Create(&review).Error; err != nil {
		logger.DBLogger.Error("Error creating review", zap.String("userId", review.UserID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("Error creating review")
	}

	if err := r.updateHostScore(ctx, review.HostID); err != nil {
		logger.DBLogger.Error("Error updating host score", zap.String("hostId", review.HostID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("Error updating host score")
	}

	return nil
}

func (r *ReviewRepository) GetUserReviews(ctx context.Context, userId string) ([]domain.UserReviews, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetUserReviews called", zap.String("request_id", requestID), zap.String("userID", userId))

	var user domain.User
	var reviews []domain.UserReviews
	if err := r.db.Where("uuid = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("User not found", zap.String("request_id", requestID), zap.String("userID", userId))
			return nil, errors.New("user not found")
		}
		logger.DBLogger.Error("Error fetching user by ID", zap.String("request_id", requestID), zap.String("userID", userId), zap.Error(err))
		return nil, err
	}

	if err := r.db.Model(&domain.Review{}).
		Select("reviews.*, users.avatar as \"UserAvatar\", users.name as \"UserName\"").
		Joins("JOIN users ON reviews.\"userId\" = users.uuid").
		Where("reviews.\"hostId\" = ?", userId).
		Order("reviews.\"createdAt\" ASC").
		Find(&reviews).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("No reviews found", zap.String("request_id", requestID), zap.String("userID", userId))
			return nil, nil
		}
		logger.DBLogger.Error("Error fetching reviews", zap.String("request_id", requestID), zap.String("userID", userId), zap.Error(err))
		return nil, err
	}

	logger.DBLogger.Info("Successfully fetched user by ID", zap.String("request_id", requestID), zap.String("userID", userId))
	return reviews, nil
}

func (r *ReviewRepository) updateHostScore(ctx context.Context, hostID string) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("updateHostScore called", zap.String("HostId", hostID), zap.String("request_id", requestID))
	var reviews []domain.Review
	if err := r.db.Where("\"hostId\" = ?", hostID).Find(&reviews).Error; err != nil {
		logger.DBLogger.Error("Failed to fetch reviews for host", zap.String("hostId", hostID), zap.String("request_id", requestID), zap.Error(err))
		return fmt.Errorf("failed to fetch reviews for host: %w", err)
	}

	if len(reviews) == 0 {
		return r.db.Model(&domain.User{}).Where("uuid = ?", hostID).Update("score", 0).Error
	}

	var totalScore int
	for _, review := range reviews {
		totalScore += review.Rating
	}
	averageScore := float64(totalScore) / float64(len(reviews))
	averageScore = math.Round(averageScore*10) / 10

	if err := r.db.Model(&domain.User{}).Where("uuid = ?", hostID).Update("score", averageScore).Error; err != nil {
		return fmt.Errorf("failed to update host score: %w", err)
	}

	return nil
}
