# Сборка
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

    # Копируем исходники
    COPY . .

    # Обновляем зависимости
    RUN go mod tidy

    # Собираем бинарник
    RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seo-service ./cmd/seo-service# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/seo-service .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations

# Устанавливаем таймзону для России
ENV TZ=Europe/Moscow

EXPOSE 8093

CMD ["./seo-service"]
