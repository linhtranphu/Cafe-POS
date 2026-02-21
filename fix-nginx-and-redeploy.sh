#!/bin/bash
set -e

echo "🔧 Fix Nginx Config và Redeploy Frontend"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

echo "1️⃣  Test nginx config syntax..."
docker run --rm -v $(pwd)/frontend/nginx.conf:/etc/nginx/conf.d/default.conf:ro nginx:alpine nginx -t

if [ $? -ne 0 ]; then
    echo "❌ Nginx config có lỗi syntax!"
    exit 1
fi

echo "✅ Nginx config OK!"
echo ""

echo "2️⃣  Rebuild frontend image..."
docker build -t linhtranphu/cafe-pos-frontend:latest ./frontend

if [ $? -ne 0 ]; then
    echo "❌ Build thất bại!"
    exit 1
fi

echo "✅ Build thành công!"
echo ""

echo "3️⃣  Restart frontend container..."
docker-compose -f docker-compose.prod.yml up -d frontend

echo ""
echo "4️⃣  Đợi container khởi động (10 giây)..."
sleep 10

echo ""
echo "5️⃣  Kiểm tra container status..."
docker ps | grep frontend

echo ""
echo "6️⃣  Kiểm tra logs..."
docker logs cafe-pos-frontend --tail 20

echo ""
echo "✅ Hoàn tất!"
