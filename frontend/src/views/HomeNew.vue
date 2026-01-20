<template>
  <div class="min-h-screen bg-zinc-50 flex flex-col">
    <!-- Header -->
    <div class="bg-white border-b border-zinc-200 sticky top-0 z-10 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center space-x-8">
            <router-link to="/" class="flex items-center space-x-2">
              <div class="w-8 h-8 bg-zinc-900 rounded-lg flex items-center justify-center">
                <span class="text-white font-bold text-sm">NL</span>
              </div>
              <span class="text-lg font-bold text-zinc-900">{{ settings.site_name || 'NodeLoc 社区发卡' }}</span>
            </router-link>
          </div>
          <div class="flex items-center space-x-6">
            <template v-if="authStore.isAuthenticated">
              <router-link to="/orders" class="text-sm font-medium text-zinc-600 hover:text-zinc-900 transition">我的订单</router-link>
              <router-link v-if="authStore.isAdmin" to="/admin" class="text-sm font-medium text-zinc-600 hover:text-zinc-900 transition">管理后台</router-link>
              <button @click="handleLogout" class="text-sm font-medium text-zinc-600 hover:text-zinc-900 transition">退出</button>
            </template>
            <router-link v-else to="/login" class="inline-flex items-center px-4 py-2 bg-zinc-900 text-white text-sm font-medium rounded-lg hover:bg-zinc-800 transition">
              登录
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 w-full">
      <div class="flex gap-6">
        <!-- Sidebar - Categories -->
        <div class="w-64 flex-shrink-0">
          <div class="bg-white rounded-lg border border-zinc-200 overflow-hidden sticky top-24 shadow-sm">
            <div class="px-4 py-3 bg-zinc-50 border-b border-zinc-200">
              <h3 class="text-sm font-semibold text-zinc-900 uppercase tracking-wide">商品分类</h3>
            </div>
            <nav class="p-2">
              <button
                @click="selectedCategoryId = null"
                :class="[
                  'w-full text-left px-4 py-2.5 text-sm font-medium rounded-md transition-colors',
                  selectedCategoryId === null
                    ? 'bg-zinc-900 text-white'
                    : 'text-zinc-700 hover:bg-zinc-100'
                ]"
              >
                <div class="flex items-center justify-between">
                  <span>全部商品</span>
                  <span class="text-xs">{{ totalProducts }}</span>
                </div>
              </button>
              
              <button
                v-for="category in categories"
                :key="category.id"
                @click="selectedCategoryId = category.id"
                :class="[
                  'w-full text-left px-4 py-2.5 text-sm font-medium rounded-md transition-colors mt-1',
                  selectedCategoryId === category.id
                    ? 'bg-zinc-900 text-white'
                    : 'text-zinc-700 hover:bg-zinc-100'
                ]"
              >
                <div class="flex items-center justify-between">
                  <span>{{ category.name }}</span>
                  <span class="text-xs">{{ category.product_count || 0 }}</span>
                </div>
              </button>
            </nav>
          </div>
        </div>

        <!-- Products Grid -->
        <div class="flex-1 min-w-0">
          <!-- Loading -->
          <div v-if="loading" class="flex justify-center items-center py-20">
            <Loader2 class="w-8 h-8 animate-spin text-zinc-400" />
          </div>

          <!-- Products -->
          <div v-else-if="filteredProducts.length > 0">
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <router-link
                v-for="product in filteredProducts"
                :key="product.id"
                :to="`/product/${product.id}`"
                class="group bg-white border border-zinc-200 rounded-lg hover:border-zinc-900 hover:shadow-lg transition-all duration-200"
              >
                <!-- Product Image -->
                <div v-if="product.image" class="aspect-video bg-zinc-100 rounded-t-lg overflow-hidden">
                  <img :src="product.image" :alt="product.name" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-200">
                </div>
                <div v-else class="aspect-video bg-gradient-to-br from-zinc-100 to-zinc-200 rounded-t-lg flex items-center justify-center">
                  <Package class="w-12 h-12 text-zinc-400" />
                </div>

                <!-- Product Info -->
                <div class="p-4">
                  <h3 class="font-semibold text-zinc-900 mb-2 line-clamp-1">{{ product.name }}</h3>
                  <p v-if="product.description" class="text-sm text-zinc-600 line-clamp-2 mb-3">{{ product.description }}</p>
                  
                  <!-- Price & Stock -->
                  <div class="flex items-center justify-between pt-3 border-t border-zinc-100">
                    <div>
                      <div class="text-lg font-bold text-zinc-900 font-mono">{{ formatPrice(product.price) }}</div>
                      <div v-if="product.orig_price && product.orig_price > product.price" class="text-xs text-zinc-500 line-through font-mono">
                        {{ formatPrice(product.orig_price) }}
                      </div>
                    </div>
                    <div class="text-right">
                      <div class="text-xs text-zinc-500">库存</div>
                      <div :class="[
                        'text-sm font-semibold font-mono',
                        product.stock_count > 10 ? 'text-green-600' : product.stock_count > 0 ? 'text-orange-600' : 'text-red-600'
                      ]">
                        {{ product.stock_count > 0 ? product.stock_count : '已售罄' }}
                      </div>
                    </div>
                  </div>
                </div>
              </router-link>
            </div>
          </div>

          <!-- Empty State -->
          <div v-else class="bg-white rounded-lg border border-zinc-200 p-12 text-center">
            <Package class="w-16 h-16 text-zinc-300 mx-auto mb-4" />
            <h3 class="text-lg font-semibold text-zinc-900 mb-2">暂无商品</h3>
            <p class="text-sm text-zinc-600">{{ selectedCategoryId ? '该分类下' : '' }}暂时没有可售商品</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <footer class="bg-white border-t border-zinc-200 mt-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="text-center text-sm text-zinc-600">
          <p>{{ settings.site_footer || '© 2026 NodeLoc 社区发卡系统' }}</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Package, Loader2 } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'
import { formatPrice } from '@/utils/helpers'

const router = useRouter()
const authStore = useAuthStore()

const settings = ref({})
const categories = ref([])
const products = ref([])
const loading = ref(true)
const selectedCategoryId = ref(null)

const totalProducts = computed(() => products.value.length)

const filteredProducts = computed(() => {
  if (selectedCategoryId.value === null) {
    return products.value
  }
  return products.value.filter(p => p.category_id === selectedCategoryId.value)
})

onMounted(async () => {
  try {
    // Fetch settings
    const settingsRes = await api.get('/api/settings')
    settings.value = settingsRes.data
    
    // Fetch categories with products
    const categoriesRes = await api.get('/api/categories/with-products')
    const categoriesData = categoriesRes.data || []
    
    // Extract categories and products
    categories.value = categoriesData.map(cat => ({
      id: cat.id,
      name: cat.name,
      description: cat.description,
      product_count: cat.products?.length || 0
    }))
    
    // Flatten all products
    products.value = categoriesData.flatMap(cat => 
      (cat.products || []).map(p => ({
        ...p,
        category_id: cat.id,
        category_name: cat.name
      }))
    )
  } catch (error) {
    console.error('Failed to load data', error)
  } finally {
    loading.value = false
  }
})

async function handleLogout() {
  try {
    await authStore.logout()
    router.push('/')
  } catch (error) {
    console.error('Logout failed', error)
  }
}
</script>
