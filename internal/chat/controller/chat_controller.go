package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/chat/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
	"sync"
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
	maxConnections    = 100
	messageRateLimit  = 5
)

var (
	upgrader    = websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize, CheckOrigin: func(r *http.Request) bool { return true }}
	mapUserConn = make(map[string]*Client)
	connCounter = 0
	mu          sync.Mutex
)

func (cc *ChatHandler) SetConnection(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	logger.AccessLogger.Info("Received SetConnection request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)
	var err error
	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}

	mu.Lock()
	if connCounter >= maxConnections {
		mu.Unlock()
		http.Error(w, "Too many connections", http.StatusTooManyRequests)
		return
	}
	connCounter++
	mu.Unlock()

	defer func() {
		mu.Lock()
		connCounter--
		mu.Unlock()

		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	sess, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Info("Failed to get sessionId",
			zap.String("request_id", requestID),
			zap.Error(err))
		cc.handleError(w, err, requestID)
		return
	}

	UserID, err := cc.sessionService.GetUserID(r.Context(), sess)
	if err != nil || UserID == "" {
		logger.AccessLogger.Info("Unauthorized user",
			zap.String("request_id", requestID),
			zap.Error(err))
		cc.handleError(w, err, requestID)
		return
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		cc.handleError(w, errors.New("failed to upgrade connection"), requestID)
		logger.AccessLogger.Info("Failed to get socket",
			zap.String("request_id", requestID),
			zap.Error(err))
		return
	}

	client := &Client{
		Socket:         socket,
		Receive:        make(chan *domain.Message, messageBufferSize),
		ChatController: cc,
		RateLimiter:    NewRateLimiter(messageRateLimit, time.Second, 10*time.Second),
	}

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
	logger.AccessLogger.Info("Received GetAllChats request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)
	defer cancel()
	var err error
	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	var (
		lastTimeQuery = r.URL.Query().Get("lastTime")
		lastTime      time.Time
	)

	if lastTimeQuery == "" {
		lastTime = time.Now()
	} else {
		lastTime, err = time.Parse("2006-01-02T15:04:05.000000Z", lastTimeQuery)
		if err != nil {
			logger.AccessLogger.Info("Failed to parse lastTime",
				zap.String("request_id", requestID),
				zap.Error(err))
			cc.handleError(w, errors.New("failed to parse lastTime"), requestID)
			return
		}
	}

	sess, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Info("Failed to get sessionId",
			zap.String("request_id", requestID),
			zap.Error(err))
		cc.handleError(w, err, requestID)
		return
	}
	UserID, err := cc.sessionService.GetUserID(ctx, sess)
	if err != nil {
		logger.AccessLogger.Info("Failed to get UserID",
			zap.String("request_id", requestID),
			zap.Error(err))
		cc.handleError(w, err, requestID)
		return
	}
	chats, err := cc.chatUseCase.GetAllChats(ctx, UserID, lastTime)
	if err != nil {
		logger.AccessLogger.Info("Failed to get all chats",
			zap.String("request_id", requestID),
			zap.Error(err))
		cc.handleError(w, err, requestID)
		return
	}

	body := domain.AllChats{
		Chats: chats,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := easyjson.MarshalToWriter(&body, w); err != nil {
		logger.AccessLogger.Warn("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err))
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

	var (
		lastTimeQuery = r.URL.Query().Get("lastTime")
		lastTime      time.Time
		err           error
	)

	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received RegisterUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	if lastTimeQuery == "" {
		lastTime = time.Now()
	} else {
		lastTime, err = time.Parse("2006-01-02T15:04:05.000000Z", lastTimeQuery)
		if err != nil {
			logger.AccessLogger.Info("Failed to parse lastTime",
				zap.String("request_id", requestID),
				zap.Error(err))
			cc.handleError(w, errors.New("failed to parse lastTime"), requestID)
			return
		}
	}

	sess, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Info("Failed to get sessionId",
			zap.String("request_id", requestID),
			zap.Error(err))
		cc.handleError(w, err, requestID)
		return
	}
	UserID, err := cc.sessionService.GetUserID(ctx, sess)
	if err != nil {
		logger.AccessLogger.Info("Failed to get UserID",
			zap.String("request_id", requestID),
			zap.Error(err))
		cc.handleError(w, err, requestID)
		return
	}
	messages, err := cc.chatUseCase.GetChat(ctx, UserID, id, lastTime)
	if err != nil {
		logger.AccessLogger.Info("Failed to get chat",
			zap.String("request_id", requestID),
			zap.Error(err))
		cc.handleError(w, err, requestID)
		return
	}

	body := domain.AllMessages{
		Chat: messages,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := easyjson.MarshalToWriter(&body, w); err != nil {
		logger.AccessLogger.Warn("Failed to encode response",
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

func (cc *ChatHandler) handleError(w http.ResponseWriter, err error, requestID string) int {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)

	w.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{"error": err.Error()}
	var status int
	switch err.Error() {
	case "error fetching chats", "error fetching messages",
		"failed to generate session id", "failed to save session", "error generating random bytes for session ID",
		"failed to delete session", "failed to get session id from request cookie", "failed to upgrade connection":
		w.WriteHeader(http.StatusInternalServerError)
		status = http.StatusInternalServerError
	case "error sending message",
		"failed to parse lastTime":
		w.WriteHeader(http.StatusBadRequest)
		status = http.StatusBadRequest
	case "session not found", "user ID not found in session":
		w.WriteHeader(http.StatusUnauthorized)
		status = http.StatusUnauthorized
	default:
		w.WriteHeader(http.StatusInternalServerError)
		status = http.StatusInternalServerError
	}

	if jsonErr := json.NewEncoder(w).Encode(errorResponse); jsonErr != nil {
		logger.AccessLogger.Error("Failed to encode error response",
			zap.String("request_id", requestID),
			zap.Error(jsonErr),
		)
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}
	return status
}
