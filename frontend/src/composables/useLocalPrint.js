import { ref, onMounted, onUnmounted } from 'vue'
import { localPrintService } from '../services/localPrint'
import { usePrintJobStore } from '../stores/printJob'

/**
 * Composable for managing local print bridge
 */
export function useLocalPrint() {
  const printJobStore = usePrintJobStore()
  const isAvailable = ref(false)
  const isChecking = ref(false)
  const lastCheckTime = ref(null)

  /**
   * Check if local print bridge is available
   */
  const checkAvailability = async () => {
    isChecking.value = true
    try {
      isAvailable.value = await localPrintService.checkAvailability()
      lastCheckTime.value = new Date()
    } finally {
      isChecking.value = false
    }
  }

  /**
   * Initialize local print bridge
   */
  const initialize = async () => {
    await checkAvailability()
    
    if (isAvailable.value) {
      localPrintService.startHealthCheck()
      console.log('[useLocalPrint] Local print bridge initialized')
    }
  }

  /**
   * Test printer connection
   */
  const testPrinterConnection = async (printerIP, printerPort = 9100) => {
    if (!isAvailable.value) {
      throw new Error('Local print bridge is not available')
    }

    return await localPrintService.testConnection(printerIP, printerPort)
  }

  /**
   * Get bridge status
   */
  const getBridgeStatus = async () => {
    if (!isAvailable.value) {
      throw new Error('Local print bridge is not available')
    }

    return await localPrintService.getStatus()
  }

  /**
   * Cleanup
   */
  const cleanup = () => {
    localPrintService.stopHealthCheck()
  }

  // Auto-initialize on mount
  onMounted(() => {
    initialize()
  })

  // Cleanup on unmount
  onUnmounted(() => {
    cleanup()
  })

  return {
    isAvailable,
    isChecking,
    lastCheckTime,
    checkAvailability,
    testPrinterConnection,
    getBridgeStatus,
    initialize,
    cleanup
  }
}
