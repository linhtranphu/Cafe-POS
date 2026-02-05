// Order Status Constants
// Must match backend: backend/domain/order/order.go
export const ORDER_STATUS = {
  CREATED: 'CREATED',
  PAID: 'PAID',
  QUEUED: 'QUEUED',
  IN_PROGRESS: 'IN_PROGRESS',
  READY: 'READY',
  SERVED: 'SERVED',
  CANCELLED: 'CANCELLED',
  LOCKED: 'LOCKED'
}

// Payment Method Constants
// Must match backend: backend/domain/order/order.go (PaymentMethod type)
export const PAYMENT_METHOD = {
  CASH: 'CASH',
  TRANSFER: 'TRANSFER',
  QR: 'QR'
}

// Payment Method Display
export const PAYMENT_METHOD_DISPLAY = [
  { value: PAYMENT_METHOD.CASH, label: 'Tiền mặt', icon: '💵' },
  { value: PAYMENT_METHOD.QR, label: 'QR', icon: '📱' },
  { value: PAYMENT_METHOD.TRANSFER, label: 'CK', icon: '🏦' }
]

// Order Status Display
export const ORDER_STATUS_DISPLAY = {
  [ORDER_STATUS.CREATED]: { label: 'Mới tạo', icon: '🆕', badge: 'bg-gray-100 text-gray-800' },
  [ORDER_STATUS.PAID]: { label: 'Đã thanh toán', icon: '💰', badge: 'bg-green-100 text-green-800' },
  [ORDER_STATUS.QUEUED]: { label: 'Chờ pha chế', icon: '⏳', badge: 'bg-yellow-100 text-yellow-800' },
  [ORDER_STATUS.IN_PROGRESS]: { label: 'Đang pha chế', icon: '🍹', badge: 'bg-blue-100 text-blue-800' },
  [ORDER_STATUS.READY]: { label: 'Sẵn sàng', icon: '✅', badge: 'bg-purple-100 text-purple-800' },
  [ORDER_STATUS.SERVED]: { label: 'Đã phục vụ', icon: '🎉', badge: 'bg-green-100 text-green-800' },
  [ORDER_STATUS.CANCELLED]: { label: 'Đã hủy', icon: '❌', badge: 'bg-red-100 text-red-800' },
  [ORDER_STATUS.LOCKED]: { label: 'Đã khóa', icon: '🔒', badge: 'bg-gray-100 text-gray-800' }
}

// Status Filter Options
export const STATUS_FILTER_OPTIONS = [
  { value: 'ALL', label: 'Tất cả', icon: '📋' },
  { value: ORDER_STATUS.CREATED, label: 'Mới', icon: '🆕' },
  { value: ORDER_STATUS.PAID, label: 'Đã thu', icon: '💰' },
  { value: ORDER_STATUS.QUEUED, label: 'Chờ pha', icon: '⏳' },
  { value: ORDER_STATUS.IN_PROGRESS, label: 'Đang pha', icon: '🍹' },
  { value: ORDER_STATUS.READY, label: 'Sẵn sàng', icon: '✅' },
  { value: ORDER_STATUS.SERVED, label: 'Hoàn tất', icon: '🎉' }
]
