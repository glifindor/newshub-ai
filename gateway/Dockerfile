# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

    # Copy go mod files
    COPY go.mod ./
    RUN go mod tidy
    RUN go mod download || true# Copy source code
COPY . .

# Update dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd/gateway

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/gateway .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./gateway"]
