package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"go.uber.org/zap"
	"regexp"
	"time"
)

type ReviewUsecase interface {
	CreateReview(ctx context.Context, review *domain.Review, userId string) error
	GetUserReviews(ctx context.Context, userId string) ([]domain.UserReviews, error)
	UpdateReview(ctx context.Context, userID, hostID string, updatedReview *domain.Review) error
	DeleteReview(ctx context.Context, userID, hostID string) error
}

type reviewUsecase struct {
	repository domain.ReviewRepository
}

func (r *reviewUsecase) UpdateReview(ctx context.Context, userID, hostID string, updatedReview *domain.Review) error {
	const maxLenTitle = 100
	const maxLenText = 1000
	const minScore, maxScore = 1, 5
	requestID := middleware.GetRequestID(ctx)
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9@.,\s\-!?:;_/()]*$`)
	if !validCharPattern.MatchString(updatedReview.Title) ||
		!validCharPattern.MatchString(updatedReview.Text) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return errors.New("input contains invalid characters")
	}

	if updatedReview.Rating < minScore || updatedReview.Rating > maxScore {
		logger.AccessLogger.Warn("Score out of range", zap.String("request_id", requestID))
		return errors.New("score out of range")
	}

	if len(updatedReview.Title) > maxLenTitle || len(updatedReview.Text) > maxLenText {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return errors.New("input exceeds character limit")
	}
	err := r.repository.UpdateReview(ctx, userID, hostID, updatedReview)
	if err != nil {
		return err
	}
	return nil
}

func (r *reviewUsecase) DeleteReview(ctx context.Context, userID, hostID string) error {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(hostID) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return errors.New("input contains invalid characters")
	}

	if len(hostID) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return errors.New("input exceeds character limit")
	}

	err := r.repository.DeleteReview(ctx, userID, hostID)
	if err != nil {
		return err
	}
	return nil
}

func NewReviewUsecase(repository domain.ReviewRepository) ReviewUsecase {
	return &reviewUsecase{
		repository: repository,
	}
}

func (r *reviewUsecase) CreateReview(ctx context.Context, review *domain.Review, userId string) error {
	const maxLenTitle = 100
	const maxLenText = 1000
	const minScore, maxScore = 1, 5
	requestID := middleware.GetRequestID(ctx)
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9@.,\s\-!?:;_/()]*$`)
	if !validCharPattern.MatchString(review.Title) ||
		!validCharPattern.MatchString(review.Text) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return errors.New("input contains invalid characters")
	}

	if review.Rating < minScore || review.Rating > maxScore {
		logger.AccessLogger.Warn("Score out of range", zap.String("request_id", requestID))
		return errors.New("score out of range")
	}

	if len(review.Title) > maxLenTitle || len(review.Text) > maxLenText {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return errors.New("input exceeds character limit")
	}
	if review.HostID == userId {
		return errors.New("host and user are the same")
	}

	review.UserID = userId
	review.CreatedAt = time.Now()
	err := r.repository.CreateReview(ctx, review)
	if err != nil {
		return err
	}
	return nil
}

func (r *reviewUsecase) GetUserReviews(ctx context.Context, userId string) ([]domain.UserReviews, error) {
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

	reviews, err := r.repository.GetUserReviews(ctx, userId)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
