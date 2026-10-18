<template>
  <div>
    <BannerThree />
    <PromoSection />
    <PopularCategories />
    <Marquee />
    <TrendingProducts />
    <Discount />
    <NewArrivalsThree />
    <DealSection />
    <PopularProducts />
    <TopBrands />
    <ShippingOverviewThree />
    <TestimonialSliderTwo />
    <InstagramFeed />
    <NewsletterTwo />
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import BannerThree from "~/components/layout/banner/BannerThree.vue";
import PromoSection from "~/components/containers/home/PromoSection.vue";
import PopularCategories from "~/components/containers/home/PopularCategories.vue";
import Marquee from "~/components/containers/home/Marquee.vue";
import TrendingProducts from "~/components/containers/shop/TrendingProducts.vue";
import Discount from "~/components/containers/shop/Discount.vue";
import NewArrivalsThree from "~/components/containers/shop/NewArrivalsThree.vue";
import DealSection from "~/components/containers/shop/DealSection.vue";
import PopularProducts from "~/components/containers/shop/PopularProducts.vue";
import TopBrands from "~/components/containers/shop/TopBrands.vue";
import ShippingOverviewThree from "~/components/containers/shipping/ShippingOverviewThree.vue";
import TestimonialSliderTwo from "~/components/containers/testimonials/TestimonialSliderTwo.vue";
import InstagramFeed from "~/components/containers/InstagramFeed.vue";
import NewsletterTwo from "~/components/widgets/newsletter/NewsletterTwo.vue";

definePageMeta({
  layout: "layout-three",
});

useSeo({
  title: "Home Three",
  
  keywords: ["ecommerce", "fashion", "electronics", "grocery", "shop", "store", "multipurpose"],
});
</script>
