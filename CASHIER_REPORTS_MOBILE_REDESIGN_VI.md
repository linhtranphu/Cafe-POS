# Thiết kế lại Màn hình Báo cáo Thu ngân - Mobile First

## Tổng quan
Thiết kế lại màn hình CashierReports (`/cashier/reports`) theo hướng mobile-first với UX cải thiện và hiển thị giá tiền gọn gàng.

## Các thay đổi

### 1. Cải thiện Header
- **Giảm padding**: Đổi từ `py-4` sang `py-3` cho header gọn hơn
- **Tiêu đề nhỏ hơn**: Đổi từ `text-2xl` sang `text-xl`
- **Phụ đề nhỏ hơn**: Đổi từ `text-sm` sang `text-xs`
- **Nút quay lại**: Thêm nút "← Quay lại" khi đang xem báo cáo
- **Thống kê nhanh**: Thêm 3 cột thống kê khi không xem báo cáo
  - Tổng số ca làm
  - Tổng doanh thu (định dạng gọn)
  - Tổng đơn hàng

### 2. Định dạng giá tiền gọn gàng
Thêm hàm `formatCompactPrice()` cho tiền Việt Nam:
- **≥ 1,000,000**: Hiển thị "15.5tr" (triệu)
- **≥ 1,000**: Hiển thị "500k" (nghìn)
- **< 1,000**: Hiển thị "999đ"

Áp dụng cho:
- Thống kê nhanh ở header
- Thẻ tổng kết báo cáo
- Chi tiết thanh toán
- Số tiền đối soát
- Số tiền kiểm toán
- Lịch sử báo cáo

### 3. Hiển thị có điều kiện
- **Thẻ tạo báo cáo**: Chỉ hiện khi `!currentReport`
- **Lịch sử báo cáo**: Chỉ hiện khi `!currentReport`
- **Báo cáo hiện tại**: Chỉ hiện khi `currentReport` tồn tại
- Tạo giao diện sạch sẽ, tập trung hơn

### 4. Điều chỉnh Typography
Kích thước chữ thân thiện với mobile hơn:
- Tiêu đề báo cáo: `text-lg` → `text-base`
- Số liệu thống kê: Giảm từ `text-2xl`/`text-lg` xuống `text-xl`/`text-sm`
- Chi tiết thanh toán: Thêm `text-sm` cho số tiền
- Tiêu đề lịch sử: Thêm `text-sm`

### 5. Cải thiện nút In
- Đổi từ chỉ icon sang "🖨️ In" có chữ
- Style tốt hơn: `bg-blue-500 text-white`
- Rõ ràng và nổi bật hơn

### 6. Sửa hành vi Scroll
- Đổi từ `window.scrollTo()` sang container scroll
- Tìm container `.overflow-y-auto` và scroll nó
- Hoạt động đúng với pattern container scroll

### 7. Hàm mới được thêm
```javascript
// Định dạng giá gọn cho tiền Việt Nam lớn
formatCompactPrice(value)

// Computed property thống kê nhanh
quickStats

// Xóa báo cáo hiện tại
clearCurrentReport()
```

## Tính năng Mobile-First
✅ Pattern container scroll (h-screen)
✅ Sticky header với hỗ trợ safe area
✅ Hiển thị giá gọn cho số lớn
✅ Hỗ trợ pull-to-refresh
✅ Nút thân thiện với chạm, có active state
✅ Hiển thị có điều kiện cho UX sạch hơn
✅ Bottom navigation với khoảng trống pb-24

## Checklist kiểm tra
- [ ] Thống kê nhanh hiển thị đúng ở header
- [ ] Nút quay lại xóa báo cáo hiện tại
- [ ] Giá gọn hiển thị đúng (tr, k, đ)
- [ ] Thẻ tạo báo cáo ẩn khi xem báo cáo
- [ ] Lịch sử báo cáo ẩn khi xem báo cáo
- [ ] Nút in hoạt động
- [ ] Pull-to-refresh hoạt động trong container scroll
- [ ] Scroll lên đầu khi xem báo cáo
- [ ] Tất cả chữ đọc được trên mobile
- [ ] Vùng chạm đủ lớn

## File đã sửa
- `frontend/src/views/CashierReports.vue`

## Pattern liên quan
- ExpenseManagementView: Tham khảo định dạng giá gọn
- Pattern container scroll: Dùng ở 16 views
- Pull-to-refresh: Hoạt động với container scroll

## Ví dụ hiển thị tiền Việt Nam
```
15,500,000 → 15.5tr
2,000,000 → 2tr
500,000 → 500k
1,500 → 1.5k
999 → 999đ
```

---
**Trạng thái**: ✅ Hoàn thành
**Ngày**: 2026-02-07
**Pattern**: Mobile-first, container scroll, hiển thị gọn
