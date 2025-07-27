#!/bin/bash

# proto 컴파일 스크립트
# 파일명: scripts/generate_proto.sh

# googleapis 다운로드 (gRPC Gateway 사용을 위해 필요)
if [ ! -d "third_party/googleapis" ]; then
    echo "Downloading googleapis..."
    mkdir -p third_party
    git clone https://github.com/googleapis/googleapis.git third_party/googleapis
fi

# Protocol Buffer 파일 컴파일
echo "Generating Protocol Buffer files..."

protoc \
    --proto_path=. \
    --proto_path=third_party/googleapis \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=. \
    --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=./docs \
    pkg/proto/user/user.proto

echo "Protocol Buffer generation complete!"

# 생성된 파일들
echo "Generated files:"
echo "- pkg/proto/user/user.pb.go"
echo "- pkg/proto/user/user_grpc.pb.go"
echo "- pkg/proto/user/user.pb.gw.go"
echo "- docs/user/user.swagger.json"

