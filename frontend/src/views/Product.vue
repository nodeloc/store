<template>
  <div class="max-w-3xl mx-auto space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="w-8 h-8 animate-spin text-zinc-400" />
    </div>
    
    <template v-else-if="product">
      <!-- Back Button -->
      <router-link to="/" class="inline-flex items-center space-x-2 text-sm text-zinc-600 hover:text-zinc-900">
        <ArrowLeft class="w-4 h-4" />
        <span>返回首页</span>
      </router-link>
      
      <!-- Product Detail Card -->
      <div class="bg-white border border-zinc-100 rounded-lg overflow-hidden">
        <div class="p-6 space-y-5">
          <!-- Product Header -->
          <div class="flex items-start space-x-6">
            <div class="flex-shrink-0 w-20 h-20 bg-zinc-50 rounded-lg flex items-center justify-center">
              <Package class="w-10 h-10 text-zinc-600" />
            </div>
            <div class="flex-1 space-y-2">
              <h1 class="text-3xl font-bold text-zinc-900">{{ product.name }}</h1>
              <p v-if="product.description" class="text-zinc-600">{{ product.description }}</p>
            </div>
          </div>
          
          <!-- Product Details -->
          <div class="grid grid-cols-2 gap-6 py-5 border-t border-b border-zinc-100">
            <div class="space-y-1">
              <div class="text-sm text-zinc-500">商品价格</div>
              <div class="text-3xl font-bold text-zinc-900">{{ formatPrice(product.price) }}</div>
            </div>
            <div class="space-y-1">
              <div class="text-sm text-zinc-500">可用库存</div>
              <div class="text-3xl font-bold text-zinc-900">{{ product.stock_count }}</div>
            </div>
          </div>
          
          <!-- Action Button -->
          <router-link
            v-if="product.stock_count > 0"
            :to="`/purchase/${product.id}`"
            class="inline-flex items-center justify-center space-x-2 w-full px-6 py-3 bg-zinc-900 text-white font-medium rounded-lg hover:bg-zinc-800 transition"
          >
            <ShoppingCart class="w-5 h-5" />
            <span>立即购买</span>
          </router-link>
          <button
            v-else
            disabled
            class="inline-flex items-center justify-center space-x-2 w-full px-6 py-3 bg-zinc-100 text-zinc-400 font-medium rounded-lg cursor-not-allowed"
          >
            <PackageX class="w-5 h-5" />
            <span>暂无库存</span>
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Package, PackageX, ShoppingCart, ArrowLeft, Loader2 } from 'lucide-vue-next'
import api from '@/utils/api'
import { formatPrice } from '@/utils/helpers'

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const product = ref(null)

onMounted(async () => {
  try {
    const response = await api.get(`/api/products/${route.params.id}`)
    product.value = response.data
  } catch (error) {
    console.error('Failed to load product', error)
    router.push({ name: 'NotFound' })
  } finally {
    loading.value = false
  }
})
</script>
