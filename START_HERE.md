# ===================================================================
# ✅ ФИНАЛЬНАЯ ИНСТРУКЦИЯ - ЗАПУСК НОВОСТНОГО ПОРТАЛА
# ===================================================================

## 📊 ТЕКУЩИЙ СТАТУС

✅ **Выполнено:**
- Проект загружен на сервер `151.241.228.203` в `/opt/news-portal`
- Docker 28.5.1 установлен
- PostgreSQL, Redis, MinIO запущены
- Firewall настроен (порты 22, 80, 443, 8091-8094)
- Dockerfiles обновлены (Go 1.23)
- Пароли сгенерированы в `/opt/news-portal/PASSWORDS.txt`

⏳ **Сейчас выполняется:**
- Сборка Docker образов для 3 сервисов (процесс ID: 10207)
- Ожидаемое время: 3-5 минут

---

## 🎯 ЧТО ДЕЛАТЬ ДАЛЬШЕ

### ШАГ 1: Дождаться завершения сборки (3-5 минут)

Подключитесь к серверу:

```powershell
ssh root@151.241.228.203
```

Проверьте статус сборки:

```bash
# Проверить, работает ли процесс сборки
ps aux | grep "docker compose build" | grep -v grep

# Посмотреть лог сборки (последние 30 строк)
tail -30 /tmp/full-build.log

# Следить за логом в реальном времени
tail -f /tmp/full-build.log
# Нажмите Ctrl+C когда увидите "FINISHED" или "Successfully"
```

**Признаки успеха:**
- В логе появится текст `Successfully built` или `FINISHED`
- Процесс `docker compose build` завершится

**Признаки ошибки:**
- В логе появится `ERROR` или `failed to solve`
- Нужно будет отправить мне лог ошибки

---

### ШАГ 2: Запустить сервисы

После завершения сборки:

```bash
cd /opt/news-portal

# Запустить ВСЕ сервисы
docker compose up -d

# ИЛИ запустить только микросервисы (если инфраструктура уже работает)
docker compose up -d auth-service news-service media-service
```

---

### ШАГ 3: Проверить работу

```bash
# Проверить статус всех контейнеров
docker compose ps

# Должны быть запущены (STATUS: Up):
# - news-postgres
# - news-redis  
# - news-minio
# - news-portal-auth-service-1
# - news-portal-news-service-1
# - news-portal-media-service-1
```

Проверить логи:

```bash
# Все логи
docker compose logs --tail=50

# Логи конкретного сервиса
docker compose logs auth-service
docker compose logs news-service
docker compose logs media-service
```

Проверить здоровье API:

```bash
# На сервере
curl http://localhost:8091/health    # Auth Service
curl http://localhost:8092/health    # News Service  
curl http://localhost:8094/health    # Media Service
```

**Ожидаемый ответ:** `{"status":"ok"}` или `{"status":"healthy"}`

---

### ШАГ 4: Проверить доступ извне

С вашего Windows компьютера (PowerShell):

```powershell
curl http://151.241.228.203:8091/health
curl http://151.241.228.203:8092/health
curl http://151.241.228.203:8094/health
```

**Если не работает** - откройте порты:

```bash
sudo ufw allow 8091/tcp
sudo ufw allow 8092/tcp
sudo ufw allow 8094/tcp
sudo ufw reload
```

---

## 🔐 ДОСТУПЫ И ПАРОЛИ

Все пароли сохранены на сервере:

```bash
cat /opt/news-portal/PASSWORDS.txt
```

### PostgreSQL
- **Host:** `151.241.228.203:5432`
- **Database:** `newsportal_db`
- **User:** `newsportal`
- **Password:** (см. PASSWORDS.txt)

### Redis
- **Host:** `151.241.228.203:6379`
- **Password:** (см. PASSWORDS.txt)

### MinIO (S3 Storage)
- **Console:** `http://151.241.228.203:9001`
- **API:** `http://151.241.228.203:9000`
- **User:** `newsportal_admin`
- **Password:** (см. PASSWORDS.txt)

---

## 🌐 API ЭНДПОИНТЫ

### Auth Service (`:8091`)
```
POST   /api/v1/register              - Регистрация
POST   /api/v1/login                 - Вход
POST   /api/v1/logout                - Выход
GET    /api/v1/profile               - Профиль (требует JWT)
POST   /api/v1/refresh-token         - Обновить токен
POST   /api/v1/change-password       - Смена пароля (требует JWT)
```

### News Service (`:8092`)
```
# Категории
GET    /api/v1/categories            - Список категорий
POST   /api/v1/categories            - Создать категорию (admin)
GET    /api/v1/categories/:id        - Получить категорию
PUT    /api/v1/categories/:id        - Обновить категорию (admin)
DELETE /api/v1/categories/:id        - Удалить категорию (admin)

# Теги
GET    /api/v1/tags                  - Список тегов
POST   /api/v1/tags                  - Создать тег (editor)
GET    /api/v1/tags/:id              - Получить тег
PUT    /api/v1/tags/:id              - Обновить тег (editor)
DELETE /api/v1/tags/:id              - Удалить тег (editor)

# Новости
GET    /api/v1/news                  - Список новостей
POST   /api/v1/news                  - Создать новость (editor)
GET    /api/v1/news/:id              - Получить новость
PUT    /api/v1/news/:id              - Обновить новость (editor)
DELETE /api/v1/news/:id              - Удалить новость (editor)
GET    /api/v1/news/featured         - Избранные новости
GET    /api/v1/news/breaking         - Срочные новости
```

### Media Service (`:8094`)
```
POST   /api/v1/upload                - Загрузить файл (требует JWT)
GET    /api/v1/files/:id             - Информация о файле
DELETE /api/v1/files/:id             - Удалить файл (требует JWT)
GET    /api/v1/files/:id/url         - Получить URL файла
GET    /api/v1/files                 - Список файлов (требует JWT)
```

---

## 🧪 ТЕСТИРОВАНИЕ API

### 1. Регистрация пользователя

```bash
curl -X POST http://151.241.228.203:8091/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "SecurePass123!",
    "full_name": "Test User"
  }'
```

### 2. Вход

```bash
curl -X POST http://151.241.228.203:8091/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'
```

**Сохраните `access_token` из ответа!**

### 3. Создать категорию (требует admin роли)

```bash
curl -X POST http://151.241.228.203:8092/api/v1/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "Технологии",
    "slug": "tech",
    "description": "Новости технологий"
  }'
```

### 4. Создать новость

```bash
curl -X POST http://151.241.228.203:8092/api/v1/news \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "title": "Первая новость",
    "slug": "first-news",
    "content": "Содержимое первой новости",
    "excerpt": "Краткое описание",
    "category_id": 1,
    "status": "published"
  }'
```

---

## 🛠️ УПРАВЛЕНИЕ СЕРВИСАМИ

```bash
# Остановить все
docker compose down

# Запустить все
docker compose up -d

# Перезапустить конкретный сервис
docker compose restart auth-service

# Просмотр логов
docker compose logs -f auth-service

# Статус
docker compose ps

# Использование ресурсов
docker stats
```

---

## 🆘 TROUBLESHOOTING

### Проблема: Сервис не запускается

```bash
# Посмотреть подробные логи
docker compose logs auth-service --tail=100

# Проверить конфигурацию
docker compose config

# Пересобрать образ
docker compose build auth-service
docker compose up -d auth-service
```

### Проблема: База данных недоступна

```bash
# Проверить PostgreSQL
docker compose exec postgres psql -U newsportal -d newsportal_db -c "SELECT 1;"

# Проверить Redis
docker compose exec redis redis-cli PING

# Проверить подключения
docker compose exec auth-service ping postgres
```

### Проблема: Порт занят

```bash
# Проверить какой процесс использует порт
ss -tulpn | grep :8091

# Остановить контейнер
docker compose stop auth-service
```

---

## 📋 СЛЕДУЮЩИЕ ШАГИ

После успешного запуска:

1. ✅ Настроить Nginx reverse proxy
2. ✅ Установить SSL/HTTPS (Let's Encrypt)
3. ✅ Настроить доменное имя
4. ✅ Настроить мониторинг (Prometheus + Grafana)
5. ✅ Создать frontend (Next.js)
6. ✅ Создать Admin Panel (React)

Подробные инструкции в `DEPLOYMENT_GUIDE.md`

---

## 💾 БЭКАПЫ

Создать бэкап:

```bash
cd /opt/news-portal
./deploy/backup.sh
```

Бэкапы сохраняются в `/opt/backups/`

---

## 📞 НУЖНА ПОМОЩЬ?

Если возникли проблемы, отправьте мне:

```bash
# 1. Статус контейнеров
docker compose ps

# 2. Логи проблемного сервиса
docker compose logs auth-service --tail=50

# 3. Лог сборки (если была ошибка)
cat /tmp/full-build.log | tail -100
```

Удачи! 🚀
