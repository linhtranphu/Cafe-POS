import api from './api'

export const authService = {
  async login(username, password) {
    const response = await api.post('/login', { username, password })
    return response.data
  }
}

export const login = async (credentials) => {
  const response = await api.post('/login', credentials)
  return response.data
}
