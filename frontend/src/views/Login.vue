<template>
  <div class="min-h-screen flex items-center justify-center bg-zinc-50 p-4">
    <div class="w-full max-w-md space-y-8">
      <!-- Logo -->
      <div class="text-center space-y-2">
        <div class="flex justify-center">
          <div class="w-16 h-16 bg-zinc-900 rounded-2xl flex items-center justify-center">
            <Ticket class="w-8 h-8 text-white" />
          </div>
        </div>
        <h1 class="text-2xl font-bold text-zinc-900">NodeLoc 社区发卡</h1>
        <p class="text-sm text-zinc-600">使用 NodeLoc 账号登录</p>
      </div>
      
      <!-- Login Card -->
      <div class="bg-white border border-zinc-200 rounded-lg p-8 space-y-6">
        <div class="space-y-4">
          <p class="text-sm text-zinc-600 text-center">
            点击下方按钮将跳转到 NodeLoc 进行授权登录
          </p>
          
          <button
            @click="handleLogin"
            class="w-full flex items-center justify-center space-x-2 px-6 py-3 bg-zinc-900 text-white font-medium rounded-lg hover:bg-zinc-800 transition"
          >
            <LogIn class="w-5 h-5" />
            <span>使用 NodeLoc 登录</span>
          </button>
        </div>
        
        <div class="pt-4 border-t border-zinc-200">
          <p class="text-xs text-zinc-500 text-center">
            登录即表示您同意我们的服务条款和隐私政策
          </p>
        </div>
      </div>
      
      <!-- Footer -->
      <div class="text-center">
        <router-link to="/" class="text-sm text-zinc-600 hover:text-zinc-900 flex items-center justify-center space-x-1">
          <ArrowLeft class="w-4 h-4" />
          <span>返回首页</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Ticket, LogIn, ArrowLeft } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

onMounted(() => {
  // Check if already logged in
  // If so, redirect to original page or home
})

function handleLogin() {
  const redirect = route.query.redirect || '/'
  window.location.href = `/auth/login?redirect=${encodeURIComponent(redirect)}`
}
</script>
