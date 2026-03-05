#!/bin/bash

# Script để chẩn đoán vấn đề print jobs trên EC2

echo "=== Chẩn đoán Print Jobs ==="
echo ""

# 1. Kiểm tra Print Bridge URL
echo "1️⃣  Print Bridge URL trong settings:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108tret --eval '
db.shop_settings.findOne({}, {
  shop_name: 1,
  print_bridge_url: 1,
  auto_print_enabled: 1,
  _id: 0
})
' 2>/dev/null || echo "❌ Không thể kết nối MongoDB"

echo ""

# 2. Kiểm tra printer configs
echo "2️⃣  Cấu hình máy in:"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval '
db.printer_configs.find({}, {
  name: 1,
  type: 1,
  ip_address: 1,
  port: 1,
  is_enabled: 1,
  is_default: 1,
  _id: 1
}).forEach(function(p) {
  print("Printer: " + p.name);
  print("  ID: " + p._id);
  print("  Type: " + p.type);
  print("  IP: " + p.ip_address + ":" + p.port);
  print("  Enabled: " + p.is_enabled);
  print("  Default: " + p.is_default);
  print("---");
})
' 2>/dev/null

echo ""

# 3. Kiểm tra print jobs gần đây
echo "3️⃣  Print jobs gần đây (10 jobs cuối):"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval '
db.print_jobs.find({}).sort({created_at: -1}).limit(10).forEach(function(job) {
  print("Job ID: " + job._id);
  print("  Order: " + job.order_number);
  print("  Type: " + job.type);
  print("  Status: " + job.status);
  print("  Printer ID: " + job.printer_id);
  print("  Content Type: " + (job.content_type || "text"));
  print("  Content Size: " + (job.content ? job.content.length : 0) + " bytes");
  print("  Retry: " + job.retry_count + "/" + job.max_retries);
  print("  Created: " + job.created_at);
  print("  Updated: " + job.updated_at);
  if (job.error_message) {
    print("  Error: " + job.error_message);
  }
  print("---");
})
' 2>/dev/null

echo ""

# 4. Đếm print jobs theo status
echo "4️⃣  Thống kê print jobs:"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval '
print("PENDING: " + db.print_jobs.countDocuments({status: "PENDING"}));
print("PRINTING: " + db.print_jobs.countDocuments({status: "PRINTING"}));
print("COMPLETED: " + db.print_jobs.countDocuments({status: "COMPLETED"}));
print("FAILED: " + db.print_jobs.countDocuments({status: "FAILED"}));
print("TOTAL: " + db.print_jobs.countDocuments({}));
' 2>/dev/null

echo ""

# 5. Kiểm tra backend logs
echo "5️⃣  Backend logs (print-related, 30 dòng cuối):"
docker logs cafe-pos-backend 2>&1 | grep -i "print\|bridge" | tail -30

echo ""

# 6. Kiểm tra print worker có chạy không
echo "6️⃣  Print worker status:"
docker logs cafe-pos-backend 2>&1 | grep -i "print worker" | tail -5

echo ""

# 7. Test kết nối tới Print Bridge
echo "7️⃣  Test kết nối Print Bridge:"
echo -n "   https://print.tacafe.store/health: "
curl -s -w "HTTP %{http_code} - %{time_total}s\n" https://print.tacafe.store/health -o /dev/null

echo ""

# 8. Kiểm tra có print job bị stuck không
echo "8️⃣  Print jobs bị stuck (PENDING > 5 phút):"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval '
var fiveMinutesAgo = new Date(Date.now() - 5*60*1000);
db.print_jobs.find({
  status: "PENDING",
  created_at: {$lt: fiveMinutesAgo}
}).forEach(function(job) {
  print("⚠️  Job ID: " + job._id);
  print("   Order: " + job.order_number);
  print("   Created: " + job.created_at);
  print("   Age: " + Math.round((Date.now() - job.created_at.getTime())/60000) + " minutes");
  print("   Retry: " + job.retry_count + "/" + job.max_retries);
  if (job.error_message) {
    print("   Last Error: " + job.error_message);
  }
  print("---");
})
' 2>/dev/null

echo ""
echo "=== Kết thúc chẩn đoán ==="
echo ""

# Gợi ý
echo "💡 Gợi ý:"
echo ""
echo "Nếu có print jobs PENDING:"
echo "  - Kiểm tra print_bridge_url có đúng không"
echo "  - Kiểm tra print bridge có chạy không"
echo "  - Kiểm tra backend logs có lỗi gì không"
echo "  - Xem print worker có đang chạy không"
echo ""
echo "Nếu có print jobs FAILED:"
echo "  - Xem error_message để biết nguyên nhân"
echo "  - Kiểm tra IP và port máy in có đúng không"
echo "  - Test kết nối trực tiếp tới máy in"
echo ""
echo "Để xóa print jobs cũ:"
echo "  docker exec cafe-pos-mongodb mongosh cafe_pos_db --eval 'db.print_jobs.deleteMany({status: \"FAILED\"})'"
echo ""
