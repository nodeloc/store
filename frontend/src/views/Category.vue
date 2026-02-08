<template>
  <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20 gap-3">
      <Loader2 class="w-8 h-8 animate-spin text-brand-green" />
      <span class="text-sm text-zinc-400">加载中...</span>
    </div>
    
    <template v-else-if="category">
      <!-- Category Header -->
      <div class="space-y-3">
        <router-link to="/" class="inline-flex items-center gap-1.5 text-sm text-zinc-500 hover:text-brand-green transition-colors">
          <ArrowLeft class="w-4 h-4" />
          <span>返回首页</span>
        </router-link>
        <h1 class="text-2xl sm:text-3xl font-bold text-zinc-900">{{ category.name }}</h1>
        <p v-if="category.description" class="text-base text-zinc-500 leading-relaxed">{{ category.description }}</p>
      </div>
      
      <!-- Products Grid -->
      <div v-if="products.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-5 stagger-in">
        <router-link
          v-for="product in products"
          :key="product.id"
          :to="`/product/${product.id}`"
          class="product-card group"
        >
          <div class="p-5 sm:p-6 space-y-4">
            <div class="w-12 h-12 rounded-2xl bg-brand-gradient-subtle flex items-center justify-center group-hover:bg-brand-gradient-hover transition-colors">
              <Package class="w-6 h-6 text-brand-green/70" />
            </div>
            <div class="space-y-1.5">
              <h3 class="text-base font-semibold text-zinc-900 group-hover:text-brand-green transition-colors">{{ product.name }}</h3>
              <p v-if="product.description" class="text-sm text-zinc-500 line-clamp-2 leading-relaxed">{{ product.description }}</p>
            </div>
            <div class="flex items-end justify-between pt-4 border-t border-zinc-100/80">
              <div class="space-y-1">
                <div class="text-xl font-bold text-zinc-900 font-mono tracking-tight">{{ formatPrice(product.price) }}</div>
                <div
                  :class="[
                    'inline-flex items-center gap-1 text-xs font-medium',
                    product.stock_count > 0 ? 'text-zinc-400' : 'text-red-500'
                  ]"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :class="[
                    product.stock_count > 10 ? 'bg-emerald-500' : product.stock_count > 0 ? 'bg-amber-500' : 'bg-red-500'
                  ]"></span>
                  库存: {{ product.stock_count }}
                </div>
              </div>
              <div class="flex items-center gap-1 text-sm text-zinc-400 group-hover:text-brand-green transition-colors">
                <span>购买</span>
                <ArrowUpRight class="w-4 h-4 transition-transform group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
              </div>
            </div>
          </div>
        </router-link>
      </div>
      
      <!-- Empty State -->
      <div v-else class="flex flex-col items-center justify-center py-20 space-y-4">
        <div class="w-20 h-20 rounded-3xl bg-brand-gradient-subtle flex items-center justify-center">
          <PackageOpen class="w-10 h-10 text-brand-green/50" />
        </div>
        <div class="text-center space-y-1.5">
          <h3 class="text-lg font-semibold text-zinc-900">该分类下暂无商品</h3>
          <p class="text-sm text-zinc-500">请查看其他分类</p>
        </div>
        <router-link
          to="/"
          class="inline-flex items-center gap-2 px-5 py-2.5 bg-brand-gradient text-white text-sm font-medium rounded-xl hover:shadow-glow transition-all duration-300"
        >
          <ArrowLeft class="w-4 h-4" />
          <span>返回首页</span>
        </router-link>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Package, PackageOpen, ArrowUpRight, ArrowLeft, Loader2 } from 'lucide-vue-next'
import api from '@/utils/api'
import { formatPrice } from '@/utils/helpers'

const route = useRoute()
const loading = ref(true)
const category = ref(null)
const products = ref([])

onMounted(async () => {
  try {
    const [categoryRes, productsRes] = await Promise.all([
      api.get(`/api/categories/${route.params.id}`),
      api.get(`/api/products?category_id=${route.params.id}`)
    ])
    category.value = categoryRes.data
    products.value = productsRes.data
  } catch (error) {
    console.error('Failed to load category', error)
  } finally {
    loading.value = false
  }
})
</script>
