package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"fmt"
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

func (r *csatRepository) GetStatistics(ctx context.Context) ([]domain.GetStatictics, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetStatistics called", zap.String("request_id", requestID))

	var questions []domain.Question
	if err := r.db.Find(&questions).Error; err != nil {
		logger.DBLogger.Error("Error fetching questions", zap.String("request_id", requestID), zap.Error(err))
		return nil, fmt.Errorf("error fetching questions: %w", err)
	}

	var statistics []domain.GetStatictics

	for _, question := range questions {
		var answers []domain.Answer
		if err := r.db.Where("\"questionId\" = ?", question.ID).Find(&answers).Error; err != nil {
			logger.DBLogger.Error("Error fetching answers", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("error fetching answers for question")
		}

		stat := domain.GetStatictics{
			Title:         question.Title,
			AnswerNumbers: make(map[int]int),
		}

		maxValue := 0
		switch question.Type {
		case "STARS":
			maxValue = 5
		case "SMILE":
			maxValue = 5
		case "RATE":
			maxValue = 10
		default:
			logger.DBLogger.Warn("Unknown question type", zap.String("request_id", requestID), zap.String("type", question.Type))
			continue
		}

		for i := 1; i <= maxValue; i++ {
			stat.AnswerNumbers[i] = 0
		}

		var totalValue int
		for _, answer := range answers {
			if answer.Value >= 1 && answer.Value <= maxValue {
				stat.AnswerNumbers[answer.Value]++
				totalValue += answer.Value
			}
		}

		if len(answers) > 0 {
			stat.Avg = float32(totalValue) / float32(len(answers))
		} else {
			stat.Avg = 0
		}

		statistics = append(statistics, stat)
	}

	return statistics, nil
}
