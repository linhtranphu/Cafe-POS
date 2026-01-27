import api from './api'

export const facilityService = {
  async getFacilities() {
    const response = await api.get('/manager/facilities')
    return response.data
  },

  async getFacility(id) {
    const response = await api.get(`/manager/facilities/${id}`)
    return response.data
  },

  async createFacility(facility) {
    const response = await api.post('/manager/facilities', facility)
    return response.data
  },

  async updateFacility(id, facility) {
    const response = await api.put(`/manager/facilities/${id}`, facility)
    return response.data
  },

  // FR-FM-03: Update facility information with change tracking
  async updateFacilityWithHistory(id, updates) {
    const response = await api.put(`/manager/facilities/${id}/update-with-history`, updates)
    return response.data
  },

  async moveAsset(id, newArea, reason) {
    const response = await api.patch(`/manager/facilities/${id}/move`, {
      new_area: newArea,
      reason: reason
    })
    return response.data
  },

  async checkDeletionEligibility(id) {
    const response = await api.get(`/manager/facilities/${id}/can-delete`)
    return response.data
  },

  async deleteFacility(id) {
    const response = await api.delete(`/manager/facilities/${id}`)
    return response.data
  },

  // FR-FM-01: Asset tracking and inventory
  async getFacilityHistory(id) {
    const response = await api.get(`/manager/facilities/${id}/history`)
    return response.data
  },

  // FR-FM-04: Status management
  async getStatusAlerts() {
    const response = await api.get('/manager/facilities/status-alerts')
    return response.data
  },

  async getStatusHistory(id) {
    console.log('getStatusHistory called with id:', id)
    console.log('Auth token:', localStorage.getItem('token'))
    try {
      const response = await api.get(`/manager/facilities/${id}/status-history`)
      console.log('getStatusHistory response:', response.data)
      return response.data
    } catch (error) {
      console.error('getStatusHistory error:', error)
      throw error
    }
  },

  async updateStatusWithDetails(id, statusData) {
    const response = await api.patch(`/manager/facilities/${id}/status-detailed`, statusData)
    return response.data
  },

  async resolveStatusAlert(alertId) {
    const response = await api.patch(`/manager/facilities/alerts/${alertId}/resolve`, {})
    return response.data
  },

  async getFacilitiesByStatus(status) {
    const response = await api.get(`/manager/facilities/status/${status}`)
    return response.data
  },

  async getFacilitiesByArea(area) {
    const response = await api.get(`/manager/facilities/area/${area}`)
    return response.data
  },

  // FR-FM-02: Maintenance scheduling and tracking
  async getMaintenanceHistory(id) {
    const response = await api.get(`/manager/facilities/${id}/maintenance`)
    return response.data
  },

  // FR-FM-06: Enhanced maintenance management
  async getMaintenanceStats(facilityId) {
    try {
      const response = await api.get(`/manager/facilities/${facilityId}/maintenance-stats`)
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
      const response = await api.get(`/manager/facilities/${facilityId}/next-maintenance`)
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
    const response = await api.put(`/manager/maintenance/${recordId}`, recordData)
    return response.data
  },

  async getMaintenanceCostAnalysis(facilityId, period = '12months') {
    const response = await api.get(`/manager/facilities/${facilityId}/maintenance-costs`, {
      params: { period }
    })
    return response.data
  },

  async scheduleRecurringMaintenance(facilityId, schedule) {
    const response = await api.post(`/manager/facilities/${facilityId}/recurring-maintenance`, schedule)
    return response.data
  },

  async scheduleMaintenanceTask(task) {
    const response = await api.post('/manager/maintenance/schedule', task)
    return response.data
  },

  async getScheduledMaintenance() {
    try {
      const response = await api.get('/manager/maintenance/scheduled')
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
    const response = await api.patch(`/manager/maintenance/${id}`, updates)
    return response.data
  },

  async getMaintenanceDue() {
    try {
      const response = await api.get('/manager/maintenance/due')
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
    const response = await api.get('/manager/issue-reports')
    return response.data
  },

  async createMaintenanceRecord(record) {
    const response = await api.post('/manager/maintenance', record)
    return response.data
  },

  // FR-FM-05: Issue reporting
  async createIssueReport(reportData) {
    const response = await api.post('/staff/issue-reports', reportData)
    return response.data
  },

  async updateIssueReportStatus(reportId, status) {
    const response = await api.patch(`/manager/issue-reports/${reportId}/status`, {
      status
    })
    return response.data
  },

  async addReportComment(reportId, comment) {
    const response = await api.post(`/manager/issue-reports/${reportId}/comments`, {
      comment
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
    
    const response = await api.get(`/manager/facilities/search?${params}`)
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
    
    const response = await api.get(`/waiter/facilities/search?${params}`)
    return response.data
  },
}