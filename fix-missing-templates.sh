#!/bin/bash

echo "=== Fix: Thiếu Print Templates ==="
echo ""

echo "1️⃣  Kiểm tra templates hiện có:"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --eval '
db.print_templates.find({}, {
  name: 1,
  type: 1,
  is_default: 1,
  _id: 0
}).forEach(function(t) {
  print(t.type + ": " + t.name + " (default: " + t.is_default + ")");
})
'

echo ""
echo "2️⃣  Tạo default bill template nếu chưa có:"

docker exec cafe-pos-mongodb mongosh cafe_pos_db --eval '
// Kiểm tra có bill template default không
var existingBill = db.print_templates.findOne({type: "BILL", is_default: true});

if (!existingBill) {
  print("⚠️  Không có default bill template, đang tạo...");
  
  var billTemplate = {
    name: "Default Bill Template",
    type: "BILL",
    content: `================================
{{.ShopName}}
{{if .ShowAddress}}{{.ShopAddress}}{{end}}
{{if .ShowPhone}}Tel: {{.ShopPhone}}{{end}}
================================

Order: {{.OrderNumber}}
Waiter: {{.WaiterName}}
Date: {{.CreatedDate}}

--------------------------------
{{range .Items}}
{{.Name}}
  {{.Quantity}} x {{.UnitPrice}} = {{.Total}}
{{end}}
--------------------------------

TOTAL: {{.Total}} VND

Payment: {{.PaymentMethod}}

{{if .ShowCustomMsg}}
{{.CustomMessage}}
{{end}}

Thank you!
================================`,
    is_default: true,
    paper_width: 80,
    created_at: new Date(),
    updated_at: new Date()
  };
  
  db.print_templates.insertOne(billTemplate);
  print("✅ Đã tạo default bill template");
} else {
  print("✅ Đã có default bill template: " + existingBill.name);
}
'

echo ""
echo "3️⃣  Tạo default label template nếu chưa có:"

docker exec cafe-pos-mongodb mongosh cafe_pos_db --eval '
// Kiểm tra có label template default không
var existingLabel = db.print_templates.findOne({type: "LABEL", is_default: true});

if (!existingLabel) {
  print("⚠️  Không có default label template, đang tạo...");
  
  var labelTemplate = {
    name: "Default Label Template",
    type: "LABEL",
    content: `================================
Order: {{.OrderNumber}}
--------------------------------
{{.ItemName}}
{{if .VariantName}}({{.VariantName}}){{end}}

Quantity: {{.Quantity}}
--------------------------------
Waiter: {{.WaiterName}}
Time: {{.CreatedTime}}
================================`,
    is_default: true,
    paper_width: 80,
    created_at: new Date(),
    updated_at: new Date()
  };
  
  db.print_templates.insertOne(labelTemplate);
  print("✅ Đã tạo default label template");
} else {
  print("✅ Đã có default label template: " + existingLabel.name);
}
'

echo ""
echo "4️⃣  Verify templates:"
docker exec cafe-pos-mongodb mongosh cafe_pos_db --eval '
var billCount = db.print_templates.countDocuments({type: "BILL", is_default: true});
var labelCount = db.print_templates.countDocuments({type: "LABEL", is_default: true});

print("Bill templates (default): " + billCount);
print("Label templates (default): " + labelCount);

if (billCount > 0 && labelCount > 0) {
  print("\n✅ Templates OK!");
} else {
  print("\n❌ Vẫn thiếu templates!");
}
'

echo ""
echo "=== Fix hoàn tất ==="
echo ""
echo "Bây giờ thử collect payment lại để test!"
