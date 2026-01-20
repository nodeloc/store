import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.is_admin === true)

  async function fetchUser() {
    try {
      loading.value = true
      const response = await api.get('/api/user/info')
      user.value = response.data.user
      return user.value
    } catch (error) {
      user.value = null
      throw error
    } finally {
      loading.value = false
    }
  }

  function logout() {
    user.value = null
    window.location.href = '/auth/logout'
  }

  return {
    user,
    loading,
    isAuthenticated,
    isAdmin,
    fetchUser,
    logout
  }
})
