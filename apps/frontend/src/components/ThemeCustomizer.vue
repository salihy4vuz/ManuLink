<template>
  <div class="theme-customizer glass-panel">
    <h3>Theme Settings</h3>
    
    <div class="form-group">
      <label>Portal Name (B2B)</label>
      <input type="text" v-model="formData.name" class="neon-input text-left" />
    </div>

    <div class="form-group color-pickers">
      <div>
        <label>Primary Color</label>
        <input type="color" v-model="formData.primaryColor" />
        <span>{{ formData.primaryColor }}</span>
      </div>
      <div>
        <label>Secondary Color</label>
        <input type="color" v-model="formData.secondaryColor" />
        <span>{{ formData.secondaryColor }}</span>
      </div>
    </div>

    <div class="form-group">
      <label>Logo URL (SVG/PNG)</label>
      <input type="text" v-model="formData.logoUrl" class="neon-input text-left" placeholder="https://..." />
    </div>

    <button class="btn btn-primary" @click="saveTheme" :disabled="isSaving">
      {{ isSaving ? 'Saving...' : 'Save Theme' }}
    </button>
    <p v-if="message" class="status-msg">{{ message }}</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useThemeStore } from '../stores/themeStore';
import { updateTenantConfig } from '../api/tenant';

const themeStore = useThemeStore();

const formData = ref({
  name: '',
  primaryColor: '',
  secondaryColor: '',
  logoUrl: ''
});

const isSaving = ref(false);
const message = ref('');

onMounted(() => {
  formData.value = {
    name: themeStore.name,
    primaryColor: themeStore.primaryColor,
    secondaryColor: themeStore.secondaryColor,
    logoUrl: themeStore.logoUrl
  };
});

const saveTheme = async () => {
  isSaving.value = true;
  message.value = '';
  try {
    // Optionally call API to persist
    // await updateTenantConfig(formData.value);
    
    // Update local store immediately to reflect changes
    themeStore.$patch(formData.value);
    themeStore.applyTheme();
    message.value = 'Theme updated successfully!';
  } catch (error) {
    message.value = 'Error updating theme';
  } finally {
    isSaving.value = false;
    setTimeout(() => { message.value = '' }, 3000);
  }
};
</script>

<style scoped>
.theme-customizer {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.color-pickers {
  flex-direction: row;
  gap: 2rem;
}
.color-pickers div {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.text-left {
  text-align: left;
  width: 100%;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: bold;
}
.btn-primary {
  background: var(--color-primary);
  color: #fff;
}
.btn:disabled {
  opacity: 0.7;
}

.status-msg {
  color: #10b981;
  font-size: 0.9rem;
}
</style>
