<template>
  <div class="min-h-screen flex items-center justify-center bg-white px-4">
    <div class="max-w-xl mx-auto text-center py-16 space-y-6">
      <div class="flex justify-center">
        <div class="w-32 h-32 bg-zinc-50 rounded-full flex items-center justify-center">
          <FileQuestion class="w-16 h-16 text-zinc-400" />
        </div>
      </div>
      <div class="space-y-2">
        <h1 class="text-6xl font-bold font-mono text-zinc-900">404</h1>
        <h2 class="text-2xl font-semibold text-zinc-900">页面未找到</h2>
        <p class="text-zinc-600">抱歉，您访问的页面不存在</p>
      </div>
      <router-link
        to="/"
        class="inline-flex items-center justify-center space-x-2 px-6 py-3 bg-zinc-900 text-white font-medium rounded-lg hover:bg-zinc-800 transition"
      >
        <Home class="w-5 h-5" />
        <span>返回首页</span>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { FileQuestion, Home } from 'lucide-vue-next'
</script>
