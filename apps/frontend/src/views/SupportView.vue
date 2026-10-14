<template>
  <div class="page-container">
    <NavBar />

    <div class="support-container">
      
      <!-- Left: Ticket List -->
      <div class="ticket-list glass-panel">
        <div class="list-header">
          <h2>My Tickets</h2>
          <span class="badge">{{ tickets.length }} Active</span>
        </div>
        
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>Loading tickets...</p>
        </div>
        
        <div v-else class="tickets">
          <div v-for="ticket in tickets" :key="ticket.id" class="ticket-card" :class="ticket.status.toLowerCase().replace(' ', '-')">
            <div class="ticket-header">
              <span class="ticket-id">#{{ ticket.id }}</span>
              <span class="ticket-status">{{ ticket.status }}</span>
            </div>
            <h4>{{ ticket.subject }}</h4>
            <p class="ticket-date">{{ new Date(ticket.created_at).toLocaleString() }}</p>
            <div class="ticket-priority" :class="ticket.priority.toLowerCase()">
              {{ ticket.priority }} Priority
            </div>
          </div>
          
          <div v-if="tickets.length === 0" class="empty-state">
            <span class="empty-icon">🏖️</span>
            <p>No support tickets found.</p>
            <span class="text-muted">You are all caught up!</span>
          </div>
        </div>
      </div>

      <!-- Right: New Ticket Form -->
      <div class="ticket-form glass-panel">
        <h2 class="title-gradient">Open a New Ticket</h2>
        <p class="subtitle">Our engineering team will investigate your issue immediately.</p>
        
        <form @submit.prevent="submitTicket" class="form-grid">
          <div class="input-group">
            <label>Subject</label>
            <input v-model="form.subject" type="text" placeholder="e.g., DXF bending lines not recognized" required />
          </div>
          
          <div class="input-group">
            <label>Priority</label>
            <select v-model="form.priority">
              <option value="Low">Low (General Inquiry)</option>
              <option value="Medium">Medium (Production Issue)</option>
              <option value="High">High (System Outage / Critical)</option>
            </select>
          </div>

          <div class="input-group">
            <label>Message</label>
            <textarea v-model="form.message" rows="5" placeholder="Describe your issue in detail..." required></textarea>
          </div>

          <button type="submit" class="btn-primary glow-effect" :disabled="submitting">
            <span v-if="submitting" class="spinner-small"></span>
            <span v-else>Submit Ticket</span>
          </button>
        </form>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import NavBar from '../components/NavBar.vue'

const tickets = ref([])
const loading = ref(true)
const submitting = ref(false)

const form = ref({
  subject: '',
  priority: 'Medium',
  message: ''
})

const fetchTickets = async () => {
  loading.value = true
  try {
    const res = await fetch('/api/tickets', {
      headers: { 'X-API-Key': 'admin-key' }
    })
    const data = await res.json()
    if (data.status === 'success') {
      tickets.value = data.tickets
    }
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

const submitTicket = async () => {
  submitting.value = true
  try {
    const res = await fetch('/api/tickets/new', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': 'admin-key'
      },
      body: JSON.stringify(form.value)
    })
    const data = await res.json()
    
    if (data.status === 'success') {
      alert("Ticket submitted successfully!")
      form.value = { subject: '', priority: 'Medium', message: '' }
      await fetchTickets() // Refresh the list
    } else {
      alert("Error: " + data.error)
    }
  } catch (err) {
    alert("Network error.")
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchTickets()
})
</script>

<style scoped>
.page-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.support-container {
  flex: 1;
  display: flex;
  gap: 2rem;
  padding: 2rem;
  background: transparent;
  overflow-y: hidden;
}

.glass-panel {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

.ticket-list {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 2rem;
  overflow-y: auto;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.list-header h2 {
  margin: 0;
  color: #0f172a;
  font-size: 1.8rem;
  font-weight: 800;
}

.badge {
  background: #eff6ff;
  color: #1d6bf3;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  border: 1px solid #bfdbfe;
}

.ticket-form {
  flex: 1;
  padding: 2.5rem;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.title-gradient {
  color: #0f172a;
  font-size: 2rem;
  font-weight: 800;
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.5px;
}

.subtitle { color: #64748b; margin: 0 0 2rem 0; font-size: 1rem; }

/* Ticket Cards */
.tickets { display: flex; flex-direction: column; gap: 1rem; }
.ticket-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 1.5rem;
  border-left: 4px solid #1d6bf3;
  transition: transform 0.2s ease, background 0.2s ease;
}
.ticket-card:hover {
  transform: translateX(5px);
  background: #f1f5f9;
}

.ticket-card.resolved { border-left-color: #10b981; opacity: 0.7; }
.ticket-card.in-progress { border-left-color: #f59e0b; }

.ticket-header { display: flex; justify-content: space-between; margin-bottom: 0.5rem; font-size: 0.85rem; }
.ticket-id { color: #64748b; }
.ticket-status { font-weight: 700; text-transform: uppercase; letter-spacing: 0.05em; }

.resolved .ticket-status { color: #10b981; }
.open .ticket-status { color: #1d6bf3; }
.in-progress .ticket-status { color: #f59e0b; }

.ticket-card h4 { font-size: 1.1rem; margin: 0 0 0.5rem 0; color: #0f172a; font-weight: 700; }
.ticket-date { color: #64748b; font-size: 0.85rem; margin: 0 0 1rem 0; }

.ticket-priority { 
  display: inline-block; 
  padding: 4px 10px; 
  border-radius: 8px; 
  font-size: 0.75rem; 
  font-weight: bold; 
}
.ticket-priority.high { background: #fef2f2; color: #ef4444; border: 1px solid #fca5a5; }
.ticket-priority.medium { background: #fffbeb; color: #f59e0b; border: 1px solid #fde68a; }
.ticket-priority.low { background: #eff6ff; color: #1d6bf3; border: 1px solid #bfdbfe; }

/* Empty & Loading States */
.loading-state, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #94a3b8;
  gap: 1rem;
}
.empty-icon { font-size: 4rem; }
.empty-state p { margin: 0; font-size: 1.2rem; color: #0f172a; font-weight: 600; }
.empty-state .text-muted { font-size: 0.9rem; }

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e2e8f0;
  border-radius: 50%;
  border-top-color: #1d6bf3;
  animation: spin 1s ease-in-out infinite;
}

.spinner-small {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255,255,255,0.3);
  border-radius: 50%;
  border-top-color: #fff;
  animation: spin 1s ease-in-out infinite;
  display: inline-block;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Form Styles */
.form-grid {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.input-group { display: flex; flex-direction: column; gap: 0.5rem; }
label { color: #475569; font-weight: 600; font-size: 0.9rem; }

input, select, textarea {
  background: #fff;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  color: #0f172a;
  padding: 1rem;
  font-family: inherit;
  font-size: 1rem;
  transition: all 0.2s ease;
}

input:focus, select:focus, textarea:focus { 
  outline: none; 
  border-color: #1d6bf3; 
  box-shadow: 0 0 0 4px rgba(29, 107, 243, 0.1);
}

.btn-primary {
  width: 100%;
  background: #1d6bf3;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 1rem;
  font-size: 1.1rem;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 56px;
  margin-top: 1rem;
}

.glow-effect:hover:not(:disabled) { 
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(29, 107, 243, 0.3);
}
.btn-primary:disabled { opacity: 0.7; cursor: not-allowed; }
</style>
