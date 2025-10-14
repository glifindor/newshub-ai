# 🎯 ПЛАН РЕАЛИЗАЦИИ SEO-SERVICE

## 📋 Обзор

**SEO-Service** - микросервис для управления SEO-оптимизацией новостного портала.

### Основные функции:
- 📄 Автогенерация `sitemap.xml` и `robots.txt`
- 📊 Управление Open Graph, Twitter Cards, Schema.org
- 🔍 Хранение и обновление SEO-метаданных
- 🧭 API для получения SEO-данных по slug
- 🤝 Интеграция с news-service

---

## 🗂️ Структура проекта

```
seo-service/
├── cmd/
│   └── seo-service/
│       └── main.go                    # Точка входа приложения
│
├── internal/
│   ├── config/
│   │   └── config.go                  # Конфигурация (БД, Redis, порты)
│   │
│   ├── models/
│   │   ├── seo_meta.go                # Модель SEO метаданных
│   │   ├── sitemap.go                 # Модель для sitemap
│   │   └── robots.go                  # Модель для robots.txt
│   │
│   ├── repository/
│   │   ├── seo_repository.go          # Работа с БД (CRUD для seo_meta)
│   │   └── news_client.go             # gRPC клиент для news-service
│   │
│   ├── service/
│   │   ├── seo_service.go             # Бизнес-логика SEO
│   │   ├── sitemap_service.go         # Генерация sitemap.xml
│   │   ├── robots_service.go          # Генерация robots.txt
│   │   └── opengraph_service.go       # Open Graph + Twitter Cards
│   │
│   └── handler/
│       ├── http_handler.go            # HTTP endpoints
│       └── grpc_handler.go            # gRPC server (опционально)
│
├── pkg/
│   ├── generator/
│   │   ├── sitemap.go                 # Генератор sitemap.xml
│   │   ├── robots.go                  # Генератор robots.txt
│   │   └── structured_data.go         # JSON-LD генератор
│   │
│   ├── database/
│   │   └── postgres.go                # Подключение к PostgreSQL
│   │
│   └── logger/
│       └── logger.go                  # Zap logger
│
├── migrations/
│   └── 001_create_seo_meta.sql        # Миграция БД
│
├── proto/
│   └── seo.proto                      # gRPC контракт (опционально)
│
├── Dockerfile                          # Docker образ
├── docker-compose.yml                  # Для локальной разработки
├── go.mod
├── go.sum
└── README.md
```

---

## 🗄️ База данных

### Таблица: `seo_meta`

```sql
CREATE TABLE IF NOT EXISTS seo_meta (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    news_id UUID NOT NULL UNIQUE,
    
    -- Основные SEO поля
    title VARCHAR(70) NOT NULL,                    -- SEO title (макс 60-70 символов)
    description VARCHAR(160) NOT NULL,             -- Meta description (макс 160 символов)
    keywords VARCHAR(255),                         -- Keywords через запятую
    slug VARCHAR(500) NOT NULL UNIQUE,             -- URL slug для news
    canonical_url VARCHAR(500),                    -- Canonical URL
    
    -- Open Graph (Facebook)
    og_title VARCHAR(100),                         -- OG title (может отличаться от SEO title)
    og_description VARCHAR(200),                   -- OG description
    og_image VARCHAR(500),                         -- OG image URL
    og_type VARCHAR(50) DEFAULT 'article',         -- article, website, video, etc.
    og_locale VARCHAR(10) DEFAULT 'en_US',         -- Локаль (en_US, ru_RU)
    
    -- Twitter Cards
    twitter_card VARCHAR(50) DEFAULT 'summary_large_image',  -- summary, summary_large_image
    twitter_title VARCHAR(100),                    -- Twitter title
    twitter_description VARCHAR(200),              -- Twitter description
    twitter_image VARCHAR(500),                    -- Twitter image URL
    twitter_creator VARCHAR(50),                   -- @username автора
    
    -- Schema.org (Structured Data)
    schema_type VARCHAR(50) DEFAULT 'NewsArticle', -- NewsArticle, BlogPosting, etc.
    schema_data JSONB,                             -- Полные данные JSON-LD
    
    -- Индексация
    robots_index BOOLEAN DEFAULT true,             -- Разрешить индексацию
    robots_follow BOOLEAN DEFAULT true,            -- Разрешить переход по ссылкам
    sitemap_priority DECIMAL(2,1) DEFAULT 0.5,     -- Приоритет в sitemap (0.0-1.0)
    sitemap_changefreq VARCHAR(20) DEFAULT 'weekly', -- always, hourly, daily, weekly, monthly, yearly, never
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign Key
    CONSTRAINT fk_news FOREIGN KEY (news_id) REFERENCES news(id) ON DELETE CASCADE
);

-- Индексы
CREATE INDEX idx_seo_meta_slug ON seo_meta(slug);
CREATE INDEX idx_seo_meta_news_id ON seo_meta(news_id);
CREATE INDEX idx_seo_meta_updated_at ON seo_meta(updated_at DESC);
CREATE INDEX idx_seo_meta_sitemap ON seo_meta(sitemap_priority, updated_at) WHERE robots_index = true;

-- Trigger для автообновления updated_at
CREATE OR REPLACE FUNCTION update_seo_meta_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER seo_meta_updated_at_trigger
    BEFORE UPDATE ON seo_meta
    FOR EACH ROW
    EXECUTE FUNCTION update_seo_meta_updated_at();
```

---

## 🔌 API Endpoints

### HTTP REST API

#### 1. Получить SEO данные по slug
```http
GET /api/v1/seo/:slug

Response 200:
{
  "id": "uuid",
  "news_id": "uuid",
  "title": "SEO Title",
  "description": "SEO Description",
  "keywords": "keyword1, keyword2",
  "slug": "news-article-slug",
  "canonical_url": "https://example.com/news/news-article-slug",
  "open_graph": {
    "title": "OG Title",
    "description": "OG Description",
    "image": "https://example.com/images/og.jpg",
    "type": "article",
    "locale": "en_US"
  },
  "twitter": {
    "card": "summary_large_image",
    "title": "Twitter Title",
    "description": "Twitter Description",
    "image": "https://example.com/images/twitter.jpg",
    "creator": "@author"
  },
  "schema": {
    "@context": "https://schema.org",
    "@type": "NewsArticle",
    "headline": "Article Headline",
    "datePublished": "2025-10-14T10:00:00Z",
    "author": {...},
    "publisher": {...}
  },
  "robots": {
    "index": true,
    "follow": true
  }
}
```

#### 2. Создать/Обновить SEO метаданные
```http
POST /api/v1/seo/update

Request Body:
{
  "news_id": "uuid",
  "title": "SEO Title",
  "description": "SEO Description",
  "keywords": "keyword1, keyword2",
  "slug": "article-slug",
  "og_title": "Open Graph Title",
  "og_image": "https://example.com/image.jpg",
  "twitter_card": "summary_large_image",
  "sitemap_priority": 0.8
}

Response 201:
{
  "id": "uuid",
  "news_id": "uuid",
  "message": "SEO metadata created successfully"
}
```

#### 3. Получить sitemap.xml
```http
GET /sitemap.xml

Response 200 (application/xml):
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/news/article-slug</loc>
    <lastmod>2025-10-14</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
  ...
</urlset>
```

#### 4. Получить robots.txt
```http
GET /robots.txt

Response 200 (text/plain):
User-agent: *
Allow: /
Disallow: /admin/
Disallow: /api/

Sitemap: https://example.com/sitemap.xml
```

#### 5. Удалить SEO метаданные
```http
DELETE /api/v1/seo/:news_id

Response 204 No Content
```

#### 6. Health Check
```http
GET /health

Response 200:
{
  "status": "healthy",
  "service": "seo-service"
}
```

---

## 🧩 Внутренние компоненты

### 1. **Models** (`internal/models/`)

#### `seo_meta.go`
```go
type SEOMeta struct {
    ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
    NewsID        uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"news_id"`
    
    // Basic SEO
    Title         string     `gorm:"size:70;not null" json:"title"`
    Description   string     `gorm:"size:160;not null" json:"description"`
    Keywords      string     `gorm:"size:255" json:"keywords"`
    Slug          string     `gorm:"size:500;not null;uniqueIndex" json:"slug"`
    CanonicalURL  string     `gorm:"size:500" json:"canonical_url"`
    
    // Open Graph
    OGTitle       string     `gorm:"size:100" json:"og_title"`
    OGDescription string     `gorm:"size:200" json:"og_description"`
    OGImage       string     `gorm:"size:500" json:"og_image"`
    OGType        string     `gorm:"size:50;default:'article'" json:"og_type"`
    OGLocale      string     `gorm:"size:10;default:'en_US'" json:"og_locale"`
    
    // Twitter
    TwitterCard        string `gorm:"size:50;default:'summary_large_image'" json:"twitter_card"`
    TwitterTitle       string `gorm:"size:100" json:"twitter_title"`
    TwitterDescription string `gorm:"size:200" json:"twitter_description"`
    TwitterImage       string `gorm:"size:500" json:"twitter_image"`
    TwitterCreator     string `gorm:"size:50" json:"twitter_creator"`
    
    // Schema.org
    SchemaType string         `gorm:"size:50;default:'NewsArticle'" json:"schema_type"`
    SchemaData datatypes.JSON `gorm:"type:jsonb" json:"schema_data"`
    
    // Indexing
    RobotsIndex       bool    `gorm:"default:true" json:"robots_index"`
    RobotsFollow      bool    `gorm:"default:true" json:"robots_follow"`
    SitemapPriority   float32 `gorm:"type:decimal(2,1);default:0.5" json:"sitemap_priority"`
    SitemapChangefreq string  `gorm:"size:20;default:'weekly'" json:"sitemap_changefreq"`
    
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### `sitemap.go`
```go
type SitemapURL struct {
    Loc        string  `xml:"loc"`
    LastMod    string  `xml:"lastmod,omitempty"`
    ChangeFreq string  `xml:"changefreq,omitempty"`
    Priority   float32 `xml:"priority,omitempty"`
}

type Sitemap struct {
    XMLName xml.Name     `xml:"urlset"`
    XMLNS   string       `xml:"xmlns,attr"`
    URLs    []SitemapURL `xml:"url"`
}
```

---

### 2. **Repository** (`internal/repository/`)

#### `seo_repository.go`
```go
type SEORepository interface {
    Create(ctx context.Context, meta *models.SEOMeta) error
    GetBySlug(ctx context.Context, slug string) (*models.SEOMeta, error)
    GetByNewsID(ctx context.Context, newsID uuid.UUID) (*models.SEOMeta, error)
    Update(ctx context.Context, meta *models.SEOMeta) error
    Delete(ctx context.Context, newsID uuid.UUID) error
    
    // Для sitemap
    GetAllIndexable(ctx context.Context) ([]*models.SEOMeta, error)
    GetRecentlyUpdated(ctx context.Context, limit int) ([]*models.SEOMeta, error)
}
```

#### `news_client.go` (gRPC клиент)
```go
type NewsClient interface {
    GetNewsByID(ctx context.Context, newsID uuid.UUID) (*NewsData, error)
    SubscribeToPublishEvents() (<-chan NewsPublishEvent, error)
}

// Автоматическое создание SEO при публикации новости
```

---

### 3. **Service** (`internal/service/`)

#### `seo_service.go`
```go
type SEOService interface {
    // CRUD операции
    CreateOrUpdate(ctx context.Context, req *CreateSEORequest) (*models.SEOMeta, error)
    GetBySlug(ctx context.Context, slug string) (*SEOResponse, error)
    Delete(ctx context.Context, newsID uuid.UUID) error
    
    // Автогенерация
    GenerateFromNews(ctx context.Context, newsID uuid.UUID) (*models.SEOMeta, error)
    RegenerateAll(ctx context.Context) error
}
```

#### `sitemap_service.go`
```go
type SitemapService interface {
    GenerateSitemap(ctx context.Context) ([]byte, error)
    GetSitemapXML(ctx context.Context) (string, error)
    
    // Кэширование sitemap
    InvalidateCache(ctx context.Context) error
}
```

#### `robots_service.go`
```go
type RobotsService interface {
    GenerateRobotsTxt(ctx context.Context) (string, error)
    GetRobotsTxt(ctx context.Context) (string, error)
}
```

#### `opengraph_service.go`
```go
type OpenGraphService interface {
    GenerateOGTags(news *NewsData) *OpenGraphData
    GenerateTwitterTags(news *NewsData) *TwitterCardData
    GenerateStructuredData(news *NewsData) *StructuredData
}
```

---

### 4. **Generators** (`pkg/generator/`)

#### `sitemap.go`
```go
func GenerateSitemap(urls []*models.SEOMeta, baseURL string) ([]byte, error)
func GenerateSitemapIndex(sitemaps []string) ([]byte, error) // Для больших сайтов
```

#### `robots.go`
```go
func GenerateRobotsTxt(config RobotsConfig) string

type RobotsConfig struct {
    UserAgent       string
    Allow           []string
    Disallow        []string
    CrawlDelay      int
    SitemapURL      string
}
```

#### `structured_data.go`
```go
func GenerateNewsArticleSchema(news *NewsData) *NewsArticleSchema
func GenerateBreadcrumbSchema(news *NewsData) *BreadcrumbSchema
func GenerateOrganizationSchema(config OrgConfig) *OrganizationSchema
```

---

## 🔗 Интеграция с News-Service

### Вариант 1: Event-Driven (рекомендуется)

**News-Service** публикует события в Redis Pub/Sub или RabbitMQ:

```go
// Событие при публикации новости
type NewsPublishedEvent struct {
    NewsID      uuid.UUID
    Title       string
    Slug        string
    Content     string
    FeaturedImage string
    PublishedAt time.Time
}

// SEO-Service подписывается на события
func (s *SEOService) HandleNewsPublished(event NewsPublishedEvent) {
    // Автоматически создаем SEO метаданные
    meta := s.GenerateFromEvent(event)
    s.repository.Create(ctx, meta)
}
```

### Вариант 2: gRPC вызов

**News-Service** вызывает SEO-Service при публикации:

```protobuf
service SEOService {
    rpc CreateSEOMeta(CreateSEORequest) returns (SEOResponse);
    rpc UpdateSEOMeta(UpdateSEORequest) returns (SEOResponse);
}
```

### Вариант 3: HTTP Webhook

**News-Service** отправляет HTTP POST на SEO-Service:

```go
// В news-service после публикации
func (s *NewsService) Publish(newsID uuid.UUID) {
    // ... публикация логика ...
    
    // Уведомляем SEO-Service
    go s.notifySEOService(newsID)
}
```

---

## 📊 Open Graph & Twitter Cards

### Open Graph Tags (для Facebook, LinkedIn, etc.)
```html
<meta property="og:title" content="Article Title">
<meta property="og:description" content="Article Description">
<meta property="og:image" content="https://example.com/image.jpg">
<meta property="og:url" content="https://example.com/news/slug">
<meta property="og:type" content="article">
<meta property="og:locale" content="en_US">
<meta property="article:published_time" content="2025-10-14T10:00:00Z">
<meta property="article:author" content="Author Name">
```

### Twitter Cards
```html
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="Article Title">
<meta name="twitter:description" content="Article Description">
<meta name="twitter:image" content="https://example.com/image.jpg">
<meta name="twitter:creator" content="@author">
```

### Schema.org (JSON-LD)
```json
{
  "@context": "https://schema.org",
  "@type": "NewsArticle",
  "headline": "Article Headline",
  "image": "https://example.com/image.jpg",
  "datePublished": "2025-10-14T10:00:00Z",
  "dateModified": "2025-10-14T12:00:00Z",
  "author": {
    "@type": "Person",
    "name": "Author Name"
  },
  "publisher": {
    "@type": "Organization",
    "name": "News Portal",
    "logo": {
      "@type": "ImageObject",
      "url": "https://example.com/logo.png"
    }
  },
  "description": "Article description"
}
```

---

## 🚀 Конфигурация

### `.env`
```env
# Service
ENVIRONMENT=production
HTTP_PORT=8093
GRPC_PORT=50053

# Database
DB_HOST=news-postgres
DB_PORT=5432
DB_USER=newsportal
DB_PASSWORD=SimplePass123
DB_NAME=newsportal_db

# Redis
REDIS_ADDR=news-redis:6379
REDIS_PASSWORD=
REDIS_DB=0

# SEO Settings
BASE_URL=https://newsportal.com
SITEMAP_CACHE_TTL=3600  # 1 час
ROBOTS_ALLOW=/
ROBOTS_DISALLOW=/admin/,/api/

# News Service (для интеграции)
NEWS_SERVICE_GRPC=news-service:50052
NEWS_SERVICE_HTTP=http://news-service:8092
```

---

## 📦 Docker

### `Dockerfile`
```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o seo-service ./cmd/seo-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
COPY --from=builder /app/seo-service .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8093 50053

CMD ["./seo-service"]
```

### Добавить в `docker-compose.yml`
```yaml
  seo-service:
    build: ./seo-service
    container_name: seo-service
    ports:
      - "8093:8093"  # HTTP
      - "50053:50053"  # gRPC
    environment:
      - ENVIRONMENT=production
      - HTTP_PORT=8093
      - GRPC_PORT=50053
      - DB_HOST=news-postgres
      - DB_PORT=5432
      - DB_USER=newsportal
      - DB_PASSWORD=SimplePass123
      - DB_NAME=newsportal_db
      - REDIS_ADDR=news-redis:6379
      - BASE_URL=https://newsportal.com
      - NEWS_SERVICE_GRPC=news-service:50052
    depends_on:
      - news-postgres
      - news-redis
      - news-service
    networks:
      - news-network
    restart: unless-stopped
```

---

## 🎯 ПОШАГОВЫЙ ПЛАН РЕАЛИЗАЦИИ

---

### ШАГ 1: Создание структуры проекта и базовой конфигурации

**Задачи:**
- [ ] Создать директории проекта
- [ ] Инициализировать Go модуль
- [ ] Создать `internal/config/config.go` для загрузки переменных окружения
- [ ] Создать `.env` файл с конфигурацией
- [ ] Создать `pkg/logger/logger.go` (Zap logger)

**Результат:** Базовая структура проекта готова

---

### ШАГ 2: Создание моделей данных

**Задачи:**
- [ ] Создать `internal/models/seo_meta.go` с полной структурой SEOMeta (24 поля)
- [ ] Создать `internal/models/sitemap.go` для XML структуры sitemap
- [ ] Создать `internal/models/robots.go` для конфигурации robots.txt
- [ ] Создать `internal/models/requests.go` для DTO (CreateSEORequest, UpdateSEORequest)
- [ ] Создать `internal/models/responses.go` для DTO (SEOResponse, OpenGraphResponse)

**Результат:** Все модели данных определены

---

### ШАГ 3: Создание миграции базы данных

**Задачи:**
- [ ] Создать `migrations/001_create_seo_meta.sql`
- [ ] Добавить таблицу `seo_meta` с 24+ полями
- [ ] Добавить индексы (slug, news_id, updated_at, sitemap)
- [ ] Добавить триггер для auto-update `updated_at`
- [ ] Добавить foreign key constraint на `news.id`

**Результат:** SQL миграция готова

---

### ШАГ 4: Подключение к базе данных

**Задачи:**
- [ ] Создать `pkg/database/postgres.go` для подключения через GORM
- [ ] Добавить функцию `NewPostgresDB(config)`
- [ ] Добавить функцию `AutoMigrate(db)` для применения миграций
- [ ] Протестировать подключение

**Результат:** PostgreSQL подключение работает

---

### ШАГ 5: Создание Repository слоя

**Задачи:**
- [ ] Создать `internal/repository/seo_repository.go`
- [ ] Реализовать интерфейс `SEORepository`:
  - `Create(ctx, meta) error`
  - `GetBySlug(ctx, slug) (*SEOMeta, error)`
  - `GetByNewsID(ctx, newsID) (*SEOMeta, error)`
  - `Update(ctx, meta) error`
  - `Delete(ctx, newsID) error`
  - `GetAllIndexable(ctx) ([]*SEOMeta, error)` - для sitemap
  - `GetRecentlyUpdated(ctx, limit) ([]*SEOMeta, error)`

**Результат:** Все CRUD операции для seo_meta реализованы

---

### ШАГ 6: Создание генераторов (pkg/generator)

**Задачи:**
- [ ] Создать `pkg/generator/sitemap.go`:
  - `GenerateSitemap(urls, baseURL) ([]byte, error)` - генерация XML
  - Использовать encoding/xml
- [ ] Создать `pkg/generator/robots.go`:
  - `GenerateRobotsTxt(config) string` - генерация robots.txt
- [ ] Создать `pkg/generator/structured_data.go`:
  - `GenerateNewsArticleSchema(data) *Schema` - JSON-LD для NewsArticle
  - `GenerateOrganizationSchema(config) *Schema` - JSON-LD для Organization

**Результат:** Генераторы для sitemap, robots.txt и Schema.org готовы

---

### ШАГ 7: Создание Service слоя - часть 1 (SEO Service)

**Задачи:**
- [ ] Создать `internal/service/seo_service.go`
- [ ] Реализовать интерфейс `SEOService`:
  - `CreateOrUpdate(ctx, req) (*SEOMeta, error)` - создание/обновление SEO
  - `GetBySlug(ctx, slug) (*SEOResponse, error)` - получение полных данных
  - `Delete(ctx, newsID) error` - удаление
  - `GenerateFromNewsID(ctx, newsID) (*SEOMeta, error)` - автогенерация
- [ ] Добавить логику автозаполнения полей (og_title = title если пусто)
- [ ] Добавить валидацию (title <= 70, description <= 160)

**Результат:** Основной SEO сервис реализован

---

### ШАГ 8: Создание Service слоя - часть 2 (Open Graph)

**Задачи:**
- [ ] Создать `internal/service/opengraph_service.go`
- [ ] Реализовать генерацию Open Graph тегов из SEOMeta
- [ ] Реализовать генерацию Twitter Cards из SEOMeta
- [ ] Создать метод `BuildMetaTags(meta) map[string]string` - возвращает все теги

**Результат:** Open Graph и Twitter Cards сервис готов

---

### ШАГ 9: Создание Service слоя - часть 3 (Sitemap & Robots)

**Задачи:**
- [ ] Создать `internal/service/sitemap_service.go`
- [ ] Реализовать интерфейс `SitemapService`:
  - `GenerateSitemap(ctx) ([]byte, error)` - генерация sitemap.xml
  - `GetCachedSitemap(ctx) ([]byte, error)` - получение из кэша
  - `InvalidateCache(ctx) error` - сброс кэша
- [ ] Добавить Redis кэширование (TTL 1 час)
- [ ] Создать `internal/service/robots_service.go`
- [ ] Реализовать `GenerateRobotsTxt(ctx) (string, error)`

**Результат:** Sitemap и Robots сервисы с кэшированием готовы

---

### ШАГ 10: Создание HTTP Handlers

**Задачи:**
- [ ] Создать `internal/handler/http_handler.go`
- [ ] Реализовать endpoints:
  - `GET /api/v1/seo/:slug` - получить SEO данные
  - `POST /api/v1/seo/update` - создать/обновить SEO
  - `DELETE /api/v1/seo/:news_id` - удалить SEO
  - `GET /sitemap.xml` - получить sitemap
  - `GET /robots.txt` - получить robots.txt
  - `GET /health` - health check
- [ ] Добавить валидацию request body
- [ ] Добавить обработку ошибок

**Результат:** Все HTTP endpoints реализованы

---

### ШАГ 11: Создание main.go и инициализация

**Задачи:**
- [ ] Создать `cmd/seo-service/main.go`
- [ ] Инициализировать:
  - Logger (Zap)
  - Config (из .env)
  - PostgreSQL (GORM)
  - Redis
  - Repository
  - Services
  - HTTP Handlers
- [ ] Настроить Gin router с маршрутами
- [ ] Добавить Graceful Shutdown
- [ ] Запустить HTTP сервер на порту 8093

**Результат:** Сервис запускается и работает локально

---

### ШАГ 12: Интеграция с News-Service (HTTP Webhook)

**Задачи:**
- [ ] Создать endpoint `POST /api/v1/webhook/news-published` в SEO-Service
- [ ] Реализовать обработчик:
  - Получает данные новости (news_id, title, slug, content, image)
  - Автоматически создает SEO метаданные
  - Генерирует Open Graph, Twitter, Schema.org
  - Сохраняет в БД
- [ ] Обновить News-Service:
  - После публикации отправлять POST запрос на SEO-Service
  - Добавить в `news_service.go` метод `notifySEOService(newsID)`

**Результат:** Автоматическое создание SEO при публикации новости

---

### ШАГ 13: Docker контейнеризация

**Задачи:**
- [ ] Создать `Dockerfile` для SEO-Service
- [ ] Создать multi-stage build (builder + runtime)
- [ ] Обновить `docker-compose.yml`:
  - Добавить сервис `seo-service`
  - Порты: 8093 (HTTP)
  - Зависимости: news-postgres, news-redis, news-service
- [ ] Создать `.dockerignore`
- [ ] Протестировать локально с `docker-compose up`

**Результат:** SEO-Service работает в Docker контейнере

---

### ШАГ 14: Деплой на сервер

**Задачи:**
- [ ] Упаковать проект: `tar -czf seo-service.tar.gz seo-service/`
- [ ] Загрузить на сервер: `scp seo-service.tar.gz root@151.241.228.203:/opt/news-portal/`
- [ ] На сервере распаковать: `tar -xzf seo-service.tar.gz`
- [ ] Применить миграции:
  ```bash
  docker exec -i news-postgres psql -U newsportal -d newsportal_db < seo-service/migrations/001_create_seo_meta.sql
  ```
- [ ] Пересобрать docker-compose: `docker-compose up -d --build seo-service`
- [ ] Проверить логи: `docker logs seo-service -f`

**Результат:** SEO-Service развернут на production сервере

---

### ШАГ 15: Тестирование и финализация

**Задачи:**
- [ ] **Тест 1**: Создать SEO метаданные вручную
  ```bash
  curl -X POST http://localhost:8093/api/v1/seo/update \
    -H "Content-Type: application/json" \
    -d '{
      "news_id": "uuid",
      "title": "Test Title",
      "slug": "test-slug"
    }'
  ```
- [ ] **Тест 2**: Получить SEO данные по slug
  ```bash
  curl http://localhost:8093/api/v1/seo/test-slug
  ```
- [ ] **Тест 3**: Проверить sitemap.xml
  ```bash
  curl http://localhost:8093/sitemap.xml
  ```
- [ ] **Тест 4**: Проверить robots.txt
  ```bash
  curl http://localhost:8093/robots.txt
  ```
- [ ] **Тест 5**: Протестировать автоматическое создание SEO при публикации новости
- [ ] Проверить Redis кэширование sitemap
- [ ] Проверить обновление sitemap после создания новых SEO
- [ ] Написать README.md с документацией

**Результат:** Все функции работают, сервис протестирован

---

## 🔄 Сценарии использования

### Сценарий 1: Публикация новости

```
1. User публикует новость через news-service
2. News-service создает запись в таблице news
3. News-service отправляет событие "NewsPublished" 
   → Redis Pub/Sub / HTTP Webhook / gRPC
4. SEO-service получает событие
5. SEO-service генерирует SEO метаданные:
   - title, description, keywords
   - Open Graph tags
   - Twitter Cards
   - Schema.org JSON-LD
6. SEO-service сохраняет в seo_meta таблицу
7. SEO-service инвалидирует кэш sitemap
8. Sitemap автоматически обновится при следующем запросе
```

### Сценарий 2: Запрос страницы новости

```
1. Frontend делает запрос: GET /api/v1/seo/article-slug
2. SEO-service ищет в БД по slug
3. Возвращает все SEO данные:
   - Meta tags (title, description)
   - Open Graph
   - Twitter Cards
   - Structured Data (JSON-LD)
4. Frontend рендерит теги в <head>
```

### Сценарий 3: Обновление sitemap.xml

```
1. Поисковый бот запрашивает: GET /sitemap.xml
2. SEO-service проверяет Redis кэш
3. Если кэша нет:
   - Запрашивает все индексируемые новости из БД
   - Генерирует XML
   - Кэширует в Redis на 1 час
4. Возвращает sitemap.xml
```

---

## 📈 Метрики и мониторинг

### Prometheus метрики:
```
seo_service_requests_total           # Общее количество запросов
seo_service_request_duration_seconds # Длительность запросов
seo_service_sitemap_generation_time  # Время генерации sitemap
seo_service_cache_hits_total         # Попадания в кэш
seo_service_db_queries_total         # Количество запросов к БД
```

### Health Checks:
```
GET /health        # Простой health check
GET /health/ready  # Готовность (БД, Redis доступны)
GET /health/live   # Живость сервиса
```

---

## 🔒 Безопасность

1. **API ключи** для POST/DELETE endpoints
2. **Rate limiting** для публичных endpoints
3. **CORS** настройки для frontend
4. **Валидация** всех входящих данных
5. **SQL injection** защита через GORM
6. **XSS** защита при генерации HTML тегов

---

## ✅ Финальный чеклист

### Обязательные компоненты:
- [ ] Таблица `seo_meta` с полным набором полей
- [ ] CRUD операции для SEO метаданных
- [ ] Генератор sitemap.xml
- [ ] Генератор robots.txt
- [ ] Open Graph теги
- [ ] Twitter Cards
- [ ] Schema.org JSON-LD
- [ ] Интеграция с news-service
- [ ] Redis кэширование
- [ ] Docker контейнер

### Эндпоинты:
- [ ] GET `/api/v1/seo/:slug`
- [ ] POST `/api/v1/seo/update`
- [ ] DELETE `/api/v1/seo/:news_id`
- [ ] GET `/sitemap.xml`
- [ ] GET `/robots.txt`
- [ ] GET `/health`

---

## 📊 Краткая сводка по шагам

| Шаг | Задача | Ключевые файлы |
|-----|--------|----------------|
| 1 | Структура проекта | директории, go.mod, config.go, .env |
| 2 | Модели данных | seo_meta.go, sitemap.go, requests.go |
| 3 | Миграция БД | 001_create_seo_meta.sql |
| 4 | Подключение БД | postgres.go, migrations |
| 5 | Repository | seo_repository.go (7 методов) |
| 6 | Генераторы | sitemap.go, robots.go, structured_data.go |
| 7 | SEO Service | seo_service.go (4 метода) |
| 8 | Open Graph Service | opengraph_service.go |
| 9 | Sitemap/Robots Service | sitemap_service.go, robots_service.go |
| 10 | HTTP Handlers | http_handler.go (6 endpoints) |
| 11 | Main Application | main.go, router, graceful shutdown |
| 12 | Интеграция News | webhook endpoint, auto-create SEO |
| 13 | Docker | Dockerfile, docker-compose.yml |
| 14 | Деплой | загрузка, миграции, запуск |
| 15 | Тестирование | 5 тестовых сценариев, документация |

---

## 🎯 Итого

**Всего шагов**: 15 (четко структурированных)

**Технологический стек**:
- Go 1.23
- Gin (HTTP framework)
- GORM (ORM)
- PostgreSQL (shared DB)
- Redis (caching)
- Docker

**Ключевые компоненты**:
- 3 модели данных
- 1 SQL миграция
- 1 repository (7 методов)
- 4 service (10+ методов)
- 3 генератора
- 6 HTTP endpoints
- 1 webhook для интеграции

**Порты**:
- 8093 - HTTP API

**Зависимости**:
- PostgreSQL (общая БД с news-service)
- Redis (общий кэш)
- News-Service (для webhook интеграции)

---

## 🚀 Начинаем реализацию?

**План готов!** Все шаги четко определены по порядку.

**Подтвердите**, и я начну с **ШАГ 1**: Создание структуры проекта! 🎯
