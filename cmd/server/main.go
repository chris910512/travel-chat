package main

import (
	"log"
	"os"

	"github.com/chris910512/travel-chat/internal/infrastructure/database"
	"github.com/chris910512/travel-chat/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
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

	// TODO: 의존성 주입 및 라우터 설정
	// 현재는 기본 라우터만 설정
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Travel Chat API Server is running",
			"version": "1.0.0",
		})
	})

	// JWT 테스트 엔드포인트 추가
	r.POST("/test/token", func(c *gin.Context) {
		// 테스트용 토큰 생성
		accessToken, err := jwtService.GenerateToken(1, "test@example.com")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		refreshToken, err := jwtService.GenerateRefreshToken(1, "test@example.com")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    86400, // 24시간
		})
	})

	// 보호된 엔드포인트 테스트
	protected := r.Group("/test/protected")
	protected.Use(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token := authHeader[7:] // "Bearer " 제거
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	})

	protected.GET("/me", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		userEmail, _ := c.Get("user_email")

		c.JSON(200, gin.H{
			"user_id": userID,
			"email":   userEmail,
			"message": "You are authenticated!",
		})
	})

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
