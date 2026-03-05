#!/bin/bash

echo "=== Thêm LABEL Printer trên EC2 ==="
echo ""

# Lấy thông tin từ BILL printer hiện có
echo "Đang lấy thông tin từ BILL printer..."

# Thay đổi connection string cho EC2
MONGO_URI="mongodb://admin:password123@localhost:27017/cafe_pos?authSource=admin"

# Lấy IP và Port từ BILL printer
BILL_INFO=$(mongosh "$MONGO_URI" --quiet --eval "
  var bill = db.printers.findOne({type: 'BILL'});
  if (bill) {
    print(bill.ip_address + ':' + bill.port);
  } else {
    print('NOT_FOUND');
  }
")

if [ "$BILL_INFO" = "NOT_FOUND" ]; then
  echo "❌ Không tìm thấy BILL printer"
  echo ""
  echo "Vui lòng nhập thông tin thủ công:"
  read -p "IP Address: " PRINTER_IP
  read -p "Port: " PRINTER_PORT
else
  PRINTER_IP=$(echo $BILL_INFO | cut -d':' -f1)
  PRINTER_PORT=$(echo $BILL_INFO | cut -d':' -f2)
  echo "✅ Tìm thấy BILL printer: $PRINTER_IP:$PRINTER_PORT"
fi

echo ""
echo "Đang tạo LABEL printer với cùng IP/Port..."

# Tạo LABEL printer
mongosh "$MONGO_URI" --quiet --eval "
  var now = new Date();
  
  // Kiểm tra xem đã có LABEL printer chưa
  var existing = db.printers.findOne({type: 'LABEL'});
  if (existing) {
    print('⚠️  LABEL printer đã tồn tại: ' + existing.name);
    print('Đang cập nhật...');
    db.printers.updateOne(
      {type: 'LABEL'},
      {\$set: {
        ip_address: '$PRINTER_IP',
        port: $PRINTER_PORT,
        paper_width: 58,
        enabled: true,
        is_default: true,
        updated_at: now
      }}
    );
    print('✅ Đã cập nhật LABEL printer');
  } else {
    // Tạo mới
    db.printers.insertOne({
      name: 'Máy in Label',
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
    print('✅ Đã tạo LABEL printer mới');
  }
"

echo ""
echo "Kiểm tra lại printers..."
mongosh "$MONGO_URI" --quiet --eval "
  db.printers.find({}, {name: 1, type: 1, enabled: 1, is_default: 1, ip_address: 1, port: 1, paper_width: 1, _id: 0}).forEach(function(p) {
    print('  ✓ ' + p.name + ' (Type: ' + p.type + ', IP: ' + p.ip_address + ':' + p.port + ', Paper: ' + p.paper_width + 'mm)');
  });
"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ HOÀN TẤT!"
echo ""
echo "Bây giờ auto-print sẽ hoạt động:"
echo "  1. Tạo order"
echo "  2. Collect payment"
echo "  3. Bill và label tự động in"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
