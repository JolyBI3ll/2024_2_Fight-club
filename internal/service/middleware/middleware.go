package middleware

import (
	"2024_2_FIGHT-CLUB/internal/service/dsn"
	"2024_2_FIGHT-CLUB/internal/service/images"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type contextKey string

const (
	loggerKey contextKey = "logger"
)

const requestTimeout = 5 * time.Second

func WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, requestTimeout)
}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) (*zap.Logger, error) {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return nil, fmt.Errorf("failed to get logger from context")
	}
	return logger, nil
}

type key int

const RequestIDKey key = 0

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

const (
	requestsPerSecond = 5  // Лимит запросов в секунду для каждого IP
	burstLimit        = 10 // Максимальный «всплеск» запросов
)

var clientLimiters = sync.Map{}

func getLimiter(ip string) *rate.Limiter {
	limiter, exists := clientLimiters.Load(ip)
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(requestsPerSecond), burstLimit)
		clientLimiters.Store(ip, limiter)
	}
	return limiter.(*rate.Limiter)
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		limiter := getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Set-Cookie, X-CSRFToken, x-csrftoken, X-CSRF-Token")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MinioConnect() images.MinioServiceInterface {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	minioService, err := images.NewMinioService(endpoint, accessKey, secretKey, bucketName, useSSL)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}
	fmt.Println("Connected to minio")
	return minioService
}

func DbConnect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to database")
	return db
}

func DbCSATConnect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn.FromEnvCSAT()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to database")
	return db
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
