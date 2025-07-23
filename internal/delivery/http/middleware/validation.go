package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chris910512/travel-chat/internal/delivery/http/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidateRequest - 요청 바인딩 및 검증 미들웨어
func ValidateRequest(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(model); err != nil {
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				errorList := make(map[string]string)
				for _, fieldError := range validationErrors {
					field := strings.ToLower(fieldError.Field())
					switch fieldError.Tag() {
					case "required":
						errorList[field] = fmt.Sprintf("%s는 필수 항목입니다", field)
					case "email":
						errorList[field] = "올바른 이메일 형식이 아닙니다"
					case "min":
						errorList[field] = fmt.Sprintf("%s는 최소 %s자 이상이어야 합니다", field, fieldError.Param())
					case "max":
						errorList[field] = fmt.Sprintf("%s는 최대 %s자 이하여야 합니다", field, fieldError.Param())
					case "oneof":
						errorList[field] = fmt.Sprintf("%s는 허용된 값이 아닙니다", field)
					default:
						errorList[field] = fmt.Sprintf("%s가 올바르지 않습니다", field)
					}
				}
				response.ValidationError(c, errorList)
				c.Abort()
				return
			}

			response.BadRequest(c, "요청 형식이 올바르지 않습니다", err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}

// BindQueryParams - 쿼리 파라미터 바인딩
func BindQueryParams(c *gin.Context, model interface{}) error {
	if err := c.ShouldBindQuery(model); err != nil {
		return err
	}
	return nil
}
