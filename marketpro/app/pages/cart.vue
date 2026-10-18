<template>
  <div>
    <Breadcrumb title="购物车" sub="确认商品与数量后前往结算。" />
    <CartSection />
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import CartSection from "~/components/containers/shop/CartSection.vue";

definePageMeta({ layout: "layout-three" });

useSeo({ title: "购物车" });
</script>
