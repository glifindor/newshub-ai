# News Service

Микросервис управления новостями, категориями и тегами.

## Функциональность

- ✅ CRUD операции с новостями
- ✅ Управление категориями
- ✅ Управление тегами
- ✅ Публикация/черновики/архив
- ✅ Генерация slug для SEO
- ✅ Кеширование популярных новостей
- ✅ Счетчик просмотров

## API (gRPC)

### News Operations
- `CreateNews` - создание новости
- `GetNews` - получение новости по ID
- `GetNewsBySlug` - получение новости по slug
- `ListNews` - список новостей с фильтрацией
- `UpdateNews` - обновление новости
- `DeleteNews` - удаление новости
- `PublishNews` - публикация новости

### Categories
- `CreateCategory` - создание категории
- `GetCategory` - получение категории
- `ListCategories` - список категорий

### Tags
- `CreateTag` - создание тега
- `GetTag` - получение тега
- `ListTags` - список тегов

## Установка

```bash
go mod download
protoc --go_out=. --go-grpc_out=. proto/news.proto
go run cmd/news-service/main.go
```

## Docker

```bash
docker build -t news-service .
docker run -p 8082:8082 --env-file .env news-service
```
