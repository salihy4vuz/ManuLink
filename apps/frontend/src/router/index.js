import { createRouter, createWebHistory } from 'vue-router'
import LandingView from '../views/LandingView.vue'
import Login from '../views/Login.vue'
import Home from '../views/Home.vue'
import DashboardView from '../views/DashboardView.vue'
import Orders from '../views/Orders.vue'
import Admin from '../views/Admin.vue'
import SettingsView from '../views/SettingsView.vue'
import NestingView from '../views/NestingView.vue'
import QuoteView from '../views/QuoteView.vue'
import OperatorDashboard from '../views/OperatorDashboard.vue'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    name: 'Landing',
    component: LandingView
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/home',
    name: 'Home',
    component: Home
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: { requiresAuth: true }
  },
  {
    path: '/orders',
    name: 'Orders',
    component: Orders,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: Admin,
    meta: { requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: SettingsView
  },
  {
    path: '/nesting',
    name: 'Nesting',
    component: NestingView
  },
  {
    path: '/quote',
    name: 'Quote',
    component: QuoteView
  },
  {
    path: '/operator',
    name: 'OperatorDashboard',
    component: OperatorDashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/support',
    name: 'Support',
    component: () => import('../views/SupportView.vue')
  },
  {
    path: '/subscription',
    name: 'Subscription',
    component: () => import('../views/SubscriptionView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation Guard for JWT protection
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'Login' })
  } else {
    next()
  }
})

export default router
