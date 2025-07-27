package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/chris910512/travel-chat/internal/delivery/grpc/handler"
	"github.com/chris910512/travel-chat/internal/pkg/jwt"
	usecaseInterface "github.com/chris910512/travel-chat/internal/usecase/interface"
	pb "github.com/chris910512/travel-chat/pkg/proto/user"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	grpcServer  *grpc.Server
	gatewayMux  *runtime.ServeMux
	userHandler *handler.UserGRPCHandler
	grpcPort    string
	gatewayPort string
}

// NewGRPCServer - gRPC 서버 생성자
func NewGRPCServer(
	userUsecase usecaseInterface.UserUsecase,
	jwtService *jwt.JWTService,
	grpcPort, gatewayPort string,
) *GRPCServer {
	// gRPC 서버 생성
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	// 핸들러 생성
	userHandler := handler.NewUserGRPCHandler(userUsecase, jwtService)

	// 서비스 등록
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	// gRPC reflection 등록 (개발용)
	reflection.Register(grpcServer)

	// gRPC Gateway 설정
	gatewayMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
	)

	return &GRPCServer{
		grpcServer:  grpcServer,
		gatewayMux:  gatewayMux,
		userHandler: userHandler,
		grpcPort:    grpcPort,
		gatewayPort: gatewayPort,
	}
}

// StartGRPCServer - gRPC 서버 시작
func (s *GRPCServer) StartGRPCServer() error {
	lis, err := net.Listen("tcp", ":"+s.grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", s.grpcPort, err)
	}

	log.Printf("gRPC server starting on port %s", s.grpcPort)
	return s.grpcServer.Serve(lis)
}

// StartGatewayServer - gRPC Gateway 서버 시작
func (s *GRPCServer) StartGatewayServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// gRPC 서버에 연결
	grpcEndpoint := "localhost:" + s.grpcPort
	conn, err := grpc.DialContext(
		ctx,
		grpcEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial gRPC server: %v", err)
	}

	// gRPC Gateway에 서비스 등록
	err = pb.RegisterUserServiceHandler(ctx, s.gatewayMux, conn)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %v", err)
	}

	// CORS 설정을 위한 래퍼
	corsHandler := corsWrapper(s.gatewayMux)

	log.Printf("gRPC Gateway server starting on port %s", s.gatewayPort)
	return http.ListenAndServe(":"+s.gatewayPort, corsHandler)
}

// Stop - 서버 정지
func (s *GRPCServer) Stop() {
	log.Println("Stopping gRPC server...")
	s.grpcServer.GracefulStop()
}

// 미들웨어 및 헬퍼 함수들

// loggingInterceptor - gRPC 로깅 인터셉터
func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Printf("gRPC call: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("gRPC error: %s - %v", info.FullMethod, err)
	}

	return resp, err
}

// customHeaderMatcher - 헤더 매칭 함수
func customHeaderMatcher(key string) (string, bool) {
	switch key {
	case "Authorization":
		return key, true
	case "Content-Type":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

// corsWrapper - CORS 설정
func corsWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
