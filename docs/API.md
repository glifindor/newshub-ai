# News Portal API Documentation

## Базовый URL
```
http://localhost:8080/api
```

## Аутентификация

### Регистрация
**POST** `/auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "full_name": "John Doe",
  "role": "user"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "user_id": "uuid"
}
```

### Вход
**POST** `/auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900
}
```

### Обновление токена
**POST** `/auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900
}
```

### Выход
**POST** `/auth/logout` 🔒

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

---

## Новости

### Список новостей
**GET** `/news`

**Query Parameters:**
- `page` (int) - номер страницы (default: 1)
- `page_size` (int) - размер страницы (default: 20)
- `category_id` (string) - фильтр по категории
- `status` (string) - фильтр по статусу (draft, published, archived)
- `search` (string) - поиск по заголовку и содержимому

**Response:**
```json
{
  "news": [
    {
      "id": "uuid",
      "title": "Breaking News",
      "slug": "breaking-news",
      "content": "Full article content...",
      "summary": "Short summary...",
      "author_id": "uuid",
      "category_id": "uuid",
      "tag_ids": ["uuid1", "uuid2"],
      "featured_image": "https://cdn.example.com/image.jpg",
      "status": "published",
      "views_count": 1500,
      "published_at": "2025-10-14T10:00:00Z",
      "created_at": "2025-10-14T09:00:00Z",
      "updated_at": "2025-10-14T09:30:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

### Получить новость по slug
**GET** `/news/{slug}`

**Response:**
```json
{
  "news": {
    "id": "uuid",
    "title": "Breaking News",
    "slug": "breaking-news",
    "content": "Full article content...",
    "summary": "Short summary...",
    "author_id": "uuid",
    "category_id": "uuid",
    "tag_ids": ["uuid1", "uuid2"],
    "featured_image": "https://cdn.example.com/image.jpg",
    "status": "published",
    "views_count": 1500,
    "published_at": "2025-10-14T10:00:00Z",
    "created_at": "2025-10-14T09:00:00Z",
    "updated_at": "2025-10-14T09:30:00Z"
  }
}
```

### Создать новость 🔒
**POST** `/news`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Required Role:** `editor`, `admin`

**Request Body:**
```json
{
  "title": "New Article Title",
  "content": "Full article content...",
  "summary": "Short summary...",
  "category_id": "uuid",
  "tag_ids": ["uuid1", "uuid2"],
  "featured_image": "https://cdn.example.com/image.jpg"
}
```

**Response:**
```json
{
  "news": {
    "id": "uuid",
    "title": "New Article Title",
    "slug": "new-article-title",
    "status": "draft",
    ...
  }
}
```

### Обновить новость 🔒
**PUT** `/news/{id}`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Required Role:** `editor`, `admin`

**Request Body:**
```json
{
  "title": "Updated Title",
  "content": "Updated content...",
  "summary": "Updated summary...",
  "category_id": "uuid",
  "tag_ids": ["uuid1", "uuid2"],
  "featured_image": "https://cdn.example.com/image.jpg"
}
```

### Публиковать новость 🔒
**POST** `/news/{id}/publish`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Required Role:** `editor`, `admin`

**Response:**
```json
{
  "news": {
    "id": "uuid",
    "status": "published",
    "published_at": "2025-10-14T10:00:00Z",
    ...
  }
}
```

### Удалить новость 🔒
**DELETE** `/news/{id}`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Required Role:** `admin`

**Response:**
```json
{
  "success": true,
  "message": "News deleted successfully"
}
```

---

## Категории

### Список категорий
**GET** `/categories`

**Query Parameters:**
- `page` (int)
- `page_size` (int)

**Response:**
```json
{
  "categories": [
    {
      "id": "uuid",
      "name": "Technology",
      "slug": "technology",
      "description": "Tech news and articles",
      "created_at": "2025-10-14T09:00:00Z"
    }
  ],
  "total": 10
}
```

### Создать категорию 🔒
**POST** `/categories`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Required Role:** `admin`

**Request Body:**
```json
{
  "name": "Technology",
  "description": "Tech news and articles"
}
```

---

## Теги

### Список тегов
**GET** `/tags`

**Query Parameters:**
- `page` (int)
- `page_size` (int)

**Response:**
```json
{
  "tags": [
    {
      "id": "uuid",
      "name": "AI",
      "slug": "ai",
      "created_at": "2025-10-14T09:00:00Z"
    }
  ],
  "total": 50
}
```

### Создать тег 🔒
**POST** `/tags`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Required Role:** `admin`

**Request Body:**
```json
{
  "name": "Artificial Intelligence"
}
```

---

## Медиа

### Загрузить файл 🔒
**POST** `/media/upload`

**Headers:**
```
Authorization: Bearer {access_token}
Content-Type: multipart/form-data
```

**Request Body:**
```
file: [binary data]
```

**Response:**
```json
{
  "url": "https://cdn.example.com/uploads/image-uuid.jpg",
  "filename": "image.jpg",
  "size": 102400,
  "mime_type": "image/jpeg"
}
```

---

## Админ-панель

### Список пользователей 🔒
**GET** `/admin/users`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Required Role:** `admin`, `moderator`

**Response:**
```json
{
  "users": [
    {
      "id": "uuid",
      "email": "user@example.com",
      "full_name": "John Doe",
      "role": "editor",
      "created_at": "2025-10-14T09:00:00Z"
    }
  ],
  "total": 50
}
```

### Статистика 🔒
**GET** `/admin/statistics`

**Headers:**
```
Authorization: Bearer {access_token}
```

**Required Role:** `admin`, `moderator`

**Response:**
```json
{
  "total_news": 1500,
  "published_news": 1200,
  "draft_news": 300,
  "total_users": 500,
  "total_views": 1000000,
  "popular_news": [...]
}
```

---

## SEO

### Получить sitemap.xml
**GET** `/sitemap.xml`

**Response:** XML sitemap

### Получить robots.txt
**GET** `/robots.txt`

**Response:** Plain text robots.txt

---

## Коды ответов

- **200** - Успешный запрос
- **201** - Ресурс создан
- **400** - Неверный запрос
- **401** - Не авторизован
- **403** - Доступ запрещен
- **404** - Ресурс не найден
- **429** - Слишком много запросов
- **500** - Внутренняя ошибка сервера

## Rate Limiting

- **Лимит:** 100 запросов в минуту на IP
- **Header:** `X-RateLimit-Remaining`

---

🔒 - Требуется авторизация (Bearer Token)
