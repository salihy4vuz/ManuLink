<template>
  <div class="page-container">
    <NavBar />
    
    <div class="content-wrapper">
      <div class="header">
        <h2 class="title-gradient">Order Tracking & History</h2>
        <p class="text-muted">Enterprise Data Grid for real-time sheet metal production monitoring.</p>
      </div>

      <div class="glass-panel data-grid-container">
        <!-- Toolbar -->
        <div class="grid-toolbar">
          <div class="search-box">
            <Search class="lucide-icon" size="18" />
            <input type="text" v-model="searchQuery" placeholder="Search Order ID / Part..." />
          </div>
          <div class="filter-box">
            <Filter class="lucide-icon" size="18" />
            <select v-model="statusFilter">
              <option value="All">All Statuses</option>
              <option value="Pending">Pending</option>
              <option value="In Production">In Production</option>
              <option value="Completed">Completed</option>
            </select>
          </div>
          <button class="btn-outline-primary" @click="exportCSV">
            <Download class="lucide-icon" size="16" /> Export CSV
          </button>
        </div>

        <!-- Data Table -->
        <div class="table-responsive">
          <table class="enterprise-table">
            <thead>
              <tr>
                <th>Order ID</th>
                <th>Part Name</th>
                <th>Date</th>
                <th>Status</th>
                <th>Progress</th>
                <th>Total (USD)</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="filteredOrders.length === 0">
                <td colspan="7" class="empty-state">No orders match your criteria.</td>
              </tr>
              <tr v-for="order in filteredOrders" :key="order.ID" class="table-row">
                <td class="order-id mono">#{{ String(order.ID).padStart(6, '0') }}</td>
                <td class="part-name">{{ order.PartName }}</td>
                <td>{{ new Date(order.CreatedAt).toLocaleDateString() }}</td>
                <td>
                  <span class="status-badge" :class="getStatusClass(order.Status)">
                    {{ order.Status }}
                  </span>
                </td>
                <td class="progress-cell">
                  <div class="progress-bar">
                    <div class="progress-fill" 
                         :style="{ width: getProgress(order.Status) + '%' }" 
                         :class="getStatusClass(order.Status)">
                    </div>
                  </div>
                </td>
                <td class="price mono">${{ order.TotalPriceUSD.toFixed(2) }}</td>
                <td class="actions">
                  <button 
                    class="btn-neon outline small" 
                    v-if="order.Status === 'Completed'"
                    @click="downloadPDF(order.ID)"
                    title="Download Proforma Invoice"
                  >
                    PDF
                  </button>
                  <span v-else class="text-muted small-text">Processing</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import NavBar from '../components/NavBar.vue'
import api from '../api/client.js'
import { Search, Filter, Download } from 'lucide-vue-next'

const orders = ref([])
const searchQuery = ref('')
const statusFilter = ref('All')
let ws = null

const fetchOrders = async () => {
  try {
    const res = await api.get('/orders')
    const data = res.data
    if (data.success || data.orders) {
      orders.value = data.orders || data
    }
  } catch (err) {
    console.error('Failed to fetch orders:', err)
  }
}

const connectWebSocket = () => {
  const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:3000/ws'
  const token = localStorage.getItem('token') || ''
  ws = new WebSocket(`${wsUrl}?token=${token}`)
  ws.onmessage = (event) => {
    const payload = JSON.parse(event.data)
    if (payload.event === 'order.status_updated') {
      const updatedOrder = payload.data
      const idx = orders.value.findIndex(o => o.ID === updatedOrder.ID)
      if (idx !== -1) {
        orders.value[idx].Status = updatedOrder.Status
      } else {
        orders.value.unshift(updatedOrder)
      }
    }
  }
  ws.onclose = () => {
    setTimeout(connectWebSocket, 3000)
  }
}

onMounted(() => {
  fetchOrders()
  connectWebSocket()
})

onUnmounted(() => {
  if (ws) ws.close()
})

// Advanced Filtering Logic
const filteredOrders = computed(() => {
  let result = orders.value

  if (statusFilter.value !== 'All') {
    result = result.filter(o => o.Status === statusFilter.value)
  }

  if (searchQuery.value.trim() !== '') {
    const term = searchQuery.value.toLowerCase()
    result = result.filter(o => 
      o.PartName.toLowerCase().includes(term) || 
      o.ID.toString().includes(term)
    )
  }

  return result
})

const downloadPDF = async (orderId) => {
  try {
    const res = await api.get(`/quotes/${orderId}/pdf`, {
      responseType: 'blob'
    })
    const blob = res.data
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `Quote_${orderId}.pdf`
    document.body.appendChild(a)
    a.click()
    a.remove()
    window.URL.revokeObjectURL(url)
  } catch (err) {
    alert('Failed to download invoice PDF.')
  }
}

const exportCSV = () => {
  let csvContent = "data:text/csv;charset=utf-8,Order ID,Part Name,Status,Total USD\n"
  filteredOrders.value.forEach(order => {
    csvContent += `${order.ID},${order.PartName},${order.Status},${order.TotalPriceUSD}\n`
  })
  const encodedUri = encodeURI(csvContent)
  const a = document.createElement("a")
  a.setAttribute("href", encodedUri)
  a.setAttribute("download", "orders_export.csv")
  document.body.appendChild(a)
  a.click()
  a.remove()
}

const getStatusClass = (status) => {
  switch (status) {
    case 'Pending': return 'status-warning'
    case 'In Production': return 'status-info'
    case 'Completed': return 'status-success'
    default: return 'status-default'
  }
}

const getProgress = (status) => {
  switch (status) {
    case 'Pending': return 10
    case 'In Production': return 60
    case 'Completed': return 100
    default: return 0
  }
}
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
  background: transparent;
}

.title-gradient {
  color: var(--accent-neon);
  font-size: 2.2rem;
  font-weight: 800;
  margin-bottom: 0.5rem;
  letter-spacing: -0.5px;
  text-shadow: 0 0 10px var(--accent-glow);
}

.header {
  margin-bottom: 30px;
}

.text-muted {
  color: var(--text-muted);
}

/* .glass-panel uses global CSS */

.data-grid-container {
  padding: 1.5rem;
  max-width: 1400px;
}

/* Toolbar */
.grid-toolbar {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
}

.search-box, .filter-box {
  display: flex;
  align-items: center;
  background: rgba(0,0,0,0.2);
  border: 1px solid var(--border-glass);
  border-radius: 8px;
  padding: 0.5rem 1rem;
}

.search-box input, .filter-box select {
  background: transparent;
  border: none;
  color: var(--text-main);
  outline: none;
  padding-left: 0.5rem;
  font-size: 0.90rem;
}

.search-box input { width: 300px; }

.lucide-icon {
  color: var(--text-muted);
}

.btn-outline-primary {
  background: transparent;
  border: 1px solid var(--accent-neon);
  color: var(--accent-neon);
  border-radius: 6px;
  padding: 0.4rem 1rem;
  cursor: pointer;
  transition: all 0.2s ease;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-left: auto;
  font-size: 0.85rem;
}

.btn-outline-primary:hover {
  background: rgba(0, 229, 255, 0.1);
}

/* Table */
.table-responsive {
  overflow-x: auto;
}

.enterprise-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.enterprise-table th {
  padding: 0.75rem 1rem;
  color: var(--text-muted);
  font-size: 0.75rem;
  text-transform: uppercase;
  font-weight: 700;
  letter-spacing: 0.5px;
  border-bottom: 1px solid var(--border-glass);
}

.enterprise-table td {
  padding: 1rem;
  color: var(--text-main);
  border-bottom: 1px solid var(--border-glass);
  vertical-align: middle;
  font-size: 0.9rem;
}

.table-row:hover td {
  background: rgba(255, 255, 255, 0.05);
}

.order-id { font-family: 'Fira Code', monospace; color: var(--accent-neon); font-weight: 600; }
.part-name { font-weight: 600; color: var(--text-main); }
.price { font-weight: 700; color: var(--accent-success); font-family: 'Fira Code', monospace; }
.empty-state { text-align: center; padding: 3rem !important; color: #64748b; font-style: italic; }

/* Status Badges */
.status-badge {
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.status-warning { background: rgba(245, 158, 11, 0.2); color: #f59e0b; border: 1px solid rgba(245, 158, 11, 0.5); }
.status-info { background: rgba(0, 229, 255, 0.2); color: var(--accent-neon); border: 1px solid rgba(0, 229, 255, 0.5); }
.status-success { background: rgba(16, 185, 129, 0.2); color: var(--accent-success); border: 1px solid rgba(16, 185, 129, 0.5); }

/* Progress - Industrial Gauge Style */
.progress-cell { width: 120px; }
.progress-bar {
  height: 6px;
  background: rgba(255,255,255,0.1);
  border-radius: 4px;
  overflow: hidden;
}
.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 1s ease;
}
.progress-fill.status-warning { background: #f59e0b; }
.progress-fill.status-info { background: var(--accent-neon); }
.progress-fill.status-success { background: var(--accent-success); }

/* Buttons */
.btn-neon.outline {
  background: transparent;
  border: 1px solid var(--accent-neon);
  color: var(--accent-neon);
  border-radius: 6px;
  transition: all 0.2s ease;
  cursor: pointer;
  font-weight: 600;
}
.btn-neon.outline:hover {
  background: rgba(0, 229, 255, 0.1);
}
.btn-neon.small { padding: 4px 12px; font-size: 0.8rem; }
.small-text { font-size: 0.8rem; color: var(--text-muted); }
</style>
