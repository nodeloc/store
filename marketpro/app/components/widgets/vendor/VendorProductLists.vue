<script setup lang="ts">
import { products } from "~/data/vendor-list";

const formatRatingCount = (count: number): string => {
  return count >= 1000 ? (count / 1000).toFixed(0) + "k" : count.toString();
};
</script>

<template>
  <div
    class="row row-cols-1 row-cols-sm-2 row-cols-md-3 row-cols-lg-3 row-cols-xl-4 row-cols-xxl-5 g-12"
  >
    <div
      v-for="(product, index) in products"
      :key="product.id"
      class="col"
      :data-aos="'fade-up'"
      :data-aos-duration="(index + 1) * 200"
    >
      <div
        class="product-card h-100 p-8 border border-gray-100 hover-border-main-600 rounded-16 position-relative transition-2"
      >
        <span
          v-if="product.badge"
          :class="['product-card__badge px-8 py-4 text-sm text-white', product.badgeColor]"
          >{{ product.badge }}</span
        >

        <NuxtLink to="/product-details" class="product-card__thumb flex-center overflow-hidden">
          <NuxtImg :src="product.image" :alt="product.name" />
        </NuxtLink>

        <div class="product-card__content p-sm-2 w-100">
          <h6 class="title text-lg fw-semibold mt-12 mb-8">
            <NuxtLink to="/product-details" class="link text-line-2">{{ product.name }}</NuxtLink>
          </h6>

          <div class="flex-align gap-4">
            <span class="text-main-600 text-md d-flex">
              <i class="ph-fill ph-storefront"></i>
            </span>
            <span class="text-gray-500 text-xs">By {{ product.store }}</span>
          </div>

          <div class="product-card__content mt-12">
            <div class="product-card__price mb-8">
              <span class="text-heading text-md fw-semibold"
                >${{ product.price.toFixed(2) }}
                <span class="text-gray-500 fw-normal">/Qty</span>
              </span>
              <span
                v-if="product.originalPrice"
                class="text-gray-400 text-md fw-semibold text-decoration-line-through"
                >${{ product.originalPrice.toFixed(2) }}</span
              >
            </div>

            <div class="flex-align gap-6">
              <span class="text-xs fw-bold text-gray-600">{{ product.rating }}</span>
              <span class="text-15 fw-bold text-warning-600 d-flex">
                <i class="ph-fill ph-star"></i>
              </span>
              <span class="text-xs fw-bold text-gray-600">
                ({{ formatRatingCount(product.ratingCount) }})
              </span>
            </div>

            <NuxtLink
              to="/cart"
              class="product-card__cart btn bg-main-50 text-main-600 hover-bg-main-600 hover-text-white py-11 px-24 rounded-pill flex-align gap-8 mt-24 w-100 justify-content-center"
            >
              Add To Cart <i class="ph ph-shopping-cart"></i>
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
