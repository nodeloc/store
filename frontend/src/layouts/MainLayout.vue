<template>
  <div class="min-h-screen flex flex-col bg-zinc-50/50">
    <!-- Header -->
    <header class="glass-header sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-14 sm:h-16">
          <!-- Logo + Site Name -->
          <router-link to="/" class="flex items-center gap-2.5 min-w-0 group">
            <img
              src="https://s.rmimg.com/original/2X/f/f917abb33b14bc40ecaf1defce5d1903d186b393.svg"
              alt="Logo"
              class="h-7 sm:h-8 w-auto flex-shrink-0 transition-transform duration-300 group-hover:scale-105"
            />
            <div class="min-w-0">
              <span class="text-base sm:text-lg font-bold text-zinc-900 truncate block leading-tight">
                {{ settings.site_name || '社区发卡' }}
              </span>
              <span v-if="settings.site_description" class="hidden sm:block text-xs text-zinc-400 truncate leading-tight">
                {{ settings.site_description }}
              </span>
            </div>
          </router-link>
          
          <!-- Navigation -->
          <div class="flex items-center gap-2 sm:gap-5 flex-shrink-0">
            <router-link to="/" class="hidden sm:inline-flex items-center text-sm text-zinc-500 hover:text-zinc-900 transition-colors">
              首页
            </router-link>
            
            <template v-if="authStore.isAuthenticated">
              <router-link
                to="/orders"
                class="hidden sm:inline-flex items-center gap-1.5 text-sm text-zinc-500 hover:text-zinc-900 transition-colors"
              >
                <ShoppingBag class="w-4 h-4" />
                <span>我的订单</span>
              </router-link>
              
              <router-link
                v-if="authStore.isAdmin"
                to="/admin"
                class="hidden sm:inline-flex items-center gap-1.5 text-sm text-zinc-500 hover:text-zinc-900 transition-colors"
              >
                <LayoutDashboard class="w-4 h-4" />
                <span>管理后台</span>
              </router-link>
              
              <!-- User Menu -->
              <div class="relative" ref="userMenuRef">
                <button
                  @click="showUserMenu = !showUserMenu"
                  class="flex items-center gap-2 px-2 py-1.5 rounded-xl hover:bg-zinc-100 transition-colors"
                >
                  <div class="w-7 h-7 rounded-full bg-brand-gradient flex items-center justify-center text-white text-xs font-bold flex-shrink-0">
                    {{ authStore.user.name?.charAt(0)?.toUpperCase() || 'U' }}
                  </div>
                  <span class="hidden sm:inline text-sm font-medium text-zinc-700 max-w-[100px] truncate">
                    {{ authStore.user.name }}
                  </span>
                  <ChevronDown class="w-3.5 h-3.5 text-zinc-400" />
                </button>

                <transition name="dropdown">
                  <div
                    v-if="showUserMenu"
                    class="absolute right-0 mt-2 w-52 bg-white rounded-xl shadow-lg border border-zinc-100 py-1.5 z-50"
                  >
                    <div class="px-4 py-2.5 border-b border-zinc-100">
                      <div class="text-sm font-semibold text-zinc-900">{{ authStore.user.name }}</div>
                      <div v-if="authStore.user.email" class="text-xs text-zinc-400 mt-0.5 truncate">{{ authStore.user.email }}</div>
                    </div>
                    <div class="py-1">
                      <router-link to="/" class="flex items-center gap-2.5 px-4 py-2.5 text-sm text-zinc-600 hover:bg-zinc-50 hover:text-zinc-900 transition-colors sm:hidden" @click="showUserMenu = false">
                        <Home class="w-4 h-4" />
                        首页
                      </router-link>
                      <router-link to="/orders" class="flex items-center gap-2.5 px-4 py-2.5 text-sm text-zinc-600 hover:bg-zinc-50 hover:text-zinc-900 transition-colors sm:hidden" @click="showUserMenu = false">
                        <ShoppingBag class="w-4 h-4" />
                        我的订单
                      </router-link>
                      <router-link v-if="authStore.isAdmin" to="/admin" class="flex items-center gap-2.5 px-4 py-2.5 text-sm text-zinc-600 hover:bg-zinc-50 hover:text-zinc-900 transition-colors sm:hidden" @click="showUserMenu = false">
                        <LayoutDashboard class="w-4 h-4" />
                        管理后台
                      </router-link>
                      <router-link to="/profile" class="flex items-center gap-2.5 px-4 py-2.5 text-sm text-zinc-600 hover:bg-zinc-50 hover:text-zinc-900 transition-colors" @click="showUserMenu = false">
                        <User class="w-4 h-4" />
                        个人中心
                      </router-link>
                    </div>
                    <div class="border-t border-zinc-100 pt-1">
                      <button @click="authStore.logout" class="flex items-center gap-2.5 w-full px-4 py-2.5 text-sm text-zinc-600 hover:bg-red-50 hover:text-red-600 transition-colors">
                        <LogOut class="w-4 h-4" />
                        退出登录
                      </button>
                    </div>
                  </div>
                </transition>
              </div>
            </template>
            
            <template v-else>
              <button
                @click="handleLogin"
                class="inline-flex items-center gap-1.5 px-4 py-2 bg-brand-gradient text-white text-sm font-medium rounded-xl hover:shadow-glow transition-all duration-300 hover:scale-[1.02]"
              >
                <LogIn class="w-4 h-4" />
                <span>登录</span>
              </button>
            </template>
          </div>
        </div>
      </div>
    </header>
    
    <!-- Main Content (no max-w, each child handles its own) -->
    <main class="flex-1">
      <router-view />
    </main>
    
    <!-- Footer -->
    <footer class="bg-white border-t border-zinc-100 mt-auto">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
        <div class="flex flex-col items-center gap-4">
          <div class="flex items-center gap-2.5">
            <img
              src="https://s.rmimg.com/original/2X/f/f917abb33b14bc40ecaf1defce5d1903d186b393.svg"
              alt="Logo"
              class="h-6 w-auto opacity-60"
            />
            <span class="text-sm font-medium text-zinc-400">{{ settings.site_name || '社区发卡' }}</span>
          </div>
          <p class="text-xs text-zinc-400 text-center">
            {{ settings.footer_text || '© 2026 All rights reserved.' }}
          </p>
          <a href="https://www.nodeloc.com" target="_blank" rel="noopener" class="text-xs text-zinc-400 hover:text-brand-green transition-colors">
            NodeLoc 社区
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { LogIn, LogOut, ChevronDown, LayoutDashboard, User, ShoppingBag, Home } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'

const authStore = useAuthStore()
const settings = ref({})
const showUserMenu = ref(false)
const userMenuRef = ref(null)

function handleClickOutside(event) {
  if (userMenuRef.value && !userMenuRef.value.contains(event.target)) {
    showUserMenu.value = false
  }
}

onMounted(async () => {
  document.addEventListener('click', handleClickOutside)

  try {
    const response = await api.get('/api/settings')
    settings.value = response.data
    if (settings.value.site_name) {
      document.title = settings.value.site_name
    }
  } catch (error) {
    console.error('Failed to fetch settings', error)
  }
  
  try {
    await authStore.fetchUser()
  } catch (error) {
    // User not logged in
  }
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

function handleLogin() {
  window.location.href = '/auth/login'
}
</script>
