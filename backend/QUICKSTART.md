# 🎯 Quick Start - Быстрый запуск backend

## Локальная разработка (без Docker)

### 1. Установка зависимостей

```bash
cd backend
pip install -r requirements.txt
```

### 2. Настройка .env

Скопируйте `.env.example` в `.env` и заполните ключи (они уже есть в вашем .env)

### 3. Запуск PostgreSQL локально

```bash
# Через Docker
docker run -d \
  --name newshub_postgres \
  -e POSTGRES_USER=newsadmin \
  -e POSTGRES_PASSWORD=ваш_пароль \
  -e POSTGRES_DB=newshub_db \
  -p 5432:5432 \
  postgres:15-alpine
```

### 4. Запуск backend

```bash
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

### 5. Открыть Swagger UI

http://localhost:8000/docs

---

## Тестирование API

### Сбор новостей

```bash
curl -X POST "http://localhost:8000/api/v1/pipeline/collect?channel=crypto"
```

### Просмотр новостей

```bash
curl "http://localhost:8000/api/v1/news/?status=pending"
```

### AI-анализ

```bash
curl -X POST "http://localhost:8000/api/v1/pipeline/analyze?limit=5"
```

### Публикация

```bash
curl -X POST "http://localhost:8000/api/v1/pipeline/post?limit=3"
```

---

## Автоматический режим

Backend автоматически запускает задачи:
- Сбор новостей: каждые 10 минут
- AI-анализ: каждые 5 минут
- Публикация: каждые 7 минут

Просто запустите сервер и всё работает! 🚀

---

## Логи

```bash
# Просмотр логов
tail -f logs/app.log

# Через Docker
docker-compose logs -f backend
```

---

**Готово! Backend работает! ✅**
