import api from './api'

export const shiftService = {
  async startShift(shiftData) {
    const response = await api.post('/waiter/shifts/start', shiftData)
    return response.data
  },

  async endShift(id, endCash) {
    const response = await api.post(`/waiter/shifts/${id}/end`, { end_cash: endCash })
    return response.data
  },

  async closeShift(id, endCash) {
    const response = await api.post(`/cashier/shifts/${id}/close`, { end_cash: endCash })
    return response.data
  },

  async getCurrentShift() {
    const response = await api.get('/waiter/shifts/current')
    return response.data
  },

  async getMyShifts() {
    const response = await api.get('/waiter/shifts')
    return response.data
  },

  async getAllShifts() {
    const response = await api.get('/cashier/shifts')
    return response.data
  },

  async getShift(id) {
    const response = await api.get(`/cashier/shifts/${id}`)
    return response.data
  }
}
