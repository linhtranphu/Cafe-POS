// Shift Status Constants
// Must match backend: backend/domain/order/shift.go (ShiftStatus type)
export const SHIFT_STATUS = {
  OPEN: 'OPEN',
  CLOSED: 'CLOSED'
}

// Cashier Shift Status Constants
// Must match backend: backend/domain/cashier/cashier_shift.go (CashierShiftStatus type)
export const CASHIER_SHIFT_STATUS = {
  OPEN: 'OPEN',
  CLOSURE_INITIATED: 'CLOSURE_INITIATED',
  CLOSED: 'CLOSED'
}

// Shift Type Constants
// Must match backend: backend/domain/order/shift.go (ShiftType type)
export const SHIFT_TYPE = {
  MORNING: 'MORNING',
  AFTERNOON: 'AFTERNOON',
  EVENING: 'EVENING'
}

// Role Type Constants
// Must match backend: backend/domain/order/shift.go (RoleType type)
export const ROLE_TYPE = {
  WAITER: 'waiter',
  BARISTA: 'barista'
}

// Shift Status Display
export const SHIFT_STATUS_DISPLAY = {
  [SHIFT_STATUS.OPEN]: { 
    label: 'Đang mở', 
    icon: '🟢', 
    badge: 'bg-green-100 text-green-800' 
  },
  [SHIFT_STATUS.CLOSED]: { 
    label: 'Đã đóng', 
    icon: '🔴', 
    badge: 'bg-gray-100 text-gray-800' 
  }
}

// Cashier Shift Status Display
export const CASHIER_SHIFT_STATUS_DISPLAY = {
  [CASHIER_SHIFT_STATUS.OPEN]: { 
    label: 'Đang mở', 
    icon: '🟢', 
    badge: 'bg-green-100 text-green-800' 
  },
  [CASHIER_SHIFT_STATUS.CLOSURE_INITIATED]: { 
    label: 'Đang đóng', 
    icon: '🟡', 
    badge: 'bg-yellow-100 text-yellow-800' 
  },
  [CASHIER_SHIFT_STATUS.CLOSED]: { 
    label: 'Đã đóng', 
    icon: '🔴', 
    badge: 'bg-gray-100 text-gray-800' 
  }
}

// Shift Type Display
export const SHIFT_TYPE_DISPLAY = {
  [SHIFT_TYPE.MORNING]: { 
    label: 'Ca sáng', 
    icon: '☀️', 
    badge: 'bg-yellow-100 text-yellow-800' 
  },
  [SHIFT_TYPE.AFTERNOON]: { 
    label: 'Ca chiều', 
    icon: '🌤️', 
    badge: 'bg-orange-100 text-orange-800' 
  },
  [SHIFT_TYPE.EVENING]: { 
    label: 'Ca tối', 
    icon: '🌙', 
    badge: 'bg-blue-100 text-blue-800' 
  }
}

// Role Type Display
export const ROLE_TYPE_DISPLAY = {
  [ROLE_TYPE.WAITER]: { 
    label: 'Phục vụ', 
    icon: '👨‍🍳', 
    badge: 'bg-blue-100 text-blue-800' 
  },
  [ROLE_TYPE.BARISTA]: { 
    label: 'Pha chế', 
    icon: '☕', 
    badge: 'bg-purple-100 text-purple-800' 
  }
}
