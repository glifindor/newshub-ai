# Быстрый старт Admin Service

## 1. Запуск Backend

```powershell
cd "c:\Users\Grifindor\Desktop\НОВСТНОЙ САЙТ\admin-service"

# Убедитесь, что .env настроен правильно
# Запустить сервер
go run cmd/admin-service/main.go
```

Backend будет доступен на **http://localhost:8084**

## 2. Запуск Frontend (Development)

```powershell
cd "c:\Users\Grifindor\Desktop\НОВСТНОЙ САЙТ\admin-service\frontend"

# Dev server с HMR
npm run dev
```

Frontend будет доступен на **http://localhost:5173**

## 3. Тестирование

### Проверка health
```powershell
curl http://localhost:8084/health
```

### Логин (получить JWT токен)
```powershell
curl -X POST http://localhost:8084/api/admin/login `
  -H "Content-Type: application/json" `
  -d '{\"username\": \"admin\", \"password\": \"admin123\"}'
```

### Получить дашборд
```powershell
curl http://localhost:8084/api/admin/dashboard `
  -H "Authorization: Bearer <ваш_токен>"
```

## 4. Production Build

```powershell
# Собрать frontend
cd frontend
npm run build

# Запустить backend (он будет раздавать frontend/dist)
cd ..
go run cmd/admin-service/main.go

# Открыть http://localhost:8084
```

## 5. Docker

```powershell
# Из корневой директории проекта
cd "c:\Users\Grifindor\Desktop\НОВСТНОЙ САЙТ"

# Запустить все сервисы
docker-compose up -d admin-service

# Логи
docker logs -f admin-service
```

## Переменные окружения (.env)

```env
# Минимальная конфигурация для локальной разработки
SERVER_PORT=8084
SERVER_ENV=development
AUTH_SERVICE_URL=http://localhost:8081
NEWS_SERVICE_URL=http://localhost:8082
SEO_SERVICE_URL=http://localhost:8093
JWT_SECRET=dev-secret-key
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:8084
```

## Требуемые сервисы

Для работы admin-service должны быть запущены:

- ✅ **auth-service** на :8081 (авторизация)
- ✅ **news-service** на :8082 (CRUD новостей)
- ✅ **seo-service** на :8093 (SEO метаданные)

## Аккаунт по умолчанию

- **Username**: admin
- **Password**: admin123
- **Role**: admin

(Настраивается в auth-service)

## Структура URL

### Backend API
- http://localhost:8084/api/admin/login
- http://localhost:8084/api/admin/dashboard
- http://localhost:8084/api/admin/news
- http://localhost:8084/health

### Frontend (dev)
- http://localhost:5173/login
- http://localhost:5173/ (dashboard)
- http://localhost:5173/news
- http://localhost:5173/news/new

### Frontend (production)
- http://localhost:8084/ - админ-панель обслуживается backend

## Troubleshooting

### Backend не запускается
```powershell
# Проверить зависимости
go mod tidy

# Проверить .env
cat .env
```

### Frontend ошибки компиляции
```powershell
# Переустановить зависимости
rm -r node_modules
rm package-lock.json
npm install
```

### Ошибка 401 Unauthorized
- Проверить, что auth-service запущен на :8081
- Проверить JWT_SECRET в .env (должен совпадать с auth-service)
- Проверить токен (возможно истек)

### CORS ошибка
- Добавить http://localhost:5173 в CORS_ALLOWED_ORIGINS в .env
- Перезапустить backend

## Полезные команды

```powershell
# Проверка версий
go version  # >= 1.23
node -v     # >= 20

# Компиляция backend
go build -o admin-service.exe .\cmd\admin-service\main.go

# Запуск скомпилированного
.\admin-service.exe

# Frontend production build
cd frontend
npm run build
npm run preview  # preview build локально
```
