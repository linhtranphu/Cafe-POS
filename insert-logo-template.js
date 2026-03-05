// MongoDB script to insert logo template
// Run: mongosh cafe_pos insert-logo-template.js

const template = {
  name: "Bill với Logo",
  type: "BILL",
  description: "Template hóa đơn với logo và bảng món",
  is_default: true, // Set as default
  content: `{{if .ShowLogo}}[LOGO]

{{end}}{{.ShopName}}
{{if .ShowAddress}}{{.ShopAddress}}
{{end}}{{if .ShowPhone}}Hotline: {{.ShopPhone}}
{{end}}
HÓA ĐƠN THANH TOÁN
Số HĐ: {{.Order.OrderNumber}}

Mã HĐ: #{{.Order.OrderNumber}}     TN: {{if .Order.WaiterName}}{{.Order.WaiterName}}{{else}}N/A{{end}}
Bàn: {{if .Order.TableNumber}}{{.Order.TableNumber}}{{else}}Mang về{{end}}     Ngày: {{formatTime .Order.CreatedAt "02/01/2006"}}
Giờ vào: {{formatTime .Order.CreatedAt "15:04"}}     Giờ ra: {{formatTime .Order.CreatedAt "15:04"}}

[TABLE_START]
STT | Tên món           | SL | Đơn giá  | Thành tiền
----|-------------------|----|-----------|-----------
{{range $index, $item := .Order.Items}}{{add $index 1}} | {{truncate $item.Name 17}} | {{$item.Quantity}} | {{formatPrice $item.Price}} | {{formatPrice $item.Subtotal}}{{if $item.VariantName}}
    | - {{truncate $item.VariantName 15}} | {{$item.Quantity}} | 0 | 0{{end}}
{{end}}[TABLE_END]

Thành tiền:                    {{formatPrice .Order.Subtotal}} đ
{{if gt .Order.Discount 0.0}}
Giảm giá:                      -{{formatPrice .Order.Discount}} đ
{{end}}
Tổng tiền:                     {{formatPrice .Order.Total}} đ

{{if .ShowCustomMessage}}{{.CustomMessage}}
{{end}}Cảm ơn Quý Khách
Powered by iPOS.vn`,
  created_at: new Date(),
  updated_at: new Date()
};

// Unset other templates as default
db.print_templates.updateMany(
  { type: "BILL", is_default: true },
  { $set: { is_default: false } }
);

// Insert new template
const result = db.print_templates.insertOne(template);
print("✅ Template inserted with ID:", result.insertedId);
print("📝 Template name: Bill với Logo");
print("🎯 Set as default: true");
