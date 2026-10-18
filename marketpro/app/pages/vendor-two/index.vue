<template>
  <div>
    <Breadcrumb title="All Vendors" />
    <VendorSectionTwo />
    <ShippingOverviewTwo />
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import VendorSectionTwo from "~/components/containers/vendor/VendorSectionTwo.vue";
import ShippingOverviewTwo from "~/components/containers/shipping/ShippingOverviewTwo.vue";

definePageMeta({
  layout: "layout-three",
});

useSeo({
  title: "Vendor Two",
  
  keywords: ["ecommerce", "fashion", "electronics", "grocery", "shop", "store", "multipurpose"],
});
</script>
