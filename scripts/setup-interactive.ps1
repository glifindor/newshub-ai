# ============================================
# 🎯 INTERACTIVE SETUP для NewsHub AI
# Интерактивная настройка сервера с нуля
# ============================================

param(
    [string]$Password = ""
)

$SERVER_HOST = "151.241.228.203"
$SERVER_USER = "root"
$DEPLOY_DIR = "/opt/newshub"
$REPO_URL = "https://github.com/glifindor/newsportal.git"

# ============================================
# Функции для вывода
# ============================================

function Write-Header {
    param([string]$Message)
    Write-Host ""
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

function Write-Step {
    param([int]$Step, [string]$Message)
    Write-Host ""
    Write-Host "[$Step/10] $Message" -ForegroundColor Magenta
    Write-Host "-------------------------------------------" -ForegroundColor DarkGray
}

# Выполнение команды на сервере
function Invoke-SSHCommand {
    param(
        [string]$Command,
        [string]$Password
    )
    
    if ($Password) {
        $result = echo y | plink -ssh -pw $Password ${SERVER_USER}@${SERVER_HOST} $Command 2>&1
    } else {
        $result = ssh ${SERVER_USER}@${SERVER_HOST} $Command 2>&1
    }
    
    return $result
}

# ============================================
# Шаги установки
# ============================================

function Step1-TestConnection {
    param([string]$Password)
    
    Write-Step 1 "Проверка подключения к серверу"
    
    try {
        $result = Invoke-SSHCommand -Command "echo connected" -Password $Password
        if ($result -match "connected") {
            Write-Success "Подключение установлено"
            return $true
        } else {
            Write-ErrorMsg "Не удалось подключиться"
            return $false
        }
    } catch {
        Write-ErrorMsg "Ошибка: $_"
        return $false
    }
}

function Step2-UpdateSystem {
    param([string]$Password)
    
    Write-Step 2 "Обновление системы"
    Write-Info "Это может занять несколько минут..."
    
    $commands = @"
export DEBIAN_FRONTEND=noninteractive
apt-get update -qq
apt-get upgrade -y -qq
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Success "Система обновлена"
}

function Step3-InstallDocker {
    param([string]$Password)
    
    Write-Step 3 "Установка Docker и Docker Compose"
    
    $commands = @"
# Установка Docker
if ! command -v docker &> /dev/null; then
    apt-get install -y docker.io
    systemctl start docker
    systemctl enable docker
    echo "Docker installed"
else
    echo "Docker already installed"
fi

# Установка Docker Compose
if ! command -v docker-compose &> /dev/null; then
    apt-get install -y docker-compose
    echo "Docker Compose installed"
else
    echo "Docker Compose already installed"
fi

# Проверка версий
docker --version
docker-compose --version
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
    Write-Success "Docker и Docker Compose установлены"
}

function Step4-InstallTools {
    param([string]$Password)
    
    Write-Step 4 "Установка дополнительных инструментов"
    
    $commands = @"
apt-get install -y git nano curl wget htop ufw jq
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Success "Инструменты установлены"
}

function Step5-SetupFirewall {
    param([string]$Password)
    
    Write-Step 5 "Настройка Firewall (UFW)"
    
    $commands = @"
# Разрешить SSH (ВАЖНО!)
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp

# Включить firewall
echo "y" | ufw enable

# Статус
ufw status
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
    Write-Success "Firewall настроен"
}

function Step6-CloneRepository {
    param([string]$Password)
    
    Write-Step 6 "Клонирование репозитория"
    
    $commands = @"
# Создать директорию
mkdir -p $DEPLOY_DIR
cd $DEPLOY_DIR

# Клонировать репозиторий
if [ -d ".git" ]; then
    echo "Repository already exists, pulling latest..."
    git pull origin main
else
    echo "Cloning repository..."
    git clone $REPO_URL .
fi

# Проверка
ls -la
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
    Write-Success "Репозиторий склонирован"
}

function Step7-CreateEnvFile {
    param([string]$Password)
    
    Write-Step 7 "Создание .env файла"
    
    Write-Info "Сейчас вам нужно ввести конфиденциальные данные"
    Write-Warning "Эти данные будут сохранены ТОЛЬКО на сервере"
    Write-Host ""
    
    # Сбор данных
    $POSTGRES_PASSWORD = Read-Host "Пароль PostgreSQL (придумайте сложный)"
    $REDIS_PASSWORD = Read-Host "Пароль Redis"
    $JWT_SECRET = -join ((48..57) + (65..90) + (97..122) | Get-Random -Count 32 | ForEach-Object {[char]$_})
    
    Write-Host ""
    Write-Info "JWT Secret сгенерирован автоматически: $JWT_SECRET"
    Write-Host ""
    
    $OPENROUTER_API_KEY = Read-Host "OpenRouter API Key"
    $FREEPIK_API_KEY = Read-Host "Freepik API Key (или Enter для пропуска)"
    $NEWSAPI_KEY = Read-Host "NewsAPI Key (или Enter для пропуска)"
    $TELEGRAM_BOT_TOKEN = Read-Host "Telegram Bot Token"
    $TELEGRAM_ADMIN_CHAT_ID = Read-Host "Telegram Admin Chat ID"
    
    # Создание .env на сервере
    $envContent = @"
# Database
POSTGRES_USER=newsadmin
POSTGRES_PASSWORD=$POSTGRES_PASSWORD
POSTGRES_DB=newshub_db

# Redis
REDIS_PASSWORD=$REDIS_PASSWORD

# RabbitMQ
RABBITMQ_USER=newshub
RABBITMQ_PASS=$REDIS_PASSWORD

# JWT
JWT_SECRET_KEY=$JWT_SECRET

# APIs
OPENROUTER_API_KEY=$OPENROUTER_API_KEY
FREEPIK_API_KEY=$FREEPIK_API_KEY
NEWSAPI_KEY=$NEWSAPI_KEY

# Telegram
TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN
TELEGRAM_CRYPTO_CHANNEL=@crypto_ainews
TELEGRAM_POLITICS_CHANNEL=@kremlin_digest
TELEGRAM_ADMIN_CHAT_ID=$TELEGRAM_ADMIN_CHAT_ID

# Monitoring
GRAFANA_USER=admin
GRAFANA_PASSWORD=admin123
FLOWER_USER=admin
FLOWER_PASSWORD=admin123

# Frontend
NEXT_PUBLIC_API_URL=http://$SERVER_HOST/api
NEXT_PUBLIC_WS_URL=ws://$SERVER_HOST

# Environment
ENVIRONMENT=production
LOG_LEVEL=INFO
"@
    
    # Сохранение на сервере
    $envContent | Out-File -FilePath "temp_env.txt" -Encoding UTF8
    
    # Копирование на сервер
    if ($Password) {
        pscp -pw $Password temp_env.txt ${SERVER_USER}@${SERVER_HOST}:${DEPLOY_DIR}/.env
    } else {
        scp temp_env.txt ${SERVER_USER}@${SERVER_HOST}:${DEPLOY_DIR}/.env
    }
    
    Remove-Item temp_env.txt
    
    Write-Success ".env файл создан на сервере"
}

function Step8-BuildImages {
    param([string]$Password)
    
    Write-Step 8 "Сборка Docker образов"
    Write-Info "Это займет 5-15 минут. Подождите..."
    
    $commands = @"
cd $DEPLOY_DIR
docker-compose -f docker-compose.prod.yml build
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
    Write-Success "Образы собраны"
}

function Step9-StartContainers {
    param([string]$Password)
    
    Write-Step 9 "Запуск контейнеров"
    
    $commands = @"
cd $DEPLOY_DIR

# Создать директории
mkdir -p backups logs

# Запустить контейнеры
docker-compose -f docker-compose.prod.yml up -d

# Подождать 30 секунд
sleep 30

# Запустить миграции
docker-compose exec -T backend alembic upgrade head

# Статус
docker-compose ps
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
    Write-Success "Контейнеры запущены"
}

function Step10-CreateAdmin {
    param([string]$Password)
    
    Write-Step 10 "Создание администратора"
    
    Write-Info "Введите данные администратора:"
    $ADMIN_USERNAME = Read-Host "Username"
    $ADMIN_EMAIL = Read-Host "Email"
    $ADMIN_PASSWORD = Read-Host "Password" -AsSecureString
    $BSTR = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($ADMIN_PASSWORD)
    $ADMIN_PASSWORD_TEXT = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)
    
    $commands = @"
cd $DEPLOY_DIR
docker-compose exec -T backend python -c "
from app.database import get_db
from app.models import User
from app.auth import get_password_hash
from sqlalchemy.orm import Session

db = next(get_db())
user = User(
    username='$ADMIN_USERNAME',
    email='$ADMIN_EMAIL',
    hashed_password=get_password_hash('$ADMIN_PASSWORD_TEXT'),
    is_active=True,
    is_superuser=True
)
db.add(user)
db.commit()
print('Admin created!')
"
"@
    
    $result = Invoke-SSHCommand -Command $commands -Password $Password
    Write-Host $result
    Write-Success "Администратор создан"
}

# ============================================
# Главная функция
# ============================================

function Start-InteractiveSetup {
    param([string]$Password)
    
    Write-Header "🎯 ИНТЕРАКТИВНАЯ УСТАНОВКА NEWSHUB AI"
    Write-Info "Сервер: $SERVER_HOST"
    Write-Info "Пользователь: $SERVER_USER"
    Write-Host ""
    
    Write-Warning "Этот скрипт настроит сервер с нуля!"
    Write-Warning "Убедитесь, что у вас есть все необходимые API ключи"
    Write-Host ""
    
    $confirm = Read-Host "Продолжить? (yes/no)"
    if ($confirm -ne "yes") {
        Write-Info "Установка отменена"
        exit
    }
    
    # Выполнение шагов
    if (-not (Step1-TestConnection -Password $Password)) { exit 1 }
    Step2-UpdateSystem -Password $Password
    Step3-InstallDocker -Password $Password
    Step4-InstallTools -Password $Password
    Step5-SetupFirewall -Password $Password
    Step6-CloneRepository -Password $Password
    Step7-CreateEnvFile -Password $Password
    Step8-BuildImages -Password $Password
    Step9-StartContainers -Password $Password
    Step10-CreateAdmin -Password $Password
    
    # Финал
    Write-Header "🎉 УСТАНОВКА ЗАВЕРШЕНА!"
    Write-Host ""
    Write-Success "NewsHub AI успешно установлен!"
    Write-Host ""
    Write-Info "🌐 Доступные сервисы:"
    Write-Host "  Frontend:    http://$SERVER_HOST" -ForegroundColor Cyan
    Write-Host "  API Docs:    http://$SERVER_HOST/docs" -ForegroundColor Cyan
    Write-Host "  Grafana:     http://$SERVER_HOST:3001 (admin/admin123)" -ForegroundColor Cyan
    Write-Host "  RabbitMQ:    http://$SERVER_HOST:15672 (guest/guest)" -ForegroundColor Cyan
    Write-Host ""
    Write-Info "📝 Следующие шаги:"
    Write-Host "  1. Откройте http://$SERVER_HOST и проверьте работу" -ForegroundColor Yellow
    Write-Host "  2. Войдите в админ-панель с созданными учетными данными" -ForegroundColor Yellow
    Write-Host "  3. Настройте SSL сертификат (см. DEPLOYMENT.md)" -ForegroundColor Yellow
    Write-Host "  4. Настройте регулярные backup'ы" -ForegroundColor Yellow
    Write-Host ""
    Write-Info "📚 Документация: $DEPLOY_DIR/DEPLOYMENT.md"
    Write-Host ""
}

# ============================================
# Точка входа
# ============================================

if (-not $Password) {
    Write-Info "Введите пароль от сервера $SERVER_HOST"
    $SecurePassword = Read-Host -AsSecureString
    $BSTR = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($SecurePassword)
    $Password = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)
}

Start-InteractiveSetup -Password $Password
