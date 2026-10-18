<template>
  <div>
    <Breadcrumb title="Contact" />
    <ContactSection />
    <ShippingOverviewTwo />
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import ContactSection from "~/components/containers/ContactSection.vue";
import ShippingOverviewTwo from "~/components/containers/shipping/ShippingOverviewTwo.vue";

definePageMeta({
  layout: "layout-three",
});

useSeo({
  title: "联系我们",
  
  keywords: ["ecommerce", "fashion", "electronics", "grocery", "shop", "store", "multipurpose"],
});
</script>
