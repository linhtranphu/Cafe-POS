const net = require('net')
const logger = require('../utils/logger')

// Statistics
const stats = {
  totalPrints: 0,
  successfulPrints: 0,
  failedPrints: 0,
  lastPrintTime: null
}

/**
 * Send content to thermal printer via TCP socket
 * @param {string} content - ESC/POS formatted content
 * @param {string} printerIP - Printer IP address
 * @param {number} printerPort - Printer port (default: 9100)
 * @returns {Promise<void>}
 */
function print(content, printerIP, printerPort = 9100) {
  return new Promise((resolve, reject) => {
    stats.totalPrints++
    
    const timeout = parseInt(process.env.PRINTER_TIMEOUT) || 5000
    let isResolved = false

    // Create TCP client
    const client = new net.Socket()
    
    // Set timeout
    client.setTimeout(timeout)

    // Connection established
    client.on('connect', () => {
      logger.debug(`Connected to printer ${printerIP}:${printerPort}`)
      
      try {
        // Send content to printer
        client.write(content, 'utf8', (err) => {
          if (err) {
            logger.error(`Write error to ${printerIP}:`, err.message)
            if (!isResolved) {
              isResolved = true
              stats.failedPrints++
              client.destroy()
              reject(new Error(`Failed to write to printer: ${err.message}`))
            }
          } else {
            logger.debug(`Content sent to printer ${printerIP}`)
            // Close connection after sending
            client.end()
          }
        })
      } catch (error) {
        logger.error(`Error sending to ${printerIP}:`, error.message)
        if (!isResolved) {
          isResolved = true
          stats.failedPrints++
          client.destroy()
          reject(error)
        }
      }
    })

    // Connection closed successfully
    client.on('end', () => {
      logger.debug(`Connection closed to ${printerIP}`)
      if (!isResolved) {
        isResolved = true
        stats.successfulPrints++
        stats.lastPrintTime = new Date().toISOString()
        resolve()
      }
    })

    // Connection error
    client.on('error', (err) => {
      logger.error(`Connection error to ${printerIP}:${printerPort}:`, err.message)
      if (!isResolved) {
        isResolved = true
        stats.failedPrints++
        
        let errorMessage = err.message
        if (err.code === 'ECONNREFUSED') {
          errorMessage = `Printer offline or unreachable at ${printerIP}:${printerPort}`
        } else if (err.code === 'ETIMEDOUT') {
          errorMessage = `Connection timeout to ${printerIP}:${printerPort}`
        } else if (err.code === 'EHOSTUNREACH') {
          errorMessage = `Host unreachable: ${printerIP}`
        }
        
        reject(new Error(errorMessage))
      }
    })

    // Timeout
    client.on('timeout', () => {
      logger.error(`Timeout connecting to ${printerIP}:${printerPort}`)
      if (!isResolved) {
        isResolved = true
        stats.failedPrints++
        client.destroy()
        reject(new Error(`Connection timeout to ${printerIP}:${printerPort}`))
      }
    })

    // Connect to printer
    try {
      client.connect(printerPort, printerIP)
    } catch (error) {
      if (!isResolved) {
        isResolved = true
        stats.failedPrints++
        reject(error)
      }
    }
  })
}

/**
 * Test printer connection
 * @param {string} printerIP - Printer IP address
 * @param {number} printerPort - Printer port (default: 9100)
 * @returns {Promise<void>}
 */
function testConnection(printerIP, printerPort = 9100) {
  return new Promise((resolve, reject) => {
    const timeout = parseInt(process.env.PRINTER_TIMEOUT) || 5000
    let isResolved = false

    const client = new net.Socket()
    client.setTimeout(timeout)

    client.on('connect', () => {
      logger.info(`✅ Printer ${printerIP}:${printerPort} is online`)
      if (!isResolved) {
        isResolved = true
        client.destroy()
        resolve()
      }
    })

    client.on('error', (err) => {
      logger.error(`❌ Printer ${printerIP}:${printerPort} connection failed:`, err.message)
      if (!isResolved) {
        isResolved = true
        
        let errorMessage = err.message
        if (err.code === 'ECONNREFUSED') {
          errorMessage = `Printer offline at ${printerIP}:${printerPort}`
        } else if (err.code === 'ETIMEDOUT') {
          errorMessage = `Connection timeout to ${printerIP}:${printerPort}`
        } else if (err.code === 'EHOSTUNREACH') {
          errorMessage = `Host unreachable: ${printerIP}`
        }
        
        reject(new Error(errorMessage))
      }
    })

    client.on('timeout', () => {
      if (!isResolved) {
        isResolved = true
        client.destroy()
        reject(new Error(`Connection timeout to ${printerIP}:${printerPort}`))
      }
    })

    try {
      client.connect(printerPort, printerIP)
    } catch (error) {
      if (!isResolved) {
        isResolved = true
        reject(error)
      }
    }
  })
}

/**
 * Get printing statistics
 * @returns {object} Statistics object
 */
function getStats() {
  return {
    ...stats,
    successRate: stats.totalPrints > 0 
      ? ((stats.successfulPrints / stats.totalPrints) * 100).toFixed(2) + '%'
      : 'N/A'
  }
}

module.exports = {
  print,
  testConnection,
  getStats
}
