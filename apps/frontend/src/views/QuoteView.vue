<script setup>
import { ref, computed } from 'vue'
import { TresCanvas, useRenderLoop } from '@tresjs/core'
import { useRouter } from 'vue-router'
import Uploader from '../components/Uploader.vue'
import api from '../api/client.js'
import { ShoppingCart, CheckCircle2 } from 'lucide-vue-next'

const router = useRouter()

const form = ref({ material: 'steel', thickness: 2, bounding_box_x: 200, bounding_box_y: 150, quantity: 10, bend_count: 2 })
const result = ref(null)
const loading = ref(false)
const isUploading = ref(false)
const fileUrl = ref(null)
const fileName = ref('')
const orderSuccess = ref(false)

const sheetRef = ref(null)
const rotationY = ref(0)
const { onLoop } = useRenderLoop()

onLoop(() => {
  rotationY.value += 0.005
})

const metalColor = computed(() => form.value.material === 'aluminum' ? '#d0d4dc' : '#475569')

const onFileDropped = async (file) => {
  isUploading.value = true
  const formData = new FormData()
  formData.append('file', file)

  try {
    const res = await api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    const data = res.data
    if (data.success || data.file_url) {
      fileUrl.value = data.file_url || data.url
      fileName.value = data.file_name || file.name
    }
  } catch (err) {
    alert('Upload failed')
  } finally {
    isUploading.value = false
  }
}

const getQuote = async () => {
  loading.value = true
  try {
    const payload = {
      ...form.value,
      machine_type: 'laser_cutting', // Default to laser/sheet metal
      file_url: fileUrl.value,
      file_name: fileName.value
    }
    const res = await api.post('/quotes', payload)
    const data = res.data
    result.value = data.quote || data
  } catch (err) {
    if (err.response && err.response.status === 429) {
      alert(`QUOTA EXCEEDED: ${err.response.data.message || err.message}`)
    } else {
      alert("API Error")
    }
  } finally {
    loading.value = false
  }
}

const placeOrder = async () => {
  try {
    const res = await api.post('/orders', {
      customer_id: 1, // Will be read from JWT in backend
      part_name: fileName.value || 'manual_entry.dxf',
      quote_data: form.value,
      shipping_address: '123 Fake St, Industry City'
    })
    const data = res.data
    if (data.success || data.order_id) {
      orderSuccess.value = true
      setTimeout(() => {
        router.push('/orders')
      }, 3000)
    }
  } catch(e) {
    alert("Failed to place order")
  }
}
</script>

<template>
  <div class="quote-view">
    <div class="canvas-container glass-panel">
      <div v-if="!fileUrl" class="upload-overlay">
        <Uploader :is-calculating="isUploading" @file-dropped="onFileDropped" />
      </div>
      <TresCanvas clear-color="transparent" v-else>
        <TresPerspectiveCamera :position="[400, 400, 500]" :look-at="[0, 0, 0]" />
        <TresAmbientLight :intensity="0.5" />
        <TresDirectionalLight :position="[200, 500, 300]" :intensity="1" />
        <TresGridHelper :args="[1000, 20]" :position="[0, -50, 0]" />
        
        <TresMesh ref="sheetRef" :rotation-y="rotationY">
          <TresBoxGeometry :args="[form.bounding_box_x, form.thickness, form.bounding_box_y]" />
          <TresMeshStandardMaterial :color="metalColor" :metalness="0.8" :roughness="0.3" />
        </TresMesh>
      </TresCanvas>
    </div>
    
    <div class="form-container glass-panel">
      <h2 class="title-gradient">Instant Quote</h2>
      
      <form @submit.prevent="getQuote" class="form-grid">
        <div class="input-group">
          <label>Material</label>
          <select v-model="form.material">
            <option value="steel">Carbon Steel</option>
            <option value="aluminum">Aluminum 6061</option>
            <option value="stainless">Stainless 304</option>
          </select>
        </div>
        <div class="input-row">
          <div class="input-group">
            <label>Thickness</label>
            <input v-model.number="form.thickness" type="number" step="0.1" />
          </div>
          <div class="input-group">
            <label>Quantity</label>
            <input v-model.number="form.quantity" type="number" />
          </div>
        </div>
        
        <button type="submit" class="btn-neon" :disabled="loading">
          <span v-if="loading">Calculating...</span>
          <span v-else>Analyze & Price</span>
        </button>
      </form>
      
      <transition name="fade-slide">
        <div v-if="result" class="result-card glass-panel mt-4">
          <div class="price-header">
            <h3 class="mono accent-success">{{ result.final_selling_price.toFixed(2) }} {{ result.currency }}</h3>
            <span class="unit-price mono">{{ result.unit_cost_usd.toFixed(2) }} / unit</span>
          </div>
          
          <div class="price-details">
            <p>Setup Cost: <span>{{ result.setup_cost_usd.toFixed(2) }}</span></p>
            <p>Material: <span>{{ result.material_cost_usd?.toFixed(2) || '0.00' }}</span></p>
            <p>Laser Time: <span>{{ result.laser_cost_usd?.toFixed(2) || '0.00' }}</span></p>
          </div>
          
          <button @click="placeOrder" class="btn-neon mt-4" style="width: 100%;">
            <ShoppingCart class="lucide-icon inline-icon" size="20" /> Place Order Now
          </button>
        </div>
      </transition>
    </div>

    <!-- Success Modal -->
    <dialog class="success-modal" :open="orderSuccess">
      <div class="modal-content glass-panel">
        <div class="check-icon-wrap">
          <CheckCircle2 size="64" class="success-icon" />
        </div>
        <h2>Order Confirmed!</h2>
        <p>Your B2B order has been sent to production.</p>
        <p class="text-muted">Redirecting to dashboard...</p>
      </div>
    </dialog>
  </div>
</template>

<style scoped>
.quote-view { 
  display: flex; 
  gap: 2rem; 
  height: calc(100vh - 6rem); 
  padding: 2rem;
  color: var(--text-main);
}

.canvas-container { 
  flex: 1; 
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
}

.upload-overlay {
  position: absolute;
  width: 60%;
  max-width: 500px;
  z-index: 10;
}

.form-container { 
  width: 450px; 
  padding: 2.5rem; 
  display: flex;
  flex-direction: column;
}

.title-gradient {
  color: var(--accent-neon);
  font-size: 1.8rem;
  font-weight: 800;
  margin-bottom: 2rem;
  letter-spacing: -0.5px;
  text-shadow: 0 0 10px var(--accent-glow);
}

.form-grid { 
  display: flex; 
  flex-direction: column; 
  gap: 1.5rem; 
}

.input-row {
  display: flex;
  gap: 1rem;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  flex: 1;
}

label {
  font-size: 0.85rem;
  color: var(--text-muted);
  font-weight: 600;
}

input, select { 
  padding: 0.75rem 1rem; 
  background: rgba(0, 0, 0, 0.2); 
  border: 1px solid var(--border-glass); 
  color: var(--text-main); 
  border-radius: 6px;
  font-size: 0.95rem;
  transition: all 0.2s ease;
}

input:focus, select:focus {
  outline: none;
  border-color: var(--accent-neon);
  box-shadow: 0 0 0 3px var(--accent-glow);
}

.result-card { 
  margin-top: 2rem; 
  padding: 1.5rem; 
  border-radius: 8px; 
}

.price-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 1rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border-glass);
}

.price-header h3 {
  font-size: 2.2rem;
  margin: 0;
  font-weight: 800;
}

.accent-success {
  color: var(--accent-success);
  text-shadow: 0 0 10px rgba(0, 230, 118, 0.3);
}

.unit-price {
  color: var(--text-muted);
  font-size: 0.9rem;
  font-weight: 500;
}

.price-details p {
  display: flex;
  justify-content: space-between;
  color: var(--text-muted);
  margin-bottom: 0.5rem;
  font-size: 0.95rem;
}

.price-details span {
  font-weight: 700;
  color: var(--text-main);
}

.mt-4 { margin-top: 1rem; }

/* View Transitions & Animations */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(20px);
}
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* Dialog Animation @starting-style */
.success-modal {
  margin: auto;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: 24px;
}

.success-modal::backdrop {
  background: rgba(15, 23, 42, 0.5);
  backdrop-filter: blur(4px);
  transition: opacity 0.3s ease;
}

.modal-content {
  padding: 3rem;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.modal-content h2 {
  color: var(--text-main);
  margin: 0;
}

.modal-content p {
  color: var(--text-muted);
}

.check-icon-wrap {
  margin-bottom: 1rem;
  animation: scale-up 0.5s cubic-bezier(0.16, 1, 0.3, 1);
}

.success-icon {
  color: #10b981;
}

@keyframes scale-up {
  0% { transform: scale(0); opacity: 0; }
  100% { transform: scale(1); opacity: 1; }
}

.success-modal[open] {
  animation: modal-enter 0.4s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes modal-enter {
  0% { opacity: 0; transform: scale(0.95) translateY(20px); }
  100% { opacity: 1; transform: scale(1) translateY(0); }
}
</style>
