package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/payment/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type PaymentHandler struct {
	usecase        usecase.PaymentUseCase
	sessionService session.InterfaceSession
	jwtToken       middleware.JwtTokenService
}

func NewReviewHandler(usecase usecase.PaymentUseCase, sessionService session.InterfaceSession, jwtToken middleware.JwtTokenService) *PaymentHandler {
	return &PaymentHandler{
		usecase:        usecase,
		sessionService: sessionService,
		jwtToken:       jwtToken,
	}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
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

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	authHeader := r.Header.Get("X-CSRF-Token")
	if authHeader == "" {
		logger.AccessLogger.Warn("Missing X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("missing X-CSRF-Token header")),
		)
		statusCode = h.handleError(w, errors.New("missing X-CSRF-Token header"), requestID)
		return
	}

	tokenString := authHeader[len("Bearer "):]
	_, err = h.jwtToken.Validate(tokenString, sessionID)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, errors.New("invalid JWT token"), requestID)
		return
	}

	userId, err := h.sessionService.GetUserID(ctx, sessionID)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get user ID", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	var request domain.PaymentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.AccessLogger.Warn("Failed to unmarshal review", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	adId := request.AdId
	amount := request.Amount
	result, err := h.usecase.PaymentCreate(ctx, adId, userId, amount)
	if err != nil {
		logger.AccessLogger.Warn("Failed to create payment", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	reponse := map[string]interface{}{
		"response":     "success",
		"redirect_url": result,
	}

	if err := json.NewEncoder(w).Encode(reponse); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed CreatePayment request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *PaymentHandler) handleError(w http.ResponseWriter, err error, requestID string) int {
	return 0
}
