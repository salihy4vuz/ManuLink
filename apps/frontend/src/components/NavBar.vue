<template>
  <nav class="glass-navbar">
    <div class="nav-container">
      
      <!-- Brand Logo -->
      <router-link to="/" class="logo">
        <Hexagon class="icon lucide-icon" size="28" />
        <span class="brand-text">Apex</span>
        <span class="brand-sub">FactoryOS</span>
      </router-link>

      <!-- Navigation Links -->
      <div class="nav-links">
        <router-link to="/home" class="nav-item" exact-active-class="active">Portal</router-link>
        <router-link to="/quote" class="nav-item" active-class="active">New Quote</router-link>
        <router-link to="/orders" class="nav-item" active-class="active">My Orders</router-link>
        <router-link to="/operator" class="nav-item" active-class="active">Operator</router-link>
        <router-link to="/dashboard" class="nav-item" active-class="active">Analytics</router-link>
        <router-link to="/admin" class="nav-item" active-class="active">Telemetry</router-link>
        <router-link to="/support" class="nav-item" active-class="active">Support</router-link>
        <router-link to="/subscription" class="nav-item" active-class="active">Upgrade</router-link>
      </div>

      <!-- User Profile & Auth -->
      <div class="user-profile">
        <div class="user-info">
          <div class="avatar">A</div>
          <span class="tenant-id">admin-tenant</span>
        </div>
        <button @click="handleLogout" class="btn-logout">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
        </button>
      </div>

    </div>
  </nav>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Hexagon } from 'lucide-vue-next'

const router = useRouter()
const auth = useAuthStore()

const handleLogout = () => {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.glass-navbar {
  background: #fff;
  border-bottom: 1px solid #e2e8f0;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.nav-container {
  max-width: 1600px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.8rem 2rem;
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  text-decoration: none;
}

.icon {
  color: #1d6bf3;
}

.brand-text {
  font-size: 1.2rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.5px;
}

.brand-sub {
  font-size: 1.2rem;
  font-weight: 300;
  color: #1d6bf3;
}

.nav-links {
  display: flex;
  gap: 0.25rem;
  background: transparent;
  padding: 0.4rem;
  border-radius: 12px;
}

.nav-item {
  text-decoration: none;
  color: #475569;
  font-weight: 500;
  font-size: 0.9rem;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.nav-item:hover {
  color: #1d6bf3;
  background: #f8fafc;
}

.nav-item.active {
  color: #1d6bf3;
  background: #eff6ff;
  font-weight: 600;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.4rem 1rem 0.4rem 0.4rem;
  background: #f8fafc;
  border-radius: 20px;
  border: 1px solid #e2e8f0;
}

.avatar {
  width: 32px;
  height: 32px;
  background: #1d6bf3;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  color: white;
  font-size: 0.9rem;
}

.tenant-id {
  font-size: 0.85rem;
  color: #475569;
  font-weight: 600;
}

.btn-logout {
  background: #fef2f2;
  color: #ef4444;
  border: 1px solid #fca5a5;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-logout:hover {
  background: #fee2e2;
  border-color: #f87171;
  transform: translateY(-2px);
}
</style>
