import { defineStore } from 'pinia'
import { shiftService } from '../services/shift'

export const useShiftStore = defineStore('shift', {
  state: () => ({
    currentShift: null,
    shifts: [],
    loading: false,
    error: null
  }),

  actions: {
    async startShift(shiftData) {
      this.error = null
      try {
        const shift = await shiftService.startShift(shiftData)
        this.currentShift = shift
        return shift
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi mở ca'
        throw error
      }
    },

    async endShift(id, endCash) {
      this.error = null
      try {
        const shift = await shiftService.endShift(id, endCash)
        this.currentShift = null
        return shift
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi kết ca'
        throw error
      }
    },

    async closeShift(id, endCash) {
      this.error = null
      try {
        const shift = await shiftService.closeShift(id, endCash)
        return shift
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi chốt ca'
        throw error
      }
    },

    async fetchCurrentShift() {
      this.loading = true
      this.error = null
      try {
        this.currentShift = await shiftService.getCurrentShift()
      } catch (error) {
        this.currentShift = null
        if (error.response?.data?.error === 'no open shift found') {
          this.error = null
        } else {
          this.error = error.response?.data?.error || 'Lỗi tải ca'
        }
      } finally {
        this.loading = false
      }
    },

    async fetchMyShifts() {
      this.loading = true
      this.error = null
      try {
        this.shifts = await shiftService.getMyShifts() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải shifts'
        this.shifts = []
      } finally {
        this.loading = false
      }
    },

    async fetchAllShifts() {
      this.loading = true
      this.error = null
      try {
        this.shifts = await shiftService.getAllShifts() || []
      } catch (error) {
        this.error = error.response?.data?.error || 'Lỗi tải shifts'
        this.shifts = []
      } finally {
        this.loading = false
      }
    }
  },

  getters: {
    hasOpenShift: (state) => {
      return state.currentShift !== null && state.currentShift.status === 'OPEN'
    },

    openShifts: (state) => {
      return state.shifts.filter(s => s.status === 'OPEN')
    },

    closedShifts: (state) => {
      return state.shifts.filter(s => s.status === 'CLOSED')
    }
  }
})
