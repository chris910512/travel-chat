package usecase

import (
	"context"
	"github.com/chris910512/travel-chat/internal/usecase/dto"
)

// UserUsecase 인터페이스 정의
type UserUsecase interface {
	// 사용자 인증
	Register(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)

	// 사용자 조회
	GetByID(ctx context.Context, id uint) (*dto.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	GetUsers(ctx context.Context, req *dto.GetUsersRequest) (*dto.GetUsersResponse, error)
	GetUsersByDestination(ctx context.Context, country, city string) ([]dto.UserResponse, error)

	// 사용자 관리
	UpdateProfile(ctx context.Context, userID uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	UpdateLastActive(ctx context.Context, userID uint) error
	DeleteUser(ctx context.Context, userID uint) error

	// 유틸리티
	ValidateUserExists(ctx context.Context, userID uint) error
}
