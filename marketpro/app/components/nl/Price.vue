<template>
  <span class="nl-price" :class="{ 'nl-price--lg': lg, 'nl-price--promo': promo }">
    <span class="nl-price__now">
      <span class="nl-price__bolt" aria-hidden="true"><i class="ph-fill ph-lightning"></i></span>
      <span class="nl-sr">积分</span>{{ money(now) }}
    </span>
    <span v-if="was && was > now" class="nl-price__was"><span class="nl-sr">原价</span>{{ money(was) }}</span>
  </span>
</template>

<script setup lang="ts">
import { money } from "~/utils/nl";

defineProps<{ now: number; was?: number; promo?: boolean; lg?: boolean }>();
</script>
