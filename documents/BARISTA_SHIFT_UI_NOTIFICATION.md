# Barista Shift UI Notification

## Overview

Implemented UI notifications to inform barista when they need to open a shift before accepting orders.

## Features

### 1. Warning Banner

**Location**: Top of BaristaView, above tabs

**Appearance**:
- Gradient background: orange to red
- Warning icon: ⚠️
- Clear message: "Chưa mở ca làm việc"
- Call-to-action button: "Mở ca ngay →"

**Behavior**:
- Shows when `hasOpenShift = false`
- Hides when shift is opened
- Button redirects to `/shifts` page

**Code**:
```vue
<div v-if="!hasOpenShift" class="mb-4 bg-gradient-to-r from-orange-500 to-red-500 text-white rounded-2xl p-4 shadow-lg">
  <div class="flex items-start gap-3">
    <div class="text-3xl">⚠️</div>
    <div class="flex-1">
      <h3 class="font-bold text-lg mb-1">Chưa mở ca làm việc</h3>
      <p class="text-sm opacity-90 mb-3">Bạn cần mở ca trước khi nhận order từ queue</p>
      <button @click="$router.push('/shifts')"
        class="bg-white text-orange-600 px-4 py-2 rounded-lg font-medium text-sm active:scale-95 transition-transform">
        Mở ca ngay →
      </button>
    </div>
  </div>
</div>
```

### 2. Disabled Button State

**Location**: "Nhận order" button on each queued order

**Appearance**:
- **Without shift**: Gray background, disabled cursor, text "🔒 Cần mở ca"
- **With shift**: Blue background, active cursor, text "👍 Nhận order"

**Behavior**:
- Button is disabled when `hasOpenShift = false`
- Button is enabled when `hasOpenShift = true`
- Visual feedback prevents confusion

**Code**:
```vue
<button @click="acceptOrder(order.id)"
  :disabled="!hasOpenShift"
  :class="[
    'w-full py-3 rounded-xl font-bold transition-all',
    hasOpenShift 
      ? 'bg-blue-500 text-white active:scale-95' 
      : 'bg-gray-300 text-gray-500 cursor-not-allowed'
  ]">
  {{ hasOpenShift ? '👍 Nhận order' : '🔒 Cần mở ca' }}
</button>
```

### 3. Pre-Accept Validation

**Location**: `acceptOrder()` function

**Behavior**:
- Checks shift status before API call
- Shows confirmation dialog if no shift
- Offers to redirect to shift page

**Code**:
```javascript
const acceptOrder = async (id) => {
  // Check shift before accepting
  if (!hasOpenShift.value) {
    if (confirm('Bạn chưa mở ca làm việc. Bạn có muốn mở ca ngay không?')) {
      router.push('/shifts')
    }
    return
  }

  try {
    await baristaStore.acceptOrder(id)
    activeTab.value = 'working'
  } catch (error) {
    // Handle error...
  }
}
```

### 4. Error Handling

**Location**: `acceptOrder()` catch block

**Behavior**:
- Catches API error from backend
- Detects shift-related errors
- Shows user-friendly message
- Offers to redirect to shift page

**Code**:
```javascript
catch (error) {
  const errorMsg = error.response?.data?.error || error.message
  
  // Handle specific error for shift requirement
  if (errorMsg.includes('shift')) {
    alert('⚠️ Bạn phải mở ca trước khi nhận order.\n\nVui lòng vào "Ca làm việc" để mở ca.')
    if (confirm('Chuyển đến trang Ca làm việc?')) {
      router.push('/shifts')
    }
  } else {
    alert('Lỗi: ' + errorMsg)
  }
}
```

### 5. Auto-Refresh Shift Status

**Location**: `refreshAll()` function

**Behavior**:
- Fetches current shift status every 10 seconds
- Updates `hasOpenShift` computed property
- Banner appears/disappears automatically

**Code**:
```javascript
const refreshAll = async () => {
  await Promise.all([
    baristaStore.fetchQueuedOrders(),
    baristaStore.fetchMyOrders(),
    shiftStore.fetchCurrentShift()  // ← Added
  ])
}

onMounted(async () => {
  await refreshAll()
  refreshInterval = setInterval(refreshAll, 10000)
})
```

## User Flow

### Scenario 1: Barista Without Shift

1. **Login** → Barista logs in
2. **Navigate to Barista tab** → See warning banner
3. **See disabled buttons** → All "Nhận order" buttons are gray and disabled
4. **Click "Mở ca ngay"** → Redirected to `/shifts`
5. **Open shift** → Fill form and submit
6. **Return to Barista tab** → Warning banner disappears, buttons enabled

### Scenario 2: Try to Accept Without Shift

1. **Barista sees queue** → Warning banner visible
2. **Clicks disabled button** → Nothing happens (button disabled)
3. **Somehow bypasses UI** → Backend rejects with 400 error
4. **Error dialog shows** → "Bạn phải mở ca trước khi nhận order"
5. **Confirm dialog** → "Chuyển đến trang Ca làm việc?"
6. **Click OK** → Redirected to `/shifts`

### Scenario 3: Shift Already Open

1. **Login** → Barista logs in
2. **Navigate to Barista tab** → No warning banner
3. **See enabled buttons** → All "Nhận order" buttons are blue
4. **Click "Nhận order"** → Order accepted successfully
5. **Switch to "Đang pha" tab** → See accepted order

## Visual Design

### Warning Banner
```
┌─────────────────────────────────────────┐
│ ⚠️  Chưa mở ca làm việc                 │
│                                          │
│     Bạn cần mở ca trước khi nhận order  │
│     từ queue                             │
│                                          │
│     [ Mở ca ngay → ]                     │
└─────────────────────────────────────────┘
```

### Button States
```
Without Shift:
┌─────────────────────────┐
│   🔒 Cần mở ca          │  (Gray, disabled)
└─────────────────────────┘

With Shift:
┌─────────────────────────┐
│   👍 Nhận order         │  (Blue, active)
└─────────────────────────┘
```

## State Management

### Shift Store Integration

```javascript
import { useShiftStore } from '../stores/shift'

const shiftStore = useShiftStore()
const hasOpenShift = computed(() => shiftStore.hasOpenShift)
```

### Computed Property

```javascript
// In shift store
hasOpenShift: (state) => {
  return state.currentShift && state.currentShift.status === 'OPEN'
}
```

## Testing

### Manual Test Cases

#### Test 1: Warning Banner Visibility
- [ ] Login as barista without shift
- [ ] Navigate to Barista tab
- [ ] Verify warning banner is visible
- [ ] Open shift
- [ ] Verify warning banner disappears

#### Test 2: Button State
- [ ] Login as barista without shift
- [ ] Navigate to Barista tab
- [ ] Verify all "Nhận order" buttons show "🔒 Cần mở ca"
- [ ] Verify buttons are gray and disabled
- [ ] Open shift
- [ ] Verify buttons change to "👍 Nhận order"
- [ ] Verify buttons are blue and enabled

#### Test 3: Pre-Accept Validation
- [ ] Login as barista without shift
- [ ] Try to click disabled button (should do nothing)
- [ ] If somehow able to trigger, verify confirmation dialog
- [ ] Click OK in dialog
- [ ] Verify redirected to `/shifts`

#### Test 4: Error Handling
- [ ] Login as barista without shift
- [ ] Use browser console to bypass UI and call API
- [ ] Verify error alert shows
- [ ] Verify redirect offer appears

#### Test 5: Auto-Refresh
- [ ] Login as barista without shift
- [ ] Keep Barista tab open
- [ ] In another tab, open shift
- [ ] Wait 10 seconds
- [ ] Verify warning banner disappears automatically

### Automated Tests (Future)

```javascript
// BaristaView.spec.js
describe('BaristaView Shift Notification', () => {
  it('shows warning banner when no shift', () => {
    // Mock hasOpenShift = false
    // Render component
    // Assert banner is visible
  })

  it('hides warning banner when shift open', () => {
    // Mock hasOpenShift = true
    // Render component
    // Assert banner is not visible
  })

  it('disables accept button when no shift', () => {
    // Mock hasOpenShift = false
    // Render component
    // Assert button is disabled
  })

  it('redirects to shifts page when clicking banner button', () => {
    // Mock router
    // Click banner button
    // Assert router.push('/shifts') called
  })
})
```

## Accessibility

- ✅ Clear visual indicators (color, icon, text)
- ✅ Disabled state prevents accidental clicks
- ✅ Confirmation dialogs for important actions
- ✅ Descriptive button text
- ✅ Touch-friendly button sizes (44px minimum)

## Performance

- ✅ Computed properties for reactive updates
- ✅ Minimal re-renders (only when shift status changes)
- ✅ Auto-refresh every 10 seconds (not too frequent)
- ✅ No unnecessary API calls

## Browser Compatibility

- ✅ Modern browsers (Chrome, Firefox, Safari, Edge)
- ✅ Mobile browsers (iOS Safari, Chrome Mobile)
- ✅ Responsive design (mobile-first)

## Future Enhancements

1. **Toast Notifications**: Replace alerts with toast messages
2. **Shift Status Indicator**: Show shift info in header
3. **Quick Open Shift**: Modal to open shift without leaving page
4. **Shift Reminder**: Notification on login if no shift
5. **Sound Alert**: Audio notification for shift requirement

## Related Files

- `frontend/src/views/BaristaView.vue` - Main implementation
- `frontend/src/stores/shift.js` - Shift state management
- `frontend/src/stores/barista.js` - Barista actions
- `backend/application/services/order_service.go` - Backend validation
- `documents/BR13_BARISTA_SHIFT_VALIDATION.md` - Business rule documentation

## Screenshots

### Before Opening Shift
```
┌─────────────────────────────────────────┐
│ 🍹 Barista                        🔄    │
├─────────────────────────────────────────┤
│ ⏳ Queue (2) │ 🍹 Đang pha (0) │ ✅ ... │
├─────────────────────────────────────────┤
│                                          │
│ ⚠️  Chưa mở ca làm việc                 │
│     Bạn cần mở ca trước khi nhận order  │
│     [ Mở ca ngay → ]                     │
│                                          │
├─────────────────────────────────────────┤
│ ORD-001                    ⏳ Chờ pha   │
│ Khách lẻ                                 │
│ ☕ Coffee x1                             │
│ [ 🔒 Cần mở ca ]                         │
└─────────────────────────────────────────┘
```

### After Opening Shift
```
┌─────────────────────────────────────────┐
│ 🍹 Barista                        🔄    │
├─────────────────────────────────────────┤
│ ⏳ Queue (2) │ 🍹 Đang pha (0) │ ✅ ... │
├─────────────────────────────────────────┤
│ ORD-001                    ⏳ Chờ pha   │
│ Khách lẻ                                 │
│ ☕ Coffee x1                             │
│ [ 👍 Nhận order ]                        │
└─────────────────────────────────────────┘
```

## Summary

The UI notification system provides:
- ✅ Clear visual warning when shift not open
- ✅ Disabled buttons to prevent errors
- ✅ Easy navigation to shift page
- ✅ Graceful error handling
- ✅ Auto-refresh for real-time updates
- ✅ User-friendly messages in Vietnamese
- ✅ Mobile-optimized design

This ensures baristas always know they need to open a shift before accepting orders, preventing confusion and errors.
