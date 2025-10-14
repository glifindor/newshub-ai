#!/bin/bash

# Скрипт для генерации всех protobuf файлов

set -e

BLUE='\033[0;34m'
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${BLUE}Generating protobuf files...${NC}\n"

# Функция генерации proto
generate_proto() {
    local service=$1
    local proto_file=$2
    
    echo -e "${BLUE}Generating $service...${NC}"
    cd "$service"
    protoc --go_out=. --go-grpc_out=. "$proto_file"
    cd ..
    echo -e "${GREEN}✓ $service done${NC}\n"
}

# Генерация для каждого сервиса
generate_proto "auth-service" "proto/auth.proto"
generate_proto "news-service" "proto/news.proto"
generate_proto "seo-service" "proto/seo.proto"
generate_proto "admin-service" "proto/admin.proto"
generate_proto "media-service" "proto/media.proto"

echo -e "${GREEN}All protobuf files generated successfully!${NC}"
