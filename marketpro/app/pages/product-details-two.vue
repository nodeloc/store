<template>
  <div>
    <Breadcrumb title="Potato Chips" />
    <ProductDetailsTwo />
    <NewArrival />
    <ShippingOverviewTwo />
    <NewsletterThree />
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import ProductDetailsTwo from "~/components/containers/shop/ProductDetailsTwo.vue";
import NewArrival from "~/components/containers/shop/NewArrival.vue";
import ShippingOverviewTwo from "~/components/containers/shipping/ShippingOverviewTwo.vue";
import NewsletterThree from "~/components/widgets/newsletter/NewsletterThree.vue";

definePageMeta({
  layout: "layout-five",
});

useSeo({
  title: "Product Details",
  
  keywords: ["ecommerce", "fashion", "electronics", "grocery", "shop", "store", "multipurpose"],
});
</script>
