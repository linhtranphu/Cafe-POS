#!/bin/bash

echo "=== Kiểm tra Database Collections ==="
echo ""

echo "1️⃣  List tất cả databases:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin --quiet --eval '
db.adminCommand("listDatabases").databases.forEach(function(d) {
  print(d.name);
})
'

echo ""
echo "2️⃣  List collections trong cafe_pos_db:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin cafe_pos_db --quiet --eval '
db.getCollectionNames().forEach(function(c) {
  var count = db[c].countDocuments({});
  print(c + ": " + count + " documents");
})
'

echo ""
echo "2b️⃣  List collections trong cafe_pos (không có _db):"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin cafe_pos --quiet --eval '
db.getCollectionNames().forEach(function(c) {
  var count = db[c].countDocuments({});
  print(c + ": " + count + " documents");
})
'

echo ""
echo "3️⃣  Kiểm tra printer_configs chi tiết:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin cafe_pos_db --quiet --eval '
print("Total printers: " + db.printer_configs.countDocuments({}));
print("\nAll printers:");
db.printer_configs.find({}).forEach(function(p) {
  print("\nPrinter: " + p.name);
  print("  Type: " + p.type);
  print("  IP: " + p.ip_address + ":" + p.port);
  print("  Default: " + p.is_default);
  print("  Enabled: " + p.is_enabled);
})
'

echo ""
echo "4️⃣  Kiểm tra shop_settings:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin cafe_pos_db --quiet --eval '
var count = db.shop_settings.countDocuments({});
print("Total shop_settings: " + count);

if (count > 0) {
  var settings = db.shop_settings.findOne({});
  print("\nSettings:");
  print("  Shop Name: " + (settings.shop_name || "N/A"));
  print("  Print Bridge URL: " + (settings.print_bridge_url || "NOT SET"));
  print("  Auto Print: " + (settings.auto_print_enabled !== undefined ? settings.auto_print_enabled : "undefined"));
}
'

echo ""
echo "5️⃣  Kiểm tra orders:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin cafe_pos_db --quiet --eval '
var total = db.orders.countDocuments({});
var paid = db.orders.countDocuments({status: "PAID"});
print("Total orders: " + total);
print("Paid orders: " + paid);

if (total > 0) {
  print("\nRecent order:");
  var order = db.orders.findOne({}, {sort: {created_at: -1}});
  print("  Order: " + order.order_number);
  print("  Status: " + order.status);
}
'

echo ""
echo "=== Kết thúc kiểm tra ==="
