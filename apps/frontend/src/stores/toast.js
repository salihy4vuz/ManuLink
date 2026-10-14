import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useToastStore = defineStore('toast', () => {
  const toasts = ref([])

  function add(message, type = 'info', duration = 3000) {
    const id = Date.now() + Math.random()
    toasts.value.push({ id, message, type })
    setTimeout(() => {
      remove(id)
    }, duration)
  }

  function remove(id) {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }

  function error(message, duration = 4000) { add(message, 'error', duration) }
  function success(message, duration = 3000) { add(message, 'success', duration) }
  function info(message, duration = 3000) { add(message, 'info', duration) }

  return { toasts, add, remove, error, success, info }
})
