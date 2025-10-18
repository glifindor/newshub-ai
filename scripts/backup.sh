#!/bin/bash

# ============================================
# 📦 BACKUP SCRIPT для NewsHub AI
# Автоматический backup всех важных данных
# ============================================

set -e

# Конфигурация
BACKUP_DIR="/opt/newshub/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=7  # Хранить backup'ы 7 дней

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# ============================================
# Функции
# ============================================

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Создание директорий
create_backup_dirs() {
    log_info "Создание директорий для backup..."
    mkdir -p "$BACKUP_DIR/database"
    mkdir -p "$BACKUP_DIR/redis"
    mkdir -p "$BACKUP_DIR/logs"
    mkdir -p "$BACKUP_DIR/configs"
}

# Backup PostgreSQL
backup_postgres() {
    log_info "Backup PostgreSQL..."
    
    local backup_file="$BACKUP_DIR/database/postgres_${TIMESTAMP}.sql"
    local backup_file_gz="$backup_file.gz"
    
    cd /opt/newshub
    
    # Создать SQL dump
    docker-compose exec -T postgres pg_dump -U ${POSTGRES_USER:-newsadmin} ${POSTGRES_DB:-newshub_db} > "$backup_file"
    
    if [ $? -eq 0 ]; then
        # Сжать backup
        gzip "$backup_file"
        
        local size=$(du -h "$backup_file_gz" | cut -f1)
        log_info "PostgreSQL backup создан: $backup_file_gz ($size)"
        
        # Создать latest symlink
        ln -sf "$(basename $backup_file_gz)" "$BACKUP_DIR/database/postgres_latest.sql.gz"
    else
        log_error "Ошибка создания PostgreSQL backup"
        return 1
    fi
}

# Backup Redis
backup_redis() {
    log_info "Backup Redis..."
    
    local backup_file="$BACKUP_DIR/redis/redis_${TIMESTAMP}.rdb"
    
    cd /opt/newshub
    
    # Сохранить Redis RDB
    docker-compose exec -T redis redis-cli SAVE
    
    # Скопировать RDB файл
    docker cp newshub_redis:/data/dump.rdb "$backup_file"
    
    if [ $? -eq 0 ]; then
        local size=$(du -h "$backup_file" | cut -f1)
        log_info "Redis backup создан: $backup_file ($size)"
        
        # Создать latest symlink
        ln -sf "$(basename $backup_file)" "$BACKUP_DIR/redis/redis_latest.rdb"
    else
        log_error "Ошибка создания Redis backup"
        return 1
    fi
}

# Backup логов
backup_logs() {
    log_info "Backup логов..."
    
    local backup_file="$BACKUP_DIR/logs/logs_${TIMESTAMP}.tar.gz"
    
    cd /opt/newshub
    
    # Архивировать логи
    tar -czf "$backup_file" \
        backend/logs/ \
        2>/dev/null || true
    
    if [ -f "$backup_file" ]; then
        local size=$(du -h "$backup_file" | cut -f1)
        log_info "Логи заархивированы: $backup_file ($size)"
    else
        log_warning "Логи не найдены или пусты"
    fi
}

# Backup конфигураций
backup_configs() {
    log_info "Backup конфигураций..."
    
    local backup_file="$BACKUP_DIR/configs/configs_${TIMESTAMP}.tar.gz"
    
    cd /opt/newshub
    
    # Архивировать конфиги (без .env!)
    tar -czf "$backup_file" \
        docker-compose.yml \
        docker-compose.prod.yml \
        nginx/ \
        monitoring/ \
        2>/dev/null || true
    
    if [ -f "$backup_file" ]; then
        local size=$(du -h "$backup_file" | cut -f1)
        log_info "Конфигурации заархивированы: $backup_file ($size)"
    else
        log_warning "Конфигурации не найдены"
    fi
}

# Очистка старых backup'ов
cleanup_old_backups() {
    log_info "Очистка старых backup'ов (старше $RETENTION_DAYS дней)..."
    
    local deleted_count=0
    
    # Удалить старые файлы
    find "$BACKUP_DIR" -type f \( -name "*.sql.gz" -o -name "*.rdb" -o -name "*.tar.gz" \) -mtime +$RETENTION_DAYS -delete
    
    log_info "Старые backup'ы удалены"
}

# Показать статистику
show_statistics() {
    log_info "Статистика backup'ов:"
    echo ""
    
    echo "📁 Директория: $BACKUP_DIR"
    echo "📊 Общий размер: $(du -sh $BACKUP_DIR | cut -f1)"
    echo ""
    
    echo "🗄️  PostgreSQL backup'ы:"
    ls -lh "$BACKUP_DIR/database/" | grep "\.sql\.gz" | tail -5
    echo ""
    
    echo "💾 Redis backup'ы:"
    ls -lh "$BACKUP_DIR/redis/" | grep "\.rdb" | tail -5
    echo ""
    
    echo "📝 Логи:"
    ls -lh "$BACKUP_DIR/logs/" | grep "\.tar\.gz" | tail -5
    echo ""
}

# Отправка уведомления в Telegram
send_telegram_notification() {
    local message=$1
    local status=$2  # success или error
    
    # Проверить, установлен ли TELEGRAM_BOT_TOKEN
    if [ -z "$TELEGRAM_BOT_TOKEN" ] || [ -z "$TELEGRAM_ADMIN_CHAT_ID" ]; then
        return
    fi
    
    local emoji
    if [ "$status" == "success" ]; then
        emoji="✅"
    else
        emoji="❌"
    fi
    
    local full_message="${emoji} BACKUP NOTIFICATION%0A%0A${message}%0A%0ATime: $(date '+%Y-%m-%d %H:%M:%S')"
    
    curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
        -d "chat_id=${TELEGRAM_ADMIN_CHAT_ID}" \
        -d "text=${full_message}" \
        > /dev/null 2>&1
}

# ============================================
# Главная функция
# ============================================

main() {
    echo ""
    echo "============================================"
    echo "📦 BACKUP NEWSHUB AI"
    echo "============================================"
    echo ""
    
    # Загрузить переменные окружения
    if [ -f "/opt/newshub/.env" ]; then
        source /opt/newshub/.env
    fi
    
    # Создать директории
    create_backup_dirs
    
    # Выполнить backup'ы
    local errors=0
    
    if ! backup_postgres; then
        ((errors++))
    fi
    
    if ! backup_redis; then
        ((errors++))
    fi
    
    backup_logs
    backup_configs
    
    # Очистить старые backup'ы
    cleanup_old_backups
    
    # Показать статистику
    show_statistics
    
    # Результат
    if [ $errors -eq 0 ]; then
        log_info "✅ Все backup'ы созданы успешно!"
        send_telegram_notification "All backups completed successfully" "success"
        exit 0
    else
        log_error "❌ Backup завершен с ошибками ($errors)"
        send_telegram_notification "Backup completed with $errors errors" "error"
        exit 1
    fi
}

# ============================================
# Запуск
# ============================================

main
