<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-xl font-bold text-zinc-900">系统设置</h2>
      <p class="text-sm text-zinc-500 mt-1">配置网站基本信息，修改后前台页面将实时更新</p>
    </div>
    
    <form @submit.prevent="saveSettings" class="space-y-6">
      <!-- Site Settings -->
      <div class="bg-white rounded-2xl border border-zinc-100 shadow-card overflow-hidden">
        <div class="px-6 py-4 border-b border-zinc-100 bg-zinc-50/50">
          <h3 class="text-sm font-semibold text-zinc-900">网站信息</h3>
          <p class="text-xs text-zinc-400 mt-0.5">设置展示在前台页面的网站名称、描述和页脚信息</p>
        </div>
        <div class="p-6 space-y-5">
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1.5">网站名称</label>
            <input
              v-model="settings.site_name"
              type="text"
              placeholder="例如：NodeLoc 社区发卡"
              class="w-full px-4 py-2.5 border border-zinc-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-green/20 focus:border-brand-green transition-colors"
            />
            <p class="text-xs text-zinc-400 mt-1">显示在页头导航栏和页脚处</p>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1.5">网站描述</label>
            <textarea
              v-model="settings.site_description"
              rows="2"
              placeholder="例如：安全便捷的数字商品自助购买平台"
              class="w-full px-4 py-2.5 border border-zinc-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-green/20 focus:border-brand-green transition-colors resize-none"
            ></textarea>
            <p class="text-xs text-zinc-400 mt-1">显示在首页导航栏的站名下方（桌面端）和信息栏（移动端）</p>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1.5">页脚文本</label>
            <input
              v-model="settings.footer_text"
              type="text"
              placeholder="例如：© 2026 NodeLoc 社区发卡系统"
              class="w-full px-4 py-2.5 border border-zinc-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-green/20 focus:border-brand-green transition-colors"
            />
            <p class="text-xs text-zinc-400 mt-1">显示在所有页面底部的版权信息</p>
          </div>

          <div>
            <label class="block text-sm font-medium text-zinc-700 mb-1.5">公告信息</label>
            <input
              v-model="settings.announcement"
              type="text"
              placeholder="输入公告内容，留空则不显示"
              class="w-full px-4 py-2.5 border border-zinc-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-brand-green/20 focus:border-brand-green transition-colors"
            />
            <p class="text-xs text-zinc-400 mt-1">显示在首页信息栏，支持一行文字公告</p>
          </div>
        </div>
      </div>
      
      <!-- Info Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="bg-white rounded-2xl border border-zinc-100 shadow-card p-5">
          <div class="flex items-start gap-3">
            <div class="w-9 h-9 rounded-xl bg-blue-50 flex items-center justify-center flex-shrink-0">
              <Info class="w-4.5 h-4.5 text-blue-600" />
            </div>
            <div class="space-y-1">
              <h4 class="text-sm font-semibold text-zinc-900">OAuth 登录</h4>
              <p class="text-xs text-zinc-500 leading-relaxed">系统默认使用 NodeLoc OAuth 2.0 登录，需在 .env 中配置 Client ID 和 Secret</p>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-2xl border border-zinc-100 shadow-card p-5">
          <div class="flex items-start gap-3">
            <div class="w-9 h-9 rounded-xl bg-emerald-50 flex items-center justify-center flex-shrink-0">
              <CreditCard class="w-4.5 h-4.5 text-emerald-600" />
            </div>
            <div class="space-y-1">
              <h4 class="text-sm font-semibold text-zinc-900">社区支付</h4>
              <p class="text-xs text-zinc-500 leading-relaxed">系统默认使用 NodeLoc 社区积分支付，需在 .env 中配置支付 ID 和 Secret</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Configuration Hint -->
      <div class="bg-white rounded-2xl border border-zinc-100 shadow-card overflow-hidden">
        <div class="px-6 py-4 border-b border-zinc-100 bg-zinc-50/50">
          <h3 class="text-sm font-semibold text-zinc-900">配置说明</h3>
        </div>
        <div class="p-6 space-y-2.5 text-sm text-zinc-500">
          <p class="flex items-start gap-2">
            <span class="text-zinc-300 mt-0.5">•</span>
            <span>OAuth 和支付参数请在 <code class="bg-zinc-100 px-2 py-0.5 rounded-md border border-zinc-200 font-mono text-xs text-zinc-700">.env</code> 文件中配置</span>
          </p>
          <p class="flex items-start gap-2">
            <span class="text-zinc-300 mt-0.5">•</span>
            <span>必需配置项：<code class="font-mono text-xs text-zinc-700">NODELOC_CLIENT_ID</code>、<code class="font-mono text-xs text-zinc-700">NODELOC_CLIENT_SECRET</code>、<code class="font-mono text-xs text-zinc-700">PAYMENT_ID</code>、<code class="font-mono text-xs text-zinc-700">PAYMENT_SECRET</code></span>
          </p>
          <p class="flex items-start gap-2">
            <span class="text-zinc-300 mt-0.5">•</span>
            <span>配置完成后需要重启服务才能生效</span>
          </p>
        </div>
      </div>
      
      <div class="flex justify-end">
        <button
          type="submit"
          :disabled="saving"
          class="inline-flex items-center gap-2 px-6 py-2.5 bg-brand-gradient text-white text-sm font-medium rounded-xl hover:shadow-glow disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300"
        >
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
  footer_text: '',
  announcement: ''
})

onMounted(async () => {
  try {
    const response = await api.get('/api/admin/settings')
    const data = response.data.settings
    settings.value = {
      site_name: data.site_name || '',
      site_description: data.site_description || '',
      footer_text: data.footer_text || '',
      announcement: data.announcement || ''
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
