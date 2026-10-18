<template>
  <div v-if="items.length" class="nl-ticker" aria-label="实时成交">
    <div class="nl-ticker__track" :style="{ animationDuration: `${Math.max(24, items.length * 6)}s` }">
      <template v-for="copy in 2" :key="copy">
        <NuxtLink
          v-for="(it, i) in items"
          :key="`${copy}-${i}`"
          :to="`/product/${it.product_id}`"
          class="nl-ticker__item"
          :aria-hidden="copy === 2 ? 'true' : undefined"
          :tabindex="copy === 2 ? -1 : undefined"
        >
          <span class="nl-live"></span>
          <i :class="iconOf(it)"></i>
          {{ it.buyer }} 刚刚购买 <b>{{ it.product_name }}</b>
          <span class="nl-mono">{{ relTime(it.paid_at) }}</span>
        </NuxtLink>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { categoryIcon, relTime } from "~/utils/nl";

interface Purchase {
  buyer: string;
  product_id: number;
  product_name: string;
  category_icon?: string;
  paid_at?: string;
}

const items = ref<Purchase[]>([]);

function iconOf(it: Purchase) {
  const r = categoryIcon({ icon: it.category_icon, name: it.product_name });
  return r.img ? "ph ph-package" : r.cls;
}

async function load() {
  try {
    const res = await $fetch<{ items: Purchase[] }>("/api/recent-purchases");
    items.value = res.items || [];
  } catch {
    /* ticker is decorative; keep the last list */
  }
}

let timer: ReturnType<typeof setInterval> | null = null;
onMounted(() => {
  load();
  timer = setInterval(load, 60_000);
});
onBeforeUnmount(() => timer && clearInterval(timer));
</script>
