# Shift View Redesign - Implementation Guide

## Overview
This guide provides step-by-step instructions to implement the Cashier view redesign for ShiftView.vue.

## Current Status
The file `frontend/src/views/ShiftView.vue` has been partially modified. Due to the file's length (793 lines), we need to complete the implementation systematically.

## Implementation Approach

### Option 1: Manual Implementation (Recommended)
Follow the tasks in `tasks.md` and implement each section step by step.

### Option 2: Automated Script
Create a backup of the current file and apply changes programmatically.

## Key Changes Needed

### 1. Add New Imports
```javascript
import { useCashierStore } from '../stores/cashier'
import OverridePaymentModal from '../components/OverridePaymentModal.vue'
import DiscrepancyModal from '../components/DiscrepancyModal.vue'
import { ORDER_STATUS, PAYMENT_METHOD } from '../constants/order'
```

### 2. Add New State Variables
```javascript
// Cashier-specific state
const activeTab = ref('waiter') // 'waiter' | 'barista'
const selectedShiftId = ref('')
const selectedDate = ref(new Date().toISOString().split('T')[0])
const loadingShiftDetails = ref(false)
const loadingPayments = ref(false)

// Modal state
const showOverride = ref(false)
const showDiscrepancy = ref(false)
const selectedPayment = ref(null)

// Store
const cashierStore = useCashierStore()
```

### 3. Add Computed Properties
```javascript
// Filter shifts by date and role type
const filteredShifts = computed(() => {
  if (!isCashier.value) return shifts.value
  
  return shifts.value.filter(shift => {
    const shiftDate = new Date(shift.started_at).toISOString().split('T')[0]
    if (shiftDate !== selectedDate.value) return false
    
    if (activeTab.value === 'waiter') {
      return shift.role_type === 'waiter'
    } else {
      return shift.role_type === 'barista'
    }
  })
})

const selectedShift = computed(() => {
  if (!selectedShiftId.value) return null
  return shifts.value.find(s => s.id === selectedShiftId.value)
})

const isShiftOpen = computed(() => {
  return selectedShift.value?.status === SHIFT_STATUS.OPEN
})

const isWaiterShift = computed(() => {
  return selectedShift.value?.role_type === 'waiter'
})

const shiftStatus = computed(() => cashierStore.shiftStatus)
const payments = computed(() => cashierStore.payments)
```

### 4. Add Methods
```javascript
const onTabChange = (tab) => {
  activeTab.value = tab
  selectedShiftId.value = ''
}

const onDateChange = () => {
  selectedShiftId.value = ''
}

const onShiftSelect = async () => {
  if (!selectedShiftId.value) return
  
  loadingShiftDetails.value = true
  try {
    await cashierStore.getShiftStatus(selectedShiftId.value)
    
    if (isWaiterShift.value) {
      loadingPayments.value = true
      await cashierStore.getPaymentsByShift(selectedShiftId.value)
      loadingPayments.value = false
    }
  } catch (error) {
    console.error('Error loading shift details:', error)
    alert('Lỗi khi tải thông tin ca: ' + error.message)
  } finally {
    loadingShiftDetails.value = false
  }
}

const showCloseShiftModal = () => {
  selectedShift.value = selectedShift.value
  showCloseForm.value = true
}

const showOverrideModal = (payment) => {
  selectedPayment.value = payment
  showOverride.value = true
}

const showDiscrepancyModal = (payment) => {
  selectedPayment.value = payment
  showDiscrepancy.value = true
}

const handleOverridePayment = async (reason) => {
  try {
    await cashierStore.overridePayment(selectedPayment.value.order_id, reason)
    showOverride.value = false
    await cashierStore.getPaymentsByShift(selectedShiftId.value)
    alert('✅ Đã điều chỉnh thanh toán')
  } catch (error) {
    console.error('Override failed:', error)
    alert('❌ Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const handleReportDiscrepancy = async (data) => {
  try {
    await cashierStore.reportDiscrepancy({
      order_id: selectedPayment.value.order_id,
      ...data
    })
    showDiscrepancy.value = false
    alert('✅ Đã báo lỗi')
  } catch (error) {
    console.error('Report discrepancy failed:', error)
    alert('❌ Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

const lockOrder = async (orderId) => {
  if (confirm('Bạn có chắc muốn khóa order này? Không thể hoàn tác!')) {
    try {
      await cashierStore.lockOrder(orderId)
      await cashierStore.getPaymentsByShift(selectedShiftId.value)
      alert('✅ Đã khóa order')
    } catch (error) {
      console.error('Lock order failed:', error)
      alert('❌ Lỗi: ' + (error.response?.data?.error || error.message))
    }
  }
}

// Helper functions
const getPaymentMethodBadge = (method) => {
  const badges = {
    CASH: 'px-2 py-1 text-xs rounded-full bg-green-100 text-green-700 font-medium',
    TRANSFER: 'px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-700 font-medium',
    QR: 'px-2 py-1 text-xs rounded-full bg-purple-100 text-purple-700 font-medium'
  }
  return badges[method] || 'px-2 py-1 text-xs rounded-full bg-gray-100 text-gray-700 font-medium'
}

const getPaymentMethodText = (method) => {
  const texts = {
    CASH: '💵 Tiền mặt',
    TRANSFER: '💳 CK',
    QR: '📱 QR'
  }
  return texts[method] || method
}

const getStatusBadge = (status) => {
  const badges = {
    PAID: 'px-3 py-1 text-xs rounded-full bg-green-100 text-green-700 font-medium inline-block',
    QUEUED: 'px-3 py-1 text-xs rounded-full bg-yellow-100 text-yellow-700 font-medium inline-block',
    IN_PROGRESS: 'px-3 py-1 text-xs rounded-full bg-blue-100 text-blue-700 font-medium inline-block',
    READY: 'px-3 py-1 text-xs rounded-full bg-purple-100 text-purple-700 font-medium inline-block',
    SERVED: 'px-3 py-1 text-xs rounded-full bg-green-100 text-green-700 font-medium inline-block',
    LOCKED: 'px-3 py-1 text-xs rounded-full bg-red-100 text-red-700 font-medium inline-block'
  }
  return badges[status] || 'px-3 py-1 text-xs rounded-full bg-gray-100 text-gray-700 font-medium inline-block'
}

const getStatusText = (status) => {
  const texts = {
    PAID: '✓ Đã thu',
    QUEUED: '⏳ Chờ pha',
    IN_PROGRESS: '🍹 Đang pha',
    READY: '✅ Sẵn sàng',
    SERVED: '🎯 Hoàn tất',
    LOCKED: '🔒 Đã khóa'
  }
  return texts[status] || status
}

const formatDateTime = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleString('vi-VN', {
    day: '2-digit',
    month: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}
```

### 5. Add Modals to Template
Add these modals before the closing `</template>` tag:

```vue
<!-- Override Payment Modal -->
<OverridePaymentModal
  :show="showOverride"
  :payment="selectedPayment"
  @close="showOverride = false"
  @confirm="handleOverridePayment"
/>

<!-- Discrepancy Modal -->
<DiscrepancyModal
  :show="showDiscrepancy"
  :payment="selectedPayment"
  @close="showDiscrepancy = false"
  @confirm="handleReportDiscrepancy"
/>
```

## Next Steps

1. Complete the template section with all waiter/barista UI (copy from original file)
2. Add all modals (handover modals, close shift modal, etc.)
3. Complete the script section with all methods
4. Test the implementation
5. Fix any bugs or issues

## Testing Checklist

- [ ] Cashier can switch between waiter and barista tabs
- [ ] Date filter works correctly
- [ ] Shift selector shows filtered shifts
- [ ] Shift summary card displays correct data
- [ ] Payment list shows for waiter shifts only
- [ ] Payment actions work (override, discrepancy, lock)
- [ ] Close shift works for open shifts
- [ ] Waiter view remains unchanged
- [ ] Barista view remains unchanged
- [ ] No console errors
- [ ] Mobile responsive

## Notes

- The current file has been partially created but is incomplete
- You may need to restore from the original file and apply changes incrementally
- Consider creating a backup before making changes
- Test thoroughly after each major change
