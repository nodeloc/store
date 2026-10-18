<template>
  <section class="vendor-two py-80">
    <div class="container container-lg">

      <!-- 搜索栏 -->
      <div class="d-flex align-items-center justify-content-between flex-wrap mb-48 gap-16">
        <form class="input-group w-100 max-w-418" novalidate @submit.prevent="doSearch">
          <input
            v-model="keyword"
            type="text"
            class="form-control common-input rounded-start-3"
            placeholder="搜索店铺名称..."
          />
          <button
            type="submit"
            class="input-group-text border-0 bg-main-two-600 rounded-end-3 text-white text-2xl hover-bg-main-two-700 px-24"
          >
            <i class="ph ph-magnifying-glass"></i>
          </button>
        </form>

        <div class="d-flex align-items-center justify-content-between justify-content-sm-end gap-16 flex-grow-1">
          <div class="text-gray-600 text-md flex-shrink-0">
            <span class="text-neutral-900 fw-semibold">{{ total }}</span> 家店铺
          </div>
          <div class="d-flex align-items-center gap-8 d-sm-flex d-none">
            <button type="button"
              :class="['grid-btn text-2xl d-flex w-48 h-48 border rounded-8 justify-content-center align-items-center',
                !isListView ? 'border-main-600 text-white bg-main-600' : 'border-gray-100']"
              aria-label="Grid view" @click="isListView = false">
              <i class="ph ph-squares-four"></i>
            </button>
            <button type="button"
              :class="['list-btn text-2xl d-flex w-48 h-48 border rounded-8 justify-content-center align-items-center',
                isListView ? 'border-main-600 text-white bg-main-600' : 'border-gray-100']"
              aria-label="List view" @click="isListView = true">
              <i class="ph ph-list-bullets"></i>
            </button>
          </div>
        </div>
      </div>

      <!-- 店铺列表 -->
      <div :class="['list-grid-wrapper vendors-two-item-wrapper grid-cols-3', { 'list-view': isListView }]">
        <div v-if="loading" class="text-center py-48 text-gray-400 w-100">
          <div class="spinner-border text-main-600" role="status"></div>
        </div>
        <div v-else-if="shops.length === 0" class="text-center py-48 text-gray-400 w-100">
          <i class="ph ph-storefront d-block mb-12" style="font-size:3rem"></i>
          <p>{{ keyword ? `未找到包含"${keyword}"的店铺` : '暂无店铺' }}</p>
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
              <div class="position-absolute top-0 inset-inline-start-0 w-100 h-100 p-24 z-1 d-flex flex-column justify-content-between">
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
                  <span class="flex-center text-main-two-600 text-2xl flex-shrink-0"><i class="ph ph-info"></i></span>
                  <p class="text-md text-gray-900 mb-0 text-line-2">{{ shop.description }}</p>
                </div>
                <div v-if="shop.contact" class="flex-align gap-8">
                  <span class="flex-center text-main-two-600 text-2xl flex-shrink-0"><i class="ph ph-chat-circle"></i></span>
                  <span class="text-md text-gray-900">{{ shop.contact }}</span>
                </div>
              </div>
              <NuxtLink :to="`/vendor/${shop.id}`"
                class="btn bg-neutral-600 hover-bg-neutral-700 text-white py-12 px-24 rounded-8 flex-center gap-8 fw-medium mt-24">
                进入店铺
                <span class="text-xl d-flex text-main-two-600"><i class="ph ph-storefront"></i></span>
              </NuxtLink>
            </div>
          </div>
        </template>
      </div>

      <!-- 分页 -->
      <ul v-if="pageCount > 1" class="pagination flex-center flex-wrap gap-16 mt-48">
        <li class="page-item">
          <button class="page-link h-64 w-64 flex-center text-xxl rounded-circle fw-medium text-neutral-600 border border-gray-100"
            :disabled="page <= 1" @click="goPage(page - 1)">
            <i class="ph-bold ph-arrow-left"></i>
          </button>
        </li>
        <li v-for="n in pageCount" :key="n" class="page-item">
          <button class="page-link h-64 w-64 flex-center text-md rounded-circle fw-medium text-neutral-600 border border-gray-100"
            :class="{ 'bg-main-600 text-white border-main-600': n === page }"
            @click="goPage(n)">
            {{ String(n).padStart(2, '0') }}
          </button>
        </li>
        <li class="page-item">
          <button class="page-link h-64 w-64 flex-center text-xxl rounded-circle fw-medium text-neutral-600 border border-gray-100"
            :disabled="page >= pageCount" @click="goPage(page + 1)">
            <i class="ph-bold ph-arrow-right"></i>
          </button>
        </li>
      </ul>

    </div>
  </section>
</template>

<script setup lang="ts">
const shops = ref<any[]>([]);
const loading = ref(true);
const total = ref(0);
const page = ref(1);
const pageSize = 12;
const keyword = ref("");
const isListView = ref(false);

const pageCount = computed(() => total.value > 0 ? Math.ceil(total.value / pageSize) : 0);

async function fetchShops() {
  loading.value = true;
  try {
    const params: Record<string, string> = {
      page: String(page.value),
      page_size: String(pageSize),
    };
    if (keyword.value.trim()) params.keyword = keyword.value.trim();
    const data = await $fetch<any>(`/api/shops?${new URLSearchParams(params)}`);
    shops.value = data.data ?? [];
    total.value = data.total ?? 0;
  } catch {
    shops.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
}

function doSearch() {
  page.value = 1;
  fetchShops();
}

function goPage(n: number) {
  if (n < 1 || n > pageCount.value) return;
  page.value = n;
  fetchShops();
  window.scrollTo({ top: 0, behavior: "smooth" });
}

onMounted(fetchShops);
</script>
