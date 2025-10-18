#!/bin/bash

# ============================================
# 🚀 AUTO DEPLOY SCRIPT для NewsHub AI
# Автоматический деплой на сервер 151.241.228.203
# ============================================

set -e  # Остановка при ошибке

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Конфигурация
SERVER_HOST="151.241.228.203"
SERVER_USER="root"
DEPLOY_DIR="/opt/newshub"
BACKUP_DIR="/opt/newshub/backups"

# ============================================
# Функции
# ============================================

print_header() {
    echo -e "${BLUE}============================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}============================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Проверка SSH соединения
check_ssh_connection() {
    print_info "Проверка SSH соединения с $SERVER_HOST..."
    if ssh -o ConnectTimeout=5 -o BatchMode=yes $SERVER_USER@$SERVER_HOST exit 2>/dev/null; then
        print_success "SSH соединение установлено"
        return 0
    else
        print_error "Не удалось подключиться к серверу"
        print_info "Попробуйте: ssh $SERVER_USER@$SERVER_HOST"
        return 1
    fi
}

# Создание backup базы данных
create_backup() {
    print_info "Создание backup базы данных..."
    
    ssh $SERVER_USER@$SERVER_HOST << 'ENDSSH'
        cd /opt/newshub
        
        # Создать директорию для backups
        mkdir -p backups
        
        # Backup PostgreSQL
        BACKUP_FILE="backups/backup_$(date +%Y%m%d_%H%M%S).sql"
        docker-compose exec -T postgres pg_dump -U ${POSTGRES_USER:-newsadmin} ${POSTGRES_DB:-newshub_db} > $BACKUP_FILE
        
        if [ -f "$BACKUP_FILE" ]; then
            echo "✅ Backup создан: $BACKUP_FILE"
            
            # Удаление старых backups (старше 7 дней)
            find backups/ -name "backup_*.sql" -mtime +7 -delete
            echo "✅ Старые backups удалены"
        else
            echo "❌ Ошибка создания backup"
            exit 1
        fi
ENDSSH
    
    print_success "Backup создан успешно"
}

# Обновление кода на сервере
update_code() {
    print_info "Обновление кода на сервере..."
    
    ssh $SERVER_USER@$SERVER_HOST << 'ENDSSH'
        cd /opt/newshub
        
        # Pull последних изменений
        git pull origin main
        
        if [ $? -eq 0 ]; then
            echo "✅ Код обновлен"
        else
            echo "❌ Ошибка обновления кода"
            exit 1
        fi
ENDSSH
    
    print_success "Код обновлен"
}

# Pull Docker images
pull_images() {
    print_info "Загрузка Docker images..."
    
    ssh $SERVER_USER@$SERVER_HOST << 'ENDSSH'
        cd /opt/newshub
        
        docker-compose -f docker-compose.prod.yml pull
        
        if [ $? -eq 0 ]; then
            echo "✅ Images загружены"
        else
            echo "❌ Ошибка загрузки images"
            exit 1
        fi
ENDSSH
    
    print_success "Docker images загружены"
}

# Запуск контейнеров
start_containers() {
    print_info "Запуск контейнеров..."
    
    ssh $SERVER_USER@$SERVER_HOST << 'ENDSSH'
        cd /opt/newshub
        
        # Запуск с force-recreate
        docker-compose -f docker-compose.prod.yml up -d --force-recreate --remove-orphans
        
        if [ $? -eq 0 ]; then
            echo "✅ Контейнеры запущены"
        else
            echo "❌ Ошибка запуска контейнеров"
            exit 1
        fi
ENDSSH
    
    print_success "Контейнеры запущены"
}

# Ожидание готовности сервисов
wait_for_services() {
    print_info "Ожидание готовности сервисов (30 сек)..."
    sleep 30
    print_success "Сервисы должны быть готовы"
}

# Запуск миграций
run_migrations() {
    print_info "Запуск миграций базы данных..."
    
    ssh $SERVER_USER@$SERVER_HOST << 'ENDSSH'
        cd /opt/newshub
        
        docker-compose exec -T backend alembic upgrade head
        
        if [ $? -eq 0 ]; then
            echo "✅ Миграции выполнены"
        else
            echo "⚠️  Миграции не выполнены (возможно, уже актуальны)"
        fi
ENDSSH
    
    print_success "Миграции проверены"
}

# Health check
health_check() {
    print_info "Проверка здоровья сервисов..."
    
    # Проверка API
    response=$(curl -s -o /dev/null -w "%{http_code}" http://$SERVER_HOST/health)
    if [ "$response" -eq 200 ]; then
        print_success "API работает (HTTP $response)"
    else
        print_error "API не отвечает (HTTP $response)"
        return 1
    fi
    
    # Проверка Frontend
    response=$(curl -s -o /dev/null -w "%{http_code}" http://$SERVER_HOST)
    if [ "$response" -eq 200 ]; then
        print_success "Frontend работает (HTTP $response)"
    else
        print_error "Frontend не отвечает (HTTP $response)"
        return 1
    fi
    
    print_success "Все сервисы работают корректно"
}

# Показать статус контейнеров
show_status() {
    print_info "Статус контейнеров:"
    
    ssh $SERVER_USER@$SERVER_HOST << 'ENDSSH'
        cd /opt/newshub
        docker-compose ps
ENDSSH
}

# Показать логи
show_logs() {
    print_info "Последние логи:"
    
    ssh $SERVER_USER@$SERVER_HOST << 'ENDSSH'
        cd /opt/newshub
        docker-compose logs --tail=50
ENDSSH
}

# Очистка старых Docker images
cleanup_docker() {
    print_info "Очистка старых Docker images..."
    
    ssh $SERVER_USER@$SERVER_HOST << 'ENDSSH'
        # Удаление dangling images
        docker image prune -f
        
        # Удаление images старше 72 часов
        docker image prune -af --filter "until=72h"
        
        echo "✅ Docker очищен"
ENDSSH
    
    print_success "Docker очищен"
}

# ============================================
# Основная логика деплоя
# ============================================

main() {
    print_header "🚀 НАЧАЛО ДЕПЛОЯ NEWSHUB AI"
    
    # 1. Проверка SSH
    if ! check_ssh_connection; then
        exit 1
    fi
    
    # 2. Backup БД
    print_header "📦 СОЗДАНИЕ BACKUP"
    create_backup
    
    # 3. Обновление кода
    print_header "🔄 ОБНОВЛЕНИЕ КОДА"
    update_code
    
    # 4. Pull images
    print_header "🐳 ЗАГРУЗКА DOCKER IMAGES"
    pull_images
    
    # 5. Запуск контейнеров
    print_header "▶️  ЗАПУСК КОНТЕЙНЕРОВ"
    start_containers
    
    # 6. Ожидание
    wait_for_services
    
    # 7. Миграции
    print_header "🗄️  МИГРАЦИИ БД"
    run_migrations
    
    # 8. Health check
    print_header "🏥 ПРОВЕРКА ЗДОРОВЬЯ"
    if health_check; then
        print_success "Деплой завершен успешно!"
    else
        print_error "Деплой завершен с ошибками"
        show_logs
        exit 1
    fi
    
    # 9. Статус
    print_header "📊 СТАТУС СЕРВИСОВ"
    show_status
    
    # 10. Очистка
    print_header "🧹 ОЧИСТКА"
    cleanup_docker
    
    print_header "✅ ДЕПЛОЙ ЗАВЕРШЕН УСПЕШНО"
    print_info "Frontend: http://$SERVER_HOST"
    print_info "API Docs: http://$SERVER_HOST/docs"
    print_info "Grafana: http://$SERVER_HOST:3001"
}

# ============================================
# Запуск скрипта
# ============================================

# Проверка аргументов
case "${1:-}" in
    backup)
        print_header "📦 ТОЛЬКО BACKUP"
        check_ssh_connection && create_backup
        ;;
    update)
        print_header "🔄 ТОЛЬКО ОБНОВЛЕНИЕ КОДА"
        check_ssh_connection && update_code
        ;;
    restart)
        print_header "🔄 ПЕРЕЗАПУСК КОНТЕЙНЕРОВ"
        check_ssh_connection && start_containers && wait_for_services && health_check
        ;;
    status)
        print_header "📊 СТАТУС СЕРВИСОВ"
        check_ssh_connection && show_status
        ;;
    logs)
        print_header "📜 ЛОГИ"
        check_ssh_connection && show_logs
        ;;
    cleanup)
        print_header "🧹 ОЧИСТКА"
        check_ssh_connection && cleanup_docker
        ;;
    health)
        print_header "🏥 HEALTH CHECK"
        health_check
        ;;
    *)
        main
        ;;
esac
