<script setup>
import { useToastStore } from '../stores/toast'
const toastStore = useToastStore()
</script>

<template>
  <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 pointer-events-none">
    <TransitionGroup name="toast">
      <div 
        v-for="toast in toastStore.toasts" 
        :key="toast.id"
        class="min-w-[250px] p-4 rounded-lg shadow-lg text-white font-medium flex items-center justify-between gap-4 pointer-events-auto backdrop-blur-md"
        :class="{
          'bg-red-500/90 border border-red-400': toast.type === 'error',
          'bg-emerald-500/90 border border-emerald-400': toast.type === 'success',
          'bg-blue-500/90 border border-blue-400': toast.type === 'info'
        }"
      >
        <span>{{ toast.message }}</span>
        <button @click="toastStore.remove(toast.id)" class="hover:opacity-75 transition-opacity font-bold">
          &times;
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(30px);
}
.toast-leave-to {
  opacity: 0;
  transform: scale(0.9);
}
</style>
