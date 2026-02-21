import { onMounted, onUnmounted } from 'vue'
import { websocketService } from '../services/websocket'
import { usePrintJobStore } from '../stores/printJob'
import { useNotifications } from './useNotifications'

/**
 * Composable for handling print job WebSocket events
 * Listens to:
 * - print-job-created: When a new print job is created
 * - print-job-status-changed: When a print job status changes
 * - print-job-failed: When a print job fails
 */
export function usePrintJobWebSocket() {
  const printJobStore = usePrintJobStore()
  const { showPrintJobFailed, showSuccess } = useNotifications()

  // Event handlers
  const handlePrintJobCreated = (data) => {
    console.log('[WebSocket] Print job created:', data)
    
    // Add new job to store
    if (data.job) {
      printJobStore.addJob(data.job)
    }
  }

  const handlePrintJobStatusChanged = (data) => {
    console.log('[WebSocket] Print job status changed:', data)
    
    // Update job status in store
    if (data.job_id && data.status) {
      printJobStore.updateJobStatus(data.job_id, data.status, data.error_msg)
      
      // Show success notification when job completes
      if (data.status === 'COMPLETED') {
        showSuccess(`In ${data.type === 'BILL' ? 'bill' : 'tem'} thành công cho đơn ${data.order_number}`)
      }
    }
  }

  const handlePrintJobFailed = (data) => {
    console.log('[WebSocket] Print job failed:', data)
    
    // Update job status in store
    if (data.job_id) {
      printJobStore.updateJobStatus(data.job_id, 'FAILED', data.error_msg)
      
      // Show error notification with retry action
      showPrintJobFailed({
        type: data.type,
        order_number: data.order_number,
        onRetry: async () => {
          await printJobStore.retryJob(data.job_id)
        }
      })
    }
  }

  // Setup listeners
  const setupListeners = () => {
    if (!websocketService.isConnected()) {
      console.warn('[WebSocket] Not connected, cannot setup print job listeners')
      return
    }

    websocketService.on('print-job-created', handlePrintJobCreated)
    websocketService.on('print-job-status-changed', handlePrintJobStatusChanged)
    websocketService.on('print-job-failed', handlePrintJobFailed)
  }

  // Cleanup listeners
  const cleanupListeners = () => {
    websocketService.off('print-job-created', handlePrintJobCreated)
    websocketService.off('print-job-status-changed', handlePrintJobStatusChanged)
    websocketService.off('print-job-failed', handlePrintJobFailed)
  }

  // Lifecycle hooks
  onMounted(() => {
    setupListeners()
  })

  onUnmounted(() => {
    cleanupListeners()
  })

  return {
    setupListeners,
    cleanupListeners
  }
}
