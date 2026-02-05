import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor to handle 401 errors
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      const token = localStorage.getItem('token')
      if (token) {
        // Clear auth and redirect to login
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        
        // Show message
        alert('Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.')
        
        // Redirect to login
        window.location.href = '/#/login'
      }
    }
    return Promise.reject(error)
  }
)

export default api