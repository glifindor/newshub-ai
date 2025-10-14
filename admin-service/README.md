# Admin Service - Административная панель

Микросервис для управления новостным сайтом с React UI и TipTap WYSIWYG редактором.

## 🚀 Функциональность

### Backend (Go + Gin)
- ✅ Авторизация через JWT (интеграция с auth-service)
- ✅ RBAC - доступ только для ролей `admin` и `editor`
- ✅ Управление новостями (CRUD)
- ✅ Управление SEO-метаданными
- ✅ Дашборд со статистикой (Prometheus метрики)
- ✅ Фильтрация по категориям и тегам
- ✅ CORS middleware для SPA

### Frontend (React + TypeScript + Vite)
- ✅ **TipTap WYSIWYG редактор** для создания контента
- ✅ React Router с защищенными маршрутами
- ✅ Zustand для управления состоянием (persist в localStorage)
- ✅ Tailwind CSS для стилизации
- ✅ Recharts для визуализации статистики
- ✅ SEO-форма с лимитами символов (title ≤70, description ≤160)
- ✅ Автоматическая генерация slug из заголовка

## 📁 Структура проекта

```
admin-service/
├── cmd/admin-service/
│   └── main.go                 # Точка входа
├── internal/
│   ├── config/
│   │   └── config.go           # Конфигурация
│   ├── handler/
│   │   └── admin_handler.go    # HTTP handlers (15 endpoints)
│   ├── middleware/
│   │   └── auth.go             # Auth, RBAC, CORS
│   ├── client/
│   │   ├── auth_client.go      # HTTP клиент для auth-service
│   │   ├── news_client.go      # HTTP клиент для news-service
│   │   └── seo_client.go       # HTTP клиент для seo-service
│   └── models/
│       └── models.go           # Data models (15+ structs)
├── pkg/logger/
│   └── logger.go               # Zap logger
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Layout.tsx      # Main layout с навигацией
│   │   │   ├── TipTapEditor.tsx# WYSIWYG редактор
│   │   │   └── SEOForm.tsx     # SEO метаданные
│   │   ├── pages/
│   │   │   ├── LoginPage.tsx   # Страница авторизации
│   │   │   ├── Dashboard.tsx   # Статистика + графики
│   │   │   ├── NewsList.tsx    # Список новостей + фильтры
│   │   │   └── NewsEdit.tsx    # Редактор новости
│   │   ├── api/
│   │   │   └── api.ts          # Axios клиент (200+ lines)
│   │   ├── store/
│   │   │   └── authStore.ts    # Zustand auth store
│   │   ├── App.tsx             # Router + Protected Routes
│   │   ├── main.tsx            # React entry point
│   │   └── index.css           # TipTap styles + Tailwind
│   ├── package.json
│   ├── vite.config.ts          # Dev server + proxy
│   ├── tailwind.config.js
│   └── tsconfig.json
├── .env                        # Переменные окружения
├── go.mod
└── README.md
```

## ⚙️ Установка

### Требования
- **Go 1.23+**
- **Node.js 20+**
- Запущенные сервисы: `auth-service`, `news-service`, `seo-service`

### 1. Backend

```bash
# Перейти в директорию
cd admin-service

# Установить зависимости Go
go mod tidy

# Настроить .env (скопировать .env.example)
cp .env.example .env
```

**Конфигурация .env:**
```env
# Server
SERVER_PORT=8084
SERVER_ENV=development

# Services
AUTH_SERVICE_URL=http://localhost:8081
NEWS_SERVICE_URL=http://localhost:8082
SEO_SERVICE_URL=http://localhost:8093
MEDIA_SERVICE_URL=http://localhost:8094

# JWT
JWT_SECRET=your-secret-key

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:8084

# Pagination
DEFAULT_PAGE_SIZE=20
MAX_PAGE_SIZE=100

# Upload
MAX_FILE_SIZE=10485760
ALLOWED_FILE_TYPES=image/jpeg,image/png,image/webp
```

### 2. Frontend

```bash
# Перейти в директорию frontend
cd frontend

# Установить npm зависимости
npm install

# Собрать для production
npm run build
```

## 🚀 Запуск

### Development

**Backend + Frontend одновременно:**

1. **Терминал 1** - Backend:
```bash
cd admin-service
go run cmd/admin-service/main.go
# Server слушает на :8084
```

2. **Терминал 2** - Frontend dev server:
```bash
cd admin-service/frontend
npm run dev
# Dev server на :5173 с proxy /api -> :8084
```

3. Открыть http://localhost:5173

### Production

**Backend обслуживает статику:**

```bash
# Собрать frontend
cd admin-service/frontend
npm run build

# Запустить backend (будет раздавать dist/)
cd ..
go run cmd/admin-service/main.go

# Открыть http://localhost:8084
```

## 📡 API Endpoints

### Авторизация
- `POST /api/admin/login` - Вход (username, password)
- `GET /api/admin/me` - Текущий пользователь

### Дашборд
- `GET /api/admin/dashboard` - Статистика (total_news, views, popular_news, recent_news, views_trend)

### Новости
- `GET /api/admin/news` - Список (фильтры: search, status, category, tag, page, page_size)
- `GET /api/admin/news/:id` - Одна новость по ID
- `POST /api/admin/news` - Создать (requires auth + admin/editor role)
- `PUT /api/admin/news/:id` - Обновить (requires auth + admin/editor role)
- `DELETE /api/admin/news/:id` - Удалить (requires auth + admin role)

### SEO
- `GET /api/admin/news/:id/seo` - SEO по ID новости
- `PUT /api/admin/news/:id/seo` - Обновить SEO (requires auth + admin/editor role)

### Health
- `GET /health` - Health check

## 🔐 Авторизация

**Логин:**
```bash
curl -X POST http://localhost:8084/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

**Ответ:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin"
  }
}
```

**Использование токена:**
```bash
curl http://localhost:8084/api/admin/news \
  -H "Authorization: Bearer <token>"
```

## 🎨 Компоненты Frontend

### 1. **LoginPage** (`/login`)
- Форма авторизации (username, password)
- Сохранение JWT в localStorage через Zustand
- Редирект на Dashboard после успеха

### 2. **Dashboard** (`/`)
- 4 карточки статистики (всего новостей, опубликовано, черновики, просмотры)
- График динамики просмотров (Recharts LineChart)
- Топ-5 популярных новостей
- Топ-5 последних новостей

### 3. **NewsList** (`/news`)
- Таблица с пагинацией (20 элементов)
- Фильтры: поиск, статус (draft/published/archived)
- Действия: Редактировать, Удалить
- Кнопка "Создать новость"

### 4. **NewsEdit** (`/news/new`, `/news/:id`)
- Поля: Заголовок, Slug, Краткое описание
- **TipTap WYSIWYG редактор** для содержания
  * Форматирование: Bold, Italic, Strikethrough
  * Заголовки: H1, H2, H3
  * Списки: маркированные, нумерованные
  * Blockquote, Code Block
  * Undo/Redo
- Категория (выбор из списка)
- Теги (добавление/удаление)
- Обложка (URL изображения + preview)
- **SEO-форма** (раскрывающаяся секция):
  * Meta Title (≤70 символов)
  * Meta Description (≤160 символов)
  * Meta Keywords
  * Open Graph: Title, Description, Image
- Кнопки: "Сохранить черновик", "Опубликовать"

### 5. **TipTapEditor**
- StarterKit extension (базовый набор)
- Toolbar с 15+ кнопками
- Автосохранение контента в родительский state
- Стили ProseMirror в `index.css`

### 6. **SEOForm**
- Загрузка текущих SEO-данных из seo-service
- Счетчики символов для title/description
- Preview изображения OG:image
- Отдельная кнопка "Сохранить SEO"

## 🐳 Docker

```dockerfile
# Dockerfile (пример)
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o admin-service ./cmd/admin-service

FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /app/admin-service .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
COPY .env .
EXPOSE 8084
CMD ["./admin-service"]
```

**Docker Compose:**
```yaml
admin-service:
  build: ./admin-service
  ports:
    - "8084:8084"
  environment:
    SERVER_PORT: "8084"
    AUTH_SERVICE_URL: "http://auth-service:8081"
    NEWS_SERVICE_URL: "http://news-service:8082"
    SEO_SERVICE_URL: "http://seo-service:8093"
    MEDIA_SERVICE_URL: "http://media-service:8094"
  depends_on:
    - auth-service
    - news-service
    - seo-service
  networks:
    - news-network
```

## 🧪 Тестирование

### Backend
```bash
# Health check
curl http://localhost:8084/health

# Dashboard stats
curl http://localhost:8084/api/admin/dashboard \
  -H "Authorization: Bearer <token>"
```

### Frontend
```bash
# Dev server
npm run dev

# Build
npm run build

# Preview production build
npm run preview
```

## 🔧 Технологии

### Backend
- **Go 1.23** - язык программирования
- **Gin** - HTTP framework
- **Zap** - структурированное логирование
- **godotenv** - загрузка .env
- **uuid** - генерация UUID

### Frontend
- **React 18** - UI библиотека
- **TypeScript 5** - типизация
- **Vite 5** - сборщик + dev server
- **React Router 6** - маршрутизация
- **TailwindCSS 3** - стилизация
- **@tiptap/react** - WYSIWYG редактор
- **Zustand 4** - state management
- **Axios 1** - HTTP клиент
- **Recharts 2** - графики

## 📝 Примеры использования

### Создание новости через UI
1. Перейти на `/news`
2. Нажать "Создать новость"
3. Заполнить заголовок (slug сгенерируется автоматически)
4. Написать содержание в TipTap редакторе:
   - Форматировать текст (жирный, курсив)
   - Добавить заголовки (H1-H3)
   - Вставить списки, цитаты, код
5. Выбрать категорию, добавить теги
6. Указать URL обложки (или загрузить через media-service)
7. Раскрыть "SEO-настройки", заполнить метаданные
8. Нажать "Опубликовать" или "Сохранить черновик"

### Просмотр статистики
1. Перейти на Dashboard (`/`)
2. Увидеть:
   - Общее количество новостей
   - Количество опубликованных/черновиков
   - Всего просмотров
   - График динамики за последние 30 дней
   - Топ-5 популярных новостей
   - Последние 5 новостей

## 🛡️ Безопасность

- ✅ JWT Bearer токены для авторизации
- ✅ RBAC - проверка ролей `admin`/`editor` на backend
- ✅ CORS - ограничение разрешенных origins
- ✅ Protected Routes - редирект на `/login` при 401
- ✅ Token refresh через interceptor (auto-logout при 401)
- ✅ SQL injection protection (параметризованные запросы в news-service)

## 📄 Лицензия

MIT

## 👥 Авторы

Новостной портал - Admin Service
