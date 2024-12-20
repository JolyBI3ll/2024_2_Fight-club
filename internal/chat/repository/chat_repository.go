package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"

	"go.uber.org/zap"
)

type Repo struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) domain.ChatRepository {
	return &Repo{
		db: db,
	}
}

func (cr *Repo) GetChats(ctx context.Context, userID string, lastUpdateTime time.Time) ([]*domain.Chat, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetChats called", zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetChats", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetChats", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetChats").Observe(duration)
	}()
	var chats []*domain.Chat

	subquery := cr.db.
		Table("messages").
		Select(`
		CASE WHEN "senderId" = ? THEN "receiverId" ELSE "senderId" END AS related_user,
		content,
		"createdAt"`, userID).
		Where("\"senderId\" = ? OR \"receiverId\" = ?", userID, userID).
		Order("\"createdAt\" DESC")

	latestMessagesQuery := cr.db.
		Table("(?) as all_messages", subquery).
		Select("related_user, MAX(\"createdAt\") as max_date").
		Group("related_user")

	// Соединяем это с основной таблицей для извлечения нужной информации
	err = cr.db.
		Table("(?) as latest_messages", latestMessagesQuery).
		Joins("INNER JOIN messages ON messages.\"createdAt\" = latest_messages.max_date AND messages.\"senderId\" IN (?, latest_messages.related_user) AND messages.\"receiverId\" IN (?, latest_messages.related_user)", userID, userID).
		Joins("INNER JOIN users ON latest_messages.related_user = users.uuid").
		Select(`
		latest_messages.related_user,
		users.name AS "authorName",
		users.avatar AS "authorAvatar",
		users.uuid AS "authorUuid",
		messages.content AS "lastMessage",
		messages."createdAt" AS "lastDate"`).
		Order("\"lastDate\" DESC").
		Limit(15).
		Scan(&chats).Error

	if err != nil {
		logger.DBLogger.Error("Error fetching chats", zap.String("request_id", requestID), zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("error fetching chats")
		}
		return nil, errors.New("error fetching chats")
	}

	logger.DBLogger.Info("Successfully fetched chats", zap.String("request_id", requestID), zap.Int("count", len(chats)))
	return chats, nil
}

func (cr *Repo) GetMessages(ctx context.Context, userID1 string, userID2 string, lastSentTime time.Time) ([]*domain.Message, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetMessages called", zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetMessages", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetMessages", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetMessages").Observe(duration)
	}()
	var messages []*domain.Message
	err = cr.db.
		Where("(\"senderId\" = ? AND \"receiverId\" = ?) OR (\"senderId\" = ? AND \"receiverId\" = ?)", userID1, userID2, userID2, userID1).
		Where("\"createdAt\" < ?", lastSentTime).
		Order("\"createdAt\" ASC").
		Find(&messages).Error

	if err != nil {
		logger.DBLogger.Error("Error fetching messages", zap.String("request_id", requestID), zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("error fetching messages")
		}
		return nil, errors.New("error fetching messages")
	}

	logger.DBLogger.Info("Successfully fetched messages", zap.String("request_id", requestID), zap.Int("count", len(messages)))
	return messages, nil
}

func (cr *Repo) SendNewMessage(ctx context.Context, receiver string, sender string, message string) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("SendNewMessage called", zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("SendNewMessage", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("SendNewMessage", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("SendNewMessage").Observe(duration)
	}()
	newMessage := &domain.Message{
		ReceiverID: receiver,
		SenderID:   sender,
		Content:    message,
		CreatedAt:  time.Now(),
	}

	err = cr.db.Create(newMessage).Error
	if err != nil {
		logger.DBLogger.Error("Error sending message", zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error sending message")
	}

	logger.DBLogger.Info("Successfully sent message", zap.String("request_id", requestID))
	return nil
}
