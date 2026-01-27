import api from './api'

export const tableService = {
  async getTables() {
    const response = await api.get('/waiter/tables')
    return response.data
  },

  async createTable(table) {
    const response = await api.post('/manager/tables', table)
    return response.data
  },

  async updateTable(id, table) {
    const response = await api.put(`/manager/tables/${id}`, table)
    return response.data
  },

  async deleteTable(id) {
    await api.delete(`/manager/tables/${id}`)
  },

  async getTable(id) {
    const response = await api.get(`/manager/tables/${id}`)
    return response.data
  }
}
