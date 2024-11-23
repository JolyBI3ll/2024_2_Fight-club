package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type csatRepository struct {
	db *gorm.DB
}

func NewCsatRepository(db *gorm.DB) domain.CSATRepository {
	return &csatRepository{
		db: db,
	}
}

func (r *csatRepository) GetSurvey(ctx context.Context, surveyId int) (domain.SurveyResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetSurvey called",
		zap.Int("surveyId", surveyId),
		zap.String("request_id", requestID),
	)

	var response domain.SurveyResponse
	var survey domain.Survey
	var questions []domain.Question

	if err := r.db.First(&survey, surveyId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Info("Survey not found", zap.Int("surveyId", surveyId))
		}
		return response, err
	}
	response.ID = survey.ID
	response.Title = survey.Title
	if err := r.db.Find(&questions).Where("\"surveyId\" = ?", surveyId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Info("Question not found", zap.Int("surveyId", surveyId))
		}
		return response, err
	}
	for _, question := range questions {
		response.Questions = append(response.Questions, question)
	}

	return response, nil
}

func (r *csatRepository) PostSurvey(ctx context.Context, answers []domain.PostSurvey, userId string) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("PostSurvey called", zap.String("request_id", requestID))

	for _, answer := range answers {
		if err := r.db.Model(&domain.Answer{}).
			Create(&domain.Answer{
				QuestionId: answer.QuestionId,
				Value:      answer.Value,
				UserId:     userId,
			}).Error; err != nil {
			logger.AccessLogger.Warn("Failed to insert", zap.Error(err))
			return err
		}
	}
	return nil
}
