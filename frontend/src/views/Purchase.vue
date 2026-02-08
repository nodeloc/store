<template>
  <div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20 gap-3">
      <Loader2 class="w-8 h-8 animate-spin text-brand-green" />
      <span class="text-sm text-zinc-400">加载中...</span>
    </div>
    
    <template v-else-if="product">
      <!-- Back -->
      <router-link :to="`/product/${product.id}`" class="inline-flex items-center gap-1.5 text-sm text-zinc-500 hover:text-brand-green transition-colors">
        <ArrowLeft class="w-4 h-4" />
        <span>返回商品详情</span>
      </router-link>
      
      <h1 class="text-2xl sm:text-3xl font-bold text-zinc-900">确认购买</h1>
      
      <!-- Product Summary -->
      <div class="bg-white rounded-2xl border border-zinc-100 shadow-card p-5 sm:p-6">
        <div class="flex items-center gap-4">
          <div class="flex-shrink-0 w-14 h-14 rounded-2xl bg-brand-gradient-subtle flex items-center justify-center">
            <Package class="w-7 h-7 text-brand-green/60" />
          </div>
          <div class="flex-1 min-w-0">
            <h2 class="text-lg font-semibold text-zinc-900 truncate">{{ product.name }}</h2>
            <div class="text-2xl font-bold text-zinc-900 font-mono tracking-tight mt-1">{{ formatPrice(product.price) }}</div>
            <div class="flex items-center gap-1 text-xs text-zinc-400 mt-1">
              <span class="w-1.5 h-1.5 rounded-full" :class="[product.stock_count > 10 ? 'bg-emerald-500' : product.stock_count > 0 ? 'bg-amber-500' : 'bg-red-500']"></span>
              库存: {{ product.stock_count }}
            </div>
          </div>
        </div>
      </div>
      
      <!-- Purchase Form -->
      <form @submit.prevent="handleSubmit" class="space-y-5">
        <div class="bg-white rounded-2xl border border-zinc-100 shadow-card p-5 sm:p-6 space-y-5">
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1.5">购买数量</label>
            <input
              v-model.number="form.quantity"
              type="number"
              min="1"
              :max="product.stock_count"
              class="block w-full px-4 py-2.5 border border-zinc-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-green/20 focus:border-brand-green transition-colors"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1.5">联系方式<span class="text-zinc-400 font-normal">（选填）</span></label>
            <input
              v-model="form.contact"
              type="text"
              placeholder="邮箱或其他联系方式"
              class="block w-full px-4 py-2.5 border border-zinc-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-green/20 focus:border-brand-green transition-colors"
            />
            <p class="text-xs text-zinc-400 mt-1">方便我们联系您处理售后问题</p>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1.5">备注<span class="text-zinc-400 font-normal">（选填）</span></label>
            <textarea
              v-model="form.remark"
              rows="3"
              placeholder="有什么需要备注的吗？"
              class="block w-full px-4 py-2.5 border border-zinc-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-green/20 focus:border-brand-green transition-colors resize-none"
            ></textarea>
          </div>
        </div>
        
        <!-- Total -->
        <div class="bg-white rounded-2xl border border-zinc-100 shadow-card p-5 sm:p-6 space-y-3">
          <div class="flex justify-between text-sm">
            <span class="text-zinc-500">商品单价</span>
            <span class="font-mono text-zinc-900">{{ formatPrice(product.price) }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-zinc-500">购买数量</span>
            <span class="font-mono text-zinc-900">× {{ form.quantity }}</span>
          </div>
          <div class="flex justify-between text-lg font-bold pt-3 border-t border-zinc-100">
            <span class="text-zinc-900">应付金额</span>
            <span class="font-mono text-gradient">{{ formatPrice(totalAmount) }}</span>
          </div>
        </div>
        
        <button
          type="submit"
          :disabled="submitting"
          class="w-full flex items-center justify-center gap-2 px-6 py-3.5 bg-brand-gradient text-white font-medium rounded-xl hover:shadow-glow transition-all duration-300 hover:scale-[1.01] disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100"
        >
          <Loader2 v-if="submitting" class="w-5 h-5 animate-spin" />
          <CreditCard v-else class="w-5 h-5" />
          <span>{{ submitting ? '创建订单中...' : '确认购买' }}</span>
        </button>
      </form>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Package, ArrowLeft, CreditCard, Loader2 } from 'lucide-vue-next'
import api from '@/utils/api'
import { formatPrice } from '@/utils/helpers'
import { useToastStore } from '@/stores/toast'

const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const loading = ref(true)
const submitting = ref(false)
const product = ref(null)
const form = ref({ quantity: 1, contact: '', remark: '' })

const totalAmount = computed(() => product.value ? product.value.price * form.value.quantity : 0)

onMounted(async () => {
  try {
    const response = await api.get(`/api/products/${route.params.id}`)
    product.value = response.data
    if (!product.value.is_active || product.value.stock_count <= 0) {
      toast.error('商品暂不可购买')
      router.push({ name: 'Product', params: { id: route.params.id } })
    }
  } catch (error) {
    console.error('Failed to load product', error)
    router.push({ name: 'NotFound' })
  } finally {
    loading.value = false
  }
})

async function handleSubmit() {
  try {
    submitting.value = true
    const response = await api.post('/api/orders/create', {
      product_id: product.value.id,
      quantity: form.value.quantity,
      contact: form.value.contact,
      remark: form.value.remark
    })
    if (response.data.payment_url) {
      window.location.href = response.data.payment_url
    } else {
      router.push({ name: 'OrderDetail', params: { orderNo: response.data.order_no } })
    }
  } catch (error) {
    toast.error(error.response?.data?.error || '创建订单失败')
  } finally {
    submitting.value = false
  }
}
</script>
