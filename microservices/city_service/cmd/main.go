package main

import (
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	grpcCity "2024_2_FIGHT-CLUB/microservices/city_service/controller"
	generatedCity "2024_2_FIGHT-CLUB/microservices/city_service/controller/gen"
	cityRepository "2024_2_FIGHT-CLUB/microservices/city_service/repository"
	cityUseCase "2024_2_FIGHT-CLUB/microservices/city_service/usecase"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Инициализация зависимостей
	middleware.InitRedis()
	db := middleware.DbConnect()

	// Инициализация метрик
	metrics.InitMetrics()
	metrics.InitRepoMetric()
	// Экспозиция метрик на порту 9093
	go func() {
		http.Handle("/api/metrics", promhttp.Handler())
		log.Println("Metrics server is running on :9093")
		if err := http.ListenAndServe(":9093", nil); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// Инициализация логгеров
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		if err := logger.SyncLoggers(); err != nil {
			log.Fatalf("Failed to sync loggers: %v", err)
		}
	}()

	citiesRepository := cityRepository.NewCityRepository(db)
	citiesUseCase := cityUseCase.NewCityUseCase(citiesRepository)
	cityServer := grpcCity.NewGrpcCityHandler(citiesUseCase)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryMetricsInterceptor),
	)
	generatedCity.RegisterCityServiceServer(grpcServer, cityServer)

	// Запуск gRPC сервера
	listener, err := net.Listen("tcp", os.Getenv("CITY_SERVICE_ADDRESS"))
	if err != nil {
		log.Fatalf("Failed to listen on address: %s %v", os.Getenv("CITY_SERVICE_ADDRESS"), err)
	}

	log.Printf("CityService is running on address: %s", os.Getenv("CITY_SERVICE_ADDRESS"))
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
