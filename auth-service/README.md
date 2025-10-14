# Auth Service

Микросервис аутентификации и авторизации для новостного портала.

## Функциональность

- ✅ Регистрация пользователей
- ✅ Авторизация (JWT токены)
- ✅ Управление ролями и правами
- ✅ Refresh токены
- ✅ Валидация токенов
- ✅ Logout

## Технологии

- **Golang** - основной язык
- **gRPC** - межсервисное взаимодействие
- **PostgreSQL** - хранение пользователей
- **Redis** - хранение сессий и refresh токенов
- **JWT** - токены доступа

## Установка и запуск

### Требования

- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### Установка зависимостей

```bash
go mod download
```

### Генерация protobuf

```bash
protoc --go_out=. --go-grpc_out=. proto/auth.proto
```

### Миграции базы данных

```bash
# Используйте golang-migrate или выполните вручную
psql -U postgres -d auth_db -f migrations/001_create_users_table.up.sql
```

### Запуск

```bash
# Настройте .env файл
cp .env.example .env

# Запустите сервис
go run cmd/auth-service/main.go
```

### Docker

```bash
# Сборка образа
docker build -t auth-service .

# Запуск контейнера
docker run -p 8081:8081 --env-file .env auth-service
```

## API (gRPC)

### Register
```protobuf
rpc Register(RegisterRequest) returns (RegisterResponse)
```

### Login
```protobuf
rpc Login(LoginRequest) returns (LoginResponse)
```

### ValidateToken
```protobuf
rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse)
```

### RefreshToken
```protobuf
rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse)
```

### Logout
```protobuf
rpc Logout(LogoutRequest) returns (LogoutResponse)
```

## Роли

- **admin** - полный доступ ко всем функциям
- **editor** - создание и редактирование новостей
- **moderator** - модерация контента
- **user** - чтение и комментирование

## Разработка

### Структура проекта

```
/auth-service
├── cmd/auth-service/main.go      # Точка входа
├── internal/
│   ├── config/                   # Конфигурация
│   ├── handler/                  # gRPC обработчики
│   ├── service/                  # Бизнес-логика
│   ├── repository/               # Работа с БД
│   └── models/                   # Модели данных
├── pkg/                          # Переиспользуемые пакеты
├── proto/                        # Protobuf определения
├── migrations/                   # SQL миграции
└── Dockerfile
```

## Тестирование

```bash
go test ./...
```

## Лицензия

MIT
