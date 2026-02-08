<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20 gap-3">
      <Loader2 class="w-8 h-8 animate-spin text-brand-green" />
      <span class="text-sm text-zinc-400">加载中...</span>
    </div>
    
    <template v-else-if="order">
      <!-- Header -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <div class="space-y-1">
          <router-link to="/orders" class="inline-flex items-center gap-1.5 text-sm text-zinc-500 hover:text-brand-green transition-colors mb-2">
            <ArrowLeft class="w-4 h-4" />
            <span>返回订单列表</span>
          </router-link>
          <h1 class="text-2xl sm:text-3xl font-bold text-zinc-900">订单详情</h1>
          <p class="text-sm text-zinc-400">
            订单号: <code class="font-mono text-zinc-600 bg-zinc-100 px-2 py-0.5 rounded-lg text-xs">{{ order.order_no }}</code>
          </p>
        </div>
        <span :class="['inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-sm font-medium self-start', getStatusClass(order.status)]">
          <span class="w-2 h-2 rounded-full" :class="getStatusDot(order.status)"></span>
          {{ getStatusText(order.status) }}
        </span>
      </div>
      
      <!-- Order Info -->
      <div class="bg-white rounded-2xl border border-zinc-100 shadow-card overflow-hidden">
        <div class="px-5 sm:px-6 py-4 border-b border-zinc-100 bg-zinc-50/50">
          <h2 class="text-xs font-semibold text-zinc-400 uppercase tracking-wider">订单信息</h2>
        </div>
        <div class="p-5 sm:p-6">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
            <div class="space-y-1">
              <div class="text-xs text-zinc-400 uppercase tracking-wider">商品名称</div>
              <div class="text-sm text-zinc-900 font-medium">{{ order.product?.name }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-xs text-zinc-400 uppercase tracking-wider">购买数量</div>
              <div class="text-sm text-zinc-900 font-medium font-mono">{{ order.quantity }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-xs text-zinc-400 uppercase tracking-wider">订单金额</div>
              <div class="text-xl text-zinc-900 font-bold font-mono tracking-tight">{{ formatPrice(order.total_amount) }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-xs text-zinc-400 uppercase tracking-wider">下单时间</div>
              <div class="text-sm text-zinc-900 font-mono">{{ formatDate(order.created_at) }}</div>
            </div>
            <div v-if="order.contact" class="space-y-1 sm:col-span-2">
              <div class="text-xs text-zinc-400 uppercase tracking-wider">联系方式</div>
              <div class="text-sm text-zinc-900">{{ order.contact }}</div>
            </div>
            <div v-if="order.remark" class="space-y-1 sm:col-span-2">
              <div class="text-xs text-zinc-400 uppercase tracking-wider">备注</div>
              <div class="text-sm text-zinc-900">{{ order.remark }}</div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Card Keys -->
      <div v-if="order.card_keys && order.card_keys.length > 0" class="bg-white rounded-2xl border border-zinc-100 shadow-card overflow-hidden">
        <div class="px-5 sm:px-6 py-4 border-b border-zinc-100 bg-zinc-50/50 flex items-center justify-between">
          <h2 class="text-xs font-semibold text-zinc-400 uppercase tracking-wider flex items-center gap-1.5">
            <Key class="w-3.5 h-3.5" />
            卡密信息
          </h2>
          <div class="flex gap-2">
            <button @click="exportCards('json')" class="text-xs text-zinc-500 hover:text-brand-green flex items-center gap-1 transition-colors">
              <Download class="w-3 h-3" />
              JSON
            </button>
            <button @click="exportCards('csv')" class="text-xs text-zinc-500 hover:text-brand-green flex items-center gap-1 transition-colors">
              <Download class="w-3 h-3" />
              CSV
            </button>
          </div>
        </div>
        <div class="p-5 sm:p-6 space-y-4">
          <div class="bg-blue-50 border border-blue-100 text-blue-700 px-4 py-3 rounded-xl flex items-start gap-2 text-sm">
            <Info class="w-4 h-4 flex-shrink-0 mt-0.5" />
            <span>请妥善保存卡密信息，点击可快速复制</span>
          </div>
          
          <div class="space-y-3">
            <div v-for="card in order.card_keys" :key="card.id" class="bg-zinc-50 rounded-xl p-4 space-y-3">
              <div class="space-y-1.5">
                <div class="text-xs text-zinc-400 uppercase tracking-wider">卡号</div>
                <div
                  @click="handleCopy(card.card_no)"
                  class="font-mono text-sm text-zinc-900 bg-white px-3 py-2.5 rounded-xl border border-zinc-200 cursor-pointer hover:border-brand-green/30 hover:bg-brand-gradient-subtle transition-all flex items-center justify-between gap-2"
                >
                  <span class="truncate">{{ card.card_no }}</span>
                  <Copy class="w-4 h-4 text-zinc-300 flex-shrink-0" />
                </div>
              </div>
              <div v-if="card.card_pwd" class="space-y-1.5">
                <div class="text-xs text-zinc-400 uppercase tracking-wider">密码</div>
                <div
                  @click="handleCopy(card.card_pwd)"
                  class="font-mono text-sm text-zinc-900 bg-white px-3 py-2.5 rounded-xl border border-zinc-200 cursor-pointer hover:border-brand-green/30 hover:bg-brand-gradient-subtle transition-all flex items-center justify-between gap-2"
                >
                  <span class="truncate">{{ card.card_pwd }}</span>
                  <Copy class="w-4 h-4 text-zinc-300 flex-shrink-0" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Actions -->
      <div class="flex flex-col sm:flex-row justify-center gap-3">
        <button
          v-if="order.status === 0"
          @click="handleRepay"
          :disabled="repaying"
          class="inline-flex items-center justify-center gap-2 px-6 py-3 bg-brand-gradient text-white font-medium rounded-xl hover:shadow-glow transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <Loader2 v-if="repaying" class="w-5 h-5 animate-spin" />
          <CreditCard v-else class="w-5 h-5" />
          <span>{{ repaying ? '处理中...' : '重新支付' }}</span>
        </button>
        
        <router-link
          to="/orders"
          class="inline-flex items-center justify-center gap-2 px-6 py-3 bg-white border border-zinc-200 text-zinc-700 font-medium rounded-xl hover:bg-zinc-50 transition-colors"
        >
          <ArrowLeft class="w-4 h-4" />
          <span>返回订单列表</span>
        </router-link>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Loader2, CreditCard, Key, Download, Info, Copy, ArrowLeft } from 'lucide-vue-next'
import api from '@/utils/api'
import { formatPrice, formatDate, copyToClipboard, exportJSON, exportCSV } from '@/utils/helpers'
import { useToastStore } from '@/stores/toast'

const route = useRoute()
const toast = useToastStore()
const loading = ref(true)
const order = ref(null)
const repaying = ref(false)

onMounted(async () => {
  try {
    const response = await api.get(`/api/orders/${route.params.orderNo}`)
    order.value = response.data
  } catch (error) {
    console.error('Failed to load order', error)
  } finally {
    loading.value = false
  }
})

function getStatusClass(s) {
  return ['bg-amber-50 text-amber-700', 'bg-blue-50 text-blue-700', 'bg-emerald-50 text-emerald-700', 'bg-zinc-100 text-zinc-500'][s] || 'bg-zinc-100 text-zinc-500'
}
function getStatusDot(s) {
  return ['bg-amber-500', 'bg-blue-500', 'bg-emerald-500', 'bg-zinc-400'][s] || 'bg-zinc-400'
}
function getStatusText(s) {
  return ['待支付', '已支付', '已完成', '已取消'][s] || '未知'
}

async function handleCopy(text) {
  await copyToClipboard(text)
  toast.success('已复制到剪贴板')
}

async function handleRepay() {
  if (repaying.value) return
  repaying.value = true
  try {
    const response = await api.post(`/api/orders/${order.value.order_no}/repay`)
    if (response.data.payment_url) {
      window.location.href = response.data.payment_url
    } else {
      toast.error('支付链接获取失败')
    }
  } catch (error) {
    toast.error(error.response?.data?.error || '发起支付失败')
  } finally {
    repaying.value = false
  }
}

function exportCards(format) {
  const cards = order.value.card_keys.map(c => ({ cardNo: c.card_no, cardPwd: c.card_pwd }))
  if (format === 'json') { exportJSON(cards, `cards_${order.value.order_no}.json`); toast.success('JSON 已导出') }
  else if (format === 'csv') { exportCSV(cards, `cards_${order.value.order_no}.csv`, ['Card No', 'Card Password']); toast.success('CSV 已导出') }
}
</script>
