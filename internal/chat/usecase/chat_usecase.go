package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"context"
	"time"
)

type ChatUseCase interface {
	GetAllChats(ctx context.Context, userID string, lastUpdateTime time.Time) ([]*domain.Chat, error)
	SendNewMessage(ctx context.Context, receiver string, sender string, message string) error
	GetChat(ctx context.Context, userID1 string, userID2 string, lastSentTime time.Time) ([]*domain.Message, error)
}

type chatUseCase struct {
	repo domain.ChatRepository
}

func NewChatService(repo domain.ChatRepository) ChatUseCase {
	return &chatUseCase{
		repo: repo,
	}
}

func (cs *chatUseCase) GetAllChats(ctx context.Context, userID string, lastUpdateTime time.Time) ([]*domain.Chat, error) {
	chats, err := cs.repo.GetChats(ctx, userID, lastUpdateTime)

	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (cs *chatUseCase) GetChat(ctx context.Context, userID1 string, userID2 string, lastSent time.Time) ([]*domain.Message, error) {
	messages, err := cs.repo.GetMessages(ctx, userID1, userID2, lastSent)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (cs *chatUseCase) SendNewMessage(ctx context.Context, receiver string, sender string, message string) error {
	err := cs.repo.SendNewMessage(ctx, receiver, sender, message)
	if err != nil {
		return err
	}
	return nil
}
