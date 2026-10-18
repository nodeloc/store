<template>
  <div :class="['list-grid-wrapper vendors-two-item-wrapper grid-cols-3', { 'list-view': isListView }]">
    <div v-if="loading" class="text-center py-48 text-gray-400 w-100">
      <div class="spinner-border text-main-600" role="status"></div>
    </div>
    <template v-else>
      <div
        v-for="shop in shops"
        :key="shop.id"
        class="vendors-two-item rounded-12 overflow-hidden bg-color-three border border-neutral-50 hover-border-main-two-600 transition-2"
      >
        <div class="vendors-two-item__top bg-overlay style-two position-relative">
          <div class="vendors-two-item__thumbs h-210">
            <img v-if="shop.logo" :src="shop.logo" :alt="shop.name" class="cover-img" />
            <div v-else class="cover-img bg-main-50 d-flex align-items-center justify-content-center">
              <i class="ph ph-storefront text-main-200" style="font-size:4rem"></i>
            </div>
          </div>
          <div
            class="position-absolute top-0 inset-inline-start-0 w-100 h-100 p-24 z-1 d-flex flex-column justify-content-between"
          >
            <div class="d-flex align-items-center justify-content-between">
              <span class="w-80 h-80 flex-center bg-white rounded-circle flex-shrink-0 overflow-hidden border border-gray-100">
                <img v-if="shop.logo" :src="shop.logo" :alt="shop.name" class="w-100 h-100 object-fit-cover" />
                <i v-else class="ph ph-storefront text-gray-400 text-2xl"></i>
              </span>
              <span v-if="shop.is_official" class="border border-white px-16 py-8 rounded-pill text-white text-sm">
                官方认证
              </span>
            </div>
            <div class="mt-16">
              <h6 class="text-white fw-semibold mb-8">
                <NuxtLink :to="`/vendor/${shop.id}`" class="text-white hover-text-main-two-600">
                  {{ shop.name }}
                </NuxtLink>
              </h6>
              <span class="text-sm text-white opacity-75">{{ shop.product_count ?? 0 }} 件商品</span>
            </div>
          </div>
        </div>

        <div class="vendors-two-item__content p-24 flex-grow-1">
          <div class="d-flex flex-column gap-14">
            <div v-if="shop.description" class="flex-align gap-8">
              <span class="flex-center text-main-two-600 text-2xl flex-shrink-0">
                <i class="ph ph-info"></i>
              </span>
              <p class="text-md text-gray-900 mb-0 text-line-2">{{ shop.description }}</p>
            </div>
            <div v-if="shop.contact" class="flex-align gap-8">
              <span class="flex-center text-main-two-600 text-2xl flex-shrink-0">
                <i class="ph ph-chat-circle"></i>
              </span>
              <span class="text-md text-gray-900">{{ shop.contact }}</span>
            </div>
          </div>

          <NuxtLink
            :to="`/vendor/${shop.id}`"
            class="btn bg-neutral-600 hover-bg-neutral-700 text-white py-12 px-24 rounded-8 flex-center gap-8 fw-medium mt-24"
          >
            进入店铺
            <span class="text-xl d-flex text-main-two-600">
              <i class="ph ph-storefront"></i>
            </span>
          </NuxtLink>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useViewToggle } from "~/composables/useViewToggle";

const { isListView } = useViewToggle();
const shops = ref<any[]>([]);
const loading = ref(true);

onMounted(async () => {
  try {
    const data = await $fetch<any>("/api/shops?page=1&page_size=200");
    shops.value = (data as any).data ?? [];
  } catch {
    shops.value = [];
  } finally {
    loading.value = false;
  }
});
</script>
