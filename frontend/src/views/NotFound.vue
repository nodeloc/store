<template>
  <div class="min-h-[60vh] flex items-center justify-center px-4">
    <div class="max-w-md mx-auto text-center space-y-6">
      <div class="flex justify-center">
        <div class="w-24 h-24 rounded-3xl bg-brand-gradient-subtle flex items-center justify-center">
          <FileQuestion class="w-12 h-12 text-brand-green/50" />
        </div>
      </div>
      <div class="space-y-2">
        <h1 class="text-5xl font-bold font-mono text-gradient">404</h1>
        <h2 class="text-xl font-semibold text-zinc-900">页面未找到</h2>
        <p class="text-sm text-zinc-500">抱歉，您访问的页面不存在</p>
      </div>
      <router-link
        to="/"
        class="inline-flex items-center gap-2 px-5 py-2.5 bg-brand-gradient text-white text-sm font-medium rounded-xl hover:shadow-glow transition-all duration-300"
      >
        <Home class="w-4 h-4" />
        <span>返回首页</span>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { FileQuestion, Home } from 'lucide-vue-next'
</script>
