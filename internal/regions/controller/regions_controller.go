package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/regions/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type RegionHandler struct {
	usecase        usecase.RegionUsecase
	sessionService session.InterfaceSession
	jwtToken       middleware.JwtTokenService
}

func NewRegionHandler(usecase usecase.RegionUsecase, sessionService session.InterfaceSession, jwtToken middleware.JwtTokenService) *RegionHandler {
	return &RegionHandler{
		usecase:        usecase,
		sessionService: sessionService,
		jwtToken:       jwtToken,
	}
}

func (rh *RegionHandler) GetVisitedRegions(w http.ResponseWriter, r *http.Request) {
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

	logger.AccessLogger.Info("Received GetVisitedRegions request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var regions domain.VisitedRegionsList

	regions, err = rh.usecase.GetVisitedRegions(ctx, userId)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get visited regions", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = easyjson.MarshalToWriter(&regions, w); err != nil {
		logger.AccessLogger.Warn("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		statusCode = rh.handleError(w, err, requestID)
		return
	}
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetVisitedRegions request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK))
}

func (rh *RegionHandler) handleError(w http.ResponseWriter, err error, requestID string) int {
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
