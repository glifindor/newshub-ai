# 🚀 Деплой новостного портала на production сервер

**Сервер:** 151.241.228.203  
**Дата:** 14 октября 2025

---

## 📋 Предварительные требования

### На локальной машине:
- Git
- SSH клиент

### На сервере (будет установлено автоматически):
- Ubuntu 20.04+ / Debian 11+
- Docker 24+
- Docker Compose v2
- Nginx
- Certbot (для SSL)

---

## 🔧 Шаг 1: Подключение к серверу

```bash
# Подключение по SSH
ssh root@151.241.228.203

# Или если есть пользователь
ssh your_user@151.241.228.203
```

---

## 🛠 Шаг 2: Автоматическая установка (РЕКОМЕНДУЕТСЯ)

Скопируйте на сервер и запустите скрипт установки:

```bash
# На сервере
wget https://raw.githubusercontent.com/YOUR_REPO/main/deploy/install.sh
chmod +x install.sh
./install.sh
```

**ИЛИ** выполните ручную установку (следующие шаги).

---

## 📦 Шаг 3: Ручная установка зависимостей

### 3.1 Обновление системы
```bash
sudo apt update && sudo apt upgrade -y
```

### 3.2 Установка Docker
```bash
# Удаление старых версий
sudo apt remove docker docker-engine docker.io containerd runc

# Установка зависимостей
sudo apt install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# Добавление Docker GPG ключа
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Добавление Docker репозитория
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Установка Docker
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Проверка установки
docker --version
docker compose version
```

### 3.3 Установка дополнительных инструментов
```bash
sudo apt install -y git nginx certbot python3-certbot-nginx ufw
```

### 3.4 Настройка Firewall
```bash
# Разрешить SSH
sudo ufw allow OpenSSH

# Разрешить HTTP и HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Включить firewall
sudo ufw --force enable

# Проверка статуса
sudo ufw status
```

---

## 📥 Шаг 4: Загрузка проекта

### 4.1 Создание директории
```bash
sudo mkdir -p /opt/news-portal
sudo chown -R $USER:$USER /opt/news-portal
cd /opt/news-portal
```

### 4.2 Клонирование репозитория

**Вариант A: Из Git репозитория**
```bash
git clone https://github.com/YOUR_USERNAME/news-portal.git .
```

**Вариант B: Загрузка с локальной машины**
На вашей машине:
```powershell
# Упаковка проекта
cd "C:\Users\Grifindor\Desktop\НОВСТНОЙ САЙТ"
tar -czf news-portal.tar.gz .

# Загрузка на сервер
scp news-portal.tar.gz root@151.241.228.203:/opt/news-portal/
```

На сервере:
```bash
cd /opt/news-portal
tar -xzf news-portal.tar.gz
rm news-portal.tar.gz
```

---

## 🔐 Шаг 5: Настройка production конфигурации

### 5.1 Создание .env файла для production
```bash
cat > /opt/news-portal/.env.production << 'EOF'
# Environment
ENVIRONMENT=production

# Server
SERVER_IP=151.241.228.203
DOMAIN=your-domain.com

# Database
POSTGRES_USER=newsportal_user
POSTGRES_PASSWORD=CHANGE_ME_STRONG_PASSWORD_HERE
POSTGRES_DB=news_portal

# Redis
REDIS_PASSWORD=CHANGE_ME_REDIS_PASSWORD

# MinIO
MINIO_ROOT_USER=newsportal_admin
MINIO_ROOT_PASSWORD=CHANGE_ME_MINIO_PASSWORD

# RabbitMQ
RABBITMQ_DEFAULT_USER=newsportal_admin
RABBITMQ_DEFAULT_PASS=CHANGE_ME_RABBITMQ_PASSWORD

# JWT Secret (generate with: openssl rand -base64 64)
JWT_SECRET=CHANGE_ME_JWT_SECRET_64_CHARS_MINIMUM

# Auth Service
AUTH_SERVICE_PORT=8091
AUTH_GRPC_PORT=8081

# News Service
NEWS_SERVICE_PORT=8092
NEWS_GRPC_PORT=8082

# Media Service
MEDIA_SERVICE_PORT=8094

# Gateway
GATEWAY_PORT=8080

# Frontend
FRONTEND_PORT=3000
NEXT_PUBLIC_API_URL=https://api.your-domain.com

# SSL
SSL_EMAIL=your-email@example.com

# Monitoring
GRAFANA_ADMIN_PASSWORD=CHANGE_ME_GRAFANA_PASSWORD
EOF
```

### 5.2 Генерация безопасных паролей
```bash
# Генерация JWT секрета
openssl rand -base64 64

# Генерация случайных паролей
openssl rand -base64 32  # Для PostgreSQL
openssl rand -base64 32  # Для Redis
openssl rand -base64 32  # Для MinIO
openssl rand -base64 32  # Для RabbitMQ
```

**ВАЖНО:** Замените все `CHANGE_ME_*` значения в `.env.production`!

---

## 🐳 Шаг 6: Создание Dockerfiles

Создайте Dockerfile для каждого сервиса (если еще не созданы):

### Auth Service
```bash
cat > /opt/news-portal/auth-service/Dockerfile << 'EOF'
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/auth-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8091 8081
CMD ["./main"]
EOF
```

### News Service
```bash
cat > /opt/news-portal/news-service/Dockerfile << 'EOF'
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/news-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8092 8082
CMD ["./main"]
EOF
```

### Media Service
```bash
cat > /opt/news-portal/media-service/Dockerfile << 'EOF'
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/media-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8094
CMD ["./main"]
EOF
```

---

## 🚀 Шаг 7: Запуск проекта

### 7.1 Сборка образов
```bash
cd /opt/news-portal

# Сборка всех сервисов
docker compose -f docker-compose.yml --env-file .env.production build

# Или по отдельности
docker compose build auth-service
docker compose build news-service
docker compose build media-service
```

### 7.2 Запуск инфраструктуры
```bash
# Запуск PostgreSQL, Redis, MinIO
docker compose up -d postgres redis minio

# Ожидание готовности БД (30 секунд)
sleep 30

# Проверка статуса
docker compose ps
```

### 7.3 Запуск микросервисов
```bash
# Запуск Auth Service
docker compose up -d auth-service

# Запуск News Service
docker compose up -d news-service

# Запуск Media Service
docker compose up -d media-service

# Проверка логов
docker compose logs -f auth-service
docker compose logs -f news-service
docker compose logs -f media-service
```

### 7.4 Проверка работоспособности
```bash
# Health checks
curl http://localhost:8091/health  # Auth Service
curl http://localhost:8092/health  # News Service
curl http://localhost:8094/health  # Media Service

# Проверка всех контейнеров
docker compose ps
```

---

## 🌐 Шаг 8: Настройка Nginx (Reverse Proxy)

### 8.1 Создание конфигурации Nginx
```bash
sudo cat > /etc/nginx/sites-available/news-portal << 'EOF'
# Upstream серверы
upstream auth_backend {
    server localhost:8091;
}

upstream news_backend {
    server localhost:8092;
}

upstream media_backend {
    server localhost:8094;
}

upstream gateway {
    server localhost:8080;
}

# HTTP -> HTTPS редирект
server {
    listen 80;
    listen [::]:80;
    server_name your-domain.com www.your-domain.com;
    
    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }
    
    location / {
        return 301 https://$host$request_uri;
    }
}

# HTTPS сервер
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name your-domain.com www.your-domain.com;
    
    # SSL сертификаты (будут созданы Certbot)
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    
    # SSL настройки
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    
    # Безопасность
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    
    # Логи
    access_log /var/log/nginx/news-portal-access.log;
    error_log /var/log/nginx/news-portal-error.log;
    
    # Максимальный размер загружаемого файла
    client_max_body_size 10M;
    
    # API Gateway
    location /api/ {
        proxy_pass http://gateway;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
    
    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
    
    # MinIO (опционально для прямого доступа)
    location /minio/ {
        proxy_pass http://localhost:9000/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# API поддомен (опционально)
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name api.your-domain.com;
    
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    
    location / {
        proxy_pass http://gateway;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF
```

### 8.2 Активация конфигурации
```bash
# Создание символической ссылки
sudo ln -s /etc/nginx/sites-available/news-portal /etc/nginx/sites-enabled/

# Удаление дефолтной конфигурации
sudo rm -f /etc/nginx/sites-enabled/default

# Проверка конфигурации
sudo nginx -t

# Перезагрузка Nginx
sudo systemctl reload nginx
```

---

## 🔒 Шаг 9: Настройка SSL (HTTPS)

### 9.1 Получение SSL сертификата (Let's Encrypt)
```bash
# Замените your-domain.com и your-email@example.com
sudo certbot --nginx -d your-domain.com -d www.your-domain.com --email your-email@example.com --agree-tos --no-eff-email

# Автоматическое обновление
sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer
```

### 9.2 Тест обновления сертификата
```bash
sudo certbot renew --dry-run
```

---

## 📊 Шаг 10: Мониторинг и логи

### 10.1 Просмотр логов
```bash
# Логи конкретного сервиса
docker compose logs -f auth-service
docker compose logs -f news-service
docker compose logs -f media-service

# Все логи
docker compose logs -f

# Последние 100 строк
docker compose logs --tail=100 auth-service

# Логи Nginx
sudo tail -f /var/log/nginx/news-portal-access.log
sudo tail -f /var/log/nginx/news-portal-error.log
```

### 10.2 Мониторинг ресурсов
```bash
# Статистика контейнеров
docker stats

# Использование диска
docker system df

# Список контейнеров
docker compose ps
```

### 10.3 Grafana (опционально)
```bash
# Запуск мониторинга
docker compose up -d prometheus grafana

# Доступ к Grafana
# http://151.241.228.203:3001
# Login: admin / Password: из .env.production
```

---

## 🔄 Шаг 11: Управление сервисами

### Остановка
```bash
cd /opt/news-portal

# Остановка всех сервисов
docker compose down

# Остановка с удалением volumes (ОСТОРОЖНО!)
docker compose down -v
```

### Перезапуск
```bash
# Перезапуск всех сервисов
docker compose restart

# Перезапуск конкретного сервиса
docker compose restart auth-service
docker compose restart news-service
```

### Обновление
```bash
# Pull новых изменений
git pull origin main

# Пересборка и перезапуск
docker compose build
docker compose up -d
```

### Очистка
```bash
# Удаление неиспользуемых образов
docker image prune -a

# Удаление неиспользуемых volumes
docker volume prune

# Полная очистка
docker system prune -a --volumes
```

---

## 🧪 Шаг 12: Тестирование

### 12.1 Проверка API
```bash
# Регистрация пользователя
curl -X POST https://your-domain.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "AdminPass123",
    "full_name": "Admin User",
    "role": "admin"
  }'

# Вход
curl -X POST https://your-domain.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "AdminPass123"
  }'
```

### 12.2 Health Checks
```bash
curl https://your-domain.com/api/v1/auth/health
curl https://your-domain.com/api/v1/news/health
curl https://your-domain.com/api/v1/media/health
```

---

## 🛡 Шаг 13: Безопасность

### 13.1 Изменение SSH порта (рекомендуется)
```bash
sudo nano /etc/ssh/sshd_config
# Измените Port 22 на Port 2222

sudo systemctl restart sshd

# Не забудьте открыть новый порт в firewall!
sudo ufw allow 2222/tcp
```

### 13.2 Отключение root логина
```bash
sudo nano /etc/ssh/sshd_config
# PermitRootLogin no

sudo systemctl restart sshd
```

### 13.3 Fail2Ban (защита от брутфорса)
```bash
sudo apt install fail2ban -y
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

---

## 📋 Шаг 14: Автоматический запуск при перезагрузке

### 14.1 Создание systemd service
```bash
sudo cat > /etc/systemd/system/news-portal.service << 'EOF'
[Unit]
Description=News Portal Docker Compose
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/news-portal
ExecStart=/usr/bin/docker compose up -d
ExecStop=/usr/bin/docker compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
EOF
```

### 14.2 Активация сервиса
```bash
sudo systemctl daemon-reload
sudo systemctl enable news-portal.service
sudo systemctl start news-portal.service

# Проверка статуса
sudo systemctl status news-portal.service
```

---

## 🔍 Troubleshooting

### Проблема: Контейнеры не запускаются
```bash
# Проверка логов
docker compose logs

# Проверка ресурсов
docker stats
free -h
df -h
```

### Проблема: База данных недоступна
```bash
# Проверка PostgreSQL
docker compose exec postgres psql -U postgres -c "\l"

# Пересоздание БД
docker compose down
docker volume rm news-portal_postgres_data
docker compose up -d postgres
```

### Проблема: 502 Bad Gateway в Nginx
```bash
# Проверка, что сервисы запущены
docker compose ps

# Проверка портов
sudo netstat -tulpn | grep 809

# Логи Nginx
sudo tail -f /var/log/nginx/error.log
```

---

## 📞 Полезные команды

```bash
# Статус всех контейнеров
docker compose ps

# Логи последних 5 минут
docker compose logs --since 5m

# Вход в контейнер
docker compose exec auth-service sh

# Проверка сети
docker network ls
docker network inspect news-portal_news-network

# Очистка логов
sudo truncate -s 0 /var/log/nginx/news-portal-access.log
sudo truncate -s 0 /var/log/nginx/news-portal-error.log
```

---

## ✅ Checklist деплоя

- [ ] Сервер подготовлен (Docker, Nginx установлены)
- [ ] Проект загружен на сервер
- [ ] .env.production создан и заполнен
- [ ] Все пароли изменены
- [ ] Dockerfiles созданы для всех сервисов
- [ ] Docker Compose запущен
- [ ] Nginx настроен
- [ ] SSL сертификат получен
- [ ] Firewall настроен
- [ ] Автозапуск настроен
- [ ] Мониторинг работает
- [ ] API протестировано
- [ ] Резервное копирование настроено

---

## 🎉 Готово!

Ваш новостной портал теперь доступен по адресу:
- **HTTPS:** https://your-domain.com
- **API:** https://api.your-domain.com
- **Admin Panel:** https://your-domain.com/admin

**IP адрес:** http://151.241.228.203 (до настройки домена)

---

**Автор:** GitHub Copilot  
**Дата:** 14 октября 2025
