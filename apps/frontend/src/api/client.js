import axios from 'axios'
import { useAuthStore } from '../stores/auth'

// Create an Axios instance to connect to the Golang Fiber Gateway
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Intercept requests to dynamically inject JWT token
api.interceptors.request.use(config => {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers['Authorization'] = `Bearer ${authStore.token}`
  }
  return config
}, error => {
  return Promise.reject(error)
})

// Intercept responses for global error handling
api.interceptors.response.use(
  response => response,
  error => {
    // Auto logout on 401 Unauthorized
    if (error.response && error.response.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
    }
    console.error('API Error:', error.response?.data?.error || error.message)
    return Promise.reject(error)
  }
)

export default api
