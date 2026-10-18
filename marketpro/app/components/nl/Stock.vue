<template>
  <div v-reveal class="nl-stock" :class="{ 'nl-stock--low': s.low, 'nl-stock--out': s.out }">
    <div class="nl-stock__track" role="presentation">
      <div class="nl-stock__fill" :style="{ width: (s.out ? 100 : s.pct) + '%' }"></div>
    </div>
    <div class="nl-stock__meta">
      <span>{{ s.out ? 100 : s.pct }}%</span>
      <span v-if="s.out">已售罄</span>
      <b v-else-if="s.low">仅剩 {{ s.left }} 件</b>
      <span v-else>剩余 {{ s.left }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { stockOf } from "~/utils/nl";

const props = defineProps<{ sales?: number; stock: number }>();
const s = computed(() => stockOf({ sales_count: props.sales, stock_count: props.stock }));
</script>
