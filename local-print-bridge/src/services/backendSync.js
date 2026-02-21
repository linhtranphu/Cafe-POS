const axios = require('axios')
const logger = require('../utils/logger')

const BACKEND_URL = process.env.BACKEND_URL

/**
 * Update print job status in backend
 * @param {string} jobId - Print job ID
 * @param {string} status - Job status (COMPLETED, FAILED)
 * @param {string} errorMsg - Error message (optional)
 * @returns {Promise<void>}
 */
async function updateJobStatus(jobId, status, errorMsg = '') {
  if (!BACKEND_URL) {
    logger.warn('BACKEND_URL not configured, skipping status update')
    return
  }

  const url = `${BACKEND_URL}/api/print-jobs/${jobId}/status`
  
  try {
    const response = await axios.put(url, {
      status,
      error_msg: errorMsg
    }, {
      timeout: 5000,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    logger.debug(`Backend status updated for job ${jobId}: ${status}`)
    return response.data
  } catch (error) {
    if (error.response) {
      // Server responded with error status
      logger.error(`Backend update failed (${error.response.status}):`, error.response.data)
      throw new Error(`Backend returned ${error.response.status}: ${error.response.data.error || 'Unknown error'}`)
    } else if (error.request) {
      // Request made but no response
      logger.error('Backend update failed: No response from server')
      throw new Error('Cannot reach backend server')
    } else {
      // Error setting up request
      logger.error('Backend update failed:', error.message)
      throw error
    }
  }
}

/**
 * Fetch pending print jobs from backend
 * @returns {Promise<Array>} Array of pending jobs
 */
async function fetchPendingJobs() {
  if (!BACKEND_URL) {
    logger.warn('BACKEND_URL not configured, cannot fetch jobs')
    return []
  }

  const url = `${BACKEND_URL}/api/print-jobs/pending`
  
  try {
    const response = await axios.get(url, {
      timeout: 5000
    })

    return response.data || []
  } catch (error) {
    logger.error('Failed to fetch pending jobs:', error.message)
    return []
  }
}

module.exports = {
  updateJobStatus,
  fetchPendingJobs
}
