# 📡 Использование Telegram API в NewsHub AI

## Обзор

Telegram Bot API используется для автоматической публикации новостей в каналы `@crypto_ainews` и `@kremlin_digest`.

---

## 🔧 Конфигурация

### .env файл

```bash
# Telegram Bot
TELEGRAM_BOT_TOKEN=8286012057:AAG7YfZlvgij4aS-7Z9QzMBFfDhUsHphj9o
TELEGRAM_CRYPTO_CHANNEL=@crypto_ainews
TELEGRAM_POLITICS_CHANNEL=@kremlin_digest
TELEGRAM_ADMIN_CHAT_ID=433868823
```

### Settings (app/core/config.py)

```python
class Settings(BaseSettings):
    TELEGRAM_BOT_TOKEN: str
    TELEGRAM_CRYPTO_CHANNEL: str
    TELEGRAM_POLITICS_CHANNEL: str
    TELEGRAM_ADMIN_CHAT_ID: int
```

---

## 📦 Основные компоненты

### 1. RateLimiter

Ограничение скорости отправки сообщений (20 msg/min):

```python
from app.services.telegram_poster import RateLimiter

limiter = RateLimiter(max_messages=20, time_window=60)

# Использование
await limiter.acquire()  # Блокируется если превышен лимит
await bot.send_message(...)
```

**Параметры:**
- `max_messages`: Максимум сообщений (по умолчанию 20)
- `time_window`: Временное окно в секундах (по умолчанию 60)

**Логика:**
- Хранит timestamp каждой отправки в deque
- Удаляет старые записи (> time_window)
- Блокирует если len(messages) >= max_messages
- Ожидает до истечения time_window

---

### 2. TelegramPoster

Основной сервис для публикации новостей.

#### Инициализация

```python
from app.services.telegram_poster import TelegramPoster

poster = TelegramPoster(db=session)
```

#### Методы

##### format_message(news, parse_mode="HTML")

Форматирует NewsItem в Telegram сообщение.

```python
message = poster.format_message(news, parse_mode="HTML")
```

**Формат HTML:**
```html
🔐 <b>Заголовок новости</b>

📝 [Краткий тизер от AI]

🔍 <b>AI-инсайт:</b>
• Пункт 1
• Пункт 2
• Пункт 3

🔗 <a href='URL'>Читать подробнее</a>

#хэштег1 #хэштег2 #хэштег3
```

**Лимиты:**
- Текст: 4096 символов
- Caption (с фото): 1024 символа
- Автоматическая обрезка если превышено

---

##### post_news(news, notify_admin=True)

Публикация одной новости.

```python
success = await poster.post_news(news, notify_admin=True)
```

**Процесс:**
1. Проверка канала (crypto/politics)
2. Форматирование сообщения
3. Rate limiting
4. Отправка с retry (3 попытки, exponential backoff)
5. Обновление БД (status=PUBLISHED)
6. Уведомление админу (если notify_admin=True)

**Параметры:**
- `news`: NewsItem объект
- `notify_admin`: Отправить уведомление админу

**Возвращает:**
- `True` если успешно
- `False` если ошибка

---

##### post_analyzed_news(limit=5)

Публикация всех проанализированных новостей.

```python
result = await poster.post_analyzed_news(limit=10)
```

**Возвращает:**
```python
{
    'total': 10,      # Всего найдено
    'posted': 8,      # Успешно опубликовано
    'failed': 2,      # Ошибки
}
```

**Логика:**
- Выбирает новости со статусом `ANALYZED`
- Фильтрует `requires_moderation == False`
- Сортирует по `relevance_score DESC`
- Публикует с интервалом 2 секунды (антиспам)

---

##### send_to_admin(news, action="moderation")

Отправка уведомления админу.

```python
await poster.send_to_admin(news, action="moderation")
```

**Кнопки модерации:**
```
✅ Одобрить   ❌ Отклонить
```

**Формат:**
```
🚨 Требуется модерация

Заголовок: [title]
Категория: [category]
Релевантность: 8/10

AI-тизер:
[ai_summary]

Одобрить эту новость для публикации?
```

---

##### handle_moderation_requests()

Обработка всех новостей требующих модерации.

```python
result = await poster.handle_moderation_requests()
```

**Возвращает:**
```python
{
    'total': 5,
    'notified': 5,
    'failed': 0,
}
```

---

## 🔄 Retry логика

### Exponential Backoff

```python
@retry(
    stop=stop_after_attempt(3),
    wait=wait_exponential(multiplier=1, min=2, max=60),
    retry=retry_if_exception_type((RetryAfter, TimedOut, NetworkError)),
)
async def _send_with_retry(channel, message_text, image_url=None):
    ...
```

**Параметры:**
- `stop_after_attempt(3)`: Максимум 3 попытки
- `wait_exponential`: 2s → 4s → 8s → 16s → 32s → 60s
- `retry_if_exception_type`: Retry только при специфичных ошибках

**Обрабатываемые ошибки:**
- `RetryAfter`: Telegram просит подождать X секунд
- `TimedOut`: Timeout запроса
- `NetworkError`: Проблемы с сетью

---

## 🎯 Примеры использования

### 1. Простая публикация

```python
from app.services.telegram_poster import TelegramPoster
from app.models import NewsItem, NewsStatus, NewsChannel

# Создаём новость
news = NewsItem(
    title="Bitcoin $100k",
    content="...",
    url="https://...",
    category=NewsChannel.CRYPTO,
    status=NewsStatus.ANALYZED,
    ai_summary="Краткий тизер...",
    ai_insights=["Инсайт 1", "Инсайт 2"],
    ai_hashtags=["#Bitcoin", "#Crypto"],
)

# Публикуем
poster = TelegramPoster(db=session)
success = await poster.post_news(news)

if success:
    print(f"✅ Опубликовано: {news.telegram_message_id}")
```

---

### 2. Batch публикация

```python
# Публикуем 10 новостей
result = await poster.post_analyzed_news(limit=10)

print(f"Опубликовано: {result['posted']}/{result['total']}")
```

---

### 3. Модерация

```python
# Отправить новость на модерацию
await poster.send_to_admin(news, action="moderation")

# Обработать все запросы модерации
result = await poster.handle_moderation_requests()
print(f"Отправлено запросов: {result['notified']}")
```

---

### 4. Публикация с изображением

```python
news.image_url = "https://example.com/bitcoin.jpg"

# Автоматически попытается send_photo
# Если ошибка → fallback на send_message
success = await poster.post_news(news)
```

---

### 5. Использование через API

```bash
# 1. Собрать новости
curl -X POST "http://localhost:8000/api/v1/pipeline/collect"

# 2. Анализировать
curl -X POST "http://localhost:8000/api/v1/pipeline/analyze?limit=10"

# 3. Опубликовать
curl -X POST "http://localhost:8000/api/v1/pipeline/post?limit=5"
```

---

### 6. Полный pipeline

```bash
# Весь цикл: collect → analyze → post
curl -X POST "http://localhost:8000/api/v1/pipeline/pipeline?channel=crypto"
```

---

## 📊 Мониторинг

### Статистика

```python
poster = TelegramPoster(db=session)

# Публикуем новости...
await poster.post_analyzed_news(limit=10)

# Проверяем статистику
print(poster.stats)
```

**Вывод:**
```python
{
    'total_sent': 8,
    'total_failed': 2,
    'last_error': 'RetryAfter: retry after 30s',
}
```

---

### Логирование

```python
from app.core.logger import get_logger

logger = get_logger(__name__)

# Логи в TelegramPoster
logger.info("Posting news {id}", id=news.id)
logger.error("Error posting news: {error}", error=str(e))
```

**Просмотр логов:**
```bash
# Все логи
docker-compose logs -f backend

# Только Telegram
docker-compose logs backend | grep "TelegramPoster"

# Ошибки
docker-compose logs backend | grep "Error posting"
```

---

## 🛡️ Обработка ошибок

### 1. Unauthorized (401)

**Причина:** Неверный токен

**Решение:**
```bash
# Проверить .env
cat .env | grep TELEGRAM_BOT_TOKEN

# Получить новый токен от @BotFather
```

---

### 2. Chat not found (400)

**Причина:** Неверный chat_id или бот не в канале

**Решение:**
```bash
# Проверить chat_id через getUpdates
curl "https://api.telegram.org/bot{TOKEN}/getUpdates"

# Добавить бота в канал как администратора
```

---

### 3. CHAT_WRITE_FORBIDDEN (403)

**Причина:** Бот не имеет прав на публикацию

**Решение:**
1. Открыть настройки канала
2. Администраторы → Выбрать бота
3. Включить "Публиковать сообщения"

---

### 4. Too Many Requests (429)

**Причина:** Превышен rate limit

**Автоматическая обработка:**
```python
# RetryAfter exception обрабатывается автоматически
try:
    await poster.post_news(news)
except RetryAfter as e:
    await asyncio.sleep(e.retry_after)
    # Retry...
```

**Rate limiter предотвращает:**
```python
# Максимум 20 сообщений/минуту
limiter = RateLimiter(max_messages=20, time_window=60)
```

---

### 5. NetworkError / TimedOut

**Причина:** Проблемы с сетью или блокировка Telegram

**Решение:**
```python
# Использовать прокси
from telegram.request import HTTPXRequest

request = HTTPXRequest(
    proxy_url="http://proxy:8080"
)

bot = Bot(token=TOKEN, request=request)
```

---

## 🔐 Безопасность

### 1. Хранение токена

```bash
# ❌ НЕПРАВИЛЬНО: токен в коде
BOT_TOKEN = "123456:ABC-DEF..."

# ✅ ПРАВИЛЬНО: токен в .env
TELEGRAM_BOT_TOKEN=123456:ABC-DEF...
```

```python
# Загрузка из .env
from app.core.config import settings

bot = Bot(token=settings.TELEGRAM_BOT_TOKEN)
```

---

### 2. Rate limiting

```python
# Встроенный rate limiter
limiter = RateLimiter(max_messages=20, time_window=60)

# Telegram лимиты:
# - 20 msg/min в группы/каналы
# - 30 msg/sec в личку
# - 20 msg/min в одного пользователя
```

---

### 3. Валидация данных

```python
# Валидация перед отправкой
if not news.title or len(news.title) < 3:
    raise ValueError("Title too short")

if not news.url or not news.url.startswith("http"):
    raise ValueError("Invalid URL")

# Экранирование HTML
from html import escape

title = escape(news.title)
```

---

## 📈 Производительность

### Batch отправка

```python
# ❌ МЕДЛЕННО: последовательно
for news in news_items:
    await poster.post_news(news)
    await asyncio.sleep(3)  # 3 секунды задержка

# ✅ БЫСТРО: с rate limiter
await poster.post_analyzed_news(limit=10)
# Автоматический rate limiting, задержка 2 секунды
```

---

### Async публикация

```python
# Используем asyncio.gather для параллельной отправки
tasks = [
    poster.post_news(news)
    for news in news_items[:5]
]

results = await asyncio.gather(*tasks, return_exceptions=True)

# Проверяем результаты
successes = sum(1 for r in results if r is True)
print(f"Успешно: {successes}/{len(results)}")
```

---

## 🧪 Тестирование

### Unit тесты

```bash
# Запустить все тесты
pytest backend/tests/test_telegram_poster.py -v

# Конкретный тест
pytest backend/tests/test_telegram_poster.py::TestRateLimiter::test_rate_limiter_basic -v

# С покрытием
pytest backend/tests/test_telegram_poster.py --cov=app.services.telegram_poster
```

---

### Интеграционный тест

```bash
# Быстрый тест бота
python backend/scripts/test_telegram_quick.py

# Полный интеграционный тест
python backend/scripts/test_integration.py
```

---

### Mock тестирование

```python
from unittest.mock import AsyncMock

# Mock Bot
mock_bot = AsyncMock()
mock_bot.send_message.return_value = Mock(message_id=123)

poster = TelegramPoster(db=mock_db)
poster.bot = mock_bot

# Тест
await poster.post_news(news)

# Проверка
mock_bot.send_message.assert_called_once()
```

---

## 📚 Полезные ссылки

- **Telegram Bot API:** https://core.telegram.org/bots/api
- **python-telegram-bot:** https://docs.python-telegram-bot.org/
- **Tenacity (retry):** https://tenacity.readthedocs.io/
- **Rate Limiting:** https://core.telegram.org/bots/faq#broadcasting-to-users

---

## ✅ Чеклист

- [ ] Токен бота добавлен в `.env`
- [ ] Бот добавлен в каналы как администратор
- [ ] Chat ID каналов получены и настроены
- [ ] Admin Chat ID настроен
- [ ] `test_telegram_quick.py` выполнен успешно
- [ ] Rate limiter работает (20 msg/min)
- [ ] Retry логика тестирована
- [ ] Логирование настроено
- [ ] Мониторинг настроен

**Готово! 🎉**
