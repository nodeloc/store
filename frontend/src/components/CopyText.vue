<template>
  <span 
    class="inline-flex items-center space-x-2 copy-click font-mono-data relative group"
    @click="handleCopy"
    :title="copied ? '已复制!' : '点击复制'"
  >
    <span>{{ text }}</span>
    <Check v-if="copied" class="w-4 h-4 text-green-600" />
    <Copy v-else class="w-4 h-4 text-zinc-400 opacity-0 group-hover:opacity-100 transition-opacity" />
  </span>
</template>

<script setup>
import { Copy, Check } from 'lucide-vue-next'
import { useCopy } from '@/composables/useCopy'

const props = defineProps({
  text: {
    type: String,
    required: true
  }
})

const { copied, copy } = useCopy()

async function handleCopy() {
  await copy(props.text)
}
</script>
