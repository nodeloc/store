<template>
  <div>
    <Breadcrumb :title="shop?.name || '店铺'" compact :trail="[{ label: '全部店铺', to: '/vendor' }]" />
    <div class="nl-container" style="padding-bottom: 96px">
      <div v-if="loadingShop" class="nl-skel" style="min-height: 220px"></div>
      <NlEmpty v-else-if="!shop" icon="ph ph-storefront" title="店铺不存在或已关闭" text="它可能还在审核中，或者已被下线。">
        <NuxtLink to="/vendor" class="nl-btn nl-btn--secondary nl-btn--sm">返回全部店铺</NuxtLink>
      </NlEmpty>

      <template v-else>
        <section class="nl-vhead">
          <span class="nl-vhead__logo">
            <img v-if="shop.logo" :src="shop.logo" :alt="shop.name" />
            <i v-else class="ph ph-storefront" aria-hidden="true"></i>
          </span>
          <div class="nl-vhead__info">
            <div class="nl-row" style="gap: var(--space-2); flex-wrap: wrap">
              <span v-if="shop.is_official" class="nl-badge nl-badge--volt"><i class="ph-fill ph-seal-check"></i>官方店铺</span>
              <span v-else class="nl-badge nl-badge--onink"><i class="ph ph-check-circle"></i>已审核入驻</span>
            </div>
            <h1 class="nl-vhead__name">{{ shop.name }}</h1>
            <p v-if="shop.description" class="nl-vhead__desc">{{ shop.description }}</p>
            <ul class="nl-vhead__meta">
              <li><b>{{ products.length }}</b> 件在售商品</li>
              <li><b>{{ totalSales.toLocaleString("en-US") }}</b> 累计售出</li>
              <li v-if="joined"><b>{{ joined }}</b> 入驻</li>
            </ul>
          </div>
          <div class="nl-vhead__aside">
            <a
              v-if="owner"
              :href="`https://www.nodeloc.com/u/${owner}`"
              target="_blank"
              rel="noopener noreferrer"
              class="nl-btn nl-btn--outline-ink nl-btn--sm"
            ><i class="ph ph-user-circle"></i>@{{ owner }}</a>
            <span v-if="shop.contact" class="nl-vhead__contact"><i class="ph ph-chat-circle"></i>{{ shop.contact }}</span>
          </div>
        </section>

        <div v-if="features.length" class="nl-features nl-vfeatures">
          <div v-for="f in features" :key="f.title" class="nl-feature">
            <i :class="f.icon || 'ph ph-seal-check'"></i>
            <div><b>{{ f.title }}</b><span>{{ f.text }}</span></div>
          </div>
        </div>

        <div class="nl-toolbar">
          <span class="nl-toolbar__count">{{ visible.length }} 件商品</span>
          <div class="nl-select">
            <label for="vendor-sort" class="nl-sr">排序</label>
            <select id="vendor-sort" v-model="sort">
              <option value="default">默认排序</option>
              <option value="sales">销量优先</option>
              <option value="price_asc">价格从低到高</option>
              <option value="price_desc">价格从高到低</option>
            </select>
            <i class="ph ph-caret-down"></i>
          </div>
        </div>
        <NlChips v-if="chipItems.length > 2" :items="chipItems" :active="cat" style="margin-bottom: var(--space-5)" @select="cat = $event" />

        <div v-if="loadingProducts" class="nl-grid">
          <div v-for="n in 4" :key="n" class="nl-skel"></div>
        </div>
        <NlEmpty v-else-if="!visible.length" title="这家店暂时没有在售商品" text="过段时间再来看看。" />
        <div v-else class="nl-grid">
          <NlProductCard v-for="(p, i) in visible" :key="p.id" :product="p" :index="i" />
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import type { NlChipItem } from "~/components/nl/Chips.vue";
import { categoryIcon, priceOf, type NlProduct } from "~/utils/nl";

definePageMeta({ layout: "layout-three" });

interface ShopDetail {
  id: number;
  name: string;
  description: string;
  logo: string;
  contact: string;
  features: string;
  is_official: boolean;
  created_at?: string;
  user?: { username?: string };
}

const route = useRoute();
const site = useSiteStore();
const shop = ref<ShopDetail | null>(null);
const products = ref<NlProduct[]>([]);
const loadingShop = ref(true);
const loadingProducts = ref(true);
const sort = ref("default");
const cat = ref("all");

const owner = computed(() => shop.value?.user?.username || "");
const joined = computed(() => {
  const d = shop.value?.created_at;
  return d ? new Date(d).toLocaleDateString("zh-CN", { year: "numeric", month: "2-digit" }) : "";
});
const totalSales = computed(() => products.value.reduce((n, p) => n + (p.sales_count || 0), 0));

const features = computed<{ icon?: string; title: string; text?: string }[]>(() => {
  try {
    const v = JSON.parse(shop.value?.features || "[]");
    return Array.isArray(v) ? v.filter((f) => f && f.title) : [];
  } catch {
    return [];
  }
});

const chipItems = computed<NlChipItem[]>(() => {
  const seen = new Map<string, NlChipItem>();
  for (const p of products.value) {
    const c = p.category;
    if (!c?.id) continue;
    const key = String(c.id);
    const hit = seen.get(key);
    if (hit) hit.count = (hit.count || 0) + 1;
    else {
      const ic = categoryIcon(c);
      seen.set(key, { key, label: c.name, icon: ic.cls, img: ic.img, count: 1 });
    }
  }
  return [{ key: "all", label: "全部", icon: "ph ph-squares-four", count: products.value.length }, ...seen.values()];
});

const visible = computed(() => {
  let list = cat.value === "all" ? [...products.value] : products.value.filter((p) => String(p.category?.id) === cat.value);
  if (sort.value === "sales") list.sort((a, b) => (b.sales_count || 0) - (a.sales_count || 0));
  else if (sort.value === "price_asc") list.sort((a, b) => priceOf(a).now - priceOf(b).now);
  else if (sort.value === "price_desc") list.sort((a, b) => priceOf(b).now - priceOf(a).now);
  return list;
});

async function load(id: string) {
  loadingShop.value = true;
  loadingProducts.value = true;
  cat.value = "all";
  try {
    shop.value = await $fetch<ShopDetail>(`/api/shops/${id}`);
  } catch {
    shop.value = null;
  } finally {
    loadingShop.value = false;
  }
  if (!shop.value) {
    loadingProducts.value = false;
    return;
  }
  try {
    const data = await $fetch<NlProduct[]>(`/api/shops/${id}/products`);
    products.value = Array.isArray(data) ? data : [];
  } catch {
    products.value = [];
  } finally {
    loadingProducts.value = false;
  }
}

watch(() => route.params.id, (id) => id && load(String(id)), { immediate: true });

watch(shop, (s) => {
  if (!s) return;
  const siteName = (site.settings.site_name as string) || "NodeLoc";
  useHead({
    title: `${s.name} | ${siteName}`,
    meta: [
      { name: "description", content: s.description || s.name },
      { property: "og:title", content: s.name },
      ...(s.logo ? [{ property: "og:image", content: s.logo }] : []),
    ],
  });
});
</script>
