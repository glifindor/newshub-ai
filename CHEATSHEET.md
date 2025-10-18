# 📝 NewsHub AI - Шпаргалка команд

> Быстрый справочник по всем командам

---

## 🚀 Deployment

### Автоматический деплой на production (15 минут)
```powershell
.\scripts\setup-interactive.ps1
```

### Деплой через Docker (локально)
```bash
docker-compose up -d --build
```

### Обновление production
```bash
cd /opt/newshub
git pull
docker-compose -f docker-compose.prod.yml up -d --build
```

---

## 🐳 Docker команды

### Запуск
```bash
# Development
docker-compose up -d

# Production
docker-compose -f docker-compose.prod.yml up -d

# С пересборкой
docker-compose up -d --build

# В foreground (с логами)
docker-compose up
```

### Остановка
```bash
# Остановить
docker-compose stop

# Остановить и удалить контейнеры
docker-compose down

# Удалить с volumes (ОСТОРОЖНО!)
docker-compose down -v
```

### Логи
```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
docker-compose logs -f celery_worker

# Последние 100 строк
docker-compose logs --tail=100 backend
```

### Статус
```bash
# Список контейнеров
docker-compose ps

# Детальная информация
docker-compose ps -a

# Использование ресурсов
docker stats
```

### Вход в контейнер
```bash
# Backend
docker-compose exec backend bash

# Frontend
docker-compose exec frontend sh

# PostgreSQL
docker-compose exec postgres psql -U newshub
```

### Очистка
```bash
# Удалить остановленные контейнеры
docker container prune

# Удалить неиспользуемые images
docker image prune -a

# Удалить volumes
docker volume prune

# Полная очистка (ОСТОРОЖНО!)
docker system prune -a --volumes
```

---

## 🗄️ База данных

### Подключение к PostgreSQL
```bash
# Из контейнера
docker-compose exec postgres psql -U newshub

# С хоста (если порт пробросан)
psql -h localhost -p 5432 -U newshub -d newshub
```

### Миграции
```bash
# Создать миграцию
docker-compose exec backend alembic revision --autogenerate -m "description"

# Применить миграции
docker-compose exec backend alembic upgrade head

# Откатить на 1 шаг
docker-compose exec backend alembic downgrade -1

# История миграций
docker-compose exec backend alembic history

# Текущая версия
docker-compose exec backend alembic current
```

### Backup
```bash
# Создать backup
./scripts/backup.sh

# Production backup через скрипт
ssh root@151.241.228.203 "/opt/newshub/scripts/backup.sh"

# Manual backup
docker-compose exec postgres pg_dump -U newshub newshub > backup_$(date +%Y%m%d_%H%M%S).sql

# Restore
docker-compose exec -T postgres psql -U newshub newshub < backup.sql
```

---

## 🔧 Backend

### Запуск backend локально (без Docker)
```bash
cd backend

# Создать venv
python -m venv venv

# Активировать
source venv/bin/activate  # Linux/Mac
venv\Scripts\activate     # Windows

# Установить зависимости
pip install -r requirements.txt

# Миграции
alembic upgrade head

# Запуск
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

### Celery
```bash
# Worker
celery -A app.celery_app worker --loglevel=info

# Beat (scheduler)
celery -A app.celery_app beat --loglevel=info

# Flower (monitoring)
celery -A app.celery_app flower

# Через Docker
docker-compose exec celery_worker celery -A app.celery_app inspect active
docker-compose exec celery_worker celery -A app.celery_app inspect stats
```

### Тестирование
```bash
# Все тесты
pytest

# С coverage
pytest --cov=app --cov-report=html

# Конкретный файл
pytest tests/test_api.py

# Конкретный тест
pytest tests/test_api.py::test_health_check

# С выводом print
pytest -s
```

### Линтеры и форматирование
```bash
# Black (форматирование)
black .

# isort (импорты)
isort .

# Flake8 (линтер)
flake8 .

# MyPy (типы)
mypy app/

# Всё сразу
black . && isort . && flake8 . && pytest
```

### Создание администратора
```bash
# Интерактивно
docker-compose exec backend python scripts/create_admin.py

# С параметрами
docker-compose exec backend python -c "
from app.database import SessionLocal
from app.models import User
from app.utils.auth import hash_password

db = SessionLocal()
admin = User(
    username='admin',
    email='admin@example.com',
    hashed_password=hash_password('admin123'),
    is_active=True,
    is_superuser=True
)
db.add(admin)
db.commit()
print('Admin created!')
"
```

---

## 🎨 Frontend

### Запуск frontend локально
```bash
cd frontend

# Установить зависимости
npm install

# Development
npm run dev

# Production build
npm run build
npm start

# Lint
npm run lint

# Type check
npm run type-check
```

### Переменные окружения
```bash
# Development
cp .env.example .env.local

# Production
cp .env.example .env.production
```

---

## 📊 Мониторинг

### Prometheus
```bash
# Открыть в браузере
http://localhost:9090

# Query examples
up
rate(http_requests_total[5m])
celery_task_queue_length

# Config check
docker-compose exec prometheus promtool check config /etc/prometheus/prometheus.yml
```

### Grafana
```bash
# Открыть в браузере
http://localhost:3001

# Логин: admin
# Пароль: admin123 (или из .env)

# Reset password
docker-compose exec grafana grafana-cli admin reset-admin-password newpassword
```

### Flower (Celery)
```bash
# Открыть в браузере
http://localhost:5555
```

### Логи приложения
```bash
# Realtime логи
docker-compose logs -f backend

# Grep ошибки
docker-compose logs backend | grep ERROR

# Последние ошибки
docker-compose logs --tail=100 backend | grep ERROR

# Экспорт логов
docker-compose logs backend > backend.log
```

---

## 🌐 Nginx

### Проверка конфигурации
```bash
# Test config
docker-compose exec nginx nginx -t

# Reload config (без downtime)
docker-compose exec nginx nginx -s reload

# Restart
docker-compose restart nginx
```

### Логи Nginx
```bash
# Access log
docker-compose exec nginx tail -f /var/log/nginx/access.log

# Error log
docker-compose exec nginx tail -f /var/log/nginx/error.log

# Grep 404
docker-compose exec nginx grep "404" /var/log/nginx/access.log
```

---

## 🔐 SSL/HTTPS

### Let's Encrypt (автоматический)
```bash
# Установить Certbot
sudo apt-get install certbot python3-certbot-nginx

# Получить сертификат
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# Auto-renewal check
sudo certbot renew --dry-run

# Cron для auto-renewal (добавить в crontab)
0 0 * * * certbot renew --quiet
```

### Self-signed (для тестирования)
```bash
# Создать сертификат
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /etc/nginx/ssl/nginx.key \
  -out /etc/nginx/ssl/nginx.crt
```

---

## 🤖 CI/CD

### GitHub Actions

#### Локальный запуск workflow (act)
```bash
# Установить act
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Запустить workflow
act -j test-backend

# С secrets
act -j deploy --secret-file .secrets
```

#### Триггеры
```bash
# Manual trigger
gh workflow run deploy.yml

# Проверить статус
gh run list

# Логи последнего run
gh run view --log
```

---

## 📦 Telegram бот

### Проверка бота
```bash
# Получить информацию о боте
curl "https://api.telegram.org/bot<BOT_TOKEN>/getMe"

# Получить обновления
curl "https://api.telegram.org/bot<BOT_TOKEN>/getUpdates"

# Отправить сообщение
curl -X POST "https://api.telegram.org/bot<BOT_TOKEN>/sendMessage" \
  -d "chat_id=<CHAT_ID>&text=Hello"
```

### Webhook (альтернатива polling)
```bash
# Установить webhook
curl -X POST "https://api.telegram.org/bot<BOT_TOKEN>/setWebhook" \
  -d "url=https://yourdomain.com/api/telegram/webhook"

# Удалить webhook
curl -X POST "https://api.telegram.org/bot<BOT_TOKEN>/deleteWebhook"

# Проверить webhook
curl "https://api.telegram.org/bot<BOT_TOKEN>/getWebhookInfo"
```

---

## 🐛 Debugging

### Проверка здоровья сервисов
```bash
# Backend health
curl http://localhost:8000/health

# Frontend health
curl http://localhost:3000/api/health

# PostgreSQL
docker-compose exec postgres pg_isready -U newshub

# Redis
docker-compose exec redis redis-cli ping
```

### Проверка портов
```bash
# Какие порты слушают
netstat -tulpn | grep LISTEN

# Конкретный порт
lsof -i :8000

# Windows
netstat -ano | findstr :8000
```

### Memory и CPU
```bash
# Docker stats
docker stats

# Конкретный контейнер
docker stats newsportal_backend_1

# Top в контейнере
docker-compose exec backend top

# Memory usage
docker-compose exec backend free -h
```

### Network debugging
```bash
# Ping между контейнерами
docker-compose exec backend ping postgres

# Curl между контейнерами
docker-compose exec backend curl http://frontend:3000

# DNS lookup
docker-compose exec backend nslookup postgres
```

---

## 🔧 Утилиты

### Генерация секретов
```bash
# JWT secret
openssl rand -hex 32

# Random password (32 chars)
openssl rand -base64 32

# UUID
python -c "import uuid; print(uuid.uuid4())"
```

### Проверка .env
```bash
# Вывести все переменные
docker-compose config

# Проверить конкретную переменную
docker-compose exec backend env | grep DATABASE_URL
```

### Git
```bash
# Текущая ветка
git branch

# Последние коммиты
git log --oneline -10

# Изменения
git status

# Pull и deploy
git pull && docker-compose up -d --build
```

---

## 📊 Production мониторинг

### Проверка сервисов production
```bash
# SSH на сервер
ssh root@151.241.228.203

# Статус контейнеров
cd /opt/newshub
docker-compose -f docker-compose.prod.yml ps

# Логи
docker-compose -f docker-compose.prod.yml logs -f --tail=100

# Использование ресурсов
docker stats

# Disk space
df -h

# Memory
free -h
```

### Бэкапы production
```bash
# Создать backup
/opt/newshub/scripts/backup.sh

# Список backup'ов
ls -lh /opt/newshub/backups/

# Скачать backup на локальную машину
scp root@151.241.228.203:/opt/newshub/backups/latest/postgres.sql.gz ./
```

### Restart сервисов
```bash
# Restart всех сервисов
docker-compose -f docker-compose.prod.yml restart

# Restart конкретного сервиса
docker-compose -f docker-compose.prod.yml restart backend

# Полная перезагрузка (с пересборкой)
docker-compose -f docker-compose.prod.yml up -d --build --force-recreate
```

---

## 🚨 Emergency команды

### Откат к предыдущей версии
```bash
# Git rollback
git reset --hard HEAD~1
docker-compose up -d --build

# Docker rollback
docker-compose down
docker-compose pull newsportal_backend:previous
docker-compose up -d
```

### Восстановление из backup
```bash
# Остановить контейнеры
docker-compose down

# Восстановить БД
docker-compose up -d postgres
sleep 10
docker-compose exec -T postgres psql -U newshub newshub < backup.sql

# Запустить остальное
docker-compose up -d
```

### Очистка места на диске
```bash
# Docker cleanup
docker system prune -a --volumes

# Логи
find /var/log -type f -name "*.log" -mtime +7 -delete

# Старые backups
find /opt/newshub/backups -type f -mtime +30 -delete
```

---

## 📚 Полезные ссылки

### Документация
- [START_HERE.md](./START_HERE.md) - Начните отсюда
- [QUICK_DEPLOY.md](./QUICK_DEPLOY.md) - Быстрый деплой
- [PRODUCTION_DEPLOYMENT.md](./PRODUCTION_DEPLOYMENT.md) - Production гайд
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Архитектура

### API Endpoints
```
GET  /health                    - Health check
GET  /docs                      - OpenAPI docs
GET  /api/news                  - Список новостей
POST /api/news                  - Создать новость
GET  /api/news/{id}             - Детали новости
POST /api/pipeline/pipeline     - Запуск pipeline
GET  /api/telegram/channels     - Telegram каналы
```

### Web интерфейсы
```
http://localhost:3000           - Frontend
http://localhost:8000           - Backend API
http://localhost:8000/docs      - API документация
http://localhost:3001           - Grafana (admin/admin123)
http://localhost:5555           - Flower (Celery monitoring)
http://localhost:9090           - Prometheus
http://localhost:15672          - RabbitMQ (guest/guest)
```

---

## 💡 Советы

### Частые ошибки

**Port already in use**
```bash
# Найти процесс
lsof -i :8000  # или netstat -ano | findstr :8000

# Убить процесс
kill -9 <PID>  # или taskkill /PID <PID> /F
```

**Database connection refused**
```bash
# Проверить, запущен ли PostgreSQL
docker-compose ps postgres

# Перезапустить
docker-compose restart postgres

# Проверить логи
docker-compose logs postgres
```

**Out of memory**
```bash
# Увеличить swap
sudo fallocate -l 4G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# Или уменьшить workers
# В .env: CELERY_WORKER_CONCURRENCY=2
```

### Best practices

- ✅ Всегда делайте backup перед обновлением
- ✅ Проверяйте логи после deployment
- ✅ Используйте `.env` файлы, не hardcode секреты
- ✅ Следите за использованием ресурсов
- ✅ Настройте alerts в Grafana
- ✅ Делайте регулярные backup'ы
- ✅ Обновляйте зависимости (`pip list --outdated`)

---

<p align="center">
  <b>Сохраните эту шпаргалку в закладки! 📌</b>
</p>
