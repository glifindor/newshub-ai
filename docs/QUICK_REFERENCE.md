# 📊 Архитектура новостного портала - Краткое резюме

## 🎯 Общая концепция

Микросервисная архитектура новостного портала на **Golang** с **Next.js** frontend, использующая **gRPC** для межсервисного взаимодействия и **REST API** для клиентских запросов через API Gateway.

---

## 🏗️ Компоненты системы

```
┌─────────────┐
│  Frontend   │ Next.js (Port 3000)
│  (Next.js)  │
└──────┬──────┘
       │ HTTP/REST
       ▼
┌─────────────────────────────────────┐
│        API Gateway (Port 8080)      │
│  • Routing • Auth • Rate Limiting   │
└──┬────┬────┬────┬────┬─────────────┘
   │    │    │    │    │
   │ gRPC communication
   ▼    ▼    ▼    ▼    ▼
┌────┐┌────┐┌────┐┌────┐┌─────┐
│Auth││News││SEO ││Admn││Media│
│8081││8082││8083││8084││8085 │
└─┬──┘└─┬──┘└─┬──┘└─┬──┘└──┬──┘
  │     │     │     │      │
  └─────┴─────┴─────┴──────┘
         │
    ┌────┴─────┐
    ▼          ▼
┌──────┐   ┌──────┐
│Postgr│   │Redis │
│SQL   │   │Cache │
└──────┘   └──────┘
```

---

## 📦 Микросервисы

| Сервис | Port | Функции | Технологии |
|--------|------|---------|-----------|
| **auth-service** | 8081 | • Регистрация<br>• Авторизация<br>• JWT токены<br>• Управление ролями | Go, gRPC, PostgreSQL, Redis, JWT |
| **news-service** | 8082 | • CRUD новостей<br>• Категории<br>• Теги<br>• Поиск | Go, gRPC, PostgreSQL, Redis |
| **seo-service** | 8083 | • Метатеги<br>• Open Graph<br>• Sitemap<br>• Schema.org | Go, gRPC, Redis |
| **admin-service** | 8084 | • Модерация<br>• Аналитика<br>• Управление | Go, gRPC |
| **media-service** | 8085 | • Загрузка файлов<br>• Обработка изображений<br>• S3 storage | Go, gRPC, MinIO |
| **gateway** | 8080 | • API Gateway<br>• Auth middleware<br>• Rate limiting | Go, Gin, gRPC clients |
| **frontend** | 3000 | • UI<br>• SSR<br>• SEO | Next.js, React, TypeScript |

---

## 🔄 Поток данных

### Создание новости (пример)

```
1. User → Frontend
   POST /create-news

2. Frontend → Gateway
   POST /api/news
   Header: Authorization: Bearer <JWT>

3. Gateway → Auth Service (gRPC)
   ValidateToken(token)
   ← { user_id, role, permissions }

4. Gateway → News Service (gRPC)
   CreateNews(user_id, news_data)
   ← { news_id, slug }

5. Gateway → SEO Service (gRPC)
   GenerateMeta(news_id)
   ← { meta_tags, og_tags, schema }

6. Gateway → Frontend
   { success: true, news_id, slug }

7. Frontend → User
   Redirect to /news/[slug]
```

---

## 🔐 Аутентификация

**JWT Token Flow:**

```
Login → Access Token (15 min) + Refresh Token (7 days)
      ↓
Каждый запрос → Authorization: Bearer <access_token>
      ↓
Gateway → Auth Service: ValidateToken()
      ↓
      ├─ Valid → Proceed to service
      └─ Invalid → 401 Unauthorized

Refresh → New Access Token + New Refresh Token
```

**Роли:**
- `admin` - полный доступ
- `editor` - создание/редактирование новостей
- `moderator` - модерация контента
- `user` - чтение, комментарии

---

## 🗄️ Базы данных

**PostgreSQL (отдельные БД для каждого сервиса):**
- `auth_db` - пользователи, роли
- `news_db` - новости, категории, теги
- `seo_db` - SEO метаданные
- `media_db` - информация о медиафайлах

**Redis:**
- Кеш (популярные новости)
- Сессии (refresh токены)
- Rate limiting

**MinIO (S3-совместимое):**
- Хранение изображений и видео
- CDN интеграция

---

## 📡 Межсервисное взаимодействие

**gRPC (основной протокол):**
- Высокая производительность (HTTP/2 + Protobuf)
- Строгая типизация через .proto файлы
- Bi-directional streaming
- Автогенерация клиентов

**REST/HTTP (Gateway ↔ Frontend):**
- Стандартный HTTP/JSON
- Легко отлаживается
- Кеширование

---

## 🚦 API Gateway функции

1. **Маршрутизация** - проксирование к микросервисам
2. **Аутентификация** - проверка JWT токенов
3. **Авторизация** - проверка ролей и прав
4. **Rate Limiting** - 100 req/min на IP
5. **CORS** - настройка cross-origin запросов
6. **Логирование** - централизованные логи
7. **Трансформация** - gRPC ↔ REST

---

## 🌐 Service Discovery

**Consul:**
- Автоматическая регистрация сервисов
- Health checks
- DNS-based discovery
- Load balancing

**Альтернатива:** статическая конфигурация через env переменные

---

## 📈 Масштабирование

### Горизонтальное масштабирование:
```bash
# Docker Compose
docker-compose up --scale news-service=3

# Kubernetes
kubectl scale deployment news-service --replicas=5
```

### Вертикальное масштабирование:
```yaml
resources:
  limits:
    memory: "1Gi"
    cpu: "1000m"
  requests:
    memory: "512Mi"
    cpu: "500m"
```

### Кеширование:
- **L1:** Redis (популярные новости - 5 мин)
- **L2:** CDN (статика - 24 часа)
- **L3:** Browser cache

---

## 🔒 Безопасность

1. **JWT токены** с коротким TTL
2. **HTTPS/TLS** для всех соединений
3. **Rate limiting** на уровне Gateway
4. **SQL injection защита** через prepared statements
5. **Input validation** на всех входах
6. **CORS политики**
7. **Секреты** через environment variables / Vault

---

## 📊 Мониторинг и наблюдаемость

**Метрики (Prometheus):**
- HTTP request count/latency
- gRPC call statistics
- Database connection pool
- Memory/CPU usage

**Логирование:**
- Структурированные JSON логи
- Централизованное хранение (ELK)
- Корреляция по request_id

**Трейсинг (Jaeger):**
- Distributed tracing
- Визуализация цепочки вызовов
- Bottleneck detection

**Дашборды (Grafana):**
- System overview
- Service-specific metrics
- Alerting

---

## 🚀 Развертывание

### Development:
```bash
docker-compose up -d
```

### Production (Kubernetes):
```bash
kubectl apply -f k8s/
helm install news-portal ./charts/news-portal
```

### CI/CD:
```
GitHub Actions → Build → Test → Push to Registry → Deploy to K8s
```

---

## 📝 Файловая структура (пример Auth Service)

```
auth-service/
├── cmd/
│   └── auth-service/
│       └── main.go              # Entry point
├── internal/
│   ├── config/                  # Configuration
│   ├── handler/                 # gRPC handlers
│   ├── service/                 # Business logic
│   ├── repository/              # Database layer
│   └── models/                  # Data models
├── pkg/                         # Reusable packages
├── proto/                       # Protobuf definitions
├── migrations/                  # SQL migrations
├── .env                         # Environment variables
├── go.mod                       # Go dependencies
├── Dockerfile                   # Docker image
└── README.md
```

---

## 🎯 Ключевые особенности архитектуры

✅ **Независимое развертывание** сервисов  
✅ **Fault isolation** - падение одного сервиса не роняет систему  
✅ **Technology diversity** - можно использовать разные языки/БД  
✅ **Scalability** - масштабируем только нужные компоненты  
✅ **Team autonomy** - команды работают независимо  
✅ **Clear boundaries** - четкое разделение ответственности  

⚠️ **Trade-offs:**
- Сложность инфраструктуры
- Distributed transactions
- Network latency
- Debugging complexity

---

## 📚 Технологический стек

**Backend:**
- Go 1.21+
- gRPC + Protocol Buffers
- PostgreSQL 15
- Redis 7
- MinIO

**Frontend:**
- Next.js 14
- React 18
- TypeScript
- TailwindCSS

**Infrastructure:**
- Docker & Docker Compose
- Kubernetes
- Consul
- RabbitMQ
- Nginx/Traefik

**Monitoring:**
- Prometheus
- Grafana
- Jaeger
- ELK Stack

---

## 🔢 Порты сервисов

| Сервис | Port | Описание |
|--------|------|----------|
| Frontend | 3000 | Next.js app |
| Gateway | 8080 | REST API |
| Auth | 8081 | gRPC |
| News | 8082 | gRPC |
| SEO | 8083 | gRPC |
| Admin | 8084 | gRPC |
| Media | 8085 | gRPC |
| PostgreSQL | 5432 | Database |
| Redis | 6379 | Cache |
| Consul | 8500 | Service Discovery |
| RabbitMQ | 5672/15672 | Message Queue |
| MinIO | 9000/9001 | Object Storage |
| Prometheus | 9090 | Metrics |
| Grafana | 3001 | Dashboards |

---

## 🎓 Дальнейшие улучшения

1. **Event Sourcing** для истории изменений
2. **CQRS** для разделения чтения/записи
3. **GraphQL Gateway** как альтернатива REST
4. **WebSockets** для real-time уведомлений
5. **Elasticsearch** для продвинутого поиска
6. **Service Mesh (Istio)** для advanced networking
7. **API Versioning** для backward compatibility
8. **Multi-tenancy** для поддержки нескольких сайтов

---

## 📞 Быстрые команды

```bash
# Запуск
docker-compose up -d

# Остановка
docker-compose down

# Логи
docker-compose logs -f gateway

# Проверка здоровья
curl http://localhost:8080/health

# Тесты
cd auth-service && go test ./...

# Миграции
make migrate-up

# Генерация proto
make proto
```

---

**Версия:** 1.0  
**Дата:** 2025-10-14  
**Статус:** Production Ready ✅
