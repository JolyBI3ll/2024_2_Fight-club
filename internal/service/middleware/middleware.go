package middleware

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/dsn"
	"2024_2_FIGHT-CLUB/internal/service/images"
	"2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"runtime/debug"
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

func ConvertRoomsToGRPC(rooms []domain.AdRoomsResponse) []*gen.AdRooms {
	var grpcRooms []*gen.AdRooms
	for _, room := range rooms {
		grpcRooms = append(grpcRooms, &gen.AdRooms{
			Type:         room.Type,
			SquareMeters: int32(room.SquareMeters),
		})
	}
	return grpcRooms
}

func ConvertGRPCToRooms(grpc []*gen.AdRooms) []domain.AdRoomsResponse {
	var Rooms []domain.AdRoomsResponse
	for _, room := range grpc {
		Rooms = append(Rooms, domain.AdRoomsResponse{
			Type:         room.Type,
			SquareMeters: int(room.SquareMeters),
		})
	}
	return Rooms
}

type responseWriterWrapper struct {
	http.ResponseWriter
	written bool
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	if !w.written {
		w.written = true
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

func RecoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedWriter := &responseWriterWrapper{ResponseWriter: w}

		defer func() {
			if r := recover(); r != nil {
				var err error
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				if !wrappedWriter.written {
					http.Error(wrappedWriter, err.Error(), http.StatusInternalServerError)
				}
			}
		}()
		h.ServeHTTP(wrappedWriter, r)
	})
}

func RecoveryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic occurred: %v\n", r)
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "internal server error: %v", r)
		}
	}()
	return handler(ctx, req)
}

func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		current := len(interceptors) - 1
		var chain grpc.UnaryHandler
		chain = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
			if current < 0 {
				return handler(currentCtx, currentReq)
			}
			interceptor := interceptors[current]
			current--
			return interceptor(currentCtx, currentReq, info, chain)
		}
		return chain(ctx, req)
	}
}
