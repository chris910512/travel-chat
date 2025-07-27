package handler

import (
	"context"

	"github.com/chris910512/travel-chat/internal/pkg/jwt"
	"github.com/chris910512/travel-chat/internal/usecase/dto"
	usecaseInterface "github.com/chris910512/travel-chat/internal/usecase/interface"
	pb "github.com/chris910512/travel-chat/pkg/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

type UserGRPCHandler struct {
	pb.UnimplementedUserServiceServer
	userUsecase usecaseInterface.UserUsecase
	jwtService  *jwt.JWTService
}

func NewUserGRPCHandler(userUsecase usecaseInterface.UserUsecase, jwtService *jwt.JWTService) *UserGRPCHandler {
	return &UserGRPCHandler{
		userUsecase: userUsecase,
		jwtService:  jwtService,
	}
}

// Register - 사용자 등록
func (h *UserGRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Proto 메시지를 DTO로 변환
	createReq := &dto.CreateUserRequest{
		Email:         req.Email,
		Password:      req.Password,
		Name:          req.Name,
		Age:           int(req.Age),
		Gender:        protoGenderToString(req.Gender),
		Country:       req.Country,
		City:          req.City,
		TravelStart:   req.TravelStart.AsTime(),
		TravelEnd:     req.TravelEnd.AsTime(),
		Bio:           req.Bio,
		TravelPurpose: protoTravelPurposeToString(req.TravelPurpose),
		TravelBudget:  int(req.TravelBudget),
		TravelStyle:   protoTravelStyleToString(req.TravelStyle),
	}

	// Usecase 호출
	userResp, err := h.userUsecase.Register(ctx, createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "사용자 등록 실패: %v", err)
	}

	// DTO를 Proto 메시지로 변환
	protoUser := userDtoToProto(userResp)

	return &pb.RegisterResponse{
		User:    protoUser,
		Message: "사용자가 성공적으로 등록되었습니다",
	}, nil
}

// Login - 사용자 로그인
func (h *UserGRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginReq := &dto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	loginResp, err := h.userUsecase.Login(ctx, loginReq)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "로그인 실패: %v", err)
	}

	protoUser := userDtoToProto(&loginResp.User)

	return &pb.LoginResponse{
		User:         protoUser,
		AccessToken:  loginResp.AccessToken,
		RefreshToken: loginResp.RefreshToken,
		TokenType:    loginResp.TokenType,
		ExpiresIn:    uint32(loginResp.ExpiresIn),
		Message:      "로그인이 성공했습니다",
	}, nil
}

// RefreshToken - 토큰 갱신
func (h *UserGRPCHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	refreshReq := &dto.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	refreshResp, err := h.userUsecase.RefreshToken(ctx, refreshReq)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "토큰 갱신 실패: %v", err)
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  refreshResp.AccessToken,
		RefreshToken: refreshResp.RefreshToken,
		TokenType:    refreshResp.TokenType,
		ExpiresIn:    uint32(refreshResp.ExpiresIn),
		Message:      "토큰이 갱신되었습니다",
	}, nil
}

// GetProfile - 프로필 조회
func (h *UserGRPCHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	userResp, err := h.userUsecase.GetByID(ctx, uint(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "사용자를 찾을 수 없습니다: %v", err)
	}

	protoUser := userDtoToProto(userResp)

	return &pb.GetProfileResponse{
		User:    protoUser,
		Message: "사용자 정보를 조회했습니다",
	}, nil
}

// GetMyProfile - 내 프로필 조회
func (h *UserGRPCHandler) GetMyProfile(ctx context.Context, req *pb.GetMyProfileRequest) (*pb.GetProfileResponse, error) {
	// gRPC 메타데이터에서 JWT 토큰 추출
	userID, err := h.extractUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "인증이 필요합니다: %v", err)
	}

	userResp, err := h.userUsecase.GetByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "사용자를 찾을 수 없습니다: %v", err)
	}

	protoUser := userDtoToProto(userResp)

	return &pb.GetProfileResponse{
		User:    protoUser,
		Message: "내 프로필 정보를 조회했습니다",
	}, nil
}

// GetUsers - 사용자 목록 조회
func (h *UserGRPCHandler) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	getUsersReq := &dto.GetUsersRequest{
		Page:    int(req.Page),
		Limit:   int(req.Limit),
		Country: req.Country,
		City:    req.City,
	}

	usersResp, err := h.userUsecase.GetUsers(ctx, getUsersReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "사용자 목록 조회 실패: %v", err)
	}

	// DTO 슬라이스를 Proto 슬라이스로 변환
	protoUsers := make([]*pb.User, len(usersResp.Users))
	for i, user := range usersResp.Users {
		protoUsers[i] = userDtoToProto(&user)
	}

	return &pb.GetUsersResponse{
		Users:      protoUsers,
		Page:       uint32(usersResp.Page),
		Limit:      uint32(usersResp.Limit),
		TotalCount: uint64(usersResp.TotalCount),
		TotalPages: uint32(usersResp.TotalPages),
		Message:    "사용자 목록을 조회했습니다",
	}, nil
}

// GetUsersByDestination - 목적지별 사용자 조회
func (h *UserGRPCHandler) GetUsersByDestination(ctx context.Context, req *pb.GetUsersByDestinationRequest) (*pb.GetUsersResponse, error) {
	users, err := h.userUsecase.GetUsersByDestination(ctx, req.Country, req.City)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "목적지별 사용자 조회 실패: %v", err)
	}

	protoUsers := make([]*pb.User, len(users))
	for i, user := range users {
		protoUsers[i] = userDtoToProto(&user)
	}

	return &pb.GetUsersResponse{
		Users:   protoUsers,
		Message: "목적지별 사용자 목록을 조회했습니다",
	}, nil
}

// UpdateProfile - 프로필 업데이트
func (h *UserGRPCHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	updateReq := &dto.UpdateUserRequest{}

	// Optional 필드들 처리
	if req.Name != nil {
		updateReq.Name = req.Name
	}
	if req.Age != nil {
		age := int(*req.Age)
		updateReq.Age = &age
	}
	if req.Gender != nil {
		gender := protoGenderToString(*req.Gender)
		updateReq.Gender = &gender
	}
	if req.ProfilePic != nil {
		updateReq.ProfilePic = req.ProfilePic
	}
	if req.Country != nil {
		updateReq.Country = req.Country
	}
	if req.City != nil {
		updateReq.City = req.City
	}
	if req.TravelStart != nil {
		travelStart := req.TravelStart.AsTime()
		updateReq.TravelStart = &travelStart
	}
	if req.TravelEnd != nil {
		travelEnd := req.TravelEnd.AsTime()
		updateReq.TravelEnd = &travelEnd
	}
	if req.Bio != nil {
		updateReq.Bio = req.Bio
	}
	if req.TravelPurpose != nil {
		travelPurpose := protoTravelPurposeToString(*req.TravelPurpose)
		updateReq.TravelPurpose = &travelPurpose
	}
	if req.TravelBudget != nil {
		travelBudget := int(*req.TravelBudget)
		updateReq.TravelBudget = &travelBudget
	}
	if req.TravelStyle != nil {
		travelStyle := protoTravelStyleToString(*req.TravelStyle)
		updateReq.TravelStyle = &travelStyle
	}

	userResp, err := h.userUsecase.UpdateProfile(ctx, uint(req.UserId), updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "프로필 업데이트 실패: %v", err)
	}

	protoUser := userDtoToProto(userResp)

	return &pb.UpdateProfileResponse{
		User:    protoUser,
		Message: "프로필이 업데이트되었습니다",
	}, nil
}

// UpdateLastActive - 활동 시간 업데이트
func (h *UserGRPCHandler) UpdateLastActive(ctx context.Context, req *pb.UpdateLastActiveRequest) (*pb.UpdateLastActiveResponse, error) {
	err := h.userUsecase.UpdateLastActive(ctx, uint(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "활동 시간 업데이트 실패: %v", err)
	}

	return &pb.UpdateLastActiveResponse{
		Message: "활동 시간이 업데이트되었습니다",
	}, nil
}

// Helper 함수들

// JWT 토큰에서 사용자 ID 추출
func (h *UserGRPCHandler) extractUserIDFromContext(ctx context.Context) (uint, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "메타데이터가 없습니다")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return 0, status.Error(codes.Unauthenticated, "Authorization 헤더가 없습니다")
	}

	authHeader := authHeaders[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, status.Error(codes.Unauthenticated, "Bearer 토큰이 아닙니다")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		return 0, status.Errorf(codes.Unauthenticated, "토큰 검증 실패: %v", err)
	}

	return claims.UserID, nil
}

// DTO를 Proto 메시지로 변환
func userDtoToProto(userDto *dto.UserResponse) *pb.User {
	return &pb.User{
		Id:             uint32(userDto.ID),
		Email:          userDto.Email,
		Name:           userDto.Name,
		Age:            uint32(userDto.Age),
		Gender:         stringToProtoGender(userDto.Gender),
		ProfilePic:     userDto.ProfilePic,
		Country:        userDto.Country,
		City:           userDto.City,
		Destination:    userDto.Destination,
		TravelStart:    timestamppb.New(userDto.TravelStart),
		TravelEnd:      timestamppb.New(userDto.TravelEnd),
		Bio:            userDto.Bio,
		TravelPurpose:  stringToProtoTravelPurpose(userDto.TravelPurpose),
		TravelBudget:   uint32(userDto.TravelBudget),
		TravelStyle:    stringToProtoTravelStyle(userDto.TravelStyle),
		ActivityStatus: userDto.ActivityStatus,
		CreatedAt:      timestamppb.New(userDto.CreatedAt),
		UpdatedAt:      timestamppb.New(userDto.UpdatedAt),
	}
}

// Enum 변환 함수들
func protoGenderToString(gender pb.Gender) string {
	switch gender {
	case pb.Gender_GENDER_MALE:
		return "male"
	case pb.Gender_GENDER_FEMALE:
		return "female"
	case pb.Gender_GENDER_OTHER:
		return "other"
	default:
		return "male"
	}
}

func stringToProtoGender(gender string) pb.Gender {
	switch gender {
	case "male":
		return pb.Gender_GENDER_MALE
	case "female":
		return pb.Gender_GENDER_FEMALE
	case "other":
		return pb.Gender_GENDER_OTHER
	default:
		return pb.Gender_GENDER_MALE
	}
}

func protoTravelPurposeToString(purpose pb.TravelPurpose) string {
	switch purpose {
	case pb.TravelPurpose_TRAVEL_PURPOSE_TOURISM:
		return "tourism"
	case pb.TravelPurpose_TRAVEL_PURPOSE_BUSINESS:
		return "business"
	case pb.TravelPurpose_TRAVEL_PURPOSE_BACKPACKING:
		return "backpacking"
	case pb.TravelPurpose_TRAVEL_PURPOSE_FOOD_TOUR:
		return "food_tour"
	case pb.TravelPurpose_TRAVEL_PURPOSE_CULTURE:
		return "culture"
	case pb.TravelPurpose_TRAVEL_PURPOSE_ACTIVITY:
		return "activity"
	case pb.TravelPurpose_TRAVEL_PURPOSE_RELAXATION:
		return "relaxation"
	default:
		return "tourism"
	}
}

func stringToProtoTravelPurpose(purpose string) pb.TravelPurpose {
	switch purpose {
	case "tourism":
		return pb.TravelPurpose_TRAVEL_PURPOSE_TOURISM
	case "business":
		return pb.TravelPurpose_TRAVEL_PURPOSE_BUSINESS
	case "backpacking":
		return pb.TravelPurpose_TRAVEL_PURPOSE_BACKPACKING
	case "food_tour":
		return pb.TravelPurpose_TRAVEL_PURPOSE_FOOD_TOUR
	case "culture":
		return pb.TravelPurpose_TRAVEL_PURPOSE_CULTURE
	case "activity":
		return pb.TravelPurpose_TRAVEL_PURPOSE_ACTIVITY
	case "relaxation":
		return pb.TravelPurpose_TRAVEL_PURPOSE_RELAXATION
	default:
		return pb.TravelPurpose_TRAVEL_PURPOSE_TOURISM
	}
}

func protoTravelStyleToString(style pb.TravelStyle) string {
	switch style {
	case pb.TravelStyle_TRAVEL_STYLE_PLANNED:
		return "planned"
	case pb.TravelStyle_TRAVEL_STYLE_SPONTANEOUS:
		return "spontaneous"
	case pb.TravelStyle_TRAVEL_STYLE_LUXURY:
		return "luxury"
	case pb.TravelStyle_TRAVEL_STYLE_BUDGET:
		return "budget"
	case pb.TravelStyle_TRAVEL_STYLE_ADVENTURE:
		return "adventure"
	case pb.TravelStyle_TRAVEL_STYLE_LEISURELY:
		return "leisurely"
	default:
		return "planned"
	}
}

func stringToProtoTravelStyle(style string) pb.TravelStyle {
	switch style {
	case "planned":
		return pb.TravelStyle_TRAVEL_STYLE_PLANNED
	case "spontaneous":
		return pb.TravelStyle_TRAVEL_STYLE_SPONTANEOUS
	case "luxury":
		return pb.TravelStyle_TRAVEL_STYLE_LUXURY
	case "budget":
		return pb.TravelStyle_TRAVEL_STYLE_BUDGET
	case "adventure":
		return pb.TravelStyle_TRAVEL_STYLE_ADVENTURE
	case "leisurely":
		return pb.TravelStyle_TRAVEL_STYLE_LEISURELY
	default:
		return pb.TravelStyle_TRAVEL_STYLE_PLANNED
	}
}
