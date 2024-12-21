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
	"math"
	"time"
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
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("CreateReview called", zap.String("HostId", review.HostID), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("CreateReview", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("CreateReview", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("CreateReview").Observe(duration)
	}()
	if err := r.db.Where("uuid = ?", review.HostID).First(&domain.User{}).Error; err != nil {
		logger.DBLogger.Error("Error finding host", zap.String("userId", review.HostID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error finding host")
	}

	//Проверка существует ли отзыва
	var query domain.Review
	if err := r.db.Where("\"userId\" = ? AND \"hostId\" = ?", review.UserID, review.HostID).First(&query).Error; err == nil {
		logger.DBLogger.Warn("Review already exists", zap.String("userId", review.UserID), zap.String("hostId", review.HostID), zap.String("request_id", requestID))
		return errors.New("review already exist")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.DBLogger.Error("Error finding review",
			zap.String("userId", review.UserID),
			zap.String("hostId", review.HostID),
			zap.String("request_id", requestID),
			zap.Error(err))
		return errors.New("error finding review")
	}

	if err := r.db.Create(&review).Error; err != nil {
		logger.DBLogger.Error("Error creating review", zap.String("userId", review.UserID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error creating review")
	}

	if err := r.updateHostScore(ctx, review.HostID); err != nil {
		logger.DBLogger.Error("Error updating host score", zap.String("hostId", review.HostID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error updating host score")
	}

	return nil
}

func (r *ReviewRepository) GetUserReviews(ctx context.Context, userId string) ([]domain.UserReviews, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetUserReviews called", zap.String("request_id", requestID), zap.String("userID", userId))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetUserReviews", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetUserReviews", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetUserReviews").Observe(duration)
	}()
	var user domain.User
	var reviews []domain.UserReviews
	if err := r.db.Where("uuid = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("User not found", zap.String("request_id", requestID), zap.String("userID", userId))
			return nil, errors.New("user not found")
		}
		logger.DBLogger.Error("Error fetching user by ID", zap.String("request_id", requestID), zap.String("userID", userId), zap.Error(err))
		return nil, errors.New("error fetching user by ID")
	}

	if err := r.db.Model(&domain.Review{}).
		Select("reviews.*, users.avatar as \"UserAvatar\", users.name as \"UserName\"").
		Joins("JOIN users ON reviews.\"userId\" = users.uuid").
		Where("reviews.\"hostId\" = ?", userId).
		Order("reviews.\"createdAt\" DESC").
		Find(&reviews).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("No reviews found", zap.String("request_id", requestID), zap.String("userID", userId))
			return nil, errors.New("no reviews found")
		}
		logger.DBLogger.Error("Error fetching reviews", zap.String("request_id", requestID), zap.String("userID", userId), zap.Error(err))
		return nil, errors.New("error fetching reviews")
	}

	logger.DBLogger.Info("Successfully fetched user by ID", zap.String("request_id", requestID), zap.String("userID", userId))
	return reviews, nil
}

func (r *ReviewRepository) DeleteReview(ctx context.Context, userID, hostID string) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("DeleteReview called", zap.String("userID", userID), zap.String("hostID", hostID), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("DeleteReview", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("DeleteReview", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("DeleteReview").Observe(duration)
	}()
	var review domain.Review
	if err := r.db.Where("\"userId\" = ? AND \"hostId\" = ?", userID, hostID).First(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("Review not found", zap.String("userID", userID), zap.String("hostID", hostID), zap.String("request_id", requestID))
			return errors.New("review not found")
		}
		logger.DBLogger.Error("Error finding review", zap.String("userID", userID), zap.String("hostID", hostID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error finding review")
	}

	if err := r.db.Delete(&review).Error; err != nil {
		logger.DBLogger.Error("Error deleting review", zap.String("userID", userID), zap.String("hostID", hostID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error deleting review")
	}

	if err := r.updateHostScore(ctx, hostID); err != nil {
		logger.DBLogger.Error("Error updating host score after review deletion", zap.String("hostID", hostID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error updating host score")
	}

	logger.DBLogger.Info("Review successfully deleted", zap.String("userID", userID), zap.String("hostID", hostID), zap.String("request_id", requestID))
	return nil
}

func (r *ReviewRepository) UpdateReview(ctx context.Context, userID, hostID string, updatedReview *domain.Review) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("UpdateReview called", zap.String("userID", userID), zap.String("hostID", hostID), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("UpdateReview", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("UpdateReview", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("UpdateReview").Observe(duration)
	}()
	var existingReview domain.Review
	if err := r.db.Where("\"userId\" = ? AND \"hostId\" = ?", userID, hostID).First(&existingReview).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("Review not found", zap.String("userID", userID), zap.String("hostID", hostID), zap.String("request_id", requestID))
			return errors.New("review not found")
		}
		logger.DBLogger.Error("Error finding review", zap.String("userID", userID), zap.String("hostID", hostID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error finding review")
	}

	// Обновляем только текст и рейтинг.
	existingReview.Title = updatedReview.Title
	existingReview.Text = updatedReview.Text
	existingReview.Rating = updatedReview.Rating

	if err := r.db.Save(&existingReview).Error; err != nil {
		logger.DBLogger.Error("Error updating review", zap.String("userID", userID), zap.String("hostID", hostID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error updating review")
	}

	if err := r.updateHostScore(ctx, hostID); err != nil {
		logger.DBLogger.Error("Error updating host score after review update", zap.String("hostID", hostID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error updating host score")
	}

	logger.DBLogger.Info("Review successfully updated", zap.String("userID", userID), zap.String("hostID", hostID), zap.String("request_id", requestID))
	return nil
}

func (r *ReviewRepository) updateHostScore(ctx context.Context, hostID string) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("updateHostScore called", zap.String("HostId", hostID), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("updateHostScore", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("updateHostScore", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("updateHostScore").Observe(duration)
	}()
	var reviews []domain.Review
	if err := r.db.Where("\"hostId\" = ?", hostID).Find(&reviews).Error; err != nil {
		logger.DBLogger.Error("Failed to fetch reviews for host", zap.String("hostId", hostID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("failed to fetch reviews for host")
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
		return errors.New("failed to update host score")
	}
	return nil
}
