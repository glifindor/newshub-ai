# 🤖 NewsHub AI - Автоматизированная новостная платформа

> Умная система для сбора, обработки и публикации новостей с использованием искусственного интеллекта

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Python](https://img.shields.io/badge/python-3.11+-blue.svg)](https://www.python.org/)
[![FastAPI](https://img.shields.io/badge/FastAPI-0.104+-green.svg)](https://fastapi.tiangolo.com/)
[![Next.js](https://img.shields.io/badge/Next.js-14-black.svg)](https://nextjs.org/)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://www.docker.com/)

---

## ✨ Особенности

- **🤖 AI-обработка** - Автоматический анализ с OpenRouter (GPT-4, Claude 3.5, Gemini)
- **📰 Мультиисточники** - RSS, API, web scraping с rate limiting
- **🎯 Категоризация** - Автоматическое определение тематики (криптовалюты, политика)
- **📊 Админ-панель** - Modern UI на Next.js 14 + Tailwind + TanStack Table
- **🔄 Автопубликация** - Telegram каналы с retry logic и fallback
- **📈 Мониторинг** - Prometheus + Grafana + алерты в Telegram

---

## 🚀 Демо

- 🌐 **Frontend:** http://151.241.228.203
- 📚 **API Docs:** http://151.241.228.203/docs
- 📊 **Grafana:** http://151.241.228.203:3001
- 🔐 **Crypto Channel:** [@crypto_ainews](https://t.me/crypto_ainews)
- 🏛️ **Politics Channel:** [@kremlin_digest](https://t.me/kremlin_digest)

---

## 📚 Документация

### 📖 Для быстрого старта

- **[QUICK_DEPLOY.md](./QUICK_DEPLOY.md)** - 🚀 Быстрый деплой за 15 минут
- **[QUICK_START.md](./QUICK_START.md)** - ⚡ Локальная разработка

### 📚 Подробная документация

- **[PRODUCTION_DEPLOYMENT.md](./PRODUCTION_DEPLOYMENT.md)** - 🏗️ Production deploy (полный гайд)
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** - 🔧 Базовая инструкция по деплою
- **[ARCHITECTURE.md](./ARCHITECTURE.md)** - 🏛️ Архитектура системы

### 📡 Специализированные гайды

- **[backend/TELEGRAM_BOT_SETUP.md](./backend/TELEGRAM_BOT_SETUP.md)** - 📱 Настройка Telegram бота
- **[frontend/README.md](./frontend/README.md)** - 🎨 Frontend документация
- **[frontend/SETUP.md](./frontend/SETUP.md)** - ⚙️ Frontend setup

---

## ⚡ Быстрый старт

### Вариант 1: Автоматический деплой (Windows)

```powershell
# 1. Клонировать
git clone https://github.com/glifindor/newsportal.git
cd newsportal

# 2. Запустить интерактивную установку
.\scripts\setup-interactive.ps1

# Скрипт запросит пароль от сервера и API ключи
# Подождите 10-15 минут ☕
# Готово! 🎉
```

### Вариант 2: Docker локально

```bash
# 1. Клонировать
git clone https://github.com/glifindor/newsportal.git
cd newsportal

# 2. Настроить .env
cp .env.example .env
nano .env  # Заполнить API ключи

# 3. Запустить
docker-compose up -d --build

# 4. Открыть
# http://localhost:3000 - Frontend
# http://localhost:8000/docs - API
# http://localhost:3001 - Grafana
```

### Вариант 3: Локальная разработка

```bash
# Backend
cd backend
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --reload

# Frontend
cd frontend
npm install
npm run dev
```

---

## 🛠️ Технологический стек

### Backend
```
FastAPI 0.104+      SQLAlchemy 2.0      Celery 5.3
PostgreSQL 15       Redis 7             RabbitMQ 3.12
Alembic             Pydantic V2         Python 3.11
```

### Frontend
```
Next.js 14          React 18            TypeScript 5
Tailwind CSS 3      TanStack Query      TanStack Table
NextAuth.js         Socket.IO           Framer Motion
```

### Infrastructure
```
Docker & Compose    Nginx               Prometheus
Grafana             GitHub Actions      Let's Encrypt
```

---

## 🔑 Необходимые API ключи

1. ✅ **OpenRouter** - https://openrouter.ai/keys
2. ✅ **Telegram Bot** - https://t.me/BotFather
3. ✅ **Telegram Chat ID** - https://t.me/userinfobot
4. ⚪ **NewsAPI** - https://newsapi.org/ (опционально)
5. ⚪ **Freepik** - https://freepik.com/api (опционально)

---

## 📊 Архитектура

```
Client → Nginx → Frontend (Next.js) + Backend (FastAPI)
                      ↓
           PostgreSQL + Redis + RabbitMQ
                      ↓
              Celery Workers
                      ↓
           External APIs (OpenRouter, Telegram)
```

Подробнее: [ARCHITECTURE.md](./ARCHITECTURE.md)

---

## 📈 Мониторинг

- **Prometheus:** http://151.241.228.203:9090
- **Grafana:** http://151.241.228.203:3001 (admin/admin123)
- **Flower:** http://151.241.228.203:5555 (Celery)
- **RabbitMQ:** http://151.241.228.203:15672 (guest/guest)

---

## 💰 Стоимость

- **VPS:** $22-40/мес (Hetzner/Vultr)
- **OpenRouter API:** $5-50/мес
- **Домен:** $10/год

**Итого:** ~$30-50/месяц на старт

---

## 🤝 Contributing

1. Fork репозиторий
2. Создайте feature branch
3. Commit изменения
4. Push в branch
5. Создайте Pull Request

---

## 📄 License

MIT License - см. [LICENSE](LICENSE)

---

## 📞 Контакты

- 📧 Email: support@newshub.ai
- 💬 Telegram: [@newshub_support](https://t.me/newshub_support)
- 🐛 Issues: [GitHub Issues](https://github.com/glifindor/newsportal/issues)

---

<p align="center">
  <b>Сделано с ❤️ и ☕ командой NewsHub AI</b>
</p>

- **GitHub:** https://github.com/glifindor/newsportal
- **Telegram:** @crypto_ainews, @kremlin_digest

## 📄 Лицензия

MIT
