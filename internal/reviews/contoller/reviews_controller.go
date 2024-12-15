package contoller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/reviews/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"errors"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ReviewHandler struct {
	usecase        usecase.ReviewUsecase
	sessionService session.InterfaceSession
	jwtToken       middleware.JwtTokenService
}

func NewReviewHandler(usecase usecase.ReviewUsecase, sessionService session.InterfaceSession, jwtToken middleware.JwtTokenService) *ReviewHandler {
	return &ReviewHandler{
		usecase:        usecase,
		sessionService: sessionService,
		jwtToken:       jwtToken,
	}
}

func (rh *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	sanitizer := bluemonday.UGCPolicy()
	defer cancel()
	var err error
	statusCode := http.StatusCreated
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusCreated {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received CreateReview request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	authHeader := r.Header.Get("X-CSRF-Token")
	if authHeader == "" {
		logger.AccessLogger.Warn("Missing X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("missing X-CSRF-Token header")),
		)
		err = errors.New("missing X-CSRF-Token header")
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	tokenString := authHeader[len("Bearer "):]
	_, err = rh.jwtToken.Validate(tokenString, sessionID)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, errors.New("invalid JWT token"), requestID)
		return
	}

	userId, err := rh.sessionService.GetUserID(ctx, sessionID)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get user ID", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	var review domain.Review
	if err := easyjson.UnmarshalFromReader(r.Body, &review); err != nil {
		logger.AccessLogger.Warn("Failed to unmarshal review", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	review.Text = sanitizer.Sanitize(review.Text)
	review.Title = sanitizer.Sanitize(review.Title)
	review.HostID = sanitizer.Sanitize(review.HostID)
	review.UserID = sanitizer.Sanitize(review.UserID)

	err = rh.usecase.CreateReview(ctx, &review, userId)
	if err != nil {
		logger.AccessLogger.Warn("Failed to create review", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	w.WriteHeader(http.StatusCreated)
	body := domain.ReviewBody{
		Review: review,
	}
	if _, err = easyjson.MarshalToWriter(&body, w); err != nil {
		logger.AccessLogger.Warn("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed CreateReview request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (rh *ReviewHandler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())

	defer cancel()
	var err error
	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()
	userId := mux.Vars(r)["userId"]
	sanitizer := bluemonday.UGCPolicy()
	userId = sanitizer.Sanitize(userId)

	logger.AccessLogger.Info("Received GetUserReviews request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var reviews domain.UserReviewsList

	reviews, err = rh.usecase.GetUserReviews(ctx, userId)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get user reviews", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = easyjson.MarshalToWriter(&reviews, w); err != nil {
		logger.AccessLogger.Warn("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetUserReviews request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK))
}

func (rh *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	sanitizer := bluemonday.UGCPolicy()
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()
	var err error
	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	hostId := mux.Vars(r)["hostId"]
	hostId = sanitizer.Sanitize(hostId)

	logger.AccessLogger.Info("Received DeleteReview request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	authHeader := r.Header.Get("X-CSRF-Token")
	if authHeader == "" {
		logger.AccessLogger.Warn("Missing X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("missing X-CSRF-Token header")),
		)
		err = errors.New("missing X-CSRF-Token header")
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	tokenString := authHeader[len("Bearer "):]
	_, err = rh.jwtToken.Validate(tokenString, sessionID)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, errors.New("invalid JWT token"), requestID)
		return
	}

	userId, err := rh.sessionService.GetUserID(ctx, sessionID)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get user ID", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	err = rh.usecase.DeleteReview(ctx, userId, hostId)
	if err != nil {
		logger.AccessLogger.Warn("Failed to delete review", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	response := domain.ResponseMessage{Message: "deleted successfully"}
	if _, err = easyjson.MarshalToWriter(&response, w); err != nil {
		logger.AccessLogger.Warn("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed DeleteReview request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK))
}

func (rh *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()
	var err error
	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()
	sanitizer := bluemonday.UGCPolicy()

	hostId := mux.Vars(r)["hostId"]
	hostId = sanitizer.Sanitize(hostId)

	logger.AccessLogger.Info("Received UpdateReview request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	authHeader := r.Header.Get("X-CSRF-Token")
	if authHeader == "" {
		logger.AccessLogger.Warn("Missing X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("missing X-CSRF-Token header")),
		)
		err = errors.New("missing X-CSRF-Token header")
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	tokenString := authHeader[len("Bearer "):]
	_, err = rh.jwtToken.Validate(tokenString, sessionID)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, errors.New("invalid JWT token"), requestID)
		return
	}

	userId, err := rh.sessionService.GetUserID(ctx, sessionID)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get user ID", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	var updatedReview domain.Review
	if err := easyjson.UnmarshalFromReader(r.Body, &updatedReview); err != nil {
		logger.AccessLogger.Warn("Failed to unmarshal review", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	updatedReview.Text = sanitizer.Sanitize(updatedReview.Text)
	updatedReview.Title = sanitizer.Sanitize(updatedReview.Title)
	updatedReview.HostID = sanitizer.Sanitize(updatedReview.HostID)
	updatedReview.UserID = sanitizer.Sanitize(updatedReview.UserID)

	err = rh.usecase.UpdateReview(ctx, userId, hostId, &updatedReview)
	if err != nil {
		logger.AccessLogger.Warn("Failed to update review", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	response := domain.ResponseMessage{Message: "updated successfully"}
	if _, err = easyjson.MarshalToWriter(&response, w); err != nil {
		logger.AccessLogger.Warn("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed UpdateReview request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (rh *ReviewHandler) handleError(w http.ResponseWriter, err error, requestID string) int {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)

	var statusCode int
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	errorResponse := domain.ErrorResponse{
		Error: err.Error(),
	}

	switch err.Error() {
	case "input contains invalid characters",
		"score out of range",
		"input exceeds character limit":

		statusCode = http.StatusBadRequest

	case "host and user are the same",
		"review already exist":
		statusCode = http.StatusConflict

	case "user not found",
		"review not found",
		"session not found",
		"no reviews found":
		statusCode = http.StatusNotFound

	case "token invalid",
		"token expired",
		"bad sign method",
		"missing X-CSRF-Token header",
		"invalid JWT token":
		statusCode = http.StatusUnauthorized

	case "failed to generate session id",
		"failed to save session",
		"failed to delete session",
		"error generating random bytes for session ID",
		"failed to fetch reviews for host",
		"failed to update host score",
		"error creating review",
		"error updating review",
		"error finding review",
		"error finding host",
		"error updating host score",
		"error fetching reviews",
		"error fetching user by ID":
		statusCode = http.StatusInternalServerError

	default:
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	if _, jsonErr := easyjson.MarshalToWriter(&errorResponse, w); jsonErr != nil {
		logger.AccessLogger.Error("Failed to encode error response",
			zap.String("request_id", requestID),
			zap.Error(jsonErr),
		)
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}

	return statusCode
}
