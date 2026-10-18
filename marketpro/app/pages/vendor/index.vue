<template>
  <div>
    <Breadcrumb title="全部店铺" sub="社区商家开的小店，积分付款，多数商品付款后自动发货。" />
    <div class="nl-container" style="padding-bottom: 96px">
      <div class="nl-toolbar" style="margin-top: 0">
        <span class="nl-toolbar__count">{{ loading ? "加载中…" : `${total} 家店铺` }}</span>
        <form class="nl-toolbar__search" role="search" @submit.prevent="search">
          <label for="vendor-q" class="nl-sr">搜索店铺</label>
          <input id="vendor-q" v-model="keyword" class="nl-input" type="search" placeholder="搜索店铺名称或简介" />
          <button type="submit" aria-label="搜索"><i class="ph ph-magnifying-glass"></i></button>
        </form>
        <NuxtLink to="/become-seller" class="nl-btn nl-btn--secondary nl-btn--sm"><i class="ph ph-storefront"></i>我也要开店</NuxtLink>
      </div>

      <div v-if="loading" class="nl-vgrid">
        <div v-for="n in 6" :key="n" class="nl-skel" style="min-height: 228px"></div>
      </div>
      <NlEmpty
        v-else-if="!shops.length"
        icon="ph ph-storefront"
        :title="query ? `没有找到「${query}」相关的店铺` : '还没有店铺'"
        :text="query ? '换个关键词试试。' : '成为第一家入驻的店铺吧。'"
      >
        <NuxtLink v-if="query" to="/vendor" class="nl-btn nl-btn--secondary nl-btn--sm" @click="reset">查看全部店铺</NuxtLink>
        <NuxtLink v-else to="/become-seller" class="nl-btn nl-btn--volt nl-btn--sm">申请开店</NuxtLink>
      </NlEmpty>
      <div v-else class="nl-vgrid">
        <NuxtLink
          v-for="(s, i) in shops"
          :key="s.id"
          v-reveal
          :to="`/vendor/${s.id}`"
          class="nl-vcard"
          data-reveal
          :style="{ '--i': i % 8 }"
        >
          <div class="nl-vcard__band">
            <span v-if="s.is_official" class="nl-badge nl-badge--volt"><i class="ph-fill ph-seal-check"></i>官方</span>
          </div>
          <span class="nl-vcard__logo">
            <img v-if="s.logo" :src="s.logo" :alt="s.name" loading="lazy" />
            <i v-else class="ph ph-storefront" aria-hidden="true"></i>
          </span>
          <div class="nl-vcard__body">
            <h3 class="nl-vcard__name">{{ s.name }}</h3>
            <p class="nl-vcard__desc" :class="{ 'is-empty': !s.description }">{{ s.description || "店主还没有填写简介" }}</p>
            <div class="nl-vcard__foot">
              <span class="nl-mono">{{ s.product_count ?? 0 }} 件商品</span>
              <span class="nl-vcard__go">进店<i class="ph ph-arrow-right"></i></span>
            </div>
          </div>
        </NuxtLink>
      </div>

      <NlPager :page="page" :total="pageCount" @go="goPage" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";

definePageMeta({ layout: "layout-three" });
useSeo({ title: "全部店铺" });

interface ShopItem {
  id: number;
  name: string;
  description: string;
  logo: string;
  is_official: boolean;
  product_count?: number;
}

const route = useRoute();
const router = useRouter();
const pageSize = 12;

const shops = ref<ShopItem[]>([]);
const total = ref(0);
const loading = ref(true);
const query = computed(() => (typeof route.query.q === "string" ? route.query.q.trim() : ""));
const page = computed(() => Math.max(1, Number(route.query.page) || 1));
const pageCount = computed(() => Math.ceil(total.value / pageSize));
const keyword = ref(query.value);

async function load() {
  loading.value = true;
  try {
    const res = await $fetch<{ data: ShopItem[]; total: number }>("/api/shops", {
      query: { page: page.value, page_size: pageSize, keyword: query.value || undefined },
    });
    shops.value = res.data ?? [];
    total.value = res.total ?? 0;
  } catch {
    shops.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
}

function search() {
  router.push({ query: { q: keyword.value.trim() || undefined } });
}

function reset() {
  keyword.value = "";
}

function goPage(p: number) {
  router.push({ query: { ...route.query, page: p > 1 ? p : undefined } });
  window.scrollTo({ top: 0, behavior: "smooth" });
}

watch(() => [route.query.q, route.query.page], load, { immediate: true });
</script>
