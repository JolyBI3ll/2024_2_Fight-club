package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"2024_2_FIGHT-CLUB/microservices/csat_service/controller/gen"
	"2024_2_FIGHT-CLUB/microservices/csat_service/usecase"
	"context"
	"errors"
	"go.uber.org/zap"
)

type GrpcCsatHandler struct {
	gen.CsatServer
	usecase        usecase.CsatUseCase
	sessionService session.InterfaceSession
	jwtToken       middleware.JwtTokenService
}

func NewGrpcCsatHandler(usecase usecase.CsatUseCase, sessionService session.InterfaceSession, jwtToken middleware.JwtTokenService) *GrpcCsatHandler {
	return &GrpcCsatHandler{
		usecase:        usecase,
		sessionService: sessionService,
		jwtToken:       jwtToken,
	}
}

func (adh *GrpcCsatHandler) GetSurvey(ctx context.Context, in *gen.GetSurveyRequest) (*gen.GetSurveyResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("Received GetSurvey request in microservice",
		zap.String("request_id", requestID),
	)

	survey, err := adh.usecase.GetSurvey(ctx, int(in.SurveyId))
	if err != nil {
		logger.AccessLogger.Warn("GetSurvey failed", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("get survey failed")
	}

	grpcSurvey := &gen.Survey{
		Id:    int32(survey.ID),
		Title: survey.Title,
	}
	for _, question := range survey.Questions {
		grpcSurvey.Ques = append(grpcSurvey.Ques, &gen.Question{
			Id:    int32(question.ID),
			Title: question.Title,
		})
	}

	return &gen.GetSurveyResponse{
		Survey: grpcSurvey,
	}, nil
}

func (adh *GrpcCsatHandler) PostAnswers(ctx context.Context, in *gen.PostAnswersRequest) (*gen.PostResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("Received GetSurvey request in microservice",
		zap.String("request_id", requestID),
	)
	if in.AuthHeader == "" {
		logger.AccessLogger.Warn("Missing X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("missing X-CSRF-Token header")),
		)
		return nil, errors.New("missing X-CSRF-Token header")
	}

	tokenString := in.AuthHeader[len("Bearer "):]
	_, err := adh.jwtToken.Validate(tokenString, in.SessionId)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("invalid JWT token")
	}

	userId, err := adh.sessionService.GetUserID(ctx, in.SessionId)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		return nil, errors.New("no active session")
	}

	var creds []domain.PostSurvey
	for _, el := range in.Answer {
		creds = append(creds, domain.PostSurvey{
			QuestionId: int(el.QuestionId),
			Value:      int(el.Values),
		})
	}

	err = adh.usecase.PostSurvey(ctx, creds, userId)
	if err != nil {
		logger.AccessLogger.Warn("GetSurvey failed", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("get survey failed")
	}

	return &gen.PostResponse{
		Response: "Success",
	}, nil
}

func (adh *GrpcCsatHandler) GetStatistics(ctx context.Context, in *gen.Empty) (*gen.GetStatisticsResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("Received GetStatistics request in microservice",
		zap.String("request_id", requestID),
	)

	statistics, err := adh.usecase.GetStatistics(ctx)
	if err != nil {
		logger.AccessLogger.Warn("GetStatistics failed", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("get statistics failed")
	}
	
	var grpcStats []*gen.GetStatistics
	for _, stat := range statistics {
		var grpcMap []*gen.Map
		for key, value := range stat.AnswerNumbers {
			grpcMap = append(grpcMap, &gen.Map{
				Key:    int32(key),
				Values: int32(value),
			})
		}
		grpcStats = append(grpcStats, &gen.GetStatistics{
			Average: stat.Avg,
			Title:   stat.Title,
			Map:     grpcMap,
		})
	}

	return &gen.GetStatisticsResponse{
		Statistics: grpcStats,
	}, nil
}
