<template>
  <div class="max-w-3xl mx-auto space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="w-8 h-8 animate-spin text-zinc-400" />
    </div>
    
    <template v-else-if="order">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div class="space-y-1">
          <h1 class="text-3xl font-bold text-zinc-900">订单详情</h1>
          <p class="text-sm text-zinc-500">
            订单号: <code class="font-mono text-zinc-900 bg-zinc-50 px-2 py-1 rounded">{{ order.order_no }}</code>
          </p>
        </div>
        <div>
          <span
            class="inline-flex items-center px-3 py-1.5 rounded-full text-sm font-medium"
            :class="getStatusClass(order.status)"
          >
            <component :is="getStatusIcon(order.status)" class="w-4 h-4 mr-1.5" />
            {{ getStatusText(order.status) }}
          </span>
        </div>
      </div>
      
      <!-- Order Information -->
      <div class="bg-white border border-zinc-100 rounded-lg overflow-hidden">
        <div class="px-6 py-4 bg-zinc-50 border-b border-zinc-100">
          <h2 class="text-sm font-medium text-zinc-900 uppercase tracking-wide">订单信息</h2>
        </div>
        <div class="p-6 space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1">
              <div class="text-xs text-zinc-500 uppercase tracking-wide">商品名称</div>
              <div class="text-sm text-zinc-900 font-medium">{{ order.product.name }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-xs text-zinc-500 uppercase tracking-wide">购买数量</div>
              <div class="text-sm text-zinc-900 font-medium font-mono">{{ order.quantity }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-xs text-zinc-500 uppercase tracking-wide">订单金额</div>
              <div class="text-lg text-zinc-900 font-bold font-mono">{{ formatPrice(order.total_amount) }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-xs text-zinc-500 uppercase tracking-wide">下单时间</div>
              <div class="text-sm text-zinc-900 font-mono">{{ formatDate(order.created_at) }}</div>
            </div>
            <div v-if="order.contact" class="space-y-1 col-span-2">
              <div class="text-xs text-zinc-500 uppercase tracking-wide">联系方式</div>
              <div class="text-sm text-zinc-900">{{ order.contact }}</div>
            </div>
            <div v-if="order.remark" class="space-y-1 col-span-2">
              <div class="text-xs text-zinc-500 uppercase tracking-wide">备注</div>
              <div class="text-sm text-zinc-900">{{ order.remark }}</div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Card Keys -->
      <div v-if="order.card_keys && order.card_keys.length > 0" class="bg-white border border-zinc-100 rounded-lg overflow-hidden">
        <div class="px-6 py-4 bg-gradient-to-r from-zinc-50 to-zinc-100 border-b border-zinc-100 flex items-center justify-between">
          <h2 class="text-sm font-medium text-zinc-900 uppercase tracking-wide flex items-center space-x-2">
            <Key class="w-4 h-4" />
            <span>卡密信息</span>
          </h2>
          <div class="flex space-x-2">
            <button @click="exportCards('json')" class="text-xs text-zinc-600 hover:text-zinc-900 flex items-center space-x-1">
              <Download class="w-3 h-3" />
              <span>导出JSON</span>
            </button>
            <button @click="exportCards('csv')" class="text-xs text-zinc-600 hover:text-zinc-900 flex items-center space-x-1">
              <Download class="w-3 h-3" />
              <span>导出CSV</span>
            </button>
          </div>
        </div>
        <div class="p-6 space-y-4">
          <div class="bg-blue-50 border border-blue-200 text-blue-800 px-4 py-3 rounded-lg flex items-start space-x-2 text-sm">
            <Info class="w-5 h-5 flex-shrink-0 mt-0.5" />
            <span>请妥善保存以下卡密信息，点击卡密可快速复制</span>
          </div>
          
          <div class="space-y-3">
            <div v-for="card in order.card_keys" :key="card.id" class="group bg-zinc-50 border border-zinc-200 rounded-lg p-4 hover:border-zinc-300 transition">
              <div class="space-y-3">
                <div class="space-y-1.5">
                  <div class="text-xs text-zinc-500 uppercase tracking-wide">卡号</div>
                  <div
                    @click="handleCopy(card.card_no)"
                    class="font-mono text-sm text-zinc-900 bg-white px-3 py-2 rounded border border-zinc-200 cursor-pointer hover:border-zinc-400 transition flex items-center justify-between"
                  >
                    <span>{{ card.card_no }}</span>
                    <Copy class="w-4 h-4 text-zinc-400" />
                  </div>
                </div>
                <div v-if="card.card_pwd" class="space-y-1.5">
                  <div class="text-xs text-zinc-500 uppercase tracking-wide">密码</div>
                  <div
                    @click="handleCopy(card.card_pwd)"
                    class="font-mono text-sm text-zinc-900 bg-white px-3 py-2 rounded border border-zinc-200 cursor-pointer hover:border-zinc-400 transition flex items-center justify-between"
                  >
                    <span>{{ card.card_pwd }}</span>
                    <Copy class="w-4 h-4 text-zinc-400" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Actions -->
      <div class="flex justify-center space-x-4">
        <!-- Repay Button (for pending payment orders) -->
        <button
          v-if="order.status === 0"
          @click="handleRepay"
          :disabled="repaying"
          class="inline-flex items-center space-x-2 px-6 py-3 bg-zinc-900 text-white font-medium rounded-lg hover:bg-zinc-800 transition disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <Loader2 v-if="repaying" class="w-5 h-5 animate-spin" />
          <CreditCard v-else class="w-5 h-5" />
          <span>{{ repaying ? '处理中...' : '重新支付' }}</span>
        </button>
        
        <router-link
          to="/orders"
          class="inline-flex items-center space-x-2 px-6 py-3 bg-zinc-100 text-zinc-900 font-medium rounded-lg hover:bg-zinc-200 transition"
        >
          <ArrowLeft class="w-5 h-5" />
          <span>返回订单列表</span>
        </router-link>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import {
  Loader2, Clock, CreditCard, CheckCircle, XCircle, Key, Download, Info, Copy, ArrowLeft
} from 'lucide-vue-next'
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
      // Redirect to payment gateway
      window.location.href = response.data.payment_url
    } else {
      toast.error('支付链接获取失败')
    }
  } catch (error) {
    console.error('Failed to repay', error)
    toast.error(error.response?.data?.error || '发起支付失败')
  } finally {
    repaying.value = false
  }
}

function exportCards(format) {
  const cards = order.value.card_keys.map(card => ({
    cardNo: card.card_no,
    cardPwd: card.card_pwd
  }))
  
  if (format === 'json') {
    exportJSON(cards, `cards_${order.value.order_no}.json`)
    toast.success('JSON 文件已导出')
  } else if (format === 'csv') {
    exportCSV(cards, `cards_${order.value.order_no}.csv`, ['Card No', 'Card Password'])
    toast.success('CSV 文件已导出')
  }
}
</script>
