package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"context"
	"errors"
)

type CsatUseCase interface {
	GetSurvey(ctx context.Context, surveyId int) (domain.SurveyResponse, error)
	PostSurvey(ctx context.Context, answers []domain.PostSurvey, userId string) error
	GetStatistics(ctx context.Context) ([]domain.GetStatictics, error)
}

type csatUseCase struct {
	csatRepository domain.CSATRepository
}

func NewCSATUseCase(csatRepository domain.CSATRepository) CsatUseCase {
	return &csatUseCase{
		csatRepository: csatRepository,
	}
}

func (uc *csatUseCase) GetSurvey(ctx context.Context, surveyId int) (domain.SurveyResponse, error) {
	survey, err := uc.csatRepository.GetSurvey(ctx, surveyId)
	if err != nil {
		return survey, errors.New("survey not found")
	}
	return survey, nil
}

func (uc *csatUseCase) PostSurvey(ctx context.Context, answers []domain.PostSurvey, userId string) error {
	err := uc.csatRepository.PostSurvey(ctx, answers, userId)
	if err != nil {
		return errors.New("failed to insert")
	}
	return nil
}

func (uc *csatUseCase) GetStatistics(ctx context.Context) ([]domain.GetStatictics, error) {
	statistics, err := uc.csatRepository.GetStatistics(ctx)
	if err != nil {
		return statistics, errors.New("statistics not found")
	}
	return statistics, nil
}
