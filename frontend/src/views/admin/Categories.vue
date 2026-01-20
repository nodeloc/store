<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <div>
        <h2 class="text-xl font-semibold text-zinc-900">商品分类</h2>
        <p class="text-sm text-zinc-600 mt-1">管理商品分类</p>
      </div>
      <button @click="showCreateModal = true" class="inline-flex items-center space-x-2 px-4 py-2 bg-zinc-900 text-white text-sm font-medium rounded-lg hover:bg-zinc-800 transition">
        <Plus class="w-4 h-4" />
        <span>新建分类</span>
      </button>
    </div>
    
    <!-- Categories List -->
    <div class="bg-white rounded-lg border border-zinc-100">
      <div v-if="loading" class="p-12 text-center">
        <Loader2 class="w-8 h-8 animate-spin text-zinc-400 mx-auto" />
      </div>
      
      <div v-else-if="categories.length === 0" class="p-12 text-center">
        <Folder class="w-12 h-12 text-zinc-300 mx-auto mb-4" />
        <p class="text-zinc-600">暂无分类</p>
      </div>
      
      <table v-else class="w-full">
        <thead class="bg-zinc-50 border-b border-zinc-100">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">名称</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">描述</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">排序</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">状态</th>
            <th class="px-6 py-3 text-right text-xs font-medium text-zinc-500 uppercase">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-zinc-100">
          <tr v-for="category in categories" :key="category.id" class="hover:bg-zinc-50">
            <td class="px-6 py-4 text-sm text-zinc-900 font-medium">{{ category.name }}</td>
            <td class="px-6 py-4 text-sm text-zinc-600">{{ category.description || '-' }}</td>
            <td class="px-6 py-4 text-sm text-zinc-600">{{ category.sort }}</td>
            <td class="px-6 py-4">
              <span :class="category.is_active ? 'bg-green-100 text-green-800' : 'bg-zinc-100 text-zinc-800'" class="px-2 py-1 text-xs font-medium rounded">
                {{ category.is_active ? '启用' : '禁用' }}
              </span>
            </td>
            <td class="px-6 py-4 text-right space-x-2">
              <button @click="editCategory(category)" class="text-zinc-600 hover:text-zinc-900">
                <Edit2 class="w-4 h-4" />
              </button>
              <button @click="deleteCategory(category.id)" class="text-red-600 hover:text-red-900">
                <Trash2 class="w-4 h-4" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal || editingCategory" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="closeModal">
      <div class="bg-white rounded-lg w-full max-w-md p-6">
        <h3 class="text-lg font-semibold text-zinc-900 mb-4">{{ editingCategory ? '编辑分类' : '新建分类' }}</h3>
        
        <form @submit.prevent="saveCategory" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1">分类名称</label>
            <input v-model="form.name" type="text" required class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900">
          </div>
          
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1">描述</label>
            <textarea v-model="form.description" rows="3" class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900"></textarea>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1">排序</label>
            <input v-model.number="form.sort" type="number" class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900">
          </div>
          
          <div class="flex items-center">
            <input v-model="form.is_active" type="checkbox" class="w-4 h-4 text-zinc-900 border-zinc-300 rounded focus:ring-zinc-900">
            <label class="ml-2 text-sm text-zinc-700">启用</label>
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
import { ref, onMounted } from 'vue'
import { Plus, Edit2, Trash2, Folder, Loader2 } from 'lucide-vue-next'
import api from '@/utils/api'
import { useToast } from '@/stores/toast'

const toast = useToast()
const categories = ref([])
const loading = ref(false)
const showCreateModal = ref(false)
const editingCategory = ref(null)
const form = ref({
  name: '',
  description: '',
  sort: 0,
  is_active: true
})

onMounted(() => {
  fetchCategories()
})

async function fetchCategories() {
  try {
    loading.value = true
    const response = await api.get('/api/admin/categories')
    categories.value = response.data.categories || []
  } catch (error) {
    toast.error('加载分类失败')
  } finally {
    loading.value = false
  }
}

function editCategory(category) {
  editingCategory.value = category
  form.value = { ...category }
}

function closeModal() {
  showCreateModal.value = false
  editingCategory.value = null
  form.value = {
    name: '',
    description: '',
    sort: 0,
    is_active: true
  }
}

async function saveCategory() {
  try {
    if (editingCategory.value) {
      await api.put(`/api/admin/categories/${editingCategory.value.id}`, form.value)
      toast.success('更新成功')
    } else {
      await api.post('/api/admin/categories', form.value)
      toast.success('创建成功')
    }
    closeModal()
    fetchCategories()
  } catch (error) {
    toast.error('保存失败')
  }
}

async function deleteCategory(id) {
  if (!confirm('确定要删除此分类吗？')) return
  
  try {
    await api.delete(`/api/admin/categories/${id}`)
    toast.success('删除成功')
    fetchCategories()
  } catch (error) {
    toast.error('删除失败')
  }
}
</script>
