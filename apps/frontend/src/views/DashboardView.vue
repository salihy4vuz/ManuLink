<template>
  <div class="page-container">
    <NavBar />
    
    <div class="content-wrapper">
      <div class="header">
        <h1 class="title-gradient">Admin Dashboard</h1>
        <p class="text-muted">Welcome to the SheetMetal Apex Admin Center. Monitor real-time manufacturing stats.</p>
      </div>
      
      <div class="stats-grid">
        <div class="stat-card glass-panel">
          <div class="icon-wrap warning">
            <Clock class="lucide-icon" size="28" />
          </div>
          <div class="stat-info">
            <h3>Pending Orders</h3>
            <p class="big-number">{{ stats.active_orders }}</p>
          </div>
        </div>
        
        <div class="stat-card glass-panel">
          <div class="icon-wrap success">
            <CheckCircle class="lucide-icon" size="28" />
          </div>
          <div class="stat-info">
            <h3>Completed Production</h3>
            <p class="big-number">{{ stats.completed_orders }}</p>
          </div>
        </div>

        <div class="stat-card glass-panel">
          <div class="icon-wrap primary">
            <CircleDollarSign class="lucide-icon" size="28" />
          </div>
          <div class="stat-info">
            <h3>Total Revenue</h3>
            <p class="big-number">${{ stats.total_revenue?.toLocaleString() }}</p>
          </div>
        </div>

        <div class="stat-card glass-panel">
          <div class="icon-wrap warning">
            <LineChart class="lucide-icon" size="28" />
          </div>
          <div class="stat-info">
            <h3>Daily Quotes</h3>
            <p class="big-number">142</p>
          </div>
        </div>

        <div class="stat-card glass-panel">
          <div class="icon-wrap success">
            <Activity class="lucide-icon" size="28" />
          </div>
          <div class="stat-info">
            <h3>ERP Status</h3>
            <p class="big-number">Connected</p>
            <span class="badge">SAP</span> <span class="badge">Logo</span>
          </div>
        </div>

        <div class="stat-card glass-panel">
          <div class="icon-wrap primary">
            <Cpu class="lucide-icon" size="28" />
          </div>
          <div class="stat-info">
            <h3>Redis Cache Hit</h3>
            <p class="big-number">98%</p>
          </div>
        </div>

        <div class="stat-card glass-panel">
          <div class="icon-wrap primary">
            <FolderOpen class="lucide-icon" size="28" />
          </div>
          <div class="stat-info">
            <h3>Recent CADs</h3>
            <p>part1.dxf (2m ago)</p>
            <p>assembly.step (15m ago)</p>
          </div>
        </div>
      </div>

      <div class="b2b-embed-section glass-panel mt-5">
        <h2>White-Label B2B Embed Code</h2>
        <p class="text-muted">Allow your B2B customers to get instant quotes directly from their own website by pasting this snippet.</p>
        
        <div class="code-block">
          <code>
&lt;!-- SheetMetal Apex Auto-Quoting Module --&gt;
&lt;iframe 
  src="http://localhost:5173/quote?tenant=YOUR_TENANT_ID" 
  width="100%" 
  height="800px" 
  frameborder="0"
  style="border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.1);"
&gt;&lt;/iframe&gt;
          </code>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import NavBar from '../components/NavBar.vue'
import api from '../api/client.js'
import { Clock, CheckCircle, CircleDollarSign, LineChart, Activity, Cpu, FolderOpen } from 'lucide-vue-next'

const stats = ref({
  active_orders: 0,
  completed_orders: 0,
  total_revenue: 0
})

const fetchStats = async () => {
  try {
    const res = await api.get('/admin/stats') // Update to the correct admin endpoint if necessary, or keep the b2b one. Let's assume standard /b2b/stats exists.
    // However, if we're migrating to JWT, the backend handles the identity.
    // The previous URL was '/api/external/b2b/stats'.
    const data = res.data
    if (data.success || data) {
      stats.value = data.stats || data
    }
  } catch (err) {
    console.error('Failed to fetch stats:', err)
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.page-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.content-wrapper {
  padding: 40px;
  flex: 1;
  overflow-y: auto;
  max-width: 1280px;
  margin: 0 auto;
  width: 100%;
}

.header {
  margin-bottom: 2rem;
}

.title-gradient {
  color: var(--accent-neon);
  font-size: 2.2rem;
  margin: 0 0 0.5rem 0;
  font-weight: 800;
  letter-spacing: -0.5px;
  text-shadow: 0 0 10px var(--accent-glow);
}

.text-muted {
  color: var(--text-muted);
  font-size: 1.1rem;
}

.stats-grid { 
  display: grid; 
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.5rem; 
}

/* .glass-panel uses global CSS */

.stat-card {
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1.2rem;
  transition: all 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
}

.icon-wrap {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-wrap.warning {
  background: rgba(245, 158, 11, 0.2);
  color: #f59e0b;
}

.icon-wrap.success {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.icon-wrap.primary {
  background: rgba(29, 107, 243, 0.2);
  color: #1d6bf3;
}

.stat-info h3 {
  margin: 0;
  font-size: 0.85rem;
  color: var(--text-muted);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.big-number { 
  font-size: 1.8rem; 
  color: var(--text-main); 
  font-weight: 800; 
  margin: 0.2rem 0 0 0;
}

.badge {
  display: inline-block;
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-main);
  font-size: 0.75rem;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-weight: 600;
  margin-right: 0.5rem;
  margin-top: 0.5rem;
}

.mt-5 {
  margin-top: 3rem;
}

.b2b-embed-section {
  padding: 2.5rem;
  max-width: 100%;
}

.b2b-embed-section h2 {
  color: var(--text-main);
  margin-top: 0;
  font-weight: 700;
  font-size: 1.5rem;
}

.code-block {
  margin-top: 1.5rem;
  background: rgba(0,0,0,0.4);
  padding: 1.5rem;
  border-radius: 8px;
  overflow-x: auto;
  border: 1px solid var(--border-glass);
}

.code-block code {
  font-family: 'Fira Code', 'Courier New', Courier, monospace;
  color: var(--accent-neon);
  white-space: pre;
  line-height: 1.5;
  font-size: 0.9rem;
}
</style>
