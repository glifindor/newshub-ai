# NEWS-SERVICE: Инженерный план (обновлённый)

## 📋 Общая информация

**Версия:** 2.0 (без gRPC, только HTTP REST API)  
**Статус:** ✅ Реализовано, требуется исправление импортов  
**Порты:** 8092 (HTTP)

---

## 🏗️ Архитектура

### Стек технологий
- **Framework:** Gin 1.9.1
- **ORM:** GORM 1.25.5 (PostgreSQL)
- **Cache:** Redis 7
- **Logger:** Zap 1.26.0
- **Validator:** go-playground/validator 10.16.0
- **Utils:** gosimple/slug (SEO-friendly URLs)

### Структура проекта

```
news-service/
├── cmd/
│   └── news-service/
│       └── main.go                    # Entry point
├── internal/
│   ├── config/
│   │   └── config.go                  # Env configuration
│   ├── handler/
│   │   └── http_handler.go            # 18 HTTP endpoints
│   ├── middleware/
│   │   └── auth_middleware.go         # JWT validation (опционально)
│   ├── model/
│   │   ├── news.go                    # News model
│   │   ├── category.go                # Category model
│   │   └── tag.go                     # Tag model
│   ├── repository/
│   │   ├── news_repository.go         # GORM DAO
│   │   ├── category_repository.go
│   │   └── tag_repository.go
│   └── service/
│       ├── news_service.go            # Business logic
│       ├── category_service.go
│       ├── tag_service.go
│       └── seo_service_client.go      # (будущее) gRPC client to SEO service
├── pkg/
│   ├── database/
│   │   └── postgres.go                # GORM connection
│   ├── logger/
│   │   └── logger.go                  # Zap logger
│   └── utils/
│       ├── response.go                # JSON response helpers
│       ├── pagination.go              # Pagination helper
│       └── validation.go              # Custom validators
├── go.mod
├── go.sum
└── Dockerfile
```

---

## 🗄️ База данных

### Таблица: `news`

```sql
CREATE TABLE news (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    slug VARCHAR(500) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    
    -- SEO
    seo_title VARCHAR(60),
    seo_description VARCHAR(160),
    seo_keywords VARCHAR(255),
    
    -- Relations
    category_id BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    author_id BIGINT,  -- FK to auth-service users
    
    -- Media
    featured_image VARCHAR(500),
    images JSONB,  -- массив URL изображений
    
    -- Status & Visibility
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    is_featured BOOLEAN DEFAULT FALSE,
    is_breaking BOOLEAN DEFAULT FALSE,
    
    -- Stats
    views INTEGER DEFAULT 0,
    
    -- Timestamps
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP  -- Soft delete
);

CREATE INDEX idx_news_slug ON news(slug);
CREATE INDEX idx_news_status ON news(status);
CREATE INDEX idx_news_category ON news(category_id);
CREATE INDEX idx_news_published ON news(published_at);
CREATE INDEX idx_news_featured ON news(is_featured);
```

### Таблица: `categories`

```sql
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    parent_id BIGINT REFERENCES categories(id),  -- Иерархия
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Таблица: `tags`

```sql
CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Таблица связи: `news_tags`

```sql
CREATE TABLE news_tags (
    news_id BIGINT REFERENCES news(id) ON DELETE CASCADE,
    tag_id BIGINT REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (news_id, tag_id)
);
```

---

## 📡 REST API Endpoints

### News Management

| Method | Endpoint | Handler | Auth | Описание |
|--------|----------|---------|------|----------|
| GET | `/api/v1/news` | ListNews | - | Список новостей с фильтрами |
| POST | `/api/v1/news` | CreateNews | ✅ | Создать новость |
| GET | `/api/v1/news/:id` | GetNews | - | Получить по ID |
| GET | `/api/v1/news/slug/:slug` | GetNewsBySlug | - | Получить по slug |
| PUT | `/api/v1/news/:id` | UpdateNews | ✅ | Обновить |
| DELETE | `/api/v1/news/:id` | DeleteNews | ✅ | Удалить (soft) |
| PATCH | `/api/v1/news/:id/publish` | PublishNews | ✅ | Опубликовать |
| GET | `/api/v1/news/featured` | GetFeaturedNews | - | Избранные |
| GET | `/api/v1/news/breaking` | GetBreakingNews | - | Срочные |

### Categories

| Method | Endpoint | Handler | Auth | Описание |
|--------|----------|---------|------|----------|
| GET | `/api/v1/categories` | ListCategories | - | Все категории |
| POST | `/api/v1/categories` | CreateCategory | ✅ | Создать |
| GET | `/api/v1/categories/:id` | GetCategory | - | Получить |
| PUT | `/api/v1/categories/:id` | UpdateCategory | ✅ | Обновить |
| DELETE | `/api/v1/categories/:id` | DeleteCategory | ✅ | Удалить |

### Tags

| Method | Endpoint | Handler | Auth | Описание |
|--------|----------|---------|------|----------|
| GET | `/api/v1/tags` | ListTags | - | Все теги |
| POST | `/api/v1/tags` | CreateTag | ✅ | Создать |
| GET | `/api/v1/tags/:id` | GetTag | - | Получить |
| PUT | `/api/v1/tags/:id` | UpdateTag | ✅ | Обновить |
| DELETE | `/api/v1/tags/:id` | DeleteTag | ✅ | Удалить |

---

## 🔍 Full-Text Search

### PostgreSQL FTS

```sql
-- Добавить tsvector колонку
ALTER TABLE news ADD COLUMN search_vector tsvector;

-- Создать GIN индекс
CREATE INDEX idx_news_search ON news USING GIN(search_vector);

-- Триггер для автообновления
CREATE TRIGGER news_search_update 
BEFORE INSERT OR UPDATE ON news
FOR EACH ROW EXECUTE FUNCTION
tsvector_update_trigger(search_vector, 'pg_catalog.russian', title, content);
```

### Search endpoint

```go
// GET /api/v1/news/search?q=keyword&page=1&limit=20
func (r *NewsRepository) Search(ctx context.Context, query string, page, limit int) ([]*News, int64, error) {
    var news []*News
    var total int64
    
    offset := (page - 1) * limit
    
    err := r.db.WithContext(ctx).
        Where("search_vector @@ plainto_tsquery('russian', ?)", query).
        Order("ts_rank(search_vector, plainto_tsquery('russian', ?)) DESC", query).
        Offset(offset).
        Limit(limit).
        Find(&news).Error
    
    // Count total
    r.db.Model(&News{}).
        Where("search_vector @@ plainto_tsquery('russian', ?)", query).
        Count(&total)
    
    return news, total, err
}
```

---

## 🎯 Фильтрация и пагинация

### Query Parameters

```
GET /api/v1/news?
    category=tech           # Filter by category slug
    &tag=ai,blockchain      # Filter by tags (comma-separated)
    &status=published       # Filter by status
    &author_id=123          # Filter by author
    &is_featured=true       # Only featured
    &page=1                 # Page number (default: 1)
    &limit=20               # Items per page (default: 20)
    &sort=published_at      # Sort field
    &order=desc             # Sort direction (asc/desc)
```

### Repository реализация

```go
type NewsFilter struct {
    CategorySlug string
    Tags         []string
    Status       string
    AuthorID     *int64
    IsFeatured   *bool
    IsBreaking   *bool
    Page         int
    Limit        int
    SortBy       string
    SortOrder    string
}

func (r *NewsRepository) List(ctx context.Context, filter NewsFilter) ([]*News, int64, error) {
    query := r.db.WithContext(ctx).Preload("Category").Preload("Tags")
    
    // Apply filters
    if filter.CategorySlug != "" {
        query = query.Joins("JOIN categories ON news.category_id = categories.id").
            Where("categories.slug = ?", filter.CategorySlug)
    }
    
    if len(filter.Tags) > 0 {
        query = query.Joins("JOIN news_tags ON news.id = news_tags.news_id").
            Joins("JOIN tags ON news_tags.tag_id = tags.id").
            Where("tags.slug IN ?", filter.Tags).
            Group("news.id").
            Having("COUNT(DISTINCT tags.id) = ?", len(filter.Tags))
    }
    
    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }
    
    if filter.IsFeatured != nil {
        query = query.Where("is_featured = ?", *filter.IsFeatured)
    }
    
    // Pagination
    offset := (filter.Page - 1) * filter.Limit
    
    var total int64
    query.Model(&News{}).Count(&total)
    
    var news []*News
    err := query.Offset(offset).Limit(filter.Limit).
        Order(fmt.Sprintf("%s %s", filter.SortBy, filter.SortOrder)).
        Find(&news).Error
    
    return news, total, err
}
```

---

## 🌐 Интеграция с SEO Service

### gRPC Client (будущая реализация)

```go
// internal/service/seo_service_client.go
type SEOServiceClient struct {
    conn   *grpc.ClientConn
    client pb.SEOServiceClient
}

func (s *SEOServiceClient) GenerateMetaTags(ctx context.Context, newsID int64, title, content string) (*SEOMetaTags, error) {
    resp, err := s.client.GenerateMetaTags(ctx, &pb.GenerateMetaTagsRequest{
        NewsId:  newsID,
        Title:   title,
        Content: content,
    })
    if err != nil {
        return nil, err
    }
    
    return &SEOMetaTags{
        Title:       resp.Title,
        Description: resp.Description,
        Keywords:    resp.Keywords,
    }, nil
}
```

### Вызов при публикации

```go
func (s *NewsService) PublishNews(ctx context.Context, id int64) error {
    news, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    
    // Генерация SEO метатегов (если не заданы вручную)
    if news.SEOTitle == "" || news.SEODescription == "" {
        seoMeta, err := s.seoClient.GenerateMetaTags(ctx, id, news.Title, news.Content)
        if err != nil {
            logger.Warn("Failed to generate SEO meta tags", zap.Error(err))
        } else {
            news.SEOTitle = seoMeta.Title
            news.SEODescription = seoMeta.Description
            news.SEOKeywords = seoMeta.Keywords
        }
    }
    
    news.Status = "published"
    news.PublishedAt = time.Now()
    
    return s.repo.Update(ctx, news)
}
```

---

## 📊 Статусы новостей

### Enum Status

```go
const (
    StatusDraft     = "draft"      // Черновик
    StatusPublished = "published"  // Опубликовано
    StatusArchived  = "archived"   // Архив
    StatusScheduled = "scheduled"  // Запланировано (будущее)
)
```

### Валидация

```go
type CreateNewsRequest struct {
    Title      string   `json:"title" binding:"required,min=3,max=500"`
    Slug       string   `json:"slug" binding:"omitempty,max=500"`
    Content    string   `json:"content" binding:"required"`
    Excerpt    string   `json:"excerpt" binding:"max=500"`
    CategoryID *int64   `json:"category_id"`
    TagIDs     []int64  `json:"tag_ids"`
    Status     string   `json:"status" binding:"required,oneof=draft published archived"`
    IsFeatured bool     `json:"is_featured"`
    IsBreaking bool     `json:"is_breaking"`
}
```

---

## 🚀 Текущее состояние (реализовано)

### ✅ Что работает

1. **Модели (GORM):**
   - ✅ News (с полями SEO)
   - ✅ Category (иерархические)
   - ✅ Tag
   - ✅ Many-to-Many связь news_tags

2. **Repository (DAO):**
   - ✅ CRUD операции для всех сущностей
   - ✅ Фильтрация (по категории, тегам, статусу)
   - ✅ Пагинация
   - ✅ Soft delete

3. **Service (бизнес-логика):**
   - ✅ NewsService с Redis кэшированием
   - ✅ CategoryService с иерархией
   - ✅ TagService
   - ✅ Генерация slug из title

4. **HTTP Handler:**
   - ✅ 18 REST endpoints
   - ✅ Валидация входных данных
   - ✅ Обработка ошибок
   - ✅ JSON responses

5. **Дополнительно:**
   - ✅ Health check endpoint
   - ✅ Zap структурированное логирование
   - ✅ Graceful shutdown
   - ✅ CORS middleware

### ⚠️ Проблемы (требуют исправления)

1. **Import proto пакета** - нужно удалить:
   ```go
   // Удалить эту строку из main.go:
   pb "news-service/proto"
   
   // Удалить gRPC сервер из main.go
   ```

2. **grpc_handler.go** - удалить файл или закомментировать

3. **После исправления:**
   - Пересобрать Docker образ
   - Запустить сервис
   - Протестировать API

---

## 🛠️ План исправления

### Шаг 1: Удалить gRPC части

```bash
# На локальном компьютере
cd "C:\Users\Grifindor\Desktop\НОВСТНОЙ САЙТ\news-service"

# Удалить или переименовать grpc_handler.go
mv internal/handler/grpc_handler.go internal/handler/grpc_handler.go.disabled

# Также для auth-service и media-service
```

### Шаг 2: Обновить main.go

Удалить импорты и код, связанный с gRPC.

### Шаг 3: Пересобрать

```bash
# Загрузить на сервер
scp -r news-service root@151.241.228.203:/opt/news-portal/

# Пересобрать
ssh root@151.241.228.203
cd /opt/news-portal
docker compose build news-service
docker compose up -d news-service
```

---

## 📋 Следующие этапы (после исправления)

1. **SEO Service** - создать отдельный микросервис
2. **gRPC** - добавить proto файлы и реализацию
3. **Full-Text Search** - настроить PostgreSQL FTS
4. **Scheduled Posts** - публикация по расписанию
5. **Версионирование** - история изменений новостей
6. **Комментарии** - интеграция с comments-service
7. **Analytics** - счётчики просмотров, популярность

---

## 🔗 API Примеры

### Создать новость

```bash
curl -X POST http://151.241.228.203:8092/api/v1/news \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer JWT_TOKEN" \
  -d '{
    "title": "Новая технология AI",
    "content": "Подробное описание...",
    "excerpt": "Краткое описание",
    "category_id": 1,
    "tag_ids": [1, 2, 3],
    "status": "published",
    "is_featured": true
  }'
```

### Получить список новостей

```bash
curl "http://151.241.228.203:8092/api/v1/news?category=tech&page=1&limit=10"
```

### Поиск

```bash
curl "http://151.241.228.203:8092/api/v1/news/search?q=искусственный+интеллект"
```

---

**Статус:** Готов к исправлению импортов и повторной сборке! 🚀
