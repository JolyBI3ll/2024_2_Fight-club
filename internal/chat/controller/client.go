package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	mu         sync.Mutex
	tokens     int           // Текущее количество доступных токенов
	maxTokens  int           // Максимальное количество токенов
	interval   time.Duration // Интервал пополнения токенов
	lastFill   time.Time     // Последнее время пополнения токенов
	blockUntil time.Time     // Время, до которого пользователь заблокирован
	blockTime  time.Duration // Длительность блокировки
}

// NewRateLimiter создает RateLimiter с заданным лимитом, интервалом и временем блокировки.
func NewRateLimiter(maxTokens int, refillInterval, blockTime time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:    maxTokens,
		maxTokens: maxTokens,
		interval:  refillInterval,
		blockTime: blockTime,
		lastFill:  time.Now(),
	}
}

// Allow проверяет, можно ли совершить действие.
func (rl *RateLimiter) Allow() (bool, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Проверяем, заблокирован ли пользователь
	if now.Before(rl.blockUntil) {
		return false, fmt.Errorf("user is blocked until %s", rl.blockUntil.Format(time.RFC3339))
	}

	// Пополняем токены
	elapsed := now.Sub(rl.lastFill)
	if elapsed > rl.interval {
		newTokens := int(elapsed / rl.interval)
		rl.tokens += newTokens
		if rl.tokens > rl.maxTokens {
			rl.tokens = rl.maxTokens
		}
		rl.lastFill = now
	}

	// Проверяем наличие токенов
	if rl.tokens > 0 {
		rl.tokens-- // Используем токен
		return true, nil
	}

	// Если нет токенов, блокируем пользователя
	rl.blockUntil = now.Add(rl.blockTime)
	return false, fmt.Errorf("rate limit exceeded, user blocked for %s", rl.blockTime)
}

type Client struct {
	Socket         *websocket.Conn
	Receive        chan *domain.Message
	ChatController *ChatHandler
	RateLimiter    *RateLimiter
}

func (c *Client) Read(userID string) {
	defer c.Socket.Close()

	// Создаем таймер для отслеживания времени бездействия
	idleTimeout := 20 * time.Minute
	lastActive := time.Now()
	closeOnIdle := make(chan struct{})

	go func() {
		// Закрываем соединение при исчерпании таймера
		for {
			select {
			case <-closeOnIdle:
				return
			case <-time.After(10 * time.Second): // Проверяем тайм-аут каждые 10 секунд
				if time.Since(lastActive) > idleTimeout {
					logger.AccessLogger.Info("Connection closed due to inactivity",
						zap.String("user_id", userID))
					c.Socket.Close()
					return
				}
			}
		}
	}()

	for {
		msg := &domain.Message{}

		// Читаем JSON-сообщение из сокета
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.AccessLogger.Info("Unexpected socket closure",
					zap.String("user_id", userID),
					zap.Error(err))
			}
			break
		}

		// Проверяем, является ли сообщение пустым или состоит только из пробелов/переводов строк/табуляций
		if strings.TrimSpace(msg.Content) == "" {
			// Логируем и отправляем клиенту сообщение об ошибке
			logger.AccessLogger.Warn("Invalid message content received (only whitespace characters)",
				zap.String("user_id", userID))
			errMsg := map[string]string{
				"response": "Invalid message content: only whitespace characters, tabs, or line breaks are not allowed.",
				"sent":     "false",
			}
			if writeErr := c.Socket.WriteJSON(errMsg); writeErr != nil {
				logger.AccessLogger.Error("Failed to send invalid message content error to client",
					zap.String("user_id", userID),
					zap.Error(writeErr))
			}
			continue
		}

		// Обновляем последнее время активности
		lastActive = time.Now()

		// Проверяем лимит сообщений
		allowed, rateErr := c.RateLimiter.Allow()
		if !allowed {
			logger.AccessLogger.Error("Rate limit exceeded",
				zap.String("user_id", userID),
				zap.Error(rateErr))

			// Сообщаем клиенту о блокировке
			errMsg := map[string]string{
				"response": rateErr.Error(),
				"sent":     "false",
			}
			if writeErr := c.Socket.WriteJSON(errMsg); writeErr != nil {
				logger.AccessLogger.Error("Failed to send rate limit error to client",
					zap.String("user_id", userID),
					zap.Error(writeErr))
			}
			continue
		}

		// Устанавливаем SenderID на основе текущего пользователя
		msg.SenderID = userID

		// Отправляем сообщение в канал для обработки
		select {
		case c.ChatController.Messages <- msg:
			// Возвращаем успешный ответ клиенту
			successMsg := map[string]string{
				"response": "Message delivered successfully",
				"sent":     "true",
			}
			if writeErr := c.Socket.WriteJSON(successMsg); writeErr != nil {
				logger.AccessLogger.Error("Failed to send success message to client",
					zap.String("user_id", userID),
					zap.Error(writeErr))
			}
		default:
			logger.AccessLogger.Warn("Message channel is full, dropping message",
				zap.String("user_id", userID))
			errMsg := map[string]string{
				"response": "Message channel is full. Please try again later.",
				"sent":     "false",
			}
			if writeErr := c.Socket.WriteJSON(errMsg); writeErr != nil {
				logger.AccessLogger.Error("Failed to send channel full error to client",
					zap.String("user_id", userID),
					zap.Error(writeErr))
			}
		}
	}

	// Завершаем горутину таймера при выходе из цикла
	close(closeOnIdle)
}

func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Receive {
		msg.CreatedAt = time.Now()
		jsonForSend, err := easyjson.Marshal(msg)
		if err != nil {
			return
		}
		err = c.Socket.WriteMessage(websocket.TextMessage, jsonForSend)
		if err != nil {
			return
		}
	}
}
