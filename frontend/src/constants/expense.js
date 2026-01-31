// Expense Constants - Synced with backend/domain/expense/expense.go

// Payment Methods
export const PAYMENT_METHODS = {
  CASH: 'cash',
  BANK: 'bank',
  CARD: 'card'
}

export const PAYMENT_METHOD_OPTIONS = [
  { value: PAYMENT_METHODS.CASH, label: 'Tiền mặt' },
  { value: PAYMENT_METHODS.BANK, label: 'Chuyển khoản' },
  { value: PAYMENT_METHODS.CARD, label: 'Thẻ' }
]

// Recurring Expense Frequencies
export const RECURRING_FREQUENCIES = {
  DAILY: 'daily',
  WEEKLY: 'weekly',
  MONTHLY: 'monthly',
  QUARTERLY: 'quarterly',
  YEARLY: 'yearly'
}

export const RECURRING_FREQUENCY_OPTIONS = [
  { value: RECURRING_FREQUENCIES.DAILY, label: 'Hàng ngày' },
  { value: RECURRING_FREQUENCIES.WEEKLY, label: 'Hàng tuần' },
  { value: RECURRING_FREQUENCIES.MONTHLY, label: 'Hàng tháng' },
  { value: RECURRING_FREQUENCIES.QUARTERLY, label: 'Hàng quý' },
  { value: RECURRING_FREQUENCIES.YEARLY, label: 'Hàng năm' }
]

// Helper Functions
export const getPaymentMethodLabel = (method) => {
  const option = PAYMENT_METHOD_OPTIONS.find(opt => opt.value === method)
  return option ? option.label : method
}

export const getFrequencyLabel = (frequency) => {
  const option = RECURRING_FREQUENCY_OPTIONS.find(opt => opt.value === frequency)
  return option ? option.label : frequency
}
