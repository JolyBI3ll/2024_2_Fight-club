package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/reviews/mocks"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateReview(t *testing.T) {
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	validReview := &domain.Review{
		Title:  "Great Place!",
		Text:   "I loved staying here. The hosts were wonderful.",
		Rating: 5,
		HostID: "host123",
	}
	mockRepo.MockCreateReview = func(ctx context.Context, review *domain.Review) error {
		return nil
	}

	userID := "user123"
	ctx := context.Background()
	err := reviewUsecase.CreateReview(ctx, validReview, userID)
	assert.NoError(t, err)
}

func TestCreateReview_InvalidInput(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	invalidReview := &domain.Review{
		Title:  "Bad Title#$%",
		Text:   "Bad Text%@#",
		Rating: 10,
		HostID: "host123",
	}

	userID := "user123"
	ctx := context.Background()
	err := reviewUsecase.CreateReview(ctx, invalidReview, userID)
	assert.Error(t, err)
}

func TestUpdateReview(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	validReview := &domain.Review{
		Title:  "Updated Title",
		Text:   "Updated Text",
		Rating: 4,
	}
	mockRepo.MockUpdateReview = func(ctx context.Context, userID, hostID string, updatedReview *domain.Review) error {
		assert.Equal(t, "user123", userID)
		assert.Equal(t, "host123", hostID)
		assert.Equal(t, validReview.Title, updatedReview.Title)
		assert.Equal(t, validReview.Text, updatedReview.Text)
		assert.Equal(t, validReview.Rating, updatedReview.Rating)
		return nil
	}

	ctx := context.Background()
	err := reviewUsecase.UpdateReview(ctx, "user123", "host123", validReview)
	assert.NoError(t, err)
}

func TestUpdateReview_InvalidInput(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	invalidReview := &domain.Review{
		Title:  "Invalid Title#$%",
		Text:   "Invalid Text%@#",
		Rating: 0,
	}

	ctx := context.Background()
	err := reviewUsecase.UpdateReview(ctx, "user123", "host123", invalidReview)
	assert.Error(t, err)
}

func TestDeleteReview(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	mockRepo.MockDeleteReview = func(ctx context.Context, userID, hostID string) error {
		assert.Equal(t, "user123", userID)
		assert.Equal(t, "host123", hostID)
		return nil
	}

	ctx := context.Background()
	err := reviewUsecase.DeleteReview(ctx, "user123", "host123")
	assert.NoError(t, err)
}

func TestDeleteReview_InvalidInput(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	ctx := context.Background()
	err := reviewUsecase.DeleteReview(ctx, "user123", "!nv@lidH0stID")
	assert.Error(t, err)
}

func TestGetUserReviews(t *testing.T) {
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	expectedReviews := []domain.UserReviews{
		{ID: 1, Title: "Review 1", Text: "Text 1"},
		{ID: 2, Title: "Review 2", Text: "Text 2"},
	}
	mockRepo.MockGetUserReviews = func(ctx context.Context, userID string) ([]domain.UserReviews, error) {
		assert.Equal(t, "user123", userID)
		return expectedReviews, nil
	}

	ctx := context.Background()
	reviews, err := reviewUsecase.GetUserReviews(ctx, "user123")
	assert.NoError(t, err)
	assert.Equal(t, expectedReviews, reviews)
}

func TestGetUserReviews_InvalidInput(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	ctx := context.Background()
	_, err := reviewUsecase.GetUserReviews(ctx, "user#123") // Invalid userID
	assert.Error(t, err)
}

func TestUpdateReview_ScoreOutOfRange(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	outOfRangeReview := &domain.Review{
		Title:  "Valid Title",
		Text:   "Valid Text",
		Rating: 6, // Вне допустимого диапазона [1, 5]
	}

	ctx := context.Background()
	err := reviewUsecase.UpdateReview(ctx, "user123", "host123", outOfRangeReview)
	assert.EqualError(t, err, "score out of range")
}

func TestUpdateReview_InputExceedsCharacterLimit(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	longText := make([]byte, 1001) // Больше 1000 символов
	for i := range longText {
		longText[i] = 'a'
	}

	longReview := &domain.Review{
		Title:  "Valid Title",
		Text:   string(longText), // Текст превышает лимит
		Rating: 5,
	}

	ctx := context.Background()
	err := reviewUsecase.UpdateReview(ctx, "user123", "host123", longReview)
	assert.EqualError(t, err, "input exceeds character limit")
}

func TestCreateReview_ScoreOutOfRange(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	outOfRangeReview := &domain.Review{
		Title:  "Valid Title",
		Text:   "Valid Text",
		Rating: 0, // Вне допустимого диапазона [1, 5]
		HostID: "host123",
	}

	ctx := context.Background()
	err := reviewUsecase.CreateReview(ctx, outOfRangeReview, "user123")
	assert.EqualError(t, err, "score out of range")
}

func TestCreateReview_InputExceedsCharacterLimit(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	longTitle := make([]byte, 101) // Больше 100 символов
	for i := range longTitle {
		longTitle[i] = 'a'
	}

	longReview := &domain.Review{
		Title:  string(longTitle), // Название превышает лимит
		Text:   "Valid Text",
		Rating: 5,
		HostID: "host123",
	}

	ctx := context.Background()
	err := reviewUsecase.CreateReview(ctx, longReview, "user123")
	assert.EqualError(t, err, "input exceeds character limit")
}

func TestCreateReview_HostAndUserSame(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockRepo := &mocks.MockReviewsRepository{}
	reviewUsecase := NewReviewUsecase(mockRepo)

	review := &domain.Review{
		Title:  "Valid Title",
		Text:   "Valid Text",
		Rating: 5,
		HostID: "user123", // Хост совпадает с пользователем
	}

	ctx := context.Background()
	err := reviewUsecase.CreateReview(ctx, review, "user123")
	assert.EqualError(t, err, "host and user are the same")
}
