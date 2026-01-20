<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="w-8 h-8 animate-spin text-zinc-400" />
    </div>
    
    <template v-else-if="product">
      <!-- Back Button -->
      <router-link :to="`/product/${product.id}`" class="inline-flex items-center space-x-2 text-sm text-zinc-600 hover:text-zinc-900">
        <ArrowLeft class="w-4 h-4" />
        <span>返回商品详情</span>
      </router-link>
      
      <h1 class="text-3xl font-bold text-zinc-900">确认购买</h1>
      
      <!-- Product Summary -->
      <div class="bg-zinc-50 border border-zinc-100 rounded-lg p-6">
        <div class="flex items-center space-x-4">
          <div class="flex-shrink-0 w-16 h-16 bg-white rounded-lg flex items-center justify-center border border-zinc-200">
            <Package class="w-8 h-8 text-zinc-600" />
          </div>
          <div class="flex-1">
            <h2 class="text-lg font-semibold text-zinc-900">{{ product.name }}</h2>
            <div class="text-2xl font-bold text-zinc-900 mt-1">{{ formatPrice(product.price) }}</div>
            <div class="text-sm text-zinc-500 mt-1">库存: {{ product.stock_count }}</div>
          </div>
        </div>
      </div>
      
      <!-- Purchase Form -->
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <!-- Quantity -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-zinc-900">购买数量</label>
          <input
            v-model.number="form.quantity"
            type="number"
            min="1"
            :max="product.stock_count"
            class="block w-full px-4 py-2 border border-zinc-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900 focus:border-transparent"
          >
        </div>
        
        <!-- Contact -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-zinc-900">联系方式（选填）</label>
          <input
            v-model="form.contact"
            type="text"
            placeholder="邮箱或其他联系方式"
            class="block w-full px-4 py-2 border border-zinc-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900 focus:border-transparent"
          >
          <p class="text-xs text-zinc-500">方便我们联系您处理售后问题</p>
        </div>
        
        <!-- Remark -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-zinc-900">备注（选填）</label>
          <textarea
            v-model="form.remark"
            rows="3"
            placeholder="有什么需要备注的吗？"
            class="block w-full px-4 py-2 border border-zinc-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900 focus:border-transparent resize-none"
          ></textarea>
        </div>
        
        <!-- Total Amount -->
        <div class="bg-gradient-to-br from-zinc-50 to-zinc-100 border border-zinc-200 rounded-lg p-6 space-y-3">
          <div class="flex justify-between text-sm">
            <span class="text-zinc-600">商品单价</span>
            <span class="font-mono text-zinc-900">{{ formatPrice(product.price) }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-zinc-600">购买数量</span>
            <span class="font-mono text-zinc-900">{{ form.quantity }}</span>
          </div>
          <div class="flex justify-between text-lg font-bold pt-3 border-t border-zinc-200">
            <span class="text-zinc-900">应付金额</span>
            <span class="font-mono text-zinc-900">{{ formatPrice(totalAmount) }}</span>
          </div>
        </div>
        
        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="submitting"
          class="w-full inline-flex items-center justify-center space-x-2 px-6 py-3 bg-zinc-900 text-white font-medium rounded-lg hover:bg-zinc-800 transition disabled:opacity-50 disabled:cursor-not-allowed"
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
const form = ref({
  quantity: 1,
  contact: '',
  remark: ''
})

const totalAmount = computed(() => {
  return product.value ? product.value.price * form.value.quantity : 0
})

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
      // Redirect to payment
      window.location.href = response.data.payment_url
    } else {
      // Free or completed order
      router.push({ name: 'OrderDetail', params: { orderNo: response.data.order_no } })
    }
  } catch (error) {
    toast.error(error.response?.data?.error || '创建订单失败')
  } finally {
    submitting.value = false
  }
}
</script>
