#!/bin/bash

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=================================${NC}"
echo -e "${BLUE}News Portal - Service Health Check${NC}"
echo -e "${BLUE}=================================${NC}\n"

# Функция проверки сервиса
check_service() {
    local name=$1
    local url=$2
    
    if curl -s -o /dev/null -w "%{http_code}" "$url" | grep -q "200\|201"; then
        echo -e "${GREEN}✓${NC} $name is healthy"
    else
        echo -e "${RED}✗${NC} $name is down"
    fi
}

# Проверка сервисов
check_service "API Gateway     " "http://localhost:8080/health"
check_service "Frontend        " "http://localhost:3000"
check_service "Consul          " "http://localhost:8500/v1/status/leader"
check_service "RabbitMQ        " "http://localhost:15672"
check_service "MinIO           " "http://localhost:9000/minio/health/live"
check_service "Prometheus      " "http://localhost:9090/-/healthy"
check_service "Grafana         " "http://localhost:3001/api/health"

echo ""
echo -e "${BLUE}=================================${NC}"
echo -e "${BLUE}Database Connections${NC}"
echo -e "${BLUE}=================================${NC}\n"

# Проверка PostgreSQL
if docker exec news-postgres pg_isready -U postgres > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} PostgreSQL is ready"
else
    echo -e "${RED}✗${NC} PostgreSQL is down"
fi

# Проверка Redis
if docker exec news-redis redis-cli ping | grep -q "PONG"; then
    echo -e "${GREEN}✓${NC} Redis is ready"
else
    echo -e "${RED}✗${NC} Redis is down"
fi

echo ""
echo -e "${BLUE}=================================${NC}"
echo -e "${BLUE}Useful URLs${NC}"
echo -e "${BLUE}=================================${NC}"
echo -e "Frontend:        http://localhost:3000"
echo -e "API Gateway:     http://localhost:8080"
echo -e "Consul UI:       http://localhost:8500"
echo -e "RabbitMQ:        http://localhost:15672"
echo -e "MinIO Console:   http://localhost:9001"
echo -e "Grafana:         http://localhost:3001"
echo -e "Prometheus:      http://localhost:9090"
echo ""
