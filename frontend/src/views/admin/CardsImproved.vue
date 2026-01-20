<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h2 class="text-xl font-semibold text-zinc-900">卡密管理</h2>
        <p class="text-sm text-zinc-600 mt-1">管理商品卡密库存</p>
      </div>
      <div class="flex space-x-3">
        <button 
          v-if="selectedCards.length > 0"
          @click="exportSelected"
          class="inline-flex items-center space-x-2 px-4 py-2 bg-zinc-100 text-zinc-900 text-sm font-medium rounded-lg hover:bg-zinc-200 transition"
        >
          <Download class="w-4 h-4" />
          <span>导出选中 ({{ selectedCards.length }})</span>
        </button>
        <button 
          v-if="selectedCards.length > 0"
          @click="deleteSelected"
          class="inline-flex items-center space-x-2 px-4 py-2 bg-red-100 text-red-900 text-sm font-medium rounded-lg hover:bg-red-200 transition"
        >
          <Trash2 class="w-4 h-4" />
          <span>删除选中</span>
        </button>
        <button @click="showAddModal = true" class="inline-flex items-center space-x-2 px-4 py-2 bg-zinc-900 text-white text-sm font-medium rounded-lg hover:bg-zinc-800 transition">
          <Plus class="w-4 h-4" />
          <span>批量添加卡密</span>
        </button>
      </div>
    </div>
    
    <!-- Product Filter -->
    <div class="bg-white rounded-lg border border-zinc-100 p-4">
      <label class="block text-sm font-medium text-zinc-700 mb-2">选择商品</label>
      <select v-model="selectedProductId" @change="fetchCards" class="w-full max-w-md px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900">
        <option value="">全部商品</option>
        <option v-for="product in products" :key="product.id" :value="product.id">{{ product.name }}</option>
      </select>
    </div>
    
    <!-- Cards List -->
    <div class="bg-white rounded-lg border border-zinc-100 overflow-hidden">
      <div v-if="loading" class="p-12 text-center">
        <Loader2 class="w-8 h-8 animate-spin text-zinc-400 mx-auto" />
      </div>
      
      <div v-else-if="cards.length === 0" class="p-12 text-center">
        <CreditCard class="w-12 h-12 text-zinc-300 mx-auto mb-4" />
        <p class="text-zinc-600">暂无卡密</p>
      </div>
      
      <div v-else>
      <table class="w-full font-mono text-sm">
        <thead class="bg-zinc-50 border-b border-zinc-100">
          <tr>
            <th class="px-6 py-3 text-left">
              <input 
                type="checkbox" 
                @change="toggleSelectAll"
                :checked="selectedCards.length === availableCards.length && availableCards.length > 0"
                class="w-4 h-4 text-zinc-900 border-zinc-300 rounded focus:ring-zinc-900"
              >
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">卡号</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">密码</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">商品</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">状态</th>
            <th class="px-6 py-3 text-right text-xs font-medium text-zinc-500 uppercase">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-zinc-100">
          <tr v-for="card in cards" :key="card.id" class="hover:bg-zinc-50 transition">
            <td class="px-6 py-4">
              <input 
                v-if="card.status === 0"
                type="checkbox"
                :value="card.id"
                v-model="selectedCards"
                class="w-4 h-4 text-zinc-900 border-zinc-300 rounded focus:ring-zinc-900"
              >
            </td>
            <td class="px-6 py-4">
              <CopyText :text="card.card_no" />
            </td>
            <td class="px-6 py-4 text-zinc-600">
              <CopyText v-if="card.card_pwd" :text="card.card_pwd" />
              <span v-else class="text-zinc-400">-</span>
            </td>
            <td class="px-6 py-4 text-zinc-600">{{ card.product?.name || '-' }}</td>
            <td class="px-6 py-4">
              <span :class="getStatusClass(card.status)" class="px-2 py-1 text-xs font-medium rounded">
                {{ getStatusText(card.status) }}
              </span>
            </td>
            <td class="px-6 py-4 text-right">
              <button v-if="card.status === 0" @click="deleteCard(card.id)" class="text-red-600 hover:text-red-900 transition">
                <Trash2 class="w-4 h-4" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <!-- Pagination -->
      <Pagination
        :current-page="currentPage"
        :total-pages="totalPages"
        :total="total"
        :page-size="pageSize"
        @change="handlePageChange"
      />
      </div>
    </div>
    
    <!-- Add Modal -->
    <div v-if="showAddModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="closeAddModal">
      <div class="bg-white rounded-lg w-full max-w-2xl p-6">
        <h3 class="text-lg font-semibold text-zinc-900 mb-4">批量添加卡密</h3>
        
        <form @submit.prevent="addCards" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1">选择商品*</label>
            <select v-model="addForm.product_id" required class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900">
              <option value="">请选择</option>
              <option v-for="product in products" :key="product.id" :value="product.id">{{ product.name }}</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1">卡密列表*</label>
            <p class="text-xs text-zinc-500 mb-2">每行一个卡密，格式：<code class="bg-zinc-100 px-1 py-0.5 rounded font-mono">卡号----密码</code> 或只填卡号</p>
            <textarea v-model="addForm.cards_text" rows="10" required placeholder="例如：&#10;ABC123----XYZ789&#10;DEF456----UVW012&#10;GHI789" class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900 font-mono text-sm"></textarea>
          </div>
          
          <div class="flex justify-end space-x-3 pt-4">
            <button type="button" @click="closeAddModal" class="px-4 py-2 text-sm font-medium text-zinc-700 hover:text-zinc-900">取消</button>
            <button type="submit" :disabled="saving" class="px-4 py-2 bg-zinc-900 text-white text-sm font-medium rounded-lg hover:bg-zinc-800 disabled:opacity-50 inline-flex items-center space-x-2">
              <Loader2 v-if="saving" class="w-4 h-4 animate-spin" />
              <span>{{ saving ? '添加中...' : '添加' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Plus, Trash2, CreditCard, Loader2, Download } from 'lucide-vue-next'
import CopyText from '@/components/CopyText.vue'
import Pagination from '@/components/Pagination.vue'
import api from '@/utils/api'
import { useToast } from '@/stores/toast'

const route = useRoute()
const toast = useToast()
const cards = ref([])
const products = ref([])
const loading = ref(false)
const saving = ref(false)
const showAddModal = ref(false)
const selectedProductId = ref(route.query.product_id || '')
const selectedCards = ref([])

const addForm = ref({
  product_id: '',
  cards_text: ''
})

// 分页相关
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const availableCards = computed(() => cards.value.filter(c => c.status === 0))

onMounted(() => {
  fetchProducts()
  fetchCards() // 无论是否选择商品，都加载卡密列表
})

async function fetchProducts() {
  try {
    const response = await api.get('/api/admin/products', { params: { page: 1, page_size: 999 } })
    products.value = response.data.products || []
  } catch (error) {
    toast.error('加载商品失败')
  }
}

async function fetchCards() {
  try {
    loading.value = true
    selectedCards.value = []
    const params = { page: currentPage.value, page_size: pageSize.value }
    if (selectedProductId.value) {
      params.product_id = selectedProductId.value
    }
    const response = await api.get('/api/admin/card-keys', { params })
    cards.value = response.data.card_keys || []
    total.value = response.data.total || 0
  } catch (error) {
    toast.error('加载卡密失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange(page) {
  currentPage.value = page
  fetchCards()
}

function toggleSelectAll(event) {
  if (event.target.checked) {
    selectedCards.value = availableCards.value.map(c => c.id)
  } else {
    selectedCards.value = []
  }
}

function closeAddModal() {
  showAddModal.value = false
  addForm.value = {
    product_id: '',
    cards_text: ''
  }
}

async function addCards() {
  try {
    saving.value = true
    await api.post('/api/admin/card-keys', addForm.value)
    toast.success('添加成功')
    closeAddModal()
    fetchCards()
  } catch (error) {
    toast.error('添加失败')
  } finally {
    saving.value = false
  }
}

async function deleteCard(id) {
  if (!confirm('确定要删除此卡密吗？')) return
  
  try {
    await api.delete(`/api/admin/card-keys/${id}`)
    toast.success('删除成功')
    fetchCards()
  } catch (error) {
    toast.error('删除失败')
  }
}

async function deleteSelected() {
  if (!confirm(`确定要删除选中的 ${selectedCards.value.length} 个卡密吗？`)) return
  
  try {
    await Promise.all(selectedCards.value.map(id => api.delete(`/api/admin/card-keys/${id}`)))
    toast.success('批量删除成功')
    selectedCards.value = []
    fetchCards()
  } catch (error) {
    toast.error('批量删除失败')
  }
}

function exportSelected() {
  const selectedData = cards.value.filter(c => selectedCards.value.includes(c.id))
  exportCards(selectedData, 'json')
}

function exportCards(data, format) {
  const cardsToExport = data || cards.value
  
  if (format === 'json') {
    const json = JSON.stringify(cardsToExport, null, 2)
    downloadFile(json, `cards_${Date.now()}.json`, 'application/json')
  } else if (format === 'csv') {
    const csv = convertToCSV(cardsToExport)
    downloadFile(csv, `cards_${Date.now()}.csv`, 'text/csv')
  }
  
  toast.success('导出成功')
}

function convertToCSV(data) {
  const headers = ['ID', '卡号', '密码', '商品', '状态']
  const rows = data.map(card => [
    card.id,
    card.card_no,
    card.card_pwd || '',
    card.product?.name || '',
    getStatusText(card.status)
  ])
  
  return [headers, ...rows].map(row => row.map(cell => `"${cell}"`).join(',')).join('\n')
}

function downloadFile(content, filename, mimeType) {
  const blob = new Blob([content], { type: mimeType })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

function getStatusText(status) {
  const statusMap = { 0: '可售', 1: '已售出', 2: '已锁定' }
  return statusMap[status] || '未知'
}

function getStatusClass(status) {
  const classMap = {
    0: 'bg-green-100 text-green-800',
    1: 'bg-zinc-100 text-zinc-800',
    2: 'bg-yellow-100 text-yellow-800'
  }
  return classMap[status] || 'bg-zinc-100 text-zinc-800'
}
</script>
