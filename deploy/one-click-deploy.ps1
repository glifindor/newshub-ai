# ============================================================
# Скрипт полного деплоя одной командой (Windows PowerShell)
# ============================================================

$SERVER_IP = "151.241.228.203"
$SERVER_USER = "root"
$PROJECT_DIR = "C:\Users\Grifindor\Desktop\НОВСТНОЙ САЙТ"

Write-Host "=================================================="
Write-Host "  News Portal - Automatic Deploy"
Write-Host "=================================================="
Write-Host ""

Set-Location $PROJECT_DIR

# ============================================================
# 1. Package
# ============================================================
Write-Host "Packaging project..." -ForegroundColor Cyan
$archiveName = "news-portal.tar.gz"

tar -czf $archiveName auth-service news-service media-service docker-compose.yml deploy *.md

if (-not (Test-Path $archiveName)) {
    Write-Host "Packaging error" -ForegroundColor Red
    exit 1
}

Write-Host "Packaged successfully" -ForegroundColor Green
Write-Host ""

# ============================================================
# 2. Upload
# ============================================================
Write-Host "Uploading to server..." -ForegroundColor Cyan

scp $archiveName "${SERVER_USER}@${SERVER_IP}:/root/"
scp "deploy\install.sh" "${SERVER_USER}@${SERVER_IP}:/root/"

Write-Host "Uploaded successfully" -ForegroundColor Green
Write-Host ""

# Cleanup
Remove-Item $archiveName

# ============================================================
# 3. Install on server
# ============================================================
Write-Host "Installing on server..." -ForegroundColor Cyan
Write-Host ""

$installCommands = @"
# Запуск автоматической установки
chmod +x /root/install.sh
/root/install.sh

# Создание директории проекта
mkdir -p /opt/news-portal
cd /opt/news-portal

# Распаковка
tar -xzf /root/news-portal.tar.gz

# Make scripts executable
chmod +x deploy/*.sh

echo ""
echo "=================================================="
echo "  Installation completed!"
echo "=================================================="
echo ""
echo "Next step - start services:"
echo "   cd /opt/news-portal"
echo "   ./deploy/start.sh"
echo ""
"@

ssh "${SERVER_USER}@${SERVER_IP}" $installCommands

Write-Host ""
Write-Host "=================================================="
Write-Host "  Deploy completed successfully!"
Write-Host "=================================================="
Write-Host ""
Write-Host "Services available at:"
Write-Host "   Auth:  http://$SERVER_IP:8091"
Write-Host "   News:  http://$SERVER_IP:8092"
Write-Host "   Media: http://$SERVER_IP:8094"
Write-Host ""
Write-Host "Management commands:"
Write-Host "   ssh $SERVER_USER@$SERVER_IP"
Write-Host "   cd /opt/news-portal"
Write-Host "   ./deploy/start.sh   - Start services"
Write-Host "   ./deploy/stop.sh    - Stop services"
Write-Host "   ./deploy/backup.sh  - Backup data"
Write-Host ""
