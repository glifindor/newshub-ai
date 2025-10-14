#!/bin/bash

# ============================================================
# Скрипт быстрого запуска News Portal
# ============================================================

set -e

PROJECT_DIR="/opt/news-portal"

echo "🚀 Запуск News Portal..."
echo ""

cd $PROJECT_DIR

# Проверка наличия .env.production
if [ ! -f .env.production ]; then
    echo "❌ Файл .env.production не найден!"
    echo "Запустите сначала install.sh"
    exit 1
fi

# Остановка существующих контейнеров
echo "🛑 Остановка существующих контейнеров..."
docker compose down 2>/dev/null || true

# Запуск инфраструктуры
echo "🏗️  Запуск инфраструктуры (PostgreSQL, Redis, MinIO)..."
docker compose --env-file .env.production up -d postgres redis minio

echo "⏳ Ожидание готовности баз данных (30 секунд)..."
sleep 30

# Сборка сервисов
echo "🔨 Сборка микросервисов..."
docker compose --env-file .env.production build auth-service news-service media-service

# Запуск сервисов
echo "🚀 Запуск микросервисов..."
docker compose --env-file .env.production up -d auth-service news-service media-service

# Ожидание запуска
echo "⏳ Ожидание запуска сервисов (10 секунд)..."
sleep 10

# Проверка статуса
echo ""
echo "📊 Статус контейнеров:"
docker compose ps

echo ""
echo "🧪 Проверка health endpoints..."
sleep 5

curl -s http://localhost:8091/health && echo " ✅ Auth Service: OK" || echo " ❌ Auth Service: FAILED"
curl -s http://localhost:8092/health && echo " ✅ News Service: OK" || echo " ❌ News Service: FAILED"
curl -s http://localhost:8094/health && echo " ✅ Media Service: OK" || echo " ❌ Media Service: FAILED"

echo ""
echo "=================================================="
echo "  ✅ News Portal запущен!"
echo "=================================================="
echo ""
echo "📍 Сервисы доступны по адресам:"
echo "   Auth Service:  http://151.241.228.203:8091"
echo "   News Service:  http://151.241.228.203:8092"
echo "   Media Service: http://151.241.228.203:8094"
echo ""
echo "🗄️  Инфраструктура:"
echo "   PostgreSQL:    localhost:5432"
echo "   Redis:         localhost:6379"
echo "   MinIO:         http://151.241.228.203:9001"
echo ""
echo "📋 Полезные команды:"
echo "   docker compose logs -f              # Все логи"
echo "   docker compose logs -f auth-service # Логи Auth Service"
echo "   docker compose ps                   # Статус контейнеров"
echo "   docker compose restart              # Перезапуск"
echo "   docker compose down                 # Остановка"
echo ""
