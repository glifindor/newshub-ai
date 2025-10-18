# 🚀 NewsHub AI Backend - Полная Инструкция

## 📋 Описание

Backend модуль NewsHub AI для автоматического сбора, AI-анализа и публикации новостей в Telegram каналы.

---

## 🏗️ Архитектура

```
┌─────────────────┐
│  RSS/API Sources│
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  NewsCollector  │ ← Сбор новостей каждые 10 мин
│  (scheduler)    │   Фильтрация по keywords
└────────┬────────┘   Дедупликация по MD5 hash
         │
         ▼
┌─────────────────┐
│   PostgreSQL    │ ← Сохранение с status=PENDING
│   NewsItem      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   AIAnalyzer    │ ← AI-анализ каждые 5 мин
│  (OpenRouter)   │   GPT-4/Claude для анализа
└────────┬────────┘   Relevance score 0-10
         │
         ▼
┌─────────────────┐
│   PostgreSQL    │ ← Обновление: status=ANALYZED
│   (updated)     │   или REJECTED (score < 7)
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
    ▼         ▼
┌────────┐ ┌──────────┐
│ Auto   │ │ Manual   │
│ Post   │ │Moderation│
└───┬────┘ └─────┬────┘
    │            │
    ▼            ▼
┌──────────────────┐
│ TelegramPoster   │ ← Публикация каждые 7 мин
│ @crypto_ainews   │   HTML форматирование
│ @kremlin_digest  │   С изображениями
└──────────────────┘
```

---

## 📦 Установка

### 1. Зависимости

```bash
cd backend
pip install -r requirements.txt
```

### 2. Настройка .env

Создайте файл `.env` в корне проекта (уже есть в примере):

```env
# Database
POSTGRES_USER=newsadmin
POSTGRES_PASSWORD=ваш_пароль
POSTGRES_DB=newshub_db
DATABASE_URL=postgresql+asyncpg://newsadmin:пароль@postgres:5432/newshub_db

# OpenRouter API
OPENROUTER_API_KEY=sk-or-v1-ваш_ключ
OPENROUTER_API_URL=https://openrouter.ai/api/v1
OPENROUTER_MODEL=openai/gpt-4-turbo-preview

# Freepik API
FREEPIK_API_KEY=ваш_ключ

# NewsAPI
NEWSAPI_KEY=ваш_ключ

# Telegram
TELEGRAM_BOT_TOKEN=ваш_токен
TELEGRAM_CRYPTO_CHANNEL=@crypto_ainews
TELEGRAM_POLITICS_CHANNEL=@kremlin_digest
TELEGRAM_ADMIN_CHAT_ID=ваш_id

# Settings
COLLECT_INTERVAL_MINUTES=10
AI_TIMEOUT_SECONDS=30
AI_IMPORTANCE_THRESHOLD=7
```

### 3. Инициализация БД

```bash
# Через Alembic (рекомендуется для продакшена)
alembic upgrade head

# Или автоматически при запуске (в коде lifespan)
```

---

## 🚀 Запуск

### Локально (для разработки)

```bash
# Запуск API сервера
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

# Откройте в браузере:
# http://localhost:8000/docs - Swagger UI
# http://localhost:8000/redoc - ReDoc
```

### Docker (для продакшена)

```bash
# Из корня проекта
docker-compose up -d backend

# Проверка логов
docker-compose logs -f backend
```

---

## 📡 API Endpoints

### 🔹 Pipeline (основные)

#### POST `/api/v1/pipeline/collect`
Запустить сбор новостей вручную

**Query параметры:**
- `channel` (optional): `crypto` или `politics`

**Пример:**
```bash
curl -X POST "http://localhost:8000/api/v1/pipeline/collect?channel=crypto"
```

**Ответ:**
```json
{
  "message": "News collection completed",
  "result": {
    "total_collected": 25,
    "sources": {
      "CoinTelegraph RSS": 10,
      "CoinDesk RSS": 15
    },
    "timestamp": "2025-10-18T12:00:00"
  }
}
```

---

#### POST `/api/v1/pipeline/analyze`
Запустить AI-анализ pending новостей

**Query параметры:**
- `limit` (optional): количество новостей (default: 10)

**Пример:**
```bash
curl -X POST "http://localhost:8000/api/v1/pipeline/analyze?limit=20"
```

**Ответ:**
```json
{
  "message": "News analysis completed",
  "result": {
    "total": 20,
    "analyzed": 15,
    "rejected": 5,
    "failed": 0
  }
}
```

---

#### POST `/api/v1/pipeline/post`
Запустить публикацию в Telegram

**Query параметры:**
- `limit` (optional): количество новостей (default: 5)

**Пример:**
```bash
curl -X POST "http://localhost:8000/api/v1/pipeline/post?limit=10"
```

---

#### POST `/api/v1/pipeline/pipeline`
Запустить полный цикл: collect → analyze → post

**Пример:**
```bash
curl -X POST "http://localhost:8000/api/v1/pipeline/pipeline?channel=crypto"
```

---

### 🔹 News

#### GET `/api/v1/news/`
Получить список новостей

**Query параметры:**
- `page`: номер страницы (default: 1)
- `per_page`: новостей на странице (default: 20, max: 100)
- `category`: `crypto` или `politics`
- `status`: `pending`, `analyzed`, `published`, `rejected`
- `search`: поиск по заголовку/содержимому

**Пример:**
```bash
curl "http://localhost:8000/api/v1/news/?category=crypto&status=analyzed&page=1&per_page=20"
```

---

#### GET `/api/v1/news/{news_id}`
Получить детали новости

**Пример:**
```bash
curl "http://localhost:8000/api/v1/news/550e8400-e29b-41d4-a716-446655440000"
```

---

#### POST `/api/v1/news/{news_id}/approve`
Одобрить новость для публикации

**Пример:**
```bash
curl -X POST "http://localhost:8000/api/v1/news/550e8400-e29b-41d4-a716-446655440000/approve"
```

---

#### POST `/api/v1/news/{news_id}/reject`
Отклонить новость

---

### 🔹 Sources

#### GET `/api/v1/sources/`
Получить список источников новостей

#### POST `/api/v1/sources/`
Добавить новый источник

**Body:**
```json
{
  "name": "Custom RSS Feed",
  "type": "rss",
  "url": "https://example.com/rss",
  "category": "crypto"
}
```

---

## 🤖 Автоматизация (Scheduler)

Backend автоматически запускает задачи по расписанию:

| Задача | Интервал | Описание |
|--------|----------|----------|
| `collect_news_job` | 10 минут | Сбор новостей из всех источников |
| `analyze_news_job` | 5 минут | AI-анализ pending новостей |
| `post_news_job` | 7 минут | Публикация analyzed новостей + модерация |

**Конфигурация:**
```python
# app/services/scheduler.py

# Изменить интервалы:
COLLECT_INTERVAL_MINUTES = 10  # в .env
```

---

## 🧪 Тестирование

### Ручное тестирование через API

```bash
# 1. Запустить сбор новостей для криптовалют
curl -X POST "http://localhost:8000/api/v1/pipeline/collect?channel=crypto"

# 2. Проверить что новости появились (статус: pending)
curl "http://localhost:8000/api/v1/news/?status=pending&category=crypto"

# 3. Запустить AI-анализ
curl -X POST "http://localhost:8000/api/v1/pipeline/analyze?limit=5"

# 4. Проверить проанализированные новости
curl "http://localhost:8000/api/v1/news/?status=analyzed&category=crypto"

# 5. Опубликовать в Telegram
curl -X POST "http://localhost:8000/api/v1/pipeline/post?limit=3"

# 6. Проверить опубликованные
curl "http://localhost:8000/api/v1/news/?status=published&category=crypto"
```

### Тестовые скрипты

```bash
# Тест OpenRouter API
python scripts/test_openrouter.py

# Тест Telegram Bot
python scripts/test_telegram.py

# Тест сбора RSS
python scripts/test_collector.py
```

---

## 📊 Структура БД

### Таблица `news_items`

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID | Первичный ключ |
| `source_id` | Integer | ID источника |
| `title` | String(500) | Заголовок |
| `content` | Text | Полный текст |
| `url` | Text | Ссылка на оригинал |
| `image_url` | Text | URL изображения |
| `category` | Enum | `crypto` или `politics` |
| `ai_summary` | Text | AI-тизер (100 слов) |
| `ai_insights` | Text | AI-инсайты (пункты) |
| `ai_hashtags` | Array[String] | Хэштеги |
| `relevance_score` | Float | Рейтинг 0-10 |
| `status` | Enum | `pending/analyzed/published/rejected` |
| `requires_moderation` | Boolean | Требует ручного одобрения |
| `content_hash` | String(64) | MD5 для дедупликации |
| `created_at` | DateTime | Дата создания |
| `published_at` | DateTime | Дата публикации |

---

## 🔧 Логика работы

### 1. Сбор новостей (NewsCollector)

```python
# app/services/collector.py

# Источники:
- RSS feeds (feedparser)
- NewsAPI.org (httpx)

# Фильтрация по keywords:
crypto_keywords = ["bitcoin", "ethereum", "crypto", "blockchain", ...]
politics_keywords = ["kremlin", "putin", "russia", "ukraine", ...]

# Дедупликация:
content_hash = md5(f"{title}|{content}")
```

### 2. AI-анализ (AIAnalyzer)

```python
# app/services/ai_analyzer.py

# Промпт для OpenRouter:
"""
Анализируй новость [текст].
Верни JSON:
{
  "teaser": "краткий тизер 80-120 слов",
  "insights": ["инсайт 1", "инсайт 2"],
  "relevance_score": 8,
  "hashtags": ["#Crypto", "#Bitcoin"]
}
"""

# Логика:
if relevance_score >= 7:
    status = ANALYZED
    if relevance_score >= 8:
        requires_moderation = True
else:
    status = REJECTED
```

### 3. Публикация (TelegramPoster)

```python
# app/services/telegram_poster.py

# Форматирование:
message = f"""
🔐 <b>{title}</b>

📝 {ai_summary}

💡 <b>AI-инсайт:</b>
{ai_insights}

🔗 Читать подробнее: {url}

{hashtags}
"""

# Отправка:
bot.send_photo(channel, image_url, caption=message, parse_mode=HTML)
```

---

## ⚙️ Конфигурация

### Изменение источников новостей

```python
# app/services/collector.py

# Редактировать default_sources в функции initialize_default_sources()

default_sources = [
    {
        'name': 'Ваш RSS',
        'type': SourceType.RSS,
        'url': 'https://example.com/rss',
        'category': NewsChannel.CRYPTO,
        'keywords': ['bitcoin', 'crypto'],
    },
]
```

### Изменение промптов AI

```python
# app/services/ai_analyzer.py

# Редактировать функцию create_analysis_prompt()

prompt = f"""
Ваш кастомный промпт...
"""
```

### Изменение форматирования Telegram

```python
# app/services/telegram_poster.py

# Редактировать функцию format_message()

message = f"Ваш кастомный формат..."
```

---

## 🐛 Отладка

### Логи

```bash
# Просмотр логов в реальном времени
docker-compose logs -f backend

# Поиск ошибок
docker-compose logs backend | grep ERROR

# Логи конкретного сервиса
docker-compose logs backend | grep "NewsCollector"
```

### Проверка БД

```bash
# Подключение к PostgreSQL
docker-compose exec postgres psql -U newsadmin newshub_db

# SQL запросы:
SELECT COUNT(*) FROM news_items WHERE status = 'pending';
SELECT COUNT(*) FROM news_items WHERE status = 'analyzed';
SELECT COUNT(*) FROM news_items WHERE status = 'published';

# Последние 10 новостей
SELECT title, status, relevance_score, created_at 
FROM news_items 
ORDER BY created_at DESC 
LIMIT 10;
```

---

## 🔒 Безопасность

- ✅ API ключи в `.env` (не коммитить!)
- ✅ JWT аутентификация для админских endpoints
- ✅ Rate limiting (10 req/sec)
- ✅ Input validation (Pydantic)
- ✅ SQL injection protection (SQLAlchemy ORM)

---

## 📈 Масштабирование

### Оптимизация производительности

```python
# Увеличить количество анализируемых новостей за раз
await analyzer.analyze_pending_news(limit=50)

# Параллельная обработка
import asyncio
tasks = [collector.collect_from_source(s) for s in sources]
results = await asyncio.gather(*tasks)
```

### Мониторинг

```bash
# Prometheus метрики доступны на /metrics
curl http://localhost:8000/metrics

# Grafana dashboard на порту 3001
```

---

## 🆘 Частые проблемы

### "OpenRouter API error 401"
- Проверьте API ключ в `.env`
- Убедитесь что есть баланс на openrouter.ai

### "Telegram Bot not responding"
- Проверьте токен бота
- Убедитесь что бот добавлен в каналы как админ

### "No news collected"
- Проверьте доступность RSS feeds
- Проверьте keywords фильтрацию

### "Database connection failed"
- Проверьте DATABASE_URL в `.env`
- Убедитесь что PostgreSQL запущен

---

## 📞 Поддержка

- GitHub Issues: https://github.com/glifindor/newsportal/issues
- Telegram: @your_admin

---

**Готово! Backend полностью функционален! 🎉**
