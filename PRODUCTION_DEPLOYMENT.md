# 🚀 PRODUCTION DEPLOYMENT GUIDE - NewsHub AI

## Полное руководство по деплою в production с автоматизацией

---

## 📋 Содержание

1. [Обзор инфраструктуры](#обзор-инфраструктуры)
2. [Требования](#требования)
3. [Быстрый старт (автоматический)](#быстрый-старт-автоматический)
4. [Ручная установка](#ручная-установка)
5. [CI/CD с GitHub Actions](#cicd-с-github-actions)
6. [Мониторинг и алерты](#мониторинг-и-алерты)
7. [Backup стратегия](#backup-стратегия)
8. [SSL/HTTPS настройка](#sslhttps-настройка)
9. [Масштабирование](#масштабирование)
10. [Безопасность](#безопасность)
11. [Оценка стоимости](#оценка-стоимости)
12. [Troubleshooting](#troubleshooting)

---

## 🏗️ Обзор инфраструктуры

### Архитектура

```
┌─────────────────────────────────────────────────────────────┐
│                         Internet                             │
└───────────────────────┬─────────────────────────────────────┘
                        │
                   ┌────▼────┐
                   │  Nginx  │ (Reverse Proxy + SSL + Rate Limit)
                   │  :80/443│
                   └────┬────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
   ┌────▼────┐    ┌────▼────┐    ┌────▼────┐
   │Frontend │    │ Backend │    │   API   │
   │ Next.js │    │ FastAPI │    │  Docs   │
   │  :3000  │    │  :8000  │    │  /docs  │
   └─────────┘    └────┬────┘    └─────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
   ┌────▼────┐    ┌───▼────┐    ┌───▼────┐
   │ Postgres│    │  Redis │    │RabbitMQ│
   │  :5432  │    │  :6379 │    │  :5672 │
   └─────────┘    └────────┘    └────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
   ┌────▼────┐    ┌───▼────┐    ┌───▼────┐
   │ Celery  │    │ Celery │    │ Flower │
   │ Worker  │    │  Beat  │    │ :5555  │
   └─────────┘    └────────┘    └────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
   ┌────▼────┐    ┌───▼────┐    ┌───▼────┐
   │Prometheus│   │ Grafana│    │ Cadvisor│
   │  :9090  │    │ :3001  │    │  :8080 │
   └─────────┘    └────────┘    └────────┘
```

### Технологический стек

**Backend:**
- FastAPI 0.104+ (Python 3.11)
- PostgreSQL 15 (Database)
- Redis 7 (Cache & Sessions)
- RabbitMQ 3.12 (Message Broker)
- Celery 5.3 (Task Queue)

**Frontend:**
- Next.js 14 (React 18, TypeScript)
- Tailwind CSS 3
- TanStack Query & Table
- NextAuth.js

**Infrastructure:**
- Docker & Docker Compose
- Nginx (Reverse Proxy)
- Let's Encrypt (SSL)
- Prometheus & Grafana (Monitoring)
- GitHub Actions (CI/CD)

---

## ✅ Требования

### Сервер

- **OS:** Ubuntu 22.04 LTS или новее
- **RAM:** Минимум 4 GB (рекомендуется 8 GB)
- **CPU:** 2+ cores (рекомендуется 4+ cores)
- **Disk:** 50 GB+ SSD
- **Network:** Публичный IP, открытые порты 80, 443, 22

### Локальная машина

- **Windows:** PowerShell 5.1+ или OpenSSH
- **Linux/Mac:** Bash, SSH client
- **Git:** Для клонирования репозитория
- **Putty/OpenSSH:** Для SSH соединения

### API Ключи

- ✅ OpenRouter API Key (для AI обработки)
- ✅ Telegram Bot Token (для постинга)
- ✅ NewsAPI Key (опционально)
- ✅ Freepik API Key (опционально)

---

## 🎯 Быстрый старт (автоматический)

### Вариант 1: PowerShell скрипт (Windows)

**Полная автоматическая установка с нуля:**

```powershell
# 1. Клонировать репозиторий
git clone https://github.com/glifindor/newsportal.git
cd newsportal

# 2. Запустить интерактивную установку
.\scripts\setup-interactive.ps1

# Скрипт запросит:
# - Пароль от сервера (root@151.241.228.203)
# - API ключи (OpenRouter, Telegram, etc.)
# - Пароли для БД и Redis
# - Данные администратора

# 3. Подождать 10-15 минут
# Скрипт автоматически:
# - Обновит систему
# - Установит Docker и Docker Compose
# - Настроит Firewall
# - Склонирует проект
# - Создаст .env файл
# - Соберет Docker images
# - Запустит все контейнеры
# - Создаст администратора
```

**После завершения:**

```powershell
# Откройте в браузере:
http://151.241.228.203       # Frontend
http://151.241.228.203/docs  # API Docs
http://151.241.228.203:3001  # Grafana
```

---

### Вариант 2: Bash скрипт (Linux/Mac)

```bash
# 1. Клонировать репозиторий
git clone https://github.com/glifindor/newsportal.git
cd newsportal

# 2. Сделать скрипты исполняемыми
chmod +x scripts/*.sh

# 3. Запустить интерактивную установку
./scripts/setup-interactive.sh

# Следуйте инструкциям на экране
```

---

## 📖 Ручная установка

### Шаг 1: Подключение к серверу

```bash
ssh root@151.241.228.203
```

### Шаг 2: Обновление системы

```bash
apt update && apt upgrade -y
```

### Шаг 3: Установка Docker

```bash
# Установка Docker
apt install -y docker.io

# Запуск и автозагрузка
systemctl start docker
systemctl enable docker

# Проверка
docker --version
```

### Шаг 4: Установка Docker Compose

```bash
apt install -y docker-compose

# Проверка
docker-compose --version
```

### Шаг 5: Установка дополнительных инструментов

```bash
apt install -y git nano curl wget htop ufw jq
```

### Шаг 6: Настройка Firewall

```bash
# Разрешить SSH, HTTP, HTTPS
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp

# Включить firewall
ufw enable

# Проверка
ufw status
```

### Шаг 7: Клонирование проекта

```bash
# Создать директорию
mkdir -p /opt/newshub
cd /opt/newshub

# Клонировать репозиторий
git clone https://github.com/glifindor/newsportal.git .

# Проверка
ls -la
```

### Шаг 8: Создание .env файла

```bash
# Копировать шаблон
cp .env.example .env

# Редактировать
nano .env
```

**Заполнить значения:**

```env
# Database
POSTGRES_USER=newsadmin
POSTGRES_PASSWORD=ВАШ_СЛОЖНЫЙ_ПАРОЛЬ_1
POSTGRES_DB=newshub_db

# Redis
REDIS_PASSWORD=ВАШ_СЛОЖНЫЙ_ПАРОЛЬ_2

# RabbitMQ
RABBITMQ_USER=newshub
RABBITMQ_PASS=ВАШ_СЛОЖНЫЙ_ПАРОЛЬ_3

# JWT (сгенерировать: openssl rand -hex 32)
JWT_SECRET_KEY=ВАШ_JWT_SECRET_32_СИМВОЛА

# APIs
OPENROUTER_API_KEY=sk-or-v1-ваш_ключ
FREEPIK_API_KEY=ваш_ключ
NEWSAPI_KEY=ваш_ключ

# Telegram
TELEGRAM_BOT_TOKEN=123456:ABCdef_ваш_токен
TELEGRAM_CRYPTO_CHANNEL=@crypto_ainews
TELEGRAM_POLITICS_CHANNEL=@kremlin_digest
TELEGRAM_ADMIN_CHAT_ID=ваш_chat_id

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

### Шаг 9: Генерация JWT Secret

```bash
openssl rand -hex 32
# Скопировать вывод и вставить в .env как JWT_SECRET_KEY
```

### Шаг 10: Сборка и запуск

```bash
# Сборка images (5-15 минут)
docker-compose -f docker-compose.prod.yml build

# Запуск всех контейнеров
docker-compose -f docker-compose.prod.yml up -d

# Проверка статуса
docker-compose ps
```

### Шаг 11: Миграции БД

```bash
# Подождать 30 секунд для инициализации БД
sleep 30

# Запустить миграции
docker-compose exec backend alembic upgrade head
```

### Шаг 12: Создание администратора

```bash
docker-compose exec backend python scripts/create_admin.py

# Ввести:
# Username: admin
# Email: your@email.com
# Password: (надежный пароль)
```

### Шаг 13: Проверка работы

```bash
# Health check
curl http://151.241.228.203/health

# Ожидаемый ответ:
# {"status":"healthy","database":"connected","redis":"connected"}
```

---

## 🤖 CI/CD с GitHub Actions

### Настройка GitHub Secrets

Перейдите в **Settings → Secrets and variables → Actions** и добавьте:

```yaml
# Docker Hub
DOCKER_HUB_USERNAME: ваш_username
DOCKER_HUB_TOKEN: ваш_token

# SSH
SSH_PRIVATE_KEY: ваш_приватный_ключ

# Database
POSTGRES_USER: newsadmin
POSTGRES_PASSWORD: пароль_бд
POSTGRES_DB: newshub_db

# Redis
REDIS_PASSWORD: пароль_redis

# RabbitMQ
RABBITMQ_USER: newshub
RABBITMQ_PASS: пароль_rabbitmq

# JWT
JWT_SECRET_KEY: ваш_jwt_secret

# APIs
OPENROUTER_API_KEY: ваш_ключ
FREEPIK_API_KEY: ваш_ключ
NEWSAPI_KEY: ваш_ключ

# Telegram
TELEGRAM_BOT_TOKEN: ваш_токен
TELEGRAM_CRYPTO_CHANNEL: @crypto_ainews
TELEGRAM_POLITICS_CHANNEL: @kremlin_digest
TELEGRAM_ADMIN_CHAT_ID: ваш_id

# Monitoring
GRAFANA_USER: admin
GRAFANA_PASSWORD: пароль
FLOWER_USER: admin
FLOWER_PASSWORD: пароль
```

### Workflow Pipeline

Файл `.github/workflows/deploy.yml` уже создан и включает:

1. **Test Backend** - pytest, flake8, black
2. **Test Frontend** - npm test, lint, type-check
3. **Build & Push** - Docker images в Docker Hub
4. **Deploy** - Автоматический деплой на сервер
5. **Smoke Tests** - Проверка после деплоя

### Запуск деплоя

```bash
# Автоматически при push в main:
git push origin main

# Или вручную через GitHub:
Actions → Deploy → Run workflow
```

---

## 📊 Мониторинг и алерты

### Prometheus

**URL:** http://151.241.228.203:9090

**Метрики:**

- `up{job="backend"}` - статус backend
- `http_requests_total` - количество запросов
- `http_request_duration_seconds` - latency
- `celery_task_failed_total` - ошибки Celery
- `container_memory_usage_bytes` - использование памяти
- `node_cpu_seconds_total` - использование CPU

### Grafana

**URL:** http://151.241.228.203:3001  
**Логин:** admin / admin123

**Dashboards:**

1. **System Overview** - CPU, Memory, Disk, Network
2. **Application Metrics** - Requests, Latency, Errors
3. **Database** - Connections, Queries, Slow queries
4. **Celery** - Tasks, Queue, Workers
5. **Containers** - Docker stats

### Алерты в Telegram

**Настройка:**

1. Алерты автоматически отправляются в `TELEGRAM_ADMIN_CHAT_ID`
2. Критические алерты:
   - API Down (2+ минут)
   - Database Down
   - High error rate (>5%)
   - Disk space <5%

**Пример сообщения:**

```
🚨 ALERT: APIDown
Severity: critical
Component: backend

Backend API has been down for more than 2 minutes

Time: 2025-01-18 15:30:00
```

### Flower (Celery мониторинг)

**URL:** http://151.241.228.203:5555  
**Логин:** admin / admin123

**Функции:**

- Real-time мониторинг Celery workers
- Статистика задач
- Контроль очередей
- Перезапуск задач

---

## 💾 Backup стратегия

### Автоматические backup'ы

**Скрипт:** `/opt/newshub/scripts/backup.sh`

**Что сохраняется:**

- PostgreSQL database (полный dump)
- Redis data (RDB файл)
- Логи приложения
- Конфигурации (nginx, docker-compose)

**Расписание:** Каждый день в 3:00 AM

**Хранение:** 7 дней

### Настройка cron

```bash
# На сервере
cd /opt/newshub
chmod +x scripts/backup.sh
chmod +x scripts/setup-cron.sh

# Установка cron задачи
sudo ./scripts/setup-cron.sh
```

### Ручной backup

```bash
# На сервере
/opt/newshub/scripts/backup.sh
```

### Восстановление из backup

```bash
# PostgreSQL
gunzip /opt/newshub/backups/database/postgres_latest.sql.gz
cat /opt/newshub/backups/database/postgres_latest.sql | \
  docker-compose exec -T postgres psql -U newsadmin newshub_db

# Redis
docker cp /opt/newshub/backups/redis/redis_latest.rdb newshub_redis:/data/dump.rdb
docker-compose restart redis
```

### Offsite backup (рекомендуется)

```bash
# Установка rclone для S3/Google Drive
curl https://rclone.org/install.sh | sudo bash

# Настройка remote
rclone config

# Синхронизация backup'ов
rclone sync /opt/newshub/backups remote:newshub-backups
```

---

## 🔒 SSL/HTTPS настройка

### Вариант 1: Let's Encrypt (с доменом)

```bash
# Остановить Nginx
docker-compose stop nginx

# Установить Certbot
apt install -y certbot

# Получить сертификат (ЗАМЕНИТЕ домен!)
certbot certonly --standalone -d newshub.example.com

# Сертификаты будут в:
# /etc/letsencrypt/live/newshub.example.com/

# Обновить nginx.prod.conf (раскомментировать HTTPS server block)
nano nginx/nginx.prod.conf

# Запустить Nginx
docker-compose start nginx
```

### Вариант 2: Self-signed (для тестирования)

```bash
# Создать директорию
mkdir -p nginx/ssl

# Генерация сертификата
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/nginx-selfsigned.key \
  -out nginx/ssl/nginx-selfsigned.crt

# Перезапустить Nginx
docker-compose restart nginx
```

### Авто-обновление Let's Encrypt

```bash
# Добавить в crontab
crontab -e

# Добавить строку:
0 0 1 * * certbot renew --quiet && docker-compose restart nginx
```

---

## ⚡ Масштабирование

### Horizontal scaling

**Увеличение Celery workers:**

```yaml
# docker-compose.prod.yml
celery_worker:
  deploy:
    replicas: 4  # Увеличить количество
```

**Load balancing backend:**

```yaml
backend:
  deploy:
    replicas: 3
```

### Vertical scaling

**Увеличение ресурсов контейнеров:**

```yaml
backend:
  deploy:
    resources:
      limits:
        cpus: '4'    # Было 2
        memory: 4G   # Было 2G
```

### Database optimization

```sql
-- Индексы для быстрого поиска
CREATE INDEX idx_news_created_at ON news(created_at DESC);
CREATE INDEX idx_news_status ON news(status);
CREATE INDEX idx_news_channel ON news(channel);

-- Партиционирование по дате
CREATE TABLE news_2025_01 PARTITION OF news
  FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
```

### CDN для статики

**CloudFlare / AWS CloudFront:**

1. Зарегистрироваться на CloudFlare
2. Добавить домен
3. Настроить DNS
4. Включить Caching для `/_next/static/`

---

## 🛡️ Безопасность

### 1. Firewall (UFW)

```bash
# Только необходимые порты
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp   # SSH
ufw allow 80/tcp   # HTTP
ufw allow 443/tcp  # HTTPS
ufw enable
```

### 2. SSH hardening

```bash
# Отключить root login
nano /etc/ssh/sshd_config

# Изменить:
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes

# Перезапустить SSH
systemctl restart sshd
```

### 3. Fail2Ban

```bash
# Установка
apt install -y fail2ban

# Настройка
cat > /etc/fail2ban/jail.local << EOF
[sshd]
enabled = true
port = 22
maxretry = 3
bantime = 3600
EOF

# Запуск
systemctl enable fail2ban
systemctl start fail2ban
```

### 4. Security headers (Nginx)

Уже настроено в `nginx.prod.conf`:

```nginx
add_header Strict-Transport-Security "max-age=31536000" always;
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
```

### 5. Rate limiting

Настроено в Nginx:

- API: 20 req/sec
- Auth: 5 req/min
- Web: 50 req/sec

### 6. Secrets management

```bash
# НЕ коммитить .env в Git!
echo ".env" >> .gitignore

# Использовать Docker secrets (рекомендуется)
docker secret create postgres_password /path/to/password.txt
```

---

## 💰 Оценка стоимости

### Вариант 1: VPS (151.241.228.203)

**Конфигурация:**
- 4 vCPU, 8 GB RAM, 100 GB SSD
- Bandwidth: 1 TB/month

**Провайдеры:**

| Провайдер | Стоимость/месяц | Примечание |
|-----------|-----------------|------------|
| Hetzner | €20 ($22) | Европа |
| DigitalOcean | $48 | Droplet 8GB |
| Vultr | $40 | High Performance |
| AWS EC2 (t3.large) | $60 | On-demand |
| Linode | $40 | Dedicated CPU |

### Вариант 2: Managed Services

**Backend:** AWS ECS Fargate
- 2 tasks × $30 = $60/month

**Database:** AWS RDS PostgreSQL
- db.t3.medium = $75/month

**Frontend:** Vercel Pro
- $20/month

**CDN:** CloudFlare
- Free tier

**Итого:** ~$155/month

### Вариант 3: Full AWS

- EC2 t3.large: $60
- RDS PostgreSQL: $75
- ElastiCache Redis: $50
- S3 + CloudFront: $10
- Load Balancer: $20

**Итого:** ~$215/month

### Дополнительные расходы

- **Домен:** $10-15/год
- **SSL:** Free (Let's Encrypt)
- **APIs:**
  - OpenRouter: $5-50/month (зависит от использования)
  - NewsAPI: Free tier или $449/month (Business)
  - Freepik: Free или $10/month

### Рекомендация для старта

**VPS (Hetzner/Vultr):** $22-40/month
- Достаточно для 1000-10000 пользователей/день
- Полный контроль
- Простое масштабирование

**Итого на старт:** ~$50/месяц

---

## 🔧 Troubleshooting

### Проблема: Port 80 уже занят

```bash
# Найти процесс
lsof -i :80

# Остановить Apache (если установлен)
systemctl stop apache2
systemctl disable apache2

# Перезапустить Nginx
docker-compose restart nginx
```

### Проблема: Не могу подключиться к БД

```bash
# Проверить статус PostgreSQL
docker-compose exec postgres pg_isready -U newsadmin

# Посмотреть логи
docker-compose logs postgres

# Перезапустить
docker-compose restart postgres

# Проверить переменные окружения
docker-compose exec backend env | grep DATABASE_URL
```

### Проблема: Telegram Bot не отвечает

1. Проверить токен в `.env`
2. Убедиться, что бот добавлен в каналы как админ
3. Отправить боту `/start` в личку
4. Проверить `TELEGRAM_ADMIN_CHAT_ID`

```bash
# Тест Telegram
docker-compose exec backend python scripts/test_telegram.py
```

### Проблема: OpenRouter ошибка 401

1. Проверить API ключ в `.env`
2. Убедиться, что есть баланс на openrouter.ai
3. Попробовать другую модель (более дешевую)

```bash
# Тест OpenRouter
docker-compose exec backend python scripts/test_openrouter.py
```

### Проблема: Out of memory

```bash
# Проверить RAM
free -h

# Добавить swap
fallocate -l 4G /swapfile
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile

# Сделать постоянным
echo '/swapfile none swap sw 0 0' >> /etc/fstab
```

### Проблема: Docker image не собирается

```bash
# Очистить Docker cache
docker system prune -a

# Пересобрать с нуля
docker-compose -f docker-compose.prod.yml build --no-cache

# Проверить disk space
df -h
```

### Проблема: Celery worker не запускается

```bash
# Логи worker
docker-compose logs celery_worker

# Проверить RabbitMQ
docker-compose exec rabbitmq rabbitmq-diagnostics check_running

# Перезапустить worker
docker-compose restart celery_worker
```

### Проблема: Frontend не загружается

```bash
# Логи frontend
docker-compose logs frontend

# Пересобрать frontend
docker-compose -f docker-compose.prod.yml build frontend
docker-compose up -d frontend

# Проверить Nginx
docker-compose logs nginx
```

---

## 📞 Поддержка

### Логи

```bash
# Все логи
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f backend

# Последние 100 строк
docker-compose logs --tail=100 backend

# За последний час
docker-compose logs --since 1h backend
```

### Статус сервисов

```bash
# Статус контейнеров
docker-compose ps

# Использование ресурсов
docker stats

# Системные ресурсы
htop
```

### Полезные команды

```bash
# Перезапуск всех сервисов
docker-compose restart

# Перезапуск конкретного сервиса
docker-compose restart backend

# Остановка
docker-compose stop

# Удаление контейнеров (сохраняя данные)
docker-compose down

# Удаление всего (включая volumes!)
docker-compose down -v

# Обновление проекта
cd /opt/newshub
git pull
docker-compose -f docker-compose.prod.yml up -d --build
```

---

## ✅ Чеклист Production Readiness

- [ ] Сервер настроен и доступен
- [ ] Docker и Docker Compose установлены
- [ ] Проект склонирован
- [ ] .env файл создан и заполнен
- [ ] Все контейнеры запущены
- [ ] Миграции выполнены
- [ ] Администратор создан
- [ ] API отвечает (curl http://server/health)
- [ ] Frontend открывается в браузере
- [ ] Telegram Bot работает
- [ ] OpenRouter отвечает
- [ ] Firewall настроен
- [ ] SSL сертификат установлен (опционально)
- [ ] Backup'ы настроены
- [ ] Мониторинг работает (Grafana)
- [ ] Алерты настроены
- [ ] GitHub Actions CI/CD работает
- [ ] Документация обновлена

---

## 🎉 Готово!

Ваш **NewsHub AI** успешно развернут в production!

**Доступные сервисы:**

- 🌐 **Frontend:** http://151.241.228.203
- 📚 **API Docs:** http://151.241.228.203/docs
- 📊 **Grafana:** http://151.241.228.203:3001
- 🐰 **RabbitMQ:** http://151.241.228.203:15672
- 🌸 **Flower:** http://151.241.228.203:5555
- 📈 **Prometheus:** http://151.241.228.203:9090

**Следующие шаги:**

1. Настройте SSL для HTTPS
2. Настройте домен (опционально)
3. Проверьте backup'ы
4. Настройте алерты в Telegram
5. Оптимизируйте производительность
6. Масштабируйте при необходимости

**Документация:**

- ARCHITECTURE.md - Архитектура проекта
- README.md - Общая информация
- DEPLOYMENT.md - Базовая инструкция
- PRODUCTION_DEPLOYMENT.md - Этот файл

**Поддержка:**

- GitHub Issues: https://github.com/glifindor/newsportal/issues
- Email: support@newshub.ai
- Telegram: @newshub_support

---

**Удачи! 🚀**
