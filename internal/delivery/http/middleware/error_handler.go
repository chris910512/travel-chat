package middleware

import (
	"github.com/chris910512/travel-chat/internal/delivery/http/response"
	"github.com/chris910512/travel-chat/internal/usecase/errors"
	"github.com/gin-gonic/gin"
	"log"
)

// ErrorHandler - 전역 에러 처리 미들웨어
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(error); ok {
			HandleError(c, err)
		} else {
			response.InternalServerError(c, "알 수 없는 오류가 발생했습니다")
		}
		c.Abort()
	})
}

// HandleError - 에러 타입별 처리
func HandleError(c *gin.Context, err error) {
	log.Printf("Error occurred: %v", err)

	switch {
	case errors.IsUserNotFound(err):
		response.NotFound(c, err.Error())
	case errors.IsEmailAlreadyExists(err):
		response.Conflict(c, err.Error())
	case errors.IsInvalidCredentials(err):
		response.Unauthorized(c, err.Error())
	case errors.IsWeakPassword(err):
		response.BadRequest(c, err.Error())
	case errors.IsInvalidEmail(err):
		response.BadRequest(c, err.Error())
	case errors.IsInvalidTravelDates(err):
		response.BadRequest(c, err.Error())
	case errors.IsPastTravelDate(err):
		response.BadRequest(c, err.Error())
	case errors.IsUnauthorized(err):
		response.Unauthorized(c, err.Error())
	case errors.IsForbidden(err):
		response.Forbidden(c, err.Error())
	default:
		response.InternalServerError(c, "서버 내부 오류가 발생했습니다", err.Error())
	}
}
