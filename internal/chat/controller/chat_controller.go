package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/chat/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ChatHandler struct {
	chatUseCase    usecase.ChatUseCase
	sessionService session.InterfaceSession
	Messages       chan *domain.Message
}

func NewChatController(chatUseCase usecase.ChatUseCase, sessionService session.InterfaceSession) *ChatHandler {
	return &ChatHandler{
		chatUseCase:    chatUseCase,
		sessionService: sessionService,
		Messages:       make(chan *domain.Message),
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var (
	upgrader    = websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize, CheckOrigin: func(r *http.Request) bool { return true }}
	mapUserConn = make(map[string]*Client)
)

func (cc *ChatHandler) SetConnection(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	logger.AccessLogger.Info("Received RegisterUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	sess, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Info("Failed to get sessionId",
			zap.String("request_id", requestID),
			zap.Error(err))
		return
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.AccessLogger.Info("Failed to get socket",
			zap.String("request_id", requestID),
			zap.Error(err))
		return
	}
	client := &Client{
		Socket:         socket,
		Receive:        make(chan *domain.Message, messageBufferSize),
		ChatController: cc,
	}

	UserID, err := cc.sessionService.GetUserID(r.Context(), sess)
	mapUserConn[UserID] = client
	defer func() {
		delete(mapUserConn, UserID)
		close(client.Receive)
	}()
	go client.Write()
	go client.Read(UserID)
	cc.SendChatMsg(r.Context(), requestID)

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed SetConnection request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (cc *ChatHandler) SendChatMsg(ctx context.Context, reqID string) {
	start := time.Now()
	ctx, cancel := middleware.WithTimeout(ctx)
	defer cancel()

	logger.AccessLogger.Info("Received SendChatMsg request",
		zap.String("request_id", reqID))

	for msg := range cc.Messages {
		err := cc.chatUseCase.SendNewMessage(ctx, msg.ReceiverID, msg.SenderID, msg.Content)
		if err != nil {
			logger.AccessLogger.Info("Error sending message",
				zap.String("request_id", reqID),
				zap.Error(err))
			return
		}

		resConn, ok := mapUserConn[msg.ReceiverID]
		if ok {
			resConn.Receive <- msg
		}
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed SendChatMsg request",
		zap.String("request_id", reqID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)

	return
}

func (cc *ChatHandler) GetAllChats(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received RegisterUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	var (
		lastTimeQuery = r.URL.Query().Get("lastTime")
		lastTime      time.Time
		err           error
	)

	if lastTimeQuery == "" {
		lastTime = time.Now()
	} else {
		lastTime, err = time.Parse("2006-01-02T15:04:05.000000Z", lastTimeQuery)
		if err != nil {
			logger.AccessLogger.Info("Failed to parse lastTime",
				zap.String("request_id", requestID),
				zap.Error(err))
			return
		}
	}

	sess, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Info("Failed to get sessionId",
			zap.String("request_id", requestID),
			zap.Error(err))
		return
	}
	UserID, err := cc.sessionService.GetUserID(ctx, sess)
	chats, err := cc.chatUseCase.GetAllChats(ctx, UserID, lastTime)
	if err != nil {
		logger.AccessLogger.Info("Failed to get all chats",
			zap.String("request_id", requestID),
			zap.Error(err))
		return
	}

	body := map[string]interface{}{
		"chats": chats,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		return
	}
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetAllChats request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (cc *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received RegisterUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	var (
		lastTimeQuery = r.URL.Query().Get("lastTime")
		lastTime      time.Time
		err           error
	)

	if lastTimeQuery == "" {
		lastTime = time.Now()
	} else {
		lastTime, err = time.Parse("2006-01-02T15:04:05.000000Z", lastTimeQuery)
		if err != nil {
			logger.AccessLogger.Info("Failed to parse lastTime",
				zap.String("request_id", requestID),
				zap.Error(err))
			return
		}
	}

	sess, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Info("Failed to get sessionId",
			zap.String("request_id", requestID),
			zap.Error(err))
		return
	}
	UserID, err := cc.sessionService.GetUserID(ctx, sess)
	messages, err := cc.chatUseCase.GetChat(ctx, UserID, id, lastTime)
	if err != nil {
		logger.AccessLogger.Info("Failed to get chat",
			zap.String("request_id", requestID),
			zap.Error(err))
		return
	}

	body := map[string]interface{}{
		"chat": messages,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		return
	}
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetAllChats request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}
