<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20 gap-3">
      <Loader2 class="w-8 h-8 animate-spin text-brand-green" />
      <span class="text-sm text-zinc-400">加载中...</span>
    </div>
    
    <template v-else-if="user">
      <!-- Profile Header -->
      <div class="bg-white rounded-2xl border border-zinc-100 shadow-card overflow-hidden">
        <div class="h-20 sm:h-28 bg-brand-gradient"></div>
        <div class="px-5 sm:px-8 pb-6 -mt-10 sm:-mt-12">
          <div class="flex flex-col sm:flex-row sm:items-end gap-4">
            <div class="w-20 h-20 sm:w-24 sm:h-24 rounded-2xl bg-white border-4 border-white shadow-lg overflow-hidden flex-shrink-0">
              <img v-if="user.avatar_url" :src="user.avatar_url" alt="" class="w-full h-full object-cover" />
              <div v-else class="w-full h-full bg-brand-gradient flex items-center justify-center text-white text-2xl font-bold">
                {{ user.name?.charAt(0)?.toUpperCase() || 'U' }}
              </div>
            </div>
            <div class="flex-1 min-w-0 space-y-1">
              <h1 class="text-xl sm:text-2xl font-bold text-zinc-900">{{ user.name }}</h1>
              <div class="flex flex-wrap items-center gap-3 text-sm text-zinc-500">
                <span class="flex items-center gap-1">
                  <AtSign class="w-3.5 h-3.5" />
                  {{ user.username }}
                </span>
                <span class="flex items-center gap-1">
                  <Shield class="w-3.5 h-3.5" />
                  信任等级 {{ user.trust_level }}
                </span>
                <span v-if="user.email" class="flex items-center gap-1">
                  <Mail class="w-3.5 h-3.5" />
                  {{ user.email }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Stats -->
      <div class="grid grid-cols-3 gap-3 sm:gap-4">
        <div class="bg-white rounded-2xl border border-zinc-100 shadow-card p-4 sm:p-6 text-center space-y-1">
          <div class="text-2xl sm:text-3xl font-bold font-mono text-zinc-900">{{ orders.length }}</div>
          <div class="text-xs text-zinc-400 uppercase tracking-wider">订单总数</div>
        </div>
        <div class="bg-white rounded-2xl border border-zinc-100 shadow-card p-4 sm:p-6 text-center space-y-1">
          <div class="text-2xl sm:text-3xl font-bold font-mono text-zinc-900">{{ formatPrice(user.balance) }}</div>
          <div class="text-xs text-zinc-400 uppercase tracking-wider">账户余额</div>
        </div>
        <div class="bg-white rounded-2xl border border-zinc-100 shadow-card p-4 sm:p-6 text-center space-y-1">
          <div class="text-2xl sm:text-3xl font-bold font-mono text-zinc-900">{{ user.trust_level }}</div>
          <div class="text-xs text-zinc-400 uppercase tracking-wider">信任等级</div>
        </div>
      </div>
      
      <!-- Recent Orders -->
      <div class="bg-white rounded-2xl border border-zinc-100 shadow-card overflow-hidden">
        <div class="px-5 sm:px-6 py-4 border-b border-zinc-100 bg-zinc-50/50 flex items-center justify-between">
          <h2 class="text-xs font-semibold text-zinc-400 uppercase tracking-wider">最近订单</h2>
          <router-link to="/orders" class="text-xs text-zinc-500 hover:text-brand-green flex items-center gap-1 transition-colors">
            <span>查看全部</span>
            <ArrowRight class="w-3.5 h-3.5" />
          </router-link>
        </div>
        <div class="divide-y divide-zinc-100">
          <template v-if="orders.length > 0">
            <router-link
              v-for="order in orders.slice(0, 5)"
              :key="order.id"
              :to="`/order/${order.order_no}`"
              class="flex items-center justify-between px-5 sm:px-6 py-4 hover:bg-zinc-50/50 transition-colors"
            >
              <div class="space-y-0.5 min-w-0">
                <div class="text-sm font-medium text-zinc-900 truncate">{{ order.product?.name }}</div>
                <div class="text-xs text-zinc-400 font-mono">{{ formatDate(order.created_at, 'YYYY-MM-DD HH:mm') }}</div>
              </div>
              <div class="flex items-center gap-3 flex-shrink-0">
                <span class="text-sm font-bold font-mono text-zinc-900">{{ formatPrice(order.total_amount) }}</span>
                <ChevronRight class="w-4 h-4 text-zinc-300" />
              </div>
            </router-link>
          </template>
          <div v-else class="px-6 py-12 text-center text-sm text-zinc-400">
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
    orders.value = ordersRes.data.orders || ordersRes.data || []
  } catch (error) {
    console.error('Failed to load profile', error)
  } finally {
    loading.value = false
  }
})
</script>
