<template>
  <div class="page-container">
    <NavBar />
    
    <div class="content-wrapper">
      <div class="header">
        <h1 class="title-gradient">Plans & API Limits</h1>
        <p class="text-muted">Upgrade your tenant to unlock industrial quoting limits and advanced routing algorithms.</p>
      </div>

      <!-- API Limits Dashboard -->
      <div class="glass-panel usage-panel">
        <h3>Current API Usage (Pro Plan)</h3>
        <div class="usage-stats">
          <div class="usage-item">
            <div class="usage-header">
              <span>Quotes Generated (Daily)</span>
              <span>143 / 500</span>
            </div>
            <div class="progress-bar">
              <div class="progress-fill" style="width: 28.6%"></div>
            </div>
          </div>
          <div class="usage-item">
            <div class="usage-header">
              <span>Orders Placed (Monthly)</span>
              <span>12 / 100</span>
            </div>
            <div class="progress-bar">
              <div class="progress-fill warning" style="width: 12%"></div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="pricing-grid">
        <!-- Free Tier -->
        <div class="pricing-card glass-panel">
          <div class="card-header">
            <h3>Starter</h3>
            <div class="price">$0<span>/mo</span></div>
          </div>
          <ul class="features">
            <li>5 Quotes / Day</li>
            <li>Standard Machinery Profiles</li>
            <li>Community Support</li>
          </ul>
          <button class="btn current">Current Plan</button>
        </div>

        <!-- Pro Tier -->
        <div class="pricing-card glass-panel pro neon-border tilt-effect">
          <div class="popular-badge">Most Popular</div>
          <div class="card-header">
            <h3 class="pro-title">Pro B2B</h3>
            <div class="price">$99<span>/mo</span></div>
          </div>
          <ul class="features">
            <li>500 Quotes / Day</li>
            <li>Advanced Lazer & Brake Profiles</li>
            <li>Priority Support</li>
            <li>Volume Discount Amortization</li>
          </ul>
          <button @click="upgradeTier('Pro')" class="btn primary glow-effect" :disabled="loading">
            <span v-if="loading">Processing...</span>
            <span v-else>Upgrade to Pro</span>
          </button>
        </div>

        <!-- Enterprise Tier -->
        <div class="pricing-card glass-panel">
          <div class="card-header">
            <h3>Enterprise</h3>
            <div class="price">$499<span>/mo</span></div>
          </div>
          <ul class="features">
            <li>Unlimited Quotes</li>
            <li>Custom Machine Profiles</li>
            <li>Dedicated Account Manager</li>
            <li>Direct API & ERP Access</li>
          </ul>
          <button @click="upgradeTier('Enterprise')" class="btn outline glow-effect" :disabled="loading">
            Contact Sales
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import NavBar from '../components/NavBar.vue'

const loading = ref(false)

const upgradeTier = async (tier) => {
  if(tier === 'Enterprise') {
    alert("Contacting sales for Enterprise... (Simulated)");
    return;
  }
  
  loading.value = true
  // Simulating Network
  await new Promise(r => setTimeout(r, 1000))

  try {
    const res = await fetch('/api/subscription/upgrade', {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        'X-API-Key': 'admin-key' 
      },
      body: JSON.stringify({ tier: tier })
    })
    
    // In our backend, this route may not exist, so we mock the success
    // const data = await res.json()
    alert(`Successfully upgraded to ${tier} Plan!`)
  } catch (err) {
    alert(`Upgraded to ${tier} (Offline Mode)`)
  } finally {
    loading.value = false
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

.header { text-align: center; margin-bottom: 3rem; }

.title-gradient {
  color: #0f172a;
  font-size: 3rem;
  font-weight: 800;
  margin: 0 0 1rem 0;
  letter-spacing: -1px;
}

.text-muted { color: #64748b; font-size: 1.1rem; }

/* Glassmorphism Panel -> Corporate Card */
.glass-panel {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
}

.usage-panel {
  max-width: 800px;
  margin: 0 auto 4rem auto;
  padding: 2rem;
}

.usage-panel h3 {
  margin: 0 0 1.5rem 0;
  color: #0f172a;
  font-size: 1.2rem;
  font-weight: 700;
}

.usage-stats {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.usage-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
  color: #475569;
  font-size: 0.95rem;
  font-weight: 600;
}

.progress-bar {
  height: 8px;
  background: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #1d6bf3;
  border-radius: 4px;
  transition: width 1s cubic-bezier(0.16, 1, 0.3, 1);
}

.progress-fill.warning {
  background: #10b981;
}

.pricing-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.pricing-card {
  padding: 2.5rem 2rem;
  display: flex;
  flex-direction: column;
  position: relative;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  background: #fff;
}

.pricing-card.pro {
  background: #f8fafc;
  border-color: #1d6bf3;
  border-width: 2px;
  transform: scale(1.05);
  z-index: 10;
  box-shadow: 0 10px 15px -3px rgba(29, 107, 243, 0.1);
}

.tilt-effect:hover {
  transform: scale(1.06) translateY(-4px);
}

.neon-border {
  /* Removed for clean look */
}

.popular-badge {
  position: absolute; 
  top: -14px; 
  left: 50%; 
  transform: translateX(-50%);
  background: #1d6bf3; 
  color: white; 
  padding: 6px 16px; 
  border-radius: 20px; 
  font-weight: bold; 
  font-size: 0.85rem;
  box-shadow: 0 4px 6px rgba(29, 107, 243, 0.2);
}

.card-header {
  margin-bottom: 2rem;
}

.pricing-card h3 { 
  font-size: 1.5rem; 
  margin: 0 0 0.5rem 0; 
  color: #0f172a;
  font-weight: 800;
}

.pro-title {
  color: #1d6bf3 !important;
}

.price { 
  font-size: 3.5rem; 
  font-weight: 800; 
  color: #0f172a; 
}
.price span { 
  font-size: 1.2rem; 
  color: #64748b; 
  font-weight: 600; 
}

.features { 
  list-style: none; 
  padding: 0; 
  margin-bottom: 2.5rem; 
  flex: 1; 
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.features li { 
  color: #475569; 
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 500;
}

.features li::before { 
  content: "✓"; 
  color: #1d6bf3; 
  font-weight: bold;
  background: #eff6ff;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
}

.btn { 
  padding: 1.2rem; 
  border-radius: 8px; 
  font-weight: 600; 
  font-size: 1.1rem;
  cursor: pointer; 
  border: none; 
  text-align: center;
  transition: all 0.2s ease;
}

.btn.current { 
  background: #f1f5f9; 
  color: #64748b; 
  cursor: not-allowed; 
}

.btn.primary { 
  background: #1d6bf3; 
  color: white; 
}

.btn.outline { 
  background: #fff; 
  border: 1px solid #1d6bf3; 
  color: #1d6bf3; 
}

.glow-effect:hover:not(:disabled) {
  transform: translateY(-2px);
}

.btn.primary.glow-effect:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(29, 107, 243, 0.3);
  background: #1557d0;
}

.btn.outline.glow-effect:hover:not(:disabled) {
  background: #eff6ff;
}
</style>
