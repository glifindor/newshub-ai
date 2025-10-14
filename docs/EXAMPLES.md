# 🧪 Примеры использования API

Коллекция примеров запросов к API новостного портала.

---

## 🔐 Аутентификация

### Регистрация нового пользователя

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!",
    "full_name": "John Doe",
    "role": "user"
  }'
```

**Ответ:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Вход в систему

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

**Ответ:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900
}
```

### Обновление токена

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

### Выход из системы

```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

---

## 📰 Новости

### Получить список новостей

```bash
# Все опубликованные новости
curl http://localhost:8080/api/news

# С пагинацией
curl "http://localhost:8080/api/news?page=1&page_size=10"

# Фильтр по категории
curl "http://localhost:8080/api/news?category_id=550e8400-e29b-41d4-a716-446655440000"

# Поиск
curl "http://localhost:8080/api/news?search=технологии"

# Только черновики (требуется авторизация)
curl "http://localhost:8080/api/news?status=draft" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Ответ:**
```json
{
  "news": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Новый прорыв в AI технологиях",
      "slug": "novyy-proryv-v-ai-tehnologiyah",
      "content": "Полный текст статьи...",
      "summary": "Краткое описание...",
      "author_id": "660e8400-e29b-41d4-a716-446655440001",
      "category_id": "770e8400-e29b-41d4-a716-446655440002",
      "tag_ids": ["880e8400-e29b-41d4-a716-446655440003"],
      "featured_image": "https://cdn.example.com/ai-breakthrough.jpg",
      "status": "published",
      "views_count": 1523,
      "published_at": "2025-10-14T10:00:00Z",
      "created_at": "2025-10-14T09:00:00Z",
      "updated_at": "2025-10-14T09:30:00Z"
    }
  ],
  "total": 150,
  "page": 1,
  "page_size": 10
}
```

### Получить новость по slug

```bash
curl http://localhost:8080/api/news/novyy-proryv-v-ai-tehnologiyah
```

### Создать новость (требуется роль editor/admin)

```bash
curl -X POST http://localhost:8080/api/news \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Квантовые компьютеры в 2025 году",
    "content": "Полный текст статьи о квантовых компьютерах...",
    "summary": "Обзор последних достижений в области квантовых вычислений",
    "category_id": "770e8400-e29b-41d4-a716-446655440002",
    "tag_ids": ["880e8400-e29b-41d4-a716-446655440003", "880e8400-e29b-41d4-a716-446655440004"],
    "featured_image": "https://cdn.example.com/quantum.jpg"
  }'
```

**Ответ:**
```json
{
  "news": {
    "id": "990e8400-e29b-41d4-a716-446655440005",
    "title": "Квантовые компьютеры в 2025 году",
    "slug": "kvantovye-kompyutery-v-2025-godu",
    "status": "draft",
    "created_at": "2025-10-14T11:00:00Z",
    ...
  }
}
```

### Обновить новость

```bash
curl -X PUT http://localhost:8080/api/news/990e8400-e29b-41d4-a716-446655440005 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Квантовые компьютеры в 2025: новая эра",
    "content": "Обновленный текст...",
    "summary": "Обновленное описание...",
    "category_id": "770e8400-e29b-41d4-a716-446655440002",
    "tag_ids": ["880e8400-e29b-41d4-a716-446655440003"],
    "featured_image": "https://cdn.example.com/quantum-new.jpg"
  }'
```

### Опубликовать новость

```bash
curl -X POST http://localhost:8080/api/news/990e8400-e29b-41d4-a716-446655440005/publish \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Удалить новость (требуется роль admin)

```bash
curl -X DELETE http://localhost:8080/api/news/990e8400-e29b-41d4-a716-446655440005 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## 📂 Категории

### Получить все категории

```bash
curl http://localhost:8080/api/categories
```

**Ответ:**
```json
{
  "categories": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "name": "Технологии",
      "slug": "tehnologii",
      "description": "Новости о технологиях и инновациях",
      "created_at": "2025-10-01T10:00:00Z"
    },
    {
      "id": "770e8400-e29b-41d4-a716-446655440003",
      "name": "Наука",
      "slug": "nauka",
      "description": "Научные открытия и исследования",
      "created_at": "2025-10-01T10:00:00Z"
    }
  ],
  "total": 10
}
```

### Создать категорию (требуется роль admin)

```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Искусственный интеллект",
    "description": "Все о AI и машинном обучении"
  }'
```

---

## 🏷️ Теги

### Получить все теги

```bash
curl http://localhost:8080/api/tags
```

**Ответ:**
```json
{
  "tags": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "name": "AI",
      "slug": "ai",
      "created_at": "2025-10-01T10:00:00Z"
    },
    {
      "id": "880e8400-e29b-41d4-a716-446655440004",
      "name": "Квантовые компьютеры",
      "slug": "kvantovye-kompyutery",
      "created_at": "2025-10-01T10:00:00Z"
    }
  ],
  "total": 50
}
```

### Создать тег (требуется роль admin)

```bash
curl -X POST http://localhost:8080/api/tags \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Блокчейн"
  }'
```

---

## 🖼️ Медиа

### Загрузить изображение

```bash
curl -X POST http://localhost:8080/api/media/upload \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -F "file=@/path/to/image.jpg"
```

**Ответ:**
```json
{
  "url": "https://cdn.example.com/uploads/550e8400-e29b-41d4-a716-446655440000.jpg",
  "filename": "image.jpg",
  "size": 204800,
  "mime_type": "image/jpeg"
}
```

### Загрузить несколько файлов

```bash
curl -X POST http://localhost:8080/api/media/upload \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -F "file=@/path/to/image1.jpg" \
  -F "file=@/path/to/image2.png"
```

---

## 👥 Администрирование

### Получить список пользователей (admin/moderator)

```bash
curl http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Ответ:**
```json
{
  "users": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role": "editor",
      "created_at": "2025-10-10T10:00:00Z"
    }
  ],
  "total": 25
}
```

### Получить статистику

```bash
curl http://localhost:8080/api/admin/statistics \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Ответ:**
```json
{
  "total_news": 1500,
  "published_news": 1200,
  "draft_news": 300,
  "total_users": 500,
  "total_views": 1000000,
  "popular_news": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Новый прорыв в AI технологиях",
      "views_count": 15000
    }
  ],
  "recent_registrations": 45,
  "today_views": 5000
}
```

### Модерация контента

```bash
curl -X POST http://localhost:8080/api/admin/moderate/990e8400-e29b-41d4-a716-446655440005 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "approve",
    "comment": "Content approved"
  }'
```

---

## 🔍 SEO

### Получить sitemap.xml

```bash
curl http://localhost:8080/api/sitemap.xml
```

**Ответ:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://newsportal.com/news/novyy-proryv-v-ai-tehnologiyah</loc>
    <lastmod>2025-10-14T10:00:00Z</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
</urlset>
```

### Получить robots.txt

```bash
curl http://localhost:8080/api/robots.txt
```

**Ответ:**
```
User-agent: *
Allow: /

Sitemap: https://newsportal.com/api/sitemap.xml
```

---

## 🧪 Примеры с JavaScript/TypeScript

### Fetch API

```javascript
// Вход
const login = async (email, password) => {
  const response = await fetch('http://localhost:8080/api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });
  
  const data = await response.json();
  localStorage.setItem('access_token', data.access_token);
  localStorage.setItem('refresh_token', data.refresh_token);
  return data;
};

// Получение новостей
const getNews = async (page = 1) => {
  const response = await fetch(`http://localhost:8080/api/news?page=${page}&page_size=20`);
  return await response.json();
};

// Создание новости
const createNews = async (newsData) => {
  const token = localStorage.getItem('access_token');
  
  const response = await fetch('http://localhost:8080/api/news', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(newsData),
  });
  
  return await response.json();
};

// Загрузка изображения
const uploadImage = async (file) => {
  const token = localStorage.getItem('access_token');
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await fetch('http://localhost:8080/api/media/upload', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: formData,
  });
  
  return await response.json();
};
```

### Axios

```javascript
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
});

// Interceptor для добавления токена
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Использование
const fetchNews = async () => {
  const response = await api.get('/news');
  return response.data;
};

const createNews = async (newsData) => {
  const response = await api.post('/news', newsData);
  return response.data;
};
```

---

## 🐍 Примеры с Python

```python
import requests

BASE_URL = "http://localhost:8080/api"

# Вход
def login(email, password):
    response = requests.post(
        f"{BASE_URL}/auth/login",
        json={"email": email, "password": password}
    )
    data = response.json()
    return data['access_token']

# Получение новостей
def get_news(page=1, page_size=20):
    response = requests.get(
        f"{BASE_URL}/news",
        params={"page": page, "page_size": page_size}
    )
    return response.json()

# Создание новости
def create_news(token, news_data):
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.post(
        f"{BASE_URL}/news",
        json=news_data,
        headers=headers
    )
    return response.json()

# Загрузка изображения
def upload_image(token, file_path):
    headers = {"Authorization": f"Bearer {token}"}
    with open(file_path, 'rb') as f:
        files = {'file': f}
        response = requests.post(
            f"{BASE_URL}/media/upload",
            files=files,
            headers=headers
        )
    return response.json()

# Использование
if __name__ == "__main__":
    token = login("admin@example.com", "admin123")
    news = get_news(page=1)
    print(f"Total news: {news['total']}")
```

---

## ⚠️ Обработка ошибок

### Пример обработки ошибок в JavaScript

```javascript
const fetchNewsWithErrorHandling = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/news');
    
    if (!response.ok) {
      if (response.status === 401) {
        // Токен истек, попробовать обновить
        await refreshToken();
        return fetchNewsWithErrorHandling(); // Retry
      }
      
      if (response.status === 429) {
        throw new Error('Too many requests. Please try again later.');
      }
      
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching news:', error);
    throw error;
  }
};

const refreshToken = async () => {
  const refreshToken = localStorage.getItem('refresh_token');
  
  const response = await fetch('http://localhost:8080/api/auth/refresh', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });
  
  const data = await response.json();
  localStorage.setItem('access_token', data.access_token);
  localStorage.setItem('refresh_token', data.refresh_token);
};
```

---

## 📝 Коды ответов

| Код | Значение | Описание |
|-----|----------|----------|
| 200 | OK | Успешный запрос |
| 201 | Created | Ресурс создан |
| 400 | Bad Request | Неверные данные |
| 401 | Unauthorized | Требуется аутентификация |
| 403 | Forbidden | Недостаточно прав |
| 404 | Not Found | Ресурс не найден |
| 429 | Too Many Requests | Превышен лимит запросов |
| 500 | Internal Server Error | Ошибка сервера |

---

**Дата обновления:** 2025-10-14
