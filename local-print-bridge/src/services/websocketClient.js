const io = require('socket.io-client')
const logger = require('../utils/logger')

class WebSocketClient {
  constructor() {
    this.socket = null
    this.connected = false
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 10
    this.onPrintJobCallback = null
  }

  /**
   * Connect to backend WebSocket server
   * @param {string} backendUrl - Backend URL (e.g., http://localhost:3000)
   */
  connect(backendUrl) {
    if (this.socket && this.connected) {
      logger.info('[WebSocket] Already connected')
      return
    }

    logger.info(`[WebSocket] Connecting to: ${backendUrl}`)

    this.socket = io(backendUrl, {
      path: '/socket.io/',
      transports: ['websocket'],
      reconnection: true,
      reconnectionDelay: 2000,
      reconnectionDelayMax: 10000,
      reconnectionAttempts: this.maxReconnectAttempts,
      timeout: 20000,
      forceNew: true
    })

    // Connection events
    this.socket.on('connect', () => {
      logger.info('[WebSocket] ✅ Connected to backend')
      this.connected = true
      this.reconnectAttempts = 0
    })

    this.socket.on('disconnect', (reason) => {
      logger.warn(`[WebSocket] Disconnected: ${reason}`)
      this.connected = false
    })

    this.socket.on('connect_error', (error) => {
      this.reconnectAttempts++
      logger.error(`[WebSocket] Connection error (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts}):`, error.message)
      
      if (this.reconnectAttempts >= this.maxReconnectAttempts) {
        logger.error('[WebSocket] Max reconnection attempts reached. Will retry in 30 seconds...')
        setTimeout(() => {
          this.reconnectAttempts = 0
          this.socket.connect()
        }, 30000)
      }
    })

    this.socket.on('reconnect', (attemptNumber) => {
      logger.info(`[WebSocket] ✅ Reconnected after ${attemptNumber} attempts`)
      this.connected = true
      this.reconnectAttempts = 0
    })

    this.socket.on('reconnect_error', (error) => {
      logger.error('[WebSocket] Reconnection error:', error.message)
    })

    this.socket.on('reconnect_failed', () => {
      logger.error('[WebSocket] ❌ Reconnection failed after all attempts')
      this.connected = false
    })

    // Listen for print job events
    this.socket.on('print-job-created', (data) => {
      logger.info('[WebSocket] 📨 New print job received:', data)
      if (this.onPrintJobCallback && data.job) {
        // Extract job data and convert to format expected by handler
        const jobData = {
          id: data.job.id,
          content: data.job.content,
          printer_ip: data.job.printer_ip,
          printer_port: data.job.printer_port || 9100,
          type: data.job.type
        }
        this.onPrintJobCallback(jobData)
      }
    })

    this.socket.on('print-job-status-changed', (data) => {
      logger.debug('[WebSocket] Print job status changed:', data)
    })

    this.socket.on('print-job-failed', (data) => {
      logger.warn('[WebSocket] Print job failed:', data)
    })
  }

  /**
   * Disconnect from WebSocket server
   */
  disconnect() {
    if (this.socket) {
      logger.info('[WebSocket] Disconnecting...')
      this.socket.disconnect()
      this.socket = null
      this.connected = false
    }
  }

  /**
   * Register callback for new print jobs
   * @param {Function} callback - Function to call when new job arrives
   */
  onPrintJob(callback) {
    this.onPrintJobCallback = callback
  }

  /**
   * Check if connected
   * @returns {boolean}
   */
  isConnected() {
    return this.connected
  }

  /**
   * Emit event to server
   * @param {string} event - Event name
   * @param {*} data - Event data
   */
  emit(event, data) {
    if (!this.socket || !this.connected) {
      logger.warn('[WebSocket] Cannot emit event, not connected')
      return
    }
    this.socket.emit(event, data)
  }
}

// Export singleton instance
module.exports = new WebSocketClient()
