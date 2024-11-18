package http

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/ads/controller/grpc/gen"
	"2024_2_FIGHT-CLUB/internal/ads/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
	"time"
)

type AdHandler struct {
	client         gen.AdsClient
	adUseCase      usecase.AdUseCase
	sessionService session.InterfaceSession
	jwtToken       middleware.JwtTokenService
}

func NewAdHandler(client gen.AdsClient, adUseCase usecase.AdUseCase, sessionService session.InterfaceSession, jwtToken middleware.JwtTokenService) *AdHandler {
	return &AdHandler{
		client:         client,
		adUseCase:      adUseCase,
		sessionService: sessionService,
		jwtToken:       jwtToken,
	}
}

func (h *AdHandler) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received GetAllPlaces request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	queryParams := r.URL.Query()

	response, err := h.client.GetAllPlaces(ctx, &gen.AdFilterRequest{
		Location:    queryParams.Get("location"),
		Rating:      queryParams.Get("rating"),
		NewThisWeek: queryParams.Get("new"),
		HostGender:  queryParams.Get("gender"),
		GuestCount:  queryParams.Get("guests"),
		Limit:       queryParams.Get("limit"),
		Offset:      queryParams.Get("offset"),
		DateFrom:    queryParams.Get("dateFrom"),
		DateTo:      queryParams.Get("dateTo"),
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to GetAllPlaces",
			zap.Error(err),
			zap.String("request_id", requestID),
			zap.String("method", r.Method))
		h.handleError(w, err, requestID)
	}

	body := map[string]interface{}{
		"places": response,
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
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received GetOnePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	isAuthorized := true

	if _, err := h.sessionService.GetUserID(ctx, sessionID); err != nil {
		isAuthorized = false
	}

	place, err := h.client.GetOnePlace(ctx, &gen.GetPlaceByIdRequest{
		AdId:         adId,
		IsAuthorized: isAuthorized,
	})
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
	requestID := middleware.GetRequestID(r.Context())

	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received CreatePlace request",
		zap.String("request_id", requestID),
	)

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	authHeader := r.Header.Get("X-CSRF-Token")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		logger.AccessLogger.Error("Failed to parse multipart form", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	metadata := r.FormValue("metadata")
	var newPlace domain.CreateAdRequest
	if err := json.Unmarshal([]byte(metadata), &newPlace); err != nil {
		logger.AccessLogger.Error("Failed to decode metadata", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, errors.New("Failed to decode metadata"), requestID)
		return
	}

	fileHeaders := r.MultipartForm.File["images"]
	if len(fileHeaders) == 0 {
		logger.AccessLogger.Warn("No images", zap.String("request_id", requestID))
		h.handleError(w, errors.New("No images provided"), requestID)
		return
	}

	// Преобразование файлов в [][]byte
	files := make([][]byte, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			logger.AccessLogger.Error("Failed to open file", zap.String("request_id", requestID), zap.Error(err))
			h.handleError(w, errors.New("Failed to open file"), requestID)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			logger.AccessLogger.Error("Failed to read file", zap.String("request_id", requestID), zap.Error(err))
			h.handleError(w, errors.New("Failed to read file"), requestID)
			return
		}

		files = append(files, data)
	}

	response, err := h.client.CreatePlace(ctx, &gen.CreateAdRequest{
		CityName:    newPlace.CityName,
		Description: newPlace.Description,
		Address:     newPlace.Address,
		RoomsNumber: int32(newPlace.RoomsNumber),
		DateFrom:    timestamppb.New(newPlace.DateFrom),
		DateTo:      timestamppb.New(newPlace.DateTo),
		Images:      files,
		AuthHeader:  authHeader,
		SessionID:   sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to create place", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	body := map[string]interface{}{
		"place": response,
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, errors.New("Failed to encode response"), requestID)
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

	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received UpdatePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	authHeader := r.Header.Get("X-CSRF-Token")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := r.ParseMultipartForm(10 << 20) // 10 MB
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

	fileHeaders := r.MultipartForm.File["images"]

	// Преобразование `[]*multipart.FileHeader` в `[][]byte`
	files := make([][]byte, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			logger.AccessLogger.Error("Failed to open file", zap.String("request_id", requestID), zap.Error(err))
			h.handleError(w, errors.New("Failed to open file"), requestID)
			return
		}
		defer file.Close()

		// Чтение содержимого файла в []byte
		data, err := io.ReadAll(file)
		if err != nil {
			logger.AccessLogger.Error("Failed to read file", zap.String("request_id", requestID), zap.Error(err))
			h.handleError(w, errors.New("Failed to read file"), requestID)
			return
		}
		files = append(files, data)
	}

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	response, err := h.client.UpdatePlace(ctx, &gen.UpdateAdRequest{
		AdId:        adId,
		CityName:    updatedPlace.CityName,
		Address:     updatedPlace.Address,
		Description: updatedPlace.Description,
		RoomsNumber: int32(updatedPlace.RoomsNumber),
		SessionID:   sessionID,
		AuthHeader:  authHeader,
		Images:      files,
		DateFrom:    timestamppb.New(updatedPlace.DateFrom),
		DateTo:      timestamppb.New(updatedPlace.DateTo),
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to update place", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	w.WriteHeader(http.StatusOK)
	updateResponse := map[string]string{"response": response.Response}
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

	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received DeletePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	authHeader := r.Header.Get("X-CSRF-Token")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	response, err := h.client.DeletePlace(ctx, &gen.DeletePlaceRequest{
		AdId:       adId,
		SessionID:  sessionID,
		AuthHeader: authHeader,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to delete place", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	w.WriteHeader(http.StatusOK)
	updateResponse := map[string]string{"response": response.Response}
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

	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received GetPlacesPerCity request",
		zap.String("request_id", requestID),
		zap.String("city", city),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response, err := h.client.GetPlacesPerCity(ctx, &gen.GetPlacesPerCityRequest{
		CityName: city,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to get places per city", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}
	body := map[string]interface{}{
		"places": response,
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

	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received GetUserPlaces request",
		zap.String("request_id", requestID),
		zap.String("userId", userId),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response, err := h.client.GetUserPlaces(ctx, &gen.GetUserPlacesRequest{
		UserId: userId,
	})
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	body := map[string]interface{}{
		"places": response,
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

	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received DeleteAdImage request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
		zap.String("imageId", imageId))

	authHeader := r.Header.Get("X-CSRF-Token")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	response, err := h.client.DeleteAdImage(ctx, &gen.DeleteAdImageRequest{
		AdId:       adId,
		ImageId:    imageId,
		AuthHeader: authHeader,
		SessionID:  sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to delete ad image", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	w.WriteHeader(http.StatusOK)
	updateResponse := map[string]string{"response": response.Response}
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
	case "not owner of ad", "no active session", "Missing X-CSRF-Token header", "Invalid JWT token", "User is not host":
		w.WriteHeader(http.StatusUnauthorized)
	case "Invalid metadata JSON", "Invalid multipart form", "Invalid size, type or resolution of image",
		"Input contains invalid characters", "Input exceeds character limit", "RoomsNumber out of range", "query offset not int",
		"query limit not int", "No images", "URL contains invalid characters", "URL exceeds character limit":
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
