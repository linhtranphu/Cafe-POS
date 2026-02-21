import { defineStore } from 'pinia'
import { printJobService } from '../services/printJob'
import { localPrintService } from '../services/localPrint'
import { websocketService } from '../services/websocket'

export const usePrintJobStore = defineStore('printJob', {
  state: () => ({
    printJobs: [],
    pendingJobs: [],
    failedJobs: [],
    loading: false,
    error: null,
    failedJobNotifications: [], // For real-time failure notifications
    localBridgeAvailable: false,
    websocketListenersSetup: false
  }),

  getters: {
    jobsByOrderId: (state) => (orderId) => {
      return state.printJobs.filter(job => job.order_id === orderId)
    },

    jobsByStatus: (state) => (status) => {
      return state.printJobs.filter(job => job.status === status)
    },

    getPendingCount: (state) => {
      return state.pendingJobs.length
    },

    getFailedCount: (state) => {
      return state.failedJobs.length
    }
  },

  actions: {
    /**
     * Initialize local print bridge and WebSocket listeners
     */
    async initialize() {
      // Check local bridge availability
      this.localBridgeAvailable = await localPrintService.checkAvailability()
      
      if (this.localBridgeAvailable) {
        console.log('[PrintJob] Local print bridge detected')
        localPrintService.startHealthCheck()
      } else {
        console.log('[PrintJob] Local print bridge not available')
      }

      // Setup WebSocket listeners
      this.setupWebSocketListeners()
    },

    /**
     * Setup WebSocket event listeners for print jobs
     */
    setupWebSocketListeners() {
      if (this.websocketListenersSetup) return

      // Listen for new print jobs created
      websocketService.on('print-job-created', (data) => {
        console.log('[PrintJob] New print job created:', data)
        this.addJob(data.job)
        
        // If local bridge is available, trigger local printing
        if (this.localBridgeAvailable && data.job.status === 'PENDING') {
          this.handleLocalPrint(data.job)
        }
      })

      // Listen for print job status changes
      websocketService.on('print-job-status-changed', (data) => {
        console.log('[PrintJob] Print job status changed:', data)
        this.updateJobStatus(data.job_id, data.status, data.error_msg)
      })

      // Listen for print job failures
      websocketService.on('print-job-failed', (data) => {
        console.log('[PrintJob] Print job failed:', data)
        this.updateJobStatus(data.job_id, 'FAILED', data.error_msg)
        this.addFailedJobNotification(data)
      })

      this.websocketListenersSetup = true
      console.log('[PrintJob] WebSocket listeners setup complete')
    },

    /**
     * Handle local printing via print bridge
     */
    async handleLocalPrint(job) {
      if (!this.localBridgeAvailable) {
        console.log('[PrintJob] Local bridge not available, skipping local print')
        return
      }

      try {
        console.log('[PrintJob] Sending to local printer:', job.id)
        
        // Extract printer IP from printer config
        // Assuming printer_config has connection_details with ip and port
        const printerIP = job.printer_config?.connection_details?.ip
        const printerPort = job.printer_config?.connection_details?.port || 9100

        if (!printerIP) {
          console.error('[PrintJob] No printer IP found in job config')
          return
        }

        // Send to local print bridge
        await localPrintService.print(
          job.id,
          job.content,
          printerIP,
          printerPort,
          job.type
        )

        console.log('[PrintJob] Local print successful:', job.id)
      } catch (error) {
        console.error('[PrintJob] Local print failed:', error)
        // Error will be handled by the local bridge's backend sync
      }
    },

    /**
     * Cleanup WebSocket listeners
     */
    cleanup() {
      localPrintService.stopHealthCheck()
      this.websocketListenersSetup = false
    },

    async fetchJobs(params = {}) {
      this.loading = true
      this.error = null
      try {
        this.printJobs = await printJobService.fetchPrintJobs(params)
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải danh sách print jobs'
        this.printJobs = []
      } finally {
        this.loading = false
      }
    },

    async fetchPendingJobs() {
      this.loading = true
      this.error = null
      try {
        this.pendingJobs = await printJobService.fetchPendingJobs()
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải pending jobs'
        this.pendingJobs = []
      } finally {
        this.loading = false
      }
    },

    async fetchFailedJobs() {
      this.loading = true
      this.error = null
      try {
        this.failedJobs = await printJobService.fetchFailedJobs()
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải failed jobs'
        this.failedJobs = []
      } finally {
        this.loading = false
      }
    },

    async retryJob(jobId) {
      this.error = null
      try {
        await printJobService.retryJob(jobId)
        // Refresh the jobs after retry
        await this.fetchFailedJobs()
        await this.fetchPendingJobs()
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi retry print job'
        return false
      }
    },

    async cancelJob(jobId) {
      this.error = null
      try {
        await printJobService.cancelJob(jobId)
        // Remove from local state
        this.printJobs = this.printJobs.filter(job => job.id !== jobId)
        this.pendingJobs = this.pendingJobs.filter(job => job.id !== jobId)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi hủy print job'
        return false
      }
    },

    // Update job in local state (useful for real-time updates)
    updateJobStatus(jobId, status, errorMsg = null) {
      const updateJob = (job) => {
        if (job.id === jobId) {
          job.status = status
          if (errorMsg) {
            job.error_msg = errorMsg
          }
          job.updated_at = new Date().toISOString()
        }
        return job
      }

      this.printJobs = this.printJobs.map(updateJob)
      this.pendingJobs = this.pendingJobs.map(updateJob)
      this.failedJobs = this.failedJobs.map(updateJob)

      // Move jobs between lists based on status
      if (status === 'FAILED') {
        const job = this.printJobs.find(j => j.id === jobId)
        if (job && !this.failedJobs.find(j => j.id === jobId)) {
          this.failedJobs.push(job)
        }
        this.pendingJobs = this.pendingJobs.filter(j => j.id !== jobId)
      } else if (status === 'PENDING') {
        const job = this.printJobs.find(j => j.id === jobId)
        if (job && !this.pendingJobs.find(j => j.id === jobId)) {
          this.pendingJobs.push(job)
        }
        this.failedJobs = this.failedJobs.filter(j => j.id !== jobId)
      } else if (status === 'COMPLETED') {
        this.pendingJobs = this.pendingJobs.filter(j => j.id !== jobId)
        this.failedJobs = this.failedJobs.filter(j => j.id !== jobId)
      }
    },

    // Add new job to store (for WebSocket events)
    addJob(job) {
      // Check if job already exists
      const exists = this.printJobs.find(j => j.id === job.id)
      if (!exists) {
        this.printJobs.push(job)
        
        // Add to appropriate list based on status
        if (job.status === 'PENDING') {
          this.pendingJobs.push(job)
        } else if (job.status === 'FAILED') {
          this.failedJobs.push(job)
        }
      }
    },

    // Add failed job notification
    addFailedJobNotification(data) {
      this.failedJobNotifications.push({
        id: Date.now(),
        jobId: data.job_id,
        orderNumber: data.order_number,
        type: data.type,
        errorMsg: data.error_msg,
        timestamp: new Date().toISOString()
      })
    },

    // Remove notification
    removeNotification(notificationId) {
      this.failedJobNotifications = this.failedJobNotifications.filter(
        n => n.id !== notificationId
      )
    },

    // Clear all notifications
    clearNotifications() {
      this.failedJobNotifications = []
    }
  }
})
