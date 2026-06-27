#!/bin/bash
set -e

APP_DIR="/opt/navisha-english"
COMPOSE_FILE="docker-compose.prod.yml"

echo "==> Navisha English Deploy Script"
echo "    Target: english.navisha.cloud"
echo ""

# 1. Pull latest code
echo "[1/5] Pulling latest code..."
git pull origin main

# 2. Pastikan .env.production ada
if [ ! -f ".env.production" ]; then
  echo "ERROR: .env.production tidak ditemukan!"
  echo "       Salin dari template: cp .env.production.example .env.production"
  echo "       Lalu isi nilai yang diperlukan (JWT_SECRET, TELEGRAM_BOT_TOKEN, dll)"
  exit 1
fi

# 3. Build ulang image
echo "[2/5] Building Docker images..."
docker compose -f $COMPOSE_FILE build --no-cache

# 4. Restart containers
echo "[3/5] Restarting containers..."
docker compose -f $COMPOSE_FILE up -d --remove-orphans

# 5. Tunggu containers healthy
echo "[4/5] Waiting for containers to be healthy..."
sleep 5
docker compose -f $COMPOSE_FILE ps

# 6. Setup nginx (pertama kali saja)
echo "[5/5] Checking nginx config..."
NGINX_CONF="/etc/nginx/sites-available/english.navisha.cloud"
if [ ! -f "$NGINX_CONF" ]; then
  echo "      Copying nginx config..."
  sudo cp nginx/english.navisha.cloud.conf $NGINX_CONF
  sudo ln -sf $NGINX_CONF /etc/nginx/sites-enabled/english.navisha.cloud
  sudo nginx -t && sudo systemctl reload nginx
  echo ""
  echo "      Jalankan certbot untuk SSL:"
  echo "      sudo certbot --nginx -d english.navisha.cloud"
else
  echo "      Nginx config sudah ada, skip."
  sudo nginx -t && sudo systemctl reload nginx
fi

echo ""
echo "Deploy selesai!"
echo "  Frontend : https://english.navisha.cloud"
echo "  Backend  : https://english.navisha.cloud/api/v1"
echo ""
echo "Untuk melihat logs:"
echo "  docker compose -f $COMPOSE_FILE logs -f"
