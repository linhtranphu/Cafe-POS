import { defineStore } from 'pinia'
import { facilityService } from '../services/facility'

export const useFacilityStore = defineStore('facility', {
  state: () => ({
    items: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchFacilities() {
      this.loading = true
      this.error = null
      try {
        this.items = await facilityService.getFacilities() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải cơ sở vật chất'
        this.items = []
      } finally {
        this.loading = false
      }
    },

    async createFacility(facility) {
      this.error = null
      try {
        const newItem = await facilityService.createFacility(facility)
        this.items.push(newItem)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo tài sản'
        return false
      }
    },

    async updateFacility(id, facility) {
      this.error = null
      try {
        const updatedItem = await facilityService.updateFacility(id, facility)
        const index = this.items.findIndex(i => i.id === id)
        if (index !== -1) {
          this.items[index] = updatedItem
        }
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi cập nhật tài sản'
        return false
      }
    },

    // FR-FM-03: Enhanced update with change tracking
    async updateFacilityWithHistory(id, updates) {
      this.error = null
      try {
        const updatedItem = await facilityService.updateFacilityWithHistory(id, updates)
        const index = this.items.findIndex(i => i.id === id)
        if (index !== -1) {
          this.items[index] = updatedItem
        }
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi cập nhật tài sản'
        return false
      }
    },

    async moveAsset(id, newArea, reason) {
      this.error = null
      try {
        const updatedItem = await facilityService.moveAsset(id, newArea, reason)
        const index = this.items.findIndex(i => i.id === id)
        if (index !== -1) {
          this.items[index] = updatedItem
        }
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi di chuyển tài sản'
        return false
      }
    },

    async checkCanDelete(id) {
      try {
        return await facilityService.checkDeletionEligibility(id)
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi kiểm tra quyền xóa'
        return { canDelete: false, reason: 'Lỗi hệ thống' }
      }
    },

    async deleteFacility(id) {
      this.error = null
      try {
        await facilityService.deleteFacility(id)
        this.items = this.items.filter(i => i.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa tài sản'
        return false
      }
    },

    // FR-FM-01: Asset tracking and inventory
    // FR-FM-04: Enhanced status management
    async fetchStatusAlerts() {
      try {
        return await facilityService.getStatusAlerts() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải cảnh báo trạng thái'
        return []
      }
    },

    async fetchStatusHistory(id) {
      console.log('fetchStatusHistory called with id:', id)
      try {
        const result = await facilityService.getStatusHistory(id) || []
        console.log('fetchStatusHistory result:', result)
        return result
      } catch (error) {
        console.error('fetchStatusHistory error:', error)
        this.error = error.response?.data?.error || 'Lỗi tải lịch sử trạng thái'
        return []
      }
    },

    async updateStatusWithDetails(id, statusData) {
      this.error = null
      try {
        const updatedItem = await facilityService.updateStatusWithDetails(id, statusData)
        const index = this.items.findIndex(i => i.id === id)
        if (index !== -1) {
          this.items[index] = updatedItem
        }
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi cập nhật trạng thái'
        return false
      }
    },

    async resolveStatusAlert(alertId) {
      this.error = null
      try {
        await facilityService.resolveStatusAlert(alertId)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi giải quyết cảnh báo'
        return false
      }
    },

    async fetchFacilitiesByStatus(status) {
      try {
        return await facilityService.getFacilitiesByStatus(status) || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải tài sản theo trạng thái'
        return []
      }
    },

    async fetchFacilitiesByArea(area) {
      try {
        return await facilityService.getFacilitiesByArea(area) || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải tài sản theo khu vực'
        return []
      }
    },

    async fetchFacilityHistory(id) {
      try {
        return await facilityService.getFacilityHistory(id) || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải lịch sử'
        return []
      }
    },

    async fetchMaintenanceHistory(id) {
      try {
        return await facilityService.getMaintenanceHistory(id) || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải lịch sử bảo trì'
        return []
      }
    },

    // FR-FM-06: Enhanced maintenance management
    async fetchMaintenanceStats(id) {
      try {
        return await facilityService.getMaintenanceStats(id) || {
          total: 0, totalCost: 0, avgInterval: 0, 
          scheduledCost: 0, emergencyCost: 0
        }
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải thống kê bảo trì'
        return { total: 0, totalCost: 0, avgInterval: 0, scheduledCost: 0, emergencyCost: 0 }
      }
    },

    async fetchNextMaintenanceDate(id) {
      try {
        const result = await facilityService.getNextMaintenanceDate(id)
        return result?.next_date || null
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải ngày bảo trì tiếp theo'
        return null
      }
    },

    async updateMaintenanceRecord(recordId, recordData) {
      this.error = null
      try {
        await facilityService.updateMaintenanceRecord(recordId, recordData)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi cập nhật bảo trì'
        return false
      }
    },

    async fetchMaintenanceCostAnalysis(id, period = '12months') {
      try {
        return await facilityService.getMaintenanceCostAnalysis(id, period) || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải phân tích chi phí'
        return []
      }
    },

    async createMaintenanceRecord(record) {
      this.error = null
      try {
        await facilityService.createMaintenanceRecord(record)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo bảo trì'
        return false
      }
    },

    // FR-FM-02: Maintenance scheduling
    async scheduleMaintenanceTask(task) {
      this.error = null
      try {
        await facilityService.scheduleMaintenanceTask(task)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi lên lịch bảo trì'
        return false
      }
    },

    async fetchScheduledMaintenance() {
      try {
        return await facilityService.getScheduledMaintenance() || []
      } catch (error) {
        console.warn('Scheduled maintenance endpoint not available:', error.message)
        return []
      }
    },

    // FR-FM-05: Issue reporting
    async createIssueReport(reportData) {
      this.error = null
      try {
        await facilityService.createIssueReport(reportData)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo báo cáo sự cố'
        return false
      }
    },

    async fetchIssueReports() {
      try {
        return await facilityService.getIssueReports() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải báo cáo sự cố'
        return []
      }
    },

    async updateIssueReportStatus(reportId, status) {
      this.error = null
      try {
        await facilityService.updateIssueReportStatus(reportId, status)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi cập nhật trạng thái báo cáo'
        return false
      }
    },

    async addReportComment(reportId, comment) {
      this.error = null
      try {
        await facilityService.addReportComment(reportId, comment)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi thêm bình luận'
        return false
      }
    },

    // FR-FM-08: Search & Filter
    async searchFacilities(filters = {}) {
      this.loading = true
      this.error = null
      try {
        const result = await facilityService.searchFacilities(filters)
        return result
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tìm kiếm tài sản'
        return { data: [], total: 0 }
      } finally {
        this.loading = false
      }
    },

    async searchFacilitiesForStaff(filters = {}) {
      this.loading = true
      this.error = null
      try {
        const result = await facilityService.searchFacilitiesForStaff(filters)
        return result
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tìm kiếm tài sản'
        return { data: [], total: 0 }
      } finally {
        this.loading = false
      }
    },

    async fetchMaintenanceDue() {
      try {
        return await facilityService.getMaintenanceDue() || []
      } catch (error) {
        console.warn('Maintenance due endpoint not available:', error.message)
        return []
      }
    }
  }
})