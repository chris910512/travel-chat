package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// APIResponse - 표준 API 응답 구조체
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo - 에러 정보 구조체
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// 성공 응답 헬퍼 함수들

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func NoContent(c *gin.Context, message string) {
	c.JSON(http.StatusNoContent, APIResponse{
		Success: true,
		Message: message,
	})
}

// 에러 응답 헬퍼 함수들

func BadRequest(c *gin.Context, message string, details ...string) {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}

	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: "요청이 올바르지 않습니다",
		Error: &ErrorInfo{
			Code:    "BAD_REQUEST",
			Message: message,
			Details: detail,
		},
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Message: "인증이 필요합니다",
		Error: &ErrorInfo{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, APIResponse{
		Success: false,
		Message: "접근 권한이 없습니다",
		Error: &ErrorInfo{
			Code:    "FORBIDDEN",
			Message: message,
		},
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Message: "리소스를 찾을 수 없습니다",
		Error: &ErrorInfo{
			Code:    "NOT_FOUND",
			Message: message,
		},
	})
}

func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, APIResponse{
		Success: false,
		Message: "요청이 충돌합니다",
		Error: &ErrorInfo{
			Code:    "CONFLICT",
			Message: message,
		},
	})
}

func InternalServerError(c *gin.Context, message string, details ...string) {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}

	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Message: "내부 서버 오류가 발생했습니다",
		Error: &ErrorInfo{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: message,
			Details: detail,
		},
	})
}

// ValidationError - 유효성 검사 오류 응답
func ValidationError(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: "입력값이 올바르지 않습니다",
		Error: &ErrorInfo{
			Code:    "VALIDATION_ERROR",
			Message: "입력값을 확인해주세요",
		},
		Data: errors,
	})
}
