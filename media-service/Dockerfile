# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
RUN go mod download || true

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o media-service ./cmd/media-service

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/media-service .

# Create uploads directory
RUN mkdir -p /root/uploads

# Expose port
EXPOSE 8085

# Run the binary
CMD ["./media-service"]
