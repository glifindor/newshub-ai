# 🚀 START HERE - Первые шаги после клонирования

> Вы только что склонировали **NewsHub AI**. Что дальше?

---

## 📋 Выберите ваш сценарий:

### 🎯 Сценарий 1: "Я хочу БЫСТРО развернуть на production сервере"

**Время:** 15 минут  
**Сложность:** ⭐ Легко (всё автоматически)

```powershell
# Windows PowerShell
.\scripts\setup-interactive.ps1

# Скрипт сам:
# ✅ Установит Docker на сервер
# ✅ Настроит Firewall
# ✅ Склонирует проект
# ✅ Создаст .env файл
# ✅ Запустит все контейнеры
# ✅ Создаст администратора
```

**Что нужно:**
- Пароль от сервера (root@151.241.228.203)
- API ключи (OpenRouter, Telegram)
- 10-15 минут времени ☕

**Документация:** [QUICK_DEPLOY.md](./QUICK_DEPLOY.md)

---

### 💻 Сценарий 2: "Я хочу запустить локально для разработки"

**Время:** 10 минут  
**Сложность:** ⭐⭐ Средне

```bash
# 1. Скопировать .env
cp .env.example .env

# 2. Отредактировать .env (добавить API ключи)
nano .env

# 3. Запустить с Docker
docker-compose up -d --build

# 4. Открыть в браузере
# http://localhost:3000 - Frontend
# http://localhost:8000/docs - API
```

**Что нужно:**
- Docker Desktop установлен
- API ключи (OpenRouter, Telegram)

**Документация:** [QUICK_START.md](./QUICK_START.md)

---

### 🛠️ Сценарий 3: "Я разработчик, хочу запустить backend и frontend отдельно"

**Время:** 20 минут  
**Сложность:** ⭐⭐⭐ Продвинуто

**Backend:**
```bash
cd backend
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate
pip install -r requirements.txt

# Запустить БД через Docker
docker-compose up -d postgres redis rabbitmq

# Миграции
alembic upgrade head

# Запуск
uvicorn app.main:app --reload
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

**Celery Worker:**
```bash
cd backend
celery -A app.celery_app worker --loglevel=info
```

**Что нужно:**
- Python 3.11+
- Node.js 20+
- PostgreSQL 15
- Redis 7
- RabbitMQ 3.12

**Документация:** 
- [backend/README_BACKEND.md](./backend/README_BACKEND.md)
- [frontend/README.md](./frontend/README.md)

---

### 🤖 Сценарий 4: "Я хочу настроить CI/CD с GitHub Actions"

**Время:** 30 минут  
**Сложность:** ⭐⭐⭐⭐ Эксперт

1. **Создать Docker Hub аккаунт** - https://hub.docker.com
2. **Добавить GitHub Secrets:**
   - Settings → Secrets → Actions
   - См. список в [PRODUCTION_DEPLOYMENT.md](./PRODUCTION_DEPLOYMENT.md#настройка-github-secrets)
3. **Push в main** - автоматический деплой!

```bash
git push origin main
# GitHub Actions автоматически:
# ✅ Запустит тесты
# ✅ Соберет Docker images
# ✅ Задеплоит на сервер
# ✅ Отправит уведомление в Telegram
```

**Что нужно:**
- GitHub репозиторий
- Docker Hub аккаунт
- SSH доступ к серверу

**Документация:** [PRODUCTION_DEPLOYMENT.md](./PRODUCTION_DEPLOYMENT.md#cicd-с-github-actions)

---

## 🔑 Необходимые API ключи

Независимо от сценария, вам понадобятся:

### 1. OpenRouter API Key (обязательно)
```
Получить: https://openrouter.ai/keys
Цена: от $5/месяц (зависит от модели)
Зачем: AI-обработка новостей
```

### 2. Telegram Bot Token (обязательно)
```
Получить: https://t.me/BotFather
Команда: /newbot
Зачем: Публикация в Telegram каналы
```

### 3. Telegram Admin Chat ID (обязательно)
```
Получить: https://t.me/userinfobot
Команда: /start
Зачем: Уведомления об ошибках
```

### 4. NewsAPI Key (опционально)
```
Получить: https://newsapi.org
Free tier: 100 запросов/день
Зачем: Дополнительный источник новостей
```

### 5. Freepik API Key (опционально)
```
Получить: https://www.freepik.com/api
Зачем: Изображения для новостей
```

---

## 📁 Структура проекта

```
newsportal/
├── 📁 backend/              # FastAPI приложение
├── 📁 frontend/             # Next.js приложение
├── 📁 nginx/                # Nginx конфигурация
├── 📁 monitoring/           # Prometheus + Grafana
├── 📁 scripts/              # Deployment scripts
├── 📁 .github/workflows/    # CI/CD
├── docker-compose.yml       # Development
├── docker-compose.prod.yml  # Production
├── .env.example             # Environment template
└── 📚 Документация:
    ├── README.md            # ← Вы здесь
    ├── START_HERE.md        # ← Этот файл
    ├── QUICK_DEPLOY.md      # Быстрый деплой
    ├── QUICK_START.md       # Быстрый старт
    ├── PRODUCTION_DEPLOYMENT.md  # Production гайд
    ├── ARCHITECTURE.md      # Архитектура
    └── DEPLOYMENT_COMPLETE.md    # Итоговая сводка
```

---

## 🆘 Частые проблемы

### "ModuleNotFoundError" при запуске backend

```bash
# Убедитесь, что virtual environment активирован
source venv/bin/activate  # Linux/Mac
venv\Scripts\activate     # Windows

# Установите зависимости
pip install -r requirements.txt
```

### "ENOENT: no such file or directory" в frontend

```bash
# Установите зависимости
cd frontend
npm install
```

### "Cannot connect to Docker daemon"

```bash
# Запустите Docker Desktop (Windows/Mac)
# Или Docker service (Linux)
sudo systemctl start docker
```

### "Port 8000 already in use"

```bash
# Найдите процесс
lsof -i :8000  # Linux/Mac
netstat -ano | findstr :8000  # Windows

# Остановите процесс или измените порт
```

---

## 📚 Где искать помощь

### Документация по темам

| Тема | Файл | Описание |
|------|------|----------|
| 🚀 Быстрый деплой | [QUICK_DEPLOY.md](./QUICK_DEPLOY.md) | 15 минут до production |
| ⚡ Локальный запуск | [QUICK_START.md](./QUICK_START.md) | Development setup |
| 🏗️ Production | [PRODUCTION_DEPLOYMENT.md](./PRODUCTION_DEPLOYMENT.md) | Полный гайд (1000+ строк) |
| 🏛️ Архитектура | [ARCHITECTURE.md](./ARCHITECTURE.md) | Как устроена система |
| 📱 Telegram | [backend/TELEGRAM_BOT_SETUP.md](./backend/TELEGRAM_BOT_SETUP.md) | Настройка бота |
| 🎨 Frontend | [frontend/README.md](./frontend/README.md) | Frontend документация |
| ✅ Итоги | [DEPLOYMENT_COMPLETE.md](./DEPLOYMENT_COMPLETE.md) | Что создано |

### Каналы поддержки

- 💬 **Telegram:** @newshub_support
- 📧 **Email:** support@newshub.ai
- 🐛 **GitHub Issues:** https://github.com/glifindor/newsportal/issues
- 📖 **Wiki:** https://github.com/glifindor/newsportal/wiki

---

## ✅ Чеклист "Я готов начать"

- [ ] Проект склонирован
- [ ] Документация прочитана (хотя бы README.md)
- [ ] Выбран сценарий (1-4)
- [ ] API ключи получены
- [ ] Необходимое ПО установлено
- [ ] .env файл создан и заполнен
- [ ] Docker запущен (если используется)

**Если все ✅ - можно начинать!**

---

## 🎯 Быстрые команды

### Я хочу сразу начать (Docker)

```bash
cp .env.example .env
nano .env  # Добавить API ключи
docker-compose up -d --build
```

### Я хочу развернуть на production

```powershell
.\scripts\setup-interactive.ps1
```

### Я хочу увидеть документацию

```bash
ls *.md  # Список всех .md файлов
cat README.md  # Прочитать README
```

### Я хочу увидеть логи

```bash
docker-compose logs -f
docker-compose logs -f backend
docker-compose logs -f frontend
```

---

## 🚀 Что дальше?

После успешного запуска:

1. **Откройте в браузере:**
   - Frontend: http://localhost:3000
   - API Docs: http://localhost:8000/docs
   - Grafana: http://localhost:3001

2. **Создайте администратора:**
   ```bash
   docker-compose exec backend python scripts/create_admin.py
   ```

3. **Добавьте источники новостей:**
   - Админ-панель → Sources → Add Source

4. **Запустите сбор новостей:**
   - API Docs → /api/pipeline/pipeline → Execute

5. **Проверьте Telegram каналы:**
   - @crypto_ainews
   - @kremlin_digest

---

## 💡 Советы

- 📖 **Читайте документацию** - там всё подробно расписано
- 🐛 **Проверяйте логи** - `docker-compose logs -f`
- ✅ **Следуйте чеклистам** - они в каждом .md файле
- 💬 **Задавайте вопросы** - в Issues или Telegram
- 🔄 **Обновляйте проект** - `git pull && docker-compose up -d --build`

---

## 🎉 Готовы начать?

**Выберите ваш сценарий выше и следуйте инструкциям!**

Удачи! 🚀

---

<p align="center">
  <b>Сделано с ❤️ и ☕ командой NewsHub AI</b>
</p>
