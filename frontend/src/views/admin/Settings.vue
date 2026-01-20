<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-xl font-semibold text-zinc-900">系统设置</h2>
      <p class="text-sm text-zinc-600 mt-1">配置网站基本信息</p>
    </div>
    
    <form @submit.prevent="saveSettings" class="space-y-6">
      <!-- Site Settings -->
      <div class="bg-white rounded-lg border border-zinc-100 p-6">
        <h3 class="text-lg font-semibold text-zinc-900 mb-4">网站信息</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1">网站名称</label>
            <input v-model="settings.site_name" type="text" class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900">
          </div>
          
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1">网站描述</label>
            <textarea v-model="settings.site_description" rows="3" class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900"></textarea>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1">页脚文本</label>
            <input v-model="settings.footer_text" type="text" class="w-full px-3 py-2 border border-zinc-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-zinc-900">
          </div>
        </div>
      </div>
      
      <!-- Info Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div class="flex items-start space-x-3">
            <Info class="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" />
            <div class="space-y-1">
              <h4 class="text-sm font-medium text-blue-900">OAuth 登录</h4>
              <p class="text-xs text-blue-700">系统默认使用 NodeLoc OAuth 2.0 登录</p>
            </div>
          </div>
        </div>
        
        <div class="bg-green-50 border border-green-200 rounded-lg p-4">
          <div class="flex items-start space-x-3">
            <CreditCard class="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" />
            <div class="space-y-1">
              <h4 class="text-sm font-medium text-green-900">社区支付</h4>
              <p class="text-xs text-green-700">系统默认使用 NodeLoc 社区积分支付</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Configuration Hint -->
      <div class="bg-zinc-50 border border-zinc-200 rounded-lg p-6">
        <h3 class="text-sm font-semibold text-zinc-900 mb-3">配置说明</h3>
        <div class="space-y-2 text-sm text-zinc-600">
          <p>• OAuth 和支付参数请在 <code class="bg-white px-2 py-0.5 rounded border border-zinc-300 font-mono text-xs">.env</code> 文件中配置</p>
          <p>• 必需配置项：<code class="font-mono text-xs">NODELOC_CLIENT_ID</code>、<code class="font-mono text-xs">NODELOC_CLIENT_SECRET</code>、<code class="font-mono text-xs">PAYMENT_ID</code>、<code class="font-mono text-xs">PAYMENT_SECRET</code></p>
          <p>• 配置完成后需要重启服务才能生效</p>
        </div>
      </div>
      
      <div class="flex justify-end">
        <button type="submit" :disabled="saving" class="px-6 py-2 bg-zinc-900 text-white text-sm font-medium rounded-lg hover:bg-zinc-800 disabled:opacity-50 inline-flex items-center space-x-2">
          <Loader2 v-if="saving" class="w-4 h-4 animate-spin" />
          <span>{{ saving ? '保存中...' : '保存设置' }}</span>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Info, CreditCard, Loader2 } from 'lucide-vue-next'
import api from '@/utils/api'
import { useToast } from '@/stores/toast'

const toast = useToast()
const saving = ref(false)
const settings = ref({
  site_name: '',
  site_description: '',
  footer_text: ''
})

onMounted(async () => {
  try {
    const response = await api.get('/api/admin/settings')
    const data = response.data.settings
    settings.value = {
      site_name: data.site_name || '',
      site_description: data.site_description || '',
      footer_text: data.footer_text || ''
    }
  } catch (error) {
    toast.error('加载设置失败')
  }
})

async function saveSettings() {
  try {
    saving.value = true
    await api.put('/api/admin/settings', settings.value)
    toast.success('保存成功')
  } catch (error) {
    toast.error('保存失败')
  } finally {
    saving.value = false
  }
}
</script>
