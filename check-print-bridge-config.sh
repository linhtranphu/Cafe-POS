#!/bin/bash

# Script để kiểm tra cấu hình Print Bridge URL trong database

echo "=== Kiểm tra cấu hình Print Bridge ==="
echo ""

# Kết nối MongoDB và kiểm tra settings
echo "1. Kiểm tra PrintBridgeURL trong database:"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval '
db.shop_settings.find({}, {
  shop_name: 1,
  print_bridge_url: 1,
  _id: 0
}).pretty()
'

echo ""
echo "2. Kiểm tra các printer configs:"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval '
db.printer_configs.find({}, {
  name: 1,
  type: 1,
  ip_address: 1,
  port: 1,
  is_enabled: 1,
  is_default: 1,
  _id: 0
}).pretty()
'

echo ""
echo "3. Kiểm tra print jobs gần đây:"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval '
db.print_jobs.find({}).sort({created_at: -1}).limit(5).forEach(function(job) {
  print("Job ID: " + job._id);
  print("  Order: " + job.order_number);
  print("  Type: " + job.type);
  print("  Status: " + job.status);
  print("  Printer ID: " + job.printer_id);
  print("  Created: " + job.created_at);
  print("  Error: " + (job.error_message || "none"));
  print("---");
})
'

echo ""
echo "4. Test kết nối tới Print Bridge:"
echo "   Từ EC2 → print.tacafe.store"
curl -s -w "\n  HTTP Status: %{http_code}\n  Time: %{time_total}s\n" \
  https://print.tacafe.store/health || echo "  ❌ Không thể kết nối"

echo ""
echo "=== Kết thúc kiểm tra ==="
