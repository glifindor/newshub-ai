#!/bin/bash
set -e

SERVER_IP="151.241.228.203"

echo "==================================================="
echo "  News Portal - Final Deployment"
echo "==================================================="

cd /root/newsportal/newsportal

echo "Step 1: Creating .env files..."

cat > admin-service/.env << 'EOF'
DB_HOST=postgres
DB_PORT=5432
DB_USER=newsportal
DB_PASSWORD=newsportal123
DB_NAME=newsportal_db
JWT_SECRET=super-secret-jwt-key-2024
NEWS_SERVICE_URL=http://news-service:8082
SEO_SERVICE_URL=http://seo-service:8093
MEDIA_SERVICE_URL=http://media-service:8085
PORT=8084
EOF

cat > frontend/.env.production << EOF
NEXT_PUBLIC_API_URL=http://${SERVER_IP}:8080
EOF

echo "✅ Environment files created"

echo ""
echo "Step 2: Stopping old containers..."
docker compose -f docker-compose-simple.yml down 2>/dev/null || true

echo ""
echo "Step 3: Building images (this will take 5-10 minutes)..."
docker compose -f docker-compose-simple.yml build --no-cache

echo ""
echo "Step 4: Starting services..."
docker compose -f docker-compose-simple.yml up -d

echo ""
echo "Step 5: Waiting for services (30 seconds)..."
sleep 30

echo ""
echo "Step 6: Checking containers..."
docker compose -f docker-compose-simple.yml ps

echo ""
echo "==================================================="
echo "  Deployment Complete!"
echo "==================================================="
echo ""
echo "Public Website:  http://${SERVER_IP}:3000"
echo "Admin Panel:     http://${SERVER_IP}:8084"
echo "API Gateway:     http://${SERVER_IP}:8080"
echo ""
echo "Services:"
echo "  - Auth:   http://${SERVER_IP}:8081"
echo "  - News:   http://${SERVER_IP}:8082"  
echo "  - SEO:    http://${SERVER_IP}:8093"
echo "  - Media:  http://${SERVER_IP}:8085"
echo ""
echo "Check logs: docker compose -f docker-compose-simple.yml logs -f"
echo "==================================================="
