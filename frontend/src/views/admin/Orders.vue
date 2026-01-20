<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-xl font-semibold text-zinc-900">订单管理</h2>
      <p class="text-sm text-zinc-600 mt-1">查看和管理所有订单</p>
    </div>
    
    <div class="bg-white rounded-lg border border-zinc-100">
      <div v-if="loading" class="p-12 text-center">
        <Loader2 class="w-8 h-8 animate-spin text-zinc-400 mx-auto" />
      </div>
      
      <div v-else-if="orders.length === 0" class="p-12 text-center">
        <ShoppingCart class="w-12 h-12 text-zinc-300 mx-auto mb-4" />
        <p class="text-zinc-600">暂无订单</p>
      </div>
      
      <table v-else class="w-full text-sm">
        <thead class="bg-zinc-50 border-b border-zinc-100">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">订单号</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">用户</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">商品</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">金额</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">状态</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">创建时间</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-zinc-100">
          <tr v-for="order in orders" :key="order.id" class="hover:bg-zinc-50">
            <td class="px-6 py-4 font-mono text-zinc-900">{{ order.order_no }}</td>
            <td class="px-6 py-4 text-zinc-600">{{ order.user?.name || '-' }}</td>
            <td class="px-6 py-4 text-zinc-600">{{ order.product?.name || '-' }}</td>
            <td class="px-6 py-4 text-zinc-900">¥{{ order.total_amount }}</td>
            <td class="px-6 py-4">
              <span :class="getStatusClass(order.status)" class="px-2 py-1 text-xs font-medium rounded">
                {{ getStatusText(order.status) }}
              </span>
            </td>
            <td class="px-6 py-4 text-zinc-600">{{ formatDate(order.created_at) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ShoppingCart, Loader2 } from 'lucide-vue-next'
import api from '@/utils/api'
import { useToast } from '@/stores/toast'

const toast = useToast()
const orders = ref([])
const loading = ref(false)

onMounted(() => {
  fetchOrders()
})

async function fetchOrders() {
  try {
    loading.value = true
    const response = await api.get('/api/admin/orders', { params: { page: 1, page_size: 50 } })
    orders.value = response.data.orders || []
  } catch (error) {
    toast.error('加载订单失败')
  } finally {
    loading.value = false
  }
}

function getStatusText(status) {
  const statusMap = { 0: '待支付', 1: '已支付', 2: '已完成', 3: '已取消' }
  return statusMap[status] || '未知'
}

function getStatusClass(status) {
  const classMap = {
    0: 'bg-yellow-100 text-yellow-800',
    1: 'bg-green-100 text-green-800',
    2: 'bg-blue-100 text-blue-800',
    3: 'bg-zinc-100 text-zinc-800'
  }
  return classMap[status] || 'bg-zinc-100 text-zinc-800'
}

function formatDate(date) {
  return new Date(date).toLocaleString('zh-CN', { 
    year: 'numeric', 
    month: '2-digit', 
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>
