package mocks

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"github.com/golang-jwt/jwt"
)

type MockJwtTokenService struct {
	MockCreate            func(session_id string, tokenExpTime int64) (string, error)
	MockValidate          func(tokenString string, expectedSessionId string) (*middleware.JwtCsrfClaims, error)
	MockParseSecretGetter func(token *jwt.Token) (interface{}, error)
}

func (m *MockJwtTokenService) Create(session_id string, tokenExpTime int64) (string, error) {
	return m.MockCreate(session_id, tokenExpTime)
}

func (m *MockJwtTokenService) Validate(tokenString string, expectedSessionId string) (*middleware.JwtCsrfClaims, error) {
	return m.MockValidate(tokenString, expectedSessionId)
}

func (m *MockJwtTokenService) ParseSecretGetter(token *jwt.Token) (interface{}, error) {
	return m.MockParseSecretGetter(token)
}

type MockServiceSession struct {
	MockGetUserID      func(ctx context.Context, sessionID string) (string, error)
	MockLogoutSession  func(ctx context.Context, sessionID string) error
	MockCreateSession  func(ctx context.Context, user *domain.User) (string, error)
	MockGetSessionData func(ctx context.Context, sessionID string) (*domain.SessionData, error)
}

func (m *MockServiceSession) GetUserID(ctx context.Context, sessionID string) (string, error) {
	return m.MockGetUserID(ctx, sessionID)
}

func (m *MockServiceSession) LogoutSession(ctx context.Context, sessionID string) error {
	return m.MockLogoutSession(ctx, sessionID)
}

func (m *MockServiceSession) CreateSession(ctx context.Context, user *domain.User) (string, error) {
	return m.MockCreateSession(ctx, user)
}

func (m *MockServiceSession) GetSessionData(ctx context.Context, sessionID string) (*domain.SessionData, error) {
	return m.MockGetSessionData(ctx, sessionID)
}

type MockReviewsUsecase struct {
	MockCreateReview   func(ctx context.Context, review *domain.Review, userId string) error
	MockGetUserReviews func(ctx context.Context, userId string) ([]domain.UserReviews, error)
	MockUpdateReview   func(ctx context.Context, userID, hostID string, updatedReview *domain.Review) error
	MockDeleteReview   func(ctx context.Context, userID, hostID string) error
}

func (m *MockReviewsUsecase) CreateReview(ctx context.Context, review *domain.Review, userId string) error {
	return m.MockCreateReview(ctx, review, userId)
}

func (m *MockReviewsUsecase) GetUserReviews(ctx context.Context, userId string) ([]domain.UserReviews, error) {
	return m.MockGetUserReviews(ctx, userId)
}

func (m *MockReviewsUsecase) UpdateReview(ctx context.Context, userID, hostID string, updatedReview *domain.Review) error {
	return m.MockUpdateReview(ctx, userID, hostID, updatedReview)
}

func (m *MockReviewsUsecase) DeleteReview(ctx context.Context, userID, hostID string) error {
	return m.MockDeleteReview(ctx, userID, hostID)
}

type MockReviewsRepository struct {
	MockCreateReview   func(ctx context.Context, review *domain.Review) error
	MockGetUserReviews func(ctx context.Context, userID string) ([]domain.UserReviews, error)
	MockDeleteReview   func(ctx context.Context, userID, hostID string) error
	MockUpdateReview   func(ctx context.Context, userID, hostID string, updatedReview *domain.Review) error
}

func (m *MockReviewsRepository) CreateReview(ctx context.Context, review *domain.Review) error {
	return m.MockCreateReview(ctx, review)
}

func (m *MockReviewsRepository) GetUserReviews(ctx context.Context, userID string) ([]domain.UserReviews, error) {
	return m.MockGetUserReviews(ctx, userID)
}

func (m *MockReviewsRepository) UpdateReview(ctx context.Context, userID, hostID string, updatedReview *domain.Review) error {
	return m.MockUpdateReview(ctx, userID, hostID, updatedReview)
}

func (m *MockReviewsRepository) DeleteReview(ctx context.Context, userID, hostID string) error {
	return m.MockDeleteReview(ctx, userID, hostID)
}
