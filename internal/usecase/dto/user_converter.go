package dto

import (
	"github.com/chris910512/travel-chat/internal/domain/entity/user"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserRequest를 User 엔티티로 변환
func (req *CreateUserRequest) ToEntity() (*user.User, error) {
	// 비밀번호 해싱
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &user.User{
		Email:         req.Email,
		Password:      string(hashedPassword),
		Name:          req.Name,
		Age:           req.Age,
		Gender:        user.GenderFromString(req.Gender),
		Country:       req.Country,
		City:          req.City,
		TravelStart:   req.TravelStart,
		TravelEnd:     req.TravelEnd,
		Bio:           req.Bio,
		TravelPurpose: user.TravelPurposeFromString(req.TravelPurpose),
		TravelBudget:  req.TravelBudget,
		TravelStyle:   user.TravelStyleFromString(req.TravelStyle),
	}, nil
}

// User 엔티티를 UserResponse로 변환
func FromUserEntity(u *user.User) *UserResponse {
	return &UserResponse{
		ID:             u.ID,
		Email:          u.Email,
		Name:           u.Name,
		Age:            u.Age,
		Gender:         (&u.Gender).String(),
		ProfilePic:     u.ProfilePic,
		Country:        u.Country,
		City:           u.City,
		Destination:    u.GetDestination(),
		TravelStart:    u.TravelStart,
		TravelEnd:      u.TravelEnd,
		Bio:            u.Bio,
		TravelPurpose:  (&u.TravelPurpose).String(),
		TravelBudget:   u.TravelBudget,
		TravelStyle:    (&u.TravelStyle).String(),
		ActivityStatus: u.GetActivityStatus(),
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

// User 엔티티 슬라이스를 UserResponse 슬라이스로 변환
func FromUserEntities(users []*user.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = *FromUserEntity(u)
	}
	return responses
}

// UpdateUserRequest를 기존 User 엔티티에 적용
func (req *UpdateUserRequest) ApplyToEntity(u *user.User) {
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Age != nil {
		u.Age = *req.Age
	}
	if req.Gender != nil {
		u.Gender = user.GenderFromString(*req.Gender)
	}
	if req.ProfilePic != nil {
		u.ProfilePic = *req.ProfilePic
	}
	if req.Country != nil {
		u.Country = *req.Country
	}
	if req.City != nil {
		u.City = *req.City
	}
	if req.TravelStart != nil {
		u.TravelStart = *req.TravelStart
	}
	if req.TravelEnd != nil {
		u.TravelEnd = *req.TravelEnd
	}
	if req.Bio != nil {
		u.Bio = *req.Bio
	}
	if req.TravelPurpose != nil {
		u.TravelPurpose = user.TravelPurposeFromString(*req.TravelPurpose)
	}
	if req.TravelBudget != nil {
		u.TravelBudget = *req.TravelBudget
	}
	if req.TravelStyle != nil {
		u.TravelStyle = user.TravelStyleFromString(*req.TravelStyle)
	}
}

// 페이징 계산 헬퍼
func (req *GetUsersRequest) GetOffset() int {
	return (req.Page - 1) * req.Limit
}

// 전체 페이지 수 계산
func CalculateTotalPages(totalCount int64, limit int) int {
	if totalCount == 0 {
		return 0
	}
	return int((totalCount + int64(limit) - 1) / int64(limit))
}
