<template>
  <div :class="['list-grid-wrapper grid-cols-4', { 'list-view': isListView }]">
    <div v-if="loading" class="text-center py-48 text-gray-400 w-100">
      <div class="spinner-border text-main-600" role="status"></div>
    </div>
    <div v-else-if="products.length === 0" class="text-center py-48 text-gray-400 w-100">
      <i class="ph ph-package" style="font-size:3rem"></i>
      <p class="mt-16">暂无商品</p>
    </div>
    <template v-else>
      <div
        v-for="product in products"
        :key="product.id"
        class="product-card h-100 p-16 border border-gray-100 hover-border-main-600 rounded-16 position-relative transition-2"
      >
        <div class="product-card__thumb rounded-8 bg-gray-50 position-relative">
          <NuxtLink :to="`/product/${product.id}`" class="w-100 h-100 flex-center d-block">
            <img v-if="product.image" :src="product.image" :alt="product.name" class="w-auto max-w-unset" style="max-height:180px;object-fit:contain" />
            <div v-else class="flex-center text-gray-300" style="min-height:120px;font-size:3rem">
              <i class="ph ph-package"></i>
            </div>
          </NuxtLink>
          <div v-if="product.orig_price && product.orig_price > product.price"
            class="position-absolute inset-block-start-0 inset-inline-start-0 mt-16 ms-16 z-1">
            <span class="text-white w-40 h-40 d-flex justify-content-center align-items-center bg-danger-600 rounded-circle shadow-sm text-xs fw-semibold">
              特价
            </span>
          </div>
        </div>

        <div class="product-card__content mt-16 w-100">
          <h6 class="title text-lg fw-semibold my-16">
            <NuxtLink :to="`/product/${product.id}`" class="link text-line-2">{{ product.name }}</NuxtLink>
          </h6>

          <div class="product-card__price mt-16 mb-16">
            <span class="text-heading text-md fw-semibold">
              <i class="ph-fill ph-lightning-a" style="font-size:0.85em;vertical-align:middle;margin-right:1px"></i>{{ Number(product.effective_price || product.price || 0).toFixed(2) }}
            </span>
            <span v-if="product.orig_price && product.orig_price > product.price"
              class="text-gray-400 text-md fw-semibold text-decoration-line-through ms-8">
              <i class="ph-fill ph-lightning-a" style="font-size:0.85em;vertical-align:middle;margin-right:1px"></i>{{ Number(product.orig_price).toFixed(2) }}
            </span>
          </div>

          <div v-if="product.sales_count" class="flex-align gap-6 mb-16">
            <span class="text-xs fw-bold text-gray-600">{{ product.sales_count }}</span>
            <span class="text-warning-600 d-flex"><i class="ph-fill ph-fire"></i></span>
            <span class="text-xs text-gray-500">Sold</span>
          </div>

          <NuxtLink
            :to="`/product/${product.id}`"
            class="product-card__cart btn bg-gray-50 text-heading hover-bg-main-600 hover-text-white py-11 px-24 rounded-8 flex-center gap-8 fw-medium w-100"
          >
            立即购买 <i class="ph ph-shopping-cart"></i>
          </NuxtLink>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { useViewToggle } from "~/composables/useViewToggle";

const route = useRoute();
const { isListView } = useViewToggle();
const products = ref<any[]>([]);
const loading = ref(true);

onMounted(async () => {
  const id = route.params.id as string;
  if (!id) { loading.value = false; return; }
  try {
    const data = await $fetch<any[]>(`/api/shops/${id}/products`);
    products.value = Array.isArray(data) ? data : [];
  } catch {
    products.value = [];
  } finally {
    loading.value = false;
  }
});
</script>
