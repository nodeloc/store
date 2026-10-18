<template>
  <div>
    <Breadcrumb :title="product?.name || '商品详情'" :trail="trail" compact />

    <div v-if="loading" class="nl-container">
      <div class="nl-pd">
        <div class="nl-skel" style="aspect-ratio: 4 / 3; min-height: 0"></div>
        <div class="nl-skel"></div>
      </div>
    </div>

    <div v-else-if="!product" class="nl-container">
      <NlEmpty icon="ph ph-warning-circle" title="商品不存在或已下架" text="看看其他商品吧。">
        <NuxtLink class="nl-btn nl-btn--primary nl-btn--sm" to="/shop">浏览全部商品</NuxtLink>
      </NlEmpty>
    </div>

    <template v-else>
      <section class="nl-container">
        <div class="nl-pd">
          <div class="nl-gallery">
            <div class="nl-gallery__main">
              <img v-if="images.length" :key="images[activeIdx]" :src="images[activeIdx]" :alt="product.name" />
              <i v-else :class="icon" aria-hidden="true"></i>
            </div>
            <div v-if="images.length > 1" class="nl-gallery__thumbs">
              <button
                v-for="(img, i) in images"
                :key="img"
                type="button"
                :aria-current="i === activeIdx ? 'true' : undefined"
                :aria-label="`查看第 ${i + 1} 张图片`"
                @click="activeIdx = i"
              >
                <img :src="img" alt="" />
              </button>
            </div>
          </div>

          <div class="nl-pd__info">
            <div v-if="flags.length" class="nl-pd__flags">
              <span v-for="f in flags" :key="f.label" class="nl-badge" :class="f.tone && `nl-badge--${f.tone}`">
                <i v-if="f.icon" :class="f.icon"></i>{{ f.label }}
              </span>
            </div>
            <h1 class="nl-pd__title">{{ product.name }}</h1>
            <div class="nl-pd__meta">
              <NuxtLink v-if="product.shop_id" class="nl-pd__seller" :to="`/vendor/${product.shop_id}`">
                <i class="ph ph-storefront"></i>{{ product.shop?.name || "店铺" }}
              </NuxtLink>
              <span v-if="reviewCount > 0">
                <span class="star"><i class="ph-fill ph-star"></i></span> {{ avgRating.toFixed(1) }} · {{ reviewCount }} 条评价
              </span>
              <span v-else>暂无评价</span>
              <span>已售 {{ product.sales_count || 0 }}</span>
              <span class="nl-mono">SKU {{ String(product.id).padStart(6, "0") }}</span>
            </div>

            <div class="nl-pd__pricebox">
              <NlPrice :now="price.now" :was="price.was" :promo="price.promo" lg />
              <NlStock :sales="product.sales_count" :stock="product.stock_count" />
              <div v-if="product.is_promo_active && product.promo_end" class="nl-pd__promo">
                <span><i class="ph-fill ph-lightning"></i>限时特惠 · 距结束</span>
                <NlCountdown :until="product.promo_end" @end="load" />
              </div>
            </div>

            <div v-if="unavailable" class="nl-note nl-note--warning">
              <i class="ph ph-warning"></i>
              {{ product.is_active === false ? "商品已下架" : "库存不足，暂时无法购买" }}
            </div>
            <template v-else>
              <div class="nl-pd__buy">
                <div class="nl-qty" role="group" aria-label="购买数量">
                  <button type="button" :disabled="qty <= 1" aria-label="减少" @click="setQty(qty - 1)"><i class="ph ph-minus"></i></button>
                  <input :value="qty" type="number" min="1" :max="product.stock_count" aria-label="数量" @change="setQty(Number(($event.target as HTMLInputElement).value))" />
                  <button type="button" :disabled="qty >= product.stock_count" aria-label="增加" @click="setQty(qty + 1)"><i class="ph ph-plus"></i></button>
                </div>
                <button type="button" class="nl-btn nl-btn--volt nl-btn--lg" @click="buyNow">
                  <i class="ph-fill ph-lightning"></i>立即购买
                </button>
                <button type="button" class="nl-btn nl-btn--secondary nl-btn--lg nl-btn--icon" aria-label="加入购物车" title="加入购物车" @click="addToCart">
                  <i :class="added ? 'ph ph-check' : 'ph ph-shopping-cart-simple'"></i>
                </button>
                <button
                  type="button"
                  class="nl-btn nl-btn--secondary nl-btn--lg nl-btn--icon"
                  :aria-label="fav ? '取消收藏' : '收藏'"
                  :aria-pressed="fav ? 'true' : 'false'"
                  :style="fav ? { color: 'var(--danger)' } : undefined"
                  @click="wishlist.toggle({ ...product, price: price.now })"
                >
                  <i :class="fav ? 'ph-fill ph-heart' : 'ph ph-heart'"></i>
                </button>
              </div>
              <div class="nl-pd__total">
                <span>合计 {{ qty }} 件</span>
                <b><NlPrice :now="price.now * qty" /></b>
              </div>
            </template>

            <div class="nl-features">
              <div v-for="(f, i) in features" :key="i" class="nl-feature">
                <i :class="f.icon || 'ph-fill ph-seal-check'"></i>
                <div><b>{{ f.title }}</b><span>{{ f.text }}</span></div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="nl-container nl-section">
        <ProductContent :product-id="product.id" :description="product.description" />
      </section>

      <section v-if="related.length" class="nl-container nl-section" style="padding-bottom: 96px">
        <NlSectionHead eyebrow="MORE · 同类" title="同类推荐">
          <NuxtLink v-if="product.category_id" class="nl-link" :to="`/category/${product.category_id}`">
            查看分类 <i class="ph ph-arrow-right"></i>
          </NuxtLink>
        </NlSectionHead>
        <div class="nl-grid">
          <NlProductCard v-for="(p, i) in related" :key="p.id" :product="p" :index="i" />
        </div>
      </section>
      <div v-else style="height: 96px"></div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useCartStore } from "~/stores/cart";
import { useWishlistStore } from "~/stores/wishlist";
import { useSiteStore } from "~/stores/site";
import ProductContent from "~/components/containers/shop/ProductContent.vue";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import { categoryIcon, priceOf, productFlags, productImages, type NlProduct } from "~/utils/nl";

definePageMeta({ layout: "layout-four" });

const route = useRoute();
const router = useRouter();
const cart = useCartStore();
const wishlist = useWishlistStore();
const site = useSiteStore();
const siteName = computed(() => (site.settings.site_name as string) || "NodeLoc 商店");

const product = ref<NlProduct | null>(null);
const related = ref<NlProduct[]>([]);
const loading = ref(true);
const qty = ref(1);
const avgRating = ref(0);
const reviewCount = ref(0);
const activeIdx = ref(0);
const added = ref(false);

const images = computed(() => productImages(product.value));
const price = computed(() => (product.value ? priceOf(product.value) : { now: 0, was: undefined, promo: false }));
const flags = computed(() => (product.value ? productFlags(product.value) : []));
const icon = computed(() => categoryIcon(product.value?.category).cls);
const fav = computed(() => !!product.value && wishlist.has(product.value.id));
const unavailable = computed(() => !product.value || product.value.is_active === false || product.value.stock_count <= 0);
const trail = computed(() => {
  const t: { label: string; to?: string }[] = [{ label: "全部商品", to: "/shop" }];
  if (product.value?.category) t.push({ label: product.value.category.name, to: `/category/${product.value.category.id}` });
  return t;
});

const DEFAULT_FEATURES = [
  { icon: "ph-fill ph-lightning", title: "即时发货", text: "付款后卡密自动发放，无需等待。" },
  { icon: "ph-fill ph-shield-check", title: "确认后结算", text: "确认收货后款项才给卖家。" },
  { icon: "ph-fill ph-user-circle-check", title: "社区账号", text: "NodeLoc 账号登录，积分支付。" },
  { icon: "ph-fill ph-headset", title: "售后保障", text: "遇到问题联系店铺，快速处理。" },
];
const features = computed<{ icon?: string; title: string; text: string }[]>(() => {
  try {
    const parsed = JSON.parse(product.value?.shop?.features || "[]");
    if (Array.isArray(parsed) && parsed.length) return parsed.slice(0, 4);
  } catch {
    /* fall back to defaults */
  }
  return DEFAULT_FEATURES;
});

watch(product, p => {
  if (!p) return;
  const desc = p.description ? p.description.slice(0, 160) : `${p.name} - ${siteName.value}`;
  useHead({
    title: `${p.name} | ${siteName.value}`,
    meta: [
      { name: "description", content: desc },
      { property: "og:title", content: p.name },
      { property: "og:description", content: desc },
      ...(p.image ? [{ property: "og:image", content: p.image }] : []),
      { property: "og:type", content: "product" },
    ],
  });
});

function setQty(n: number) {
  const max = product.value?.stock_count || 1;
  qty.value = Math.max(1, Math.min(Number.isFinite(n) ? Math.floor(n) : 1, max));
}

async function load() {
  loading.value = true;
  qty.value = 1;
  activeIdx.value = 0;
  related.value = [];
  try {
    product.value = await $fetch<NlProduct>(`/api/products/${route.params.id}`, { credentials: "include" });
  } catch {
    product.value = null;
  }
  loading.value = false;
  if (!product.value) return;
  const id = product.value.id;
  const [rv, rel] = await Promise.all([
    $fetch<{ stats: { avg_rating: number; total: number } }>(`/api/products/${id}/reviews`).catch(() => null),
    product.value.category_id
      ? $fetch<{ products: NlProduct[] }>(`/api/products?category_id=${product.value.category_id}&page_size=9`).catch(() => null)
      : null,
  ]);
  avgRating.value = rv?.stats?.total ? rv.stats.avg_rating : 0;
  reviewCount.value = rv?.stats?.total || 0;
  related.value = (rel?.products || []).filter(p => p.id !== id).slice(0, 8);
}

function addToCart() {
  if (!product.value) return;
  cart.add({ ...product.value, price: price.value.now }, qty.value);
  added.value = true;
  setTimeout(() => (added.value = false), 1600);
}

function buyNow() {
  if (!product.value) return;
  cart.add({ ...product.value, price: price.value.now }, qty.value);
  router.push("/checkout");
}

onMounted(load);
watch(() => route.params.id, load);
</script>
