<template>
  <div>
    <Breadcrumb :title="title" :sub="sub" />
    <NlCatalog />
    <div style="height: 96px"></div>
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";

definePageMeta({ layout: "layout-three" });

const route = useRoute();
const title = computed(() => {
  if (route.query.promo === "1") return "今日闪购";
  if (typeof route.query.keyword === "string" && route.query.keyword) return "搜索结果";
  return "全部商品";
});
const sub = computed(() =>
  route.query.promo === "1" ? "限时特惠商品，售完或到点即恢复原价。" : "VPS、AI 账号、点卡与兑换码，付款即自动发货。"
);

useSeo({ title: "商品列表" });
</script>
