# Архитектура микросервисного новостного портала

## 📋 Общая схема системы

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend                             │
│                    (Next.js - Port 3000)                     │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTP/REST
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                      API Gateway                             │
│              (Golang Gin/Echo - Port 8080)                   │
│    • Routing • Auth Middleware • Rate Limiting               │
└──┬────────┬─────────┬──────────┬─────────┬─────────────────┘
   │        │         │          │         │
   │ gRPC   │ gRPC    │ gRPC     │ gRPC    │ gRPC
   ▼        ▼         ▼          ▼         ▼
┌─────┐ ┌──────┐ ┌──────┐ ┌───────┐ ┌──────────┐
│Auth │ │News  │ │SEO   │ │Admin  │ │Media     │
│:8081│ │:8082 │ │:8083 │ │:8084  │ │:8085     │
└──┬──┘ └──┬───┘ └──┬───┘ └───┬───┘ └────┬─────┘
   │       │        │         │          │
   └───────┴────────┴─────────┴──────────┘
                    │
        ┌───────────┴──────────────┐
        ▼                          ▼
┌──────────────┐          ┌─────────────────┐
│  PostgreSQL  │          │     Redis       │
│  (Port 5432) │          │  (Port 6379)    │
│              │          │  • Cache        │
│ • Users      │          │  • Sessions     │
│ • News       │          │  • Rate Limit   │
│ • Categories │          └─────────────────┘
│ • Tags       │
│ • Media      │          ┌─────────────────┐
└──────────────┘          │   Message Queue │
                          │  RabbitMQ/NATS  │
                          │  (Port 5672)    │
                          └─────────────────┘
```

---

## 🏗️ Микросервисы и их ответственность

### 1. **Auth Service** (Port 8081)
**Ответственность:**
- Регистрация и авторизация пользователей
- Управление JWT токенами (access + refresh)
- Управление ролями и правами доступа
- OAuth2 интеграция (Google, Facebook)

**Технологии:**
- JWT для токенов
- bcrypt для хеширования паролей
- PostgreSQL для хранения пользователей
- Redis для хранения refresh токенов и сессий

---

### 2. **News Service** (Port 8082)
**Ответственность:**
- CRUD операции с новостями
- Управление категориями и тегами
- Поиск и фильтрация новостей
- Версионирование статей
- Публикация/черновики/архив

**Технологии:**
- PostgreSQL для хранения новостей
- Elasticsearch для полнотекстового поиска
- Redis для кеширования популярных новостей

---

### 3. **SEO Service** (Port 8083)
**Ответственность:**
- Генерация метатегов (title, description, keywords)
- Open Graph и Twitter Cards
- Генерация sitemap.xml
- Robots.txt управление
- Структурированные данные (JSON-LD, Schema.org)
- Канонические URL

**Технологии:**
- Взаимодействие с News Service через gRPC
- Кеширование метаданных в Redis
- Периодическая генерация sitemap

---

### 4. **Admin Service** (Port 8084)
**Ответственность:**
- Панель модерации контента
- Управление пользователями и ролями
- Аналитика и статистика
- Управление категориями/тегами
- Логи и аудит действий

**Технологии:**
- gRPC клиенты для всех сервисов
- Агрегация данных из разных источников
- Права доступа через Auth Service

---

### 5. **Media Service** (Port 8085)
**Ответственность:**
- Загрузка изображений и видео
- Ресайз и оптимизация изображений
- Хранение файлов (S3/MinIO)
- CDN интеграция
- Генерация превью

**Технологии:**
- MinIO/S3 для хранения
- ImageMagick/Sharp для обработки
- CDN для раздачи статики

---

### 6. **API Gateway** (Port 8080)
**Ответственность:**
- Единая точка входа для всех клиентов
- Маршрутизация запросов к микросервисам
- Аутентификация и авторизация (проверка JWT)
- Rate limiting и throttling
- Логирование и мониторинг
- CORS настройка

**Технологии:**
- Gin/Echo framework
- gRPC клиенты для всех сервисов
- JWT middleware
- Redis для rate limiting

---

## 📁 Файловая структура микросервисов

### Auth Service
```
/auth-service
├── /cmd
│   └── /auth-service
│       └── main.go                 # Точка входа
├── /internal
│   ├── /config
│   │   └── config.go               # Конфигурация
│   ├── /handler
│   │   ├── grpc_handler.go         # gRPC обработчики
│   │   └── http_handler.go         # HTTP обработчики (опционально)
│   ├── /service
│   │   ├── auth_service.go         # Бизнес-логика аутентификации
│   │   ├── token_service.go        # Работа с JWT
│   │   └── user_service.go         # Управление пользователями
│   ├── /repository
│   │   ├── user_repository.go      # Работа с БД пользователей
│   │   └── session_repository.go   # Redis для сессий
│   ├── /middleware
│   │   └── auth_middleware.go      # Middleware для проверки токенов
│   └── /models
│       ├── user.go                 # Модель пользователя
│       └── token.go                # Модель токена
├── /pkg
│   ├── /jwt
│   │   └── jwt.go                  # Утилиты для JWT
│   ├── /hash
│   │   └── bcrypt.go               # Хеширование паролей
│   └── /validator
│       └── validator.go            # Валидация данных
├── /proto
│   └── auth.proto                  # gRPC контракты
├── /migrations
│   ├── 001_create_users_table.up.sql
│   └── 001_create_users_table.down.sql
├── .env
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

---

### News Service
```
/news-service
├── /cmd
│   └── /news-service
│       └── main.go
├── /internal
│   ├── /config
│   │   └── config.go
│   ├── /handler
│   │   └── grpc_handler.go
│   ├── /service
│   │   ├── news_service.go         # CRUD новостей
│   │   ├── category_service.go     # Управление категориями
│   │   ├── tag_service.go          # Управление тегами
│   │   └── search_service.go       # Поиск новостей
│   ├── /repository
│   │   ├── news_repository.go
│   │   ├── category_repository.go
│   │   └── tag_repository.go
│   ├── /models
│   │   ├── news.go
│   │   ├── category.go
│   │   └── tag.go
│   └── /cache
│       └── redis_cache.go          # Кеширование
├── /pkg
│   ├── /slug
│   │   └── slug.go                 # Генерация slug'ов
│   └── /pagination
│       └── pagination.go
├── /proto
│   └── news.proto
├── /migrations
│   ├── 001_create_news_table.up.sql
│   ├── 002_create_categories_table.up.sql
│   └── 003_create_tags_table.up.sql
├── .env
├── go.mod
├── Dockerfile
└── README.md
```

---

### SEO Service
```
/seo-service
├── /cmd
│   └── /seo-service
│       └── main.go
├── /internal
│   ├── /config
│   │   └── config.go
│   ├── /handler
│   │   └── grpc_handler.go
│   ├── /service
│   │   ├── meta_service.go         # Генерация метатегов
│   │   ├── sitemap_service.go      # Генерация sitemap
│   │   ├── opengraph_service.go    # Open Graph
│   │   └── schema_service.go       # Schema.org разметка
│   ├── /repository
│   │   └── seo_repository.go
│   ├── /models
│   │   ├── meta.go
│   │   └── sitemap.go
│   └── /templates
│       ├── sitemap.xml.tmpl
│       └── schema.json.tmpl
├── /pkg
│   └── /generator
│       └── meta_generator.go
├── /proto
│   └── seo.proto
├── .env
├── go.mod
├── Dockerfile
└── README.md
```

---

### Admin Service
```
/admin-service
├── /cmd
│   └── /admin-service
│       └── main.go
├── /internal
│   ├── /config
│   │   └── config.go
│   ├── /handler
│   │   └── grpc_handler.go
│   ├── /service
│   │   ├── moderation_service.go   # Модерация контента
│   │   ├── analytics_service.go    # Аналитика
│   │   └── audit_service.go        # Логи действий
│   ├── /repository
│   │   └── audit_repository.go
│   ├── /models
│   │   ├── audit_log.go
│   │   └── statistics.go
│   └── /clients
│       ├── news_client.go          # gRPC клиент для News Service
│       ├── auth_client.go          # gRPC клиент для Auth Service
│       └── media_client.go         # gRPC клиент для Media Service
├── /proto
│   └── admin.proto
├── .env
├── go.mod
├── Dockerfile
└── README.md
```

---

### Media Service
```
/media-service
├── /cmd
│   └── /media-service
│       └── main.go
├── /internal
│   ├── /config
│   │   └── config.go
│   ├── /handler
│   │   └── grpc_handler.go
│   ├── /service
│   │   ├── upload_service.go       # Загрузка файлов
│   │   ├── resize_service.go       # Ресайз изображений
│   │   └── storage_service.go      # Работа с S3/MinIO
│   ├── /repository
│   │   └── media_repository.go
│   ├── /models
│   │   └── media.go
│   └── /processor
│       └── image_processor.go      # Обработка изображений
├── /pkg
│   ├── /s3
│   │   └── client.go               # S3/MinIO клиент
│   └── /resize
│       └── resize.go
├── /proto
│   └── media.proto
├── .env
├── go.mod
├── Dockerfile
└── README.md
```

---

### API Gateway
```
/gateway
├── /cmd
│   └── /gateway
│       └── main.go
├── /internal
│   ├── /config
│   │   └── config.go
│   ├── /handler
│   │   ├── auth_handler.go         # Проксирование к Auth Service
│   │   ├── news_handler.go         # Проксирование к News Service
│   │   ├── seo_handler.go          # Проксирование к SEO Service
│   │   ├── admin_handler.go        # Проксирование к Admin Service
│   │   └── media_handler.go        # Проксирование к Media Service
│   ├── /middleware
│   │   ├── auth_middleware.go      # JWT проверка
│   │   ├── rate_limit.go           # Rate limiting
│   │   ├── cors.go                 # CORS настройки
│   │   └── logger.go               # Логирование
│   ├── /router
│   │   └── router.go               # Настройка маршрутов
│   └── /clients
│       ├── auth_client.go          # gRPC клиент
│       ├── news_client.go
│       ├── seo_client.go
│       ├── admin_client.go
│       └── media_client.go
├── /pkg
│   └── /response
│       └── response.go             # Стандартизированные ответы
├── .env
├── go.mod
├── Dockerfile
└── README.md
```

---

### Frontend (Next.js)
```
/frontend
├── /src
│   ├── /app
│   │   ├── /news
│   │   │   ├── page.tsx            # Список новостей
│   │   │   └── /[slug]
│   │   │       └── page.tsx        # Страница новости
│   │   ├── /category
│   │   │   └── /[slug]
│   │   │       └── page.tsx
│   │   ├── /admin
│   │   │   └── page.tsx
│   │   ├── layout.tsx
│   │   └── page.tsx
│   ├── /components
│   │   ├── /news
│   │   │   ├── NewsCard.tsx
│   │   │   └── NewsList.tsx
│   │   ├── /layout
│   │   │   ├── Header.tsx
│   │   │   └── Footer.tsx
│   │   └── /admin
│   │       └── AdminPanel.tsx
│   ├── /lib
│   │   ├── /api
│   │   │   ├── auth.ts             # API клиент для Auth
│   │   │   ├── news.ts             # API клиент для News
│   │   │   └── media.ts
│   │   └── /utils
│   │       ├── jwt.ts
│   │       └── helpers.ts
│   └── /types
│       ├── news.ts
│       └── user.ts
├── /public
│   ├── /images
│   └── /icons
├── next.config.js
├── package.json
└── tsconfig.json
```

---

## 🔄 Взаимодействие между сервисами

### Протокол взаимодействия

**gRPC** - основной протокол для межсервисного взаимодействия:
- ✅ Высокая производительность (HTTP/2, protobuf)
- ✅ Строгая типизация через .proto файлы
- ✅ Встроенная поддержка streaming
- ✅ Автогенерация клиентов

**REST/HTTP** - для взаимодействия Gateway с Frontend:
- ✅ Универсальность
- ✅ Простота отладки
- ✅ Кеширование на уровне HTTP

### Схема взаимодействия

```
Frontend (Next.js)
    ↓ HTTP/REST (JSON)
API Gateway (Port 8080)
    ↓ gRPC (Protobuf)
[Auth | News | SEO | Admin | Media] Services
    ↓
[PostgreSQL | Redis | MinIO]
```

### Примеры взаимодействия

#### 1. Создание новости
```
1. Frontend → Gateway: POST /api/news
   Headers: Authorization: Bearer <JWT>
   Body: { title, content, category_id }

2. Gateway → Auth Service (gRPC): ValidateToken(token)
   Response: { user_id, role }

3. Gateway → News Service (gRPC): CreateNews(user_id, news_data)
   Response: { news_id, slug }

4. Gateway → SEO Service (gRPC): GenerateMeta(news_id)
   Response: { meta_tags, og_tags }

5. Gateway → Frontend: { success: true, news_id, slug }
```

#### 2. Получение новости с SEO
```
1. Frontend → Gateway: GET /api/news/[slug]

2. Gateway → News Service (gRPC): GetNewsBySlug(slug)
   Response: { news_data }

3. Gateway → SEO Service (gRPC): GetMeta(news_id)
   Response: { meta, og, schema }

4. Gateway → Frontend: { news, seo }
```

---

## 🔐 Аутентификация и авторизация

### JWT Token Flow

```
┌──────────┐                          ┌──────────┐
│ Frontend │                          │ Gateway  │
└────┬─────┘                          └────┬─────┘
     │                                     │
     │ POST /api/auth/login                │
     │ { email, password }                 │
     │────────────────────────────────────>│
     │                                     │ ValidateCredentials(email, password)
     │                                     │──────────────────────┐
     │                                     │                      │
     │                                     │<─────────────────────┘
     │                                     │
     │ { access_token, refresh_token }     │
     │<────────────────────────────────────│
     │                                     │
     │ GET /api/news (protected)           │
     │ Authorization: Bearer <access>      │
     │────────────────────────────────────>│
     │                                     │ ValidateToken(access_token)
     │                                     │──────────────────────┐
     │                                     │                      │
     │                                     │<─────────────────────┘
     │                                     │
     │              { news_data }          │
     │<────────────────────────────────────│
```

### JWT Структура

**Access Token** (15 минут):
```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "editor",
  "permissions": ["create_news", "edit_news"],
  "exp": 1234567890,
  "iat": 1234567000
}
```

**Refresh Token** (7 дней):
- Хранится в Redis с TTL
- Используется для обновления access token

### Роли и права

```go
// Roles
const (
    RoleAdmin     = "admin"      // Полный доступ
    RoleEditor    = "editor"     // Создание/редактирование новостей
    RoleModerator = "moderator"  // Модерация контента
    RoleUser      = "user"       // Чтение, комментарии
)

// Permissions
const (
    PermCreateNews   = "create_news"
    PermEditNews     = "edit_news"
    PermDeleteNews   = "delete_news"
    PermModerate     = "moderate"
    PermManageUsers  = "manage_users"
)
```

---

## 🌐 Service Discovery

### Consul-based Discovery

```go
// Регистрация сервиса в Consul
func RegisterService(consul *api.Client, serviceName string, port int) error {
    registration := &api.AgentServiceRegistration{
        ID:      fmt.Sprintf("%s-%s", serviceName, uuid.New()),
        Name:    serviceName,
        Port:    port,
        Address: getLocalIP(),
        Check: &api.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://%s:%d/health", getLocalIP(), port),
            Interval: "10s",
            Timeout:  "3s",
        },
    }
    return consul.Agent().ServiceRegister(registration)
}

// Обнаружение сервиса
func DiscoverService(consul *api.Client, serviceName string) (string, error) {
    services, _, err := consul.Health().Service(serviceName, "", true, nil)
    if err != nil || len(services) == 0 {
        return "", err
    }
    
    // Load balancing - random selection
    service := services[rand.Intn(len(services))]
    return fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port), nil
}
```

### Альтернатива: Static Configuration

```yaml
# gateway/config.yaml
services:
  auth:
    host: auth-service
    port: 8081
  news:
    host: news-service
    port: 8082
  seo:
    host: seo-service
    port: 8083
  admin:
    host: admin-service
    port: 8084
  media:
    host: media-service
    port: 8085
```

---

## 🚦 API Gateway - Маршрутизация

### Структура маршрутов

```go
// internal/router/router.go
func SetupRoutes(r *gin.Engine, clients *Clients) {
    // Public routes
    public := r.Group("/api")
    {
        // Auth
        public.POST("/auth/register", handlers.Register(clients.Auth))
        public.POST("/auth/login", handlers.Login(clients.Auth))
        public.POST("/auth/refresh", handlers.RefreshToken(clients.Auth))
        
        // News (public)
        public.GET("/news", handlers.GetNews(clients.News))
        public.GET("/news/:slug", handlers.GetNewsBySlug(clients.News))
        public.GET("/categories", handlers.GetCategories(clients.News))
        
        // SEO
        public.GET("/sitemap.xml", handlers.GetSitemap(clients.SEO))
        public.GET("/robots.txt", handlers.GetRobots(clients.SEO))
    }
    
    // Protected routes (JWT required)
    protected := r.Group("/api")
    protected.Use(middleware.AuthMiddleware(clients.Auth))
    {
        // News management
        protected.POST("/news", middleware.RequireRole("editor"), 
            handlers.CreateNews(clients.News))
        protected.PUT("/news/:id", middleware.RequireRole("editor"),
            handlers.UpdateNews(clients.News))
        protected.DELETE("/news/:id", middleware.RequireRole("admin"),
            handlers.DeleteNews(clients.News))
        
        // Media
        protected.POST("/media/upload", handlers.UploadMedia(clients.Media))
    }
    
    // Admin routes
    admin := r.Group("/api/admin")
    admin.Use(middleware.AuthMiddleware(clients.Auth))
    admin.Use(middleware.RequireRole("admin", "moderator"))
    {
        admin.GET("/users", handlers.GetUsers(clients.Admin))
        admin.GET("/statistics", handlers.GetStatistics(clients.Admin))
        admin.POST("/moderate/:id", handlers.ModerateContent(clients.Admin))
    }
}
```

### Rate Limiting

```go
// internal/middleware/rate_limit.go
func RateLimitMiddleware(redis *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        key := fmt.Sprintf("rate_limit:%s", ip)
        
        count, err := redis.Incr(ctx, key).Result()
        if err != nil {
            c.AbortWithStatus(500)
            return
        }
        
        if count == 1 {
            redis.Expire(ctx, key, time.Minute)
        }
        
        // 100 запросов в минуту
        if count > 100 {
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

---

## 📊 Порты сервисов

| Сервис          | HTTP Port | gRPC Port | Описание                    |
|-----------------|-----------|-----------|----------------------------|
| Frontend        | 3000      | -         | Next.js приложение         |
| API Gateway     | 8080      | -         | REST API точка входа       |
| Auth Service    | -         | 8081      | Аутентификация             |
| News Service    | -         | 8082      | Управление новостями       |
| SEO Service     | -         | 8083      | SEO метаданные             |
| Admin Service   | -         | 8084      | Админ-панель               |
| Media Service   | -         | 8085      | Медиа файлы                |
| PostgreSQL      | 5432      | -         | База данных                |
| Redis           | 6379      | -         | Кеш и сессии               |
| RabbitMQ        | 5672/15672| -         | Очередь сообщений          |
| Consul          | 8500      | -         | Service Discovery          |
| MinIO           | 9000/9001 | -         | S3-совместимое хранилище   |

---

## 🔧 Конфигурация сервисов

### Пример .env для Auth Service

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
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Service Discovery
CONSUL_HOST=localhost
CONSUL_PORT=8500

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json
```

---

## 🐳 Docker Compose

```yaml
version: '3.8'

services:
  # Databases
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: news_portal
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  # Message Queue
  rabbitmq:
    image: rabbitmq:3-management-alpine
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: password

  # Service Discovery
  consul:
    image: consul:latest
    ports:
      - "8500:8500"
    command: agent -server -ui -bootstrap-expect=1 -client=0.0.0.0

  # Storage
  minio:
    image: minio/minio
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: admin
      MINIO_ROOT_PASSWORD: password123
    command: server /data --console-address ":9001"
    volumes:
      - minio_data:/data

  # Microservices
  auth-service:
    build: ./auth-service
    ports:
      - "8081:8081"
    depends_on:
      - postgres
      - redis
      - consul
    environment:
      DB_HOST: postgres
      REDIS_HOST: redis
      CONSUL_HOST: consul

  news-service:
    build: ./news-service
    ports:
      - "8082:8082"
    depends_on:
      - postgres
      - redis
      - consul
    environment:
      DB_HOST: postgres
      REDIS_HOST: redis
      CONSUL_HOST: consul

  seo-service:
    build: ./seo-service
    ports:
      - "8083:8083"
    depends_on:
      - postgres
      - redis
      - consul
    environment:
      DB_HOST: postgres
      REDIS_HOST: redis
      CONSUL_HOST: consul

  admin-service:
    build: ./admin-service
    ports:
      - "8084:8084"
    depends_on:
      - consul
    environment:
      CONSUL_HOST: consul

  media-service:
    build: ./media-service
    ports:
      - "8085:8085"
    depends_on:
      - postgres
      - minio
      - consul
    environment:
      DB_HOST: postgres
      MINIO_HOST: minio
      CONSUL_HOST: consul

  gateway:
    build: ./gateway
    ports:
      - "8080:8080"
    depends_on:
      - auth-service
      - news-service
      - seo-service
      - admin-service
      - media-service
      - consul
    environment:
      CONSUL_HOST: consul
      REDIS_HOST: redis

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    depends_on:
      - gateway
    environment:
      NEXT_PUBLIC_API_URL: http://gateway:8080

volumes:
  postgres_data:
  redis_data:
  minio_data:
```

---

## 📈 Масштабирование

### Горизонтальное масштабирование

**Stateless сервисы** (легко масштабируются):
```bash
# Запуск нескольких экземпляров
docker-compose up --scale news-service=3 --scale auth-service=2
```

**Балансировка нагрузки:**
- Consul автоматически регистрирует все инстансы
- Gateway использует round-robin или random selection
- Health checks отсеивают неработающие инстансы

### Кеширование

**Redis кеширование на разных уровнях:**

```go
// News Service - кеш популярных новостей
func (s *NewsService) GetNews(ctx context.Context, id string) (*News, error) {
    // Проверяем кеш
    cached, err := s.cache.Get(ctx, fmt.Sprintf("news:%s", id))
    if err == nil {
        return unmarshal(cached), nil
    }
    
    // Запрос к БД
    news, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Кешируем на 5 минут
    s.cache.Set(ctx, fmt.Sprintf("news:%s", id), news, 5*time.Minute)
    return news, nil
}
```

**CDN для статики:**
- Изображения раздаются через CloudFlare/CloudFront
- Media Service генерирует signed URLs

### Database Sharding

**По категориям/регионам:**
```
DB1: Политика, Экономика
DB2: Спорт, Культура
DB3: Технологии, Наука
```

**По времени:**
```
DB1: Новости текущего года
DB2: Архив прошлых лет
```

### Message Queue для асинхронных задач

```go
// Отправка email уведомлений
func (s *NewsService) PublishNews(news *News) error {
    // Сохраняем новость
    err := s.repo.Save(news)
    if err != nil {
        return err
    }
    
    // Отправляем событие в очередь
    event := Event{
        Type: "news.published",
        Data: news,
    }
    return s.queue.Publish("notifications", event)
}

// Воркер обрабатывает события
func (w *NotificationWorker) ProcessEvents() {
    msgs := w.queue.Consume("notifications")
    
    for msg := range msgs {
        switch msg.Type {
        case "news.published":
            w.sendEmailNotifications(msg.Data)
        }
    }
}
```

### Monitoring & Observability

**Prometheus + Grafana:**
```go
// Метрики
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"service", "method", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration",
        },
        []string{"service", "method"},
    )
)
```

**Distributed Tracing (Jaeger):**
```go
import "go.opentelemetry.io/otel"

// Трейсинг запросов между сервисами
func (h *Handler) CreateNews(ctx context.Context, req *pb.CreateNewsRequest) {
    ctx, span := otel.Tracer("news-service").Start(ctx, "CreateNews")
    defer span.End()
    
    // Вызов других сервисов сохраняет trace context
    meta, err := h.seoClient.GenerateMeta(ctx, newsID)
}
```

---

## 🔒 Безопасность

### 1. **API Gateway защита**
- JWT валидация на каждом запросе
- Rate limiting (100 req/min per IP)
- CORS настройки
- Input validation

### 2. **Межсервисное взаимодействие**
- mTLS для gRPC (опционально)
- Service mesh (Istio/Linkerd)
- API keys для сервисов

### 3. **Database Security**
- Prepared statements (защита от SQL injection)
- Least privilege principle
- Encrypted connections

### 4. **Secrets Management**
- HashiCorp Vault для хранения секретов
- Environment variables через Docker secrets

---

## 🚀 Деплой и CI/CD

### GitHub Actions Pipeline

```yaml
name: Deploy Microservices

on:
  push:
    branches: [main]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        service: [auth-service, news-service, seo-service, admin-service, media-service, gateway]
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Docker image
        run: |
          docker build -t ${{ secrets.DOCKER_REGISTRY }}/${{ matrix.service }}:${{ github.sha }} ./${{ matrix.service }}
      
      - name: Push to registry
        run: |
          docker push ${{ secrets.DOCKER_REGISTRY }}/${{ matrix.service }}:${{ github.sha }}
      
      - name: Deploy to Kubernetes
        run: |
          kubectl set image deployment/${{ matrix.service }} \
            ${{ matrix.service }}=${{ secrets.DOCKER_REGISTRY }}/${{ matrix.service }}:${{ github.sha }}
```

### Kubernetes Deployment

```yaml
# news-service-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: news-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: news-service
  template:
    metadata:
      labels:
        app: news-service
    spec:
      containers:
      - name: news-service
        image: registry/news-service:latest
        ports:
        - containerPort: 8082
        env:
        - name: DB_HOST
          value: postgres-service
        - name: REDIS_HOST
          value: redis-service
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          grpc:
            port: 8082
          initialDelaySeconds: 10
          periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: news-service
spec:
  selector:
    app: news-service
  ports:
  - port: 8082
    targetPort: 8082
  type: ClusterIP
```

---

## 📝 Итоги архитектуры

### Преимущества микросервисной архитектуры

✅ **Независимое развертывание** - каждый сервис деплоится отдельно  
✅ **Масштабируемость** - масштабируем только нужные сервисы  
✅ **Технологическая гибкость** - можем использовать разные технологии  
✅ **Отказоустойчивость** - падение одного сервиса не ломает всю систему  
✅ **Командная автономия** - разные команды работают над разными сервисами  

### Компромиссы

⚠️ **Сложность инфраструктуры** - требуется DevOps экспертиза  
⚠️ **Распределенные транзакции** - сложнее обеспечить ACID  
⚠️ **Latency** - межсервисное взаимодействие добавляет задержки  
⚠️ **Debugging** - сложнее отследить проблему по всей цепочке  

---

## 🎯 Следующие шаги

1. **Генерация proto файлов** для gRPC контрактов
2. **Настройка CI/CD** пайплайна
3. **Реализация базовых CRUD** операций
4. **Интеграция мониторинга** (Prometheus, Grafana, Jaeger)
5. **Написание интеграционных тестов**
6. **Настройка Kubernetes** для production
7. **Документация API** (Swagger/OpenAPI)

---

**Автор:** GitHub Copilot  
**Дата:** 2025-10-14  
**Версия:** 1.0
