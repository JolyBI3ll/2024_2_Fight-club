package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/utils"
	"2024_2_FIGHT-CLUB/microservices/city_service/controller/gen"
	"errors"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type CityHandler struct {
	client gen.CityServiceClient
	utils  utils.UtilsInterface
}

func NewCityHandler(client gen.CityServiceClient, utils utils.UtilsInterface) *CityHandler {
	return &CityHandler{
		client: client,
		utils:  utils,
	}
}

func (h *CityHandler) GetCities(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received GetCities request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)

	statusCode := http.StatusOK
	var err error
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

	cities, err := h.client.GetCities(ctx, &gen.GetCitiesRequest{})
	if err != nil {
		logger.AccessLogger.Error("Failed to get cities data",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}
	payload, err := h.utils.ConvertAllCitiesProtoToGo(cities)
	if err != nil {
		logger.AccessLogger.Error("Failed to convert cities data",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	body := domain.AllCitiesResponse{
		Cities: payload,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = easyjson.MarshalToWriter(body, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetCities request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *CityHandler) GetOneCity(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()
	logger.AccessLogger.Info("Received GetOneCity request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)

	statusCode := http.StatusOK
	var err error
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

	cityEnName := mux.Vars(r)["city"]
	city, err := h.client.GetOneCity(ctx, &gen.GetOneCityRequest{EnName: cityEnName})
	if err != nil {
		logger.AccessLogger.Error("Failed to get city data",
			zap.String("request_id", requestID),
			zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}
	payload, err := h.utils.ConvertOneCityProtoToGo(city.City)
	if err != nil {
		logger.AccessLogger.Error("Failed to convert city data",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	body := domain.OneCityResponse{
		City: payload,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = easyjson.MarshalToWriter(body, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetOneCity request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK))
}

func (h *CityHandler) handleError(w http.ResponseWriter, err error, requestID string) int {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)

	w.Header().Set("Content-Type", "application/json")
	errorResponse := domain.ErrorResponse{
		Error: err.Error(),
	}
	var statusCode int
	switch err.Error() {
	case "input contains invalid characters",
		"input exceeds character limit":
		statusCode = http.StatusBadRequest
	case "error fetching all cities",
		"error fetching city":
		statusCode = http.StatusInternalServerError
	default:
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	if _, jsonErr := easyjson.MarshalToWriter(errorResponse, w); jsonErr != nil {
		logger.AccessLogger.Error("Failed to encode error response",
			zap.String("request_id", requestID),
			zap.Error(jsonErr),
		)
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}
	return statusCode
}
