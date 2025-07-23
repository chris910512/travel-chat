package main

import (
	"github.com/chris910512/travel-chat/internal/delivery/http/handler"
	"github.com/chris910512/travel-chat/internal/delivery/http/router"
	"github.com/chris910512/travel-chat/internal/infrastructure/database"
	"github.com/chris910512/travel-chat/internal/infrastructure/repository"
	"github.com/chris910512/travel-chat/internal/usecase"
	"github.com/joho/godotenv"
	"log"
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

	// 의존성 주입 (Dependency Injection)
	// Repository 계층
	userRepo := repository.NewUserRepository(db)

	// Usecase 계층
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Handler 계층
	userHandler := handler.NewUserHandler(userUsecase)

	// 라우터 설정
	r := router.SetupRoutes(userHandler)

	// 서버 시작
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
