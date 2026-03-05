#!/bin/bash

echo "=== Kiểm tra và Sửa Auto Print ==="
echo ""

MONGO_USER="admin"
MONGO_PASS="password123"
MONGO_DB="cafe_pos"

# 1. Kiểm tra auto_print_enabled
echo "1. Kiểm tra auto_print_enabled trong shop_settings:"
AUTO_PRINT=$(docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  var settings = db.shop_settings.findOne({});
  if (settings && settings.auto_print_enabled !== undefined) {
    print(settings.auto_print_enabled);
  } else {
    print('NOT_SET');
  }
")

if [ "$AUTO_PRINT" = "true" ]; then
  echo "   ✅ auto_print_enabled = true"
elif [ "$AUTO_PRINT" = "false" ]; then
  echo "   ❌ auto_print_enabled = false (BỊ TẮT!)"
  echo ""
  read -p "Bạn có muốn BẬT auto_print_enabled? (y/n): " -n 1 -r
  echo ""
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
      db.shop_settings.updateMany({}, {\$set: {auto_print_enabled: true}});
      print('✅ Đã BẬT auto_print_enabled');
    "
  fi
else
  echo "   ⚠️  auto_print_enabled chưa được set, đang thêm..."
  docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
    db.shop_settings.updateMany({}, {\$set: {auto_print_enabled: true}});
    print('✅ Đã thêm auto_print_enabled = true');
  "
fi
echo ""

# 2. Kiểm tra printers
echo "2. Kiểm tra Printers:"
PRINTER_COUNT=$(docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  print(db.printers.countDocuments({}));
")

if [ "$PRINTER_COUNT" = "0" ]; then
  echo "   ❌ KHÔNG CÓ PRINTER NÀO! (Đây là vấn đề chính)"
  echo ""
  echo "   Bạn cần thêm printer để auto-print hoạt động."
  echo "   Vui lòng cung cấp thông tin printer:"
  echo ""
  
  read -p "   Tên printer (ví dụ: Máy in Bill): " PRINTER_NAME
  read -p "   IP Address (ví dụ: 192.168.1.115): " PRINTER_IP
  read -p "   Port (mặc định 9100): " PRINTER_PORT
  PRINTER_PORT=${PRINTER_PORT:-9100}
  
  echo ""
  echo "   Đang tạo BILL printer và LABEL printer..."
  
  docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
    var now = new Date();
    
    // Tạo BILL printer
    db.printers.insertOne({
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
    db.printers.insertOne({
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
    
    print('✅ Đã tạo 2 printers (BILL và LABEL)');
  "
else
  echo "   ✅ Có $PRINTER_COUNT printer(s)"
  docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
    db.printers.find({}, {name: 1, type: 1, enabled: 1, is_default: 1, ip_address: 1, port: 1, _id: 0}).forEach(function(p) {
      print('      - ' + p.name + ' (Type: ' + p.type + ', Default: ' + p.is_default + ', Enabled: ' + p.enabled + ', IP: ' + p.ip_address + ':' + p.port + ')');
    });
  "
fi
echo ""

# 3. Kiểm tra print templates
echo "3. Kiểm tra Print Templates:"
BILL_TEMPLATE=$(docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  var t = db.print_templates.findOne({type: 'BILL', is_default: true});
  print(t ? 'OK' : 'MISSING');
")

LABEL_TEMPLATE=$(docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  var t = db.print_templates.findOne({type: 'LABEL', is_default: true});
  print(t ? 'OK' : 'MISSING');
")

if [ "$BILL_TEMPLATE" = "OK" ]; then
  echo "   ✅ BILL template OK"
else
  echo "   ❌ BILL template MISSING"
fi

if [ "$LABEL_TEMPLATE" = "OK" ]; then
  echo "   ✅ LABEL template OK"
else
  echo "   ❌ LABEL template MISSING"
fi
echo ""

# 4. Tổng kết
echo "=== Tổng kết ==="
echo ""

ALL_OK=true

if [ "$AUTO_PRINT" != "true" ]; then
  echo "❌ auto_print_enabled không được bật"
  ALL_OK=false
fi

if [ "$PRINTER_COUNT" = "0" ]; then
  echo "❌ Không có printer nào"
  ALL_OK=false
fi

if [ "$BILL_TEMPLATE" != "OK" ]; then
  echo "❌ Thiếu BILL template"
  ALL_OK=false
fi

if [ "$LABEL_TEMPLATE" != "OK" ]; then
  echo "❌ Thiếu LABEL template"
  ALL_OK=false
fi

if [ "$ALL_OK" = true ]; then
  echo "✅ TẤT CẢ ĐỀU OK! Auto-print sẽ hoạt động khi:"
  echo "   1. Tạo order"
  echo "   2. Collect payment (order status = PAID)"
  echo "   3. Bill và label sẽ tự động in"
else
  echo "⚠️  Vẫn còn vấn đề cần sửa (xem bên trên)"
fi
echo ""

echo "Lưu ý: Order phải có status = PAID thì mới tự động in!"
echo "Nghĩa là phải COLLECT PAYMENT trước, không phải chỉ tạo order."
