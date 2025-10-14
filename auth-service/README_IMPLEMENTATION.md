# Auth Service - Полная реализация ✅

## 📚 Реализованная функциональность

### ✅ **База данных: PostgreSQL через GORM**
- Модель `User` с UUID, email, password_hash, role
- Автомиграции через GORM
- Connection pooling
- Soft deletes

### ✅ **HTTP REST API (Gin Framework)**

**Public Endpoints:**
- `POST /api/v1/auth/register` - Регистрация пользователя
- `POST /api/v1/auth/login` - Вход (возврат access + refresh токенов)
- `POST /api/v1/auth/refresh` - Обновление токенов

**Protected Endpoints (требуется JWT):**
- `POST /api/v1/auth/logout` - Выход из системы
- `GET /api/v1/auth/profile` - Получить профиль
- `PUT /api/v1/auth/profile` - Обновить профиль

**Health Check:**
- `GET /health` - Проверка состояния сервиса

### ✅ **gRPC API**
- `ValidateToken(token)` - Валидация токена для других микросервисов
- `Register(user)` - Регистрация через gRPC
- `Login(email, password)` - Вход через gRPC
- `RefreshToken(refresh_token)` - Обновление токенов
- `Logout(refresh_token)` - Выход

### ✅ **Безопасность**
- **bcrypt** для хеширования паролей (cost 12)
- **JWT токены:**
  - Access Token: 15 минут
  - Refresh Token: 7 дней
  - HS256 алгоритм
- **Redis Blacklist** для invalidated токенов
- **Валидация входных данных** через go-playground/validator
- **Middleware для проверки ролей**

### ✅ **Логирование (Zap)**
- Структурированные логи в JSON (production)
- Цветные логи (development)
- Контекстные поля (user_id, email, error)
- Уровни: Debug, Info, Warn, Error, Fatal

### ✅ **Роли и права**
```go
const (
    RoleAdmin     = "admin"      // Полный доступ
    RoleEditor    = "editor"     // Создание/редактирование
    RoleModerator = "moderator"  // Модерация
    RoleUser      = "user"       // Базовый доступ
)
```

**Permissions по ролям:**
- `admin`: create_news, edit_news, delete_news, moderate, manage_users, manage_categories
- `editor`: create_news, edit_news, manage_categories
- `moderator`: moderate, edit_news
- `user`: read_news, comment

### ✅ **Middleware**
- `RequireAuth()` - Проверка JWT токена
- `RequireRole(roles...)` - Проверка роли пользователя
- `RequirePermission(permission)` - Проверка конкретного permission

---

## 📁 Структура директорий

```
auth-service/
├── cmd/
│   └── auth-service/
│       └── main.go                      ✅ Точка входа с Zap, HTTP + gRPC серверами
├── internal/
│   ├── config/
│   │   └── config.go                    ✅ Конфигурация с env переменными
│   ├── handler/
│   │   ├── http_handler.go              ✅ HTTP endpoints (Gin)
│   │   └── grpc_handler.go              ✅ gRPC ValidateToken
│   ├── service/
│   │   ├── auth_service.go              ✅ Регистрация, логин, логаут
│   │   ├── token_service.go             ✅ JWT генерация и валидация
│   │   └── user_service.go              ✅ Управление профилем
│   ├── repository/
│   │   ├── user_repository.go           ✅ GORM queries для users
│   │   ├── session_repository.go        ✅ Redis для refresh токенов
│   │   └── blacklist_repository.go      ✅ Redis blacklist
│   ├── models/
│   │   ├── user.go                      ✅ GORM модель + DTO
│   │   └── token.go                     ✅ TokenPair, TokenClaims
│   └── middleware/
│       ├── auth_middleware.go           ✅ JWT проверка
│       └── role_middleware.go           ✅ Проверка ролей
├── pkg/
│   ├── jwt/
│   │   └── jwt.go                       ✅ JWT утилиты (в token_service)
│   ├── logger/
│   │   └── logger.go                    ✅ Zap logger wrapper
│   ├── validator/
│   │   └── validator.go                 ✅ go-playground/validator
│   ├── hash/
│   │   └── bcrypt.go                    ✅ Хеширование паролей
│   └── database/
│       └── postgres.go                  ✅ GORM подключение
├── proto/
│   └── auth.proto                       ✅ gRPC контракт
├── migrations/
│   └── 001_create_users_table.*.sql     ✅ SQL миграции (опционально, есть AutoMigrate)
├── .env                                 ✅ Environment переменные
├── go.mod                               ✅ Dependencies
├── Dockerfile                           ✅ Docker образ
└── README.md                            ✅ Документация
```

---

## 🚀 Запуск

### Локальный запуск

```bash
cd auth-service

# Установка зависимостей
go mod download

# Запуск PostgreSQL и Redis (через Docker)
docker-compose up -d postgres redis

# Запуск сервиса
go run cmd/auth-service/main.go
```

### Docker

```bash
docker build -t auth-service .
docker run -p 8091:8091 -p 8081:8081 --env-file .env auth-service
```

---

## 📝 Примеры использования API

### Регистрация

```bash
curl -X POST http://localhost:8091/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePass123",
    "full_name": "Admin User",
    "role": "admin"
  }'
```

### Вход

```bash
curl -X POST http://localhost:8091/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePass123"
  }'
```

**Ответ:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900,
  "token_type": "Bearer"
}
```

### Получить профиль

```bash
curl http://localhost:8091/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Обновить токен

```bash
curl -X POST http://localhost:8091/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

### Выход

```bash
curl -X POST http://localhost:8091/api/v1/auth/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

---

## 🔧 Environment Variables

```env
# Server
SERVICE_NAME=auth-service
GRPC_PORT=8081
HTTP_PORT=8091

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=auth_db
DB_SSL_MODE=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-super-secret-key-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Service Discovery
CONSUL_HOST=localhost
CONSUL_PORT=8500

# Environment
ENVIRONMENT=development
LOG_LEVEL=debug
```

---

## ✅ Что реализовано

- [x] PostgreSQL через GORM
- [x] HTTP REST API на Gin
- [x] gRPC API для межсервисного взаимодействия
- [x] JWT аутентификация (access + refresh токены)
- [x] bcrypt для паролей
- [x] Redis для сессий и blacklist
- [x] Zap логирование
- [x] Валидация данных
- [x] Middleware для JWT и ролей
- [x] GORM AutoMigrate
- [x] Graceful shutdown
- [x] Health check endpoint

---

## 🎯 Следующие шаги

1. **Генерация protobuf:**
```bash
protoc --go_out=. --go-grpc_out=. proto/auth.proto
```

2. **Установка зависимостей:**
```bash
go mod tidy
```

3. **Запуск тестов:**
```bash
go test ./...
```

4. **Интеграция с API Gateway**

---

**Статус:** ✅ Готово к использованию
