package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/ads/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"time"
)

type AdHandler struct {
	adUseCase      usecase.AdUseCase
	sessionService session.InterfaceSession
}

func NewAdHandler(adUseCase usecase.AdUseCase, sessionService session.InterfaceSession) *AdHandler {
	return &AdHandler{
		adUseCase:      adUseCase,
		sessionService: sessionService,
	}
}

const requestTimeout = 5 * time.Second

func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, requestTimeout)
}

func (h *AdHandler) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received GetAllPlaces request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	queryParams := r.URL.Query()

	location := queryParams.Get("location")
	rating := queryParams.Get("rating")
	newThisWeek := queryParams.Get("new")
	hostGender := queryParams.Get("gender")
	guestCounter := queryParams.Get("guests")

	filter := domain.AdFilter{
		Location:    location,
		Rating:      rating,
		NewThisWeek: newThisWeek,
		HostGender:  hostGender,
		GuestCount:  guestCounter,
	}

	places, err := h.adUseCase.GetAllPlaces(ctx, filter)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	body := map[string]interface{}{
		"places": places,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetAllPlaces request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) GetOnePlace(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received GetOnePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	place, err := h.adUseCase.GetOnePlace(ctx, adId)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	body := map[string]interface{}{
		"place": place,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetOnePlace request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) CreatePlace(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received CreatePlace request",
		zap.String("request_id", requestID),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	r.ParseMultipartForm(10 << 20) // 10 mb

	metadata := r.FormValue("metadata")
	var place domain.Ad
	if err := json.Unmarshal([]byte(metadata), &place); err != nil {
		logger.AccessLogger.Error("Failed to decode metadata", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, "Invalid metadata JSON", http.StatusBadRequest)
		return
	}
	var files []*multipart.FileHeader
	if len(r.MultipartForm.File["images"]) > 0 {
		files = r.MultipartForm.File["images"]
	}

	userID, err := h.sessionService.GetUserID(ctx, r, w)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		h.handleError(w, errors.New("no active session"), requestID)
		return
	}
	place.AuthorUUID = userID

	err = h.adUseCase.CreatePlace(ctx, &place, files)
	if err != nil {
		logger.AccessLogger.Error("Failed to create place", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}
	body := map[string]interface{}{
		"place": place,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed CreatePlace request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received UpdatePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.AccessLogger.Error("Failed to parse multipart form", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}

	metadata := r.FormValue("metadata")
	var place domain.Ad
	if err := json.Unmarshal([]byte(metadata), &place); err != nil {
		logger.AccessLogger.Error("Failed to decode metadata", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, "Invalid metadata JSON", http.StatusBadRequest)
		return
	}

	var files []*multipart.FileHeader
	if len(r.MultipartForm.File["images"]) > 0 {
		files = r.MultipartForm.File["images"]
	}

	userID, err := h.sessionService.GetUserID(ctx, r, w)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		h.handleError(w, errors.New("no active session"), requestID)
		return
	}

	err = h.adUseCase.UpdatePlace(ctx, &place, adId, userID, files)
	if err != nil {
		logger.AccessLogger.Error("Failed to update place", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	w.WriteHeader(http.StatusOK)
	updateResponse := map[string]string{"response": "Update successfully"}
	if err := json.NewEncoder(w).Encode(updateResponse); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed UpdatePlace request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) DeletePlace(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received DeletePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	userID, err := h.sessionService.GetUserID(ctx, r, w)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		h.handleError(w, errors.New("no active session"), requestID)
		return
	}

	err = h.adUseCase.DeletePlace(ctx, adId, userID)
	if err != nil {
		logger.AccessLogger.Error("Failed to delete place", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	w.WriteHeader(http.StatusOK)
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed DeletePlace request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) GetPlacesPerCity(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	city := mux.Vars(r)["city"]

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received GetPlacesPerCity request",
		zap.String("request_id", requestID),
		zap.String("city", city),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ads, err := h.adUseCase.GetPlacesPerCity(ctx, city)
	if err != nil {
		logger.AccessLogger.Error("Failed to get places per city", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}
	body := map[string]interface{}{
		"places": ads,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetPlacesPerCity request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) handleError(w http.ResponseWriter, err error, requestID string) {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)

	w.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{"error": err.Error()}

	switch err.Error() {
	case "ad not found":
		w.WriteHeader(http.StatusNotFound)
	case "ad already exists":
		w.WriteHeader(http.StatusConflict)
	case "not owner of ad", "no active session":
		w.WriteHeader(http.StatusUnauthorized)
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
