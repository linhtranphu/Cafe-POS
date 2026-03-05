# Fix giá nguyên liệu "cfe cot"

## Vấn đề
Nguyên liệu "cfe cot" có giá = 0 ₫/ml, nên món "áddd" không tính được chi phí.

## Giải pháp

### Bước 1: Tìm nguyên liệu trong giao diện
1. Vào: http://localhost:5173/#/manager/ingredients
2. Tìm "cfe cot" hoặc "cà phê cốt"
3. Xem giá hiện tại

### Bước 2: Cập nhật giá
Ví dụ tính giá:
- Nếu mua 1 lít (1000ml) cà phê cốt giá 100,000đ
- Thì giá/ml = 100,000 / 1000 = 100 đ/ml

Hoặc:
- Nếu mua 1 kg cà phê cốt giá 200,000đ
- Và 1 kg = 1000g
- Thì giá/g = 200,000 / 1000 = 200 đ/g

### Bước 3: Nhập giá vào hệ thống
1. Click Edit nguyên liệu
2. Nhập giá vào trường "Giá/đơn vị"
3. Đảm bảo đơn vị đúng (ml, g, kg, etc.)
4. Save

### Bước 4: Kiểm tra lại
1. Quay lại trang Menu Costs: http://localhost:5173/#/manager/menu-costs
2. Tìm món "áddd"
3. Chi phí sẽ được tính tự động
4. Cost Status sẽ chuyển từ "⚠ Thiếu dữ liệu" → "✓ Chính thức"

## Lưu ý
- Giá phải > 0 để tính được chi phí
- Đơn vị phải khớp với đơn vị trong công thức món
- Nếu có nhiều nguyên liệu thiếu giá, cần cập nhật tất cả
