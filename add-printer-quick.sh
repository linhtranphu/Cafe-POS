#!/bin/bash

echo "=== Thêm Printer cho Auto Print ==="
echo ""

MONGO_USER="admin"
MONGO_PASS="password123"
MONGO_DB="cafe_pos"

# Prompt for printer info
echo "Nhập thông tin máy in:"
read -p "Tên máy in (ví dụ: Máy in Bill): " PRINTER_NAME
PRINTER_NAME=${PRINTER_NAME:-"Máy in Bill"}

read -p "IP Address (ví dụ: 192.168.1.115): " PRINTER_IP
PRINTER_IP=${PRINTER_IP:-"192.168.1.115"}

read -p "Port (mặc định 9100): " PRINTER_PORT
PRINTER_PORT=${PRINTER_PORT:-9100}

echo ""
echo "Đang tạo printers..."
echo "  - BILL printer: $PRINTER_NAME"
echo "  - LABEL printer: $PRINTER_NAME (Label)"
echo "  - IP: $PRINTER_IP:$PRINTER_PORT"
echo ""

# Create printers
docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  var now = new Date();
  
  // Xóa printers cũ nếu có
  db.printers.deleteMany({});
  
  // Tạo BILL printer
  var billResult = db.printers.insertOne({
    name: '$PRINTER_NAME',
    type: 'BILL',
    connection_type: 'NETWORK',
    ip_address: '$PRINTER_IP',
    port: $PRINTER_PORT,
    paper_width: 72,
    enabled: true,
    is_default: true,
    created_at: now,
    updated_at: now
  });
  
  // Tạo LABEL printer (cùng máy in)
  var labelResult = db.printers.insertOne({
    name: '$PRINTER_NAME (Label)',
    type: 'LABEL',
    connection_type: 'NETWORK',
    ip_address: '$PRINTER_IP',
    port: $PRINTER_PORT,
    paper_width: 58,
    enabled: true,
    is_default: true,
    created_at: now,
    updated_at: now
  });
  
  if (billResult.acknowledged && labelResult.acknowledged) {
    print('✅ Đã tạo 2 printers thành công!');
  } else {
    print('❌ Có lỗi khi tạo printers');
  }
"

echo ""
echo "Kiểm tra lại printers..."
docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  db.printers.find({}, {name: 1, type: 1, enabled: 1, is_default: 1, ip_address: 1, port: 1, _id: 0}).forEach(function(p) {
    print('  ✓ ' + p.name + ' (Type: ' + p.type + ', IP: ' + p.ip_address + ':' + p.port + ')');
  });
"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ HOÀN TẤT!"
echo ""
echo "Auto-print bây giờ sẽ hoạt động khi:"
echo "  1. Tạo order"
echo "  2. Collect payment (tiền mặt/chuyển khoản)"
echo "  3. Order status → PAID"
echo "  4. Backend tự động tạo print jobs"
echo "  5. Bill và label tự động in ra"
echo ""
echo "Lưu ý: Đảm bảo máy in đang bật và kết nối mạng!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
