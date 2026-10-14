<template>
  <div class="page-container">
    <NavBar />
    
    <div class="content-wrapper">
      <div class="header">
        <h1 class="title-gradient">Tenant Settings & ERP Configuration</h1>
        <p class="text-muted">Manage your factory's integrations, white-label branding, and security tokens securely.</p>
      </div>

      <div class="settings-grid">
        <!-- ERP Configuration Panel -->
        <div class="glass-panel settings-card">
          <div class="card-header">
            <h3><span class="icon">🔌</span> ERP Integration</h3>
          </div>
          
          <form @submit.prevent="saveSettings" class="form-grid">
            <div class="input-group">
              <label>ERP System</label>
              <select v-model="settings.erpSystem">
                <option value="logo">Logo Tiger 3</option>
                <option value="sap">SAP S/4HANA</option>
                <option value="odoo">Odoo</option>
                <option value="none">None (Standalone)</option>
              </select>
            </div>

            <div class="input-group" v-if="settings.erpSystem !== 'none'">
              <label>ERP Endpoint URL</label>
              <input type="url" v-model="settings.erpEndpoint" placeholder="https://internal-erp:8443/api" required />
            </div>

            <div class="input-group" v-if="settings.erpSystem !== 'none'">
              <label>API Key / Auth Token (AES-256 Encrypted)</label>
              <input type="password" v-model="settings.erpToken" placeholder="••••••••••••••••" />
              <small class="text-muted mt-1">Stored securely using AES-256 encryption. Never exposed to clients.</small>
            </div>
            
            <button type="submit" class="btn-primary glow-effect mt-4">
              Save ERP Configuration
            </button>
          </form>
        </div>

        <!-- White Label Branding Panel -->
        <div class="glass-panel settings-card">
          <div class="card-header">
            <h3><span class="icon">🎨</span> White-Label Branding</h3>
          </div>
          
          <div class="form-grid">
            <div class="input-group">
              <label>Portal Accent Color</label>
              <div class="color-picker-row">
                <input type="color" v-model="settings.themeColor" />
                <span class="color-hex">{{ settings.themeColor }}</span>
              </div>
            </div>

            <div class="input-group">
              <label>Company Logo URL</label>
              <input type="url" v-model="settings.logoUrl" placeholder="https://yourcompany.com/logo.png" />
            </div>

            <div class="preview-box mt-3" :style="{ borderColor: settings.themeColor }">
              <p>Your B2B customers will see this color theme.</p>
              <button class="mock-btn" :style="{ background: settings.themeColor }">Preview Button</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import NavBar from '../components/NavBar.vue'

// Simulated fetch from /api/v1/tenants/config
const settings = ref({
  erpSystem: 'sap',
  erpEndpoint: 'https://s4hana-host:443/sap/opu/odata/sap/API_SALES_ORDER_SRV',
  erpToken: '',
  themeColor: '#38bdf8',
  logoUrl: ''
})

const saveSettings = () => {
  // In reality, this posts to /api/v1/tenants/config
  alert('Settings saved securely to the database!')
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

.header {
  margin-bottom: 40px;
}

.text-muted {
  color: #475569;
}

.title-gradient {
  color: #0f172a;
  font-size: 2.5rem;
  font-weight: 800;
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.5px;
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(450px, 1fr));
  gap: 2rem;
  max-width: 1200px;
}

/* Glassmorphism Panel -> Corporate Card */
.glass-panel {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

.settings-card {
  padding: 2.5rem;
}

.card-header h3 {
  margin: 0 0 1.5rem 0;
  color: #0f172a;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.3rem;
  font-weight: 700;
}

.form-grid {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

label {
  font-size: 0.85rem;
  color: #475569;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

input, select {
  padding: 0.75rem 1rem;
  background: #fff;
  border: 1px solid #cbd5e1;
  color: #0f172a;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.2s ease;
}

input:focus, select:focus {
  outline: none;
  border-color: #1d6bf3;
  box-shadow: 0 0 0 3px rgba(29, 107, 243, 0.1);
}

.color-picker-row {
  display: flex;
  align-items: center;
  gap: 1rem;
}

input[type="color"] {
  padding: 0;
  width: 50px;
  height: 50px;
  border-radius: 8px;
  cursor: pointer;
}

.color-hex {
  color: #0f172a;
  font-family: monospace;
  font-size: 1.1rem;
}

.preview-box {
  padding: 1.5rem;
  border: 2px dashed #cbd5e1;
  border-radius: 12px;
  background: #f8fafc;
  text-align: center;
  color: #475569;
}

.mock-btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  color: white;
  font-weight: 600;
  margin-top: 1rem;
}

.btn-primary {
  padding: 1rem;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  font-size: 1.1rem;
  cursor: pointer;
  transition: all 0.2s ease;
  color: white;
  background: #1d6bf3;
}

.glow-effect:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(29, 107, 243, 0.3);
}

.mt-4 { margin-top: 1.5rem; }
.mt-3 { margin-top: 1rem; }
.mt-1 { margin-top: 0.25rem; }
</style>
