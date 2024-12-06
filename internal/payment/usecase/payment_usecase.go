package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"os"
)

type PaymentUseCase interface {
	PaymentCreate(ctx context.Context, adId string, userId string, amount string) (string, error)
}

type paymentUseCase struct {
	repository domain.PaymentRepository
}

const (
	apiURL    = "https://yoomoney.ru/api/request-payment"
	wallet    = "410014464355568"
	PatternID = "p2p"
)

func (uc *paymentUseCase) PaymentCreate(ctx context.Context, adId string, userId string, amount string) (string, error) {
	comment := fmt.Sprintf("Оплата услуги продвижения объявления %s", adId)
	client_id := os.Getenv("CLIENT_ID")
	requestID := middleware.GetRequestID(ctx)

	payment := domain.PaymentRequest{
		PatternID: PatternID,
		To:        wallet,
		Amount:    amount,
		Comment:   comment,
	}

	formData := fmt.Sprintf(
		"pattern_id=%s&to=%s&amount_due=%s&comment=%s",
		payment.PatternID, payment.To, payment.Amount, payment.Comment,
	)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+client_id).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(formData).
		Post(apiURL)
	if err != nil {
		logger.AccessLogger.Error("Failed send post request",
			zap.String("request_id", requestID),
			zap.String("error", err.Error()))
		return "", errors.New("failed to send post request")
	}

	// Обработка ответа
	var response domain.PaymentResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		logger.AccessLogger.Error("Failed unmarshal response body",
			zap.String("request_id", requestID),
			zap.String("error", err.Error()))
		return "", errors.New("failed to unmarshal response body")
	}

	if response.Status == "success" {
		return response.Confirmation.ConfirmationURL, nil
	} else {
		return "", errors.New("failed to create payment")
	}
}
