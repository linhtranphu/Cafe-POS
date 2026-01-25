import axios from 'axios'

const API_URL = 'http://localhost:8080/api'

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Add token to requests
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

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

export default api