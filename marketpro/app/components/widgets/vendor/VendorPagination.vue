<template>
  <ul v-if="pageCount > 1" class="pagination flex-center flex-wrap gap-16 mt-48" role="navigation" aria-label="Pagination Navigation">
    <li class="page-item">
      <button
        class="page-link h-64 w-64 flex-center text-xxl rounded-circle fw-medium text-neutral-600 border border-gray-100"
        :disabled="currentPage <= 1"
        aria-label="Previous page"
        @click="emit('change', currentPage - 1)"
      >
        <i class="ph-bold ph-arrow-left"></i>
      </button>
    </li>

    <li v-for="n in pageCount" :key="n" class="page-item">
      <button
        class="page-link h-64 w-64 flex-center text-md rounded-circle fw-medium text-neutral-600 border border-gray-100"
        :class="{ 'bg-main-600 text-white border-main-600': n === currentPage }"
        @click="emit('change', n)"
      >
        {{ String(n).padStart(2, '0') }}
      </button>
    </li>

    <li class="page-item">
      <button
        class="page-link h-64 w-64 flex-center text-xxl rounded-circle fw-medium text-neutral-600 border border-gray-100"
        :disabled="currentPage >= pageCount"
        aria-label="Next page"
        @click="emit('change', currentPage + 1)"
      >
        <i class="ph-bold ph-arrow-right"></i>
      </button>
    </li>
  </ul>
</template>

<script setup lang="ts">
import { computed } from "vue";

const props = withDefaults(defineProps<{
  total?: number;
  page?: number;
  pageSize?: number;
}>(), { total: 0, page: 1, pageSize: 9 });

const emit = defineEmits<{ (e: "change", page: number): void }>();

const currentPage = computed(() => props.page);
const pageCount = computed(() => props.total > 0 ? Math.ceil(props.total / props.pageSize) : 0);
</script>
