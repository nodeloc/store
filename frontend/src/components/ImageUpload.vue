<template>
  <div class="space-y-2">
    <label class="block text-sm font-medium text-zinc-700">{{ label }}</label>
    
    <!-- Preview -->
    <div v-if="imageUrl" class="relative w-32 h-32 border border-zinc-200 rounded-lg overflow-hidden group">
      <img :src="imageUrl" alt="Preview" class="w-full h-full object-cover">
      <div class="absolute inset-0 bg-black bg-opacity-50 opacity-0 group-hover:opacity-100 transition flex items-center justify-center">
        <button @click="removeImage" class="p-2 bg-red-600 text-white rounded-lg hover:bg-red-700">
          <Trash2 class="w-5 h-5" />
        </button>
      </div>
    </div>
    
    <!-- Upload Button -->
    <div v-else class="flex items-center space-x-3">
      <label class="relative cursor-pointer">
        <input
          type="file"
          accept="image/*"
          @change="handleFileSelect"
          class="hidden"
        >
        <div class="inline-flex items-center space-x-2 px-4 py-2 bg-zinc-100 text-zinc-900 text-sm font-medium rounded-lg hover:bg-zinc-200 transition">
          <Upload class="w-4 h-4" />
          <span>选择图片</span>
        </div>
      </label>
      <span class="text-xs text-zinc-500">支持 JPG、PNG、GIF、WEBP，最大 5MB</span>
    </div>
    
    <!-- Uploading -->
    <div v-if="uploading" class="flex items-center space-x-2 text-sm text-zinc-600">
      <Loader2 class="w-4 h-4 animate-spin" />
      <span>上传中...</span>
    </div>
    
    <!-- Error -->
    <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Upload, Trash2, Loader2 } from 'lucide-vue-next'
import api from '@/utils/api'
import { useToast } from '@/stores/toast'

const props = defineProps({
  label: {
    type: String,
    default: '商品图片'
  },
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])
const toast = useToast()

const imageUrl = ref(props.modelValue)
const uploading = ref(false)
const error = ref('')

watch(() => props.modelValue, (val) => {
  imageUrl.value = val
})

async function handleFileSelect(event) {
  const file = event.target.files[0]
  if (!file) return
  
  // 验证文件大小
  if (file.size > 5 * 1024 * 1024) {
    error.value = '图片大小不能超过 5MB'
    return
  }
  
  // 验证文件类型
  const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
  if (!validTypes.includes(file.type)) {
    error.value = '只支持 JPG、PNG、GIF、WEBP 格式'
    return
  }
  
  error.value = ''
  uploading.value = true
  
  try {
    const formData = new FormData()
    formData.append('file', file)
    
    const response = await api.post('/api/admin/upload/image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    
    imageUrl.value = response.data.url
    emit('update:modelValue', response.data.url)
    toast.success('上传成功')
  } catch (err) {
    error.value = err.response?.data?.error || '上传失败'
    toast.error(error.value)
  } finally {
    uploading.value = false
  }
}

function removeImage() {
  imageUrl.value = ''
  emit('update:modelValue', '')
}
</script>
