# News Portal - Микросервисная Архитектура

Новостной портал на Golang с микросервисной архитектурой, Next.js frontend и gRPC взаимодействием.

## 🏗️ Архитектура

Проект состоит из следующих микросервисов:

- **auth-service** (`:8081`) - Аутентификация и авторизация
- **news-service** (`:8082`) - Управление новостями
- **seo-service** (`:8083`) - SEO метаданные и sitemap
- **admin-service** (`:8084`) - Админ-панель
- **media-service** (`:8085`) - Загрузка и хранение медиа
- **gateway** (`:8080`) - API Gateway
- **frontend** (`:3000`) - Next.js приложение

## 🚀 Быстрый старт

### Требования

- Docker & Docker Compose
- Go 1.21+ (для локальной разработки)
- Node.js 18+ (для frontend разработки)

### Запуск всех сервисов

```bash
# Клонируйте репозиторий
cd "НОВСТНОЙ САЙТ"

# Запустите все сервисы через Docker Compose
docker-compose up -d

# Проверьте статус сервисов
docker-compose ps

# Просмотр логов
docker-compose logs -f gateway
```

### Доступ к сервисам

- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Consul UI**: http://localhost:8500
- **RabbitMQ Management**: http://localhost:15672 (admin/password)
- **MinIO Console**: http://localhost:9001 (admin/password123)
- **Grafana**: http://localhost:3001 (admin/admin)
- **Prometheus**: http://localhost:9090

## 📁 Структура проекта

```
НОВСТНОЙ САЙТ/
├── auth-service/           # Сервис аутентификации
├── news-service/           # Сервис новостей
├── seo-service/            # SEO сервис
├── admin-service/          # Админ-панель
├── media-service/          # Медиа сервис
├── gateway/                # API Gateway
├── frontend/               # Next.js приложение
├── scripts/                # Скрипты инициализации
├── monitoring/             # Конфигурация мониторинга
├── docker-compose.yml      # Docker Compose конфигурация
└── ARCHITECTURE.md         # Подробная архитектура
```

## 🔧 Разработка

### Запуск отдельного сервиса локально

```bash
# Auth Service
cd auth-service
cp .env.example .env
go mod download
go run cmd/auth-service/main.go

# News Service
cd news-service
cp .env.example .env
go run cmd/news-service/main.go

# Frontend
cd frontend
npm install
npm run dev
```

### Генерация protobuf файлов

```bash
# Установка protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Генерация для auth-service
cd auth-service
protoc --go_out=. --go-grpc_out=. proto/auth.proto

# Генерация для news-service
cd news-service
protoc --go_out=. --go-grpc_out=. proto/news.proto
```

### Миграции базы данных

```bash
# Установка golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Применение миграций
migrate -path auth-service/migrations -database "postgresql://postgres:password@localhost:5432/auth_db?sslmode=disable" up

# Откат миграций
migrate -path auth-service/migrations -database "postgresql://postgres:password@localhost:5432/auth_db?sslmode=disable" down
```

## 🔐 Аутентификация

Система использует JWT токены:

1. **Access Token** - короткоживущий (15 минут)
2. **Refresh Token** - долгоживущий (7 дней)

### Пример использования

```bash
# Регистрация
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "full_name": "John Doe",
    "role": "user"
  }'

# Вход
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Использование токена
curl http://localhost:8080/api/news \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 📊 Мониторинг

### Prometheus Metrics

Каждый сервис экспортирует метрики:
- HTTP запросов (количество, латентность)
- gRPC вызовов
- Использование памяти и CPU
- Размер очередей

### Grafana Dashboards

Предустановленные дашборды:
- Обзор системы
- Метрики по сервисам
- Database performance
- API Gateway statistics

## 🧪 Тестирование

```bash
# Unit тесты
cd auth-service
go test ./...

# Integration тесты
docker-compose -f docker-compose.test.yml up --abort-on-container-exit

# E2E тесты frontend
cd frontend
npm run test:e2e
```

## 🔄 CI/CD

Проект использует GitHub Actions для автоматического деплоя:

```yaml
# .github/workflows/deploy.yml
- Build Docker images
- Run tests
- Push to registry
- Deploy to Kubernetes
```

## 📚 Документация

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Подробная архитектура системы
- [auth-service/README.md](./auth-service/README.md) - Документация Auth Service
- [news-service/README.md](./news-service/README.md) - Документация News Service
- [gateway/README.md](./gateway/README.md) - Документация API Gateway

## 🛠️ Технологии

**Backend:**
- Go 1.21
- gRPC
- PostgreSQL 15
- Redis 7
- MinIO (S3-compatible storage)
- RabbitMQ
- Consul (Service Discovery)

**Frontend:**
- Next.js 14
- React 18
- TypeScript
- TailwindCSS

**DevOps:**
- Docker & Docker Compose
- Kubernetes (для production)
- Prometheus & Grafana
- Jaeger (Distributed Tracing)

## 🤝 Вклад в проект

1. Fork проекта
2. Создайте feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit изменения (`git commit -m 'Add some AmazingFeature'`)
4. Push в branch (`git push origin feature/AmazingFeature`)
5. Откройте Pull Request

## 📄 Лицензия

MIT License

## 👥 Авторы

- GitHub Copilot

## 📞 Поддержка

Если у вас возникли вопросы, создайте Issue в репозитории.

---

**Статус проекта:** В разработке 🚧
