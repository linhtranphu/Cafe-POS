# Shift View Redesign - Design Document

## Overview
This document describes the redesign of ShiftView component for Cashier role only. The redesign separates waiter and barista shift management into tabs, with focus on payment management for waiter shifts.

## Architecture

### Component Structure
```
ShiftView.vue (modified)
├── Cashier View (NEW)
│   ├── Tab Navigation (Waiter Shifts / Barista Shifts)
│   ├── Date Filter
│   ├── Shift Selector Dropdown
│   ├── Shift Summary Card
│   │   ├── Waiter Shift Summary (financial stats)
│   │   └── Barista Shift Summary (order stats)
│   └── Payment List (Waiter shifts only)
│       └── Payment Card with Actions
└── Waiter/Barista View (UNCHANGED)
    └── Existing UI preserved
```

### State Management

#### Cashier State
```javascript
// Tab state
const activeTab = ref('waiter') // 'waiter' | 'barista'

// Filter state
const selectedDate = ref(new Date().toISOString().split('T')[0])
const selectedShiftId = ref('')

// Data state
const shifts = computed(() => shiftStore.shifts)
const filteredShifts = computed(() => /* filter by date and role_type */)
const selectedShift = computed(() => /* find shift by selectedShiftId */)
const shiftStatus = computed(() => cashierStore.shiftStatus)
const payments = computed(() => cashierStore.payments)

// Modal state
const showOverride = ref(false)
const showDiscrepancy = ref(false)
const showCloseShift = ref(false)
const selectedPayment = ref(null)
```

## UI Design

### Cashier View Layout

#### 1. Header (Fixed)
```
┌─────────────────────────────────────┐
│ ⏰ Ca làm việc                      │
└─────────────────────────────────────┘
```

#### 2. Tab Navigation
```
┌─────────────────────────────────────┐
│ [🍽️ Ca phục vụ] [🍹 Ca pha chế]    │
└─────────────────────────────────────┘
```

#### 3. Filters Section
```
┌─────────────────────────────────────┐
│ 📅 Chọn ngày                        │
│ [Date Picker: YYYY-MM-DD]          │
│                                     │
│ Chọn ca để xem chi tiết             │
│ [Dropdown: Shift List]             │
└─────────────────────────────────────┘
```

#### 4. Shift Summary Card (Waiter)
```
┌─────────────────────────────────────┐
│ 📊 Tổng quan ca phục vụ            │
│ Nguyễn Văn A • Ca sáng             │
│ 7:00 - 15:00 • 🟢 Đang mở          │
│                                     │
│ ┌─────────┬─────────┐              │
│ │ 15 đơn  │ 2.5M₫   │              │
│ │ Tổng đơn│ Doanh thu│              │
│ └─────────┴─────────┘              │
│ ┌─────────┬─────────┐              │
│ │ 💵 1.8M │ 💳 700K │              │
│ │ Tiền mặt│ CK/QR   │              │
│ └─────────┴─────────┘              │
│                                     │
│ [🔒 Chốt ca]                       │
└─────────────────────────────────────┘
```

#### 5. Shift Summary Card (Barista)
```
┌─────────────────────────────────────┐
│ 📊 Tổng quan ca pha chế            │
│ Trần Thị B • Ca chiều              │
│ 12:00 - 18:00 • 🟢 Đang mở         │
│                                     │
│ ┌─────────┬─────────┐              │
│ │ 25 đơn  │ 8 món   │              │
│ │ Tổng đơn│ Đang pha│              │
│ └─────────┴─────────┘              │
│                                     │
│ [🔒 Chốt ca]                       │
└─────────────────────────────────────┘
```

#### 6. Payment List (Waiter shifts only)
```
┌─────────────────────────────────────┐
│ 💳 Danh sách thanh toán (15)       │
│                                     │
│ ┌─────────────────────────────────┐│
│ │ Nguyễn Văn C                    ││
│ │ 14:30 • 150,000₫ • 💵 Tiền mặt ││
│ │ [✓ Đã thu]                      ││
│ │                                 ││
│ │ [✏️ Điều chỉnh] [⚠️ Báo lỗi]   ││
│ │ [🔒 Khóa]                       ││
│ └─────────────────────────────────┘│
│                                     │
│ ┌─────────────────────────────────┐│
│ │ Khách lẻ                        ││
│ │ 14:25 • 85,000₫ • 💳 CK        ││
│ │ [✓ Đã thu]                      ││
│ │                                 ││
│ │ [✏️ Điều chỉnh] [⚠️ Báo lỗi]   ││
│ │ [🔒 Khóa]                       ││
│ └─────────────────────────────────┘│
└─────────────────────────────────────┘
```

### Color Scheme

#### Gradient Cards
- **Waiter Shift**: `from-yellow-500 to-orange-500` (warm, money-related)
- **Barista Shift**: `from-purple-500 to-indigo-500` (cool, production-related)

#### Status Colors
- **Open Shift**: Green (`bg-green-100 text-green-800`)
- **Closed Shift**: Gray (`bg-gray-100 text-gray-800`)
- **Payment Method - Cash**: Green (`bg-green-100 text-green-700`)
- **Payment Method - Transfer**: Blue (`bg-blue-100 text-blue-700`)
- **Payment Method - QR**: Purple (`bg-purple-100 text-purple-700`)

#### Action Buttons
- **Override**: Orange (`bg-orange-50 text-orange-600 border-orange-200`)
- **Discrepancy**: Yellow (`bg-yellow-50 text-yellow-600 border-yellow-200`)
- **Lock**: Red (`bg-red-50 text-red-600 border-red-200`)
- **Close Shift**: Red (`bg-red-500 text-white`)

## Data Flow

### Initial Load (Cashier)
```
1. Component mounts
2. Check user role → if Cashier, show new UI
3. Fetch all shifts: GET /api/shifts
4. Set default date to today
5. Filter shifts by date and activeTab (waiter/barista)
6. Wait for user to select a shift
```

### Shift Selection Flow
```
1. User selects shift from dropdown
2. Store selectedShiftId
3. Fetch shift status: GET /api/cashier/shifts/:id/status
4. If waiter shift: Fetch payments: GET /api/cashier/shifts/:id/payments
5. Display shift summary card
6. Display payment list (waiter only)
```

### Payment Action Flow

#### Override Payment
```
1. User clicks "Điều chỉnh" on payment
2. Show OverridePaymentModal
3. User enters reason
4. POST /api/cashier/payments/:orderId/override
5. Refresh payments list
```

#### Report Discrepancy
```
1. User clicks "Báo lỗi" on payment
2. Show DiscrepancyModal
3. User enters amount and reason
4. POST /api/cashier/discrepancies
5. Refresh payments list
```

#### Lock Order
```
1. User clicks "Khóa" on payment
2. Show confirmation dialog
3. PATCH /api/cashier/orders/:orderId/lock
4. Refresh payments list
```

### Close Shift Flow
```
1. User clicks "Chốt ca" on open shift
2. Show CloseShiftModal with end_cash input
3. User enters end_cash
4. POST /api/manager/shifts/:id/close
5. Refresh shifts list
6. Clear selectedShiftId
```

## API Integration

### Endpoints

#### Get All Shifts
```javascript
GET /api/shifts
Response: Array<Shift>
```

#### Get Shift Status
```javascript
GET /api/cashier/shifts/:id/status
Response: {
  shift_id: string
  total_orders: number
  total_revenue: number
  cash_revenue: number
  transfer_revenue: number
  qr_revenue: number
}
```

#### Get Shift Payments
```javascript
GET /api/cashier/shifts/:id/payments
Response: Array<Payment>
```

#### Override Payment
```javascript
PATCH /api/cashier/payments/:orderId/override
Body: { reason: string }
```

#### Report Discrepancy
```javascript
POST /api/cashier/discrepancies
Body: {
  order_id: string
  amount: number
  reason: string
}
```

#### Lock Order
```javascript
PATCH /api/cashier/orders/:orderId/lock
```

#### Close Shift
```javascript
POST /api/manager/shifts/:id/close
Body: { end_cash: number }
```

## Component Logic

### Computed Properties

```javascript
// Filter shifts by date and role type
const filteredShifts = computed(() => {
  if (!isCashier.value) return shifts.value
  
  return shifts.value.filter(shift => {
    // Filter by date
    const shiftDate = new Date(shift.started_at).toISOString().split('T')[0]
    if (shiftDate !== selectedDate.value) return false
    
    // Filter by role type based on active tab
    if (activeTab.value === 'waiter') {
      return shift.role_type === 'waiter'
    } else {
      return shift.role_type === 'barista'
    }
  })
})

// Get selected shift object
const selectedShift = computed(() => {
  if (!selectedShiftId.value) return null
  return shifts.value.find(s => s.id === selectedShiftId.value)
})

// Check if selected shift is open
const isShiftOpen = computed(() => {
  return selectedShift.value?.status === SHIFT_STATUS.OPEN
})

// Check if selected shift is waiter shift
const isWaiterShift = computed(() => {
  return selectedShift.value?.role_type === 'waiter'
})
```

### Methods

```javascript
// Load shift details when selected
const loadShiftDetails = async () => {
  if (!selectedShiftId.value) return
  
  // Fetch shift status
  await cashierStore.getShiftStatus(selectedShiftId.value)
  
  // Fetch payments if waiter shift
  if (isWaiterShift.value) {
    await cashierStore.getPaymentsByShift(selectedShiftId.value)
  }
}

// Handle tab change
const onTabChange = (tab) => {
  activeTab.value = tab
  selectedShiftId.value = '' // Clear selection
}

// Handle date change
const onDateChange = () => {
  selectedShiftId.value = '' // Clear selection
}

// Handle shift selection
const onShiftSelect = async () => {
  await loadShiftDetails()
}

// Close shift
const closeShift = async (shiftId, endCash) => {
  await shiftStore.closeShift(shiftId, endCash)
  selectedShiftId.value = ''
  await shiftStore.fetchAllShifts()
}
```

## Responsive Design

### Mobile Layout
- Full width cards
- Stack payment actions vertically on small screens
- Touch-friendly button sizes (min 44px height)
- Proper spacing for thumb navigation

### Tablet/Desktop
- Same layout as mobile (mobile-first)
- Slightly wider cards with max-width constraint
- Better use of horizontal space for payment actions

## Accessibility

### Keyboard Navigation
- Tab through all interactive elements
- Enter/Space to activate buttons
- Escape to close modals

### Screen Readers
- Proper ARIA labels for all buttons
- Semantic HTML structure
- Status announcements for actions

### Color Contrast
- All text meets WCAG AA standards
- Status indicators have sufficient contrast
- Focus indicators visible on all interactive elements

## Performance Considerations

### Data Loading
- Fetch shifts once on mount
- Fetch shift details only when selected
- Cache shift status to avoid repeated calls

### Rendering
- Use v-if for conditional rendering (not v-show)
- Lazy load payment list
- Virtual scrolling for large payment lists (if needed)

### State Updates
- Debounce date filter changes
- Batch state updates where possible
- Use computed properties for derived state

## Error Handling

### API Errors
- Show error toast for failed requests
- Provide retry mechanism
- Clear error state after successful retry

### Validation Errors
- Validate end_cash before closing shift
- Show inline validation errors
- Prevent submission with invalid data

### Network Errors
- Show offline indicator
- Queue actions for retry when online
- Inform user of pending actions

## Testing Strategy

### Unit Tests
- Test computed properties
- Test filter logic
- Test state management

### Integration Tests
- Test API integration
- Test data flow
- Test error handling

### E2E Tests
- Test complete user flows
- Test tab switching
- Test shift selection
- Test payment actions
- Test shift closing

## Migration Strategy

### Phase 1: Add Cashier View
- Add new UI for cashier role
- Keep existing UI for waiter/barista
- Test cashier view thoroughly

### Phase 2: Deploy
- Deploy to production
- Monitor for issues
- Gather user feedback

### Phase 3: Iterate
- Address feedback
- Optimize performance
- Add enhancements

## Future Enhancements

### Potential Features
- Export shift report to PDF
- Bulk payment actions
- Advanced filtering (by payment method, status)
- Real-time updates via WebSocket
- Shift comparison view
- Analytics dashboard

### Performance Improvements
- Virtual scrolling for large lists
- Infinite scroll for shift history
- Optimistic UI updates
- Background data refresh
