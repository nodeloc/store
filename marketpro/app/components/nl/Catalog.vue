<template>
  <div class="nl-container nl-catalog">
    <NlChips v-if="chips.length > 1" :items="chips" :active="activeChip" />
    <div class="nl-toolbar">
      <span class="nl-toolbar__count">
        <template v-if="loading">加载中…</template>
        <template v-else>共 {{ total.toLocaleString("en-US") }} 件商品</template>
        <template v-if="keyword"> · “{{ keyword }}”</template>
      </span>
      <form class="nl-toolbar__search" role="search" @submit.prevent="setQuery({ keyword: kw.trim() || undefined })">
        <input v-model="kw" class="nl-input" type="search" placeholder="搜索商品…" aria-label="搜索商品" />
        <button type="submit" aria-label="搜索"><i class="ph ph-magnifying-glass"></i></button>
      </form>
      <label class="nl-select">
        <span class="nl-sr">排序</span>
        <select :value="sort" @change="setQuery({ sort: ($event.target as HTMLSelectElement).value })">
          <option value="default">默认排序</option>
          <option value="newest">最新上架</option>
          <option value="sales">销量优先</option>
          <option value="price_asc">价格从低到高</option>
          <option value="price_desc">价格从高到低</option>
        </select>
        <i class="ph ph-caret-down"></i>
      </label>
    </div>

    <div v-if="loading" class="nl-grid">
      <div v-for="i in 8" :key="i" class="nl-skel"></div>
    </div>
    <NlEmpty v-else-if="!products.length" icon="ph ph-magnifying-glass" title="没有找到商品" text="换个关键词或分类试试。">
      <NuxtLink v-if="hasFilters" class="nl-btn nl-btn--secondary nl-btn--sm" :to="categoryId ? `/category/${categoryId}` : '/shop'">
        清除筛选
      </NuxtLink>
    </NlEmpty>
    <div v-else class="nl-grid">
      <NlProductCard v-for="(p, i) in products" :key="p.id" :product="p" :index="i" />
    </div>

    <NlPager :page="page" :total="totalPages" @go="p => setQuery({ page: p > 1 ? String(p) : undefined }, true)" />
  </div>
</template>

<script setup lang="ts">
import { useSiteStore } from "~/stores/site";
import type { NlChipItem } from "~/components/nl/Chips.vue";
import { categoryIcon, type NlCategory, type NlProduct } from "~/utils/nl";

const props = defineProps<{ categoryId?: number }>();
const emit = defineEmits<{ loaded: [total: number] }>();

const route = useRoute();
const router = useRouter();
const site = useSiteStore();

const PAGE_SIZE = 20;
const products = ref<NlProduct[]>([]);
const total = ref(0);
const totalPages = ref(1);
const loading = ref(true);

const q = (k: string) => (typeof route.query[k] === "string" ? (route.query[k] as string) : "");
const keyword = computed(() => q("keyword"));
const sort = computed(() => q("sort") || "default");
const promo = computed(() => q("promo") === "1");
const page = computed(() => Math.max(1, Number(q("page")) || 1));
const catId = computed(() => props.categoryId || Number(q("category_id")) || 0);
const hasFilters = computed(() => !!(keyword.value || promo.value || (!props.categoryId && catId.value)));

const kw = ref(keyword.value);
watch(keyword, v => (kw.value = v));

const activeChip = computed(() => (promo.value ? "promo" : catId.value ? String(catId.value) : "all"));
const chips = computed<NlChipItem[]>(() => {
  const cats = site.categories as NlCategory[];
  return [
    { key: "all", label: "全部", icon: "ph ph-squares-four", to: "/shop" },
    { key: "promo", label: "限时特惠", icon: "ph-fill ph-lightning", to: "/shop?promo=1", className: "nl-chip--promo" },
    ...cats.map(c => {
      const ic = categoryIcon(c);
      return { key: String(c.id), label: c.name, icon: ic.cls, img: ic.img, to: `/category/${c.id}`, count: c.product_count };
    }),
  ];
});

function setQuery(patch: Record<string, string | undefined>, keepPage = false) {
  const next: Record<string, string> = {};
  for (const [k, v] of Object.entries({ ...route.query, ...(keepPage ? {} : { page: undefined }), ...patch })) {
    if (typeof v === "string" && v && !(k === "sort" && v === "default")) next[k] = v;
  }
  router.push({ path: route.path, query: next });
}

let seq = 0;
async function load() {
  const id = ++seq;
  loading.value = true;
  const params: Record<string, string> = { page: String(page.value), page_size: String(PAGE_SIZE), sort: sort.value };
  if (catId.value) params.category_id = String(catId.value);
  if (keyword.value) params.keyword = keyword.value;
  if (promo.value) params.promo = "1";
  try {
    const res = await $fetch<{ products: NlProduct[]; total: number; total_pages: number }>("/api/products", { params });
    if (id !== seq) return;
    products.value = res.products || [];
    total.value = res.total || 0;
    totalPages.value = res.total_pages || 1;
  } catch {
    if (id !== seq) return;
    products.value = [];
    total.value = 0;
    totalPages.value = 1;
  } finally {
    if (id === seq) {
      loading.value = false;
      emit("loaded", total.value);
    }
  }
}

watch(
  () => [catId.value, keyword.value, sort.value, promo.value, page.value],
  (_n, old) => {
    if (old && old[4] !== page.value && import.meta.client) window.scrollTo({ top: 0, behavior: "smooth" });
    load();
  },
  { immediate: true }
);
onMounted(() => site.ensureLoaded());
</script>
