package middleware

import (
	"net/http"
	"strings"

	"github.com/chris910512/travel-chat/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware - JWT 인증 미들웨어
func AuthMiddleware(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization 헤더에서 토큰 추출
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header is required",
				"error": map[string]string{
					"code":    "MISSING_TOKEN",
					"message": "인증 토큰이 필요합니다",
				},
			})
			c.Abort()
			return
		}

		// Bearer 토큰 형식 확인
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid token format",
				"error": map[string]string{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Bearer 토큰 형식이 올바르지 않습니다",
				},
			})
			c.Abort()
			return
		}

		// 토큰 검증
		token := tokenParts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
				"error": map[string]string{
					"code":    "INVALID_TOKEN",
					"message": "토큰이 유효하지 않거나 만료되었습니다",
				},
			})
			c.Abort()
			return
		}

		// 사용자 정보를 컨텍스트에 저장
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}

// OptionalAuthMiddleware - 선택적 인증 미들웨어 (토큰이 있으면 검증, 없어도 통과)
func OptionalAuthMiddleware(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
			token := tokenParts[1]
			if claims, err := jwtService.ValidateToken(token); err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
			}
		}

		c.Next()
	}
}

// GetCurrentUserID - 현재 인증된 사용자 ID 가져오기
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}

// GetCurrentUserEmail - 현재 인증된 사용자 이메일 가져오기
func GetCurrentUserEmail(c *gin.Context) (string, bool) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	email, ok := userEmail.(string)
	return email, ok
}
