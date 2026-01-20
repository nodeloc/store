<template>
  <div class="max-w-4xl mx-auto space-y-8">
    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="w-8 h-8 animate-spin text-zinc-400" />
    </div>
    
    <template v-else-if="user">
      <!-- Profile Header -->
      <div class="bg-gradient-to-br from-zinc-50 to-zinc-100 border border-zinc-200 rounded-lg p-8">
        <div class="flex items-center space-x-6">
          <img :src="user.avatar_url" alt="" class="w-24 h-24 rounded-full border-4 border-white shadow-lg">
          <div class="flex-1 space-y-2">
            <h1 class="text-3xl font-bold text-zinc-900">{{ user.name }}</h1>
            <div class="flex items-center space-x-4 text-sm text-zinc-600">
              <span class="flex items-center space-x-1">
                <AtSign class="w-4 h-4" />
                <span>{{ user.username }}</span>
              </span>
              <span class="flex items-center space-x-1">
                <Shield class="w-4 h-4" />
                <span>信任等级 {{ user.trust_level }}</span>
              </span>
              <span v-if="user.email" class="flex items-center space-x-1">
                <Mail class="w-4 h-4" />
                <span>{{ user.email }}</span>
              </span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Stats -->
      <div class="grid grid-cols-3 gap-6">
        <div class="bg-white border border-zinc-100 rounded-lg p-6 text-center space-y-2">
          <div class="text-4xl font-bold font-mono text-zinc-900">{{ orders.length }}</div>
          <div class="text-sm text-zinc-500 uppercase tracking-wide">订单总数</div>
        </div>
        <div class="bg-white border border-zinc-100 rounded-lg p-6 text-center space-y-2">
          <div class="text-4xl font-bold font-mono text-zinc-900">{{ formatPrice(user.balance) }}</div>
          <div class="text-sm text-zinc-500 uppercase tracking-wide">账户余额</div>
        </div>
        <div class="bg-white border border-zinc-100 rounded-lg p-6 text-center space-y-2">
          <div class="text-4xl font-bold font-mono text-zinc-900">{{ user.trust_level }}</div>
          <div class="text-sm text-zinc-500 uppercase tracking-wide">信任等级</div>
        </div>
      </div>
      
      <!-- Recent Orders -->
      <div class="bg-white border border-zinc-100 rounded-lg overflow-hidden">
        <div class="px-6 py-4 bg-zinc-50 border-b border-zinc-100 flex items-center justify-between">
          <h2 class="text-sm font-medium text-zinc-900 uppercase tracking-wide">最近订单</h2>
          <router-link to="/orders" class="text-sm text-zinc-600 hover:text-zinc-900 flex items-center space-x-1">
            <span>查看全部</span>
            <ArrowRight class="w-4 h-4" />
          </router-link>
        </div>
        <div class="divide-y divide-zinc-100">
          <template v-if="orders.length > 0">
            <div
              v-for="(order, index) in orders.slice(0, 5)"
              :key="order.id"
              class="px-6 py-4 flex items-center justify-between hover:bg-zinc-50 transition"
            >
              <div class="space-y-1">
                <div class="font-medium text-zinc-900">{{ order.product.name }}</div>
                <div class="text-sm text-zinc-500 font-mono">{{ formatDate(order.created_at, 'YYYY-MM-DD HH:mm') }}</div>
              </div>
              <div class="flex items-center space-x-4">
                <div class="text-right">
                  <div class="font-mono font-semibold text-zinc-900">{{ formatPrice(order.total_amount) }}</div>
                </div>
                <router-link
                  :to="`/order/${order.order_no}`"
                  class="text-sm text-zinc-600 hover:text-zinc-900 flex items-center space-x-1"
                >
                  <span>查看</span>
                  <ChevronRight class="w-4 h-4" />
                </router-link>
              </div>
            </div>
          </template>
          <div v-else class="px-6 py-12 text-center text-zinc-500">
            暂无订单
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { AtSign, Shield, Mail, ArrowRight, ChevronRight, Loader2 } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'
import { formatPrice, formatDate } from '@/utils/helpers'

const authStore = useAuthStore()
const loading = ref(true)
const user = ref(null)
const orders = ref([])

onMounted(async () => {
  try {
    const [userRes, ordersRes] = await Promise.all([
      api.get('/api/user/info'),
      api.get('/api/orders')
    ])
    user.value = userRes.data.user
    orders.value = ordersRes.data
  } catch (error) {
    console.error('Failed to load profile', error)
  } finally {
    loading.value = false
  }
})
</script>
