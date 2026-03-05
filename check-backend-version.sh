#!/bin/bash

echo "=== Kiểm tra Backend Version ==="
echo ""

echo "1️⃣  Backend container info:"
docker ps --filter name=cafe-pos-backend --format "table {{.Names}}\t{{.Status}}\t{{.Image}}"

echo ""
echo "2️⃣  Kiểm tra code có logic mới không:"
echo "   Tìm log 'Using HTML for print bridge'..."
docker logs cafe-pos-backend 2>&1 | grep "Using HTML for print bridge" | tail -5

if [ $? -eq 0 ]; then
  echo "   ✅ Code mới đã được deploy!"
else
  echo "   ❌ CHƯA CÓ code mới!"
  echo ""
  echo "   Bạn cần deploy code mới:"
  echo "   1. Build: docker build -t cafe-pos-backend:latest backend/"
  echo "   2. Stop: docker stop cafe-pos-backend && docker rm cafe-pos-backend"
  echo "   3. Run: docker run -d --name cafe-pos-backend --network cafe-pos-network -p 8080:8080 --env-file .env.ec2 --restart unless-stopped cafe-pos-backend:latest"
fi

echo ""
echo "3️⃣  Backend start time:"
docker inspect cafe-pos-backend --format='Started: {{.State.StartedAt}}'

echo ""
echo "4️⃣  Image created time:"
docker images cafe-pos-backend:latest --format "Created: {{.CreatedAt}}"

echo ""
echo "=== Kết thúc kiểm tra ==="
