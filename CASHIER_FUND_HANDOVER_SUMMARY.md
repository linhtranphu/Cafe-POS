# Tóm tắt: Tính năng Cashier Fund Handover

## Vấn đề cần giải quyết

Hiện tại, khi Cashier nhận tiền từ các waiter thông qua handover, hệ thống ghi nhận số tiền nhưng:
1. **Không hiển thị rõ ràng** tổng số tiền Cashier đang quản lý
2. **Không có quy trình** để Cashier handover lại số tiền này về quỹ khi đóng ca
3. **Không rõ trách nhiệm** của Cashier đối với số tiền đã nhận

## Giải pháp

### 1. Hiển thị "Tiền đang quản lý" trong Dashboard

Thêm một section mới trong `CashierDashboard.vue`:

```
┌─────────────────────────────────────┐
│ 💰 Tiền đang quản lý                │
├─────────────────────────────────────┤
│ 💵 Tiền mặt đã nhận                 │
│ 1,500,000₫                          │
│                                     │
│ 💳 Tiền CK đã nhận                  │
│ 800,000₫                            │
│                                     │
│ 📊 Tổng cộng                        │
│ 2,300,000₫                          │
│                                     │
│ ⚠️ Bạn chịu trách nhiệm số tiền này│
└─────────────────────────────────────┘
```

### 2. Quy trình Handover về Quỹ khi Đóng Ca

Mở rộng `CashierShiftClosureV2.vue` với các bước:

**Bước 1: Xem tổng quan**
- Hiển thị tiền đầu ca
- Hiển thị tiền nhận từ waiter (mặt + CK)
- Tính tổng tiền mặt lý thuyết

**Bước 2: Đếm tiền mặt**
- Nhập số tiền thực tế đếm được
- Tự động tính chênh lệch (variance)

**Bước 3: Giải thích chênh lệch (nếu có)**
- Chọn lý do (lỗi đếm, mất mát, tranh chấp...)
- Nhập ghi chú chi tiết (tối thiểu 10 ký tự)

**Bước 4: Xác nhận bàn giao**
- Xem tóm tắt: tiền mặt + tiền CK + chênh lệch
- Xác nhận bàn giao về quỹ
- Đóng ca

### 3. Ghi nhận vào Database

Tạo collection mới: `fund_handovers`

```javascript
{
  _id: ObjectId,
  cashier_shift_id: ObjectId,
  cashier_id: ObjectId,
  cashier_name: String,
  
  cash_amount: Number,           // Tiền mặt thực tế bàn giao
  transfer_amount: Number,       // Tiền CK ghi nhận
  total_amount: Number,          // Tổng
  
  expected_cash: Number,         // Tiền mặt lý thuyết
  variance_amount: Number,       // Chênh lệch
  variance_reason: String,       // Lý do (nếu có)
  variance_notes: String,        // Ghi chú (nếu có)
  
  receiver_id: ObjectId,         // Người nhận (nullable - dành cho tương lai)
  receiver_name: String,         // Tên người nhận (nullable)
  
  handover_at: Date,
  created_at: Date,
  updated_at: Date
}
```

## Luồng hoạt động

### Luồng 1: Xem tiền đang quản lý

```
1. Cashier mở CashierDashboard
2. Hệ thống lấy thông tin cashier shift hiện tại
3. Hiển thị:
   - ReceivedCash (từ CashierShift)
   - ReceivedTransfer (từ CashierShift)
   - Tổng = ReceivedCash + ReceivedTransfer
4. Hiển thị cảnh báo về trách nhiệm
```

### Luồng 2: Đóng ca với Fund Handover

```
1. Cashier click "Đóng ca thu ngân"
2. Hệ thống hiển thị tổng quan:
   - Starting Float: 500,000₫
   - Received Cash: 1,500,000₫
   - Received Transfer: 800,000₫
   - Expected Cash: 2,000,000₫ (Starting + Received)

3. Cashier đếm tiền và nhập: 1,995,000₫
4. Hệ thống tính variance: -5,000₫ (thiếu)

5. Cashier chọn lý do: "Lỗi đếm tiền"
6. Cashier nhập ghi chú: "Đếm nhầm tờ 50k thành 100k"

7. Cashier xác nhận bàn giao
8. Hệ thống thực hiện TRANSACTION:
   a. Tạo FundHandover record
   b. Ghi nhận actual_cash vào CashierShift
   c. Đóng CashierShift (status = CLOSED)
   d. Ghi audit log
9. Nếu thành công → Redirect về dashboard
10. Nếu lỗi → Rollback toàn bộ, hiển thị lỗi
```

## Thiết kế mở rộng

### Receiver (Người nhận) - Dành cho tương lai

Hiện tại: `receiver_id` và `receiver_name` = `null` (bàn giao về quỹ chung)

Tương lai: Có thể chỉ định Manager làm người nhận:
```javascript
{
  receiver_id: ObjectId("manager_id"),
  receiver_name: "Nguyễn Văn Manager"
}
```

API đã thiết kế sẵn:
```javascript
POST /api/cashier/shifts/:id/close
{
  "actual_cash": 1995000,
  "variance_reason": "COUNTING_ERROR",
  "variance_notes": "...",
  "receiver_id": "manager_id"  // Optional
}
```

## Lợi ích

### Cho Cashier:
1. ✅ Biết rõ mình đang quản lý bao nhiêu tiền
2. ✅ Quy trình đóng ca rõ ràng, có hướng dẫn từng bước
3. ✅ Ghi nhận chênh lệch một cách minh bạch
4. ✅ Hoàn tất trách nhiệm khi đóng ca

### Cho Hệ thống:
1. ✅ Audit trail đầy đủ cho mọi giao dịch
2. ✅ Dữ liệu nhất quán (sử dụng transaction)
3. ✅ Dễ dàng truy vấn lịch sử handover
4. ✅ Sẵn sàng mở rộng (receiver, approval workflow...)

### Cho Quản lý:
1. ✅ Theo dõi được luồng tiền rõ ràng
2. ✅ Phát hiện được các vấn đề (variance patterns)
3. ✅ Có dữ liệu để kiểm toán
4. ✅ Chuẩn bị sẵn để thêm approval workflow

## Thời gian thực hiện

- **Backend**: ~12 giờ (domain model, repository, service, API)
- **Frontend Dashboard**: ~3 giờ (hiển thị managed funds)
- **Frontend Closure**: ~9 giờ (4 bước trong closure flow)
- **Testing**: ~10 giờ (unit, integration, E2E)
- **Documentation**: ~3 giờ
- **Deployment**: ~2 giờ

**Tổng**: ~42 giờ (~5-6 ngày làm việc)

## Files cần tạo/sửa

### Backend (NEW):
- `backend/domain/cashier/fund_handover.go` ⭐ NEW
- `backend/infrastructure/mongodb/fund_handover_repository.go` ⭐ NEW

### Backend (MODIFY):
- `backend/application/services/cashier_shift_service.go`
- `backend/api/handlers/cashier_handler.go`
- `backend/api/routes/cashier_routes.go`

### Frontend (MODIFY):
- `frontend/src/views/CashierDashboard.vue`
- `frontend/src/views/CashierShiftClosureV2.vue`

### Database:
- Create collection: `fund_handovers`
- Create indexes

## Câu hỏi thường gặp

### Q1: Tại sao cần ghi nhận transfer_amount nếu không bàn giao vật lý?
**A**: Để có audit trail đầy đủ. Cashier chịu trách nhiệm xác nhận số tiền CK đã nhận từ waiter. Khi đóng ca, cần ghi nhận rằng số tiền này đã được "chuyển giao trách nhiệm" về quỹ/kế toán.

### Q2: Nếu không có chênh lệch thì có cần ghi lý do không?
**A**: Không. Chỉ khi có chênh lệch (variance ≠ 0) thì mới bắt buộc nhập lý do và ghi chú.

### Q3: Có thể hủy việc đóng ca sau khi đã xác nhận không?
**A**: Không. Sau khi xác nhận, transaction sẽ commit và shift chuyển sang CLOSED. Đây là thiết kế có chủ đích để đảm bảo tính toàn vẹn.

### Q4: Receiver_id dùng để làm gì?
**A**: Hiện tại để null (bàn giao về quỹ chung). Tương lai có thể mở rộng để chỉ định Manager cụ thể nhận tiền, kèm theo approval workflow.

### Q5: Nếu có lỗi trong quá trình đóng ca thì sao?
**A**: Transaction sẽ rollback toàn bộ. Shift vẫn ở trạng thái CLOSURE_INITIATED, Cashier có thể thử lại.

## Kết luận

Tính năng này giải quyết triệt để vấn đề quản lý trách nhiệm của Cashier đối với số tiền nhận từ waiter, đồng thời tạo audit trail đầy đủ cho mục đích kiểm toán và quản lý. Thiết kế mở rộng (receiver) cho phép dễ dàng thêm tính năng trong tương lai mà không cần thay đổi cấu trúc dữ liệu.
