<template>
  <div class="ai-assistant" :class="{ 'ai-open': isOpen }">
    <!-- Floating Button -->
    <div v-if="!isOpen" class="ai-fab" @click="isOpen = true">
      <div class="pulse-ring"></div>
      <span class="ai-icon">🤖</span>
    </div>

    <!-- Chat Window -->
    <div v-if="isOpen" class="ai-chat-window glass-panel">
      <div class="ai-header">
        <div class="header-title">
          <span class="ai-icon-small">🤖</span>
          <h3>AI Üretim Mühendisi</h3>
        </div>
        <button class="close-btn" @click="isOpen = false">×</button>
      </div>

      <div class="ai-body">
        <div v-if="loading" class="ai-loading">
          <div class="spinner"></div>
          <p>Tasarımınız DFM standartlarında inceleniyor...</p>
        </div>
        
        <div v-else-if="aiResult" class="ai-message-bubble">
          <div class="ai-score">
            <span class="score-label">Üretilebilirlik Skoru:</span>
            <span class="score-value" :class="scoreClass">{{ aiResult.score }}/100</span>
          </div>
          <p class="ai-advice" v-html="formattedAdvice"></p>
        </div>
        
        <div v-else class="ai-empty">
          <p>CAD verilerinizi veya geometrik parametreleri göndererek analiz başlatabilirsiniz.</p>
        </div>
      </div>

      <div class="ai-footer">
        <button class="btn btn-primary" @click="analyzeDesign" :disabled="loading">
          {{ loading ? 'Analiz Ediliyor...' : 'Geometriyi Analiz Et' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const isOpen = ref(false)
const loading = ref(false)
const aiResult = ref(null)

// Computed class for score color
const scoreClass = computed(() => {
  if (!aiResult.value) return ''
  if (aiResult.value.score >= 90) return 'text-success'
  if (aiResult.value.score >= 50) return 'text-warning'
  return 'text-danger'
})

// Typewriter formatting simulation
const formattedAdvice = computed(() => {
  if (!aiResult.value || !aiResult.value.advice) return ''
  return aiResult.value.advice.replace(/\n/g, '<br>')
})

const analyzeDesign = async () => {
  loading.value = true
  
  // Dummy payload derived from typical CAD UI
  const payload = {
    material: "Paslanmaz Çelik 304",
    thickness_mm: 3.0,
    bend_count: 2,
    flange_length_mm: 8.5
  }

  try {
    const res = await fetch('/api/ai/advisor/analyze', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    const data = await res.json()
    
    if (data.success) {
      aiResult.value = data.data
    } else {
      console.error(data.error)
    }
  } catch (err) {
    console.error("AI Request Failed", err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.ai-assistant {
  position: fixed;
  bottom: 30px;
  right: 30px;
  z-index: 9999;
}

.ai-fab {
  width: 65px;
  height: 65px;
  border-radius: 50%;
  background: linear-gradient(135deg, #6366f1, #a855f7);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 10px 25px rgba(99, 102, 241, 0.4);
  transition: transform 0.3s ease;
  position: relative;
}

.ai-fab:hover {
  transform: scale(1.1);
}

.ai-icon {
  font-size: 2rem;
  z-index: 2;
}

.pulse-ring {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 50%;
  border: 2px solid #a855f7;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { transform: scale(1); opacity: 0.8; }
  100% { transform: scale(1.5); opacity: 0; }
}

.glass-panel {
  background: rgba(15, 23, 42, 0.85);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5);
}

.ai-chat-window {
  width: 380px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transform-origin: bottom right;
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px) scale(0.9); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.ai-header {
  padding: 1rem 1.5rem;
  background: linear-gradient(to right, rgba(99, 102, 241, 0.2), rgba(168, 85, 247, 0.2));
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-title h3 {
  margin: 0;
  font-size: 1.1rem;
  color: #f8fafc;
}

.close-btn {
  background: none;
  border: none;
  color: #94a3b8;
  font-size: 1.5rem;
  cursor: pointer;
  transition: color 0.2s;
}

.close-btn:hover {
  color: #f8fafc;
}

.ai-body {
  padding: 1.5rem;
  min-height: 250px;
  max-height: 400px;
  overflow-y: auto;
}

.ai-empty {
  text-align: center;
  color: #94a3b8;
  font-size: 0.95rem;
  margin-top: 2rem;
}

.ai-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  color: #cbd5e1;
  margin-top: 2rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(99, 102, 241, 0.3);
  border-top-color: #a855f7;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.ai-message-bubble {
  background: rgba(30, 41, 59, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 1.25rem;
  margin-bottom: 1rem;
  border-left: 4px solid #a855f7;
}

.ai-score {
  display: flex;
  justify-content: space-between;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.score-value {
  font-weight: 800;
  font-size: 1.1rem;
}

.text-success { color: #10b981; }
.text-warning { color: #f59e0b; }
.text-danger { color: #ef4444; }

.ai-advice {
  color: #f1f5f9;
  line-height: 1.6;
  font-size: 0.95rem;
  margin: 0;
}

.ai-footer {
  padding: 1.25rem 1.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  background: rgba(15, 23, 42, 0.95);
}

.btn-primary {
  width: 100%;
  padding: 0.8rem;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, #6366f1, #a855f7);
  color: white;
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: opacity 0.2s;
}

.btn-primary:hover:not(:disabled) {
  opacity: 0.9;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
