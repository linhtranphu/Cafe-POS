#!/bin/bash

echo "=== Verify Code Mới Đã Deploy ==="
echo ""

echo "1️⃣  Kiểm tra backend image:"
docker images cafe-pos-backend:latest --format "Created: {{.CreatedAt}}"

echo ""
echo "2️⃣  Kiểm tra backend container:"
docker inspect cafe-pos-backend --format='Started: {{.State.StartedAt}}'

echo ""
echo "3️⃣  Tìm log đặc trưng của code mới:"
echo "   Tìm 'Using HTML for print bridge'..."
docker logs cafe-pos-backend 2>&1 | grep "Using HTML for print bridge"

if [ $? -eq 0 ]; then
  echo "   ✅ Code mới đã có!"
else
  echo "   ❌ CHƯA CÓ code mới!"
  echo ""
  echo "   Code mới phải có log: 'Using HTML for print bridge'"
  echo "   Nhưng hiện tại vẫn thấy: 'failed to get default bill template'"
  echo ""
  echo "   → Backend vẫn đang chạy code cũ!"
fi

echo ""
echo "4️⃣  Kiểm tra file HTML template có trong container không:"
docker exec cafe-pos-backend ls -la /root/application/services/templates/ 2>/dev/null

if [ $? -eq 0 ]; then
  echo "   ✅ Thư mục templates tồn tại"
  docker exec cafe-pos-backend ls -la /root/application/services/templates/bill_template_optimized.html 2>/dev/null
  if [ $? -eq 0 ]; then
    echo "   ✅ File bill_template_optimized.html tồn tại"
  else
    echo "   ❌ File bill_template_optimized.html KHÔNG tồn tại!"
  fi
else
  echo "   ❌ Thư mục templates KHÔNG tồn tại!"
fi

echo ""
echo "5️⃣  Kiểm tra code trong container:"
docker exec cafe-pos-backend grep -n "Using HTML for print bridge" /root/application/services/print_service.go 2>/dev/null

if [ $? -eq 0 ]; then
  echo "   ✅ Code mới có trong file"
else
  echo "   ❌ Code mới KHÔNG có trong file!"
  echo "   → Image chưa được build với code mới"
fi

echo ""
echo "=== Kết luận ==="
echo ""
echo "Nếu thấy ❌ ở trên, cần:"
echo "1. Build lại image: docker build -t cafe-pos-backend:latest backend/"
echo "2. Stop container: docker stop cafe-pos-backend && docker rm cafe-pos-backend"
echo "3. Run lại: docker run -d --name cafe-pos-backend --network cafe-pos-network -p 8080:8080 --env-file .env.ec2 --restart unless-stopped cafe-pos-backend:latest"
