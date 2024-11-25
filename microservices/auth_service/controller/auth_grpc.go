package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"2024_2_FIGHT-CLUB/microservices/auth_service/controller/gen"
	"2024_2_FIGHT-CLUB/microservices/auth_service/usecase"
	"context"
	"errors"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type GrpcAuthHandler struct {
	gen.AuthServer
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

func (h *GrpcAuthHandler) RegisterUser(ctx context.Context, in *gen.RegisterUserRequest) (*gen.UserResponse, error) {
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

	return &gen.UserResponse{
		SessionId: userSession,
		Jwttoken:  jwtToken,
		User: &gen.User{
			Id:       payload.UUID,
			Username: payload.Username,
			Email:    payload.Email,
		},
	}, nil
}

func (h *GrpcAuthHandler) LoginUser(ctx context.Context, in *gen.LoginUserRequest) (*gen.UserResponse, error) {
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

	return &gen.UserResponse{
		SessionId: userSession,
		Jwttoken:  jwtToken,
		User: &gen.User{
			Id:       response.UUID,
			Username: response.Username,
			Email:    response.Email,
		},
	}, nil
}

func (h *GrpcAuthHandler) LogoutUser(ctx context.Context, in *gen.LogoutRequest) (*gen.LogoutUserResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("Received Logout request in microservice",
		zap.String("request_id", requestID))

	if in.AuthHeader == "" {
		return nil, errors.New("missing X-CSRF-Token header")
	}

	tokenString := in.AuthHeader[len("Bearer "):]
	_, err := h.jwtToken.Validate(tokenString, in.SessionId)
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
	return &gen.LogoutUserResponse{
		Response: "Success",
	}, nil
}

func (h *GrpcAuthHandler) PutUser(ctx context.Context, in *gen.PutUserRequest) (*gen.UpdateResponse, error) {
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
	_, err := h.jwtToken.Validate(tokenString, in.SessionId)
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
		Sex:        in.Creds.Sex,
		GuestCount: int(in.Creds.GuestCount),
		Birthdate:  (in.Creds.Birthdate).AsTime(),
		IsHost:     in.Creds.IsHost,
	}
	err = h.usecase.PutUser(ctx, Payload, userID, in.Avatar)
	if err != nil {
		logger.AccessLogger.Warn("Failed to put user",
			zap.String("request_id", requestID),
			zap.Error(err))
		return nil, err
	}
	return &gen.UpdateResponse{
		Response: "Success",
	}, nil
}

func (h *GrpcAuthHandler) GetUserById(ctx context.Context, in *gen.GetUserByIdRequest) (*gen.GetUserByIdResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	in.UserId = sanitizer.Sanitize(in.UserId)
	logger.AccessLogger.Info("Received GetUserByID request in microservice",
		zap.String("request_id", requestID))

	user, err := h.usecase.GetUserById(ctx, in.UserId)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get user",
			zap.String("request_id", requestID),
			zap.Error(err))
		return nil, err
	}

	userMetadata := &gen.MetadataOneUser{
		Uuid:       user.UUID,
		Username:   user.Username,
		Name:       user.Name,
		Score:      float32(user.Score),
		Avatar:     user.Avatar,
		Sex:        user.Sex,
		GuestCount: int32(user.GuestCount),
		Birthdate:  timestamppb.New(user.Birthdate),
		IsHost:     user.IsHost,
	}

	return &gen.GetUserByIdResponse{
		User: userMetadata,
	}, nil
}

func (h *GrpcAuthHandler) GetAllUsers(ctx context.Context, in *gen.Empty) (*gen.AllUsersResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("Received GetAllUsers request in microservice",
		zap.String("request_id", requestID))

	users, err := h.usecase.GetAllUser(ctx)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get all users",
			zap.String("request_id", requestID),
			zap.Error(err))
		return nil, err
	}

	var userMetadata []*gen.MetadataOneUser
	for _, user := range users {
		userMetadata = append(userMetadata, &gen.MetadataOneUser{
			Uuid:       user.UUID,
			Username:   user.Username,
			Name:       user.Name,
			Score:      float32(user.Score),
			Avatar:     user.Avatar,
			Sex:        user.Sex,
			GuestCount: int32(user.GuestCount),
			Birthdate:  timestamppb.New(user.Birthdate),
			IsHost:     user.IsHost,
		})
	}

	return &gen.AllUsersResponse{
		Users: userMetadata,
	}, nil
}

func (h *GrpcAuthHandler) GetSessionData(ctx context.Context, in *gen.GetSessionDataRequest) (*gen.SessionDataResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("Received GetSessionData request in microservice",
		zap.String("request_id", requestID))
	sessionData, err := h.sessionService.GetSessionData(ctx, in.SessionId)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get session data",
			zap.String("request_id", requestID),
			zap.Error(err))
		return nil, errors.New("failed to get session data")
	}
	data := *sessionData

	id, ok := data["id"].(string)
	if !ok {
		logger.AccessLogger.Warn("Invalid type for id in session data",
			zap.String("request_id", requestID))
		return nil, errors.New("invalid type for id in session data")
	}

	avatar, ok := data["avatar"].(string)
	if !ok {
		logger.AccessLogger.Warn("Invalid type for avatar in session data",
			zap.String("request_id", requestID))
		return nil, errors.New("invalid type for avatar in session data")
	}

	return &gen.SessionDataResponse{
		Id:     id,
		Avatar: avatar,
	}, nil
}

func (h *GrpcAuthHandler) RefreshCsrfToken(ctx context.Context, in *gen.RefreshCsrfTokenRequest) (*gen.RefreshCsrfTokenResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("Received RefreshCsrfToken request in microservice",
		zap.String("request_id", requestID))
	newCsrfToken, err := h.jwtToken.Create(in.SessionId, time.Now().Add(1*time.Hour).Unix())
	if err != nil {
		logger.AccessLogger.Warn("Failed to refresh csrf token",
			zap.String("request_id", requestID),
			zap.Error(err))
		return nil, errors.New("failed to refresh csrf token")
	}
	return &gen.RefreshCsrfTokenResponse{
		CsrfToken: newCsrfToken,
	}, nil
}
