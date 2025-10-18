# 🚀 Инструкция по деплою на сервер 151.241.228.203

## Шаг 1: Подключение к серверу

```bash
# Подключитесь через SSH
ssh root@151.241.228.203

# Введите пароль
```

---

## Шаг 2: Первичная настройка сервера

```bash
# Обновление системы
apt update && apt upgrade -y

# Установка необходимых пакетов
apt install -y docker.io docker-compose git nano curl wget

# Запуск Docker
systemctl start docker
systemctl enable docker

# Проверка
docker --version
docker-compose --version
```

---

## Шаг 3: Клонирование проекта

```bash
# Создание директории
mkdir -p /opt/newshub
cd /opt/newshub

# Клонирование (ЗАМЕНИТЕ на ваш репозиторий!)
git clone https://github.com/glifindor/newsportal.git .

# Проверка файлов
ls -la
```

Вы должны увидеть:
- `backend/`
- `frontend/`
- `nginx/`
- `docker-compose.yml`
- `.env.example`

---

## Шаг 4: Настройка .env файла

```bash
# Копирование шаблона
cp .env.example .env

# Редактирование
nano .env
```

**Замените значения:**

```env
# 1. Пароль БД (придумайте свой!)
POSTGRES_PASSWORD=ВАШ_СУПЕР_СЕКРЕТНЫЙ_ПАРОЛЬ

# 2. JWT Secret (сгенерируйте случайную строку)
JWT_SECRET_KEY=СЛУЧАЙНАЯ_СТРОКА_32_СИМВОЛА

# 3. OpenRouter API
OPENROUTER_API_KEY=sk-or-v1-ваш_ключ

# 4. Freepik API
FREEPIK_API_KEY=ваш_ключ

# 5. NewsAPI
NEWSAPI_KEY=ваш_ключ

# 6. Telegram Bot
TELEGRAM_BOT_TOKEN=123456:ABCdef_ваш_токен
TELEGRAM_CRYPTO_CHANNEL=@crypto_ainews
TELEGRAM_POLITICS_CHANNEL=@kremlin_digest
TELEGRAM_ADMIN_CHAT_ID=ваш_id

# 7. Frontend URL
NEXT_PUBLIC_API_URL=http://151.241.228.203:8000/api
```

**Сохранить:** `Ctrl+X` → `Y` → `Enter`

---

## Шаг 5: Генерация JWT Secret

```bash
# Генерация случайной строки
openssl rand -hex 32

# Скопируйте вывод и вставьте в .env
nano .env
# Найдите JWT_SECRET_KEY= и вставьте сгенерированную строку
```

---

## Шаг 6: Запуск проекта

```bash
# Запуск всех сервисов
docker-compose up -d --build
```

Это займёт 5-10 минут. Ждите...

**Проверка статуса:**

```bash
docker-compose ps
```

Все контейнеры должны быть `Up`:
- `newshub_postgres`
- `newshub_redis`
- `newshub_rabbitmq`
- `newshub_backend`
- `newshub_frontend`
- `newshub_nginx`
- `newshub_celery_worker`
- `newshub_celery_beat`

---

## Шаг 7: Проверка логов

```bash
# Все логи
docker-compose logs -f

# Только backend
docker-compose logs -f backend

# Только frontend
docker-compose logs -f frontend

# Выход из логов: Ctrl+C
```

---

## Шаг 8: Создание администратора

```bash
# Войти в контейнер backend
docker-compose exec backend bash

# Запустить скрипт
python scripts/create_admin.py

# Ввести данные:
# Username: admin
# Email: your@email.com
# Password: (придумайте надёжный пароль)

# Выйти из контейнера
exit
```

---

## Шаг 9: Проверка работы

Откройте в браузере:

1. **Главная страница:** http://151.241.228.203
2. **API документация:** http://151.241.228.203:8000/docs
3. **Grafana (мониторинг):** http://151.241.228.203:3001
   - Логин: `admin`
   - Пароль: `admin` (смените!)
4. **RabbitMQ:** http://151.241.228.203:15672
   - Логин: `guest`
   - Пароль: `guest`

---

## Шаг 10: Тестирование API

```bash
# Health check
curl http://151.241.228.203:8000/health

# Ожидаемый ответ:
# {"status":"healthy","database":"connected","redis":"connected","celery":"running"}
```

---

## Шаг 11: Тестирование Telegram Bot

```bash
# Войти в backend
docker-compose exec backend bash

# Запустить тест
python scripts/test_telegram.py

# Вы должны получить сообщение в Telegram!
exit
```

---

## Шаг 12: Тестирование OpenRouter

```bash
docker-compose exec backend bash
python scripts/test_openrouter.py
exit
```

---

## 🔧 Полезные команды

### Перезапуск сервисов

```bash
# Перезапустить всё
docker-compose restart

# Перезапустить только backend
docker-compose restart backend
```

### Остановка

```bash
# Остановить всё (сохраняются данные)
docker-compose stop

# Остановить и удалить контейнеры
docker-compose down

# Остановить и удалить всё (включая volumes!)
docker-compose down -v
```

### Обновление проекта

```bash
# Скачать изменения
cd /opt/newshub
git pull

# Пересобрать и перезапустить
docker-compose up -d --build
```

### Просмотр логов

```bash
# Реал-тайм логи
docker-compose logs -f

# Последние 100 строк
docker-compose logs --tail=100

# Логи конкретного сервиса за последний час
docker-compose logs --since 1h backend
```

### Backup базы данных

```bash
# Создать backup
docker-compose exec postgres pg_dump -U newsadmin newshub_db > backup_$(date +%Y%m%d).sql

# Восстановить из backup
cat backup_20250118.sql | docker-compose exec -T postgres psql -U newsadmin newshub_db
```

---

## 🔐 Настройка Firewall (безопасность)

```bash
# Установка UFW
apt install ufw -y

# Разрешить SSH (ВАЖНО! Иначе потеряете доступ)
ufw allow 22/tcp

# Разрешить HTTP и HTTPS
ufw allow 80/tcp
ufw allow 443/tcp

# Включить firewall
ufw enable

# Проверка статуса
ufw status
```

---

## 🔒 Настройка SSL (HTTPS)

### Вариант 1: С доменом (рекомендуется)

```bash
# Остановить Nginx временно
docker-compose stop nginx

# Установка Certbot
apt install certbot -y

# Получение сертификата (ЗАМЕНИТЕ домен!)
certbot certonly --standalone -d newshub.example.com

# Сертификаты будут в: /etc/letsencrypt/live/newshub.example.com/

# Обновить nginx.conf с путями к сертификатам
nano nginx/nginx.conf

# Запустить Nginx
docker-compose start nginx
```

### Вариант 2: Self-signed (для тестирования)

```bash
# Создание директории
mkdir -p nginx/ssl

# Генерация сертификата
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/nginx-selfsigned.key \
  -out nginx/ssl/nginx-selfsigned.crt

# Введите данные (можно оставить по умолчанию)

# Обновить docker-compose.yml с маунтом SSL
```

---

## 📊 Мониторинг

### Grafana Dashboard

1. Откройте http://151.241.228.203:3001
2. Логин: `admin` / Пароль: `admin`
3. Смените пароль при первом входе!
4. Добавьте Prometheus как Data Source:
   - URL: `http://prometheus:9090`
   - Access: `Server`

### Prometheus

1. Откройте http://151.241.228.203:9090
2. Проверьте Targets: Status → Targets
3. Все должны быть `UP`

---

## 🆘 Решение проблем

### "Port 80 already in use"

```bash
# Найти процесс
lsof -i :80

# Остановить Apache (если установлен)
systemctl stop apache2
systemctl disable apache2
```

### "Cannot connect to database"

```bash
# Проверить статус PostgreSQL
docker-compose exec postgres pg_isready -U newsadmin

# Посмотреть логи
docker-compose logs postgres

# Перезапустить
docker-compose restart postgres
```

### "Telegram Bot не отвечает"

1. Проверьте токен в `.env`
2. Убедитесь, что бот добавлен в каналы как админ
3. Отправьте боту `/start` в личку
4. Проверьте TELEGRAM_ADMIN_CHAT_ID

### "OpenRouter ошибка 401"

1. Проверьте API ключ в `.env`
2. Убедитесь, что есть баланс на openrouter.ai
3. Попробуйте другую модель (более дешёвую)

### "Out of memory"

```bash
# Проверить RAM
free -h

# Добавить swap
fallocate -l 4G /swapfile
chmod 600 /swapfile
mkswap /swapfile
swapon /swapfile

# Сделать постоянным
echo '/swapfile none swap sw 0 0' >> /etc/fstab
```

---

## 📞 Получить помощь

1. Проверьте логи: `docker-compose logs -f`
2. Проверьте статус: `docker-compose ps`
3. Посмотрите ARCHITECTURE.md
4. Создайте Issue на GitHub

---

## ✅ Чеклист готовности

- [ ] Сервер доступен по SSH
- [ ] Docker и Docker Compose установлены
- [ ] Проект склонирован
- [ ] .env файл настроен
- [ ] Все контейнеры запущены (docker-compose ps)
- [ ] API отвечает (curl http://151.241.228.203:8000/health)
- [ ] Frontend открывается в браузере
- [ ] Админ создан
- [ ] Telegram Bot работает
- [ ] OpenRouter отвечает
- [ ] Firewall настроен
- [ ] Backup настроен

---

**Готово! Ваш NewsHub AI работает! 🎉**
