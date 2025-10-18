# ============================================
# 🚀 AUTO DEPLOY SCRIPT для NewsHub AI (PowerShell)
# Автоматический деплой на сервер 151.241.228.203
# ============================================

param(
    [string]$Action = "deploy",
    [string]$Password = ""
)

# Конфигурация
$SERVER_HOST = "151.241.228.203"
$SERVER_USER = "root"
$DEPLOY_DIR = "/opt/newshub"

# ============================================
# Функции
# ============================================

function Write-Header {
    param([string]$Message)
    Write-Host "============================================" -ForegroundColor Blue
    Write-Host $Message -ForegroundColor Blue
    Write-Host "============================================" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "✅ $Message" -ForegroundColor Green
}

function Write-ErrorMsg {
    param([string]$Message)
    Write-Host "❌ $Message" -ForegroundColor Red
}

function Write-Warning {
    param([string]$Message)
    Write-Host "⚠️  $Message" -ForegroundColor Yellow
}

function Write-Info {
    param([string]$Message)
    Write-Host "ℹ️  $Message" -ForegroundColor Cyan
}

# Проверка наличия plink (PuTTY)
function Test-PlinkInstalled {
    try {
        $null = Get-Command plink -ErrorAction Stop
        return $true
    } catch {
        Write-ErrorMsg "plink не найден. Установите PuTTY или используйте OpenSSH"
        Write-Info "Скачать: https://www.putty.org/"
        return $false
    }
}

# Выполнение команды на сервере
function Invoke-SSHCommand {
    param(
        [string]$Command,
        [string]$Password
    )
    
    if ($Password) {
        # Используем plink с паролем
        $result = echo y | plink -ssh -pw $Password ${SERVER_USER}@${SERVER_HOST} $Command 2>&1
    } else {
        # Используем стандартный SSH (ключ)
        $result = ssh ${SERVER_USER}@${SERVER_HOST} $Command 2>&1
    }
    
    return $result
}

# Проверка SSH соединения
function Test-SSHConnection {
    param([string]$Password)
    
    Write-Info "Проверка SSH соединения с $SERVER_HOST..."
    
    try {
        $result = Invoke-SSHCommand -Command "echo connected" -Password $Password
        if ($result -match "connected") {
            Write-Success "SSH соединение установлено"
            return $true
        } else {
            Write-ErrorMsg "Не удалось подключиться к серверу"
            return $false
        }
    } catch {
        Write-ErrorMsg "Ошибка SSH: $_"
        return $false
    }
}

# Создание backup
function New-DatabaseBackup {
    param([string]$Password)
    
    Write-Info "Создание backup базы данных..."
    
    $commands = @"
cd $DEPLOY_DIR
mkdir -p backups
BACKUP_FILE="backups/backup_`$(date +%Y%m%d_%H%M%S).sql"
docker-compose exec -T postgres pg_dump -U `${POSTGRES_USER:-newsadmin} `${POSTGRES_DB:-newshub_db} > `$BACKUP_FILE
echo "Backup created: `$BACKUP_FILE"
find backups/ -name "backup_*.sql" -mtime +7 -delete
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
    Write-Success "Backup создан успешно"
}

# Обновление кода
function Update-Code {
    param([string]$Password)
    
    Write-Info "Обновление кода на сервере..."
    
    $commands = @"
cd $DEPLOY_DIR
git pull origin main
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
    Write-Success "Код обновлен"
}

# Pull Docker images
function Get-DockerImages {
    param([string]$Password)
    
    Write-Info "Загрузка Docker images..."
    
    $commands = @"
cd $DEPLOY_DIR
docker-compose -f docker-compose.prod.yml pull
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
    Write-Success "Docker images загружены"
}

# Запуск контейнеров
function Start-Containers {
    param([string]$Password)
    
    Write-Info "Запуск контейнеров..."
    
    $commands = @"
cd $DEPLOY_DIR
docker-compose -f docker-compose.prod.yml up -d --force-recreate --remove-orphans
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
    Write-Success "Контейнеры запущены"
}

# Ожидание готовности
function Wait-ForServices {
    Write-Info "Ожидание готовности сервисов (30 сек)..."
    Start-Sleep -Seconds 30
    Write-Success "Сервисы должны быть готовы"
}

# Запуск миграций
function Invoke-Migrations {
    param([string]$Password)
    
    Write-Info "Запуск миграций базы данных..."
    
    $commands = @"
cd $DEPLOY_DIR
docker-compose exec -T backend alembic upgrade head
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
    Write-Success "Миграции проверены"
}

# Health check
function Test-HealthCheck {
    Write-Info "Проверка здоровья сервисов..."
    
    try {
        # Проверка API
        $response = Invoke-WebRequest -Uri "http://${SERVER_HOST}/health" -UseBasicParsing -TimeoutSec 10
        if ($response.StatusCode -eq 200) {
            Write-Success "API работает (HTTP $($response.StatusCode))"
        } else {
            Write-ErrorMsg "API не отвечает (HTTP $($response.StatusCode))"
            return $false
        }
        
        # Проверка Frontend
        $response = Invoke-WebRequest -Uri "http://${SERVER_HOST}" -UseBasicParsing -TimeoutSec 10
        if ($response.StatusCode -eq 200) {
            Write-Success "Frontend работает (HTTP $($response.StatusCode))"
        } else {
            Write-ErrorMsg "Frontend не отвечает (HTTP $($response.StatusCode))"
            return $false
        }
        
        Write-Success "Все сервисы работают корректно"
        return $true
    } catch {
        Write-ErrorMsg "Ошибка проверки здоровья: $_"
        return $false
    }
}

# Показать статус
function Show-Status {
    param([string]$Password)
    
    Write-Info "Статус контейнеров:"
    
    $commands = @"
cd $DEPLOY_DIR
docker-compose ps
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
}

# Показать логи
function Show-Logs {
    param([string]$Password)
    
    Write-Info "Последние логи:"
    
    $commands = @"
cd $DEPLOY_DIR
docker-compose logs --tail=50
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
}

# Очистка Docker
function Clear-Docker {
    param([string]$Password)
    
    Write-Info "Очистка старых Docker images..."
    
    $commands = @"
docker image prune -f
docker image prune -af --filter "until=72h"
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
    Write-Success "Docker очищен"
}

# ============================================
# Основная логика
# ============================================

function Start-Deployment {
    param([string]$Password)
    
    Write-Header "🚀 НАЧАЛО ДЕПЛОЯ NEWSHUB AI"
    
    # 1. Проверка SSH
    if (-not (Test-SSHConnection -Password $Password)) {
        Write-ErrorMsg "Невозможно подключиться к серверу"
        exit 1
    }
    
    # 2. Backup
    Write-Header "📦 СОЗДАНИЕ BACKUP"
    New-DatabaseBackup -Password $Password
    
    # 3. Обновление кода
    Write-Header "🔄 ОБНОВЛЕНИЕ КОДА"
    Update-Code -Password $Password
    
    # 4. Pull images
    Write-Header "🐳 ЗАГРУЗКА DOCKER IMAGES"
    Get-DockerImages -Password $Password
    
    # 5. Запуск
    Write-Header "▶️  ЗАПУСК КОНТЕЙНЕРОВ"
    Start-Containers -Password $Password
    
    # 6. Ожидание
    Wait-ForServices
    
    # 7. Миграции
    Write-Header "🗄️  МИГРАЦИИ БД"
    Invoke-Migrations -Password $Password
    
    # 8. Health check
    Write-Header "🏥 ПРОВЕРКА ЗДОРОВЬЯ"
    if (Test-HealthCheck) {
        Write-Success "Деплой завершен успешно!"
    } else {
        Write-ErrorMsg "Деплой завершен с ошибками"
        Show-Logs -Password $Password
        exit 1
    }
    
    # 9. Статус
    Write-Header "📊 СТАТУС СЕРВИСОВ"
    Show-Status -Password $Password
    
    # 10. Очистка
    Write-Header "🧹 ОЧИСТКА"
    Clear-Docker -Password $Password
    
    Write-Header "✅ ДЕПЛОЙ ЗАВЕРШЕН УСПЕШНО"
    Write-Info "Frontend: http://$SERVER_HOST"
    Write-Info "API Docs: http://$SERVER_HOST/docs"
    Write-Info "Grafana: http://$SERVER_HOST:3001"
}

# ============================================
# Точка входа
# ============================================

# Если пароль не указан, запросить
if (-not $Password) {
    Write-Info "Для автоматического деплоя нужен пароль от сервера"
    Write-Info "Использование: .\deploy.ps1 -Password 'ваш_пароль'"
    Write-Info "Или: .\deploy.ps1 -Action status -Password 'ваш_пароль'"
    Write-Warning "Введите пароль (или нажмите Enter для использования SSH ключа):"
    $SecurePassword = Read-Host -AsSecureString
    $BSTR = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($SecurePassword)
    $Password = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)
}

# Выполнение действия
switch ($Action) {
    "deploy" {
        Start-Deployment -Password $Password
    }
    "backup" {
        Write-Header "📦 ТОЛЬКО BACKUP"
        if (Test-SSHConnection -Password $Password) {
            New-DatabaseBackup -Password $Password
        }
    }
    "update" {
        Write-Header "🔄 ТОЛЬКО ОБНОВЛЕНИЕ КОДА"
        if (Test-SSHConnection -Password $Password) {
            Update-Code -Password $Password
        }
    }
    "restart" {
        Write-Header "🔄 ПЕРЕЗАПУСК КОНТЕЙНЕРОВ"
        if (Test-SSHConnection -Password $Password) {
            Start-Containers -Password $Password
            Wait-ForServices
            Test-HealthCheck
        }
    }
    "status" {
        Write-Header "📊 СТАТУС СЕРВИСОВ"
        if (Test-SSHConnection -Password $Password) {
            Show-Status -Password $Password
        }
    }
    "logs" {
        Write-Header "📜 ЛОГИ"
        if (Test-SSHConnection -Password $Password) {
            Show-Logs -Password $Password
        }
    }
    "cleanup" {
        Write-Header "🧹 ОЧИСТКА"
        if (Test-SSHConnection -Password $Password) {
            Clear-Docker -Password $Password
        }
    }
    "health" {
        Write-Header "🏥 HEALTH CHECK"
        Test-HealthCheck
    }
    default {
        Write-ErrorMsg "Неизвестное действие: $Action"
        Write-Info "Доступные действия: deploy, backup, update, restart, status, logs, cleanup, health"
        exit 1
    }
}
