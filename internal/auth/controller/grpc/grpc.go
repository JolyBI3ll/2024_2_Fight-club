package grpc

import (
	"2024_2_FIGHT-CLUB/domain"
	generatedAuth "2024_2_FIGHT-CLUB/internal/auth/controller/grpc/gen"
	"2024_2_FIGHT-CLUB/internal/auth/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"context"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"time"
)

type GrpcAuthHandler struct {
	generatedAuth.AuthServer
	usecase        usecase.AuthUseCase
	sessionService session.InterfaceSession
	jwtToken       middleware.JwtTokenService
}

func NewGrpcAuthHandler(usecase usecase.AuthUseCase, sessionService session.InterfaceSession, jwtToken middleware.JwtTokenService) *GrpcAuthHandler {
	return &GrpcAuthHandler{
		usecase:        usecase,
		sessionService: sessionService,
		jwtToken:       jwtToken,
	}
}

func (h *GrpcAuthHandler) RegisterUser(ctx context.Context, in *generatedAuth.RegisterUserRequest) (*generatedAuth.UserResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	logger.AccessLogger.Info("Received RegisterUser request in microservice",
		zap.String("request_id", requestID),
	)

	in.Username = sanitizer.Sanitize(in.Username)
	in.Email = sanitizer.Sanitize(in.Email)
	in.Password = sanitizer.Sanitize(in.Password)
	in.Name = sanitizer.Sanitize(in.Name)

	payload := &domain.User{
		Username: in.Username,
		Email:    in.Email,
		Password: in.Password,
		Name:     in.Name,
	}

	userSession, err := h.sessionService.CreateSession(ctx, payload)
	if err != nil {
		logger.AccessLogger.Error("Failed create session",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		return nil, err
	}

	tokenExpTime := time.Now().Add(24 * time.Hour).Unix()
	jwtToken, err := h.jwtToken.Create(userSession, tokenExpTime)
	if err != nil {
		logger.AccessLogger.Error("Failed to create JWT token",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		return nil, err
	}

	err = h.usecase.RegisterUser(ctx, payload)
	if err != nil {
		logger.AccessLogger.Error("Failed to register user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		return nil, err
	}

	return &generatedAuth.UserResponse{
		SessionId: userSession,
		Jwttoken:  jwtToken,
		User: &generatedAuth.User{
			Id:       payload.UUID,
			Username: payload.Username,
			Email:    payload.Email,
		},
	}, nil
}
