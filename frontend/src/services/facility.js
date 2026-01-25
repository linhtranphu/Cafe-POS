import axios from 'axios'

// Cache refresh: 2024-12-30 - Added 404 error handling
const API_URL = 'http://localhost:8080/api'

const getAuthHeader = () => {
  const token = localStorage.getItem('token')
  return { Authorization: `Bearer ${token}` }
}

export const facilityService = {
  async getFacilities() {
    const response = await axios.get(`${API_URL}/manager/facilities`, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async getFacility(id) {
    const response = await axios.get(`${API_URL}/manager/facilities/${id}`, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async createFacility(facility) {
    const response = await axios.post(`${API_URL}/manager/facilities`, facility, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async updateFacility(id, facility) {
    const response = await axios.put(`${API_URL}/manager/facilities/${id}`, facility, {
      headers: getAuthHeader()
    })
    return response.data
  },

  // FR-FM-03: Update facility information with change tracking
  async updateFacilityWithHistory(id, updates) {
    const response = await axios.put(`${API_URL}/manager/facilities/${id}/update-with-history`, updates, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async moveAsset(id, newArea, reason) {
    const response = await axios.patch(`${API_URL}/manager/facilities/${id}/move`, {
      new_area: newArea,
      reason: reason
    }, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async checkDeletionEligibility(id) {
    const response = await axios.get(`${API_URL}/manager/facilities/${id}/can-delete`, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async deleteFacility(id) {
    const response = await axios.delete(`${API_URL}/manager/facilities/${id}`, {
      headers: getAuthHeader()
    })
    return response.data
  },

  // FR-FM-01: Asset tracking and inventory
  async getFacilityHistory(id) {
    const response = await axios.get(`${API_URL}/manager/facilities/${id}/history`, {
      headers: getAuthHeader()
    })
    return response.data
  },

  // FR-FM-04: Status management
  async getStatusAlerts() {
    const response = await axios.get(`${API_URL}/manager/facilities/status-alerts`, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async getStatusHistory(id) {
    console.log('getStatusHistory called with id:', id)
    console.log('Auth token:', localStorage.getItem('token'))
    try {
      const response = await axios.get(`${API_URL}/manager/facilities/${id}/status-history`, {
        headers: getAuthHeader()
      })
      console.log('getStatusHistory response:', response.data)
      return response.data
    } catch (error) {
      console.error('getStatusHistory error:', error)
      throw error
    }
  },

  async updateStatusWithDetails(id, statusData) {
    const response = await axios.patch(`${API_URL}/manager/facilities/${id}/status-detailed`, statusData, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async resolveStatusAlert(alertId) {
    const response = await axios.patch(`${API_URL}/manager/facilities/alerts/${alertId}/resolve`, {}, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async getFacilitiesByStatus(status) {
    const response = await axios.get(`${API_URL}/manager/facilities/status/${status}`, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async getFacilitiesByArea(area) {
    const response = await axios.get(`${API_URL}/manager/facilities/area/${area}`, {
      headers: getAuthHeader()
    })
    return response.data
  },

  // FR-FM-02: Maintenance scheduling and tracking
  async getMaintenanceHistory(id) {
    const response = await axios.get(`${API_URL}/manager/facilities/${id}/maintenance`, {
      headers: getAuthHeader()
    })
    return response.data
  },

  // FR-FM-06: Enhanced maintenance management
  async getMaintenanceStats(facilityId) {
    try {
      const response = await axios.get(`${API_URL}/manager/facilities/${facilityId}/maintenance-stats`, {
        headers: getAuthHeader()
      })
      return response.data
    } catch (error) {
      if (error.response?.status === 404) {
        console.warn('Maintenance stats endpoint not implemented')
        return { total: 0, totalCost: 0, avgInterval: 0, scheduledCost: 0, emergencyCost: 0 }
      }
      throw error
    }
  },

  async getNextMaintenanceDate(facilityId) {
    try {
      const response = await axios.get(`${API_URL}/manager/facilities/${facilityId}/next-maintenance`, {
        headers: getAuthHeader()
      })
      return response.data
    } catch (error) {
      if (error.response?.status === 404) {
        console.warn('Next maintenance date endpoint not implemented')
        return { next_date: null }
      }
      throw error
    }
  },

  async updateMaintenanceRecord(recordId, recordData) {
    const response = await axios.put(`${API_URL}/manager/maintenance/${recordId}`, recordData, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async getMaintenanceCostAnalysis(facilityId, period = '12months') {
    const response = await axios.get(`${API_URL}/manager/facilities/${facilityId}/maintenance-costs`, {
      params: { period },
      headers: getAuthHeader()
    })
    return response.data
  },

  async scheduleRecurringMaintenance(facilityId, schedule) {
    const response = await axios.post(`${API_URL}/manager/facilities/${facilityId}/recurring-maintenance`, schedule, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async scheduleMaintenanceTask(task) {
    const response = await axios.post(`${API_URL}/manager/maintenance/schedule`, task, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async getScheduledMaintenance() {
    try {
      const response = await axios.get(`${API_URL}/manager/maintenance/scheduled`, {
        headers: getAuthHeader()
      })
      return response.data
    } catch (error) {
      if (error.response?.status === 404) {
        console.warn('Scheduled maintenance endpoint not implemented')
        return []
      }
      throw error
    }
  },

  async updateMaintenanceTask(id, updates) {
    const response = await axios.patch(`${API_URL}/manager/maintenance/${id}`, updates, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async getMaintenanceDue() {
    try {
      const response = await axios.get(`${API_URL}/manager/maintenance/due`, {
        headers: getAuthHeader()
      })
      return response.data
    } catch (error) {
      if (error.response?.status === 404) {
        console.warn('Maintenance due endpoint not implemented')
        return []
      }
      throw error
    }
  },

  async getIssueReports() {
    const response = await axios.get(`${API_URL}/manager/issue-reports`, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async createMaintenanceRecord(record) {
    const response = await axios.post(`${API_URL}/manager/maintenance`, record, {
      headers: getAuthHeader()
    })
    return response.data
  },

  // FR-FM-05: Issue reporting
  async createIssueReport(reportData) {
    const response = await axios.post(`${API_URL}/staff/issue-reports`, reportData, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async updateIssueReportStatus(reportId, status) {
    const response = await axios.patch(`${API_URL}/manager/issue-reports/${reportId}/status`, {
      status
    }, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async addReportComment(reportId, comment) {
    const response = await axios.post(`${API_URL}/manager/issue-reports/${reportId}/comments`, {
      comment
    }, {
      headers: getAuthHeader()
    })
    return response.data
  },

  // FR-FM-08: Search & Filter
  async searchFacilities(filters = {}) {
    const params = new URLSearchParams()
    
    if (filters.name) params.append('name', filters.name)
    if (filters.type) params.append('type', filters.type)
    if (filters.area) params.append('area', filters.area)
    if (filters.status) params.append('status', filters.status)
    if (filters.limit) params.append('limit', filters.limit)
    if (filters.offset) params.append('offset', filters.offset)
    
    const response = await axios.get(`${API_URL}/manager/facilities/search?${params}`, {
      headers: getAuthHeader()
    })
    return response.data
  },

  async searchFacilitiesForStaff(filters = {}) {
    const params = new URLSearchParams()
    
    if (filters.name) params.append('name', filters.name)
    if (filters.type) params.append('type', filters.type)
    if (filters.area) params.append('area', filters.area)
    if (filters.status) params.append('status', filters.status)
    if (filters.limit) params.append('limit', filters.limit)
    if (filters.offset) params.append('offset', filters.offset)
    
    const response = await axios.get(`${API_URL}/waiter/facilities/search?${params}`, {
      headers: getAuthHeader()
    })
    return response.data
  },
}
