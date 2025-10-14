# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./
RUN go mod download || true

# Copy source code
COPY . .

# Update dependencies
RUN go mod tidy

    # Build the application (HTTP only, no gRPC)
    RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service ./cmd/auth-service/main-http-only.go# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/auth-service .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Expose gRPC port
EXPOSE 8081

CMD ["./auth-service"]
