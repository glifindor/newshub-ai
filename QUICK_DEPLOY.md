# 🚀 БЫСТРЫЙ ДЕПЛОЙ - NewsHub AI

## Автоматическая установка за 15 минут

---

## 📋 Что нужно подготовить

### 1. API Ключи

- ✅ **OpenRouter API Key** - https://openrouter.ai/keys
- ✅ **Telegram Bot Token** - https://t.me/BotFather
- ✅ **Telegram Admin Chat ID** - https://t.me/userinfobot (отправь /start)
- ⚪ NewsAPI Key (опционально) - https://newsapi.org
- ⚪ Freepik API Key (опционально) - https://freepik.com

### 2. Сервер

- ✅ **IP:** 151.241.228.203
- ✅ **Пользователь:** root
- ✅ **Пароль:** (у вас должен быть)

---

## ⚡ Вариант 1: АВТОМАТИЧЕСКАЯ УСТАНОВКА (Рекомендуется)

### Для Windows (PowerShell)

```powershell
# 1. Клонировать проект
git clone https://github.com/glifindor/newsportal.git
cd newsportal

# 2. Запустить автоматическую установку
.\scripts\setup-interactive.ps1

# Скрипт попросит:
# - Пароль от сервера
# - API ключи
# - Пароли для БД
# - Email и пароль администратора

# 3. Подождать 10-15 минут ☕
# Скрипт автоматически:
# ✅ Обновит систему
# ✅ Установит Docker
# ✅ Настроит Firewall
# ✅ Клонирует проект на сервер
# ✅ Создаст .env файл
# ✅ Соберет Docker images
# ✅ Запустит все контейнеры
# ✅ Создаст администратора

# 4. Готово! 🎉
```

**После завершения откройте в браузере:**

```
http://151.241.228.203       → Frontend
http://151.241.228.203/docs  → API Documentation
http://151.241.228.203:3001  → Grafana (Мониторинг)
```

---

## 🔧 Вариант 2: РУЧНАЯ УСТАНОВКА

### Шаг 1: Подключение к серверу

```powershell
# Windows (PowerShell)
ssh root@151.241.228.203
# Введите пароль

# Или используйте PuTTY
# Host: 151.241.228.203
# User: root
```

### Шаг 2: Быстрая установка на сервере

```bash
# Скопируйте и вставьте в терминал:

# 1. Обновление и установка Docker
apt update && apt upgrade -y
apt install -y docker.io docker-compose git nano curl wget ufw

# 2. Настройка Firewall
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

# 3. Клонирование проекта
mkdir -p /opt/newshub
cd /opt/newshub
git clone https://github.com/glifindor/newsportal.git .

# 4. Создание .env файла
cp .env.example .env
nano .env
```

### Шаг 3: Заполнить .env файл

**Нажмите в nano, вставьте и измените значения:**

```env
# Database
POSTGRES_USER=newsadmin
POSTGRES_PASSWORD=YOUR_STRONG_PASSWORD_1
POSTGRES_DB=newshub_db

# Redis
REDIS_PASSWORD=YOUR_STRONG_PASSWORD_2

# RabbitMQ
RABBITMQ_USER=newshub
RABBITMQ_PASS=YOUR_STRONG_PASSWORD_3

# JWT Secret (сгенерировать: openssl rand -hex 32)
JWT_SECRET_KEY=YOUR_32_CHARS_JWT_SECRET

# APIs
OPENROUTER_API_KEY=sk-or-v1-YOUR_KEY
FREEPIK_API_KEY=YOUR_KEY
NEWSAPI_KEY=YOUR_KEY

# Telegram
TELEGRAM_BOT_TOKEN=123456:ABCdef_YOUR_TOKEN
TELEGRAM_CRYPTO_CHANNEL=@crypto_ainews
TELEGRAM_POLITICS_CHANNEL=@kremlin_digest
TELEGRAM_ADMIN_CHAT_ID=YOUR_CHAT_ID

# Monitoring
GRAFANA_USER=admin
GRAFANA_PASSWORD=admin123
FLOWER_USER=admin
FLOWER_PASSWORD=admin123

# Frontend
NEXT_PUBLIC_API_URL=http://151.241.228.203/api
NEXT_PUBLIC_WS_URL=ws://151.241.228.203

# Environment
ENVIRONMENT=production
LOG_LEVEL=INFO
```

**Сохранить:** `Ctrl+X` → `Y` → `Enter`

### Шаг 4: Запуск

```bash
# Сборка и запуск (5-15 минут)
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d

# Подождать 30 секунд
sleep 30

# Миграции БД
docker-compose exec backend alembic upgrade head

# Проверка статуса
docker-compose ps
```

### Шаг 5: Создание администратора

```bash
docker-compose exec backend python scripts/create_admin.py

# Ввести:
# Username: admin
# Email: your@email.com
# Password: (ваш пароль)
```

### Шаг 6: Проверка

```bash
# Health check
curl http://151.241.228.203/health

# Должен вернуть:
# {"status":"healthy","database":"connected","redis":"connected"}
```

**Готово! 🎉**

---

## 🔄 Обновление проекта

### Автоматическое (из Windows)

```powershell
# Запустить скрипт деплоя
.\scripts\deploy.ps1 -Password "ваш_пароль"

# Или отдельные команды:
.\scripts\deploy.ps1 -Action status -Password "пароль"   # Статус
.\scripts\deploy.ps1 -Action restart -Password "пароль"  # Перезапуск
.\scripts\deploy.ps1 -Action logs -Password "пароль"     # Логи
```

### Ручное (на сервере)

```bash
cd /opt/newshub

# Pull изменений
git pull origin main

# Пересобрать и перезапустить
docker-compose -f docker-compose.prod.yml up -d --build

# Миграции (если нужны)
docker-compose exec backend alembic upgrade head
```

---

## 📊 Полезные команды

### Просмотр логов

```bash
# Все логи
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f backend

# Последние 100 строк
docker-compose logs --tail=100 backend
```

### Статус сервисов

```bash
# Статус контейнеров
docker-compose ps

# Использование ресурсов
docker stats
```

### Перезапуск

```bash
# Перезапустить всё
docker-compose restart

# Конкретный сервис
docker-compose restart backend
```

### Backup

```bash
# Ручной backup
/opt/newshub/scripts/backup.sh

# Настроить автоматический (каждый день в 3:00)
chmod +x /opt/newshub/scripts/setup-cron.sh
/opt/newshub/scripts/setup-cron.sh
```

---

## 🔍 Проверка работы

### 1. Frontend

```
http://151.241.228.203
```

Должна открыться главная страница с новостями

### 2. API Documentation

```
http://151.241.228.203/docs
```

Swagger UI с документацией API

### 3. Grafana (Мониторинг)

```
http://151.241.228.203:3001
Логин: admin
Пароль: admin123
```

Графики CPU, Memory, Requests

### 4. RabbitMQ (Очередь задач)

```
http://151.241.228.203:15672
Логин: guest
Пароль: guest
```

### 5. Flower (Celery мониторинг)

```
http://151.241.228.203:5555
Логин: admin
Пароль: admin123
```

---

## ❓ Частые проблемы

### "Cannot connect to server"

```bash
# Проверить, что сервер доступен
ping 151.241.228.203

# Проверить SSH
ssh root@151.241.228.203 echo "OK"
```

### "Port 80 already in use"

```bash
# Найти процесс
lsof -i :80

# Остановить Apache
systemctl stop apache2
systemctl disable apache2
```

### "Database connection failed"

```bash
# Проверить PostgreSQL
docker-compose exec postgres pg_isready -U newsadmin

# Логи
docker-compose logs postgres

# Перезапустить
docker-compose restart postgres
```

### "Telegram bot not responding"

1. Проверьте TELEGRAM_BOT_TOKEN в .env
2. Убедитесь, что бот добавлен в каналы как админ
3. Отправьте /start боту в личку
4. Проверьте TELEGRAM_ADMIN_CHAT_ID

```bash
# Тест
docker-compose exec backend python scripts/test_telegram.py
```

### "Out of memory"

```bash
# Добавить swap
fallocate -l 4G /swapfile
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile
echo '/swapfile none swap sw 0 0' >> /etc/fstab
```

---

## 📚 Документация

### Полная документация

- **PRODUCTION_DEPLOYMENT.md** - Детальное руководство по деплою
- **ARCHITECTURE.md** - Архитектура проекта
- **README.md** - Общая информация
- **TELEGRAM_BOT_SETUP.md** - Настройка Telegram бота

### Скрипты

```
scripts/
├── setup-interactive.ps1   # Автоматическая установка (Windows)
├── deploy.ps1              # Деплой (Windows)
├── deploy.sh               # Деплой (Linux)
├── backup.sh               # Backup скрипт
└── setup-cron.sh           # Настройка автоматических backup'ов
```

---

## 🎯 Чеклист готовности

- [ ] Сервер доступен (ssh root@151.241.228.203)
- [ ] Docker установлен (docker --version)
- [ ] Проект склонирован (/opt/newshub)
- [ ] .env файл создан и заполнен
- [ ] Контейнеры запущены (docker-compose ps)
- [ ] API отвечает (curl http://151.241.228.203/health)
- [ ] Frontend открывается в браузере
- [ ] Администратор создан
- [ ] Telegram Bot работает
- [ ] Grafana доступна
- [ ] Backup'ы настроены

---

## 🚀 Следующие шаги

### 1. Настройте SSL (HTTPS)

```bash
# Let's Encrypt
docker-compose stop nginx
apt install -y certbot
certbot certonly --standalone -d newshub.example.com
# Обновить nginx.prod.conf
docker-compose start nginx
```

### 2. Настройте домен

1. Купите домен (Namecheap, GoDaddy)
2. Добавьте A-запись: newshub.example.com → 151.241.228.203
3. Обновите .env:
   ```env
   NEXT_PUBLIC_API_URL=https://newshub.example.com/api
   ```

### 3. Настройте CI/CD

1. Добавьте GitHub Secrets (см. PRODUCTION_DEPLOYMENT.md)
2. При push в main → автоматический деплой

### 4. Мониторинг

- Настройте алерты в Telegram
- Проверьте дашборды в Grafana
- Настройте Sentry (опционально)

---

## 💰 Стоимость

**Минимальная конфигурация:**

- VPS (4GB RAM): $22-40/мес
- OpenRouter API: $5-50/мес
- Домен: $10/год

**Итого:** ~$30-50/месяц на старт

---

## ✅ Готово!

Ваш **NewsHub AI** работает в production! 🎉

**Доступ:**

- 🌐 Frontend: http://151.241.228.203
- 📚 API: http://151.241.228.203/docs
- 📊 Grafana: http://151.241.228.203:3001

**Поддержка:**

- GitHub: https://github.com/glifindor/newsportal/issues
- Telegram: @newshub_support
- Email: support@newshub.ai

**Удачи! 🚀**
