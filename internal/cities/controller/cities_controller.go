package controller

import (
	"2024_2_FIGHT-CLUB/internal/cities/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type CityHandler struct {
	cityUseCase usecase.CityUseCase
}

func NewCityHandler(cityUseCase usecase.CityUseCase) *CityHandler {
	return &CityHandler{
		cityUseCase: cityUseCase,
	}
}

const requestTimeout = 5 * time.Second

func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, requestTimeout)
}

func (h *CityHandler) GetCities(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received GetCities request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cities, err := h.cityUseCase.GetCities(ctx)
	if err != nil {
		logger.AccessLogger.Error("Failed to get cities data",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}
	body := map[string]interface{}{
		"cities": cities,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetCities request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *CityHandler) handleError(w http.ResponseWriter, err error, requestID string) {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)

	w.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{"error": err.Error()}

	switch err.Error() {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	if jsonErr := json.NewEncoder(w).Encode(errorResponse); jsonErr != nil {
		logger.AccessLogger.Error("Failed to encode error response",
			zap.String("request_id", requestID),
			zap.Error(jsonErr),
		)
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}
}
