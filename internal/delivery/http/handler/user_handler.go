package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/chris910512/travel-chat/internal/delivery/http/response"
	"github.com/chris910512/travel-chat/internal/usecase/dto"
	usecaseErrors "github.com/chris910512/travel-chat/internal/usecase/errors"
	usecaseInterface "github.com/chris910512/travel-chat/internal/usecase/interface"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase usecaseInterface.UserUsecase
}

// NewUserHandler - User Handler 생성자
func NewUserHandler(userUsecase usecaseInterface.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// RefreshToken - 토큰 갱신
// POST /api/auth/refresh
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	// 요청 바인딩
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "요청 형식이 올바르지 않습니다", err.Error())
		return
	}

	// 토큰 갱신 처리
	refreshResp, err := h.userUsecase.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		handleUsecaseError(c, err)
		return
	}

	response.Success(c, "토큰이 갱신되었습니다", refreshResp)
}

// Register - 사용자 등록
// POST /api/users/register
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.CreateUserRequest

	// 요청 바인딩
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "요청 형식이 올바르지 않습니다", err.Error())
		return
	}

	// 사용자 등록 처리
	user, err := h.userUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		handleUsecaseError(c, err)
		return
	}

	response.Created(c, "사용자가 성공적으로 등록되었습니다", user)
}

// Login - 사용자 로그인
// POST /api/users/login
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// 요청 바인딩
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "요청 형식이 올바르지 않습니다", err.Error())
		return
	}

	// 로그인 처리
	loginResp, err := h.userUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		handleUsecaseError(c, err)
		return
	}

	response.Success(c, "로그인이 성공했습니다", loginResp)
}

// GetProfile - 사용자 프로필 조회
// GET /api/users/:id
func (h *UserHandler) GetProfile(c *gin.Context) {
	// URL 파라미터에서 사용자 ID 추출
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "올바르지 않은 사용자 ID입니다")
		return
	}

	// 사용자 조회
	user, err := h.userUsecase.GetByID(c.Request.Context(), uint(userID))
	if err != nil {
		handleUsecaseError(c, err)
		return
	}

	response.Success(c, "사용자 정보를 조회했습니다", user)
}

// GetUsers - 사용자 목록 조회
// GET /api/users?page=1&limit=10&country=일본&city=도쿄
func (h *UserHandler) GetUsers(c *gin.Context) {
	var req dto.GetUsersRequest

	// 쿼리 파라미터 바인딩
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "쿼리 파라미터가 올바르지 않습니다", err.Error())
		return
	}

	// 기본값 설정
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// 사용자 목록 조회
	users, err := h.userUsecase.GetUsers(c.Request.Context(), &req)
	if err != nil {
		handleUsecaseError(c, err)
		return
	}

	response.Success(c, "사용자 목록을 조회했습니다", users)
}

// GetUsersByDestination - 목적지별 사용자 조회
// GET /api/users/destination/:country/:city
func (h *UserHandler) GetUsersByDestination(c *gin.Context) {
	country := c.Param("country")
	city := c.Param("city")

	if country == "" || city == "" {
		response.BadRequest(c, "국가와 도시를 모두 입력해주세요")
		return
	}

	// 목적지별 사용자 조회
	users, err := h.userUsecase.GetUsersByDestination(c.Request.Context(), country, city)
	if err != nil {
		handleUsecaseError(c, err)
		return
	}

	response.Success(c, "목적지별 사용자 목록을 조회했습니다", users)
}

// UpdateProfile - 사용자 프로필 업데이트
// PUT /api/users/:id
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// URL 파라미터에서 사용자 ID 추출
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "올바르지 않은 사용자 ID입니다")
		return
	}

	var req dto.UpdateUserRequest

	// 요청 바인딩
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "요청 형식이 올바르지 않습니다", err.Error())
		return
	}

	// 프로필 업데이트
	user, err := h.userUsecase.UpdateProfile(c.Request.Context(), uint(userID), &req)
	if err != nil {
		handleUsecaseError(c, err)
		return
	}

	response.Success(c, "프로필이 업데이트되었습니다", user)
}

// UpdateLastActive - 마지막 활동 시간 업데이트
// POST /api/users/:id/activity
func (h *UserHandler) UpdateLastActive(c *gin.Context) {
	// URL 파라미터에서 사용자 ID 추출
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "올바르지 않은 사용자 ID입니다")
		return
	}

	// 마지막 활동 시간 업데이트
	if err := h.userUsecase.UpdateLastActive(c.Request.Context(), uint(userID)); err != nil {
		handleUsecaseError(c, err)
		return
	}

	response.Success(c, "활동 시간이 업데이트되었습니다", nil)
}

// DeleteUser - 사용자 삭제
// DELETE /api/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// URL 파라미터에서 사용자 ID 추출
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "올바르지 않은 사용자 ID입니다")
		return
	}

	// 사용자 삭제
	if err := h.userUsecase.DeleteUser(c.Request.Context(), uint(userID)); err != nil {
		handleUsecaseError(c, err)
		return
	}

	response.Success(c, "사용자가 삭제되었습니다", nil)
}

// GetMe - 현재 로그인한 사용자 정보 조회 (JWT 토큰 기반)
// GET /api/users/me
func (h *UserHandler) GetMe(c *gin.Context) {
	// JWT 미들웨어에서 설정한 사용자 ID 가져오기
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "인증이 필요합니다")
		return
	}

	id, ok := userID.(uint)
	if !ok {
		response.Unauthorized(c, "올바르지 않은 인증 정보입니다")
		return
	}

	// 사용자 조회
	user, err := h.userUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		handleUsecaseError(c, err)
		return
	}

	response.Success(c, "내 프로필 정보를 조회했습니다", user)
}

// Health Check - 서버 상태 확인
// GET /api/health
func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Travel Chat API Server is running",
		"version": "1.0.0",
	})
}

// handleUsecaseError - Usecase 에러를 HTTP 응답으로 변환하는 헬퍼 함수
func handleUsecaseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecaseErrors.ErrUserNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, usecaseErrors.ErrEmailAlreadyExists):
		response.Conflict(c, err.Error())
	case errors.Is(err, usecaseErrors.ErrInvalidCredentials):
		response.Unauthorized(c, err.Error())
	case errors.Is(err, usecaseErrors.ErrWeakPassword):
		response.BadRequest(c, err.Error())
	case errors.Is(err, usecaseErrors.ErrInvalidEmail):
		response.BadRequest(c, err.Error())
	case errors.Is(err, usecaseErrors.ErrInvalidTravelDates):
		response.BadRequest(c, err.Error())
	case errors.Is(err, usecaseErrors.ErrPastTravelDate):
		response.BadRequest(c, err.Error())
	case errors.Is(err, usecaseErrors.ErrUnauthorized):
		response.Unauthorized(c, err.Error())
	case errors.Is(err, usecaseErrors.ErrForbidden):
		response.Forbidden(c, err.Error())
	default:
		response.InternalServerError(c, "서버 내부 오류가 발생했습니다", err.Error())
	}
}
