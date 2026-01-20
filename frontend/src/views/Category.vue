<template>
  <div class="space-y-8">
    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="w-8 h-8 animate-spin text-zinc-400" />
    </div>
    
    <template v-else-if="category">
      <!-- Category Header -->
      <div class="space-y-2">
        <router-link to="/" class="inline-flex items-center space-x-2 text-sm text-zinc-600 hover:text-zinc-900 mb-4">
          <ArrowLeft class="w-4 h-4" />
          <span>返回首页</span>
        </router-link>
        <h1 class="text-4xl font-bold text-zinc-900">{{ category.name }}</h1>
        <p v-if="category.description" class="text-lg text-zinc-600">{{ category.description }}</p>
      </div>
      
      <!-- Products Grid -->
      <div v-if="products.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <router-link
          v-for="product in products"
          :key="product.id"
          :to="`/product/${product.id}`"
          class="group block bg-white border border-zinc-100 rounded-lg hover:border-zinc-200 hover:shadow-md transition-all"
        >
          <div class="p-6 space-y-4">
            <div class="flex items-center justify-center w-12 h-12 bg-zinc-50 rounded-lg group-hover:bg-zinc-100 transition">
              <Package class="w-6 h-6 text-zinc-600" />
            </div>
            <div class="space-y-2">
              <h3 class="text-lg font-semibold text-zinc-900 group-hover:text-zinc-700">{{ product.name }}</h3>
              <p v-if="product.description" class="text-sm text-zinc-600 line-clamp-2">{{ product.description }}</p>
            </div>
            <div class="flex items-center justify-between pt-4 border-t border-zinc-100">
              <div class="space-y-1">
                <div class="text-2xl font-bold text-zinc-900">{{ formatPrice(product.price) }}</div>
                <div class="text-xs text-zinc-500">库存: {{ product.stock_count }}</div>
              </div>
              <div class="flex items-center space-x-1 text-sm text-zinc-600 group-hover:text-zinc-900">
                <span>购买</span>
                <ChevronRight class="w-4 h-4" />
              </div>
            </div>
          </div>
        </router-link>
      </div>
      
      <!-- Empty State -->
      <div v-else class="text-center py-16 space-y-4">
        <div class="flex justify-center">
          <div class="w-24 h-24 bg-zinc-50 rounded-full flex items-center justify-center">
            <PackageOpen class="w-12 h-12 text-zinc-400" />
          </div>
        </div>
        <div class="space-y-2">
          <h3 class="text-xl font-semibold text-zinc-900">该分类下暂无商品</h3>
          <p class="text-zinc-600">请查看其他分类</p>
        </div>
        <router-link
          to="/"
          class="inline-flex items-center space-x-2 px-6 py-3 bg-zinc-900 text-white font-medium rounded-lg hover:bg-zinc-800 transition"
        >
          <ArrowLeft class="w-5 h-5" />
          <span>返回首页</span>
        </router-link>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Package, PackageOpen, ChevronRight, ArrowLeft, Loader2 } from 'lucide-vue-next'
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
