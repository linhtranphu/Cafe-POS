// Journal event types (double-entry source of truth)
export const EVENT_TYPES = {
  CASHIER_SHIFT_START: 'cashier_shift_start',
  CASHIER_SHIFT_END: 'cashier_shift_end',
  WAITER_SHIFT_START: 'waiter_shift_start',
  WAITER_HANDOVER: 'waiter_handover',
  FUND_TRANSFER: 'fund_transfer',
  MANAGER_DEPOSIT: 'manager_deposit',
  MANAGER_WITHDRAWAL: 'manager_withdrawal',
  EXPENSE: 'expense',
  INGREDIENT_RESTOCK: 'ingredient_restock',
  FACILITY_PURCHASE: 'facility_purchase',
  ALL: 'all'
}

export const EVENT_TYPE_LABELS = {
  cashier_shift_start: 'Đầu ca thu ngân',
  cashier_shift_end: 'Cuối ca thu ngân',
  waiter_shift_start: 'Đầu ca phục vụ',
  waiter_handover: 'Bàn giao phục vụ',
  fund_transfer: 'Chuyển quỹ',
  manager_deposit: 'Nạp tiền',
  manager_withdrawal: 'Rút tiền',
  expense: 'Chi tiêu',
  ingredient_restock: 'Mua nguyên liệu',
  facility_purchase: 'Mua tài sản',
  all: 'Tất cả sự kiện'
}

export const EVENT_TYPE_ICONS = {
  cashier_shift_start: '🏦',
  cashier_shift_end: '🏦',
  waiter_shift_start: '👤',
  waiter_handover: '🤝',
  fund_transfer: '↔️',
  manager_deposit: '📥',
  manager_withdrawal: '📤',
  expense: '🧾',
  ingredient_restock: '🥦',
  facility_purchase: '🏢'
}

// Events that represent money flowing IN to a real fund
export const INFLOW_EVENTS = new Set([
  'manager_deposit',
  'cashier_shift_end',
  'waiter_handover'
])

// Events that represent money flowing OUT of a real fund
export const OUTFLOW_EVENTS = new Set([
  'manager_withdrawal',
  'cashier_shift_start',
  'waiter_shift_start',
  'expense',
  'ingredient_restock',
  'facility_purchase'
])

export const EVENT_TYPE_FILTER_OPTIONS = [
  { value: 'all', label: 'Tất cả sự kiện' },
  { value: 'manager_deposit', label: 'Nạp tiền' },
  { value: 'manager_withdrawal', label: 'Rút tiền' },
  { value: 'fund_transfer', label: 'Chuyển quỹ' },
  { value: 'expense', label: 'Chi tiêu' },
  { value: 'ingredient_restock', label: 'Mua nguyên liệu' },
  { value: 'facility_purchase', label: 'Mua tài sản' },
  { value: 'cashier_shift_start', label: 'Đầu ca thu ngân' },
  { value: 'cashier_shift_end', label: 'Cuối ca thu ngân' },
  { value: 'waiter_shift_start', label: 'Đầu ca phục vụ' },
  { value: 'waiter_handover', label: 'Bàn giao phục vụ' }
]

export const getEventTypeLabel = (type) => EVENT_TYPE_LABELS[type] || type
export const getEventTypeIcon = (type) => EVENT_TYPE_ICONS[type] || '📋'

// Fund types
export const FUND_TYPES = {
  OPERATING: 'operating',
  INVENTORY: 'inventory',
  PROFIT: 'profit',
  CASH_DRAWER: 'cash_drawer',
  WAITER_FLOAT: 'waiter_float',
  // External counterpart accounts (không phải quỹ thực)
  OWNER: 'owner',
  SUPPLIER: 'supplier',
  CUSTOMER: 'customer',
  CASH_SHORTAGE: 'cash_shortage',
  CASH_OVERAGE: 'cash_overage',
  EXTERNAL: 'external', // legacy fund type from old data
  ALL: 'all'
}

// Fund type labels (Vietnamese)
export const FUND_TYPE_LABELS = {
  [FUND_TYPES.OPERATING]: 'Quỹ vận hành',
  [FUND_TYPES.INVENTORY]: 'Quỹ hàng hóa',
  [FUND_TYPES.PROFIT]: 'Quỹ lợi nhuận',
  [FUND_TYPES.CASH_DRAWER]: 'Ngăn kéo tiền',
  [FUND_TYPES.WAITER_FLOAT]: 'Tiền phục vụ',
  // External counterparts
  [FUND_TYPES.OWNER]: 'Chủ quán',
  [FUND_TYPES.SUPPLIER]: 'Nhà cung cấp',
  [FUND_TYPES.CUSTOMER]: 'Khách hàng',
  [FUND_TYPES.CASH_SHORTAGE]: 'Thiếu tiền',
  [FUND_TYPES.CASH_OVERAGE]: 'Thừa tiền',
  [FUND_TYPES.EXTERNAL]: 'Ngoài (cũ)',
  [FUND_TYPES.ALL]: 'Tất cả quỹ'
}

// Fund type icons
export const FUND_TYPE_ICONS = {
  [FUND_TYPES.OPERATING]: '⚙️',
  [FUND_TYPES.INVENTORY]: '📦',
  [FUND_TYPES.PROFIT]: '💹',
  [FUND_TYPES.CASH_DRAWER]: '🗄️',
  [FUND_TYPES.WAITER_FLOAT]: '👤',
  [FUND_TYPES.OWNER]: '👑',
  [FUND_TYPES.SUPPLIER]: '🏭',
  [FUND_TYPES.CUSTOMER]: '🛒',
  [FUND_TYPES.CASH_SHORTAGE]: '📉',
  [FUND_TYPES.CASH_OVERAGE]: '📈',
  [FUND_TYPES.EXTERNAL]: '💰',
}

// Fund type colors (Tailwind classes)
export const FUND_TYPE_COLORS = {
  [FUND_TYPES.OPERATING]: 'blue',
  [FUND_TYPES.INVENTORY]: 'green',
  [FUND_TYPES.PROFIT]: 'yellow',
  [FUND_TYPES.CASH_DRAWER]: 'purple',
  [FUND_TYPES.WAITER_FLOAT]: 'orange',
  [FUND_TYPES.OWNER]: 'rose',
  [FUND_TYPES.SUPPLIER]: 'teal',
  [FUND_TYPES.CUSTOMER]: 'cyan',
  [FUND_TYPES.CASH_SHORTAGE]: 'red',
  [FUND_TYPES.CASH_OVERAGE]: 'emerald'
}

// Source types
export const SOURCE_TYPES = {
  EXPENSE: 'expense',
  INGREDIENT: 'ingredient',
  FACILITY: 'facility',
  HANDOVER: 'handover',
  MANUAL: 'manual',
  FUND_TRANSFER: 'fund_transfer',
  WAITER_START: 'waiter_start',
  ALL: 'all'
}

// Source type labels (Vietnamese)
export const SOURCE_TYPE_LABELS = {
  [SOURCE_TYPES.EXPENSE]: 'Chi phí',
  [SOURCE_TYPES.INGREDIENT]: 'Nguyên liệu',
  [SOURCE_TYPES.FACILITY]: 'Cơ sở vật chất',
  [SOURCE_TYPES.HANDOVER]: 'Bàn giao ca',
  [SOURCE_TYPES.MANUAL]: 'Thủ công',
  [SOURCE_TYPES.FUND_TRANSFER]: 'Chuyển quỹ',
  [SOURCE_TYPES.WAITER_START]: 'Tiền đầu ca phục vụ',
  [SOURCE_TYPES.ALL]: 'Tất cả nguồn'
}

// Source type icons
export const SOURCE_TYPE_ICONS = {
  [SOURCE_TYPES.EXPENSE]: '🧾',
  [SOURCE_TYPES.INGREDIENT]: '🥦',
  [SOURCE_TYPES.FACILITY]: '🏢',
  [SOURCE_TYPES.HANDOVER]: '🤝',
  [SOURCE_TYPES.MANUAL]: '✋',
  [SOURCE_TYPES.FUND_TRANSFER]: '↔️',
  [SOURCE_TYPES.WAITER_START]: '👤'
}

// Fund type filter options
export const FUND_TYPE_FILTER_OPTIONS = [
  { value: FUND_TYPES.ALL, label: FUND_TYPE_LABELS[FUND_TYPES.ALL] },
  { value: FUND_TYPES.OPERATING, label: FUND_TYPE_LABELS[FUND_TYPES.OPERATING] },
  { value: FUND_TYPES.INVENTORY, label: FUND_TYPE_LABELS[FUND_TYPES.INVENTORY] },
  { value: FUND_TYPES.PROFIT, label: FUND_TYPE_LABELS[FUND_TYPES.PROFIT] },
  { value: FUND_TYPES.CASH_DRAWER, label: FUND_TYPE_LABELS[FUND_TYPES.CASH_DRAWER] },
  { value: FUND_TYPES.WAITER_FLOAT, label: FUND_TYPE_LABELS[FUND_TYPES.WAITER_FLOAT] },
  { value: FUND_TYPES.CASH_SHORTAGE, label: '📉 ' + FUND_TYPE_LABELS[FUND_TYPES.CASH_SHORTAGE] },
  { value: FUND_TYPES.CASH_OVERAGE, label: '📈 ' + FUND_TYPE_LABELS[FUND_TYPES.CASH_OVERAGE] }
]

// Source type filter options
export const SOURCE_TYPE_FILTER_OPTIONS = [
  { value: SOURCE_TYPES.ALL, label: SOURCE_TYPE_LABELS[SOURCE_TYPES.ALL] },
  { value: SOURCE_TYPES.EXPENSE, label: SOURCE_TYPE_LABELS[SOURCE_TYPES.EXPENSE] },
  { value: SOURCE_TYPES.INGREDIENT, label: SOURCE_TYPE_LABELS[SOURCE_TYPES.INGREDIENT] },
  { value: SOURCE_TYPES.FACILITY, label: SOURCE_TYPE_LABELS[SOURCE_TYPES.FACILITY] },
  { value: SOURCE_TYPES.HANDOVER, label: SOURCE_TYPE_LABELS[SOURCE_TYPES.HANDOVER] },
  { value: SOURCE_TYPES.MANUAL, label: SOURCE_TYPE_LABELS[SOURCE_TYPES.MANUAL] },
  { value: SOURCE_TYPES.FUND_TRANSFER, label: SOURCE_TYPE_LABELS[SOURCE_TYPES.FUND_TRANSFER] },
  { value: SOURCE_TYPES.WAITER_START, label: SOURCE_TYPE_LABELS[SOURCE_TYPES.WAITER_START] }
]

// Helper functions
export const getFundTypeLabel = (type) => FUND_TYPE_LABELS[type] || type
export const getFundTypeIcon = (type) => FUND_TYPE_ICONS[type] || '💰'
export const getFundTypeColor = (type) => FUND_TYPE_COLORS[type] || 'gray'
export const getSourceTypeLabel = (type) => SOURCE_TYPE_LABELS[type] || type
export const getSourceTypeIcon = (type) => SOURCE_TYPE_ICONS[type] || '📋'

// Transaction types
export const TRANSACTION_TYPES = {
  DEPOSIT: 'deposit',
  WITHDRAWAL: 'withdrawal',
  FUND_HANDOVER: 'fund_handover',
  CASH_HANDOVER: 'cash_handover',
  STARTING_FLOAT: 'starting_float',
  ALL: 'all'
}

// Money types
export const MONEY_TYPES = {
  CASH: 'cash',
  TRANSFER: 'transfer',
  BOTH: 'both',
  ALL: 'all'
}

// Transaction type labels (Vietnamese)
export const TRANSACTION_TYPE_LABELS = {
  [TRANSACTION_TYPES.DEPOSIT]: 'Thêm tiền',
  [TRANSACTION_TYPES.WITHDRAWAL]: 'Rút tiền',
  [TRANSACTION_TYPES.FUND_HANDOVER]: 'Bàn giao quỹ',
  [TRANSACTION_TYPES.CASH_HANDOVER]: 'Bàn giao tiền',
  [TRANSACTION_TYPES.STARTING_FLOAT]: 'Tiền đầu ca',
  [TRANSACTION_TYPES.ALL]: 'Tất cả'
}

// Transaction type icons
export const TRANSACTION_TYPE_ICONS = {
  [TRANSACTION_TYPES.DEPOSIT]: '📥',
  [TRANSACTION_TYPES.WITHDRAWAL]: '📤',
  [TRANSACTION_TYPES.FUND_HANDOVER]: '🔄',
  [TRANSACTION_TYPES.CASH_HANDOVER]: '🔄',
  [TRANSACTION_TYPES.STARTING_FLOAT]: '💰'
}

// Money type labels (Vietnamese)
export const MONEY_TYPE_LABELS = {
  [MONEY_TYPES.CASH]: 'Tiền mặt',
  [MONEY_TYPES.TRANSFER]: 'Chuyển khoản',
  [MONEY_TYPES.BOTH]: 'Cả hai',
  [MONEY_TYPES.ALL]: 'Tất cả'
}

// Money type icons
export const MONEY_TYPE_ICONS = {
  [MONEY_TYPES.CASH]: '💵',
  [MONEY_TYPES.TRANSFER]: '💳',
  [MONEY_TYPES.BOTH]: '💰'
}

// Filter options for transaction type
export const TRANSACTION_TYPE_FILTER_OPTIONS = [
  { value: TRANSACTION_TYPES.ALL, label: TRANSACTION_TYPE_LABELS[TRANSACTION_TYPES.ALL] },
  { value: TRANSACTION_TYPES.DEPOSIT, label: TRANSACTION_TYPE_LABELS[TRANSACTION_TYPES.DEPOSIT] },
  { value: TRANSACTION_TYPES.WITHDRAWAL, label: TRANSACTION_TYPE_LABELS[TRANSACTION_TYPES.WITHDRAWAL] },
  { value: 'handover', label: 'Bàn giao' },
  { value: 'fund_transfer', label: 'Chuyển quỹ' }
]

// Filter options for money type
export const MONEY_TYPE_FILTER_OPTIONS = [
  { value: MONEY_TYPES.ALL, label: MONEY_TYPE_LABELS[MONEY_TYPES.ALL] },
  { value: MONEY_TYPES.CASH, label: MONEY_TYPE_LABELS[MONEY_TYPES.CASH] },
  { value: MONEY_TYPES.TRANSFER, label: MONEY_TYPE_LABELS[MONEY_TYPES.TRANSFER] }
]

// Date filter options
export const DATE_FILTER_OPTIONS = [
  { value: 'today', label: 'Hôm nay' },
  { value: 'yesterday', label: 'Hôm qua' },
  { value: 'this_week', label: 'Tuần này' },
  { value: 'this_month', label: 'Tháng này' },
  { value: 'custom', label: 'Tùy chỉnh' }
]

// Role labels (Vietnamese)
export const ROLE_LABELS = {
  manager: 'Quản lý',
  cashier: 'Thu ngân',
  waiter: 'Phục vụ',
  barista: 'Pha chế'
}

// Validation rules
export const VALIDATION = {
  MIN_REASON_LENGTH: 10,
  MIN_AMOUNT: 0,
  MAX_AMOUNT: 999999999
}

// Pagination
export const PAGINATION = {
  DEFAULT_LIMIT: 20,
  MAX_LIMIT: 200,
  DEFAULT_OFFSET: 0
}

// Helper functions
export const getTransactionTypeLabel = (type) => {
  return TRANSACTION_TYPE_LABELS[type] || type
}

export const getTransactionTypeIcon = (type) => {
  return TRANSACTION_TYPE_ICONS[type] || '📋'
}

export const getMoneyTypeLabel = (type) => {
  return MONEY_TYPE_LABELS[type] || type
}

export const getMoneyTypeIcon = (type) => {
  return MONEY_TYPE_ICONS[type] || '💰'
}

export const getRoleLabel = (role) => {
  return ROLE_LABELS[role] || role
}

export const isInflowTransaction = (type) => {
  return [
    TRANSACTION_TYPES.DEPOSIT,
    TRANSACTION_TYPES.FUND_HANDOVER,
    TRANSACTION_TYPES.CASH_HANDOVER
  ].includes(type)
}

export const isOutflowTransaction = (type) => {
  return [
    TRANSACTION_TYPES.WITHDRAWAL,
    TRANSACTION_TYPES.STARTING_FLOAT
  ].includes(type)
}

export const getAmountColor = (type) => {
  return isInflowTransaction(type) ? 'text-green-600' : 'text-red-600'
}

export const getAmountPrefix = (type) => {
  return isInflowTransaction(type) ? '+' : '-'
}

// Date helper functions
export const getDateRange = (filterValue) => {
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  
  switch (filterValue) {
    case 'today':
      return {
        from_date: today.toISOString(),
        to_date: now.toISOString()
      }
    
    case 'yesterday':
      const yesterday = new Date(today)
      yesterday.setDate(yesterday.getDate() - 1)
      const yesterdayEnd = new Date(yesterday)
      yesterdayEnd.setHours(23, 59, 59, 999)
      return {
        from_date: yesterday.toISOString(),
        to_date: yesterdayEnd.toISOString()
      }
    
    case 'this_week':
      const weekStart = new Date(today)
      weekStart.setDate(weekStart.getDate() - weekStart.getDay())
      return {
        from_date: weekStart.toISOString(),
        to_date: now.toISOString()
      }
    
    case 'this_month':
      const monthStart = new Date(now.getFullYear(), now.getMonth(), 1)
      return {
        from_date: monthStart.toISOString(),
        to_date: now.toISOString()
      }
    
    default:
      return {
        from_date: today.toISOString(),
        to_date: now.toISOString()
      }
  }
}

// Currency formatting
export const formatCurrency = (amount) => {
  return new Intl.NumberFormat('vi-VN').format(amount)
}

// Date formatting
export const formatDateTime = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return 'Vừa xong'
  if (diffMins < 60) return `${diffMins} phút trước`
  if (diffHours < 24) return `${diffHours} giờ trước`
  if (diffDays < 7) return `${diffDays} ngày trước`
  
  return date.toLocaleDateString('vi-VN', { 
    day: '2-digit', 
    month: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

export const formatFullDateTime = (timestamp) => {
  const date = new Date(timestamp)
  return date.toLocaleString('vi-VN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}
