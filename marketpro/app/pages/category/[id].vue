<template>
  <div>
    <Breadcrumb :title="catName || '分类商品'" :sub="sub" :trail="[{ label: '全部商品', to: '/shop' }]" />
    <NlCatalog :category-id="catId" @loaded="n => (count = n)" />
    <div style="height: 96px"></div>
  </div>
</template>

<script setup lang="ts">
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import { useSeo } from "~/composables/useSeo";
import { useSiteStore } from "~/stores/site";

definePageMeta({ layout: "layout-three" });

const route = useRoute();
const site = useSiteStore();
const count = ref<number | null>(null);

const catId = computed(() => Number(route.params.id) || 0);
const catName = computed(() => site.categories.find(c => c.id === catId.value)?.name || "");
const sub = computed(() => (count.value === null ? "" : `该分类共有 ${count.value} 件在售商品`));

useSeo({ title: "分类商品" });
watchEffect(() => catName.value && useHead({ title: catName.value }));
</script>
