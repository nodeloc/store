<template>
  <div>
    <!-- Info Bar -->
    <div v-if="!loading && (settings.announcement || settings.site_description)" class="bg-white border-b border-zinc-100">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 sm:py-4">
        <div class="flex items-center justify-between gap-4">
          <div class="min-w-0">
            <p v-if="settings.announcement" class="text-sm text-amber-700 bg-amber-50 px-3 py-1 rounded-lg truncate">
              <span class="font-medium">📢</span> {{ settings.announcement }}
            </p>
            <p v-else-if="settings.site_description" class="text-sm text-zinc-500 truncate sm:hidden">
              {{ settings.site_description }}
            </p>
          </div>
          <div class="hidden sm:flex items-center gap-4 flex-shrink-0">
            <div class="flex items-center gap-1.5 text-sm text-zinc-400">
              <Package class="w-3.5 h-3.5 text-brand-green" />
              <span><strong class="text-zinc-700 font-semibold">{{ totalProducts }}</strong> 件商品</span>
            </div>
            <div class="w-px h-3 bg-zinc-200"></div>
            <div class="flex items-center gap-1.5 text-sm text-zinc-400">
              <Grid3X3 class="w-3.5 h-3.5 text-brand-orange" />
              <span><strong class="text-zinc-700 font-semibold">{{ categories.length }}</strong> 个分类</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile Category Tabs -->
    <div class="lg:hidden bg-white/80 backdrop-blur-sm border-b border-zinc-100 sticky top-14 sm:top-16 z-20">
      <div class="max-w-7xl mx-auto">
        <div class="flex items-center gap-2 py-3 px-4 sm:px-6 overflow-x-auto scrollbar-hide">
          <button
            @click="selectedCategoryId = null"
            :class="[
              'flex-shrink-0 inline-flex items-center gap-1.5 px-4 py-2 text-sm font-medium rounded-full transition-all duration-200 whitespace-nowrap',
              selectedCategoryId === null
                ? 'cat-pill-active'
                : 'bg-zinc-100 text-zinc-600 hover:bg-zinc-200/70'
            ]"
          >
            全部
            <span :class="['text-xs', selectedCategoryId === null ? 'text-white/70' : 'text-zinc-400']">{{ totalProducts }}</span>
          </button>
          <button
            v-for="category in categories"
            :key="'mobile-' + category.id"
            @click="selectedCategoryId = category.id"
            :class="[
              'flex-shrink-0 inline-flex items-center gap-1.5 px-4 py-2 text-sm font-medium rounded-full transition-all duration-200 whitespace-nowrap',
              selectedCategoryId === category.id
                ? 'cat-pill-active'
                : 'bg-zinc-100 text-zinc-600 hover:bg-zinc-200/70'
            ]"
          >
            {{ category.name }}
            <span :class="['text-xs', selectedCategoryId === category.id ? 'text-white/70' : 'text-zinc-400']">{{ category.product_count || 0 }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="flex gap-6">
        <!-- Sidebar -->
        <div class="hidden lg:block w-60 flex-shrink-0">
          <div class="sticky top-24 space-y-4">
            <div class="bg-white rounded-2xl border border-zinc-100 overflow-hidden shadow-card">
              <div class="px-4 py-3.5 border-b border-zinc-100">
                <h3 class="text-xs font-semibold text-zinc-400 uppercase tracking-wider">商品分类</h3>
              </div>
              <nav class="p-2 space-y-0.5">
                <button
                  @click="selectedCategoryId = null"
                  :class="[
                    'w-full text-left px-3.5 py-2.5 text-sm font-medium rounded-xl transition-all duration-200',
                    selectedCategoryId === null
                      ? 'cat-sidebar-active'
                      : 'text-zinc-600 hover:bg-zinc-50 hover:text-zinc-900'
                  ]"
                >
                  <div class="flex items-center justify-between">
                    <span>全部商品</span>
                    <span :class="['text-xs font-mono', selectedCategoryId === null ? 'text-white/70' : 'text-zinc-400']">{{ totalProducts }}</span>
                  </div>
                </button>
                <button
                  v-for="category in categories"
                  :key="category.id"
                  @click="selectedCategoryId = category.id"
                  :class="[
                    'w-full text-left px-3.5 py-2.5 text-sm font-medium rounded-xl transition-all duration-200',
                    selectedCategoryId === category.id
                      ? 'cat-sidebar-active'
                      : 'text-zinc-600 hover:bg-zinc-50 hover:text-zinc-900'
                  ]"
                >
                  <div class="flex items-center justify-between">
                    <span>{{ category.name }}</span>
                    <span :class="['text-xs font-mono', selectedCategoryId === category.id ? 'text-white/70' : 'text-zinc-400']">{{ category.product_count || 0 }}</span>
                  </div>
                </button>
              </nav>
            </div>

            <div class="bg-white rounded-2xl border border-zinc-100 p-4 shadow-card">
              <div class="space-y-3">
                <div class="flex items-center justify-between">
                  <span class="text-xs text-zinc-400">商品总数</span>
                  <span class="text-sm font-bold text-zinc-900 font-mono">{{ totalProducts }}</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-xs text-zinc-400">分类总数</span>
                  <span class="text-sm font-bold text-zinc-900 font-mono">{{ categories.length }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Products Grid -->
        <div class="flex-1 min-w-0">
          <div v-if="loading" class="flex flex-col items-center justify-center py-24 gap-3">
            <Loader2 class="w-8 h-8 animate-spin text-brand-green" />
            <span class="text-sm text-zinc-400">加载中...</span>
          </div>

          <div v-else-if="filteredProducts.length > 0">
            <div v-if="selectedCategoryId && currentCategory" class="mb-5">
              <h2 class="text-lg font-bold text-zinc-900">{{ currentCategory.name }}</h2>
              <p v-if="currentCategory.description" class="text-sm text-zinc-500 mt-1">{{ currentCategory.description }}</p>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-5 stagger-in">
              <router-link
                v-for="product in filteredProducts"
                :key="product.id"
                :to="`/product/${product.id}`"
                class="product-card group"
              >
                <div v-if="product.image" class="aspect-[16/10] bg-zinc-50 overflow-hidden">
                  <img :src="product.image" :alt="product.name" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
                </div>
                <div v-else class="aspect-[16/10] bg-gradient-to-br from-zinc-50 to-zinc-100 flex items-center justify-center">
                  <div class="w-14 h-14 rounded-2xl bg-brand-gradient-subtle flex items-center justify-center">
                    <Package class="w-7 h-7 text-brand-green/60" />
                  </div>
                </div>

                <div class="p-4 sm:p-5">
                  <div class="flex items-start justify-between gap-2 mb-2">
                    <h3 class="font-semibold text-zinc-900 line-clamp-1 group-hover:text-brand-green transition-colors">{{ product.name }}</h3>
                    <ArrowUpRight class="w-4 h-4 text-zinc-300 group-hover:text-brand-green flex-shrink-0 transition-all duration-300 group-hover:translate-x-0.5 group-hover:-translate-y-0.5" />
                  </div>
                  <p v-if="product.description" class="text-sm text-zinc-500 line-clamp-2 mb-4 leading-relaxed">{{ product.description }}</p>
                  
                  <div class="flex items-end justify-between pt-3 border-t border-zinc-100/80">
                    <div>
                      <div class="text-xl font-bold text-zinc-900 font-mono tracking-tight">{{ formatPrice(product.price) }}</div>
                      <div v-if="product.orig_price && product.orig_price > product.price" class="text-xs text-zinc-400 line-through font-mono mt-0.5">{{ formatPrice(product.orig_price) }}</div>
                    </div>
                    <div :class="['inline-flex items-center gap-1 px-2 py-1 rounded-lg text-xs font-medium', product.stock_count > 10 ? 'bg-emerald-50 text-emerald-600' : product.stock_count > 0 ? 'bg-amber-50 text-amber-600' : 'bg-red-50 text-red-500']">
                      <span class="w-1.5 h-1.5 rounded-full" :class="[product.stock_count > 10 ? 'bg-emerald-500' : product.stock_count > 0 ? 'bg-amber-500' : 'bg-red-500']"></span>
                      {{ product.stock_count > 0 ? '库存 ' + product.stock_count : '已售罄' }}
                    </div>
                  </div>
                </div>
              </router-link>
            </div>
          </div>

          <div v-else class="flex flex-col items-center justify-center py-20 px-4">
            <div class="w-20 h-20 rounded-3xl bg-brand-gradient-subtle flex items-center justify-center mb-5">
              <Package class="w-10 h-10 text-brand-green/50" />
            </div>
            <h3 class="text-lg font-semibold text-zinc-900 mb-1.5">暂无商品</h3>
            <p class="text-sm text-zinc-500">{{ selectedCategoryId ? '该分类下暂时没有可售商品' : '暂时没有可售商品' }}</p>
            <button v-if="selectedCategoryId" @click="selectedCategoryId = null" class="mt-5 inline-flex items-center gap-1.5 px-4 py-2 text-sm font-medium text-brand-green hover:text-brand-green-light transition-colors">
              查看全部商品
              <ArrowUpRight class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Package, Loader2, ArrowUpRight, Grid3X3 } from 'lucide-vue-next'
import api from '@/utils/api'
import { formatPrice } from '@/utils/helpers'

const settings = ref({})
const categories = ref([])
const products = ref([])
const loading = ref(true)
const selectedCategoryId = ref(null)

const totalProducts = computed(() => products.value.length)

const currentCategory = computed(() => {
  if (!selectedCategoryId.value) return null
  return categories.value.find(c => c.id === selectedCategoryId.value)
})

const filteredProducts = computed(() => {
  if (selectedCategoryId.value === null) return products.value
  return products.value.filter(p => p.category_id === selectedCategoryId.value)
})

onMounted(async () => {
  try {
    const settingsRes = await api.get('/api/settings')
    settings.value = settingsRes.data

    const categoriesRes = await api.get('/api/categories/with-products')
    const categoriesData = categoriesRes.data || []
    
    categories.value = categoriesData.map(cat => ({
      id: cat.id,
      name: cat.name,
      description: cat.description,
      product_count: cat.products?.length || 0
    }))
    
    products.value = categoriesData.flatMap(cat => 
      (cat.products || []).map(p => ({ ...p, category_id: cat.id, category_name: cat.name }))
    )
  } catch (error) {
    console.error('Failed to load data', error)
  } finally {
    loading.value = false
  }
})
</script>
