# API Gateway

Единая точка входа для всех микросервисов новостного портала.

## Функциональность

- ✅ Маршрутизация запросов к микросервисам
- ✅ JWT аутентификация и авторизация
- ✅ Rate limiting
- ✅ CORS настройка
- ✅ Логирование запросов
- ✅ Проверка прав доступа

## Endpoints

### Auth
- `POST /api/auth/register` - регистрация
- `POST /api/auth/login` - вход
- `POST /api/auth/refresh` - обновление токена
- `POST /api/auth/logout` - выход (protected)

### News
- `GET /api/news` - список новостей
- `GET /api/news/:slug` - новость по slug
- `POST /api/news` - создание новости (protected, editor+)
- `PUT /api/news/:id` - обновление новости (protected, editor+)
- `DELETE /api/news/:id` - удаление новости (protected, admin)
- `POST /api/news/:id/publish` - публикация (protected, editor+)

### Categories & Tags
- `GET /api/categories` - список категорий
- `POST /api/categories` - создание категории (protected, admin)
- `GET /api/tags` - список тегов
- `POST /api/tags` - создание тега (protected, admin)

### Media
- `POST /api/media/upload` - загрузка файла (protected)

### Admin
- `GET /api/admin/users` - список пользователей (admin/moderator)
- `GET /api/admin/statistics` - статистика (admin/moderator)
- `POST /api/admin/moderate/:id` - модерация (admin/moderator)

## Запуск

```bash
go mod download
go run cmd/gateway/main.go
```

## Docker

```bash
docker build -t gateway .
docker run -p 8080:8080 --env-file .env gateway
```
