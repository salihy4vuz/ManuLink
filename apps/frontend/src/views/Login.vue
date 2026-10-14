<template>
  <div class="login-container">
    <div class="glass-panel login-card">
      <div class="logo-area">
        <h2>Apex Factory OS</h2>
        <p class="text-accent">Autonomous B2B Portal</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="input-group">
          <label>Email Address</label>
          <input type="email" v-model="email" placeholder="john@acmecorp.com" required />
        </div>

        <div class="input-group">
          <label>Password</label>
          <input type="password" v-model="password" placeholder="••••••••" required />
        </div>

        <button type="submit" class="btn-neon full-width" :disabled="loading">
          {{ loading ? 'AUTHENTICATING...' : 'ACCESS SECURE PORTAL' }}
        </button>
        <div v-if="errorMsg" class="error-text">{{ errorMsg }}</div>
      </form>

      <div class="footer-links">
        <a href="#" class="text-muted">Forgot Password?</a>
        <a href="#" class="text-muted">Request B2B Account</a>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMsg = ref('')

const handleLogin = async () => {
  loading.value = true
  errorMsg.value = ''
  
  const success = await authStore.login(email.value, password.value)
  loading.value = false
  
  if (success) {
    router.push('/dashboard')
  } else {
    errorMsg.value = 'Invalid credentials or API error.'
  }
}
</script>

<style scoped>
.error-text {
  color: #ef4444;
  font-size: 0.85rem;
  text-align: center;
  margin-top: 10px;
}
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
  background: #f8fafc;
}

.login-card {
  width: 400px;
  padding: 40px;
  display: flex;
  flex-direction: column;
  gap: 30px;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

.logo-area {
  text-align: center;
}

.logo-area h2 {
  font-size: 2rem;
  margin-bottom: 5px;
  color: #0f172a;
  font-weight: 800;
}

.text-accent {
  color: #1d6bf3;
  font-size: 0.9rem;
  letter-spacing: 2px;
  text-transform: uppercase;
  font-weight: 600;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.input-group label {
  font-size: 0.85rem;
  color: #475569;
  font-weight: 600;
}

.input-group input {
  background: #fff;
  border: 1px solid #cbd5e1;
  padding: 12px 15px;
  border-radius: 8px;
  color: #0f172a;
  outline: none;
  font-size: 1rem;
  transition: all 0.2s ease;
}

.input-group input:focus {
  border-color: #1d6bf3;
  box-shadow: 0 0 0 3px rgba(29, 107, 243, 0.1);
}

.btn-neon {
  background: #1d6bf3;
  color: white;
  border: none;
  padding: 14px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}
.btn-neon:hover:not(:disabled) {
  background: #1557d0;
  box-shadow: 0 4px 12px rgba(29, 107, 243, 0.3);
}

.full-width {
  width: 100%;
  margin-top: 10px;
}

.footer-links {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
}

.footer-links a {
  text-decoration: none;
  transition: color 0.3s ease;
  color: #64748b;
}

.footer-links a:hover {
  color: #1d6bf3;
}
</style>
