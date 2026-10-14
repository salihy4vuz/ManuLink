<template>
  <div class="white-label-logo">
    <img v-if="themeStore.logoUrl" :src="themeStore.logoUrl" :alt="themeStore.name" class="logo-img" />
    <span v-else class="logo-text">{{ themeStore.name }}</span>
  </div>
</template>

<script setup>
import { useThemeStore } from '../stores/themeStore';

const themeStore = useThemeStore();
</script>

<style scoped>
.white-label-logo {
  display: flex;
  align-items: center;
  gap: 10px;
}
.logo-img {
  max-height: 40px;
  max-width: 150px;
  object-fit: contain;
}
.logo-text {
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--color-primary);
}
</style>
