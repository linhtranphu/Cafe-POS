#!/bin/bash

echo "=== Fix Database Name ==="
echo ""

echo "1️⃣  Kiểm tra .env.ec2 hiện tại:"
grep MONGODB_DATABASE .env.ec2 || echo "MONGODB_DATABASE không tồn tại"

echo ""
echo "2️⃣  Backup .env.ec2:"
cp .env.ec2 .env.ec2.backup
echo "✅ Đã backup sang .env.ec2.backup"

echo ""
echo "3️⃣  Sửa MONGODB_DATABASE:"
if grep -q "MONGODB_DATABASE=" .env.ec2; then
  sed -i 's/MONGODB_DATABASE=.*/MONGODB_DATABASE=cafe_pos/' .env.ec2
  echo "✅ Đã sửa MONGODB_DATABASE=cafe_pos"
else
  echo "MONGODB_DATABASE=cafe_pos" >> .env.ec2
  echo "✅ Đã thêm MONGODB_DATABASE=cafe_pos"
fi

echo ""
echo "4️⃣  Verify .env.ec2:"
grep MONGODB .env.ec2

echo ""
echo "5️⃣  Restart backend để apply changes:"
docker stop cafe-pos-backend
docker rm cafe-pos-backend
docker run -d \
  --name cafe-pos-backend \
  --network cafe-pos-network \
  -p 8080:8080 \
  --env-file .env.ec2 \
  --restart unless-stopped \
  cafe-pos-backend:latest

echo ""
echo "6️⃣  Đợi backend khởi động..."
sleep 5

echo ""
echo "7️⃣  Kiểm tra backend logs:"
docker logs cafe-pos-backend 2>&1 | tail -20

echo ""
echo "=== Fix hoàn tất ==="
echo ""
echo "Bây giờ thử collect payment lại!"
