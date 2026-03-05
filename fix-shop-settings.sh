#!/bin/bash

echo "=== Fix Shop Settings ==="
echo ""

echo "1️⃣  Kiểm tra shop_settings hiện tại:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin cafe_pos_db --quiet --eval '
var settings = db.shop_settings.findOne({});
if (settings) {
  print("✅ Đã có shop_settings");
  print("   Shop Name: " + (settings.shop_name || "N/A"));
  print("   Print Bridge URL: " + (settings.print_bridge_url || "NOT SET"));
  print("   Auto Print: " + (settings.auto_print_enabled !== undefined ? settings.auto_print_enabled : "undefined"));
} else {
  print("❌ CHƯA CÓ shop_settings");
}
'

echo ""
echo "2️⃣  Tạo hoặc update shop_settings:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin cafe_pos_db --quiet --eval '
var result = db.shop_settings.updateOne(
  {},
  {
    $setOnInsert: {
      shop_name: "Tạ Cafe",
      shop_address: "",
      shop_phone: "",
      created_at: new Date(),
    },
    $set: {
      print_bridge_url: "https://print.tacafe.store",
      auto_print_enabled: true,
      updated_at: new Date()
    }
  },
  { upsert: true }
);

if (result.upsertedCount > 0) {
  print("✅ Đã tạo shop_settings mới");
} else if (result.modifiedCount > 0) {
  print("✅ Đã update shop_settings");
} else {
  print("ℹ️  Shop_settings không thay đổi");
}
'

echo ""
echo "3️⃣  Verify shop_settings:"
docker exec cafe-pos-mongodb mongosh -u admin -p 108trannhatduat --authenticationDatabase admin cafe_pos_db --quiet --eval '
var settings = db.shop_settings.findOne({});
print("Shop Name: " + settings.shop_name);
print("Print Bridge URL: " + settings.print_bridge_url);
print("Auto Print Enabled: " + settings.auto_print_enabled);
'

echo ""
echo "=== Fix hoàn tất ==="
echo ""
echo "Bây giờ thử collect payment lại!"
