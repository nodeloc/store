<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20 gap-3">
      <Loader2 class="w-8 h-8 animate-spin text-brand-green" />
      <span class="text-sm text-zinc-400">加载中...</span>
    </div>
    
    <template v-else-if="product">
      <!-- Back Button -->
      <router-link to="/" class="inline-flex items-center gap-1.5 text-sm text-zinc-500 hover:text-brand-green transition-colors">
        <ArrowLeft class="w-4 h-4" />
        <span>返回首页</span>
      </router-link>
      
      <!-- Product Detail Card -->
      <div class="bg-white rounded-2xl border border-zinc-100 overflow-hidden shadow-card">
        <!-- Product Image -->
        <div v-if="product.image" class="aspect-video bg-zinc-50 overflow-hidden">
          <img :src="product.image" :alt="product.name" class="w-full h-full object-cover" />
        </div>

        <div class="p-5 sm:p-8 space-y-6">
          <!-- Product Header -->
          <div class="flex items-start gap-5">
            <div v-if="!product.image" class="flex-shrink-0 w-16 h-16 rounded-2xl bg-brand-gradient-subtle flex items-center justify-center">
              <Package class="w-8 h-8 text-brand-green/60" />
            </div>
            <div class="flex-1 space-y-2">
              <h1 class="text-2xl sm:text-3xl font-bold text-zinc-900">{{ product.name }}</h1>
              <p v-if="product.description" class="text-zinc-500 leading-relaxed">{{ product.description }}</p>
            </div>
          </div>
          
          <!-- Product Details -->
          <div class="grid grid-cols-2 gap-6 py-5 border-t border-b border-zinc-100">
            <div class="space-y-1.5">
              <div class="text-xs font-medium text-zinc-400 uppercase tracking-wider">商品价格</div>
              <div class="text-3xl font-bold text-zinc-900 font-mono tracking-tight">{{ formatPrice(product.price) }}</div>
              <div v-if="product.orig_price && product.orig_price > product.price" class="text-sm text-zinc-400 line-through font-mono">
                {{ formatPrice(product.orig_price) }}
              </div>
            </div>
            <div class="space-y-1.5">
              <div class="text-xs font-medium text-zinc-400 uppercase tracking-wider">可用库存</div>
              <div class="text-3xl font-bold font-mono tracking-tight" :class="[
                product.stock_count > 10 ? 'text-zinc-900' : product.stock_count > 0 ? 'text-amber-600' : 'text-red-500'
              ]">
                {{ product.stock_count }}
              </div>
              <div
                :class="[
                  'inline-flex items-center gap-1 text-xs font-medium',
                  product.stock_count > 10 ? 'text-emerald-600' : product.stock_count > 0 ? 'text-amber-600' : 'text-red-500'
                ]"
              >
                <span class="w-1.5 h-1.5 rounded-full" :class="[
                  product.stock_count > 10 ? 'bg-emerald-500' : product.stock_count > 0 ? 'bg-amber-500' : 'bg-red-500'
                ]"></span>
                {{ product.stock_count > 10 ? '库存充足' : product.stock_count > 0 ? '库存紧张' : '已售罄' }}
              </div>
            </div>
          </div>
          
          <!-- Action Button -->
          <router-link
            v-if="product.stock_count > 0"
            :to="`/purchase/${product.id}`"
            class="flex items-center justify-center gap-2 w-full px-6 py-3.5 bg-brand-gradient text-white font-medium rounded-xl hover:shadow-glow transition-all duration-300 hover:scale-[1.01]"
          >
            <ShoppingCart class="w-5 h-5" />
            <span>立即购买</span>
          </router-link>
          <button
            v-else
            disabled
            class="flex items-center justify-center gap-2 w-full px-6 py-3.5 bg-zinc-100 text-zinc-400 font-medium rounded-xl cursor-not-allowed"
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
