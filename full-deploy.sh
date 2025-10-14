#!/bin/bash

# Full News Portal Deployment Script
# Deploys ALL services: postgres, redis, auth, news, seo, media, admin, frontend, nginx

set -e

SERVER_IP="151.241.228.203"

echo "==================================================="
echo "  News Portal - FULL Deployment"
echo "  All services: DB, Redis, Auth, News, SEO,"
echo "  Media, Admin Panel, Frontend, Gateway"
echo "==================================================="

cd /root/newsportal/newsportal

echo ""
echo "Step 1: Creating environment files..."

# Admin Service .env
cat > admin-service/.env << 'EOF'
DB_HOST=postgres
DB_PORT=5432
DB_USER=newsportal
DB_PASSWORD=newsportal123
DB_NAME=newsportal_db
JWT_SECRET=super-secret-jwt-key-change-in-production-2024
NEWS_SERVICE_URL=http://news-service:8082
SEO_SERVICE_URL=http://seo-service:8093
MEDIA_SERVICE_URL=http://media-service:8085
PORT=8084
EOF

echo "✅ admin-service/.env created"

# Frontend .env.production
cat > frontend/.env.production << EOF
NEXT_PUBLIC_API_URL=http://${SERVER_IP}:8080
EOF

echo "✅ frontend/.env.production created"

# Set BASE_URL environment variable for docker-compose
export BASE_URL="http://${SERVER_IP}"
export SERVER_IP="${SERVER_IP}"

echo ""
echo "Step 2: Stopping old containers (if any)..."
docker compose -f docker-compose-minimal.yml down 2>/dev/null || true

echo ""
echo "Step 3: Building Docker images..."
echo "⏳ This will take 10-15 minutes on first run..."
docker compose -f docker-compose-minimal.yml build

echo ""
echo "Step 4: Starting all services..."
docker compose -f docker-compose-minimal.yml up -d

echo ""
echo "Step 5: Waiting for services to start (60 seconds)..."
sleep 60

echo ""
echo "Step 6: Health checks..."
echo ""

services=(
    "nginx-gateway:8080:/health"
    "auth-service:8081:/health"
    "news-service:8082:/health"
    "admin-service:8084:/health"
    "media-service:8085:/health"
    "seo-service:8093:/health"
    "frontend:3000"
)

for service_info in "${services[@]}"; do
    IFS=':' read -r name port path <<< "$service_info"
    
    if [ -z "$path" ]; then
        url="http://localhost:${port}"
    else
        url="http://localhost:${port}${path}"
    fi
    
    echo -n "Checking ${name}... "
    if curl -s -f "${url}" > /dev/null 2>&1; then
        echo "✅ OK"
    else
        echo "❌ FAILED"
    fi
done

echo ""
echo "Step 7: Running containers:"
docker compose -f docker-compose-minimal.yml ps

echo ""
echo "==================================================="
echo "  ✅ DEPLOYMENT COMPLETE!"
echo "==================================================="
echo ""
echo "🌐 PUBLIC WEBSITE:  http://${SERVER_IP}:3000"
echo "🔐 ADMIN PANEL:     http://${SERVER_IP}:8084"
echo "🚪 API GATEWAY:     http://${SERVER_IP}:8080"
echo ""
echo "📊 Services:"
echo "  - Auth Service:   http://${SERVER_IP}:8081"
echo "  - News Service:   http://${SERVER_IP}:8082"
echo "  - SEO Service:    http://${SERVER_IP}:8093"
echo "  - Media Service:  http://${SERVER_IP}:8085"
echo "  - Admin Service:  http://${SERVER_IP}:8084"
echo ""
echo "🔑 Default admin credentials:"
echo "  Username: admin"
echo "  Password: admin123"
echo ""
echo "📝 View logs: docker compose -f docker-compose-minimal.yml logs -f"
echo "🔄 Restart: docker compose -f docker-compose-minimal.yml restart"
echo "⏹️  Stop: docker compose -f docker-compose-minimal.yml down"
echo "==================================================="
