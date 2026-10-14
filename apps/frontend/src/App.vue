<template>
  <div v-if="themeStore.isLoading" class="loading-screen">
    <div class="loader-core"></div>
    <div class="loader-text">INITIALIZING CAD ENGINE...</div>
  </div>
  <div v-else>
    <router-view />
    <ToastNotification />
    <AIAssistant />
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useThemeStore } from './stores/themeStore';
import ToastNotification from './components/ToastNotification.vue';
import AIAssistant from './components/AIAssistant.vue';

const themeStore = useThemeStore();

onMounted(async () => {
  // Use hostname to fetch tenant theme configuration
  const hostname = window.location.hostname;
  await themeStore.loadTenantTheme(hostname);
});
</script>

<style>
/* Global CSS variables fallback */
@import './assets/theme.css';

.loading-screen {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: var(--bg-dark);
  position: relative;
  overflow: hidden;
}

.loader-core {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: rgba(0, 229, 255, 0.1);
  border: 2px solid rgba(0, 229, 255, 0.5);
  box-shadow: 0 0 20px rgba(0, 229, 255, 0.4), inset 0 0 20px rgba(0, 229, 255, 0.4);
  animation: pulse 1.5s ease-in-out infinite alternate, spin 3s linear infinite;
  margin-bottom: 24px;
}

.loader-text {
  font-family: var(--font-mono);
  color: var(--accent-neon);
  font-size: 1.1rem;
  letter-spacing: 3px;
  text-transform: uppercase;
  animation: textGlow 1.5s ease-in-out infinite alternate;
}

@keyframes pulse {
  0% { transform: scale(0.9); box-shadow: 0 0 10px rgba(0, 229, 255, 0.2), inset 0 0 10px rgba(0, 229, 255, 0.2); }
  100% { transform: scale(1.1); box-shadow: 0 0 30px rgba(0, 229, 255, 0.8), inset 0 0 30px rgba(0, 229, 255, 0.8); }
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@keyframes textGlow {
  0% { opacity: 0.7; text-shadow: 0 0 5px rgba(0, 229, 255, 0.3); }
  100% { opacity: 1; text-shadow: 0 0 15px rgba(0, 229, 255, 0.8); }
}
</style>
