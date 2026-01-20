import { ref } from 'vue'

export function useCopy() {
  const copied = ref(false)
  let timeout = null

  const copy = async (text) => {
    try {
      await navigator.clipboard.writeText(text)
      copied.value = true
      
      if (timeout) clearTimeout(timeout)
      timeout = setTimeout(() => {
        copied.value = false
      }, 2000)
      
      return true
    } catch (error) {
      console.error('Failed to copy:', error)
      return false
    }
  }

  return {
    copied,
    copy
  }
}
