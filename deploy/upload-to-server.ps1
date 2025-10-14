# ============================================================
# Скрипт загрузки проекта на сервер (Windows PowerShell)
# Сервер: 151.241.228.203
# ============================================================

$SERVER_IP = "151.241.228.203"
$SERVER_USER = "root"
$PROJECT_DIR = "C:\Users\Grifindor\Desktop\НОВСТНОЙ САЙТ"

Write-Host "=================================================="
Write-Host "  News Portal - Загрузка на сервер"
Write-Host "=================================================="
Write-Host ""

# Проверка наличия директории
if (-not (Test-Path $PROJECT_DIR)) {
    Write-Host "❌ Директория проекта не найдена: $PROJECT_DIR" -ForegroundColor Red
    exit 1
}

Set-Location $PROJECT_DIR

# ============================================================
# 1. Упаковка проекта
# ============================================================
Write-Host "📦 Шаг 1: Упаковка проекта..."

$archiveName = "news-portal-$(Get-Date -Format 'yyyyMMdd-HHmmss').tar.gz"

# Создаем список файлов для упаковки
$filesToPack = @(
    "auth-service",
    "news-service", 
    "media-service",
    "docker-compose.yml",
    "deploy",
    "DEPLOYMENT_GUIDE.md",
    "QUICK_START.md",
    "PROJECT_SUMMARY.md"
)

# Используем tar (встроенный в Windows 10+)
tar -czf $archiveName $filesToPack

if (Test-Path $archiveName) {
    $fileSize = (Get-Item $archiveName).Length / 1MB
    Write-Host "✅ Архив создан: $archiveName ($([math]::Round($fileSize, 2)) MB)" -ForegroundColor Green
} else {
    Write-Host "❌ Не удалось создать архив" -ForegroundColor Red
    exit 1
}

Write-Host ""

# ============================================================
# 2. Загрузка на сервер
# ============================================================
Write-Host "📤 Шаг 2: Загрузка на сервер $SERVER_IP..."
Write-Host ""

# Проверка наличия SSH
try {
    ssh -V 2>&1 | Out-Null
} catch {
    Write-Host "❌ SSH не установлен. Установите OpenSSH Client в Windows Settings" -ForegroundColor Red
    exit 1
}

# Загрузка архива
Write-Host "Загрузка архива проекта..."
scp $archiveName "${SERVER_USER}@${SERVER_IP}:/root/"

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Архив загружен на сервер" -ForegroundColor Green
} else {
    Write-Host "❌ Ошибка при загрузке архива" -ForegroundColor Red
    exit 1
}

# Загрузка скрипта установки
Write-Host "Загрузка скрипта установки..."
scp "deploy\install.sh" "${SERVER_USER}@${SERVER_IP}:/root/"

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Скрипт установки загружен" -ForegroundColor Green
} else {
    Write-Host "⚠️  Предупреждение: не удалось загрузить install.sh" -ForegroundColor Yellow
}

Write-Host ""

# ============================================================
# 3. Удаление локального архива
# ============================================================
Write-Host "🧹 Шаг 3: Очистка..."
Remove-Item $archiveName
Write-Host "✅ Локальный архив удален" -ForegroundColor Green
Write-Host ""

# ============================================================
# 4. Инструкции
# ============================================================
Write-Host "=================================================="
Write-Host "  ✅ Загрузка завершена!"
Write-Host "=================================================="
Write-Host ""
Write-Host "📋 Следующие шаги на сервере:"
Write-Host ""
Write-Host "1️⃣  Подключитесь к серверу:" -ForegroundColor Cyan
Write-Host "   ssh $SERVER_USER@$SERVER_IP"
Write-Host ""
Write-Host "2️⃣  Запустите автоматическую установку:" -ForegroundColor Cyan
Write-Host "   chmod +x /root/install.sh"
Write-Host "   /root/install.sh"
Write-Host ""
Write-Host "3️⃣  Распакуйте проект:" -ForegroundColor Cyan
Write-Host "   cd /opt/news-portal"
Write-Host "   tar -xzf /root/$archiveName"
Write-Host ""
Write-Host "4️⃣  Запустите сервисы:" -ForegroundColor Cyan
Write-Host "   chmod +x deploy/*.sh"
Write-Host "   ./deploy/start.sh"
Write-Host ""
Write-Host "=================================================="
Write-Host "🌐 После запуска проект будет доступен:"
Write-Host "   http://$SERVER_IP:8091  (Auth Service)"
Write-Host "   http://$SERVER_IP:8092  (News Service)"
Write-Host "   http://$SERVER_IP:8094  (Media Service)"
Write-Host "=================================================="
Write-Host ""

# ============================================================
# 5. Автоматическое подключение (опционально)
# ============================================================
$connect = Read-Host "Подключиться к серверу сейчас? (y/n)"

if ($connect -eq "y" -or $connect -eq "Y") {
    Write-Host ""
    Write-Host "🔌 Подключение к серверу..." -ForegroundColor Cyan
    Write-Host ""
    ssh "${SERVER_USER}@${SERVER_IP}"
}

Write-Host ""
Write-Host "✅ Готово!" -ForegroundColor Green
