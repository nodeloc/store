<template>
  <div class="min-h-screen flex">
    <!-- Sidebar -->
    <aside class="w-64 bg-zinc-900 text-white flex flex-col">
      <div class="p-6 border-b border-zinc-800">
        <router-link to="/" class="flex items-center space-x-2 text-white hover:text-zinc-300 transition">
          <Ticket class="w-6 h-6" />
          <span class="font-semibold text-lg">管理后台</span>
        </router-link>
      </div>
      
      <nav class="flex-1 p-4 space-y-1">
        <router-link 
          v-for="item in menuItems" 
          :key="item.path"
          :to="item.path"
          class="flex items-center space-x-3 px-4 py-3 rounded-lg text-sm transition"
          :class="isActive(item.path) ? 'bg-zinc-800 text-white' : 'text-zinc-400 hover:bg-zinc-800 hover:text-white'"
        >
          <component :is="item.icon" class="w-5 h-5" />
          <span>{{ item.label }}</span>
        </router-link>
      </nav>
      
      <div class="p-4 border-t border-zinc-800">
        <button @click="goToFrontend" class="flex items-center space-x-3 px-4 py-3 rounded-lg text-sm text-zinc-400 hover:bg-zinc-800 hover:text-white transition w-full">
          <Home class="w-5 h-5" />
          <span>返回前台</span>
        </button>
      </div>
    </aside>
    
    <!-- Main Content -->
    <div class="flex-1 flex flex-col">
      <!-- Header -->
      <header class="bg-white border-b border-zinc-100">
        <div class="px-8 py-4 flex justify-between items-center">
          <h1 class="text-2xl font-semibold text-zinc-900">{{ pageTitle }}</h1>
          
          <div class="flex items-center space-x-4">
            <span class="text-sm text-zinc-600">{{ authStore.user?.name }}</span>
            <img :src="authStore.user?.avatar_url" alt="" class="w-8 h-8 rounded-full">
          </div>
        </div>
      </header>
      
      <!-- Content -->
      <main class="flex-1 p-8 bg-zinc-50">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { 
  Ticket, 
  Home, 
  LayoutDashboard, 
  Package, 
  Folder, 
  CreditCard, 
  ShoppingCart, 
  Users, 
  Settings 
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const menuItems = [
  { path: '/admin', label: '仪表板', icon: LayoutDashboard },
  { path: '/admin/categories', label: '商品分类', icon: Folder },
  { path: '/admin/products', label: '商品管理', icon: Package },
  { path: '/admin/cards', label: '卡密管理', icon: CreditCard },
  { path: '/admin/orders', label: '订单管理', icon: ShoppingCart },
  { path: '/admin/users', label: '用户管理', icon: Users },
  { path: '/admin/settings', label: '系统设置', icon: Settings },
]

const pageTitle = computed(() => {
  const item = menuItems.find(i => i.path === route.path)
  return item ? item.label : '管理后台'
})

function isActive(path) {
  return route.path === path || (path !== '/admin' && route.path.startsWith(path))
}

function goToFrontend() {
  router.push('/')
}
</script>
