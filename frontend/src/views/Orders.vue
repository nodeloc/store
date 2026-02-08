<template>
  <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <h1 class="text-2xl sm:text-3xl font-bold text-zinc-900">我的订单</h1>
    
    <!-- Loading -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20 gap-3">
      <Loader2 class="w-8 h-8 animate-spin text-brand-green" />
      <span class="text-sm text-zinc-400">加载中...</span>
    </div>
    
    <!-- Orders -->
    <template v-else-if="orders.length > 0">
      <!-- Desktop Table -->
      <div class="hidden md:block bg-white rounded-2xl border border-zinc-100 shadow-card overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-zinc-50/80 border-b border-zinc-100">
              <tr>
                <th class="px-5 py-3 text-left text-xs font-semibold text-zinc-400 uppercase tracking-wider">订单号</th>
                <th class="px-5 py-3 text-left text-xs font-semibold text-zinc-400 uppercase tracking-wider">商品</th>
                <th class="px-5 py-3 text-left text-xs font-semibold text-zinc-400 uppercase tracking-wider">金额</th>
                <th class="px-5 py-3 text-left text-xs font-semibold text-zinc-400 uppercase tracking-wider">状态</th>
                <th class="px-5 py-3 text-left text-xs font-semibold text-zinc-400 uppercase tracking-wider">时间</th>
                <th class="px-5 py-3 text-right text-xs font-semibold text-zinc-400 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-100">
              <tr v-for="order in orders" :key="order.id" class="hover:bg-zinc-50/50 transition">
                <td class="px-5 py-4 whitespace-nowrap">
                  <code class="text-xs font-mono text-zinc-700 bg-zinc-100 px-2 py-1 rounded-lg">{{ order.order_no }}</code>
                </td>
                <td class="px-5 py-4">
                  <div class="text-sm font-medium text-zinc-900">{{ order.product?.name }}</div>
                  <div class="text-xs text-zinc-400">× {{ order.quantity }}</div>
                </td>
                <td class="px-5 py-4 whitespace-nowrap">
                  <span class="text-sm font-bold font-mono text-zinc-900">{{ formatPrice(order.total_amount) }}</span>
                </td>
                <td class="px-5 py-4 whitespace-nowrap">
                  <span :class="['inline-flex items-center gap-1 px-2.5 py-1 rounded-lg text-xs font-medium', getStatusClass(order.status)]">
                    <span class="w-1.5 h-1.5 rounded-full" :class="getStatusDot(order.status)"></span>
                    {{ getStatusText(order.status) }}
                  </span>
                </td>
                <td class="px-5 py-4 whitespace-nowrap text-sm text-zinc-400 font-mono">
                  {{ formatDate(order.created_at, 'MM-DD HH:mm') }}
                </td>
                <td class="px-5 py-4 whitespace-nowrap text-right">
                  <router-link :to="`/order/${order.order_no}`" class="inline-flex items-center gap-1 text-sm text-zinc-500 hover:text-brand-green transition-colors">
                    <span>查看</span>
                    <ChevronRight class="w-4 h-4" />
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <Pagination
          :current-page="currentPage"
          :total-pages="totalPages"
          :total="total"
          :page-size="pageSize"
          @change="handlePageChange"
        />
      </div>

      <!-- Mobile Cards -->
      <div class="md:hidden space-y-3">
        <router-link
          v-for="order in orders"
          :key="'m-' + order.id"
          :to="`/order/${order.order_no}`"
          class="block bg-white rounded-2xl border border-zinc-100 shadow-card p-4 hover:border-zinc-200 transition-all"
        >
          <div class="flex items-start justify-between gap-3 mb-3">
            <div class="min-w-0">
              <div class="text-sm font-semibold text-zinc-900 truncate">{{ order.product?.name }}</div>
              <code class="text-xs font-mono text-zinc-400">{{ order.order_no }}</code>
            </div>
            <span :class="['inline-flex items-center gap-1 px-2 py-1 rounded-lg text-xs font-medium flex-shrink-0', getStatusClass(order.status)]">
              <span class="w-1.5 h-1.5 rounded-full" :class="getStatusDot(order.status)"></span>
              {{ getStatusText(order.status) }}
            </span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-lg font-bold font-mono text-zinc-900">{{ formatPrice(order.total_amount) }}</span>
            <span class="text-xs text-zinc-400 font-mono">{{ formatDate(order.created_at, 'MM-DD HH:mm') }}</span>
          </div>
        </router-link>
      </div>
    </template>
    
    <!-- Empty State -->
    <div v-else class="flex flex-col items-center justify-center py-20 space-y-4">
      <div class="w-20 h-20 rounded-3xl bg-brand-gradient-subtle flex items-center justify-center">
        <PackageOpen class="w-10 h-10 text-brand-green/50" />
      </div>
      <div class="text-center space-y-1.5">
        <h3 class="text-lg font-semibold text-zinc-900">暂无订单</h3>
        <p class="text-sm text-zinc-500">您还没有购买任何商品</p>
      </div>
      <router-link to="/" class="inline-flex items-center gap-2 px-5 py-2.5 bg-brand-gradient text-white text-sm font-medium rounded-xl hover:shadow-glow transition-all duration-300">
        <ShoppingCart class="w-4 h-4" />
        <span>去购物</span>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { PackageOpen, ShoppingCart, ChevronRight, Loader2 } from 'lucide-vue-next'
import Pagination from '@/components/Pagination.vue'
import api from '@/utils/api'
import { formatPrice, formatDate } from '@/utils/helpers'

const loading = ref(true)
const orders = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

onMounted(() => { fetchOrders() })

async function fetchOrders() {
  try {
    loading.value = true
    const response = await api.get('/api/orders', { params: { page: currentPage.value, page_size: pageSize.value } })
    orders.value = response.data.orders || response.data || []
    total.value = response.data.total || orders.value.length
  } catch (error) {
    console.error('Failed to load orders', error)
  } finally {
    loading.value = false
  }
}

function handlePageChange(page) { currentPage.value = page; fetchOrders() }

function getStatusClass(s) {
  return [null, 'bg-amber-50 text-amber-700', 'bg-blue-50 text-blue-700', 'bg-emerald-50 text-emerald-700', 'bg-zinc-100 text-zinc-500'][s + 1] || 'bg-zinc-100 text-zinc-500'
}
function getStatusDot(s) {
  return [null, 'bg-amber-500', 'bg-blue-500', 'bg-emerald-500', 'bg-zinc-400'][s + 1] || 'bg-zinc-400'
}
function getStatusText(s) {
  return ['待支付', '已支付', '已完成', '已取消'][s] || '未知'
}
</script>
