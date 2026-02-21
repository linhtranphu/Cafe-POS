import { ref, onMounted, onUnmounted } from 'vue'
import { websocketService } from '../services/websocket'
import { useNotifications } from './useNotifications'

/**
 * Composable for monitoring printer status via WebSocket
 * Listens to:
 * - printer-offline: When a printer goes offline
 * - printer-online: When a printer comes back online
 * - printer-error: When a printer has a hardware error
 */
export function usePrinterStatus() {
  const { showPrinterOffline, showWarning, showSuccess } = useNotifications()
  const printerStatuses = ref(new Map())

  // Event handlers
  const handlePrinterOffline = (data) => {
    console.log('[WebSocket] Printer offline:', data)
    
    if (data.printer_id && data.printer_name) {
      printerStatuses.value.set(data.printer_id, {
        status: 'offline',
        name: data.printer_name,
        timestamp: Date.now()
      })
      
      showPrinterOffline(data.printer_name)
    }
  }

  const handlePrinterOnline = (data) => {
    console.log('[WebSocket] Printer online:', data)
    
    if (data.printer_id && data.printer_name) {
      printerStatuses.value.set(data.printer_id, {
        status: 'online',
        name: data.printer_name,
        timestamp: Date.now()
      })
      
      showSuccess(`Máy in "${data.printer_name}" đã kết nối`)
    }
  }

  const handlePrinterError = (data) => {
    console.log('[WebSocket] Printer error:', data)
    
    if (data.printer_id && data.printer_name) {
      printerStatuses.value.set(data.printer_id, {
        status: 'error',
        name: data.printer_name,
        error: data.error_msg,
        timestamp: Date.now()
      })
      
      // Show specific error message
      const errorMessages = {
        'paper_out': 'Máy in hết giấy',
        'paper_jam': 'Máy in bị kẹt giấy',
        'cover_open': 'Nắp máy in đang mở',
        'hardware_error': 'Lỗi phần cứng máy in'
      }
      
      const message = errorMessages[data.error_type] || data.error_msg || 'Lỗi máy in'
      
      showWarning(`${data.printer_name}: ${message}`, {
        duration: 0 // Don't auto-dismiss hardware errors
      })
    }
  }

  // Setup listeners
  const setupListeners = () => {
    if (!websocketService.isConnected()) {
      console.warn('[WebSocket] Not connected, cannot setup printer status listeners')
      return
    }

    websocketService.on('printer-offline', handlePrinterOffline)
    websocketService.on('printer-online', handlePrinterOnline)
    websocketService.on('printer-error', handlePrinterError)
  }

  // Cleanup listeners
  const cleanupListeners = () => {
    websocketService.off('printer-offline', handlePrinterOffline)
    websocketService.off('printer-online', handlePrinterOnline)
    websocketService.off('printer-error', handlePrinterError)
  }

  // Get printer status
  const getPrinterStatus = (printerId) => {
    return printerStatuses.value.get(printerId) || { status: 'unknown' }
  }

  // Lifecycle hooks
  onMounted(() => {
    setupListeners()
  })

  onUnmounted(() => {
    cleanupListeners()
  })

  return {
    printerStatuses,
    getPrinterStatus,
    setupListeners,
    cleanupListeners
  }
}
