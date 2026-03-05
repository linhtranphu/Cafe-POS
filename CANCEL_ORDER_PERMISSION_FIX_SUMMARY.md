# Tóm tắt: Fix lỗi "insufficient permissions" khi hủy order

## Vấn đề
Waiter không thể hủy order vì backend đang chạy code cũ chưa có endpoint `/waiter/orders/:id/cancel`.

## Giải pháp
Backend cần được rebuild để áp dụng code mới.

## Cách fix nhanh

### Option 1: Script tự động (Khuyến nghị)
```bash
./rebuild-backend-with-cancel.sh
```

### Option 2: Thủ công
```bash
docker-compose stop backend
docker-compose rm -f backend
cd backend && docker build --no-cache -t cafe-pos-backend:latest . && cd ..
docker-compose up -d backend
```

### Option 3: Rebuild toàn bộ
```bash
docker-compose down
docker-compose up -d --build
```

## Kiểm tra sau khi fix

1. Check endpoint trong code:
```bash
grep "waiter.POST.*cancel" backend/main.go
```

2. Test trong app:
- Login với waiter
- Tạo order mới
- Tap "❌ Hủy order"
- Nhập lý do
- ✅ Thành công

## Files liên quan

- `backend/main.go` - Đã thêm endpoint
- `frontend/src/services/order.js` - Đã cập nhật
- `frontend/src/views/OrderView.vue` - Đã thêm UI
- `FIX_CANCEL_ORDER_PERMISSION.md` - Hướng dẫn chi tiết
- `rebuild-backend-with-cancel.sh` - Script rebuild
- `check-cancel-endpoint.sh` - Script kiểm tra

## Lưu ý

- ✅ Code đã đúng
- ✅ Frontend đã đúng
- ❌ Backend cần rebuild
- ⏱️ Downtime: ~5-10 giây
- 🔄 Không cần thay đổi database

---

**TL;DR:** Chạy `./rebuild-backend-with-cancel.sh` để fix.
