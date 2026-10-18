<template>
  <div class="nl-container" style="padding-bottom: 96px">
    <NlEmpty v-if="!wishlist.items.length" icon="ph ph-heart" title="收藏夹是空的" text="在商品卡片上点亮小红心即可收藏。">
      <NuxtLink class="nl-btn nl-btn--volt" to="/shop"><i class="ph-fill ph-lightning"></i>去逛逛</NuxtLink>
    </NlEmpty>
    <div v-else class="nl-grid">
      <NlProductCard v-for="(p, i) in products" :key="p.id" :product="p" :index="i" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useWishlistStore } from "~/stores/wishlist";
import type { NlProduct } from "~/utils/nl";

const wishlist = useWishlistStore();

const products = computed<NlProduct[]>(() =>
  wishlist.items.map(i => ({ id: i.id, name: i.name, price: i.price, image: i.image, stock_count: i.stock }))
);
</script>
