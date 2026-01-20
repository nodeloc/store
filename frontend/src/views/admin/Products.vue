<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h2 class="text-xl font-semibold text-zinc-900">商品管理</h2>
        <p class="text-sm text-zinc-600 mt-1">管理所有商品</p>
      </div>
      <button @click="showCreateModal = true" class="inline-flex items-center space-x-2 px-4 py-2 bg-zinc-900 text-white text-sm font-medium rounded-lg hover:bg-zinc-800 transition">
        <Plus class="w-4 h-4" />
        <span>新建商品</span>
      </button>
    </div>
    
    <!-- Products List -->
    <div class="bg-white rounded-lg border border-zinc-100">
      <div v-if="loading" class="p-12 text-center">
        <Loader2 class="w-8 h-8 animate-spin text-zinc-400 mx-auto" />
      </div>
      
      <div v-else-if="products.length === 0" class="p-12 text-center">
        <Package class="w-12 h-12 text-zinc-300 mx-auto mb-4" />
        <p class="text-zinc-600">暂无商品</p>
      </div>
      
      <table v-else class="w-full">
        <thead class="bg-zinc-50 border-b border-zinc-100">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">商品名称</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">价格</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">库存</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">销量</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">状态</th>
            <th class="px-6 py-3 text-right text-xs font-medium text-zinc-500 uppercase">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-zinc-100">
          <tr v-for="product in products" :key="product.id" class="hover:bg-zinc-50">
            <td class="px-6 py-4">
              <div class="flex items-center space-x-3">
                <img v-if="product.image" :src="product.image" alt="" class="w-10 h-10 rounded object-cover">
                <div class="w-10 h-10 bg-zinc-100 rounded flex items-center justify-center" v-else>
                  <Package class="w-5 h-5 text-zinc-400" />
                </div>
                <span class="text-sm text-zinc-900 font-medium">{{ product.name }}</span>
              </div>
            </td>
            <td class="px-6 py-4 text-sm text-zinc-900">¥{{ product.price }}</td>
            <td class="px-6 py-4 text-sm text-zinc-600">{{ product.stock_count }}</td>
            <td class="px-6 py-4 text-sm text-zinc-600">{{ product.sales_count }}</td>
            <td class="px-6 py-4">
              <span :class="product.is_active ? 'bg-green-100 text-green-800' : 'bg-zinc-100 text-zinc-800'" class="px-2 py-1 text-xs font-medium rounded">
                {{ product.is_active ? '上架' : '下架' }}
              </span>
            </td>
            <td class="px-6 py-4 text-right space-x-2">
              <button @click="$router.push(`/admin/cards?product_id=${product.id}`)" class="text-zinc-600 hover:text-zinc-900" title="管理卡密">
                <CreditCard class="w-4 h-4" />
              </button>
              <button @click="editProduct(product)" class="text-zinc-600 hover:text-zinc-900">
                <Edit2 class="w-4 h-4" />
              </button>
              <button @click="deleteProduct(product.id)" class="text-red-600 hover:text-red-900">
                <Trash2 class="w-4 h-4" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <!-- Pagination -->
    <div v-if="total > pageSize" class="flex justify-center">
      <div class="flex space-x-2">
        <button v-for="p in totalPages" :key="p" @click="page = p" :class="p === page ? 'bg-zinc-900 text-white' : 'bg-white text-zinc-700 hover:bg-zinc-100'" class="px-4 py-2 text-sm font-medium rounded-lg border border-zinc-200">
          {{ p }}
        </button>
      </div>
    </div>
    
    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal || editingProduct" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 overflow-y-auto" @click.self="closeModal">
      <div class="bg-white rounded-lg w-full max-w-2xl p-6 my-8">
        <h3 class="text-lg font-semibold text-zinc-900 mb-4">{{ editingProduct ? '编辑商品' : '新建商品' }}</h3>
        
        <form @submit.prevent="saveProduct" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-zinc-700 mb-1">商品名称*</label>
              <input v-model="form.name" type="text" required class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900">
            </div>
            
            <div>
              <label class="block text-sm font-medium text-zinc-700 mb-1">分类*</label>
              <select v-model="form.category_id" required class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900">
                <option value="">请选择</option>
                <option v-for="cat in allCategories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1">商品描述</label>
            <textarea v-model="form.description" rows="3" class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900"></textarea>
          </div>
          
          <div class="grid grid-cols-3 gap-4">
            <div>
              <label class="block text-sm font-medium text-zinc-700 mb-1">售价*</label>
              <input v-model.number="form.price" type="number" step="0.01" required class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900">
            </div>
            
            <div>
              <label class="block text-sm font-medium text-zinc-700 mb-1">原价</label>
              <input v-model.number="form.orig_price" type="number" step="0.01" class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900">
            </div>
            
            <div>
              <label class="block text-sm font-medium text-zinc-700 mb-1">排序</label>
              <input v-model.number="form.sort" type="number" class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900">
            </div>
          </div>
          
          <div>
            <ImageUpload v-model="form.image" label="商品图片" />
          </div>
          
          <div class="flex items-center">
            <input v-model="form.is_active" type="checkbox" class="w-4 h-4 text-zinc-900 border-zinc-300 rounded focus:ring-zinc-900">
            <label class="ml-2 text-sm text-zinc-700">上架</label>
          </div>
          
          <div class="flex justify-end space-x-3 pt-4">
            <button type="button" @click="closeModal" class="px-4 py-2 text-sm font-medium text-zinc-700 hover:text-zinc-900">取消</button>
            <button type="submit" class="px-4 py-2 bg-zinc-900 text-white text-sm font-medium rounded-lg hover:bg-zinc-800">保存</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { Plus, Edit2, Trash2, Package, CreditCard, Loader2 } from 'lucide-vue-next'
import ImageUpload from '@/components/ImageUpload.vue'
import api from '@/utils/api'
import { useToast } from '@/stores/toast'

const toast = useToast()
const products = ref([])
const allCategories = ref([])
const loading = ref(false)
const showCreateModal = ref(false)
const editingProduct = ref(null)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const form = ref({
  name: '',
  category_id: '',
  description: '',
  price: 0,
  orig_price: 0,
  image: '',
  sort: 0,
  is_active: true
})

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

watch(page, () => {
  fetchProducts()
})

onMounted(() => {
  fetchProducts()
  fetchCategories()
})

async function fetchProducts() {
  try {
    loading.value = true
    const response = await api.get('/api/admin/products', {
      params: { page: page.value, page_size: pageSize.value }
    })
    products.value = response.data.products || []
    total.value = response.data.total || 0
  } catch (error) {
    toast.error('加载商品失败')
  } finally {
    loading.value = false
  }
}

async function fetchCategories() {
  try {
    const response = await api.get('/api/admin/categories')
    allCategories.value = response.data.categories || []
  } catch (error) {
    console.error('加载分类失败', error)
  }
}

function editProduct(product) {
  editingProduct.value = product
  form.value = { ...product }
}

function closeModal() {
  showCreateModal.value = false
  editingProduct.value = null
  form.value = {
    name: '',
    category_id: '',
    description: '',
    price: 0,
    orig_price: 0,
    image: '',
    sort: 0,
    is_active: true
  }
}

async function saveProduct() {
  try {
    if (editingProduct.value) {
      await api.put(`/api/admin/products/${editingProduct.value.id}`, form.value)
      toast.success('更新成功')
    } else {
      await api.post('/api/admin/products', form.value)
      toast.success('创建成功')
    }
    closeModal()
    fetchProducts()
  } catch (error) {
    toast.error('保存失败')
  }
}

async function deleteProduct(id) {
  if (!confirm('确定要删除此商品吗？')) return
  
  try {
    await api.delete(`/api/admin/products/${id}`)
    toast.success('删除成功')
    fetchProducts()
  } catch (error) {
    toast.error('删除失败')
  }
}
</script>
