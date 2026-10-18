<template>
  <section class="vendor-two-details py-80">
    <div class="container container-lg">
      <div class="vendor-two-details-wrapper d-flex flex-wrap align-items-start gap-24">
        <ShopSidebarThree />

        <div class="vendor-two-details__contents">
          <VendorListThree />
          <VendorPagination />
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import ShopSidebarThree from "~/components/widgets/sidebar/ShopSidebarThree.vue";
import VendorListThree from "~/components/widgets/vendor/VendorListThree.vue";
import VendorPagination from "~/components/widgets/vendor/VendorPagination.vue";
</script>
