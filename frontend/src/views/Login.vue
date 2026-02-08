<template>
  <div class="min-h-screen flex items-center justify-center bg-zinc-50/50 px-4">
    <div class="w-full max-w-sm space-y-8">
      <!-- Logo -->
      <div class="text-center space-y-3">
        <div class="flex justify-center">
          <img
            src="https://s.rmimg.com/original/2X/f/f917abb33b14bc40ecaf1defce5d1903d186b393.svg"
            alt="Logo"
            class="h-12 w-auto"
          />
        </div>
        <div>
          <h1 class="text-xl font-bold text-zinc-900">欢迎回来</h1>
          <p class="text-sm text-zinc-500 mt-1">使用 NodeLoc 账号登录</p>
        </div>
      </div>
      
      <!-- Login Card -->
      <div class="bg-white rounded-2xl border border-zinc-100 shadow-card p-6 sm:p-8 space-y-6">
        <p class="text-sm text-zinc-500 text-center leading-relaxed">
          点击下方按钮将跳转到 NodeLoc 进行授权登录
        </p>
        
        <button
          @click="handleLogin"
          class="w-full flex items-center justify-center gap-2 px-6 py-3.5 bg-brand-gradient text-white font-medium rounded-xl hover:shadow-glow transition-all duration-300 hover:scale-[1.01]"
        >
          <LogIn class="w-5 h-5" />
          <span>使用 NodeLoc 登录</span>
        </button>
        
        <div class="pt-4 border-t border-zinc-100">
          <p class="text-xs text-zinc-400 text-center">
            登录即表示您同意我们的服务条款和隐私政策
          </p>
        </div>
      </div>
      
      <!-- Back -->
      <div class="text-center">
        <router-link to="/" class="inline-flex items-center gap-1.5 text-sm text-zinc-500 hover:text-brand-green transition-colors">
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
import { LogIn, ArrowLeft } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

onMounted(() => {})

function handleLogin() {
  const redirect = route.query.redirect || '/'
  window.location.href = `/auth/login?redirect=${encodeURIComponent(redirect)}`
}
</script>
