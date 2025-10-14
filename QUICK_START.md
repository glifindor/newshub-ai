# 🚀 Быстрый старт для сервера 151.241.228.203

## Метод 1: Автоматическая установка (РЕКОМЕНДУЕТСЯ)

### На вашей локальной машине:

```powershell
# 1. Упаковать проект
cd "C:\Users\Grifindor\Desktop\НОВСТНОЙ САЙТ"
tar -czf news-portal.tar.gz auth-service news-service media-service docker-compose.yml deploy

# 2. Загрузить на сервер
scp news-portal.tar.gz root@151.241.228.203:/root/
scp deploy/install.sh root@151.241.228.203:/root/
```

### На сервере:

```bash
# 1. Подключиться по SSH
ssh root@151.241.228.203

# 2. Запустить автоматическую установку
chmod +x install.sh
./install.sh

# 3. Распаковать проект
cd /opt/news-portal
tar -xzf /root/news-portal.tar.gz

# 4. Запустить сервисы
chmod +x deploy/*.sh
./deploy/start.sh
```

---

## Метод 2: Ручная установка

### 1. Подключение
```bash
ssh root@151.241.228.203
```

### 2. Установка Docker
```bash
# Обновление системы
apt update && apt upgrade -y

# Установка Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Установка Docker Compose
apt install docker-compose-plugin -y

# Проверка
docker --version
docker compose version
```

### 3. Установка зависимостей
```bash
apt install -y git nginx ufw
```

### 4. Настройка Firewall
```bash
ufw allow OpenSSH
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable
```

### 5. Создание директории проекта
```bash
mkdir -p /opt/news-portal
cd /opt/news-portal
```

### 6. Загрузка проекта

**Вариант A: С локальной машины**
```powershell
# На Windows (PowerShell)
cd "C:\Users\Grifindor\Desktop\НОВСТНОЙ САЙТ"
scp -r . root@151.241.228.203:/opt/news-portal/
```

**Вариант B: Из Git**
```bash
# На сервере
cd /opt/news-portal
git clone YOUR_REPOSITORY_URL .
```

### 7. Создание .env.production
```bash
cd /opt/news-portal

cat > .env.production << 'EOF'
ENVIRONMENT=production
SERVER_IP=151.241.228.203

# Сгенерируйте безопасные пароли!
POSTGRES_USER=newsportal
POSTGRES_PASSWORD=CHANGE_ME
POSTGRES_DB=news_portal

REDIS_PASSWORD=CHANGE_ME
MINIO_ROOT_USER=admin
MINIO_ROOT_PASSWORD=CHANGE_ME
JWT_SECRET=CHANGE_ME

AUTH_SERVICE_PORT=8091
NEWS_SERVICE_PORT=8092
MEDIA_SERVICE_PORT=8094
EOF

# Сгенерировать пароли
openssl rand -base64 32  # Используйте для POSTGRES_PASSWORD
openssl rand -base64 32  # Используйте для REDIS_PASSWORD
openssl rand -base64 32  # Используйте для MINIO_ROOT_PASSWORD
openssl rand -base64 64  # Используйте для JWT_SECRET
```

### 8. Запуск

```bash
# Запуск инфраструктуры
docker compose up -d postgres redis minio

# Ожидание (30 секунд)
sleep 30

# Сборка сервисов
docker compose build auth-service news-service media-service

# Запуск сервисов
docker compose up -d auth-service news-service media-service

# Проверка
docker compose ps
docker compose logs -f
```

---

## Проверка работы

```bash
# Health checks
curl http://localhost:8091/health  # Auth Service
curl http://localhost:8092/health  # News Service
curl http://localhost:8094/health  # Media Service

# Проверка с внешнего IP
curl http://151.241.228.203:8091/health
```

---

## Полезные команды

```bash
# Просмотр логов
docker compose logs -f auth-service
docker compose logs -f news-service
docker compose logs -f media-service

# Статус контейнеров
docker compose ps

# Перезапуск сервиса
docker compose restart auth-service

# Остановка всех сервисов
docker compose down

# Обновление
docker compose build
docker compose up -d
```

---

## Endpoints после запуска

### HTTP API:
- **Auth Service:** http://151.241.228.203:8091
  - POST /api/v1/auth/register
  - POST /api/v1/auth/login
  - GET  /api/v1/auth/profile

- **News Service:** http://151.241.228.203:8092
  - GET  /api/v1/news
  - GET  /api/v1/categories
  - GET  /api/v1/tags

- **Media Service:** http://151.241.228.203:8094
  - POST /api/v1/media/upload
  - GET  /api/v1/media

### Инфраструктура:
- **PostgreSQL:** localhost:5432
- **Redis:** localhost:6379
- **MinIO Console:** http://151.241.228.203:9001

---

## Troubleshooting

### Контейнеры не запускаются
```bash
docker compose logs
docker system df
free -h
```

### Порты заняты
```bash
sudo netstat -tulpn | grep 8091
sudo lsof -i :8091
```

### Проблемы с БД
```bash
docker compose exec postgres psql -U postgres
# \l  - список баз
# \dt - список таблиц
```

### Очистка и перезапуск
```bash
docker compose down -v
docker system prune -a
docker compose up -d
```

---

## 🔒 Безопасность (ВАЖНО!)

После запуска обязательно:

1. **Измените все пароли** в `.env.production`
2. **Настройте SSL** (см. DEPLOYMENT_GUIDE.md)
3. **Настройте Nginx** как reverse proxy
4. **Включите Fail2Ban**
5. **Настройте регулярные бэкапы**

---

## 📚 Полная документация

Подробная инструкция в файле: **DEPLOYMENT_GUIDE.md**

---

**IP сервера:** 151.241.228.203  
**Дата:** 14 октября 2025
