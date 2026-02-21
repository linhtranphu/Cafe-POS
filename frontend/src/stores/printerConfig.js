import { defineStore } from 'pinia'
import { printerConfigService } from '../services/printerConfig'

export const usePrinterConfigStore = defineStore('printerConfig', {
  state: () => ({
    printers: [],
    loading: false,
    error: null,
    testingConnection: false,
    testResult: null
  }),

  getters: {
    defaultBillPrinter: (state) => {
      return state.printers.find(p => p.type === 'BILL' && p.is_default) || null
    },

    defaultLabelPrinter: (state) => {
      return state.printers.find(p => p.type === 'LABEL' && p.is_default) || null
    },

    printersByType: (state) => (type) => {
      return state.printers.filter(p => p.type === type)
    },

    enabledPrinters: (state) => {
      return state.printers.filter(p => p.is_enabled)
    },

    billPrinters: (state) => {
      return state.printers.filter(p => p.type === 'BILL')
    },

    labelPrinters: (state) => {
      return state.printers.filter(p => p.type === 'LABEL')
    },

    getPrinterById: (state) => (id) => {
      return state.printers.find(p => p.id === id) || null
    }
  },

  actions: {
    async fetchPrinters(params = {}) {
      this.loading = true
      this.error = null
      try {
        this.printers = await printerConfigService.fetchPrinters(params)
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải danh sách máy in'
        this.printers = []
      } finally {
        this.loading = false
      }
    },

    async fetchPrinterById(id) {
      this.loading = true
      this.error = null
      try {
        const printer = await printerConfigService.fetchPrinterById(id)
        const index = this.printers.findIndex(p => p.id === id)
        if (index !== -1) {
          this.printers[index] = printer
        } else {
          this.printers.push(printer)
        }
        return printer
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải thông tin máy in'
        throw error
      } finally {
        this.loading = false
      }
    },

    async savePrinter(printer) {
      this.error = null
      try {
        if (printer.id) {
          // Update existing printer
          const updated = await printerConfigService.updatePrinter(printer.id, printer)
          const index = this.printers.findIndex(p => p.id === printer.id)
          if (index !== -1) {
            this.printers[index] = updated
          }
          return updated
        } else {
          // Create new printer
          const created = await printerConfigService.createPrinter(printer)
          this.printers.push(created)
          return created
        }
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi lưu cấu hình máy in'
        throw error
      }
    },

    async deletePrinter(id) {
      this.error = null
      try {
        await printerConfigService.deletePrinter(id)
        this.printers = this.printers.filter(p => p.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa máy in'
        return false
      }
    },

    async testConnection(id) {
      this.testingConnection = true
      this.testResult = null
      this.error = null
      try {
        const result = await printerConfigService.testConnection(id)
        this.testResult = {
          success: true,
          message: result.message || 'Kết nối thành công',
          status: result.status
        }
        return this.testResult
      } catch (error) {
        this.testResult = {
          success: false,
          message: error.response?.data?.error || 'Không thể kết nối đến máy in'
        }
        this.error = this.testResult.message
        return this.testResult
      } finally {
        this.testingConnection = false
      }
    },

    clearTestResult() {
      this.testResult = null
    },

    clearError() {
      this.error = null
    }
  }
})
