<template>
  <nav v-if="total > 1" class="nl-pager" aria-label="分页">
    <button type="button" :disabled="page <= 1" aria-label="上一页" @click="go(page - 1)">
      <i class="ph ph-arrow-left"></i>
    </button>
    <button
      v-for="p in pages"
      :key="p"
      type="button"
      :aria-current="p === page ? 'page' : undefined"
      @click="go(p)"
    >
      {{ String(p).padStart(2, "0") }}
    </button>
    <button type="button" :disabled="page >= total" aria-label="下一页" @click="go(page + 1)">
      <i class="ph ph-arrow-right"></i>
    </button>
  </nav>
</template>

<script setup lang="ts">
const props = defineProps<{ page: number; total: number }>();
const emit = defineEmits<{ go: [page: number] }>();

const pages = computed(() => {
  const start = Math.max(1, Math.min(props.page - 2, props.total - 4));
  const end = Math.min(props.total, start + 4);
  const out: number[] = [];
  for (let i = start; i <= end; i++) out.push(i);
  return out;
});

function go(p: number) {
  if (p >= 1 && p <= props.total && p !== props.page) emit("go", p);
}
</script>
