# 🚀 Руководство по развертыванию

Пошаговое руководство по развертыванию микросервисного новостного портала.

---

## 📋 Предварительные требования

### Для запуска через Docker (рекомендуется)
- Docker Desktop 4.0+
- Docker Compose 2.0+
- 8 GB RAM минимум
- 20 GB свободного места на диске

### Для локальной разработки
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- Protocol Buffers compiler (protoc)

---

## 🎯 Быстрый старт (Docker)

### Шаг 1: Клонирование проекта

```bash
cd "C:\Users\Grifindor\Desktop\НОВСТНОЙ САЙТ"
```

### Шаг 2: Запуск всех сервисов

```bash
# Сборка и запуск всех контейнеров
docker-compose up -d

# Проверка статуса
docker-compose ps

# Просмотр логов
docker-compose logs -f
```

### Шаг 3: Проверка работоспособности

Откройте в браузере:
- **Frontend:** http://localhost:3000
- **API Gateway:** http://localhost:8080/health
- **Consul UI:** http://localhost:8500

### Шаг 4: Первый пользователь

```bash
# Создание администратора
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123",
    "full_name": "Admin User",
    "role": "admin"
  }'
```

---

## 🔧 Локальная разработка

### 1. Установка зависимостей

```bash
# Установка Go зависимостей
cd auth-service && go mod download && cd ..
cd news-service && go mod download && cd ..
cd gateway && go mod download && cd ..

# Установка Node.js зависимостей
cd frontend && npm install && cd ..
```

### 2. Запуск инфраструктуры

```bash
# Запуск только БД и вспомогательных сервисов
docker-compose up -d postgres redis minio rabbitmq consul
```

### 3. Инициализация баз данных

```bash
# Создание баз данных
docker exec -i news-postgres psql -U postgres < scripts/init-db.sql
```

### 4. Применение миграций

```bash
# Установка golang-migrate (если еще не установлен)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Auth Service миграции
migrate -path auth-service/migrations \
  -database "postgresql://postgres:password@localhost:5432/auth_db?sslmode=disable" up

# News Service миграции
migrate -path news-service/migrations \
  -database "postgresql://postgres:password@localhost:5432/news_db?sslmode=disable" up
```

### 5. Генерация protobuf

```bash
# Установка protoc-gen-go (если еще не установлен)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Генерация proto файлов
cd auth-service && protoc --go_out=. --go-grpc_out=. proto/auth.proto && cd ..
cd news-service && protoc --go_out=. --go-grpc_out=. proto/news.proto && cd ..
```

### 6. Запуск сервисов

Откройте 4 терминала:

**Терминал 1 - Auth Service:**
```bash
cd auth-service
go run cmd/auth-service/main.go
```

**Терминал 2 - News Service:**
```bash
cd news-service
go run cmd/news-service/main.go
```

**Терминал 3 - API Gateway:**
```bash
cd gateway
go run cmd/gateway/main.go
```

**Терминал 4 - Frontend:**
```bash
cd frontend
npm run dev
```

---

## 🌐 Конфигурация для Production

### 1. Переменные окружения

Создайте `.env.production` файлы для каждого сервиса:

**auth-service/.env.production:**
```env
SERVICE_NAME=auth-service
GRPC_PORT=8081

DB_HOST=production-db.example.com
DB_PORT=5432
DB_USER=auth_user
DB_PASSWORD=STRONG_PASSWORD_HERE
DB_NAME=auth_db
DB_SSL_MODE=require

REDIS_ADDR=production-redis.example.com:6379
REDIS_PASSWORD=REDIS_PASSWORD_HERE

JWT_SECRET=CHANGE_THIS_TO_STRONG_SECRET
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

CONSUL_HOST=consul.example.com
CONSUL_PORT=8500

LOG_LEVEL=info
LOG_FORMAT=json
```

### 2. Kubernetes Deployment

**Создайте namespace:**
```bash
kubectl create namespace news-portal
```

**Деплой PostgreSQL:**
```yaml
# k8s/postgres-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: news-portal
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        env:
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: password
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc
```

**Деплой сервисов:**
```bash
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/redis-deployment.yaml
kubectl apply -f k8s/auth-service-deployment.yaml
kubectl apply -f k8s/news-service-deployment.yaml
kubectl apply -f k8s/gateway-deployment.yaml
kubectl apply -f k8s/frontend-deployment.yaml
```

### 3. Масштабирование

```bash
# Масштабирование News Service до 3 реплик
kubectl scale deployment news-service --replicas=3 -n news-portal

# Автоматическое масштабирование
kubectl autoscale deployment news-service \
  --cpu-percent=70 \
  --min=2 \
  --max=10 \
  -n news-portal
```

---

## 🔐 Безопасность

### 1. Секреты

**Создание секретов в Kubernetes:**
```bash
# PostgreSQL секрет
kubectl create secret generic postgres-secret \
  --from-literal=password=YOUR_STRONG_PASSWORD \
  -n news-portal

# JWT секрет
kubectl create secret generic jwt-secret \
  --from-literal=secret=YOUR_JWT_SECRET \
  -n news-portal

# Redis секрет
kubectl create secret generic redis-secret \
  --from-literal=password=YOUR_REDIS_PASSWORD \
  -n news-portal
```

### 2. HTTPS/TLS

**Настройка Ingress с Let's Encrypt:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: news-portal-ingress
  namespace: news-portal
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - api.newsportal.com
    - newsportal.com
    secretName: newsportal-tls
  rules:
  - host: api.newsportal.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: gateway
            port:
              number: 8080
  - host: newsportal.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend
            port:
              number: 3000
```

---

## 📊 Мониторинг

### 1. Prometheus + Grafana

```bash
# Установка через Helm
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install prometheus prometheus-community/kube-prometheus-stack \
  -n monitoring --create-namespace
```

### 2. Логирование (ELK Stack)

```bash
# Установка Elasticsearch
helm install elasticsearch elastic/elasticsearch -n logging --create-namespace

# Установка Kibana
helm install kibana elastic/kibana -n logging

# Установка Filebeat
helm install filebeat elastic/filebeat -n logging
```

### 3. Distributed Tracing (Jaeger)

```bash
# Установка Jaeger
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm install jaeger jaegertracing/jaeger -n tracing --create-namespace
```

---

## 🧪 Тестирование

### Unit тесты

```bash
# Все сервисы
cd auth-service && go test ./... && cd ..
cd news-service && go test ./... && cd ..
cd gateway && go test ./... && cd ..
```

### Integration тесты

```bash
# Запуск с тестовым окружением
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

### Load тесты (k6)

```bash
# Установка k6
choco install k6  # Windows

# Запуск load теста
k6 run tests/load/login.js
```

---

## 🔄 CI/CD Pipeline

### GitHub Actions

**.github/workflows/ci.yml:**
```yaml
name: CI/CD

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: |
          cd auth-service && go test ./...
          cd ../news-service && go test ./...
          cd ../gateway && go test ./...
      
      - name: Build
        run: |
          cd auth-service && go build ./...
          cd ../news-service && go build ./...
          cd ../gateway && go build ./...
  
  build-and-push:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}
      
      - name: Build and push
        run: |
          docker build -t username/auth-service:latest ./auth-service
          docker push username/auth-service:latest
          
          docker build -t username/news-service:latest ./news-service
          docker push username/news-service:latest
          
          docker build -t username/gateway:latest ./gateway
          docker push username/gateway:latest
  
  deploy:
    needs: build-and-push
    runs-on: ubuntu-latest
    
    steps:
      - name: Deploy to Kubernetes
        run: |
          kubectl set image deployment/auth-service \
            auth-service=username/auth-service:latest \
            -n news-portal
          
          kubectl set image deployment/news-service \
            news-service=username/news-service:latest \
            -n news-portal
          
          kubectl set image deployment/gateway \
            gateway=username/gateway:latest \
            -n news-portal
```

---

## 🆘 Troubleshooting

### Проблема: Сервисы не могут подключиться к БД

**Решение:**
```bash
# Проверьте статус PostgreSQL
docker exec news-postgres pg_isready -U postgres

# Проверьте логи
docker-compose logs postgres

# Пересоздайте контейнер
docker-compose restart postgres
```

### Проблема: gRPC ошибки между сервисами

**Решение:**
```bash
# Проверьте сетевое соединение
docker network inspect nouvstnoj-sajt_news-network

# Проверьте, что сервисы запущены
docker-compose ps

# Проверьте порты
netstat -an | findstr "8081 8082 8083"
```

### Проблема: Frontend не может подключиться к Gateway

**Решение:**
```bash
# Проверьте переменную окружения в frontend
echo $NEXT_PUBLIC_API_URL

# Проверьте CORS настройки в gateway
curl -I -X OPTIONS http://localhost:8080/api/news
```

---

## 📝 Чеклист развертывания

- [ ] Docker и Docker Compose установлены
- [ ] Порты 3000, 8080-8085, 5432, 6379, 9000-9001 свободны
- [ ] Клонирован репозиторий
- [ ] Выполнен `docker-compose up -d`
- [ ] Проверен статус всех контейнеров
- [ ] Создан первый пользователь-администратор
- [ ] Настроены секреты для production
- [ ] Настроен мониторинг
- [ ] Настроен backup баз данных
- [ ] Настроен CI/CD pipeline
- [ ] Проведены load тесты

---

## 📚 Дополнительные ресурсы

- [ARCHITECTURE.md](../ARCHITECTURE.md) - Подробная архитектура
- [API.md](./API.md) - Документация API
- [Docker Compose Reference](https://docs.docker.com/compose/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/)

---

**Дата обновления:** 2025-10-14
