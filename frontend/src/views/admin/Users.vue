<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-xl font-semibold text-zinc-900">用户管理</h2>
      <p class="text-sm text-zinc-600 mt-1">管理所有用户</p>
    </div>
    
    <div class="bg-white rounded-lg border border-zinc-100">
      <div v-if="loading" class="p-12 text-center">
        <Loader2 class="w-8 h-8 animate-spin text-zinc-400 mx-auto" />
      </div>
      
      <table v-else class="w-full text-sm">
        <thead class="bg-zinc-50 border-b border-zinc-100">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">用户</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">邮箱</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">角色</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">状态</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-zinc-500 uppercase">注册时间</th>
            <th class="px-6 py-3 text-right text-xs font-medium text-zinc-500 uppercase">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-zinc-100">
          <tr v-for="user in users" :key="user.id" class="hover:bg-zinc-50">
            <td class="px-6 py-4">
              <div class="flex items-center space-x-3">
                <img :src="user.avatar_url" alt="" class="w-8 h-8 rounded-full">
                <span class="text-zinc-900">{{ user.name }}</span>
              </div>
            </td>
            <td class="px-6 py-4 text-zinc-600">{{ user.email || '-' }}</td>
            <td class="px-6 py-4">
              <span v-if="user.is_admin" class="px-2 py-1 bg-purple-100 text-purple-800 text-xs font-medium rounded">管理员</span>
              <span v-else class="px-2 py-1 bg-zinc-100 text-zinc-800 text-xs font-medium rounded">普通用户</span>
            </td>
            <td class="px-6 py-4">
              <span :class="user.is_blocked ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'" class="px-2 py-1 text-xs font-medium rounded">
                {{ user.is_blocked ? '已封禁' : '正常' }}
              </span>
            </td>
            <td class="px-6 py-4 text-zinc-600">{{ formatDate(user.created_at) }}</td>
            <td class="px-6 py-4 text-right">
              <button @click="toggleAdmin(user)" class="text-zinc-600 hover:text-zinc-900 text-xs mr-2">
                {{ user.is_admin ? '取消管理员' : '设为管理员' }}
              </button>
              <button @click="toggleBlock(user)" :class="user.is_blocked ? 'text-green-600 hover:text-green-900' : 'text-red-600 hover:text-red-900'" class="text-xs">
                {{ user.is_blocked ? '解封' : '封禁' }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Loader2 } from 'lucide-vue-next'
import api from '@/utils/api'
import { useToast } from '@/stores/toast'

const toast = useToast()
const users = ref([])
const loading = ref(false)

onMounted(() => {
  fetchUsers()
})

async function fetchUsers() {
  try {
    loading.value = true
    const response = await api.get('/api/admin/users', { params: { page: 1, page_size: 50 } })
    users.value = response.data.users || []
  } catch (error) {
    toast.error('加载用户失败')
  } finally {
    loading.value = false
  }
}

async function toggleAdmin(user) {
  try {
    await api.put(`/api/admin/users/${user.id}`, { is_admin: !user.is_admin })
    toast.success('更新成功')
    fetchUsers()
  } catch (error) {
    toast.error('更新失败')
  }
}

async function toggleBlock(user) {
  try {
    await api.put(`/api/admin/users/${user.id}`, { is_blocked: !user.is_blocked })
    toast.success('更新成功')
    fetchUsers()
  } catch (error) {
    toast.error('更新失败')
  }
}

function formatDate(date) {
  return new Date(date).toLocaleDateString('zh-CN')
}
</script>
