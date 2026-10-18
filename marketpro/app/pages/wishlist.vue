<template>
  <div>
    <Breadcrumb title="我的收藏" :sub="wishlist.count ? `共 ${wishlist.count} 件商品` : ''" />
    <WishSection />
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import { useWishlistStore } from "~/stores/wishlist";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import WishSection from "~/components/containers/shop/WishSection.vue";

definePageMeta({ layout: "layout-three" });

const wishlist = useWishlistStore();
useSeo({ title: "收藏夹" });
</script>
