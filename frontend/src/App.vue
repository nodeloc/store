<template>
  <div id="app" class="min-h-screen bg-white text-zinc-900 antialiased">
    <router-view />
    
    <!-- Global Toast -->
    <Transition name="toast">
      <div
        v-if="toast.show"
        class="fixed bottom-4 right-4 px-4 py-3 bg-zinc-900 text-white rounded-lg shadow-lg flex items-center space-x-2 z-50"
        :class="{ 'bg-red-600': toast.type === 'error' }"
      >
        <CheckCircle v-if="toast.type === 'success'" class="w-5 h-5" />
        <AlertCircle v-if="toast.type === 'error'" class="w-5 h-5" />
        <span>{{ toast.message }}</span>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { CheckCircle, AlertCircle } from 'lucide-vue-next'
import { useToastStore } from '@/stores/toast'

const toast = useToastStore()
</script>
