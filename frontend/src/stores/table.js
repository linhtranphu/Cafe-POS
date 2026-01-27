import { defineStore } from 'pinia'
import { tableService } from '../services/table'

export const useTableStore = defineStore('table', {
  state: () => ({
    tables: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchTables() {
      this.loading = true
      this.error = null
      try {
        this.tables = await tableService.getTables() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải bàn'
        this.tables = []
      } finally {
        this.loading = false
      }
    },

    async createTable(tableData) {
      this.error = null
      try {
        const table = await tableService.createTable(tableData)
        this.tables.push(table)
        return table
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tạo bàn'
        throw error
      }
    },

    async updateTable(id, tableData) {
      this.error = null
      try {
        const table = await tableService.updateTable(id, tableData)
        const index = this.tables.findIndex(t => t.id === id)
        if (index !== -1) {
          this.tables[index] = table
        }
        return table
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi cập nhật bàn'
        throw error
      }
    },

    async deleteTable(id) {
      this.error = null
      try {
        await tableService.deleteTable(id)
        this.tables = this.tables.filter(t => t.id !== id)
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi xóa bàn'
        throw error
      }
    }
  },

  getters: {
    emptyTables: (state) => {
      return state.tables.filter(t => t.status === 'EMPTY')
    },

    occupiedTables: (state) => {
      return state.tables.filter(t => t.status === 'OCCUPIED')
    },

    tablesByArea: (state) => (area) => {
      return state.tables.filter(t => t.area === area)
    }
  }
})
