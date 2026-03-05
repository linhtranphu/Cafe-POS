#!/bin/bash

echo "=== Test Auto Print Flow ==="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

MONGO_USER="admin"
MONGO_PASS="password123"
MONGO_DB="cafe_pos"

echo "Bước 1: Kiểm tra cấu hình"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check auto_print_enabled
AUTO_PRINT=$(docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  var settings = db.shop_settings.findOne({});
  print(settings && settings.auto_print_enabled ? 'true' : 'false');
")

if [ "$AUTO_PRINT" = "true" ]; then
  echo -e "${GREEN}✓${NC} auto_print_enabled = true"
else
  echo -e "${RED}✗${NC} auto_print_enabled = false hoặc không tồn tại"
fi

# Check printers
BILL_PRINTER=$(docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  var p = db.printers.findOne({type: 'BILL', is_default: true, enabled: true});
  print(p ? 'OK' : 'MISSING');
")

LABEL_PRINTER=$(docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  var p = db.printers.findOne({type: 'LABEL', is_default: true, enabled: true});
  print(p ? 'OK' : 'MISSING');
")

if [ "$BILL_PRINTER" = "OK" ]; then
  echo -e "${GREEN}✓${NC} BILL printer configured"
else
  echo -e "${RED}✗${NC} BILL printer MISSING"
fi

if [ "$LABEL_PRINTER" = "OK" ]; then
  echo -e "${GREEN}✓${NC} LABEL printer configured"
else
  echo -e "${RED}✗${NC} LABEL printer MISSING"
fi

# Check templates
BILL_TEMPLATE=$(docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  var t = db.print_templates.findOne({type: 'BILL', is_default: true});
  print(t ? 'OK' : 'MISSING');
")

LABEL_TEMPLATE=$(docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  var t = db.print_templates.findOne({type: 'LABEL', is_default: true});
  print(t ? 'OK' : 'MISSING');
")

if [ "$BILL_TEMPLATE" = "OK" ]; then
  echo -e "${GREEN}✓${NC} BILL template configured"
else
  echo -e "${RED}✗${NC} BILL template MISSING"
fi

if [ "$LABEL_TEMPLATE" = "OK" ]; then
  echo -e "${GREEN}✓${NC} LABEL template configured"
else
  echo -e "${RED}✗${NC} LABEL template MISSING"
fi

echo ""
echo "Bước 2: Kiểm tra Backend API"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check if backend is running
if pgrep -f "backend" > /dev/null; then
  echo -e "${GREEN}✓${NC} Backend đang chạy"
  
  # Check backend port
  if curl -s http://localhost:3000/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Backend API responding (port 3000)"
  else
    echo -e "${YELLOW}⚠${NC} Backend API không response trên port 3000"
  fi
else
  echo -e "${RED}✗${NC} Backend KHÔNG chạy"
fi

echo ""
echo "Bước 3: Kiểm tra Print Jobs gần đây"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  var jobs = db.print_jobs.find({}).sort({created_at: -1}).limit(5).toArray();
  if (jobs.length === 0) {
    print('Không có print job nào');
  } else {
    print('5 print jobs gần nhất:');
    jobs.forEach(function(job) {
      var status = job.status;
      var icon = status === 'COMPLETED' ? '✓' : (status === 'FAILED' ? '✗' : '○');
      print('  ' + icon + ' Order: ' + job.order_number + ' | Type: ' + job.type + ' | Status: ' + status);
      if (job.error_message) {
        print('    Error: ' + job.error_message);
      }
    });
  }
"

echo ""
echo "Bước 4: Kiểm tra Backend Logs"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Tìm kiếm log liên quan đến print..."
echo ""

# Check if backend log file exists
if [ -f "backend.log" ]; then
  echo "Logs về auto-print (10 dòng cuối):"
  grep -i "print\|auto" backend.log | tail -10
elif docker ps | grep -q cafe-pos-backend; then
  echo "Logs từ Docker container:"
  docker logs cafe-pos-backend 2>&1 | grep -i "print\|auto" | tail -10
else
  echo -e "${YELLOW}⚠${NC} Không tìm thấy backend logs"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Tóm tắt:"
echo ""

CAN_AUTO_PRINT=true

if [ "$AUTO_PRINT" != "true" ]; then
  echo -e "${RED}✗${NC} auto_print_enabled không được bật"
  CAN_AUTO_PRINT=false
fi

if [ "$BILL_PRINTER" != "OK" ]; then
  echo -e "${RED}✗${NC} Thiếu BILL printer (ĐÂY LÀ VẤN ĐỀ CHÍNH!)"
  CAN_AUTO_PRINT=false
fi

if [ "$LABEL_PRINTER" != "OK" ]; then
  echo -e "${RED}✗${NC} Thiếu LABEL printer"
  CAN_AUTO_PRINT=false
fi

if [ "$BILL_TEMPLATE" != "OK" ]; then
  echo -e "${RED}✗${NC} Thiếu BILL template"
  CAN_AUTO_PRINT=false
fi

if [ "$LABEL_TEMPLATE" != "OK" ]; then
  echo -e "${RED}✗${NC} Thiếu LABEL template"
  CAN_AUTO_PRINT=false
fi

if [ "$CAN_AUTO_PRINT" = true ]; then
  echo -e "${GREEN}✓ TẤT CẢ CẤU HÌNH OK!${NC}"
  echo ""
  echo "Auto-print sẽ hoạt động khi:"
  echo "  1. Tạo order"
  echo "  2. Collect payment (order status → PAID)"
  echo "  3. Backend tự động tạo print jobs"
  echo "  4. Print worker xử lý và gửi đến máy in"
  echo ""
  echo "API endpoint: POST /api/waiter/orders/:id/payment"
else
  echo -e "${RED}✗ CÓ VẤN ĐỀ CẦN SỬA!${NC}"
  echo ""
  echo "Để sửa, chạy: ./fix-auto-print.sh"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
