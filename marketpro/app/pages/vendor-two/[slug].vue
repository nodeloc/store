<template>
  <div>
    <Breadcrumb title="Vendor Details" />
    <VendorDetailsTwo />
    <ShippingOverviewTwo />
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import VendorDetailsTwo from "~/components/containers/vendor/VendorDetailsTwo.vue";
import ShippingOverviewTwo from "~/components/containers/shipping/ShippingOverviewTwo.vue";

definePageMeta({
  layout: "layout-three",
});

useSeo({
  title: "Vendor Details Two",
  
  keywords: ["ecommerce", "fashion", "electronics", "grocery", "shop", "store", "multipurpose"],
});
</script>
