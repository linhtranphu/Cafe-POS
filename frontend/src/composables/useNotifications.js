import { ref, computed } from 'vue'

// Global notification state
const notifications = ref([])
let notificationId = 0

export function useNotifications() {
  const activeNotifications = computed(() => notifications.value)

  const addNotification = (notification) => {
    const id = ++notificationId
    const newNotification = {
      id,
      type: notification.type || 'info', // info, success, warning, error
      title: notification.title || '',
      message: notification.message || '',
      duration: notification.duration || 5000, // Auto-dismiss after 5 seconds
      action: notification.action || null, // { label: string, handler: function }
      dismissible: notification.dismissible !== false,
      timestamp: Date.now()
    }

    notifications.value.push(newNotification)

    // Auto-dismiss if duration is set
    if (newNotification.duration > 0) {
      setTimeout(() => {
        removeNotification(id)
      }, newNotification.duration)
    }

    return id
  }

  const removeNotification = (id) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
  }

  const clearAll = () => {
    notifications.value = []
  }

  // Convenience methods for different notification types
  const showSuccess = (message, options = {}) => {
    return addNotification({
      type: 'success',
      message,
      ...options
    })
  }

  const showError = (message, options = {}) => {
    return addNotification({
      type: 'error',
      message,
      duration: 0, // Don't auto-dismiss errors
      ...options
    })
  }

  const showWarning = (message, options = {}) => {
    return addNotification({
      type: 'warning',
      message,
      ...options
    })
  }

  const showInfo = (message, options = {}) => {
    return addNotification({
      type: 'info',
      message,
      ...options
    })
  }

  // Print-specific notifications
  const showPrintJobFailed = (data) => {
    return addNotification({
      type: 'error',
      title: 'In thất bại',
      message: `Không thể in ${data.type === 'BILL' ? 'bill' : 'tem'} cho đơn ${data.order_number}`,
      duration: 0, // Don't auto-dismiss
      action: {
        label: 'Thử lại',
        handler: () => data.onRetry && data.onRetry()
      }
    })
  }

  const showPrinterOffline = (printerName) => {
    return addNotification({
      type: 'warning',
      title: 'Máy in offline',
      message: `Máy in "${printerName}" không khả dụng`,
      duration: 0
    })
  }

  return {
    notifications: activeNotifications,
    addNotification,
    removeNotification,
    clearAll,
    showSuccess,
    showError,
    showWarning,
    showInfo,
    showPrintJobFailed,
    showPrinterOffline
  }
}
