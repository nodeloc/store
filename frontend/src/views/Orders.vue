<template>
  <div class="max-w-4xl mx-auto space-y-6">
    <h1 class="text-3xl font-bold text-zinc-900">我的订单</h1>
    
    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="w-8 h-8 animate-spin text-zinc-400" />
    </div>
    
    <!-- Orders Table -->
    <div v-else-if="orders.length > 0" class="bg-white border border-zinc-100 rounded-lg overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-zinc-50 border-b border-zinc-100">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">订单号</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">商品</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">金额</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">时间</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-100">
            <tr v-for="order in orders" :key="order.id" class="hover:bg-zinc-50 transition">
              <td class="px-6 py-4 whitespace-nowrap">
                <code class="text-sm font-mono text-zinc-900 bg-zinc-50 px-2 py-1 rounded">{{ order.order_no }}</code>
              </td>
              <td class="px-6 py-4">
                <div class="text-sm text-zinc-900">{{ order.product.name }}</div>
                <div class="text-xs text-zinc-500">x {{ order.quantity }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="text-sm font-mono font-semibold text-zinc-900">{{ formatPrice(order.total_amount) }}</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  :class="getStatusClass(order.status)"
                >
                  <component :is="getStatusIcon(order.status)" class="w-3 h-3 mr-1" />
                  {{ getStatusText(order.status) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-zinc-500">
                {{ formatDate(order.created_at, 'MM-DD HH:mm') }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <router-link
                  :to="`/order/${order.order_no}`"
                  class="inline-flex items-center space-x-1 text-sm text-zinc-600 hover:text-zinc-900"
                >
                  <span>查看</span>
                  <ChevronRight class="w-4 h-4" />
                </router-link>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- Pagination -->
      <Pagination
        :current-page="currentPage"
        :total-pages="totalPages"
        :total="total"
        :page-size="pageSize"
        @change="handlePageChange"
      />
    </div>
    
    <!-- Empty State -->
    <div v-else class="text-center py-16 space-y-4">
      <div class="flex justify-center">
        <div class="w-24 h-24 bg-zinc-50 rounded-full flex items-center justify-center">
          <PackageOpen class="w-12 h-12 text-zinc-400" />
        </div>
      </div>
      <div class="space-y-2">
        <h3 class="text-xl font-semibold text-zinc-900">暂无订单</h3>
        <p class="text-zinc-600">您还没有购买任何商品</p>
      </div>
      <router-link
        to="/"
        class="inline-flex items-center justify-center space-x-2 px-6 py-3 bg-zinc-900 text-white font-medium rounded-lg hover:bg-zinc-800 transition"
      >
        <ShoppingCart class="w-5 h-5" />
        <span>去购物</span>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { PackageOpen, ShoppingCart, ChevronRight, Loader2, Clock, CreditCard, CheckCircle, XCircle } from 'lucide-vue-next'
import Pagination from '@/components/Pagination.vue'
import api from '@/utils/api'
import { formatPrice, formatDate } from '@/utils/helpers'

const loading = ref(true)
const orders = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

onMounted(() => {
  fetchOrders()
})

async function fetchOrders() {
  try {
    loading.value = true
    const response = await api.get('/api/orders', {
      params: {
        page: currentPage.value,
        page_size: pageSize.value
      }
    })
    orders.value = response.data.orders || response.data || []
    total.value = response.data.total || orders.value.length
  } catch (error) {
    console.error('Failed to load orders', error)
  } finally {
    loading.value = false
  }
}

function handlePageChange(page) {
  currentPage.value = page
  fetchOrders()
}

function getStatusClass(status) {
  switch (status) {
    case 0: return 'bg-yellow-100 text-yellow-800'
    case 1: return 'bg-blue-100 text-blue-800'
    case 2: return 'bg-green-100 text-green-800'
    case 3: return 'bg-zinc-100 text-zinc-600'
    default: return 'bg-zinc-100 text-zinc-600'
  }
}

function getStatusIcon(status) {
  switch (status) {
    case 0: return Clock
    case 1: return CreditCard
    case 2: return CheckCircle
    case 3: return XCircle
    default: return XCircle
  }
}

function getStatusText(status) {
  switch (status) {
    case 0: return '待支付'
    case 1: return '已支付'
    case 2: return '已完成'
    case 3: return '已取消'
    default: return '未知'
  }
}
</script>
