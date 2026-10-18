<template>
  <div>
    <Breadcrumb title="结算" sub="积分付款，卡密自动发放。" :trail="[{ label: '购物车', to: '/cart' }]" />
    <CheckoutSection />
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import CheckoutSection from "~/components/containers/shop/CheckoutSection.vue";

definePageMeta({ layout: "layout-three" });

useSeo({ title: "结算" });
</script>
