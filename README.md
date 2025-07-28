# Travel Chat - 여행자 매칭 채팅 서비스

> Go 언어 학습을 위한 백엔드 프로젝트 - gRPC, WebSocket, Go Routine을 활용한 실시간 채팅 서비스

## 📖 프로젝트 개요

여행자들이 같은 목적지로 여행을 계획 중인 사람들과 매칭되어 채팅할 수 있는 서비스입니다. 사용자는 자신의 여행 정보를 등록하고, 같은 목적지의 전체 채팅방에 참여하거나 관심 있는 사용자와 1:1 채팅을 할 수 있습니다.

### 🎯 주요 기능

- **사용자 프로필 관리**: 나이, 성별, 여행 목적지, 여행 기간, 예산, 스타일 등
- **목적지별 전체 채팅**: 같은 국가-도시로 여행하는 사용자들의 공개 채팅 (메시지 6시간 보관)
- **1:1 개인 채팅**: 매칭된 사용자 간의 개인 채팅 (메시지 24시간 보관)
- **실시간 사용자 활동 상태**: 온라인, 10분 전 활동 등

### 🏗️ 기술 스택

- **언어**: Go 1.24
- **웹 프레임워크**: Gin
- **데이터베이스**: PostgreSQL (Supabase)
- **ORM**: GORM
- **인증**: JWT
- **통신**:
    - REST API
    - gRPC + gRPC Gateway
    - WebSocket (예정)
- **동시성**: Go Routines & Channels (예정)
- **컨테이너**: Docker

## 📁 프로젝트 구조

```
travel-chat/
├── cmd/
│   ├── server/main.go              # 메인 서버 (HTTP + gRPC)
│   └── grpc-client/main.go         # gRPC 테스트 클라이언트
├── internal/
│   ├── delivery/
│   │   ├── http/                   # REST API 핸들러
│   │   │   ├── handler/
│   │   │   ├── middleware/
│   │   │   ├── response/
│   │   │   └── router/
│   │   └── grpc/                   # gRPC 핸들러
│   │       ├── handler/
│   │       └── server/
│   ├── domain/
│   │   ├── entity/                 # 도메인 엔티티
│   │   └── repository/             # 레포지토리 인터페이스
│   ├── infrastructure/
│   │   ├── database/               # DB 연결 및 마이그레이션
│   │   └── repository/             # 레포지토리 구현
│   ├── usecase/                    # 비즈니스 로직
│   └── pkg/                        # 공통 패키지
├── proto/user/user.proto           # Protocol Buffer 정의
├── pkg/proto/user/                 # 생성된 gRPC 코드
├── scripts/                        # 빌드 및 실행 스크립트
└── docs/                           # API 문서
```

## 🚀 빠른 시작

### 1. 사전 요구사항

```bash
# Go 1.24+ 설치 확인
go version

# Protocol Buffers 컴파일러 설치
# macOS
brew install protobuf

# Ubuntu/Debian
sudo apt update && sudo apt install -y protobuf-compiler

# Windows (Chocolatey)
choco install protoc

# Go gRPC 플러그인 설치
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

### 2. 프로젝트 설정

```bash
# 프로젝트 클론
git clone <repository-url>
cd travel-chat

# 의존성 설치
go mod download

# 환경변수 설정
cp .env.example .env
# .env 파일을 실제 환경에 맞게 수정
```

### 3. 환경변수 설정 (.env)

```bash
# 데이터베이스 설정 
DB_HOST=your-postgres-host
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=postgres

# 서버 포트 설정
SERVER_PORT=8080
GRPC_PORT=9090
GATEWAY_PORT=8081

# JWT 설정
JWT_SECRET_KEY=your-super-secret-jwt-key-here-make-it-long-and-secure
JWT_ISSUER=travel-chat-api
```

### 4. Protocol Buffer 컴파일

```bash
# googleapis 다운로드 (최초 1회)
git clone https://github.com/googleapis/googleapis.git third_party/googleapis

# Proto 파일 컴파일
chmod +x scripts/generate_proto.sh
./scripts/generate_proto.sh
```

### 5. 서버 실행

```bash
# 서버 시작
go run cmd/server/main.go
```

성공하면 다음과 같은 출력이 나타납니다:

```
=== Travel Chat API Server Started ===
HTTP Server: http://localhost:8080
gRPC Server: localhost:9090
gRPC Gateway: http://localhost:8081
Press Ctrl+C to exit
```

## 🧪 테스트 방법

### 1. HTTP API 테스트

#### 기본 테스트 스크립트 실행
```bash
chmod +x scripts/test_api.sh
./scripts/test_api.sh
```

#### 수동 테스트 예시
```bash
# 사용자 등록
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "테스트 사용자",
    "age": 25,
    "gender": "male",
    "country": "일본",
    "city": "도쿄",
    "travel_start": "2025-08-01T00:00:00Z",
    "travel_end": "2025-08-07T00:00:00Z",
    "bio": "안녕하세요!",
    "travel_purpose": "tourism",
    "travel_budget": 100,
    "travel_style": "planned"
  }'

# 로그인
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 2. gRPC 클라이언트 테스트

```bash
# gRPC 테스트 클라이언트 실행 (별도 터미널)
go run cmd/grpc-client/main.go
```

**예상 출력:**
```
🚀 === gRPC Client Test Started ===

📝 1. Testing User Registration...
✅ Registration successful: 사용자가 성공적으로 등록되었습니다
   User ID: 1, Name: gRPC 테스트 사용자

🔐 2. Testing User Login...
✅ Login successful: 로그인이 성공했습니다
   Access Token: eyJhbGciOiJIUzI1NiIs...
   User: gRPC 테스트 사용자 (ID: 1)

👤 3. Testing Get My Profile...
✅ Get my profile successful: 내 프로필 정보를 조회했습니다
   Profile: gRPC 테스트 사용자, Destination: 일본-도쿄

🎉 === gRPC Client Test Complete ===
```

### 3. gRPC Gateway 테스트 (HTTP로 gRPC 호출)

```bash
# gRPC Gateway를 통한 사용자 등록
curl -X POST http://localhost:8081/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email":"gateway@example.com",
    "password":"password123",
    "name":"Gateway 테스트",
    "age":26,
    "gender":"GENDER_FEMALE",
    "country":"프랑스",
    "city":"파리",
    "travel_start":"2025-09-01T00:00:00Z",
    "travel_end":"2025-09-07T00:00:00Z",
    "bio":"Gateway를 통한 등록",
    "travel_purpose":"TRAVEL_PURPOSE_CULTURE",
    "travel_budget":200,
    "travel_style":"TRAVEL_STYLE_LUXURY"
  }'
```

### 4. grpcurl을 이용한 직접 gRPC 테스트 (선택사항)

```bash
# grpcurl 설치
brew install grpcurl  # macOS

# 서비스 목록 확인
grpcurl -plaintext localhost:9090 list

# 사용자 등록
grpcurl -plaintext -d '{
  "email":"grpcurl@example.com",
  "password":"password123",
  "name":"grpcurl 테스트",
  "age":25,
  "gender":"GENDER_MALE",
  "country":"한국",
  "city":"서울",
  "travel_start":"2025-08-01T00:00:00Z",
  "travel_end":"2025-08-07T00:00:00Z",
  "bio":"grpcurl로 등록",
  "travel_purpose":"TRAVEL_PURPOSE_TOURISM",
  "travel_budget":100,
  "travel_style":"TRAVEL_STYLE_PLANNED"
}' localhost:9090 user.UserService/Register
```

## 📚 API 문서

### REST API 엔드포인트

#### 인증 (Authentication)
- `POST /api/auth/register` - 사용자 등록
- `POST /api/auth/login` - 로그인
- `POST /api/auth/refresh` - 토큰 갱신

#### 사용자 관리 (Users)
- `GET /api/users` - 사용자 목록 조회 (페이징)
- `GET /api/users/:id` - 특정 사용자 조회
- `GET /api/users/me` - 내 프로필 조회 (인증 필요)
- `PUT /api/users/:id` - 프로필 업데이트 (인증 필요)
- `DELETE /api/users/:id` - 사용자 삭제 (인증 필요)
- `GET /api/users/destination/:country/:city` - 목적지별 사용자 조회

#### 유틸리티
- `GET /api/health` - 서버 상태 확인

### gRPC 서비스

동일한 기능을 gRPC로도 제공하며, gRPC Gateway를 통해 HTTP로도 접근 가능합니다.

- **gRPC 엔드포인트**: `localhost:9090`
- **gRPC Gateway 엔드포인트**: `http://localhost:8081/v1/`

## 🔧 개발 진행 상황

### ✅ 완료된 기능 (1-3주차)

- [x] **1-2주차**: 프로젝트 구조 설계 및 기본 CRUD
    - Clean Architecture 구조
    - 사용자 등록/로그인/프로필 관리
    - JWT 인증 시스템
    - PostgreSQL 연동 (GORM)
    - REST API 구현

- [x] **3주차**: gRPC 서비스 구현
    - Protocol Buffer 정의
    - gRPC 서버 및 핸들러
    - gRPC Gateway (REST ↔ gRPC 브릿지)
    - 통합 서버 (HTTP + gRPC 동시 실행)

### 🚧 예정된 기능 (4-8주차)

- [ ] **4주차**: Go Routine과 채널 활용
    - 동시성 프로그래밍 기초
    - 메시지 처리를 위한 worker pool
    - 비동기 작업 처리

- [ ] **5주차**: WebSocket 실시간 채팅
    - 실시간 메시지 송수신
    - 채팅방 관리
    - 연결 상태 관리

- [ ] **6주차**: 고급 채팅 기능
    - 자동 채팅방 생성 (국가-도시 기반)
    - 메시지 TTL (전체: 6시간, 1:1: 24시간)
    - 1:1 매칭 시스템

- [ ] **7주차**: 성능 최적화 및 테스트
    - 캐싱 시스템
    - 부하 테스트
    - 단위 테스트

- [ ] **8주차**: 배포 및 모니터링
    - Docker 컨테이너화
    - GCP/AWS 배포
    - 로깅 및 모니터링

## 📋 문제 해결

### 자주 발생하는 오류

#### 1. Protocol Buffer 컴파일 오류
```bash
# protoc 설치 확인
which protoc

# Go 플러그인 확인
ls $GOPATH/bin/ | grep protoc

# 재설치
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

#### 2. 포트 충돌
```bash
# 포트 사용 확인
lsof -i :8080
lsof -i :9090
lsof -i :8081

# 프로세스 종료
kill -9 <PID>
```

#### 3. 데이터베이스 연결 오류
- `.env` 파일의 데이터베이스 설정 확인
- Supabase 프로젝트가 활성화되어 있는지 확인
- 네트워크 연결 확인

#### 4. gRPC 연결 오류
- 서버가 먼저 실행되었는지 확인
- 방화벽에서 9090 포트가 열려있는지 확인
- proto 파일이 올바르게 컴파일되었는지 확인

## 📄 라이센스

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 연락처

프로젝트 관련 질문이나 제안사항이 있으시면 Issue를 생성해 주세요.
