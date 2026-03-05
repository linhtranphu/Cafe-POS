#!/bin/bash

echo "=== Debug Collect Payment - Print Job Creation ==="
echo ""

echo "1️⃣  Backend logs khi collect payment (50 dòng cuối):"
docker logs cafe-pos-backend 2>&1 | tail -50

echo ""
echo "2️⃣  Tìm logs liên quan đến print:"
docker logs cafe-pos-backend 2>&1 | grep -i "print\|payment\|order.*paid" | tail -30

echo ""
echo "3️⃣  Kiểm tra auto-print setting:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin cafe_pos_db --quiet --eval '
var settings = db.shop_settings.findOne({});
print("Auto Print Enabled: " + (settings.auto_print_enabled !== undefined ? settings.auto_print_enabled : "undefined (default: true)"));
print("Print Bridge URL: " + (settings.print_bridge_url || "NOT SET"));
'

echo ""
echo "4️⃣  Kiểm tra có printer configs không:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin cafe_pos_db --quiet --eval '
var billPrinter = db.printer_configs.findOne({type: "BILL", is_default: true, is_enabled: true});
if (billPrinter) {
  print("✅ Bill Printer: " + billPrinter.name);
  print("   IP: " + billPrinter.ip_address + ":" + billPrinter.port);
} else {
  print("❌ KHÔNG CÓ default bill printer enabled");
}

var labelPrinter = db.printer_configs.findOne({type: "LABEL", is_default: true, is_enabled: true});
if (labelPrinter) {
  print("✅ Label Printer: " + labelPrinter.name);
} else {
  print("❌ KHÔNG CÓ default label printer enabled");
}
'

echo ""
echo "5️⃣  Kiểm tra orders gần đây:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin cafe_pos_db --quiet --eval '
db.orders.find({}).sort({created_at: -1}).limit(3).forEach(function(o) {
  print("\nOrder: " + o.order_number);
  print("  Status: " + o.status);
  print("  Total: " + o.total);
  print("  Payment Method: " + (o.payment_method || "N/A"));
  print("  Paid At: " + (o.paid_at || "N/A"));
  print("  Created: " + o.created_at);
})
'

echo ""
echo "6️⃣  Kiểm tra print jobs:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin cafe_pos_db --quiet --eval '
var count = db.print_jobs.countDocuments({});
print("Total print jobs: " + count);

if (count > 0) {
  print("\nRecent print jobs:");
  db.print_jobs.find({}).sort({created_at: -1}).limit(3).forEach(function(j) {
    print("  Job: " + j.order_number + " (" + j.type + ")");
    print("    Status: " + j.status);
    print("    Created: " + j.created_at);
  });
} else {
  print("\n❌ KHÔNG CÓ PRINT JOB NÀO!");
}
'

echo ""
echo "=== Kết thúc debug ==="
echo ""
echo "💡 Nếu không thấy log \"Creating print jobs for order\"..."
echo "   → Auto-print có thể bị tắt"
echo "   → Hoặc không có default printer"
echo "   → Hoặc code chưa được deploy"
