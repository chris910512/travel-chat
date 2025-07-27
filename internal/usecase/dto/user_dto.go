package dto

import (
	"time"
)

// 사용자 등록 요청
type CreateUserRequest struct {
	Email         string    `json:"email" binding:"required,email"`
	Password      string    `json:"password" binding:"required,min=6"`
	Name          string    `json:"name" binding:"required,min=2"`
	Age           int       `json:"age" binding:"min=18,max=100"`
	Gender        string    `json:"gender" binding:"required,oneof=male female other"`
	Country       string    `json:"country" binding:"required"`
	City          string    `json:"city" binding:"required"`
	TravelStart   time.Time `json:"travel_start" binding:"required"`
	TravelEnd     time.Time `json:"travel_end" binding:"required"`
	Bio           string    `json:"bio" binding:"max=500"`
	TravelPurpose string    `json:"travel_purpose" binding:"required"`
	TravelBudget  int       `json:"travel_budget" binding:"min=0"`
	TravelStyle   string    `json:"travel_style" binding:"required"`
}

// 사용자 로그인 요청
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// 사용자 응답 (비밀번호 제외)
type UserResponse struct {
	ID             uint      `json:"id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	Age            int       `json:"age"`
	Gender         string    `json:"gender"`
	ProfilePic     string    `json:"profile_pic"`
	Country        string    `json:"country"`
	City           string    `json:"city"`
	Destination    string    `json:"destination"` // "국가-도시" 형식
	TravelStart    time.Time `json:"travel_start"`
	TravelEnd      time.Time `json:"travel_end"`
	Bio            string    `json:"bio"`
	TravelPurpose  string    `json:"travel_purpose"`
	TravelBudget   int       `json:"travel_budget"`
	TravelStyle    string    `json:"travel_style"`
	ActivityStatus string    `json:"activity_status"` // "온라인", "10분 전 활동" 등
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}


// 사용자 프로필 업데이트 요청
type UpdateUserRequest struct {
	Name          *string    `json:"name,omitempty"`
	Age           *int       `json:"age,omitempty"`
	Gender        *string    `json:"gender,omitempty"`
	ProfilePic    *string    `json:"profile_pic,omitempty"`
	Country       *string    `json:"country,omitempty"`
	City          *string    `json:"city,omitempty"`
	TravelStart   *time.Time `json:"travel_start,omitempty"`
	TravelEnd     *time.Time `json:"travel_end,omitempty"`
	Bio           *string    `json:"bio,omitempty"`
	TravelPurpose *string    `json:"travel_purpose,omitempty"`
	TravelBudget  *int       `json:"travel_budget,omitempty"`
	TravelStyle   *string    `json:"travel_style,omitempty"`
}

// 사용자 목록 요청 (페이징)
type GetUsersRequest struct {
	Page    int    `form:"page" binding:"min=1"`          // 페이지 번호 (1부터 시작)
	Limit   int    `form:"limit" binding:"min=1,max=100"` // 페이지 크기
	Country string `form:"country"`                       // 국가 필터
	City    string `form:"city"`                          // 도시 필터
}

// 사용자 목록 응답
type GetUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalCount int64          `json:"total_count"`
	TotalPages int            `json:"total_pages"`
}
