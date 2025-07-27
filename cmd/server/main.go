package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/chris910512/travel-chat/internal/delivery/grpc/server"
	"github.com/chris910512/travel-chat/internal/delivery/http/handler"
	"github.com/chris910512/travel-chat/internal/delivery/http/router"
	"github.com/chris910512/travel-chat/internal/infrastructure/database"
	"github.com/chris910512/travel-chat/internal/infrastructure/repository"
	"github.com/chris910512/travel-chat/internal/pkg/jwt"
	"github.com/chris910512/travel-chat/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	// 환경변수 로드
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 데이터베이스 연결
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 마이그레이션 실행
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// JWT 서비스 초기화
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is required")
	}

	jwtIssuer := os.Getenv("JWT_ISSUER")
	if jwtIssuer == "" {
		jwtIssuer = "travel-chat-api"
	}

	jwtService := jwt.NewJWTService(jwtSecret, jwtIssuer)

	// 의존성 주입 (Dependency Injection)
	// Repository 계층
	userRepo := repository.NewUserRepository(db)

	// Usecase 계층 (JWT 서비스 주입)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtService)

	// HTTP Handler 계층
	userHandler := handler.NewUserHandler(userUsecase)

	// 포트 설정
	httpPort := os.Getenv("SERVER_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}

	gatewayPort := os.Getenv("GATEWAY_PORT")
	if gatewayPort == "" {
		gatewayPort = "8081"
	}

	// HTTP 라우터 설정
	httpRouter := router.SetupRoutes(userHandler, jwtService)

	// gRPC 서버 설정
	grpcServer := server.NewGRPCServer(userUsecase, jwtService, grpcPort, gatewayPort)

	// 서버들을 고루틴으로 동시 실행
	var wg sync.WaitGroup
	wg.Add(3)

	// HTTP 서버 시작
	go func() {
		defer wg.Done()
		log.Printf("HTTP server starting on port %s", httpPort)
		if err := httpRouter.Run(":" + httpPort); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// gRPC 서버 시작
	go func() {
		defer wg.Done()
		if err := grpcServer.StartGRPCServer(); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// gRPC Gateway 서버 시작
	go func() {
		defer wg.Done()
		if err := grpcServer.StartGatewayServer(); err != nil {
			log.Printf("gRPC Gateway server error: %v", err)
		}
	}()

	// Graceful shutdown을 위한 신호 대기
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	log.Println("=== Travel Chat API Server Started ===")
	log.Printf("HTTP Server: http://localhost:%s", httpPort)
	log.Printf("gRPC Server: localhost:%s", grpcPort)
	log.Printf("gRPC Gateway: http://localhost:%s", gatewayPort)
	log.Println("Press Ctrl+C to exit")

	// 종료 신호 대기
	<-c
	log.Println("Shutting down servers...")

	// gRPC 서버 정지
	grpcServer.Stop()

	log.Println("Servers stopped")
}
