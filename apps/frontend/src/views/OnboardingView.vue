<template>
  <div class="onboarding-container">
    <h2>Welcome to Onboarding</h2>
    <div class="stepper">
      <div v-for="(step, index) in steps" :key="index" :class="['step', { active: currentStep === index }]">
        {{ index + 1 }}. {{ step.name }}
      </div>
    </div>
    
    <div class="step-content">
      <div v-if="currentStep === 0">
        <h3>Şirket Bilgileri</h3>
        <input v-model="formData.companyName" placeholder="Şirket Adı" />
      </div>
      <div v-if="currentStep === 1">
        <h3>CNC Makineleri</h3>
        <label><input type="checkbox" v-model="formData.hasLaser" /> Laser Cutter</label>
      </div>
      <!-- Diğer adımlar... -->
    </div>
    
    <div class="actions">
      <button :disabled="currentStep === 0" @click="prev">Geri</button>
      <button :disabled="currentStep === steps.length - 1" @click="next">İleri</button>
      <button v-if="currentStep === steps.length - 1" @click="finish">Tamamla</button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'

const steps = [
  { name: 'Şirket Bilgileri' },
  { name: 'CNC Makineleri' },
  { name: 'Malzemeler' },
  { name: 'ERP Konfigürasyonu' },
  { name: 'CAD Bridge İndir' }
]

const currentStep = ref(0)
const formData = ref({
  companyName: '',
  hasLaser: false
})

onMounted(() => {
  const saved = localStorage.getItem('onboardingState')
  if (saved) {
    const state = JSON.parse(saved)
    currentStep.value = state.currentStep
    formData.value = state.formData
  }
})

watch([currentStep, formData], () => {
  localStorage.setItem('onboardingState', JSON.stringify({
    currentStep: currentStep.value,
    formData: formData.value
  }))
}, { deep: true })

const next = () => { if (currentStep.value < steps.length - 1) currentStep.value++ }
const prev = () => { if (currentStep.value > 0) currentStep.value-- }
const finish = () => { alert('Onboarding completed!') }
</script>

<style scoped>
.onboarding-container { padding: 2rem; }
.stepper { display: flex; gap: 1rem; margin-bottom: 2rem; }
.step { padding: 1rem; border: 1px solid #ccc; flex: 1; text-align: center; }
.step.active { background: #4caf50; color: white; border-color: #4caf50; }
.actions { margin-top: 2rem; display: flex; gap: 1rem; }
</style>
