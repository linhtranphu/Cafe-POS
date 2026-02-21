/**
 * Local Print Bridge Service
 * Detects and communicates with the local print bridge running at the cafe
 */

const LOCAL_BRIDGE_URL = 'http://localhost:3001'
const HEALTH_CHECK_INTERVAL = 30000 // 30 seconds

class LocalPrintService {
  constructor() {
    this.available = false
    this.checking = false
    this.healthCheckTimer = null
  }

  /**
   * Check if local print bridge is available
   * @returns {Promise<boolean>}
   */
  async checkAvailability() {
    if (this.checking) return this.available

    this.checking = true
    try {
      const response = await fetch(`${LOCAL_BRIDGE_URL}/health`, {
        method: 'GET',
        signal: AbortSignal.timeout(2000) // 2 second timeout
      })
      
      if (response.ok) {
        const data = await response.json()
        this.available = data.status === 'ok'
        console.log('[LocalPrint] Bridge available:', this.available)
      } else {
        this.available = false
      }
    } catch (error) {
      this.available = false
      console.log('[LocalPrint] Bridge not available:', error.message)
    } finally {
      this.checking = false
    }

    return this.available
  }

  /**
   * Start periodic health checks
   */
  startHealthCheck() {
    if (this.healthCheckTimer) return

    // Initial check
    this.checkAvailability()

    // Periodic checks
    this.healthCheckTimer = setInterval(() => {
      this.checkAvailability()
    }, HEALTH_CHECK_INTERVAL)

    console.log('[LocalPrint] Health check started')
  }

  /**
   * Stop periodic health checks
   */
  stopHealthCheck() {
    if (this.healthCheckTimer) {
      clearInterval(this.healthCheckTimer)
      this.healthCheckTimer = null
      console.log('[LocalPrint] Health check stopped')
    }
  }

  /**
   * Send print job to local bridge
   * @param {string} jobId - Print job ID
   * @param {string} content - ESC/POS formatted content
   * @param {string} printerIP - Printer IP address
   * @param {number} printerPort - Printer port (default: 9100)
   * @param {string} type - Job type (BILL or LABEL)
   * @returns {Promise<object>}
   */
  async print(jobId, content, printerIP, printerPort = 9100, type = 'BILL') {
    if (!this.available) {
      throw new Error('Local print bridge is not available')
    }

    try {
      const response = await fetch(`${LOCAL_BRIDGE_URL}/print`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          jobId,
          content,
          printerIP,
          printerPort,
          type
        }),
        signal: AbortSignal.timeout(10000) // 10 second timeout
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Print failed')
      }

      console.log('[LocalPrint] Print successful:', jobId)
      return data
    } catch (error) {
      console.error('[LocalPrint] Print error:', error)
      throw error
    }
  }

  /**
   * Test printer connection via local bridge
   * @param {string} printerIP - Printer IP address
   * @param {number} printerPort - Printer port (default: 9100)
   * @returns {Promise<object>}
   */
  async testConnection(printerIP, printerPort = 9100) {
    if (!this.available) {
      throw new Error('Local print bridge is not available')
    }

    try {
      const response = await fetch(`${LOCAL_BRIDGE_URL}/test-connection`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          printerIP,
          printerPort
        }),
        signal: AbortSignal.timeout(5000) // 5 second timeout
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Connection test failed')
      }

      return data
    } catch (error) {
      console.error('[LocalPrint] Connection test error:', error)
      throw error
    }
  }

  /**
   * Get local bridge status and statistics
   * @returns {Promise<object>}
   */
  async getStatus() {
    if (!this.available) {
      throw new Error('Local print bridge is not available')
    }

    try {
      const response = await fetch(`${LOCAL_BRIDGE_URL}/status`, {
        method: 'GET',
        signal: AbortSignal.timeout(2000)
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error('Failed to get status')
      }

      return data
    } catch (error) {
      console.error('[LocalPrint] Status error:', error)
      throw error
    }
  }

  /**
   * Check if local bridge is available
   * @returns {boolean}
   */
  isAvailable() {
    return this.available
  }
}

// Export singleton instance
export const localPrintService = new LocalPrintService()
