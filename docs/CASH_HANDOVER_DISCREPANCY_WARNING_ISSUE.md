# Vấn đề: Không hiển thị cảnh báo Discrepancy trong quá trình bàn giao

## 🔍 Phân tích vấn đề

Trong quá trình bàn giao giữa waiter và cashier tại màn hình `/cashier/handovers`, **KHÔNG có cảnh báo discrepancy** được hiển thị khi cashier nhập số tiền thực nhận khác với số tiền khai báo.

## 🎯 Nguyên nhân chính

### 1. **UI không tính toán và hiển thị discrepancy trong real-time**

Trong file `frontend/src/views/CashierHandoverView.vue`:

```vue
<!-- Confirm Modal -->
<form @submit.prevent="confirmHandover" class="space-y-4">
  <!-- Actual Amount (only for CONFIRMED) -->
  <div v-if="confirmAction === 'CONFIRMED'">
    <label class="block text-sm font-medium mb-2">Số tiền thực nhận (VNĐ) *</label>
    <input v-model.number="confirmForm.actual_amount" 
      type="number" 
      min="0" 
      step="1000" 
      required 
      class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-blue-500">
  </div>
  
  <!-- Cashier Note -->
  <div>
    <label class="block text-sm font-medium mb-2">
      {{ confirmAction === 'CONFIRMED' ? 'Ghi chú xác nhận' : 'Lý do từ chối' }}
    </label>
    <textarea v-model="confirmForm.cashier_note" ...></textarea>
  </div>
</form>
```

**Vấn đề:**
- ❌ Không có computed property để tính `discrepancy = actual_amount - declared_amount`
- ❌ Không có cảnh báo hiển thị khi có chênh lệch
- ❌ Không yêu cầu nhập `discrepancy_reason` khi có chênh lệch
- ❌ Không yêu cầu chọn `discrepancy_responsibility` khi có chênh lệch

### 2. **Backend yêu cầu thông tin discrepancy nhưng frontend không gửi**

Backend service (`cash_handover_service.go`) xử lý discrepancy:

```go
// 4. Calculate discrepancy
discrepancy := req.ActualAmount - h.DeclaredAmount

// 6. Handle discrepancy if exists
if h.HasDiscrepancy() {
    h.DiscrepancyReason = req.DiscrepancyReason
    h.DiscrepancyResponsibility = req.DiscrepancyResponsibility
    
    // Check if requires manager approval
    if h.RequiresManagerApproval(s.discrepancyThreshold) {
        h.RequiresApproval = true
        h.Status = handover.StatusDiscrepancy
    }
    
    // Create discrepancy record
    if err := s.createDiscrepancyRecord(ctx, h); err != nil {
        return err
    }
}
```

**Backend mong đợi:**
- `discrepancy_reason` (string) - Lý do chênh lệch
- `discrepancy_responsibility` (ResponsibilityType) - Trách nhiệm: WAITER, CASHIER, SYSTEM, CUSTOMER, UNKNOWN

**Frontend hiện tại gửi:**
```javascript
const data = {
  status: confirmAction.value,
  cashier_note: confirmForm.value.cashier_note
}

if (confirmAction.value === 'CONFIRMED') {
  data.actual_amount = confirmForm.value.actual_amount
}
```

❌ **Không gửi `discrepancy_reason` và `discrepancy_responsibility`**

### 3. **Logic phát hiện discrepancy**

Backend có logic rõ ràng:

```go
// HasDiscrepancy checks if handover has discrepancy
func (h *CashHandover) HasDiscrepancy() bool {
    return h.Discrepancy != 0
}

// RequiresManagerApproval checks if requires manager approval (large discrepancy)
func (h *CashHandover) RequiresManagerApproval(threshold float64) bool {
    return h.HasDiscrepancy() && (h.Discrepancy > threshold || h.Discrepancy < -threshold)
}
```

**Threshold:** 100,000 VNĐ
- Nếu `|discrepancy| > 100,000` → Cần manager approval
- Nếu `|discrepancy| <= 100,000` → Cashier có thể xử lý

## 🛠️ Giải pháp

### Bước 1: Thêm computed property tính discrepancy

```vue
<script setup>
// ... existing code ...

const discrepancy = computed(() => {
  if (!selectedHandover.value || !confirmForm.value.actual_amount) return 0
  return confirmForm.value.actual_amount - selectedHandover.value.declared_amount
})

const hasDiscrepancy = computed(() => discrepancy.value !== 0)

const discrepancyType = computed(() => {
  if (discrepancy.value < 0) return 'SHORTAGE' // Thiếu
  if (discrepancy.value > 0) return 'OVERAGE'  // Thừa
  return null
})

const requiresManagerApproval = computed(() => {
  return Math.abs(discrepancy.value) > 100000
})
</script>
```

### Bước 2: Hiển thị cảnh báo discrepancy trong modal

```vue
<!-- Discrepancy Warning (after Actual Amount input) -->
<div v-if="hasDiscrepancy && confirmAction === 'CONFIRMED'" 
  class="p-4 rounded-xl border-2"
  :class="discrepancy > 0 ? 'bg-green-50 border-green-300' : 'bg-red-50 border-red-300'">
  
  <div class="flex items-start gap-3 mb-3">
    <span class="text-2xl">{{ discrepancy > 0 ? '📈' : '📉' }}</span>
    <div class="flex-1">
      <h4 class="font-bold" :class="discrepancy > 0 ? 'text-green-800' : 'text-red-800'">
        {{ discrepancy > 0 ? '⚠️ Thừa tiền' : '⚠️ Thiếu tiền' }}
      </h4>
      <p class="text-sm mt-1" :class="discrepancy > 0 ? 'text-green-700' : 'text-red-700'">
        Chênh lệch: <strong>{{ formatPrice(Math.abs(discrepancy)) }}</strong>
      </p>
      <p v-if="requiresManagerApproval" class="text-sm mt-2 font-medium text-orange-700">
        🔔 Chênh lệch lớn hơn 100,000₫ - Cần manager phê duyệt
      </p>
    </div>
  </div>
  
  <!-- Discrepancy Reason (Required) -->
  <div class="mb-3">
    <label class="block text-sm font-medium mb-2">Lý do chênh lệch *</label>
    <textarea v-model="confirmForm.discrepancy_reason" 
      required
      rows="2" 
      class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500"
      placeholder="Giải thích nguyên nhân chênh lệch..."></textarea>
  </div>
  
  <!-- Discrepancy Responsibility (Required) -->
  <div>
    <label class="block text-sm font-medium mb-2">Trách nhiệm *</label>
    <select v-model="confirmForm.discrepancy_responsibility" 
      required
      class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500">
      <option value="">-- Chọn người chịu trách nhiệm --</option>
      <option value="WAITER">Waiter</option>
      <option value="CASHIER">Cashier (Tôi)</option>
      <option value="CUSTOMER">Khách hàng</option>
      <option value="SYSTEM">Hệ thống</option>
      <option value="UNKNOWN">Chưa rõ</option>
    </select>
  </div>
</div>
```

### Bước 3: Cập nhật form data structure

```javascript
const confirmForm = ref({
  actual_amount: 0,
  cashier_note: '',
  discrepancy_reason: '',
  discrepancy_responsibility: ''
})
```

### Bước 4: Gửi đầy đủ thông tin discrepancy

```javascript
const confirmHandover = async () => {
  try {
    const data = {
      status: confirmAction.value,
      cashier_note: confirmForm.value.cashier_note
    }
    
    // Add actual_amount only for CONFIRMED
    if (confirmAction.value === 'CONFIRMED') {
      data.actual_amount = confirmForm.value.actual_amount
      
      // Add discrepancy info if exists
      if (hasDiscrepancy.value) {
        data.discrepancy_reason = confirmForm.value.discrepancy_reason
        data.discrepancy_responsibility = confirmForm.value.discrepancy_responsibility
      }
    }
    
    await cashierStore.confirmHandover(selectedHandover.value.id, data)
    
    // ... rest of the code
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}
```

### Bước 5: Reset form khi đóng modal

```javascript
const showConfirmModal = (handover, action) => {
  selectedHandover.value = handover
  confirmAction.value = action
  confirmForm.value = {
    actual_amount: handover.declared_amount, // Default to declared amount
    cashier_note: '',
    discrepancy_reason: '',
    discrepancy_responsibility: ''
  }
  showConfirmForm.value = true
}
```

## 📊 Flow hoàn chỉnh

```
1. Cashier mở modal xác nhận bàn giao
   ↓
2. Nhập số tiền thực nhận (actual_amount)
   ↓
3. Computed property tự động tính: discrepancy = actual_amount - declared_amount
   ↓
4. Nếu discrepancy !== 0:
   ├─ Hiển thị cảnh báo màu đỏ (thiếu) hoặc xanh (thừa)
   ├─ Yêu cầu nhập lý do (discrepancy_reason) *
   ├─ Yêu cầu chọn trách nhiệm (discrepancy_responsibility) *
   └─ Nếu |discrepancy| > 100,000₫ → Hiển thị thông báo cần manager approval
   ↓
5. Submit form với đầy đủ thông tin
   ↓
6. Backend xử lý:
   ├─ Tính discrepancy
   ├─ Lưu discrepancy_reason và discrepancy_responsibility
   ├─ Tạo CashDiscrepancy record
   ├─ Nếu |discrepancy| > 100,000₫:
   │  ├─ Set status = DISCREPANCY
   │  ├─ Set requires_approval = true
   │  └─ Chờ manager phê duyệt
   └─ Nếu |discrepancy| <= 100,000₫:
      ├─ Set status = CONFIRMED
      └─ Cập nhật cash amounts ngay
```

## 🎨 UI/UX Improvements

### Màu sắc cảnh báo:
- **Thiếu tiền (SHORTAGE):** Nền đỏ nhạt, viền đỏ, text đỏ đậm
- **Thừa tiền (OVERAGE):** Nền xanh nhạt, viền xanh, text xanh đậm
- **Cần approval:** Badge cam với icon 🔔

### Validation:
- ✅ `actual_amount` là required
- ✅ Nếu có discrepancy → `discrepancy_reason` là required
- ✅ Nếu có discrepancy → `discrepancy_responsibility` là required
- ✅ Disable submit button nếu thiếu thông tin

### Real-time feedback:
- Hiển thị discrepancy ngay khi user nhập `actual_amount`
- Tự động focus vào `discrepancy_reason` khi có chênh lệch
- Hiển thị tooltip giải thích các loại responsibility

## 🔗 Files cần sửa

1. **frontend/src/views/CashierHandoverView.vue** - Thêm UI cảnh báo và validation
2. **frontend/src/stores/cashier.js** - Đã OK, không cần sửa
3. **frontend/src/services/handover.js** - Đã OK, không cần sửa
4. **backend** - Đã OK, logic đã hoàn chỉnh

## ✅ Checklist triển khai

- [ ] Thêm computed properties (discrepancy, hasDiscrepancy, discrepancyType, requiresManagerApproval)
- [ ] Thêm UI cảnh báo discrepancy trong modal
- [ ] Thêm input fields cho discrepancy_reason và discrepancy_responsibility
- [ ] Cập nhật confirmForm structure
- [ ] Cập nhật confirmHandover method để gửi discrepancy info
- [ ] Thêm validation cho required fields
- [ ] Test với các scenarios:
  - [ ] Không có chênh lệch (actual = declared)
  - [ ] Thiếu tiền nhỏ (< 100k)
  - [ ] Thừa tiền nhỏ (< 100k)
  - [ ] Thiếu tiền lớn (> 100k) → Cần approval
  - [ ] Thừa tiền lớn (> 100k) → Cần approval

## 📝 Notes

- Backend đã xử lý đầy đủ logic discrepancy
- Vấn đề chỉ nằm ở frontend UI/UX
- Cần thêm màn hình manager approval cho discrepancy lớn
- Có thể thêm history log để track discrepancy resolution
