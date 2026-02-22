const express = require('express')
const cors = require('cors')
const dotenv = require('dotenv')
const printerService = require('./services/printerService')
const backendSync = require('./services/backendSync')
const websocketClient = require('./services/websocketClient')
const printJobHandler = require('./services/printJobHandler')
const logger = require('./utils/logger')

// Load environment variables
dotenv.config()

const app = express()
const PORT = process.env.PORT || 3001

// Middleware
app.use(cors({
  origin: '*', // Allow all origins for local development
  methods: ['GET', 'POST', 'PUT', 'DELETE'],
  credentials: true
}))
app.use(express.json())

// Request logging
app.use((req, res, next) => {
  logger.info(`${req.method} ${req.path}`)
  next()
})

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({
    status: 'ok',
    service: 'Local Print Bridge',
    version: '1.0.0',
    timestamp: new Date().toISOString()
  })
})

// Print endpoint - Main functionality
app.post('/print', async (req, res) => {
  const { jobId, content, printerIP, printerPort, type } = req.body

  // Validation
  if (!jobId || !content || !printerIP) {
    return res.status(400).json({
      success: false,
      error: 'Missing required fields: jobId, content, printerIP'
    })
  }

  logger.info(`Print request received - Job ID: ${jobId}, Printer: ${printerIP}:${printerPort || 9100}, Type: ${type || 'unknown'}`)

  try {
    // Send to printer
    await printerService.print(content, printerIP, printerPort || 9100)
    
    logger.info(`Print successful - Job ID: ${jobId}`)

    // Update backend status
    try {
      await backendSync.updateJobStatus(jobId, 'COMPLETED')
      logger.info(`Backend updated - Job ID: ${jobId} -> COMPLETED`)
    } catch (backendError) {
      logger.error(`Failed to update backend for job ${jobId}:`, backendError.message)
      // Don't fail the request if backend update fails
    }

    res.json({
      success: true,
      jobId,
      message: 'Print completed successfully',
      timestamp: new Date().toISOString()
    })

  } catch (error) {
    logger.error(`Print failed - Job ID: ${jobId}:`, error.message)

    // Update backend with failure status
    try {
      await backendSync.updateJobStatus(jobId, 'FAILED', error.message)
      logger.info(`Backend updated - Job ID: ${jobId} -> FAILED`)
    } catch (backendError) {
      logger.error(`Failed to update backend for job ${jobId}:`, backendError.message)
    }

    res.status(500).json({
      success: false,
      error: error.message,
      jobId,
      timestamp: new Date().toISOString()
    })
  }
})

// Test printer connection
app.post('/test-connection', async (req, res) => {
  const { printerIP, printerPort } = req.body

  if (!printerIP) {
    return res.status(400).json({
      success: false,
      error: 'Missing printerIP'
    })
  }

  logger.info(`Testing connection to printer: ${printerIP}:${printerPort || 9100}`)

  try {
    await printerService.testConnection(printerIP, printerPort || 9100)
    
    res.json({
      success: true,
      message: 'Printer connection successful',
      printer: `${printerIP}:${printerPort || 9100}`
    })
  } catch (error) {
    logger.error(`Connection test failed for ${printerIP}:`, error.message)
    
    res.status(500).json({
      success: false,
      error: error.message,
      printer: `${printerIP}:${printerPort || 9100}`
    })
  }
})

// Get printer status
app.get('/status', (req, res) => {
  const stats = printerService.getStats()
  
  res.json({
    success: true,
    stats,
    uptime: process.uptime(),
    timestamp: new Date().toISOString()
  })
})

// Error handling middleware
app.use((err, req, res, next) => {
  logger.error('Unhandled error:', err)
  res.status(500).json({
    success: false,
    error: 'Internal server error',
    message: err.message
  })
})

// Start server
app.listen(PORT, () => {
  logger.info('='.repeat(50))
  logger.info('🖨️  Local Print Bridge Started')
  logger.info('='.repeat(50))
  logger.info(`Server running on: http://localhost:${PORT}`)
  logger.info(`Backend URL: ${process.env.BACKEND_URL || 'Not configured'}`)
  logger.info(`Default Bill Printer: ${process.env.DEFAULT_BILL_PRINTER_IP || 'Not configured'}`)
  logger.info(`Default Label Printer: ${process.env.DEFAULT_LABEL_PRINTER_IP || 'Not configured'}`)
  logger.info('='.repeat(50))
  logger.info('Ready to accept print requests!')
  logger.info('='.repeat(50))

  // Connect to backend WebSocket for real-time job notifications
  const backendUrl = process.env.BACKEND_URL
  if (backendUrl) {
    logger.info(`[WebSocket] Connecting to backend: ${backendUrl}`)
    websocketClient.connect(backendUrl)
    
    // Register handler for incoming print jobs
    websocketClient.onPrintJob((jobData) => {
      logger.info('[WebSocket] 📨 Received print job via WebSocket:', jobData.id)
      printJobHandler.handlePrintJob(jobData)
    })
  } else {
    logger.warn('[WebSocket] ⚠️  BACKEND_URL not configured - WebSocket disabled')
    logger.warn('[WebSocket] Print jobs will only work via HTTP POST /print endpoint')
  }
})

// Graceful shutdown
process.on('SIGTERM', () => {
  logger.info('SIGTERM received, shutting down gracefully...')
  websocketClient.disconnect()
  process.exit(0)
})

process.on('SIGINT', () => {
  logger.info('SIGINT received, shutting down gracefully...')
  websocketClient.disconnect()
  process.exit(0)
})
