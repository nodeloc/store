<template>
  <div class="min-h-screen flex flex-col">
    <!-- Header -->
    <header class="sticky top-0 z-50 bg-white border-b border-zinc-100">
      <nav class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo -->
          <router-link to="/" class="flex items-center space-x-2 text-zinc-900 font-semibold text-lg hover:text-zinc-700 transition">
            <Ticket class="w-6 h-6" />
            <span>{{ settings.site_name || 'NodeLoc 社区发卡' }}</span>
          </router-link>
          
          <!-- Navigation -->
          <div class="flex items-center space-x-6">
            <router-link to="/" class="text-sm text-zinc-600 hover:text-zinc-900 transition">首页</router-link>
            
            <template v-if="authStore.isAuthenticated">
              <router-link to="/orders" class="text-sm text-zinc-600 hover:text-zinc-900 transition">我的订单</router-link>
              
              <!-- Admin Link -->
              <router-link v-if="authStore.isAdmin" to="/admin" class="inline-flex items-center space-x-1 text-sm text-zinc-900 font-medium hover:text-zinc-700 transition">
                <Settings class="w-4 h-4" />
                <span>管理后台</span>
              </router-link>
              
              <!-- User Menu -->
              <div class="relative group">
                <button class="flex items-center space-x-2 text-sm text-zinc-600 hover:text-zinc-900 transition">
                  <img :src="authStore.user.avatar_url" alt="" class="w-6 h-6 rounded-full">
                  <span>{{ authStore.user.name }}</span>
                  <ChevronDown class="w-4 h-4" />
                </button>
                <div class="absolute right-0 mt-2 w-48 bg-white border border-zinc-100 rounded-lg shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all">
                  <router-link to="/profile" class="block px-4 py-2 text-sm text-zinc-700 hover:bg-zinc-50">个人中心</router-link>
                  <button @click="authStore.logout" class="w-full text-left px-4 py-2 text-sm text-zinc-700 hover:bg-zinc-50">退出登录</button>
                </div>
              </div>
            </template>
            
            <template v-else>
              <button @click="handleLogin" class="inline-flex items-center space-x-1 px-4 py-2 bg-zinc-900 text-white text-sm font-medium rounded-lg hover:bg-zinc-800 transition">
                <LogIn class="w-4 h-4" />
                <span>登录</span>
              </button>
            </template>
          </div>
        </div>
      </nav>
    </header>
    
    <!-- Main Content -->
    <main class="flex-1 max-w-5xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <router-view />
    </main>
    
    <!-- Footer -->
    <footer class="border-t border-zinc-100 mt-16">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <p class="text-center text-sm text-zinc-500">
          {{ settings.footer_text || '© 2026 NodeLoc 社区发卡系统' }}
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Ticket, LogIn, ChevronDown, Settings } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'

const authStore = useAuthStore()
const settings = ref({})

onMounted(async () => {
  // Fetch site settings
  try {
    const response = await api.get('/api/settings')
    settings.value = response.data
  } catch (error) {
    console.error('Failed to fetch settings', error)
  }
  
  // Try to fetch user info
  try {
    await authStore.fetchUser()
  } catch (error) {
    // User not logged in
  }
})

function handleLogin() {
  window.location.href = '/auth/login'
}
</script>
