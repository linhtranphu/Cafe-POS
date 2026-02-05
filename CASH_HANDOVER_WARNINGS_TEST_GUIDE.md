# Cash Handover Warnings - Test Guide 🧪

## 🎯 Mục đích

Test 2 loại cảnh báo trong quá trình bàn giao:
1. **Discrepancy Warning** - Chênh lệch giữa declared và actual
2. **Shift Cash Warning** - Chênh lệch giữa declared và remaining cash

## 🚀 Setup

### 1. Rebuild Backend (Bắt buộc cho Shift Cash Warning)

```bash
cd backend
go build -o cafe-pos-server
```

### 2. Restart Services

```bash
docker-compose restart backend
# hoặc
docker-compose up -d
```

### 3. Open Frontend

```
http://localhost:5173
```

## 🧪 Test Cases

### Test Case 1: Discrepancy Warning - Thiếu tiền

**Mục tiêu:** Test cảnh báo khi cashier nhận ít hơn số tiền khai báo

**Steps:**
1. Login as waiter (waiter1/password123)
2. Tạo shift mới với start_cash = 500,000₫
3. Tạo order và thanh toán bằng CASH = 50,000₫
4. Tạo handover với declared_amount = 50,000₫
5. Login as cashier (cashier1/password123)
6. Start cashier shift
7. Vào `/cashier/handovers`
8. Click "Xác nhận" trên handover
9. Nhập actual_amount = 45,000₫ (thiếu 5k)

**Expected:**
- 🔴 Cảnh báo đỏ hiển thị: "⚠️ Thiếu tiền"
- Chênh lệch: 5,000₫
- Required: Lý do chênh lệch
- Required: Trách nhiệm
- Không thể submit nếu thiếu thông tin

---

### Test Case 2: Discrepancy Warning - Thừa tiền

**Mục tiêu:** Test cảnh báo khi cashier nhận nhiều hơn số tiền khai báo

**Steps:**
1. Tiếp tục từ Test Case 1
2. Tạo handover mới với declared_amount = 50,000₫
3. Cashier nhập actual_amount = 55,000₫ (thừa 5k)

**Expected:**
- 🟢 Cảnh báo xanh hiển thị: "⚠️ Thừa tiền"
- Chênh lệch: 5,000₫
- Required: Lý do chênh lệch
- Required: Trách nhiệm

---

### Test Case 3: Discrepancy Warning - Cần Manager Approval

**Mục tiêu:** Test cảnh báo khi chênh lệch > 100k

**Steps:**
1. Tạo handover với declared_amount = 200,000₫
2. Cashier nhập actual_amount = 50,000₫ (thiếu 150k)

**Expected:**
- 🔴 Cảnh báo đỏ
- 🟠 Badge cam: "Cần manager phê duyệt"
- Chênh lệch: 150,000₫
- Required: Lý do + Trách nhiệm

---

### Test Case 4: Shift Cash Warning - UNDER_DECLARED

**Mục tiêu:** Test cảnh báo khi waiter khai báo ít hơn remaining cash

**Setup:**
```
Waiter shift:
- Start cash: 500,000₫
- Cash revenue: 200,000₫ (từ orders)
- Handed over: 0₫
- Remaining cash: 700,000₫
```

**Steps:**
1. Login as waiter
2. Tạo shift với start_cash = 500,000₫
3. Tạo nhiều orders, tổng cash = 200,000₫
4. Tạo handover với declared_amount = 400,000₫ (PARTIAL)
5. Login as cashier
6. Vào `/cashier/handovers`

**Expected:**
- 🟡 Cảnh báo vàng trong danh sách:
  - "Khai báo ít hơn tiền còn lại trong ca"
  - "Tiền còn lại: 700,000₫ | Chênh: 300,000₫"
- Click "Xác nhận" → Cảnh báo lớn hơn trong modal

**Lý do hợp lệ:**
- Handover type = PARTIAL
- Waiter giữ lại 300k để tiếp tục làm việc

---

### Test Case 5: Shift Cash Warning - OVER_DECLARED

**Mục tiêu:** Test cảnh báo khi waiter khai báo nhiều hơn remaining cash

**Setup:**
```
Waiter shift:
- Start cash: 500,000₫
- Cash revenue: 200,000₫
- Handed over: 0₫
- Remaining cash: 700,000₫
```

**Steps:**
1. Waiter tạo handover với declared_amount = 900,000₫ (nhiều hơn 200k)
2. Cashier vào `/cashier/handovers`

**Expected:**
- 🟠 Cảnh báo cam trong danh sách:
  - "Khai báo nhiều hơn tiền còn lại trong ca"
  - "Tiền còn lại: 700,000₫ | Chênh: 200,000₫"
- Click "Xác nhận" → Cảnh báo lớn hơn trong modal

**Action:**
- ⚠️ DỪNG LẠI - Không xác nhận
- Gọi waiter kiểm tra
- Báo manager

---

### Test Case 6: Cả 2 loại cảnh báo cùng lúc

**Mục tiêu:** Test khi có cả shift cash warning và discrepancy warning

**Setup:**
```
Shift remaining: 700,000₫
Declared: 900,000₫ (OVER_DECLARED)
Actual: 850,000₫ (thiếu 50k so với declared)
```

**Steps:**
1. Waiter tạo handover với declared = 900,000₫
2. Cashier thấy cảnh báo cam (OVER_DECLARED)
3. Cashier nhập actual = 850,000₫
4. Thấy thêm cảnh báo đỏ (SHORTAGE)

**Expected:**
- 🟠 Cảnh báo cam: Shift cash mismatch
- 🔴 Cảnh báo đỏ: Discrepancy
- Required: Lý do + Trách nhiệm cho discrepancy

---

### Test Case 7: END_SHIFT - Phải khớp hoàn toàn

**Mục tiêu:** Test END_SHIFT handover phải bàn giao toàn bộ

**Setup:**
```
Shift remaining: 700,000₫
```

**Steps:**
1. Waiter tạo handover với type = END_SHIFT
2. Declared amount phải = 700,000₫

**Expected:**
- Nếu declared = 700,000₫ → Không có cảnh báo shift cash
- Nếu declared ≠ 700,000₫ → Cảnh báo (không hợp lệ cho END_SHIFT)

---

## 📊 Test Matrix

| Test Case | Shift Remaining | Declared | Actual | Shift Warning | Discrepancy Warning |
|-----------|----------------|----------|--------|---------------|---------------------|
| 1 | 700k | 50k | 45k | 🟡 UNDER | 🔴 SHORTAGE |
| 2 | 700k | 50k | 55k | 🟡 UNDER | 🟢 OVERAGE |
| 3 | 700k | 200k | 50k | 🟡 UNDER | 🔴 SHORTAGE + 🟠 Approval |
| 4 | 700k | 400k | 400k | 🟡 UNDER | ✅ None |
| 5 | 700k | 900k | 900k | 🟠 OVER | ✅ None |
| 6 | 700k | 900k | 850k | 🟠 OVER | 🔴 SHORTAGE |
| 7 | 700k | 700k | 700k | ✅ None | ✅ None |

## 🔍 Verification Checklist

### UI Elements
- [ ] Cảnh báo hiển thị đúng màu sắc
- [ ] Text hiển thị rõ ràng
- [ ] Số tiền format đúng (VNĐ)
- [ ] Icon hiển thị đúng (⚠️, 📉, 📈, 🔔)
- [ ] Cảnh báo trong list và modal đều hiển thị

### Validation
- [ ] Không thể submit nếu thiếu lý do (khi có discrepancy)
- [ ] Không thể submit nếu thiếu trách nhiệm (khi có discrepancy)
- [ ] Form reset sau khi submit thành công

### Backend
- [ ] API trả về shift_remaining_cash
- [ ] API trả về shift_current_cash
- [ ] API trả về shift_cash_revenue
- [ ] Discrepancy được lưu đúng
- [ ] Status được set đúng (CONFIRMED/DISCREPANCY)

### Data Integrity
- [ ] Shift cash được cập nhật đúng
- [ ] Cashier shift cash được cập nhật đúng
- [ ] Discrepancy record được tạo (nếu có)
- [ ] Audit trail đầy đủ

## 🐛 Common Issues

### Issue 1: Shift cash warning không hiển thị

**Nguyên nhân:** Backend chưa rebuild

**Fix:**
```bash
cd backend
go build -o cafe-pos-server
docker-compose restart backend
```

### Issue 2: shift_remaining_cash = undefined

**Nguyên nhân:** API response không có field này

**Check:**
```bash
curl -X GET "http://localhost:8080/cash-handovers/pending" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Expected response:**
```json
[
  {
    "id": "...",
    "declared_amount": 400000,
    "shift_remaining_cash": 700000,  // ← Phải có field này
    "shift_current_cash": 700000,
    "shift_cash_revenue": 200000
  }
]
```

### Issue 3: Cảnh báo không đúng màu

**Check:** CSS classes trong Vue component
- OVER_DECLARED → `bg-orange-50 border-orange-300`
- UNDER_DECLARED → `bg-yellow-50 border-yellow-300`
- SHORTAGE → `bg-red-50 border-red-300`
- OVERAGE → `bg-green-50 border-green-300`

## 📝 Test Report Template

```
Test Date: ___________
Tester: ___________

Test Case 1: Discrepancy - Thiếu tiền
[ ] PASS  [ ] FAIL
Notes: _______________________________

Test Case 2: Discrepancy - Thừa tiền
[ ] PASS  [ ] FAIL
Notes: _______________________________

Test Case 3: Discrepancy - Manager Approval
[ ] PASS  [ ] FAIL
Notes: _______________________________

Test Case 4: Shift Cash - UNDER_DECLARED
[ ] PASS  [ ] FAIL
Notes: _______________________________

Test Case 5: Shift Cash - OVER_DECLARED
[ ] PASS  [ ] FAIL
Notes: _______________________________

Test Case 6: Cả 2 cảnh báo
[ ] PASS  [ ] FAIL
Notes: _______________________________

Test Case 7: END_SHIFT
[ ] PASS  [ ] FAIL
Notes: _______________________________

Overall Result: [ ] PASS  [ ] FAIL
```

## 🚀 Next Steps

Sau khi test xong:
1. [ ] Fix bugs nếu có
2. [ ] Update documentation
3. [ ] Train cashiers
4. [ ] Deploy to production
5. [ ] Monitor for issues

## 📚 Related Documents

- [CASH_HANDOVER_DISCREPANCY_FIX.md](./docs/CASH_HANDOVER_DISCREPANCY_FIX.md)
- [CASH_HANDOVER_SHIFT_CASH_WARNING.md](./docs/CASH_HANDOVER_SHIFT_CASH_WARNING.md)
- [CASH_HANDOVER_DISCREPANCY_SUMMARY.md](./CASH_HANDOVER_DISCREPANCY_SUMMARY.md)
