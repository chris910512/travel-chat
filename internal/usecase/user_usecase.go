package usecase

import (
	"context"
	"github.com/chris910512/travel-chat/internal/domain/entity/user"
	"github.com/chris910512/travel-chat/internal/pkg/jwt"
	"time"

	"github.com/chris910512/travel-chat/internal/domain/repository"
	"github.com/chris910512/travel-chat/internal/usecase/dto"
	"github.com/chris910512/travel-chat/internal/usecase/errors"
	usecaseInterface "github.com/chris910512/travel-chat/internal/usecase/interface"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userUsecase struct {
	userRepo   repository.UserRepository
	jwtService *jwt.JWTService
}

// NewUserUsecase - User Usecase 생성자
func NewUserUsecase(userRepo repository.UserRepository, jwtService *jwt.JWTService) usecaseInterface.UserUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Register - 사용자 등록
func (u *userUsecase) Register(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// 1. 요청 검증
	if err := u.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	// 2. 이메일 중복 체크
	existingUser, err := u.userRepo.GetByEmail(req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.ErrEmailAlreadyExists
	}

	// 3. DTO를 엔티티로 변환
	userEntity, err := req.ToEntity()
	if err != nil {
		return nil, err
	}

	// 4. 사용자 생성
	if err := u.userRepo.Create(userEntity); err != nil {
		return nil, err
	}

	// 5. 응답 반환
	return dto.FromUserEntity(userEntity), nil
}

// Login - 사용자 로그인
func (u *userUsecase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. 이메일로 사용자 조회
	userEntity, err := u.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrInvalidCredentials
		}
		return nil, err
	}

	// 2. 비밀번호 검증
	if err := bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(req.Password)); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// 3. 마지막 활동 시간 업데이트
	userEntity.UpdateLastActive()
	if err := u.userRepo.Update(userEntity); err != nil {
		return nil, err
	}

	// 4. JWT 토큰 생성 (임시로 더미 토큰 반환, 나중에 JWT 구현)
	accessToken, err := u.jwtService.GenerateToken(userEntity.ID, userEntity.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.jwtService.GenerateRefreshToken(userEntity.ID, userEntity.Email)
	if err != nil {
		return nil, err
	}

	// 5. 응답 반환
	return &dto.LoginResponse{
		User:         *dto.FromUserEntity(userEntity),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400, // 24시간
	}, nil
}

// RefreshToken - 토큰 갱신 메서드 추가
func (u *userUsecase) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	accessToken, refreshToken, err := u.jwtService.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400,
	}, nil
}

// GetByID - ID로 사용자 조회
func (u *userUsecase) GetByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	userEntity, err := u.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}

	return dto.FromUserEntity(userEntity), nil
}

// GetByEmail - 이메일로 사용자 조회
func (u *userUsecase) GetByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	userEntity, err := u.userRepo.GetByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}

	return dto.FromUserEntity(userEntity), nil
}

// GetUsers - 사용자 목록 조회 (페이징)
func (u *userUsecase) GetUsers(ctx context.Context, req *dto.GetUsersRequest) (*dto.GetUsersResponse, error) {
	// 기본값 설정
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	var users []*user.User
	var totalCount int64
	var err error

	// 필터링 조건에 따라 조회
	if req.Country != "" && req.City != "" {
		users, err = u.userRepo.GetByDestination(req.Country, req.City)
		if err != nil {
			return nil, err
		}
		totalCount = int64(len(users))

		// 메모리에서 페이징
		offset := req.GetOffset()
		end := offset + req.Limit
		if offset >= len(users) {
			users = []*user.User{}
		} else {
			if end > len(users) {
				end = len(users)
			}
			users = users[offset:end]
		}
	} else {
		// 전체 카운트 조회
		totalCount, err = u.userRepo.Count()
		if err != nil {
			return nil, err
		}

		users, err = u.userRepo.List(req.GetOffset(), req.Limit)
		if err != nil {
			return nil, err
		}
	}

	return &dto.GetUsersResponse{
		Users:      dto.FromUserEntities(users),
		Page:       req.Page,
		Limit:      req.Limit,
		TotalCount: totalCount,
		TotalPages: dto.CalculateTotalPages(totalCount, req.Limit),
	}, nil
}

// GetUsersByDestination - 목적지별 사용자 조회
func (u *userUsecase) GetUsersByDestination(ctx context.Context, country, city string) ([]dto.UserResponse, error) {
	users, err := u.userRepo.GetByDestination(country, city)
	if err != nil {
		return nil, err
	}

	return dto.FromUserEntities(users), nil
}

// UpdateProfile - 사용자 프로필 업데이트
func (u *userUsecase) UpdateProfile(ctx context.Context, userID uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// 1. 기존 사용자 조회
	userEntity, err := u.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, err
	}

	// 2. 업데이트 요청 적용
	req.ApplyToEntity(userEntity)

	// 3. 여행 날짜 검증
	if err := u.validateTravelDates(userEntity.TravelStart, userEntity.TravelEnd); err != nil {
		return nil, err
	}

	// 4. 사용자 업데이트
	if err := u.userRepo.Update(userEntity); err != nil {
		return nil, err
	}

	return dto.FromUserEntity(userEntity), nil
}

// UpdateLastActive - 마지막 활동 시간 업데이트
func (u *userUsecase) UpdateLastActive(ctx context.Context, userID uint) error {
	return u.userRepo.UpdateLastActive(userID)
}

// DeleteUser - 사용자 삭제
func (u *userUsecase) DeleteUser(ctx context.Context, userID uint) error {
	// 사용자 존재 여부 확인
	if err := u.ValidateUserExists(ctx, userID); err != nil {
		return err
	}

	return u.userRepo.Delete(userID)
}

// ValidateUserExists - 사용자 존재 여부 검증
func (u *userUsecase) ValidateUserExists(ctx context.Context, userID uint) error {
	_, err := u.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUserNotFound
		}
		return err
	}
	return nil
}

// 비공개 헬퍼 메서드들

// validateCreateUserRequest - 사용자 생성 요청 검증
func (u *userUsecase) validateCreateUserRequest(req *dto.CreateUserRequest) error {
	// 비밀번호 길이 체크
	if len(req.Password) < 6 {
		return errors.ErrWeakPassword
	}

	// 여행 날짜 검증
	if err := u.validateTravelDates(req.TravelStart, req.TravelEnd); err != nil {
		return err
	}

	return nil
}

// validateTravelDates - 여행 날짜 검증
func (u *userUsecase) validateTravelDates(start, end time.Time) error {
	now := time.Now()

	// 여행 시작일이 현재보다 과거인지 체크
	if start.Before(now.Truncate(24 * time.Hour)) {
		return errors.ErrPastTravelDate
	}

	// 여행 시작일이 종료일보다 늦은지 체크
	if start.After(end) || start.Equal(end) {
		return errors.ErrInvalidTravelDates
	}

	return nil
}
