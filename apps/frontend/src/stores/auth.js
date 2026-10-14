import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api/client.js'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('jwt_token') || null)
  const isAuthenticated = ref(!!token.value)

  async function login(username, password) {
    try {
      const res = await api.post('/auth/login', { username, password })
      const data = res.data
      if (data.token) {
        token.value = data.token
        isAuthenticated.value = true
        localStorage.setItem('jwt_token', data.token)
        return true
      }
      return false
    } catch (err) {
      console.error('Login error:', err)
      return false
    }
  }

  function logout() {
    isAuthenticated.value = false
    token.value = null
    localStorage.removeItem('jwt_token')
  }

  return { isAuthenticated, token, login, logout }
})
