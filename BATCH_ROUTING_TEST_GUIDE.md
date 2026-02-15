# Hướng Dẫn Test Batch Routing

## 🎯 Mục đích
Kiểm tra các routes batch đã được thêm vào và hoạt động đúng.

## ✅ Các Routes Đã Thêm

### 1. Dashboard Batch Chính
```
URL: http://localhost:5173/#/batch
Quyền: Manager only
Mô tả: Trang chính với 4 tabs (Định nghĩa, Batch, Cảnh báo, Báo cáo)
```

### 2. Quản Lý Batch Records
```
URL: http://localhost:5173/#/batch/records
Quyền: Tất cả user đã đăng nhập
Mô tả: Danh sách batch records với filters và sorting
```

### 3. Tạo Batch Record Mới
```
URL: http://localhost:5173/#/batch/records/create
Quyền: Tất cả user đã đăng nhập
Mô tả: Form tạo batch record mới
```

### 4. Chi Tiết Batch Record
```
URL: http://localhost:5173/#/batch/records/:id
Quyền: Tất cả user đã đăng nhập
Mô tả: Xem chi tiết batch record, ingredients used, cost breakdown
```

### 5. Quản Lý Batch Definitions
```
URL: http://localhost:5173/#/batch/definitions
Quyền: Manager only
Mô tả: Danh sách batch definitions
```

### 6. Tạo Batch Definition
```
URL: http://localhost:5173/#/batch/definitions/create
Quyền: Manager only
Mô tả: Form tạo batch definition mới
```

### 7. Cảnh Báo Batch
```
URL: http://localhost:5173/#/batch/alerts
Quyền: Tất cả user đã đăng nhập
Mô tả: Panel hiển thị low stock, expiring, expired alerts
```

### 8. Báo Cáo Batch
```
URL: http://localhost:5173/#/batch/reports
Quyền: Manager only
Mô tả: Production, Wastage, Usage reports
```

## 🧪 Test Cases

### Test 1: Truy cập từ Dashboard
1. Đăng nhập với tài khoản Manager
2. Vào Dashboard
3. Tìm BatchStatusWidget (widget màu xanh với icon 🧪)
4. Click vào widget hoặc nút "Quản lý batch"
5. ✅ Kỳ vọng: Chuyển đến `/batch` với 4 tabs

### Test 2: Tạo Batch Record
1. Vào `/batch/records` hoặc click tab "Batch" trong `/batch`
2. Click nút "➕ Ghi Nhận Batch Mới" (màu xanh lá)
3. ✅ Kỳ vọng: Chuyển đến `/batch/records/create` và hiển thị form
4. ✅ Kỳ vọng: KHÔNG còn blank page

### Test 3: Xem Chi Tiết Batch
1. Vào `/batch/records`
2. Click nút "👁️ Xem" trên bất kỳ batch record nào
3. ✅ Kỳ vọng: Chuyển đến `/batch/records/:id` và hiển thị chi tiết

### Test 4: Navigation giữa các tabs
1. Vào `/batch`
2. Click qua các tabs: Định nghĩa → Batch → Cảnh báo → Báo cáo
3. ✅ Kỳ vọng: Mỗi tab hiển thị đúng component

### Test 5: Direct URL Access
1. Copy paste trực tiếp URL vào browser:
   - `http://localhost:5173/#/batch/records/create`
   - `http://localhost:5173/#/batch/alerts`
   - `http://localhost:5173/#/batch/reports`
2. ✅ Kỳ vọng: Mỗi URL load đúng trang

### Test 6: Authorization
1. Đăng nhập với tài khoản Barista
2. Thử truy cập `/batch/definitions`
3. ✅ Kỳ vọng: Redirect về dashboard (không có quyền)
4. Thử truy cập `/batch/records`
5. ✅ Kỳ vọng: Hiển thị danh sách (có quyền)

### Test 7: Back Navigation
1. Vào `/batch/records/create`
2. Click nút back (←) trên header
3. ✅ Kỳ vọng: Quay về `/batch/records`

### Test 8: Error States
1. Tắt backend server
2. Vào `/batch/records`
3. ✅ Kỳ vọng: Hiển thị ErrorState với nút "Thử lại"
4. Click "Thử lại"
5. ✅ Kỳ vọng: Gọi lại API (không crash)

## 🐛 Các Lỗi Đã Fix

### 1. Blank Page Issue
**Vấn đề**: `/batch/records/create` hiển thị blank page
**Nguyên nhân**: Route không tồn tại trong router config
**Giải pháp**: Đã thêm route và tạo view page `BatchRecordFormView.vue`

### 2. Vue Router Warning
**Vấn đề**: "No match found for location with path '/batch/records/create'"
**Nguyên nhân**: Route chưa được định nghĩa
**Giải pháp**: Đã thêm đầy đủ 8 routes cho batch system

### 3. ErrorState Callback
**Vấn đề**: `loadRecords` function không tồn tại
**Nguyên nhân**: Sai tên function trong ErrorState component
**Giải pháp**: Đổi thành `() => batchRecordStore.fetchRecords()`

## 📱 Mobile Testing

Test trên các kích thước màn hình:
- iPhone SE (375px)
- iPhone 12 Pro (390px)
- iPad (768px)

Kiểm tra:
- [ ] Header sticky đúng
- [ ] Bottom navigation không che nội dung
- [ ] Buttons có kích thước phù hợp để tap
- [ ] Forms dễ điền trên mobile
- [ ] Modals hiển thị đúng

## 🎨 UI/UX Checklist

- [ ] Loading states hiển thị khi fetch data
- [ ] Error states có nút retry
- [ ] Empty states có call-to-action
- [ ] Transitions mượt mà
- [ ] Colors consistent với design system
- [ ] Icons rõ ràng và có ý nghĩa

## 📊 Performance

- [ ] Routes load nhanh (< 1s)
- [ ] Lazy loading hoạt động (check Network tab)
- [ ] Không có memory leaks khi navigate
- [ ] Smooth scrolling

## ✅ Kết Luận

Sau khi test tất cả các cases trên, batch routing system đã hoàn chỉnh và sẵn sàng sử dụng.

## 🚀 Next Steps

1. Test manual tất cả routes
2. Implement Task 14.1 - MenuRecipeEditor batch support
3. Complete Task 16 - Styling & UX improvements
4. Write E2E tests cho routing flows
