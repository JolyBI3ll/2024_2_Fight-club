package controller

import (
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/microservices/city_service/controller/gen"
	"2024_2_FIGHT-CLUB/microservices/city_service/usecase"
	"context"
	"errors"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"regexp"
)

type GrpcCityHandler struct {
	gen.CityServiceServer
	usecase usecase.CityUseCase
}

func NewGrpcCityHandler(usecase usecase.CityUseCase) *GrpcCityHandler {
	return &GrpcCityHandler{
		usecase: usecase,
	}
}

func (h *GrpcCityHandler) GetCities(ctx context.Context, in *gen.GetCitiesRequest) (*gen.GetCitiesResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("Received GetCities request in microservice",
		zap.String("request_id", requestID),
	)
	cities, err := h.usecase.GetCities(ctx)
	if err != nil {
		logger.AccessLogger.Error("Failed to get cities",
			zap.String("request_id", requestID),
			zap.Error(err))
		return nil, err
	}

	var response []*gen.City
	for _, city := range cities {
		response = append(response, &gen.City{
			Id:          int32(city.ID),
			Title:       city.Title,
			Entitle:     city.EnTitle,
			Description: city.Description,
			Image:       city.Image,
		})
	}

	return &gen.GetCitiesResponse{
		Cities: response,
	}, nil
}

func (h *GrpcCityHandler) GetOneCity(ctx context.Context, in *gen.GetOneCityRequest) (*gen.GetOneCityResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	logger.AccessLogger.Info("Received GetOneCity request in microservice",
		zap.String("request_id", requestID),
	)

	in.EnName = sanitizer.Sanitize(in.EnName)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(in.EnName) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return nil, errors.New("input contains invalid characters")
	}

	if len(in.EnName) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return nil, errors.New("input exceeds character limit")
	}
	city, err := h.usecase.GetOneCity(ctx, in.EnName)
	if err != nil {
		logger.AccessLogger.Error("Failed to get city",
			zap.String("request_id", requestID),
			zap.Error(err))
		return nil, err
	}

	return &gen.GetOneCityResponse{
		City: &gen.City{
			Id:          int32(city.ID),
			Title:       city.Title,
			Entitle:     city.EnTitle,
			Description: city.Description,
			Image:       city.Image,
		},
	}, nil
}
