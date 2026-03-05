#!/bin/bash

echo "=== Test Full Print Flow ==="
echo ""

echo "📋 Checklist:"
echo ""

# 1. Print Bridge
echo "1️⃣  Print Bridge:"
BRIDGE_STATUS=$(curl -s https://print.tacafe.store/health | grep -o '"status":"ok"')
if [ -n "$BRIDGE_STATUS" ]; then
  echo "   ✅ Print Bridge đang chạy"
else
  echo "   ❌ Print Bridge KHÔNG chạy hoặc không truy cập được"
  echo "   → Kiểm tra: docker logs local-print-bridge"
fi

echo ""

# 2. Backend Print Bridge URL
echo "2️⃣  Backend Print Bridge URL:"
BRIDGE_URL=$(docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval 'db.shop_settings.findOne({}, {print_bridge_url: 1, _id: 0}).print_bridge_url' 2>/dev/null | tr -d '\n' | tr -d ' ')
if [ "$BRIDGE_URL" = "https://print.tacafe.store" ]; then
  echo "   ✅ URL đúng: $BRIDGE_URL"
elif [ -z "$BRIDGE_URL" ]; then
  echo "   ❌ URL CHƯA CẤU HÌNH"
  echo "   → Vào https://tacafe.store/#/print-management để cấu hình"
else
  echo "   ⚠️  URL: $BRIDGE_URL"
  echo "   → Nên là: https://print.tacafe.store"
fi

echo ""

# 3. Auto Print
echo "3️⃣  Auto Print:"
AUTO_PRINT=$(docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval 'var s = db.shop_settings.findOne({}); s && s.auto_print_enabled !== undefined ? s.auto_print_enabled : true' 2>/dev/null)
if [ "$AUTO_PRINT" = "true" ]; then
  echo "   ✅ Auto Print đang BẬT"
else
  echo "   ❌ Auto Print đang TẮT"
  echo "   → Vào Settings để bật Auto Print"
fi

echo ""

# 4. Printer Configs
echo "4️⃣  Printer Configs:"
BILL_PRINTER=$(docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval 'db.printer_configs.findOne({type: "BILL", is_default: true, is_enabled: true})' 2>/dev/null)
if [ -n "$BILL_PRINTER" ]; then
  echo "   ✅ Có Bill Printer mặc định"
  docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval '
  var p = db.printer_configs.findOne({type: "BILL", is_default: true, is_enabled: true});
  if (p) {
    print("      Name: " + p.name);
    print("      IP: " + p.ip_address + ":" + p.port);
  }
  ' 2>/dev/null
else
  echo "   ❌ KHÔNG CÓ Bill Printer mặc định hoặc bị disabled"
  echo "   → Vào #/printer-management để cấu hình"
fi

echo ""

# 5. Print Worker
echo "5️⃣  Print Worker:"
WORKER_RUNNING=$(docker logs cafe-pos-backend 2>&1 | grep "Print worker started" | tail -1)
if [ -n "$WORKER_RUNNING" ]; then
  echo "   ✅ Print Worker đã khởi động"
else
  echo "   ❌ Print Worker CHƯA khởi động"
  echo "   → Kiểm tra backend logs"
fi

echo ""

# 6. Print Jobs
echo "6️⃣  Print Jobs:"
TOTAL_JOBS=$(docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval 'db.print_jobs.countDocuments({})' 2>/dev/null)
PENDING_JOBS=$(docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval 'db.print_jobs.countDocuments({status: "PENDING"})' 2>/dev/null)
FAILED_JOBS=$(docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval 'db.print_jobs.countDocuments({status: "FAILED"})' 2>/dev/null)

echo "   Total: $TOTAL_JOBS"
echo "   Pending: $PENDING_JOBS"
echo "   Failed: $FAILED_JOBS"

if [ "$TOTAL_JOBS" = "0" ]; then
  echo "   ⚠️  Chưa có print job nào"
  echo "   → Thử collect payment một order để tạo print job"
fi

if [ "$PENDING_JOBS" != "0" ]; then
  echo "   ⚠️  Có $PENDING_JOBS jobs đang pending"
  echo "   → Kiểm tra print worker logs"
fi

if [ "$FAILED_JOBS" != "0" ]; then
  echo "   ❌ Có $FAILED_JOBS jobs failed"
  echo "   → Xem error: ./diagnose-print-jobs.sh"
fi

echo ""

# 7. Recent Orders
echo "7️⃣  Recent PAID Orders:"
RECENT_PAID=$(docker exec cafe-pos-mongodb mongosh cafe_pos_db --quiet --eval 'db.orders.countDocuments({status: "PAID", paid_at: {$gte: new Date(Date.now() - 24*60*60*1000)}})' 2>/dev/null)
echo "   Paid trong 24h qua: $RECENT_PAID"

if [ "$RECENT_PAID" = "0" ]; then
  echo "   ⚠️  Chưa có order nào được paid trong 24h qua"
  echo "   → Thử collect payment một order để test"
fi

echo ""
echo "=== Tổng kết ==="
echo ""

# Tính điểm
SCORE=0
[ -n "$BRIDGE_STATUS" ] && SCORE=$((SCORE+1))
[ "$BRIDGE_URL" = "https://print.tacafe.store" ] && SCORE=$((SCORE+1))
[ "$AUTO_PRINT" = "true" ] && SCORE=$((SCORE+1))
[ -n "$BILL_PRINTER" ] && SCORE=$((SCORE+1))
[ -n "$WORKER_RUNNING" ] && SCORE=$((SCORE+1))

if [ $SCORE -eq 5 ]; then
  echo "✅ Tất cả đều OK! ($SCORE/5)"
  echo ""
  echo "Nếu vẫn không in được:"
  echo "  1. Thử collect payment một order"
  echo "  2. Chạy: ./check-print-worker-detail.sh"
  echo "  3. Xem backend logs: docker logs -f cafe-pos-backend"
  echo "  4. Xem print bridge logs: docker logs -f local-print-bridge"
else
  echo "⚠️  Có vấn đề cần fix ($SCORE/5)"
  echo ""
  echo "Hãy fix các mục ❌ ở trên trước"
fi

echo ""
