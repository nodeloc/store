<template>
  <div class="space-y-6">
    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div v-for="stat in stats" :key="stat.label" class="bg-white rounded-lg border border-zinc-100 p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-zinc-600">{{ stat.label }}</p>
            <p class="text-3xl font-semibold text-zinc-900 mt-2">{{ stat.value }}</p>
          </div>
          <div class="p-3 bg-zinc-100 rounded-lg">
            <component :is="stat.icon" class="w-6 h-6 text-zinc-600" />
          </div>
        </div>
      </div>
    </div>
    
    <!-- Quick Actions -->
    <div class="bg-white rounded-lg border border-zinc-100 p-6">
      <h2 class="text-lg font-semibold text-zinc-900 mb-4">快速操作</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <router-link 
          to="/admin/categories" 
          class="flex items-center space-x-3 p-4 border border-zinc-200 rounded-lg hover:border-zinc-900 transition"
        >
          <Folder class="w-5 h-5 text-zinc-600" />
          <span class="text-sm font-medium text-zinc-900">添加分类</span>
        </router-link>
        
        <router-link 
          to="/admin/products" 
          class="flex items-center space-x-3 p-4 border border-zinc-200 rounded-lg hover:border-zinc-900 transition"
        >
          <Package class="w-5 h-5 text-zinc-600" />
          <span class="text-sm font-medium text-zinc-900">添加商品</span>
        </router-link>
        
        <router-link 
          to="/admin/cards" 
          class="flex items-center space-x-3 p-4 border border-zinc-200 rounded-lg hover:border-zinc-900 transition"
        >
          <CreditCard class="w-5 h-5 text-zinc-600" />
          <span class="text-sm font-medium text-zinc-900">添加卡密</span>
        </router-link>
      </div>
    </div>
    
    <!-- Loading/Error States -->
    <div v-if="loading" class="text-center py-12">
      <Loader2 class="w-8 h-8 animate-spin text-zinc-400 mx-auto" />
    </div>
    
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
      <p class="text-sm text-red-800">{{ error }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Package, Users, ShoppingCart, Folder, CreditCard, Loader2 } from 'lucide-vue-next'
import api from '@/utils/api'

const stats = ref([
  { label: '商品分类', value: 0, icon: Folder },
  { label: '商品总数', value: 0, icon: Package },
  { label: '订单总数', value: 0, icon: ShoppingCart },
  { label: '用户总数', value: 0, icon: Users },
])

const loading = ref(false)
const error = ref(null)

onMounted(async () => {
  try {
    loading.value = true
    const response = await api.get('/api/admin/dashboard')
    const data = response.data.stats
    
    stats.value = [
      { label: '商品分类', value: data.categories || 0, icon: Folder },
      { label: '商品总数', value: data.products || 0, icon: Package },
      { label: '订单总数', value: data.orders || 0, icon: ShoppingCart },
      { label: '用户总数', value: data.users || 0, icon: Users },
    ]
  } catch (err) {
    error.value = '加载数据失败'
    console.error(err)
  } finally {
    loading.value = false
  }
})
</script>
