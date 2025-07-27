package main

import (
	"log"
	"os"

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
	//chatRoomRepo := repository.NewChatRoomRepository(db)
	//messageRepo := repository.NewMessageRepository(db)

	// Usecase 계층 (JWT 서비스 주입)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtService)

	// Handler 계층
	userHandler := handler.NewUserHandler(userUsecase)

	// 라우터 설정 (JWT 서비스 주입)
	r := router.SetupRoutes(userHandler, jwtService)

	// 서버 시작
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
