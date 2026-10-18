<template>
  <div>
    <BannerTwo />
    <PromotionalBanner />
    <DealOfWeek />
    <TopSellingProducts />
    <TrendingProductsTwo />
    <DiscountTwo />
    <FeaturedProducts />
    <DiscountThree />
    <RecommendedProducts />
    <BigDeals />
    <TopBrandsTwo />
    <PopularProductsTwo />
    <Newsletter />
    <ShippingOverview />
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import BannerTwo from "~/components/layout/banner/BannerTwo.vue";
import PromotionalBanner from "~/components/containers/home/PromotionalBanner.vue";
import DealOfWeek from "~/components/containers/home/DealOfWeek.vue";
import TopSellingProducts from "~/components/containers/shop/TopSellingProducts.vue";
import TrendingProductsTwo from "~/components/containers/shop/TrendingProductsTwo.vue";
import DiscountTwo from "~/components/containers/shop/DiscountTwo.vue";
import FeaturedProducts from "~/components/containers/shop/FeaturedProducts.vue";
import DiscountThree from "~/components/containers/shop/DiscountThree.vue";
import RecommendedProducts from "~/components/containers/shop/RecommendedProducts.vue";
import BigDeals from "~/components/containers/shop/BigDeals.vue";
import TopBrandsTwo from "~/components/containers/shop/TopBrandsTwo.vue";
import PopularProductsTwo from "~/components/containers/shop/PopularProductsTwo.vue";
import Newsletter from "~/components/widgets/newsletter/Newsletter.vue";
import ShippingOverview from "~/components/containers/shipping/ShippingOverview.vue";

definePageMeta({
  layout: "layout-two",
});

useSeo({
  title: "Home Two",
  
  keywords: ["ecommerce", "fashion", "electronics", "grocery", "shop", "store", "multipurpose"],
});
</script>
