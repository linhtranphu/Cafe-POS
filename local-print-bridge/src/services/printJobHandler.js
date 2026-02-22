const printerService = require('./printerService')
const backendSync = require('./backendSync')
const logger = require('../utils/logger')

/**
 * Handle incoming print job from WebSocket
 * @param {Object} jobData - Print job data from backend
 */
async function handlePrintJob(jobData) {
  const { id, content, printer_ip, printer_port, type } = jobData

  if (!id || !content || !printer_ip) {
    logger.error('[PrintJobHandler] Invalid job data:', jobData)
    return
  }

  logger.info(`[PrintJobHandler] Processing job ${id} - Type: ${type}, Printer: ${printer_ip}:${printer_port || 9100}`)

  try {
    // Send to printer
    await printerService.print(content, printer_ip, printer_port || 9100)
    
    logger.info(`[PrintJobHandler] ✅ Job ${id} printed successfully`)

    // Update backend status
    try {
      await backendSync.updateJobStatus(id, 'COMPLETED')
      logger.info(`[PrintJobHandler] ✅ Backend updated - Job ${id} -> COMPLETED`)
    } catch (backendError) {
      logger.error(`[PrintJobHandler] Failed to update backend for job ${id}:`, backendError.message)
    }

  } catch (error) {
    logger.error(`[PrintJobHandler] ❌ Job ${id} failed:`, error.message)

    // Update backend with failure status
    try {
      await backendSync.updateJobStatus(id, 'FAILED', error.message)
      logger.info(`[PrintJobHandler] Backend updated - Job ${id} -> FAILED`)
    } catch (backendError) {
      logger.error(`[PrintJobHandler] Failed to update backend for job ${id}:`, backendError.message)
    }
  }
}

module.exports = {
  handlePrintJob
}
