# Debug Menu "áddd" - Không tính được chi phí

## Các nguyên nhân có thể:

### 1. Món không có nguyên liệu
- Kiểm tra: Vào trang quản lý menu, xem món "áddd" có nguyên liệu không
- Nếu không có → Thêm nguyên liệu cho món

### 2. Nguyên liệu không có giá (cost_per_unit = 0 hoặc null)
- Kiểm tra: Vào trang quản lý nguyên liệu
- Xem các nguyên liệu của món "áddd" có giá chưa
- Nếu chưa → Cập nhật giá cho nguyên liệu

### 3. Nguyên liệu bị xóa nhưng vẫn còn reference trong món
- Kiểm tra: ingredient_id trong món có tồn tại trong collection ingredients không
- Nếu không → Xóa reference cũ và thêm lại nguyên liệu đúng

### 4. Cost chưa được tính toán
- Kiểm tra: cost_status = null hoặc cost_last_calculated_at = null
- Giải pháp: Trigger tính toán lại bằng cách:
  - Cập nhật món (edit và save)
  - Hoặc cập nhật nguyên liệu
  - Hoặc gọi API calculate cost

## Cách kiểm tra từ Browser Console:

Mở trang http://localhost:5173/#/manager/menu-costs và chạy trong console:

```javascript
// 1. Tìm món "áddd" trong danh sách
const menuItems = await fetch('/api/manager/menu/costs').then(r => r.json())
const adddItem = menuItems.items.find(item => item.name.includes('áddd'))
console.log('Món áddd:', adddItem)

// 2. Xem chi tiết
if (adddItem) {
  const detail = await fetch(`/api/manager/menu/costs/${adddItem.menu_item_id}`).then(r => r.json())
  console.log('Chi tiết:', detail)
  console.log('Nguyên liệu:', detail.ingredients)
}
```

## Cách fix:

### Nếu thiếu nguyên liệu:
1. Vào http://localhost:5173/#/manager/menu
2. Tìm món "áddd"
3. Click Edit
4. Thêm nguyên liệu
5. Save

### Nếu nguyên liệu không có giá:
1. Vào http://localhost:5173/#/manager/ingredients
2. Tìm nguyên liệu của món "áddd"
3. Cập nhật giá (cost_per_unit)
4. Save
5. Hệ thống sẽ tự động tính lại chi phí món

### Force recalculate:
```bash
# Gọi API để tính lại chi phí
curl -X POST http://localhost:3000/api/manager/menu/{menu_item_id}/calculate-cost \
  -H "Authorization: Bearer YOUR_TOKEN"
```
