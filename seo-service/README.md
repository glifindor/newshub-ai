# SEO-Service - Микросервис SEO для новостного портала

Микросервис для управления SEO метаданными, генерации sitemap.xml, robots.txt и Open Graph тегов для российских соцсетей (ВКонтакте, Telegram, Одноклассники).

## 🎯 Основные возможности

- ✅ **Автоматическая генерация SEO** из новостей (через webhook)
- ✅ **Sitemap.xml** с кэшированием в Redis
- ✅ **Robots.txt** с поддержкой Яндекс, Google, Mail.ru
- ✅ **Open Graph** для VK, Telegram, OK
- ✅ **Schema.org JSON-LD** для структурированных данных
- ✅ **Full CRUD API** для управления SEO метаданными
- ✅ **Redis кэширование** (TTL: 1 час для sitemap, 24 часа для robots)

## 🛠 Технологии

- **Go 1.23**
- **Gin** - HTTP фреймворк
- **GORM** - ORM для PostgreSQL
- **PostgreSQL 15** - База данных (таблица `seo_meta`)
- **Redis 7** - Кэширование
- **Docker** - Контейнеризация

## 📁 Структура проекта

```
seo-service/
├── cmd/seo-service/          # Точка входа
│   └── main.go
├── internal/
│   ├── config/               # Конфигурация
│   ├── models/               # Модели данных
│   ├── repository/           # Слой данных
│   ├── service/              # Бизнес-логика
│   │   ├── seo_service.go
│   │   ├── opengraph_service.go
│   │   ├── sitemap_service.go
│   │   └── robots_service.go
│   └── handler/              # HTTP обработчики
├── pkg/
│   ├── database/             # PostgreSQL подключение
│   ├── logger/               # Zap logger
│   └── generator/            # Генераторы XML/JSON
│       ├── sitemap.go
│       ├── robots.go
│       └── structured_data.go
├── migrations/               # SQL миграции
│   └── 001_create_seo_meta.sql
├── .env                      # Конфигурация
├── Dockerfile
└── go.mod
```

## 🚀 Установка и запуск

### 1. Локальная разработка

```bash
# Клонировать репозиторий
cd seo-service

# Установить зависимости
go mod download

# Настроить .env
cp .env.example .env
# Отредактировать .env

# Запустить миграции
psql -h localhost -U postgres -d newsportal_db -f migrations/001_create_seo_meta.sql

# Запустить сервис
go run cmd/seo-service/main.go
```

### 2. Docker

```bash
# Собрать образ
docker build -t seo-service:latest .

# Запустить через docker-compose
docker-compose up seo-service
```

## 🌐 API Endpoints

### SEO Метаданные

#### 1. Получить SEO по slug
```http
GET /api/v1/seo/:slug
```

**Пример ответа:**
```json
{
  "id": "uuid",
  "news_id": "uuid",
  "slug": "breaking-news-today",
  "title": "Главная новость дня - Новостной портал",
  "description": "Краткое описание новости для поисковиков...",
  "keywords": "новость, политика, экономика",
  "canonical_url": "http://example.com/news/breaking-news-today",
  "og_title": "Главная новость дня",
  "og_description": "Краткое описание...",
  "og_image": "http://example.com/images/news.jpg",
  "og_type": "article",
  "og_locale": "ru_RU",
  "robots_index": true,
  "robots_follow": true,
  "sitemap_change_freq": "daily",
  "sitemap_priority": 0.8
}
```

#### 2. Создать SEO метаданные
```http
POST /api/v1/seo/create
Content-Type: application/json

{
  "news_id": "uuid",
  "slug": "my-news",
  "title": "SEO заголовок",
  "description": "SEO описание",
  "keywords": "ключ1, ключ2",
  "og_title": "OG заголовок",
  "og_description": "OG описание",
  "og_image": "http://example.com/image.jpg",
  "robots_index": true,
  "robots_follow": true
}
```

#### 3. Обновить SEO
```http
PUT /api/v1/seo/update
Content-Type: application/json
```

#### 4. Удалить SEO
```http
DELETE /api/v1/seo/:newsId
```

### Webhook

#### 5. Webhook при публикации новости
```http
POST /api/v1/webhook/news-published
Content-Type: application/json

{
  "news_id": "uuid",
  "slug": "breaking-news",
  "title": "Заголовок новости",
  "summary": "Краткое содержание",
  "content": "Полный текст новости...",
  "author_name": "Иван Иванов",
  "image_url": "http://example.com/news.jpg",
  "published_at": "2025-10-14T10:00:00Z"
}
```

**Что происходит:**
- Автоматически генерируется SEO title (≤70 символов)
- Автоматически генерируется SEO description (≤160 символов)
- Извлекаются ключевые слова (частотный анализ с фильтрацией стоп-слов)
- Создаются Open Graph теги для VK/Telegram/OK
- Генерируется Schema.org JSON-LD
- Инвалидируется кэш sitemap

### Open Graph

#### 6. Получить OG теги
```http
GET /api/v1/seo/:slug/og-tags
```

**Пример ответа:**
```json
{
  "slug": "breaking-news",
  "og_tags": {
    "og:title": "Главная новость дня",
    "og:description": "Краткое описание...",
    "og:type": "article",
    "og:url": "http://example.com/news/breaking-news",
    "og:image": "http://example.com/images/news.jpg",
    "og:image:width": "1200",
    "og:image:height": "630",
    "og:locale": "ru_RU",
    "og:site_name": "Новостной портал",
    "article:published_time": "2025-10-14T10:00:00+03:00",
    "article:modified_time": "2025-10-14T12:00:00+03:00"
  }
}
```

### Публичные endpoints

#### 7. Sitemap.xml
```http
GET /sitemap.xml
```

Возвращает XML sitemap со всеми индексируемыми новостями.

#### 8. Robots.txt
```http
GET /robots.txt
```

Возвращает robots.txt с настройками для Яндекс, Google, Mail.ru.

#### 9. Health Check
```http
GET /health
```

```json
{
  "status": "ok",
  "service": "seo-service",
  "version": "1.0.0"
}
```

## 🗄️ База данных

### Таблица `seo_meta`

Содержит 24 поля:
- **Основные:** id, news_id, slug, title, description, keywords, canonical_url
- **Open Graph:** og_title, og_description, og_image, og_type, og_locale, og_site_name
- **Robots:** robots_index, robots_follow
- **Sitemap:** sitemap_change_freq, sitemap_priority
- **Schema.org:** schema_data (JSONB)
- **Timestamps:** created_at, updated_at

### Индексы (5 шт)

1. `idx_seo_meta_slug` (UNIQUE) - Быстрый поиск по URL
2. `idx_seo_meta_news_id` (UNIQUE) - Связь с новостью
3. `idx_seo_meta_updated_at` - Сортировка по дате
4. `idx_seo_meta_sitemap` (WHERE robots_index=true) - Оптимизация sitemap
5. `idx_seo_meta_schema_data` (GIN) - Поиск по JSON-LD

### Триггер

`update_seo_meta_updated_at` - автоматически обновляет `updated_at` при изменении записи

## 🔄 Интеграция с news-service

При публикации новости в `news-service` отправляется webhook:

```go
// В news-service добавить после создания новости
func (s *newsService) PublishNews(news *News) error {
    // ... сохранение новости ...
    
    // Отправляем webhook в seo-service
    webhook := SEOWebhook{
        NewsID:      news.ID,
        Slug:        news.Slug,
        Title:       news.Title,
        Summary:     news.Summary,
        Content:     news.Content,
        AuthorName:  news.Author.Name,
        ImageURL:    news.ImageURL,
        PublishedAt: news.PublishedAt,
    }
    
    resp, err := http.Post(
        "http://seo-service:8093/api/v1/webhook/news-published",
        "application/json",
        toJSON(webhook),
    )
    
    return err
}
```

## 📝 Конфигурация (.env)

```env
# Окружение
ENVIRONMENT=production
SERVER_PORT=8093

# PostgreSQL
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=newsportal_db

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0

# SEO
BASE_URL=http://151.241.228.203
SITE_NAME=Новостной портал
NEWS_SERVICE_URL=http://news-service:8082

# Логирование
LOG_LEVEL=info
```

## 🎨 Особенности для России

### 1. Поисковики
- **Яндекс** - приоритетная оптимизация (Crawl-delay, Clean-param)
- **Google** - стандартная SEO
- **Mail.ru** - поддержка Mail.RU_Bot

### 2. Социальные сети
- **ВКонтакте** - Open Graph с `vk:image`
- **Telegram** - OG для Instant View
- **Одноклассники** - стандартный Open Graph

### 3. Локализация
- `og:locale` = `ru_RU`
- Timezone = `Europe/Moscow`
- Стоп-слова на русском и английском

## 🧪 Тестирование

### Проверка работоспособности

```bash
# Health check
curl http://localhost:8093/health

# Sitemap
curl http://localhost:8093/sitemap.xml

# Robots
curl http://localhost:8093/robots.txt

# Создать SEO
curl -X POST http://localhost:8093/api/v1/seo/create \
  -H "Content-Type: application/json" \
  -d '{
    "news_id": "uuid-здесь",
    "slug": "test-news",
    "title": "Тестовая новость",
    "description": "Описание тестовой новости"
  }'

# Получить SEO
curl http://localhost:8093/api/v1/seo/test-news

# OG теги
curl http://localhost:8093/api/v1/seo/test-news/og-tags
```

## 📊 Производительность

- **Sitemap кэш:** 1 час (3600s)
- **Robots кэш:** 24 часа (86400s)
- **DB запросы:** оптимизированы индексами
- **Генерация SEO:** ~5-10ms
- **Sitemap generation:** ~50-100ms (для 1000 новостей)

## 🔧 Deployment на сервер

```bash
# 1. Подключиться к серверу
ssh root@151.241.228.203

# 2. Перейти в директорию проекта
cd /path/to/project

# 3. Скопировать seo-service
scp -r seo-service root@151.241.228.203:/path/to/project/

# 4. Применить миграцию
docker exec -i news-postgres psql -U postgres -d newsportal_db < seo-service/migrations/001_create_seo_meta.sql

# 5. Пересобрать и запустить
docker-compose up -d --build seo-service

# 6. Проверить логи
docker logs -f seo-service

# 7. Проверить работу
curl http://151.241.228.203:8093/health
curl http://151.241.228.203:8093/sitemap.xml
curl http://151.241.228.203:8093/robots.txt
```

## 📚 Дополнительная документация

- [SEO_SERVICE_PLAN.md](./SEO_SERVICE_PLAN.md) - План разработки (15 шагов)
- [PROGRESS.md](./PROGRESS.md) - Статус выполнения
- [migrations/001_create_seo_meta.sql](./migrations/001_create_seo_meta.sql) - SQL схема

## 🤝 Вклад

Проект разработан для российского новостного портала с фокусом на:
- Яндекс SEO
- Open Graph для VK/Telegram/OK
- Schema.org для структурированных данных
- Автоматизацию SEO генерации

## 📄 Лицензия

MIT License

---

**Версия:** 1.0.0  
**Дата:** 14 октября 2025  
**Разработчик:** GitHub Copilot
