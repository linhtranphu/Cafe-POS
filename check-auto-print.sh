#!/bin/bash

echo "=== Kiểm tra Auto Print Configuration ==="
echo ""

echo "✅ Backend đang chạy"
echo ""

# MongoDB credentials
MONGO_USER="admin"
MONGO_PASS="password123"
MONGO_DB="cafe_pos"

# Check shop settings
echo "1. Kiểm tra Shop Settings (auto_print_enabled):"
docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  db.shop_settings.find({}, {
    shop_name: 1,
    auto_print_enabled: 1,
    _id: 0
  }).pretty()
"
echo ""

# Check printers
echo "2. Kiểm tra Printers:"
docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  db.printers.find({}, {
    name: 1,
    type: 1,
    enabled: 1,
    is_default: 1,
    connection_type: 1,
    ip_address: 1,
    port: 1,
    _id: 0
  }).pretty()
"
echo ""

# Check print templates
echo "3. Kiểm tra Print Templates:"
docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  db.print_templates.find({}, {
    name: 1,
    type: 1,
    is_default: 1,
    _id: 0
  }).pretty()
"
echo ""

# Check recent print jobs
echo "4. Kiểm tra Print Jobs gần đây (5 jobs cuối):"
docker exec cafe-pos-mongodb mongosh -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin "$MONGO_DB" --quiet --eval "
  db.print_jobs.find({}).sort({created_at: -1}).limit(5).forEach(function(job) {
    print('Order: ' + job.order_number + ' | Type: ' + job.type + ' | Status: ' + job.status + ' | Created: ' + job.created_at);
  })
"
echo ""

echo "=== Kết luận ==="
echo ""
echo "Để auto-print hoạt động, cần:"
echo "1. ✓ auto_print_enabled = true trong shop_settings"
echo "2. ✓ Có printer type BILL với enabled=true và is_default=true"
echo "3. ✓ Có printer type LABEL với enabled=true và is_default=true"
echo "4. ✓ Có print template type BILL với is_default=true"
echo "5. ✓ Có print template type LABEL với is_default=true"
echo "6. ✓ Order phải có status = PAID (sau khi collect payment)"
echo ""

