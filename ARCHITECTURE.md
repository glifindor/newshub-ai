# NewsHub AI - Архитектура системы

## 📋 Overview (Обзор)

**NewsHub AI** — центральный хаб для автоматического сбора, AI-анализа и публикации новостей в Telegram-каналы:
- 🔐 **@crypto_ainews** — криптовалюта, IT, AI-анализ
- 🏛️ **@kremlin_digest** — политика России и мира

### Ключевые возможности
- ✅ Автоматический сбор новостей из RSS/API
- 🤖 AI-анализ с помощью OpenRouter (GPT-4, Claude)
- 🖼️ Генерация изображений через Freepik API
- 📤 Автопостинг в Telegram с форматированием
- 👨‍💼 Админ-панель для модерации
- 🌐 Публичный архив новостей
- 🔔 Уведомления админам при спорных новостях

---

## 🏗️ Components (Компоненты системы)

### 1️⃣ **News Collector Service** (Сборщик новостей)
**Задача:** Сбор новостей из внешних источников

**Источники данных:**
- RSS-фиды (CoinDesk, Reuters, TASS, RIA)
- NewsAPI.org
- Web scraping (Selenium/BeautifulSoup)
- Twitter API (опционально)

**Функции:**
- Парсинг RSS каждые 5-15 минут (Celery Beat)
- Дедупликация по хэшу контента (MD5)
- Первичная фильтрация по ключевым словам
- Сохранение сырых данных в БД

**Endpoints:**
```
POST /api/collector/sources        # Добавить источник
GET  /api/collector/sources         # Список источников
POST /api/collector/run             # Запустить сбор вручную
GET  /api/collector/stats           # Статистика сбора
```

---

### 2️⃣ **AI Analyzer Service** (AI-анализатор)
**Задача:** Обработка и анализ новостей с помощью ИИ

**AI-провайдеры:**
- **OpenRouter API** (GPT-4, Claude 3, Llama 3)
  - Endpoint: `https://openrouter.ai/api/v1/chat/completions`
  - Промпты для: суммаризации, sentiment analysis, категоризации
- **Freepik API** (генерация обложек)
  - Endpoint: `https://api.freepik.com/v1/ai/text-to-image`

**Процесс анализа:**
1. Извлечь текст новости
2. Отправить в OpenRouter с промптом:
   ```
   "Ты — аналитик. Кратко опиши новость (3-5 предложений), 
   выдели ключевые инсайты, оцени важность (1-10), 
   определи тематику: crypto/it/politics/world"
   ```
3. Сгенерировать изображение (если нет) через Freepik
4. Сохранить результат в БД

**Endpoints:**
```
POST /api/ai/analyze/{news_id}      # Анализировать новость
POST /api/ai/generate-image         # Генерация изображения
GET  /api/ai/status/{task_id}       # Статус задачи
```

---

### 3️⃣ **Telegram Poster Service** (Постер в Telegram)
**Задача:** Публикация в каналы

**Telegram Bot API:**
- Метод: `sendMessage`, `sendPhoto`
- Форматирование: HTML/Markdown
- Rate limits: 30 сообщений/сек

**Формат поста:**
```
🔥 [ЭМОДЗИ] ЗАГОЛОВОК

📝 Краткое описание (AI-саммари)

💡 AI-инсайт: [Аналитика от ИИ]

🔗 Читать подробнее: [ссылка]

#crypto #bitcoin #ai
```

**Логика маршрутизации:**
- `category == "crypto" OR "it"` → @crypto_ainews
- `category == "politics" OR "russia"` → @kremlin_digest
- `importance > 8` → отправить админу на модерацию

**Endpoints:**
```
POST /api/telegram/post             # Опубликовать новость
POST /api/telegram/schedule         # Отложенная публикация
DELETE /api/telegram/post/{id}      # Удалить пост
GET  /api/telegram/stats            # Статистика каналов
```

---

### 4️⃣ **Admin Dashboard** (Панель администратора)
**Задача:** Управление контентом и системой

**Функционал:**
- 📊 Dashboard с метриками (новости/день, успешность AI)
- 📰 Список новостей с фильтрами (статус, категория, источник)
- ✏️ Редактор новостей (WYSIWYG)
- ✅ Одобрение/отклонение перед публикацией
- 🔧 Управление источниками и настройками
- 📈 Аналитика (просмотры, engagement в Telegram)
- 🚨 Логи и ошибки

**Аутентификация:**
- JWT токены (access + refresh)
- Роли: `admin`, `moderator`, `viewer`

**Endpoints:**
```
POST /api/auth/login                # Вход
POST /api/auth/refresh              # Обновить токен
GET  /api/admin/news                # Список новостей
PATCH /api/admin/news/{id}          # Редактировать
POST /api/admin/news/{id}/approve   # Одобрить
DELETE /api/admin/news/{id}         # Удалить
GET  /api/admin/analytics           # Аналитика
```

---

### 5️⃣ **Public Viewer** (Публичный архив)
**Задача:** Отображение новостей для пользователей

**Функционал:**
- 🌐 Главная страница с лентой новостей
- 🔍 Поиск и фильтрация (дата, категория)
- 📄 Страница отдельной новости
- 🏷️ Теги и категории
- 📱 Адаптивный дизайн (Mobile-first)
- ⚡ SSR через Next.js (SEO)

**Endpoints:**
```
GET  /api/public/news               # Лента новостей (пагинация)
GET  /api/public/news/{id}          # Детали новости
GET  /api/public/search             # Поиск
GET  /api/public/categories         # Категории
```

---

## 🔄 Data Flow (Поток данных)

### ASCII-диаграмма:

```
┌─────────────────┐
│  NEWS SOURCES   │
│ (RSS/API/Web)   │
└────────┬────────┘
         │
         ▼
┌─────────────────────────┐
│  COLLECTOR SERVICE      │
│  (Celery Worker)        │
│  • Parse RSS            │
│  • Deduplicate (MD5)    │
│  • Save to PostgreSQL   │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│     REDIS QUEUE         │
│  (Task: analyze_news)   │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│  AI ANALYZER SERVICE    │
│  (Celery Worker)        │
│  • Call OpenRouter      │
│  • Generate Image       │
│  • Categorize           │
│  • Save results         │
└────────┬────────────────┘
         │
         ▼
    ┌───┴────┐
    │        │
    ▼        ▼
┌────────┐ ┌────────────┐
│ AUTO?  │ │  MANUAL?   │
│ Post   │ │  Moderate  │
└───┬────┘ └─────┬──────┘
    │            │
    │            ▼
    │     ┌──────────────┐
    │     │ADMIN PANEL   │
    │     │(Next.js)     │
    │     │• Approve     │
    │     │• Edit        │
    │     └──────┬───────┘
    │            │
    ▼            ▼
┌──────────────────────────┐
│ TELEGRAM POSTER SERVICE  │
│ • Format message         │
│ • Post to channel        │
│ • Track stats            │
└────────┬─────────────────┘
         │
         ▼
┌────────────────────┐
│  TELEGRAM CHANNELS │
│  @crypto_ainews    │
│  @kremlin_digest   │
└────────────────────┘

         │
         ▼
┌────────────────────┐
│  PUBLIC WEBSITE    │
│  (Next.js SSR)     │
│  • News archive    │
│  • Search          │
└────────────────────┘
```

### Подробный flow:

1. **Сбор (каждые 5-15 мин):**
   - Celery Beat → Запуск задачи `collect_news`
   - Collector Service → Парсинг источников
   - PostgreSQL → Сохранение `status=pending`

2. **Анализ (асинхронно):**
   - Webhook/Trigger → Задача `analyze_news` в Redis
   - AI Analyzer → OpenRouter API (GPT-4)
   - Freepik API → Генерация изображения
   - PostgreSQL → Обновление `status=analyzed`

3. **Модерация (если нужно):**
   - IF `importance > 8` OR `category=politics`:
     - Telegram Bot → Уведомление админу
     - Admin Panel → Ожидание одобрения
   - ELSE:
     - Автопубликация

4. **Публикация:**
   - Telegram Poster → Форматирование
   - Bot API → Постинг в канал
   - PostgreSQL → `status=published`

5. **Отображение:**
   - Next.js → SSR рендеринг
   - Public Viewer → Показ архива

---

## 💻 Tech Stack (Технологии)

### Backend
- **Language:** Python 3.11
- **Framework:** FastAPI (async, high performance)
- **ORM:** SQLAlchemy 2.0 (async mode)
- **Task Queue:** Celery 5.3 + RabbitMQ
- **Scheduler:** Celery Beat

### Database & Cache
- **Primary DB:** PostgreSQL 15 (новости, пользователи, логи)
- **Cache/Queue:** Redis 7 (Celery broker, кэш запросов)
- **Search:** PostgreSQL Full-Text Search (или Elasticsearch)

### AI & External APIs
- **AI Provider:** OpenRouter (`openrouter.ai`)
  - Models: GPT-4, Claude 3 Opus, Llama 3
- **Image Generation:** Freepik API (`api.freepik.com`)
- **News APIs:**
  - NewsAPI.org
  - RSS Feeds (feedparser)
  - Web Scraping (BeautifulSoup4, httpx)
- **Telegram:** python-telegram-bot

### Frontend
- **Framework:** Next.js 14 (React 18)
- **UI Library:** Tailwind CSS + shadcn/ui
- **State:** React Query (server state) + Zustand (client state)
- **Forms:** React Hook Form + Zod validation

### DevOps & Infrastructure
- **Containerization:** Docker + Docker Compose
- **Orchestration:** Docker Swarm (или Kubernetes если нужно)
- **Web Server:** Nginx (reverse proxy)
- **SSL:** Let's Encrypt (Certbot)
- **CI/CD:** GitHub Actions
- **Monitoring:** Prometheus + Grafana + Sentry
- **Logging:** ELK Stack (Elasticsearch, Logstash, Kibana) или Loki

### Security
- **Secrets:** Environment variables (.env) + Docker secrets
- **Auth:** JWT (access 15min, refresh 7 days)
- **Rate Limiting:** Redis-based (10 req/sec per IP)
- **Input Validation:** Pydantic models
- **CORS:** Configured for specific origins

---

## 🚀 Deployment (Развертывание на сервере)

### Сервер: **151.241.228.203** (Ubuntu)

### Архитектура развертывания:

```
                    INTERNET
                       │
                       ▼
              ┌────────────────┐
              │   Cloudflare   │ (опционально, для DDoS protection)
              │   DNS Proxy    │
              └────────┬───────┘
                       │
                       ▼
              ┌────────────────┐
              │  NGINX (80/443)│
              │  Reverse Proxy │
              │  SSL Termination│
              └────────┬───────┘
                       │
           ┌───────────┴───────────┐
           │                       │
           ▼                       ▼
    ┌─────────────┐        ┌─────────────┐
    │  Frontend   │        │  Backend    │
    │  Next.js    │        │  FastAPI    │
    │  Port 3000  │        │  Port 8000  │
    └─────────────┘        └──────┬──────┘
                                  │
                    ┌─────────────┼─────────────┐
                    │             │             │
                    ▼             ▼             ▼
              ┌──────────┐  ┌──────────┐  ┌──────────┐
              │PostgreSQL│  │  Redis   │  │ RabbitMQ │
              │Port 5432 │  │Port 6379 │  │Port 5672 │
              └──────────┘  └──────────┘  └──────────┘
                    │
                    ▼
              ┌──────────┐
              │  Celery  │
              │  Workers │
              └──────────┘
```

### Шаги развертывания:

#### 1. Подготовка сервера
```bash
# SSH подключение
ssh root@151.241.228.203

# Обновление системы
apt update && apt upgrade -y

# Установка Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Установка Docker Compose
apt install docker-compose -y

# Установка Git
apt install git -y

# Создание пользователя для деплоя
useradd -m -s /bin/bash newsadmin
usermod -aG docker newsadmin
```

#### 2. Клонирование проекта
```bash
su - newsadmin
git clone https://github.com/glifindor/newsportal.git /home/newsadmin/newshub
cd /home/newsadmin/newshub
```

#### 3. Настройка переменных окружения
```bash
# Создать .env файл
nano .env
```

**Пример .env:**
```env
# Database
POSTGRES_USER=newsadmin
POSTGRES_PASSWORD=SUPER_SECRET_PASSWORD_123
POSTGRES_DB=newshub_db
DATABASE_URL=postgresql+asyncpg://newsadmin:SUPER_SECRET_PASSWORD_123@postgres:5432/newshub_db

# Redis
REDIS_URL=redis://redis:6379/0

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

# JWT
JWT_SECRET_KEY=YOUR_SECRET_KEY_CHANGE_THIS_TO_RANDOM_STRING
JWT_ALGORITHM=HS256
JWT_ACCESS_TOKEN_EXPIRE_MINUTES=15
JWT_REFRESH_TOKEN_EXPIRE_DAYS=7

# OpenRouter API
OPENROUTER_API_KEY=sk-or-v1-YOUR_API_KEY
OPENROUTER_API_URL=https://openrouter.ai/api/v1

# Freepik API
FREEPIK_API_KEY=YOUR_FREEPIK_API_KEY

# NewsAPI
NEWSAPI_KEY=YOUR_NEWSAPI_KEY

# Telegram Bot
TELEGRAM_BOT_TOKEN=YOUR_BOT_TOKEN
TELEGRAM_CRYPTO_CHANNEL=@crypto_ainews
TELEGRAM_POLITICS_CHANNEL=@kremlin_digest
TELEGRAM_ADMIN_CHAT_ID=YOUR_ADMIN_CHAT_ID

# Frontend
NEXT_PUBLIC_API_URL=https://151.241.228.203/api

# Environment
ENVIRONMENT=production
DEBUG=False
```

#### 4. Docker Compose конфигурация
```bash
nano docker-compose.yml
```

**docker-compose.yml** (см. отдельный файл в проекте)

#### 5. SSL сертификат (Let's Encrypt)
```bash
# Установка Certbot
apt install certbot python3-certbot-nginx -y

# Получение сертификата (ЗАМЕНИТЕ на ваш домен или используйте IP)
certbot certonly --standalone -d newshub.example.com

# Или для IP (без домена, self-signed):
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /etc/ssl/private/nginx-selfsigned.key \
  -out /etc/ssl/certs/nginx-selfsigned.crt
```

#### 6. Запуск проекта
```bash
# Сборка и запуск
docker-compose up -d --build

# Проверка статуса
docker-compose ps

# Логи
docker-compose logs -f

# Миграции БД
docker-compose exec backend alembic upgrade head

# Создание суперпользователя
docker-compose exec backend python scripts/create_admin.py
```

#### 7. Nginx конфигурация
```bash
nano /etc/nginx/sites-available/newshub
```

**Пример конфигурации** (см. отдельный файл)

```bash
# Активация конфигурации
ln -s /etc/nginx/sites-available/newshub /etc/nginx/sites-enabled/
nginx -t
systemctl reload nginx
```

#### 8. Мониторинг и логи
```bash
# Prometheus (порт 9090)
docker-compose exec prometheus

# Grafana (порт 3001)
# Login: admin / admin (сменить при первом входе)

# Логи Celery
docker-compose logs -f celery_worker

# Логи Backend
docker-compose logs -f backend
```

---

## 📊 Database Schema (Схема БД)

### Таблицы:

#### 1. `news_sources` (Источники новостей)
```sql
CREATE TABLE news_sources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'rss', 'api', 'scraping'
    url TEXT NOT NULL,
    category VARCHAR(50), -- 'crypto', 'politics', 'it', 'world'
    is_active BOOLEAN DEFAULT true,
    last_fetched_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### 2. `news_items` (Новости)
```sql
CREATE TABLE news_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id INTEGER REFERENCES news_sources(id),
    
    -- Контент
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    summary TEXT, -- AI-generated
    url TEXT NOT NULL UNIQUE,
    image_url TEXT,
    author VARCHAR(255),
    
    -- Метаданные
    category VARCHAR(50), -- 'crypto', 'politics', 'it', 'world'
    tags TEXT[], -- ['bitcoin', 'regulation', 'usa']
    language VARCHAR(10) DEFAULT 'ru',
    
    -- AI-анализ
    ai_summary TEXT,
    ai_insights TEXT,
    sentiment VARCHAR(20), -- 'positive', 'negative', 'neutral'
    importance_score INTEGER CHECK (importance_score BETWEEN 1 AND 10),
    
    -- Статус
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'analyzed', 'approved', 'published', 'rejected'
    requires_moderation BOOLEAN DEFAULT false,
    
    -- Публикация
    published_at TIMESTAMP,
    telegram_message_id INTEGER,
    telegram_channel VARCHAR(100),
    
    -- Дедупликация
    content_hash VARCHAR(64) UNIQUE,
    
    -- Timestamps
    source_published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_news_status ON news_items(status);
CREATE INDEX idx_news_category ON news_items(category);
CREATE INDEX idx_news_published ON news_items(published_at DESC);
CREATE INDEX idx_news_hash ON news_items(content_hash);
```

#### 3. `users` (Пользователи/Админы)
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'moderator', -- 'admin', 'moderator', 'viewer'
    is_active BOOLEAN DEFAULT true,
    telegram_chat_id BIGINT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### 4. `telegram_posts` (История постов в Telegram)
```sql
CREATE TABLE telegram_posts (
    id SERIAL PRIMARY KEY,
    news_item_id UUID REFERENCES news_items(id),
    channel VARCHAR(100) NOT NULL,
    message_id BIGINT NOT NULL,
    views_count INTEGER DEFAULT 0,
    posted_at TIMESTAMP DEFAULT NOW()
);
```

#### 5. `ai_tasks` (Задачи AI-обработки)
```sql
CREATE TABLE ai_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    news_item_id UUID REFERENCES news_items(id),
    task_type VARCHAR(50), -- 'analyze', 'generate_image', 'summarize'
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed'
    provider VARCHAR(50), -- 'openrouter', 'freepik'
    model VARCHAR(100),
    input_data JSONB,
    output_data JSONB,
    error_message TEXT,
    processing_time_ms INTEGER,
    cost_usd NUMERIC(10, 6),
    created_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);
```

#### 6. `system_logs` (Системные логи)
```sql
CREATE TABLE system_logs (
    id SERIAL PRIMARY KEY,
    level VARCHAR(20), -- 'info', 'warning', 'error', 'critical'
    service VARCHAR(100), -- 'collector', 'ai_analyzer', 'telegram_poster'
    message TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 🔒 Security (Безопасность)

### 1. API Keys & Secrets
- ✅ Все ключи в `.env` файле (не коммитить в Git!)
- ✅ Docker secrets для production
- ✅ Ротация ключей каждые 90 дней

### 2. Аутентификация
- ✅ JWT токены (HMAC SHA-256)
- ✅ Access token (15 мин) + Refresh token (7 дней)
- ✅ Хэширование паролей (bcrypt)
- ✅ Rate limiting: 10 запросов/сек на IP

### 3. Защита от дубликатов
```python
import hashlib

def calculate_content_hash(title: str, content: str) -> str:
    """Хэш для дедупликации"""
    combined = f"{title}|{content}".encode('utf-8')
    return hashlib.md5(combined).hexdigest()
```

### 4. Валидация входных данных
```python
from pydantic import BaseModel, HttpUrl, validator

class NewsCreate(BaseModel):
    title: str
    content: str
    url: HttpUrl
    
    @validator('title')
    def title_length(cls, v):
        if len(v) < 10 or len(v) > 500:
            raise ValueError('Title must be 10-500 chars')
        return v
```

### 5. CORS Configuration
```python
from fastapi.middleware.cors import CORSMiddleware

app.add_middleware(
    CORSMiddleware,
    allow_origins=["https://151.241.228.203", "https://newshub.example.com"],
    allow_methods=["GET", "POST", "PUT", "DELETE"],
    allow_headers=["*"],
)
```

---

## ⚠️ Risks & Mitigation (Риски и решения)

### 1. **Rate Limits внешних API**
**Риск:** OpenRouter, NewsAPI, Telegram имеют лимиты запросов

**Mitigation:**
- Redis-кэш для повторных запросов (TTL: 1 час)
- Exponential backoff при ошибках
- Приоритизация новостей (важные первыми)
- Fallback на альтернативные AI-модели

```python
from tenacity import retry, stop_after_attempt, wait_exponential

@retry(stop=stop_after_attempt(3), wait=wait_exponential(min=1, max=10))
async def call_openrouter_with_retry(prompt: str):
    # ...
```

### 2. **Стоимость AI-анализа**
**Риск:** GPT-4 дорогой (~$30/1M tokens)

**Mitigation:**
- Использовать Claude 3 Haiku (дешевле) для простых задач
- Кэшировать AI-ответы
- Batch processing (анализ пакетами)
- Мониторинг бюджета (alert при >$100/день)

### 3. **Спам/низкокачественный контент**
**Риск:** Источники могут публиковать мусор

**Mitigation:**
- AI-фильтрация (importance_score < 4 отбрасывать)
- Whitelist проверенных источников
- Модерация для новых источников
- Блэклист доменов

### 4. **Падение сервера/компонентов**
**Риск:** PostgreSQL/Redis/RabbitMQ может упасть

**Mitigation:**
- Docker healthchecks + автоперезапуск
- Backup БД каждые 6 часов (pg_dump)
- Redis persistence (AOF)
- Prometheus alerts (CPU > 80%, RAM > 90%)

```yaml
# docker-compose.yml
services:
  postgres:
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U $POSTGRES_USER"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
```

### 5. **Безопасность Telegram Bot**
**Риск:** Утечка токена → контроль над каналом

**Mitigation:**
- Токен в Docker secrets (не в коде!)
- Webhook вместо polling (safer)
- IP whitelist для webhook endpoint
- 2FA для Telegram-аккаунта админа

### 6. **SEO и индексация**
**Риск:** Публичный сайт не индексируется Google

**Mitigation:**
- Next.js SSR (Server-Side Rendering)
- Sitemap.xml и robots.txt
- Open Graph tags для соцсетей
- Canonical URLs

### 7. **Дублирование новостей в каналах**
**Риск:** Одна новость может подойти под обе категории

**Mitigation:**
- Жесткая категоризация (crypto OR politics, не AND)
- AI-проверка: "В какой канал больше подходит?"
- Флаг `telegram_channel` в БД (уникальность)

---

## 📈 Scaling Strategy (Масштабирование)

### Этап 1: Одиночный сервер (текущий)
- 1000 новостей/день ✅
- 1 Backend instance
- 1 PostgreSQL
- 2-4 Celery workers

### Этап 2: Горизонтальное масштабирование (>5000 новостей/день)
- Load Balancer (Nginx)
- 3+ Backend replicas
- PostgreSQL Primary + Read Replicas
- 10+ Celery workers
- Redis Cluster

### Этап 3: Микросервисы (>50k новостей/день)
- Separate services: Collector, Analyzer, Poster
- Kubernetes orchestration
- Message Queue (Kafka вместо RabbitMQ)
- Elasticsearch для поиска
- CDN (Cloudflare) для статики

---

## 🧪 Testing Strategy (Тестирование)

### 1. Unit Tests
```bash
pytest tests/unit/
```
- Тесты функций парсинга
- Валидация Pydantic-моделей
- Хэширование контента

### 2. Integration Tests
```bash
pytest tests/integration/
```
- API endpoints (FastAPI TestClient)
- Database operations (SQLAlchemy)
- Celery tasks

### 3. E2E Tests
```bash
playwright test
```
- Админ-панель (логин, создание/редактирование новостей)
- Публичный сайт (навигация, поиск)

### 4. Load Testing
```bash
locust -f tests/load/locustfile.py
```
- Симуляция 1000 одновременных запросов

---

## 📝 API Documentation (Swagger)

FastAPI автоматически генерирует документацию:
- **Swagger UI:** `http://151.241.228.203:8000/docs`
- **ReDoc:** `http://151.241.228.203:8000/redoc`
- **OpenAPI JSON:** `http://151.241.228.203:8000/openapi.json`

---

## 🔧 Maintenance (Обслуживание)

### Ежедневно:
- ✅ Проверка логов ошибок (Sentry dashboard)
- ✅ Мониторинг Grafana (CPU, RAM, disk)

### Еженедельно:
- ✅ Обзор AI costs (OpenRouter billing)
- ✅ Анализ engagement в Telegram
- ✅ Ротация логов (логи старше 30 дней удалять)

### Ежемесячно:
- ✅ Backup БД (скачать на локальный ПК)
- ✅ Обновление зависимостей (`pip-audit`, `npm audit`)
- ✅ Ревью источников новостей (удалить неактивные)

---

## 📞 Support & Contacts

- **GitHub Repo:** https://github.com/glifindor/newsportal
- **Telegram Admin:** @your_admin_username
- **Server IP:** 151.241.228.203

---

## 🎯 Roadmap (Будущие фичи)

### Q1 2025:
- ✅ Запуск MVP (сбор + постинг в Telegram)
- ✅ Админ-панель
- ✅ Публичный сайт

### Q2 2025:
- 🔄 Мобильное приложение (Flutter)
- 🔄 Голосовые новости (TTS + Podcast)
- 🔄 Интеграция с Twitter/X для постинга

### Q3 2025:
- 🔄 Sentiment trading bot (покупка крипты на основе новостей)
- 🔄 Пользовательские подписки (email-рассылка)
- 🔄 Multi-язык (EN, RU, CN)

---

## 📄 License

MIT License - свободное использование

---

**Автор документа:** AI Architect  
**Дата:** 2025-01-18  
**Версия:** 1.0.0
