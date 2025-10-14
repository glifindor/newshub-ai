# Stage 1: Build backend
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Update dependencies
RUN go mod tidy

# Build backend
RUN CGO_ENABLED=0 GOOS=linux go build -o admin-service ./cmd/admin-service

# Stage 2: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY frontend/package*.json ./

# Install dependencies (including devDependencies for build)
RUN npm ci

# Copy frontend source
COPY frontend/ ./

# Build frontend
RUN npm run build

# Stage 3: Final image
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy backend binary
COPY --from=backend-builder /app/admin-service .

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Copy .env template (override with docker-compose)
COPY .env.example .env

# Expose port
EXPOSE 8084

# Run
CMD ["./admin-service"]
