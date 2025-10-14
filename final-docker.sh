#!/bin/bash
set -e

echo "Creating FINAL Dockerfiles with go mod tidy..."

# Auth Service
cat > /opt/news-portal/auth-service/Dockerfile << 'EOF'
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service ./cmd/auth-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/auth-service .
EXPOSE 8091 8081
CMD ["./auth-service"]
EOF

echo "✅ Auth Service Dockerfile"

# News Service
cat > /opt/news-portal/news-service/Dockerfile << 'EOF'
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/news-service ./cmd/news-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/news-service .
EXPOSE 8092 8082
CMD ["./news-service"]
EOF

echo "✅ News Service Dockerfile"

# Media Service
cat > /opt/news-portal/media-service/Dockerfile << 'EOF'
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/media-service ./cmd/media-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/media-service .
EXPOSE 8094
CMD ["./media-service"]
EOF

echo "✅ Media Service Dockerfile"

echo ""
echo "==================================================" 
echo "  ALL DOCKERFILES UPDATED WITH GO MOD TIDY!"
echo "=================================================="
