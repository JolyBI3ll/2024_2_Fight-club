package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/ads/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"2024_2_FIGHT-CLUB/internal/service/validation"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type AdHandler struct {
	adUseCase      usecase.AdUseCase
	sessionService session.InterfaceSession
	jwtToken       middleware.JwtTokenService
}

func NewAdHandler(adUseCase usecase.AdUseCase, sessionService session.InterfaceSession, jwtToken middleware.JwtTokenService) *AdHandler {
	return &AdHandler{
		adUseCase:      adUseCase,
		sessionService: sessionService,
		jwtToken:       jwtToken,
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
	offset := queryParams.Get("offset")
	var offsetInt int
	if offset != "" {
		var err error
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			logger.AccessLogger.Error("Failed to parse offset as int", zap.String("request_id", requestID), zap.Error(err))
			h.handleError(w, errors.New("query offset not int"), requestID)
			return
		}
	}
	limit := queryParams.Get("limit")
	var limitInt int
	if offset != "" {
		var err error
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			logger.AccessLogger.Error("Failed to parse limit as int", zap.String("request_id", requestID), zap.Error(err))
			h.handleError(w, errors.New("query limit not int"), requestID)
			return
		}
	}
	filter := domain.AdFilter{
		Location:    location,
		Rating:      rating,
		NewThisWeek: newThisWeek,
		HostGender:  hostGender,
		GuestCount:  guestCounter,
		Limit:       limitInt,
		Offset:      offsetInt,
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
		h.handleError(w, err, requestID)
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
	sanitizer := bluemonday.UGCPolicy()
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(adId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input contains invalid characters"), requestID)
		return
	}

	if len(adId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input exceeds character limit"), requestID)
		return
	}

	adId = sanitizer.Sanitize(adId)

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
		h.handleError(w, err, requestID)
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
	sanitizer := bluemonday.UGCPolicy()
	requestID := middleware.GetRequestID(r.Context())

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received CreatePlace request",
		zap.String("request_id", requestID),
	)

	authHeader := r.Header.Get("X-CSRF-Token")
	if authHeader == "" {
		logger.AccessLogger.Warn("Failed to X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("Missing X-CSRF-Token header")),
		)
		h.handleError(w, errors.New("Missing X-CSRF-Token header"), requestID)
		return
	}

	tokenString := authHeader[len("Bearer "):]
	_, err := h.jwtToken.Validate(tokenString)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, errors.New("Invalid JWT token"), requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	r.ParseMultipartForm(10 << 20) // 10 mb

	metadata := r.FormValue("metadata")
	var newPlace domain.CreateAdRequest
	var place domain.Ad
	if err := json.Unmarshal([]byte(metadata), &newPlace); err != nil {
		logger.AccessLogger.Error("Failed to decode metadata", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, errors.New("Failed to decode metadata"), requestID)
		return
	}

	var files []*multipart.FileHeader
	if len(r.MultipartForm.File["images"]) > 0 {
		files = r.MultipartForm.File["images"]

		if err := validation.ValidateImages(files, 5<<20, []string{"image/jpeg", "image/png", "image/jpg"}, 2000, 2000); err != nil {
			logger.AccessLogger.Warn("Invalid image", zap.String("request_id", requestID), zap.Error(err))
			h.handleError(w, errors.New("Invalid size, type or resolution of image"), requestID)
			return
		}
	} else {
		logger.AccessLogger.Warn("No images", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, errors.New("No images"), requestID)
		return
	}

	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9@.,\s\-]*$`)
	if !validCharPattern.MatchString(newPlace.CityName) ||
		!validCharPattern.MatchString(newPlace.Description) ||
		!validCharPattern.MatchString(newPlace.Address) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input contains invalid characters"), requestID)
		return
	}

	if len(newPlace.CityName) > maxLen || len(newPlace.Description) > maxLen || len(newPlace.Address) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input exceeds character limit"), requestID)
		return
	}

	const minRooms, maxRooms = 1, 100
	if newPlace.RoomsNumber < minRooms || newPlace.RoomsNumber > maxRooms {
		logger.AccessLogger.Warn("RoomsNumber out of range", zap.String("request_id", requestID))
		h.handleError(w, errors.New("RoomsNumber out of range"), requestID)
		return
	}

	newPlace.CityName = sanitizer.Sanitize(newPlace.CityName)
	newPlace.Description = sanitizer.Sanitize(newPlace.Description)
	newPlace.Address = sanitizer.Sanitize(newPlace.Address)

	userID, err := h.sessionService.GetUserID(ctx, r)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		h.handleError(w, errors.New("no active session"), requestID)
		return
	}
	place.AuthorUUID = userID

	err = h.adUseCase.CreatePlace(ctx, &place, files, newPlace)
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
		h.handleError(w, errors.New("Failed to decode response"), requestID)
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
	sanitizer := bluemonday.UGCPolicy()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]
	const maxLen = 255
	validCharPatternUrl := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPatternUrl.MatchString(adId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input contains invalid characters"), requestID)
		return
	}

	if len(adId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input exceeds character limit"), requestID)
		return
	}

	adId = sanitizer.Sanitize(adId)

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received UpdatePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	authHeader := r.Header.Get("X-CSRF-Token")
	if authHeader == "" {
		logger.AccessLogger.Warn("Failed to X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("Missing X-CSRF-Token header")),
		)
		http.Error(w, "Missing X-CSRF-Token header", http.StatusUnauthorized)
		return
	}

	tokenString := authHeader[len("Bearer "):]
	_, err := h.jwtToken.Validate(tokenString)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, errors.New("Invalid JWT token"), requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.AccessLogger.Error("Failed to parse multipart form", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, errors.New("Invalid multipart form"), requestID)
		return
	}

	metadata := r.FormValue("metadata")
	var updatedPlace domain.UpdateAdRequest
	if err := json.Unmarshal([]byte(metadata), &updatedPlace); err != nil {
		logger.AccessLogger.Error("Failed to decode metadata", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, errors.New("Invalid metadata JSON"), requestID)
		return
	}

	var files []*multipart.FileHeader
	if len(r.MultipartForm.File["images"]) > 0 {
		files = r.MultipartForm.File["images"]

		if err := validation.ValidateImages(files, 5<<20, []string{"image/jpeg", "image/png", "image/jpg"}, 2000, 2000); err != nil {
			logger.AccessLogger.Warn("Invalid image", zap.String("request_id", requestID), zap.Error(err))
			h.handleError(w, errors.New("Invalid size, type or resolution of image"), requestID)
			return
		}
	}

	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9@.,\s]*$`)
	if !validCharPattern.MatchString(updatedPlace.CityName) ||
		!validCharPattern.MatchString(updatedPlace.Description) ||
		!validCharPattern.MatchString(updatedPlace.Address) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input contains invalid characters"), requestID)
		return
	}

	if len(updatedPlace.CityName) > maxLen || len(updatedPlace.Description) > maxLen || len(updatedPlace.Address) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input exceeds character limit"), requestID)
		return
	}

	const minRooms, maxRooms = 1, 100
	if updatedPlace.RoomsNumber < minRooms || updatedPlace.RoomsNumber > maxRooms {
		logger.AccessLogger.Warn("RoomsNumber out of range", zap.String("request_id", requestID))
		h.handleError(w, errors.New("RoomsNumber out of range"), requestID)
		return
	}

	updatedPlace.CityName = sanitizer.Sanitize(updatedPlace.CityName)
	updatedPlace.Description = sanitizer.Sanitize(updatedPlace.Description)
	updatedPlace.Address = sanitizer.Sanitize(updatedPlace.Address)

	userID, err := h.sessionService.GetUserID(ctx, r)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		h.handleError(w, errors.New("no active session"), requestID)
		return
	}
	var place domain.Ad
	err = h.adUseCase.UpdatePlace(ctx, &place, adId, userID, files, updatedPlace)
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
	sanitizer := bluemonday.UGCPolicy()
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(adId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input contains invalid characters"), requestID)
		return
	}

	if len(adId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input exceeds character limit"), requestID)
		return
	}

	adId = sanitizer.Sanitize(adId)

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received DeletePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	authHeader := r.Header.Get("X-CSRF-Token")
	if authHeader == "" {
		logger.AccessLogger.Warn("Failed to X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("Missing X-CSRF-Token header")),
		)
		http.Error(w, "Missing X-CSRF-Token header", http.StatusUnauthorized)
		return
	}

	tokenString := authHeader[len("Bearer "):]
	_, err := h.jwtToken.Validate(tokenString)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	userID, err := h.sessionService.GetUserID(ctx, r)
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
	updateResponse := map[string]string{"response": "Delete successfully"}
	if err := json.NewEncoder(w).Encode(updateResponse); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	sanitizer := bluemonday.UGCPolicy()
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(city) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input contains invalid characters"), requestID)
		return
	}

	if len(city) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input exceeds character limit"), requestID)
		return
	}

	city = sanitizer.Sanitize(city)

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

func (h *AdHandler) GetUserPlaces(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	userId := mux.Vars(r)["userId"]
	sanitizer := bluemonday.UGCPolicy()
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(userId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input contains invalid characters"), requestID)
		return
	}

	if len(userId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input exceeds character limit"), requestID)
		return
	}

	userId = sanitizer.Sanitize(userId)

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received GetUserPlaces request",
		zap.String("request_id", requestID),
		zap.String("userId", userId),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	places, err := h.adUseCase.GetUserPlaces(ctx, userId)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	body := map[string]interface{}{
		"places": places,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetUserPlaces request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) DeleteAdImage(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	imageId := mux.Vars(r)["imageId"]
	adId := mux.Vars(r)["adId"]
	sanitizer := bluemonday.UGCPolicy()
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(adId) || !validCharPattern.MatchString(imageId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input contains invalid characters"), requestID)
		return
	}

	if len(adId) > maxLen || len(imageId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input exceeds character limit"), requestID)
		return
	}

	adId = sanitizer.Sanitize(adId)
	imageId = sanitizer.Sanitize(imageId)

	imageIdint, err2 := strconv.Atoi(imageId)

	if err2 != nil {
		logger.AccessLogger.Warn("Failed to ATOI image url", zap.String("request_id", requestID), zap.Error(err2))
		h.handleError(w, err2, requestID)
		return
	}
	ctx, cancel := withTimeout(r.Context())
	defer cancel()
	logger.AccessLogger.Info("Received DeleteAdImage request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
		zap.String("imageId", imageId))

	authHeader := r.Header.Get("X-CSRF-Token")
	if authHeader == "" {
		logger.AccessLogger.Warn("Failed to X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("Missing X-CSRF-Token header")),
		)
		http.Error(w, "Missing X-CSRF-Token header", http.StatusUnauthorized)
		return
	}

	tokenString := authHeader[len("Bearer "):]
	_, err := h.jwtToken.Validate(tokenString)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userID, err := h.sessionService.GetUserID(ctx, r)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		h.handleError(w, errors.New("no active session"), requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = h.adUseCase.DeleteAdImage(ctx, adId, imageIdint, userID)
	if err != nil {
		logger.AccessLogger.Error("Failed to delete ad image", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	w.WriteHeader(http.StatusOK)
	updateResponse := map[string]string{"response": "Delete image successfully"}
	if err := json.NewEncoder(w).Encode(updateResponse); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed DeleteAdImage request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
		zap.String("imageId", imageId))
	zap.String("duration", duration.String())
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
	case "not owner of ad", "no active session", "Missing X-CSRF-Token header", "Invalid JWT token":
		w.WriteHeader(http.StatusUnauthorized)
	case "Invalid metadata JSON", "Invalid multipart form", "Invalid size, type or resolution of image",
		"Input contains invalid characters", "Input exceeds character limit", "RoomsNumber out of range", "query offset not int",
		"query limit not int", "No images":
		w.WriteHeader(http.StatusBadRequest)
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
