package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"2024_2_FIGHT-CLUB/internal/service/utils"
	"2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	"errors"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
	"time"
)

type AdHandler struct {
	client         gen.AdsClient
	sessionService session.InterfaceSession
	jwtToken       middleware.JwtTokenService
	utils          utils.UtilsInterface
}

func NewAdHandler(client gen.AdsClient, sessionService session.InterfaceSession, jwtToken middleware.JwtTokenService, utils utils.UtilsInterface) *AdHandler {
	return &AdHandler{
		client:         client,
		sessionService: sessionService,
		jwtToken:       jwtToken,
		utils:          utils,
	}
}

func (h *AdHandler) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error

	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeAdIdPath(r.URL.Path)
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received GetAllPlaces request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)

	sessionID, err := session.GetSessionId(r)
	if err != nil || sessionID == "" {
		logger.AccessLogger.Warn("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
	}

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
		SessionId:   sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to GetAllPlaces",
			zap.Error(err),
			zap.String("request_id", requestID),
			zap.String("method", r.Method))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}
	body, err := h.utils.ConvertGetAllAdsResponseProtoToGo(response)
	if err != nil {
		logger.AccessLogger.Error("Failed to Convert From Proto to Go",
			zap.Error(err),
			zap.String("request_id", requestID))
		statusCode = h.handleError(w, err, requestID)
		return
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

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeAdIdPath(r.URL.Path)
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received GetOnePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	isAuthorized := false

	sessionID, err := session.GetSessionId(r)
	if err != nil || sessionID == "" {
		logger.AccessLogger.Warn("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
	} else if _, err := h.sessionService.GetUserID(ctx, sessionID); err != nil {
		logger.AccessLogger.Warn("Failed to validate session",
			zap.String("request_id", requestID),
			zap.Error(err))
	} else {
		isAuthorized = true
	}

	place, err := h.client.GetOnePlace(ctx, &gen.GetPlaceByIdRequest{
		AdId:         adId,
		IsAuthorized: isAuthorized,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to GetOnePlace",
			zap.String("request_id", requestID),
			zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	payload, err := h.utils.ConvertAdProtoToGo(place)
	if err != nil {
		logger.AccessLogger.Error("Failed to Convert From Proto to Go",
			zap.Error(err),
			zap.String("request_id", requestID))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	response := domain.GetOneAdResponse{
		Place: payload,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = easyjson.MarshalToWriter(response, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
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

	statusCode := http.StatusCreated
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeAdIdPath(r.URL.Path)
		if statusCode == http.StatusCreated {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received CreatePlace request",
		zap.String("request_id", requestID),
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

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		logger.AccessLogger.Error("Failed to parse multipart form", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	metadata := r.FormValue("metadata")
	var newPlace domain.CreateAdRequest
	if err = newPlace.UnmarshalJSON([]byte(metadata)); err != nil {
		logger.AccessLogger.Error("Failed to decode metadata", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	fileHeaders := r.MultipartForm.File["images"]
	if len(fileHeaders) == 0 {
		logger.AccessLogger.Warn("No images", zap.String("request_id", requestID))
		err = errors.New("no images provided")
		statusCode = h.handleError(w, err, requestID)
		return
	}

	// Преобразование файлов в [][]byte
	files := make([][]byte, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			logger.AccessLogger.Error("Failed to open file", zap.String("request_id", requestID), zap.Error(err))
			statusCode = h.handleError(w, errors.New("failed to open file"), requestID)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			logger.AccessLogger.Error("Failed to read file", zap.String("request_id", requestID), zap.Error(err))
			statusCode = h.handleError(w, errors.New("failed to read file"), requestID)
			return
		}

		files = append(files, data)
	}

	_, err = h.client.CreatePlace(ctx, &gen.CreateAdRequest{
		CityName:     newPlace.CityName,
		Description:  newPlace.Description,
		Address:      newPlace.Address,
		RoomsNumber:  int32(newPlace.RoomsNumber),
		DateFrom:     timestamppb.New(newPlace.DateFrom),
		DateTo:       timestamppb.New(newPlace.DateTo),
		Images:       files,
		AuthHeader:   authHeader,
		SessionID:    sessionID,
		SquareMeters: int32(newPlace.SquareMeters),
		Floor:        int32(newPlace.Floor),
		BuildingType: newPlace.BuildingType,
		HasBalcony:   newPlace.HasBalcony,
		HasElevator:  newPlace.HasElevator,
		HasGas:       newPlace.HasGas,
		Rooms:        middleware.ConvertRoomsToGRPC(newPlace.Rooms),
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to create place", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	body := domain.ResponseMessage{
		Message: "Successfully created ad",
	}
	if _, err := easyjson.MarshalToWriter(body, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, errors.New("failed to encode response"), requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed CreatePlace request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusCreated),
	)
}

func (h *AdHandler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeAdIdPath(r.URL.Path)
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received UpdatePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	authHeader := r.Header.Get("X-CSRF-Token")

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		logger.AccessLogger.Error("Failed to parse multipart form", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, errors.New("invalid multipart form"), requestID)
		return
	}

	metadata := r.FormValue("metadata")
	var updatedPlace domain.UpdateAdRequest
	if err = updatedPlace.UnmarshalJSON([]byte(metadata)); err != nil {
		logger.AccessLogger.Error("Failed to decode metadata", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, errors.New("invalid metadata JSON"), requestID)
		return
	}

	fileHeaders := r.MultipartForm.File["images"]

	// Преобразование `[]*multipart.FileHeader` в `[][]byte`
	files := make([][]byte, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			logger.AccessLogger.Error("Failed to open file", zap.String("request_id", requestID), zap.Error(err))
			statusCode = h.handleError(w, errors.New("failed to open file"), requestID)
			return
		}
		defer file.Close()

		// Чтение содержимого файла в []byte
		data, err := io.ReadAll(file)
		if err != nil {
			logger.AccessLogger.Error("Failed to read file", zap.String("request_id", requestID), zap.Error(err))
			statusCode = h.handleError(w, errors.New("failed to read file"), requestID)
			return
		}
		files = append(files, data)
	}

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	_, err = h.client.UpdatePlace(ctx, &gen.UpdateAdRequest{
		AdId:         adId,
		CityName:     updatedPlace.CityName,
		Address:      updatedPlace.Address,
		Description:  updatedPlace.Description,
		RoomsNumber:  int32(updatedPlace.RoomsNumber),
		SessionID:    sessionID,
		AuthHeader:   authHeader,
		Images:       files,
		DateFrom:     timestamppb.New(updatedPlace.DateFrom),
		DateTo:       timestamppb.New(updatedPlace.DateTo),
		SquareMeters: int32(updatedPlace.SquareMeters),
		Floor:        int32(updatedPlace.Floor),
		BuildingType: updatedPlace.BuildingType,
		HasBalcony:   updatedPlace.HasBalcony,
		HasElevator:  updatedPlace.HasElevator,
		HasGas:       updatedPlace.HasGas,
		Rooms:        middleware.ConvertRoomsToGRPC(updatedPlace.Rooms),
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to update place", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	updateResponse := domain.ResponseMessage{
		Message: "Successfully updated ad",
	}
	if _, err = easyjson.MarshalToWriter(updateResponse, w); err != nil {
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

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeAdIdPath(r.URL.Path)
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received DeletePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	authHeader := r.Header.Get("X-CSRF-Token")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	_, err = h.client.DeletePlace(ctx, &gen.DeletePlaceRequest{
		AdId:       adId,
		SessionID:  sessionID,
		AuthHeader: authHeader,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to delete place", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	deleteResponse := domain.ResponseMessage{
		Message: "Successfully deleted place",
	}
	if _, err = easyjson.MarshalToWriter(deleteResponse, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AdHandler) GetPlacesPerCity(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	city := mux.Vars(r)["city"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeAdIdPath(r.URL.Path)
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received GetPlacesPerCity request",
		zap.String("request_id", requestID),
		zap.String("city", city),
	)

	response, err := h.client.GetPlacesPerCity(ctx, &gen.GetPlacesPerCityRequest{
		CityName: city,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to get places per city", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	payload, err := h.utils.ConvertGetAllAdsResponseProtoToGo(response)
	if err != nil {
		logger.AccessLogger.Error("Failed to Convert From Proto to Go",
			zap.Error(err),
			zap.String("request_id", requestID))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	body := domain.PlacesResponse{
		Places: payload,
	}
	if _, err = easyjson.MarshalToWriter(body, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
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

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeUserIdPath(r.URL.Path)
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received GetUserPlaces request",
		zap.String("request_id", requestID),
		zap.String("userId", userId),
	)

	response, err := h.client.GetUserPlaces(ctx, &gen.GetUserPlacesRequest{
		UserId: userId,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to get places per user",
			zap.String("request_id", requestID),
			zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	payload, err := h.utils.ConvertGetAllAdsResponseProtoToGo(response)
	if err != nil {
		logger.AccessLogger.Error("Failed to Convert From Proto to Go",
			zap.Error(err),
			zap.String("request_id", requestID))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	body := domain.PlacesResponse{
		Places: payload,
	}
	if _, err = easyjson.MarshalToWriter(body, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
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

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeAdIdPath(r.URL.Path)
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

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
		statusCode = h.handleError(w, err, requestID)
		return
	}

	_, err = h.client.DeleteAdImage(ctx, &gen.DeleteAdImageRequest{
		AdId:       adId,
		ImageId:    imageId,
		AuthHeader: authHeader,
		SessionID:  sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to delete ad image", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	deleteResponse := domain.ResponseMessage{
		Message: "Successfully deleted ad image",
	}
	if _, err = easyjson.MarshalToWriter(deleteResponse, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed DeleteAdImage request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
		zap.String("imageId", imageId),
		zap.Duration("duration", duration),
	)
}

func (h *AdHandler) AddToFavorites(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeAdIdPath(r.URL.Path)
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received AddToFavorites request",
		zap.String("request_id", requestID),
		zap.String("adId", adId))

	authHeader := r.Header.Get("X-CSRF-Token")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	_, err = h.client.AddToFavorites(ctx, &gen.AddToFavoritesRequest{
		AdId:       adId,
		AuthHeader: authHeader,
		SessionID:  sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to add ad to favorites", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	addToResponse := domain.ResponseMessage{
		Message: "Successfully added to favorites",
	}
	if _, err := easyjson.MarshalToWriter(addToResponse, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed AddToFavorites request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
		zap.Duration("duration", duration),
	)

}

func (h *AdHandler) DeleteFromFavorites(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeAdIdPath(r.URL.Path)
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received DeleteFromFavorites request",
		zap.String("request_id", requestID),
		zap.String("adId", adId))

	authHeader := r.Header.Get("X-CSRF-Token")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	_, err = h.client.DeleteFromFavorites(ctx, &gen.DeleteFromFavoritesRequest{
		AdId:       adId,
		AuthHeader: authHeader,
		SessionID:  sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to delete ad from favorites", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	deleteFromFavResponse := domain.ResponseMessage{
		Message: "Successfully deleted from favorites",
	}
	if _, err = easyjson.MarshalToWriter(deleteFromFavResponse, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed DeleteFromFavorites request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
		zap.Duration("duration", duration),
	)

}

func (h *AdHandler) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	userId := mux.Vars(r)["userId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeUserIdPath(r.URL.Path)
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received GetUserFavorites request",
		zap.String("request_id", requestID),
		zap.String("userId", userId))

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	response, err := h.client.GetUserFavorites(ctx, &gen.GetUserFavoritesRequest{
		UserId:    userId,
		SessionID: sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to delete ad from favorites", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	body, err := h.utils.ConvertGetAllAdsResponseProtoToGo(response)
	if err != nil {
		logger.AccessLogger.Error("Failed to Convert From Proto to Go",
			zap.Error(err),
			zap.String("request_id", requestID))
		statusCode = h.handleError(w, err, requestID)
		return
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
	logger.AccessLogger.Info("Completed GetUserFavorites request",
		zap.String("request_id", requestID),
		zap.String("userId", userId),
		zap.Duration("duration", duration),
	)

}

func (h *AdHandler) UpdatePriorityWithPayment(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeAdIdPath(r.URL.Path)
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received UpdatePriorityWithPayment request",
		zap.String("request_id", requestID),
		zap.String("userId", adId))

	var card domain.PaymentInfo
	if err = easyjson.UnmarshalFromReader(r.Body, &card); err != nil {
		logger.AccessLogger.Error("Failed to decode request body",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
		return
	}

	authHeader := r.Header.Get("X-CSRF-Token")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	_, err = h.client.UpdatePriority(ctx, &gen.UpdatePriorityRequest{
		AdId:       adId,
		AuthHeader: authHeader,
		SessionID:  sessionID,
		Amount:     card.DonationAmount,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to update priority", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	updatePriorResponse := domain.ResponseMessage{
		Message: "Successfully update ad priority",
	}
	if _, err = easyjson.MarshalToWriter(updatePriorResponse, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed UpdatePriorityWithPayment request",
		zap.String("request_id", requestID),
		zap.String("userId", adId),
		zap.Duration("duration", duration),
	)
}

func (h *AdHandler) handleError(w http.ResponseWriter, err error, requestID string) int {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)
	var statusCode int
	w.Header().Set("Content-Type", "application/json")
	errorResponse := domain.ErrorResponse{
		Error: err.Error(),
	}
	switch err.Error() {
	case "ad not found", "ad date not found", "image not found", "error fetching all places":
		statusCode = http.StatusNotFound
	case "ad already exists", "roomsNumber out of range", "not owner of ad":
		statusCode = http.StatusConflict
	case "no active session", "missing X-CSRF-Token header",
		"invalid JWT token", "user is not host", "session not found", "user ID not found in session":
		statusCode = http.StatusUnauthorized
	case "invalid metadata JSON", "invalid multipart form", "input contains invalid characters",
		"input exceeds character limit", "invalid size, type or resolution of image",
		"query offset not int", "query limit not int", "query dateFrom not int",
		"query dateTo not int", "URL contains invalid characters", "URL exceeds character limit",
		"token parse error", "token invalid", "token expired", "bad sign method",
		"failed to decode metadata", "no images provided", "failed to open file",
		"failed to read file", "failed to encode response", "invalid rating value",
		"cant access other user favorites":
		statusCode = http.StatusBadRequest
	case "error fetching images for ad", "error fetching user",
		"error finding user", "error finding city", "error creating place", "error creating date",
		"error saving place", "error updating place", "error updating date",
		"error updating views count", "error deleting place", "get places error",
		"get places per city error", "get user places error", "error creating image",
		"delete ad image error", "failed to generate session id", "failed to save session",
		"failed to delete session", "error generating random bytes for session ID",
		"failed to get session id from request cookie", "error fetching rooms for ad",
		"error counting favorites", "error updating favorites count", "error creating room", "error parsing date",
		"adAuthor is nil", "ad is nil":
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
