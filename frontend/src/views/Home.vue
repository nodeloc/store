<template>
  <div class="space-y-6">
    <!-- Hero Section with Typewriter Effect -->
    <section class="text-center py-8 space-y-4 fade-in">
      <div class="space-y-2">
        <h1 class="text-5xl font-bold text-zinc-900">
          {{ settings.site_name || 'NodeLoc 社区发卡' }}
        </h1>
        <div class="h-8 flex items-center justify-center">
          <p v-if="showTypewriter" class="text-xl text-zinc-600 typewriter inline-block">
            {{ typewriterText }}
          </p>
          <p v-else class="text-xl text-zinc-600">
            {{ settings.site_description || '基于 NodeLoc OAuth 的社区发卡系统' }}
          </p>
        </div>
      </div>
      
      <!-- Stats -->
      <div class="flex justify-center space-x-8 pt-2">
        <div class="text-center slide-up">
          <div class="text-2xl font-bold text-zinc-900 font-mono-data">{{ totalProducts }}</div>
          <div class="text-sm text-zinc-500">商品</div>
        </div>
        <div class="text-center slide-up" style="animation-delay: 0.1s">
          <div class="text-2xl font-bold text-zinc-900 font-mono-data">{{ totalCategories }}</div>
          <div class="text-sm text-zinc-500">分类</div>
        </div>
      </div>
    </section>
    
    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <Loader2 class="w-8 h-8 animate-spin text-zinc-400" />
    </div>
    
    <!-- Categories with Products -->
    <template v-else-if="categoriesWithProducts.length > 0">
      <section v-for="category in categoriesWithProducts" :key="category.id" class="space-y-4 slide-up">
        <div class="flex items-center justify-between">
          <div class="space-y-1">
            <h2 class="text-2xl font-semibold text-zinc-900">{{ category.name }}</h2>
            <p v-if="category.description" class="text-sm text-zinc-600">{{ category.description }}</p>
          </div>
          <router-link :to="`/category/${category.id}`" class="text-sm text-zinc-600 hover:text-zinc-900 flex items-center space-x-1">
            <span>查看全部</span>
            <ArrowRight class="w-4 h-4" />
          </router-link>
        </div>
        
        <!-- Products Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <router-link
            v-for="product in category.products"
            :key="product.id"
            :to="`/product/${product.id}`"
            class="group block bg-white border border-zinc-100 rounded-lg hover:border-zinc-900 hover:shadow-lg transition-all duration-300"
          >
            <div class="p-4 space-y-3">
              <!-- Product Icon -->
              <div class="flex items-center justify-center w-12 h-12 bg-zinc-50 rounded-lg group-hover:bg-zinc-900 transition-colors duration-300">
                <Package class="w-6 h-6 text-zinc-600 group-hover:text-white transition-colors duration-300" />
              </div>
              
              <!-- Product Info -->
              <div class="space-y-2">
                <h3 class="text-lg font-semibold text-zinc-900">{{ product.name }}</h3>
                <p v-if="product.description" class="text-sm text-zinc-600 line-clamp-2">{{ product.description }}</p>
              </div>
              
              <!-- Product Footer -->
              <div class="flex items-center justify-between pt-4 border-t border-zinc-100">
                <div class="space-y-1">
                  <div class="text-2xl font-bold text-zinc-900">{{ formatPrice(product.price) }}</div>
                  <div class="text-xs text-zinc-500 font-mono-data">库存: {{ product.stock_count }}</div>
                </div>
                <div class="flex items-center space-x-1 text-sm text-zinc-600 group-hover:text-zinc-900 transition">
                  <span>购买</span>
                  <ChevronRight class="w-4 h-4" />
                </div>
              </div>
            </div>
          </router-link>
        </div>
      </section>
    </template>
    
    <!-- Empty State -->
    <div v-else class="text-center py-16 space-y-4">
      <Package class="w-16 h-16 text-zinc-300 mx-auto" />
      <div class="space-y-2">
        <p class="text-lg text-zinc-600">暂无商品</p>
        <p class="text-sm text-zinc-500">管理员可以在后台添加商品</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Package, ArrowRight, ChevronRight, Loader2 } from 'lucide-vue-next'
import api from '@/utils/api'
import { formatPrice } from '@/utils/helpers'

const settings = ref({})
const categoriesWithProducts = ref([])
const loading = ref(true)
const showTypewriter = ref(true)
const typewriterText = ref('')

const totalProducts = computed(() => {
  return categoriesWithProducts.value.reduce((sum, cat) => sum + (cat.products?.length || 0), 0)
})

const totalCategories = computed(() => categoriesWithProducts.value.length)

onMounted(async () => {
  try {
    // Fetch settings
    const settingsRes = await api.get('/api/settings')
    settings.value = settingsRes.data
    
    // Start typewriter effect
    const fullText = settings.value.site_description || '基于 NodeLoc OAuth 的社区发卡系统'
    let index = 0
    const typeSpeed = 80
    
    const typeInterval = setInterval(() => {
      if (index < fullText.length) {
        typewriterText.value = fullText.substring(0, index + 1)
        index++
      } else {
        clearInterval(typeInterval)
        setTimeout(() => {
          showTypewriter.value = false
        }, 1000)
      }
    }, typeSpeed)
    
    // Fetch categories with products
    const categoriesRes = await api.get('/api/categories/with-products')
    categoriesWithProducts.value = categoriesRes.data
  } catch (error) {
    console.error('Failed to load data', error)
  } finally {
    loading.value = false
  }
})
</script>
