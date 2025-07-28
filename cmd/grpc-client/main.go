package main

import (
	"context"
	"log"
	"time"

	pb "github.com/chris910512/travel-chat/pkg/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// gRPC 서버에 연결
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// 클라이언트 생성
	client := pb.NewUserServiceClient(conn)
	ctx := context.Background()

	log.Println("🚀 === gRPC Client Test Started ===")

	// 1. 사용자 등록 테스트
	log.Println("\n📝 1. Testing User Registration...")

	travelStart := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	travelEnd := time.Date(2025, 8, 7, 0, 0, 0, 0, time.UTC)

	registerReq := &pb.RegisterRequest{
		Email:         "grpc@example.com",
		Password:      "password123",
		Name:          "gRPC 테스트 사용자",
		Age:           28,
		Gender:        pb.Gender_GENDER_MALE,
		Country:       "일본",
		City:          "도쿄",
		TravelStart:   timestamppb.New(travelStart),
		TravelEnd:     timestamppb.New(travelEnd),
		Bio:           "gRPC를 통해 등록한 사용자입니다",
		TravelPurpose: pb.TravelPurpose_TRAVEL_PURPOSE_TOURISM,
		TravelBudget:  150,
		TravelStyle:   pb.TravelStyle_TRAVEL_STYLE_PLANNED,
	}

	registerResp, err := client.Register(ctx, registerReq)
	if err != nil {
		log.Printf("❌ Registration failed: %v", err)
	} else {
		log.Printf("✅ Registration successful: %s", registerResp.Message)
		log.Printf("   User ID: %d, Name: %s", registerResp.User.Id, registerResp.User.Name)
	}

	// 2. 로그인 테스트
	log.Println("\n🔐 2. Testing User Login...")

	loginReq := &pb.LoginRequest{
		Email:    "grpc@example.com",
		Password: "password123",
	}

	loginResp, err := client.Login(ctx, loginReq)
	if err != nil {
		log.Printf("❌ Login failed: %v", err)
		return
	}

	log.Printf("✅ Login successful: %s", loginResp.Message)
	log.Printf("   Access Token: %s...", loginResp.AccessToken[:30])
	log.Printf("   User: %s (ID: %d)", loginResp.User.Name, loginResp.User.Id)

	// JWT 토큰을 메타데이터에 추가
	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+loginResp.AccessToken)

	// 3. 내 프로필 조회 테스트
	log.Println("\n👤 3. Testing Get My Profile...")

	myProfileReq := &pb.GetMyProfileRequest{}
	myProfileResp, err := client.GetMyProfile(authCtx, myProfileReq)
	if err != nil {
		log.Printf("❌ Get my profile failed: %v", err)
	} else {
		log.Printf("✅ Get my profile successful: %s", myProfileResp.Message)
		log.Printf("   Profile: %s, Destination: %s", myProfileResp.User.Name, myProfileResp.User.Destination)
		log.Printf("   Travel Budget: %d만원", myProfileResp.User.TravelBudget)
	}

	// 4. 사용자 목록 조회 테스트
	log.Println("\n📋 4. Testing Get Users...")

	getUsersReq := &pb.GetUsersRequest{
		Page:  1,
		Limit: 10,
	}

	getUsersResp, err := client.GetUsers(ctx, getUsersReq)
	if err != nil {
		log.Printf("❌ Get users failed: %v", err)
	} else {
		log.Printf("✅ Get users successful: %s", getUsersResp.Message)
		log.Printf("   Total users: %d, Page: %d/%d", getUsersResp.TotalCount, getUsersResp.Page, getUsersResp.TotalPages)

		for i, user := range getUsersResp.Users {
			log.Printf("   User %d: %s (%s)", i+1, user.Name, user.Destination)
		}
	}

	// 5. 목적지별 사용자 조회 테스트
	log.Println("\n🌏 5. Testing Get Users by Destination...")

	getByDestReq := &pb.GetUsersByDestinationRequest{
		Country: "일본",
		City:    "도쿄",
	}

	getByDestResp, err := client.GetUsersByDestination(ctx, getByDestReq)
	if err != nil {
		log.Printf("❌ Get users by destination failed: %v", err)
	} else {
		log.Printf("✅ Get users by destination successful: %s", getByDestResp.Message)
		log.Printf("   Users going to 일본-도쿄: %d", len(getByDestResp.Users))

		for i, user := range getByDestResp.Users {
			log.Printf("   User %d: %s (Budget: %d만원)", i+1, user.Name, user.TravelBudget)
		}
	}

	// 6. 프로필 업데이트 테스트
	log.Println("\n✏️ 6. Testing Update Profile...")

	newBio := "gRPC 클라이언트로 업데이트된 프로필입니다"
	newBudget := uint32(200)

	updateReq := &pb.UpdateProfileRequest{
		UserId:       loginResp.User.Id,
		Bio:          &newBio,
		TravelBudget: &newBudget,
	}

	updateResp, err := client.UpdateProfile(authCtx, updateReq)
	if err != nil {
		log.Printf("❌ Update profile failed: %v", err)
	} else {
		log.Printf("✅ Update profile successful: %s", updateResp.Message)
		log.Printf("   Updated Bio: %s", updateResp.User.Bio)
		log.Printf("   Updated Budget: %d만원", updateResp.User.TravelBudget)
	}

	// 7. 토큰 갱신 테스트
	log.Println("\n🔄 7. Testing Token Refresh...")

	refreshReq := &pb.RefreshTokenRequest{
		RefreshToken: loginResp.RefreshToken,
	}

	refreshResp, err := client.RefreshToken(ctx, refreshReq)
	if err != nil {
		log.Printf("❌ Token refresh failed: %v", err)
	} else {
		log.Printf("✅ Token refresh successful: %s", refreshResp.Message)
		log.Printf("   New Access Token: %s...", refreshResp.AccessToken[:30])
	}

	log.Println("\n🎉 === gRPC Client Test Complete ===")
}
