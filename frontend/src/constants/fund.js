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
  { value: 'handover', label: 'Bàn giao' }
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
