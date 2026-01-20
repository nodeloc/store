import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useToastStore = defineStore('toast', () => {
  const show = ref(false)
  const message = ref('')
  const type = ref('success')
  let timer = null

  function showToast(msg, toastType = 'success', duration = 2000) {
    if (timer) clearTimeout(timer)
    
    message.value = msg
    type.value = toastType
    show.value = true
    
    timer = setTimeout(() => {
      show.value = false
    }, duration)
  }

  function success(msg) {
    showToast(msg, 'success')
  }

  function error(msg) {
    showToast(msg, 'error')
  }

  return {
    show,
    message,
    type,
    showToast,
    success,
    error
  }
})

// Alias for convenience
export const useToast = useToastStore
