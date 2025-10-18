#!/bin/bash

# ============================================
# ⏰ SETUP CRON для автоматических backup'ов
# ============================================

echo "Настройка автоматических backup'ов для NewsHub AI"
echo ""

# Проверка прав root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Запустите с правами root: sudo ./setup-cron.sh"
    exit 1
fi

# Путь к скрипту backup
BACKUP_SCRIPT="/opt/newshub/scripts/backup.sh"

# Проверка существования скрипта
if [ ! -f "$BACKUP_SCRIPT" ]; then
    echo "❌ Скрипт backup не найден: $BACKUP_SCRIPT"
    exit 1
fi

# Сделать скрипт исполняемым
chmod +x "$BACKUP_SCRIPT"

# Создать cron задачу
CRON_JOB="0 3 * * * $BACKUP_SCRIPT >> /opt/newshub/logs/backup.log 2>&1"

# Проверить, существует ли уже задача
if crontab -l 2>/dev/null | grep -q "$BACKUP_SCRIPT"; then
    echo "⚠️  Cron задача уже существует"
    echo ""
    echo "Текущие cron задачи:"
    crontab -l | grep backup
    echo ""
    read -p "Обновить задачу? (y/n): " confirm
    if [ "$confirm" != "y" ]; then
        echo "Отменено"
        exit 0
    fi
    # Удалить старую задачу
    crontab -l | grep -v "$BACKUP_SCRIPT" | crontab -
fi

# Добавить новую задачу
(crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -

echo "✅ Cron задача добавлена!"
echo ""
echo "📅 Расписание: Каждый день в 3:00 AM"
echo "📝 Логи: /opt/newshub/logs/backup.log"
echo ""
echo "Текущие cron задачи:"
crontab -l
echo ""
echo "Для ручного запуска: $BACKUP_SCRIPT"
echo ""

# Создать директорию для логов
mkdir -p /opt/newshub/logs

# Тестовый запуск
read -p "Запустить тестовый backup сейчас? (y/n): " test_run
if [ "$test_run" == "y" ]; then
    echo ""
    echo "Запуск тестового backup..."
    $BACKUP_SCRIPT
fi

echo ""
echo "✅ Настройка завершена!"
