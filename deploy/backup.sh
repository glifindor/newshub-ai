#!/bin/bash

# ============================================================
# Скрипт создания резервной копии
# ============================================================

set -e

PROJECT_DIR="/opt/news-portal"
BACKUP_DIR="/opt/news-portal-backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_PATH="$BACKUP_DIR/backup_$DATE"

echo "💾 Создание резервной копии News Portal..."
echo ""

# Создание директории для бэкапов
mkdir -p $BACKUP_DIR

# Создание директории для текущего бэкапа
mkdir -p $BACKUP_PATH

echo "📦 Бэкап PostgreSQL..."
docker compose exec -T postgres pg_dumpall -U postgres > $BACKUP_PATH/postgres_dump.sql

echo "📦 Бэкап Redis..."
docker compose exec -T redis redis-cli --rdb /data/dump.rdb SAVE
docker cp news-redis:/data/dump.rdb $BACKUP_PATH/redis_dump.rdb

echo "📦 Бэкап MinIO данных..."
docker run --rm \
    --volumes-from news-minio \
    -v $BACKUP_PATH:/backup \
    alpine tar czf /backup/minio_data.tar.gz /data

echo "📦 Бэкап конфигурации..."
cp $PROJECT_DIR/.env.production $BACKUP_PATH/
cp -r $PROJECT_DIR/docker-compose.yml $BACKUP_PATH/ 2>/dev/null || true

# Создание архива
echo "📦 Создание архива..."
cd $BACKUP_DIR
tar -czf backup_$DATE.tar.gz backup_$DATE
rm -rf backup_$DATE

# Удаление старых бэкапов (старше 7 дней)
find $BACKUP_DIR -name "backup_*.tar.gz" -mtime +7 -delete

echo ""
echo "✅ Резервная копия создана: $BACKUP_DIR/backup_$DATE.tar.gz"
echo ""
echo "📋 Восстановление:"
echo "   ./restore.sh backup_$DATE.tar.gz"
echo ""
