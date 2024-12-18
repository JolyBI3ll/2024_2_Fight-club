package contoller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/reviews/mocks"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"bytes"
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateReview(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockJwtService := &mocks.MockJwtTokenService{}
	mockSessionService := &mocks.MockServiceSession{}
	mockReviewUsecase := &mocks.MockReviewsUsecase{}

	handler := NewReviewHandler(mockReviewUsecase, mockSessionService, mockJwtService)

	t.Run("Successful Review Creation", func(t *testing.T) {
		mockJwtService.MockValidate = func(tokenString string, expectedSessionId string) (*middleware.JwtCsrfClaims, error) {
			return &middleware.JwtCsrfClaims{}, nil
		}
		mockSessionService.MockGetUserID = func(ctx context.Context, sessionID string) (string, error) {
			return "test-user-id", nil
		}
		mockReviewUsecase.MockCreateReview = func(ctx context.Context, review *domain.Review, userId string) error {
			assert.Equal(t, "test-user-id", userId)
			assert.Equal(t, "Test Review", review.Title)
			return nil
		}

		reviewPayload := `{"title":"Test Review","text":"Review Content","host_id":"host1","user_id":"user1"}`
		request := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader([]byte(reviewPayload)))
		request.Header.Set("Cookie", "session_id=test-session-id")
		request.Header.Set("X-CSRF-Token", "Bearer valid-token")

		responseRecorder := httptest.NewRecorder()
		handler.CreateReview(responseRecorder, request)

		assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	})

	t.Run("Missing CSRF Token", func(t *testing.T) {
		reviewPayload := `{"title":"Test Review","text":"Review Content"}`
		request := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader([]byte(reviewPayload)))
		request.Header.Set("Cookie", "session_id=test-session-id")

		responseRecorder := httptest.NewRecorder()
		handler.CreateReview(responseRecorder, request)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
	})

	t.Run("Invalid JWT Token", func(t *testing.T) {
		mockJwtService.MockValidate = func(tokenString string, expectedSessionId string) (*middleware.JwtCsrfClaims, error) {
			return nil, errors.New("invalid JWT token")
		}

		reviewPayload := `{"title":"Test Review","text":"Review Content"}`
		request := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader([]byte(reviewPayload)))
		request.Header.Set("Cookie", "session_id=test-session-id")
		request.Header.Set("X-CSRF-Token", "Bearer invalid-token")

		responseRecorder := httptest.NewRecorder()
		handler.CreateReview(responseRecorder, request)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
	})

	t.Run("Failed to Get User ID from Session", func(t *testing.T) {
		mockJwtService.MockValidate = func(tokenString string, expectedSessionId string) (*middleware.JwtCsrfClaims, error) {
			return &middleware.JwtCsrfClaims{}, nil
		}
		mockSessionService.MockGetUserID = func(ctx context.Context, sessionID string) (string, error) {
			return "", errors.New("failed to get user ID")
		}

		reviewPayload := `{"title":"Test Review","text":"Review Content"}`
		request := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader([]byte(reviewPayload)))
		request.Header.Set("Cookie", "session_id=test-session-id")
		request.Header.Set("X-CSRF-Token", "Bearer valid-token")

		responseRecorder := httptest.NewRecorder()
		handler.CreateReview(responseRecorder, request)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	})
	t.Run("Failed to Get User ID from Session", func(t *testing.T) {
		mockJwtService.MockValidate = func(tokenString string, expectedSessionId string) (*middleware.JwtCsrfClaims, error) {
			return &middleware.JwtCsrfClaims{}, nil
		}
		mockSessionService.MockGetUserID = func(ctx context.Context, sessionID string) (string, error) {
			return "", errors.New("failed to get user ID")
		}

		reviewPayload := `{"title":"Test Review","text":"Review Content"}`
		request := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader([]byte(reviewPayload)))
		request.Header.Set("Cookie", "session_id=test-session-id")
		request.Header.Set("X-CSRF-Token", "Bearer valid-token")

		responseRecorder := httptest.NewRecorder()
		handler.CreateReview(responseRecorder, request)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	})
}

func TestGetUserReviews(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockJwtService := &mocks.MockJwtTokenService{}
	mockSessionService := &mocks.MockServiceSession{}
	mockReviewUsecase := &mocks.MockReviewsUsecase{}

	handler := NewReviewHandler(mockReviewUsecase, mockSessionService, mockJwtService)

	t.Run("Successful GetUserReviews", func(t *testing.T) {
		mockReviewUsecase.MockGetUserReviews = func(ctx context.Context, userId string) ([]domain.UserReviews, error) {
			assert.Equal(t, "", userId)
			return []domain.UserReviews{
				{
					HostID: "host1",
					Title:  "Test Review 1",
				},
				{
					HostID: "host2",
					Title:  "Test Review 2",
				},
			}, nil
		}

		request := httptest.NewRequest(http.MethodGet, "/reviews/{userId}", nil)
		responseRecorder := httptest.NewRecorder()

		handler.GetUserReviews(responseRecorder, request)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
	})

	t.Run("Error from GetUserReviews Usecase", func(t *testing.T) {
		mockReviewUsecase.MockGetUserReviews = func(ctx context.Context, userId string) ([]domain.UserReviews, error) {
			assert.Equal(t, "", userId)
			return nil, errors.New("database error")
		}

		request := httptest.NewRequest(http.MethodGet, "/reviews/{userId}", nil)
		responseRecorder := httptest.NewRecorder()

		handler.GetUserReviews(responseRecorder, request)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	})
}

func TestDeleteReview(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockJwtService := &mocks.MockJwtTokenService{}
	mockSessionService := &mocks.MockServiceSession{}
	mockReviewUsecase := &mocks.MockReviewsUsecase{}
	handler := NewReviewHandler(mockReviewUsecase, mockSessionService, mockJwtService)

	t.Run("Successful DeleteReview", func(t *testing.T) {
		mockJwtService.MockValidate = func(tokenString string, expectedSessionId string) (*middleware.JwtCsrfClaims, error) {
			return &middleware.JwtCsrfClaims{}, nil
		}
		mockSessionService.MockGetUserID = func(ctx context.Context, sessionID string) (string, error) {
			return "test-user-id", nil
		}
		mockReviewUsecase.MockDeleteReview = func(ctx context.Context, userId, hostId string) error {
			assert.Equal(t, "test-user-id", userId)
			assert.Equal(t, "test-host-id", hostId)
			return nil
		}

		request := httptest.NewRequest(http.MethodDelete, "/reviews/{adId}", nil)
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})
		request.Header.Set("Cookie", "session_id=test-session-id")
		request.Header.Set("X-CSRF-Token", "Bearer valid-token")

		rr := httptest.NewRecorder()
		handler.DeleteReview(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "deleted successfully")
	})

	t.Run("Missing CSRF Token", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "/reviews/{adId}", nil)
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})
		request.Header.Set("Cookie", "session_id=test-session-id")

		rr := httptest.NewRecorder()
		handler.DeleteReview(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid JWT Token", func(t *testing.T) {
		mockJwtService.MockValidate = func(tokenString string, expectedSessionId string) (*middleware.JwtCsrfClaims, error) {
			return nil, errors.New("invalid JWT token")
		}

		request := httptest.NewRequest(http.MethodDelete, "/reviews/{adId}", nil)
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})
		request.Header.Set("Cookie", "session_id=test-session-id")
		request.Header.Set("X-CSRF-Token", "Bearer invalid-token")

		rr := httptest.NewRecorder()
		handler.DeleteReview(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Session ID Extraction Error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "/reviews/{adId}", nil)
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})

		rr := httptest.NewRecorder()
		handler.DeleteReview(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Error Getting UserID from Session", func(t *testing.T) {
		mockJwtService.MockValidate = func(tokenString string, expectedSessionId string) (*middleware.JwtCsrfClaims, error) {
			return &middleware.JwtCsrfClaims{}, nil
		}
		mockSessionService.MockGetUserID = func(ctx context.Context, sessionID string) (string, error) {
			return "", errors.New("session error")
		}

		request := httptest.NewRequest(http.MethodDelete, "/reviews/{adId}", nil)
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})
		request.Header.Set("Cookie", "session_id=test-session-id")
		request.Header.Set("X-CSRF-Token", "Bearer valid-token")

		rr := httptest.NewRecorder()
		handler.DeleteReview(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Error Deleting Review", func(t *testing.T) {
		mockJwtService.MockValidate = func(tokenString string, expectedSessionId string) (*middleware.JwtCsrfClaims, error) {
			return &middleware.JwtCsrfClaims{}, nil
		}
		mockSessionService.MockGetUserID = func(ctx context.Context, sessionID string) (string, error) {
			return "test-user-id", nil
		}
		mockReviewUsecase.MockDeleteReview = func(ctx context.Context, userId, hostId string) error {
			return errors.New("delete error")
		}

		request := httptest.NewRequest(http.MethodDelete, "/reviews/{hostId}", nil)
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})
		request.Header.Set("Cookie", "session_id=test-session-id")
		request.Header.Set("X-CSRF-Token", "Bearer valid-token")

		rr := httptest.NewRecorder()
		handler.DeleteReview(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestUpdateReview(t *testing.T) {
	require.NoError(t, logger.InitLoggers())
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockJwtService := &mocks.MockJwtTokenService{}
	mockSessionService := &mocks.MockServiceSession{}
	mockReviewUsecase := &mocks.MockReviewsUsecase{}
	handler := NewReviewHandler(mockReviewUsecase, mockSessionService, mockJwtService)

	t.Run("Successful UpdateReview", func(t *testing.T) {
		mockJwtService.MockValidate = func(tokenString string, sessionID string) (*middleware.JwtCsrfClaims, error) {
			return &middleware.JwtCsrfClaims{}, nil
		}
		mockSessionService.MockGetUserID = func(ctx context.Context, sessionID string) (string, error) {
			return "test-user-id", nil
		}
		mockReviewUsecase.MockUpdateReview = func(ctx context.Context, userID, hostID string, review *domain.Review) error {
			assert.Equal(t, "test-host-id", hostID)
			assert.Equal(t, "test-user-id", userID)
			assert.Equal(t, "Test Title", review.Title)
			return nil
		}

		updatedReview := `{"title":"Test Title","text":"Updated text"}`

		request := httptest.NewRequest(http.MethodPut, "/reviews/{hostId}", bytes.NewBufferString(updatedReview))
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})
		request.Header.Set("X-CSRF-Token", "Bearer valid-token")
		request.Header.Set("Cookie", "session_id=test-session-id")

		rr := httptest.NewRecorder()
		handler.UpdateReview(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "updated successfully")
	})

	t.Run("Missing CSRF Token", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPut, "/reviews/{hostId}", nil)
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})
		request.Header.Set("Cookie", "session_id=test-session-id")

		rr := httptest.NewRecorder()
		handler.UpdateReview(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid JWT Token", func(t *testing.T) {
		mockJwtService.MockValidate = func(tokenString string, sessionID string) (*middleware.JwtCsrfClaims, error) {
			return nil, errors.New("invalid token")
		}

		request := httptest.NewRequest(http.MethodPut, "/reviews/{hostId}", nil)
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})
		request.Header.Set("X-CSRF-Token", "Bearer invalid-token")
		request.Header.Set("Cookie", "session_id=test-session-id")

		rr := httptest.NewRecorder()
		handler.UpdateReview(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Session ID Error", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPut, "/reviews/{hostId}", nil)
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})

		rr := httptest.NewRecorder()
		handler.UpdateReview(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Failed to get UserID", func(t *testing.T) {
		mockJwtService.MockValidate = func(tokenString string, sessionID string) (*middleware.JwtCsrfClaims, error) {
			return &middleware.JwtCsrfClaims{}, nil
		}
		mockSessionService.MockGetUserID = func(ctx context.Context, sessionID string) (string, error) {
			return "", errors.New("failed to get user ID")
		}

		request := httptest.NewRequest(http.MethodPut, "/reviews/{hostId}", nil)
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})
		request.Header.Set("X-CSRF-Token", "Bearer valid-token")
		request.Header.Set("Cookie", "session_id=test-session-id")

		rr := httptest.NewRecorder()
		handler.UpdateReview(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Failed Unmarshal", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPut, "/reviews/{hostId}", bytes.NewBufferString("invalid-json"))
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})
		request.Header.Set("X-CSRF-Token", "Bearer valid-token")
		request.Header.Set("Cookie", "session_id=test-session-id")

		rr := httptest.NewRecorder()
		handler.UpdateReview(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("UpdateReview Usecase Error", func(t *testing.T) {
		mockJwtService.MockValidate = func(tokenString string, sessionID string) (*middleware.JwtCsrfClaims, error) {
			return &middleware.JwtCsrfClaims{}, nil
		}
		mockSessionService.MockGetUserID = func(ctx context.Context, sessionID string) (string, error) {
			return "test-user-id", nil
		}
		mockReviewUsecase.MockUpdateReview = func(ctx context.Context, userID, hostID string, review *domain.Review) error {
			return errors.New("update error")
		}

		request := httptest.NewRequest(http.MethodPut, "/reviews/{hostId}", bytes.NewBufferString(`{"title":"Test"}`))
		request = mux.SetURLVars(request, map[string]string{"hostId": "test-host-id"})
		request.Header.Set("X-CSRF-Token", "Bearer valid-token")
		request.Header.Set("Cookie", "session_id=test-session-id")

		rr := httptest.NewRecorder()
		handler.UpdateReview(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
