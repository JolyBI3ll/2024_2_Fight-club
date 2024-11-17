package grpc

import (
	"2024_2_FIGHT-CLUB/domain"
	generatedAuth "2024_2_FIGHT-CLUB/internal/auth/controller/grpc/gen"
	"2024_2_FIGHT-CLUB/internal/auth/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"context"
	"errors"
	"fmt"
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

	err := h.usecase.RegisterUser(ctx, payload)
	if err != nil {
		logger.AccessLogger.Error("Failed to register user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		return nil, err
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

func (h *GrpcAuthHandler) LoginUser(ctx context.Context, in *generatedAuth.LoginUserRequest) (*generatedAuth.UserResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	logger.AccessLogger.Info("Received LoginUser request in microservice",
		zap.String("request_id", requestID))

	in.Username = sanitizer.Sanitize(in.Username)
	in.Password = sanitizer.Sanitize(in.Password)

	payload := &domain.User{
		Username: in.Username,
		Password: in.Password,
	}

	response, err := h.usecase.LoginUser(ctx, payload)
	if err != nil {
		logger.AccessLogger.Error("Failed to login user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		return nil, err
	}

	userSession, err := h.sessionService.CreateSession(ctx, response)

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

	return &generatedAuth.UserResponse{
		SessionId: userSession,
		Jwttoken:  jwtToken,
		User: &generatedAuth.User{
			Id:       response.UUID,
			Username: response.Username,
			Email:    response.Email,
		},
	}, nil
}

func (h *GrpcAuthHandler) LogoutUser(ctx context.Context, in *generatedAuth.LogoutRequest) (*generatedAuth.LogoutUserResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("Received Logout request in microservice",
		zap.String("request_id", requestID))

	if in.AuthHeader == "" {
		return nil, errors.New("missing X-CSRF-Token header")
	}

	tokenString := in.AuthHeader[len("Bearer "):]
	fmt.Printf(tokenString)
	_, err := h.jwtToken.Validate(tokenString)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("invalid JWT token")
	}
	err = h.sessionService.LogoutSession(ctx, in.SessionId)
	if err != nil {
		logger.AccessLogger.Warn("Failed to logout session",
			zap.String("request_id", requestID),
			zap.Error(err))
		return nil, err
	}
	return &generatedAuth.LogoutUserResponse{
		Response: "Success",
	}, nil
}

func (h *GrpcAuthHandler) PutUser(ctx context.Context, in *generatedAuth.PutUserRequest) (*generatedAuth.UpdateResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	logger.AccessLogger.Info("Received LoginUser request in microservice",
		zap.String("request_id", requestID))

	if in.AuthHeader == "" {
		logger.AccessLogger.Warn("Failed to X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("missing X-CSRF-Token header")),
		)
		return nil, errors.New("missing X-CSRF-Token header")
	}

	tokenString := in.AuthHeader[len("Bearer "):]
	_, err := h.jwtToken.Validate(tokenString)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("invalid JWT token")
	}

	in.Creds.Uuid = sanitizer.Sanitize(in.Creds.Uuid)
	in.Creds.Username = sanitizer.Sanitize(in.Creds.Username)
	in.Creds.Password = sanitizer.Sanitize(in.Creds.Password)
	in.Creds.Email = sanitizer.Sanitize(in.Creds.Email)
	in.Creds.Name = sanitizer.Sanitize(in.Creds.Name)
	in.Creds.Avatar = sanitizer.Sanitize(in.Creds.Avatar)
	in.Creds.Sex = sanitizer.Sanitize(in.Creds.Sex)

	userID, err := h.sessionService.GetUserID(ctx, in.SessionId)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get user ID from session",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		return nil, errors.New("failed to get user ID")
	}
	Payload := &domain.User{
		UUID:       in.Creds.Uuid,
		Username:   in.Creds.Username,
		Password:   in.Creds.Password,
		Email:      in.Creds.Email,
		Name:       in.Creds.Name,
		Score:      float64(in.Creds.Score),
		Avatar:     in.Creds.Avatar,
		GuestCount: int(in.Creds.GuestCount),
		Birthdate:  (in.Creds.Birthdate).AsTime(),
		IsHost:     in.Creds.IsHost,
	}
	err = h.usecase.PutUser(ctx, Payload, userID, in.Avatar)
	if err != nil {
		return nil, err
	}
	return &generatedAuth.UpdateResponse{
		Response: "Success",
	}, nil
}
