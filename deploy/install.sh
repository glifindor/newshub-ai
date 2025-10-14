#!/bin/bash

# ============================================================
# Автоматический скрипт установки News Portal
# Сервер: 151.241.228.203
# ============================================================

set -e  # Остановка при ошибке

echo "=================================================="
echo "  News Portal - Автоматическая установка"
echo "=================================================="
echo ""

# Проверка прав root
if [[ $EUID -ne 0 ]]; then
   echo "❌ Этот скрипт должен быть запущен с правами root (sudo)" 
   exit 1
fi

echo "✅ Проверка прав: OK"
echo ""

# Определение пользователя (не root)
if [ $SUDO_USER ]; then
    REAL_USER=$SUDO_USER
else
    REAL_USER=$(whoami)
fi

echo "👤 Пользователь: $REAL_USER"
echo ""

# ============================================================
# 1. Обновление системы
# ============================================================
echo "📦 Шаг 1: Обновление системы..."
apt update && apt upgrade -y
echo "✅ Система обновлена"
echo ""

# ============================================================
# 2. Установка зависимостей
# ============================================================
echo "📦 Шаг 2: Установка зависимостей..."
apt install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    git \
    nginx \
    ufw \
    wget

echo "✅ Зависимости установлены"
echo ""

# ============================================================
# 3. Установка Docker
# ============================================================
echo "🐳 Шаг 3: Установка Docker..."

# Удаление старых версий
apt remove -y docker docker-engine docker.io containerd runc 2>/dev/null || true

# Добавление GPG ключа Docker
mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# Добавление репозитория
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

# Установка Docker
apt update
apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Добавление пользователя в группу docker
usermod -aG docker $REAL_USER

# Проверка
docker --version
docker compose version

echo "✅ Docker установлен"
echo ""

# ============================================================
# 4. Настройка Firewall
# ============================================================
echo "🔒 Шаг 4: Настройка Firewall..."

ufw --force enable
ufw allow OpenSSH
ufw allow 80/tcp
ufw allow 443/tcp
ufw reload

echo "✅ Firewall настроен"
echo ""

# ============================================================
# 5. Создание директорий
# ============================================================
echo "📁 Шаг 5: Создание директорий..."

PROJECT_DIR="/opt/news-portal"
mkdir -p $PROJECT_DIR
chown -R $REAL_USER:$REAL_USER $PROJECT_DIR

echo "✅ Директории созданы: $PROJECT_DIR"
echo ""

# ============================================================
# 6. Генерация .env файла
# ============================================================
echo "⚙️  Шаг 6: Генерация конфигурации..."

POSTGRES_PASSWORD=$(openssl rand -base64 32)
REDIS_PASSWORD=$(openssl rand -base64 32)
MINIO_PASSWORD=$(openssl rand -base64 32)
RABBITMQ_PASSWORD=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 64)
GRAFANA_PASSWORD=$(openssl rand -base64 16)

cat > $PROJECT_DIR/.env.production << EOF
# Environment
ENVIRONMENT=production

# Server
SERVER_IP=151.241.228.203

# Database
POSTGRES_USER=newsportal
POSTGRES_PASSWORD=$POSTGRES_PASSWORD
POSTGRES_DB=news_portal

# Redis
REDIS_PASSWORD=$REDIS_PASSWORD

# MinIO
MINIO_ROOT_USER=newsportal_admin
MINIO_ROOT_PASSWORD=$MINIO_PASSWORD

# RabbitMQ
RABBITMQ_DEFAULT_USER=newsportal_admin
RABBITMQ_DEFAULT_PASS=$RABBITMQ_PASSWORD

# JWT Secret
JWT_SECRET=$JWT_SECRET

# Services Ports
AUTH_SERVICE_PORT=8091
AUTH_GRPC_PORT=8081
NEWS_SERVICE_PORT=8092
NEWS_GRPC_PORT=8082
MEDIA_SERVICE_PORT=8094
GATEWAY_PORT=8080
FRONTEND_PORT=3000

# Monitoring
GRAFANA_ADMIN_PASSWORD=$GRAFANA_PASSWORD
EOF

# Сохранение паролей в отдельный файл для безопасности
cat > $PROJECT_DIR/PASSWORDS.txt << EOF
============================================
  News Portal - Пароли и секреты
============================================

PostgreSQL:
  User: newsportal
  Password: $POSTGRES_PASSWORD

Redis:
  Password: $REDIS_PASSWORD

MinIO:
  User: newsportal_admin
  Password: $MINIO_PASSWORD

RabbitMQ:
  User: newsportal_admin
  Password: $RABBITMQ_PASSWORD

JWT Secret:
  $JWT_SECRET

Grafana:
  User: admin
  Password: $GRAFANA_PASSWORD

============================================
⚠️  ВАЖНО: Сохраните этот файл в безопасном месте!
============================================
EOF

chmod 600 $PROJECT_DIR/PASSWORDS.txt
chown $REAL_USER:$REAL_USER $PROJECT_DIR/PASSWORDS.txt

echo "✅ Конфигурация создана"
echo "📄 Пароли сохранены в: $PROJECT_DIR/PASSWORDS.txt"
echo ""

# ============================================================
# 7. Информация о следующих шагах
# ============================================================
echo "=================================================="
echo "  ✅ Базовая установка завершена!"
echo "=================================================="
echo ""
echo "📋 Следующие шаги:"
echo ""
echo "1️⃣  Загрузите проект на сервер:"
echo "   cd $PROJECT_DIR"
echo "   git clone YOUR_REPO_URL ."
echo ""
echo "   ИЛИ загрузите файлы вручную:"
echo "   scp -r project/* root@151.241.228.203:$PROJECT_DIR/"
echo ""
echo "2️⃣  Создайте Dockerfiles для сервисов (если не созданы)"
echo ""
echo "3️⃣  Запустите инфраструктуру:"
echo "   cd $PROJECT_DIR"
echo "   docker compose up -d postgres redis minio"
echo ""
echo "4️⃣  Соберите и запустите сервисы:"
echo "   docker compose build"
echo "   docker compose up -d"
echo ""
echo "5️⃣  Настройте Nginx и SSL (см. DEPLOYMENT_GUIDE.md)"
echo ""
echo "=================================================="
echo "📄 Пароли сохранены в: $PROJECT_DIR/PASSWORDS.txt"
echo "=================================================="
echo ""
echo "⚠️  НЕ ЗАБУДЬТЕ:"
echo "   - Скопировать PASSWORDS.txt в безопасное место"
echo "   - Удалить PASSWORDS.txt с сервера после копирования"
echo "   - Настроить доменное имя и SSL"
echo ""
