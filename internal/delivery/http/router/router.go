package router

import (
	"github.com/chris910512/travel-chat/internal/delivery/http/handler"
	"github.com/chris910512/travel-chat/internal/delivery/http/middleware"
	"github.com/chris910512/travel-chat/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// SetupRoutes - 라우터 설정
func SetupRoutes(
	userHandler *handler.UserHandler,
	jwtService *jwt.JWTService,
) *gin.Engine {
	// Gin 엔진 생성
	r := gin.Default()

	// 전역 미들웨어 설정
	r.Use(middleware.ErrorHandler())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS 설정 (개발용)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API 라우트 그룹
	api := r.Group("/api")
	{
		// Health Check
		api.GET("/health", userHandler.HealthCheck)

		// Auth 라우트 (인증 불필요)
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", userHandler.Register)
			authRoutes.POST("/login", userHandler.Login)
			authRoutes.POST("/refresh", userHandler.RefreshToken)
		}

		// 사용자 관련 라우트
		userRoutes := api.Group("/users")
		{
			// 사용자 조회 (공개)
			userRoutes.GET("", userHandler.GetUsers)
			userRoutes.GET("/:id", userHandler.GetProfile)
			userRoutes.GET("/destination/:country/:city", userHandler.GetUsersByDestination)

			// 인증 필요한 엔드포인트
			authenticated := userRoutes.Group("/").Use(middleware.AuthMiddleware(jwtService))
			{
				authenticated.GET("/me", userHandler.GetMe)
				authenticated.PUT("/:id", userHandler.UpdateProfile)
				authenticated.POST("/:id/activity", userHandler.UpdateLastActive)
				authenticated.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}

	return r
}
