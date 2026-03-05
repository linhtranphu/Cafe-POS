#!/bin/bash

echo "=== Kiểm tra Print Worker chi tiết ==="
echo ""

echo "1️⃣  Tìm tất cả logs của Print Worker:"
docker logs cafe-pos-backend 2>&1 | grep -i "print worker\|processing.*print\|print job" | tail -50

echo ""
echo "2️⃣  Kiểm tra có print jobs nào không:"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval '
var total = db.print_jobs.countDocuments({});
var pending = db.print_jobs.countDocuments({status: "PENDING"});
var printing = db.print_jobs.countDocuments({status: "PRINTING"});
var completed = db.print_jobs.countDocuments({status: "COMPLETED"});
var failed = db.print_jobs.countDocuments({status: "FAILED"});

print("📊 Print Jobs Summary:");
print("   Total: " + total);
print("   PENDING: " + pending);
print("   PRINTING: " + printing);
print("   COMPLETED: " + completed);
print("   FAILED: " + failed);

if (total === 0) {
  print("\n⚠️  KHÔNG CÓ PRINT JOB NÀO!");
  print("   → Có thể auto-print bị tắt hoặc chưa có order nào được collect payment");
}
'

echo ""
echo "3️⃣  Kiểm tra auto-print có bật không:"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval '
var settings = db.shop_settings.findOne({}, {auto_print_enabled: 1, _id: 0});
if (settings && settings.auto_print_enabled !== undefined) {
  print("Auto Print: " + (settings.auto_print_enabled ? "✅ BẬT" : "❌ TẮT"));
} else {
  print("Auto Print: ⚠️  CHƯA CẤU HÌNH (mặc định: BẬT)");
}
'

echo ""
echo "4️⃣  Xem 3 print jobs gần nhất (nếu có):"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval '
db.print_jobs.find({}).sort({created_at: -1}).limit(3).forEach(function(job) {
  print("\n📄 Job: " + job.order_number + " (" + job.type + ")");
  print("   Status: " + job.status);
  print("   Created: " + job.created_at);
  print("   Printer ID: " + job.printer_id);
  if (job.error_message) {
    print("   ❌ Error: " + job.error_message);
  }
})
'

echo ""
echo "5️⃣  Kiểm tra có order nào được PAID gần đây không:"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval '
var recentPaidOrders = db.orders.find({
  status: "PAID",
  paid_at: {$exists: true}
}).sort({paid_at: -1}).limit(5);

var count = 0;
recentPaidOrders.forEach(function(order) {
  count++;
  print("\n💰 Order: " + order.order_number);
  print("   Status: " + order.status);
  print("   Paid At: " + order.paid_at);
  print("   Total: " + order.total);
});

if (count === 0) {
  print("⚠️  KHÔNG CÓ ORDER NÀO ĐƯỢC PAID GÀN ĐÂY");
  print("   → Hãy thử collect payment một order để test");
}
'

echo ""
echo "6️⃣  Kiểm tra Print Bridge Client có được khởi tạo không:"
docker logs cafe-pos-backend 2>&1 | grep -i "print bridge.*configur\|bridge.*client\|bridge.*url" | head -10

echo ""
echo "=== Kết thúc kiểm tra ==="
