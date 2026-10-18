<template>
  <div>
    <div class="nl-container" style="padding-top: 16px">
      <NlHero :stats="stats" />
    </div>

    <div v-if="chips.length > 1" v-reveal class="nl-container" style="padding-top: 24px" data-reveal>
      <NlChips :items="chips" active="all" />
    </div>

    <section class="nl-container nl-section">
      <NlSectionHead v-if="flashMode === 'promo'" eyebrow="FLASH · 闪购" title="今日闪购" live>
        <NlCountdown v-if="flashEnd" :until="flashEnd" />
        <NuxtLink class="nl-link" to="/shop?promo=1">查看全部 <i class="ph ph-arrow-right"></i></NuxtLink>
      </NlSectionHead>
      <NlSectionHead v-else eyebrow="HOT · 热销" title="热销榜单" live>
        <NuxtLink class="nl-link" to="/shop?sort=sales">查看全部 <i class="ph ph-arrow-right"></i></NuxtLink>
      </NlSectionHead>
      <div class="nl-grid">
        <template v-if="loading">
          <div v-for="i in 4" :key="i" class="nl-skel"></div>
        </template>
        <NlProductCard v-for="(p, i) in flash" v-else :key="p.id" :product="p" :index="i" />
      </div>
      <NlEmpty v-if="!loading && !flash.length" title="暂无商品" text="卖家正在上架中，稍后再来看看。" />
    </section>

    <section class="nl-container nl-section">
      <div class="nl-bento">
        <div v-reveal class="nl-bento__cell nl-bento__cell--wide" data-reveal style="--i: 0">
          <span class="nl-bento__num"><NlCountUp :value="bento.n" /></span>
          <div>
            <b>{{ bento.label }}</b><br />
            <span>付款成功即从库存取卡，订单页直接可见，一键复制。</span>
          </div>
        </div>
        <div v-reveal class="nl-bento__cell" data-reveal style="--i: 1">
          <i class="ph-fill ph-user-circle-check"></i><b>社区账号登录</b><span>NodeLoc OAuth，一键进入</span>
        </div>
        <div v-reveal class="nl-bento__cell" data-reveal style="--i: 2">
          <i class="ph-fill ph-shield-check"></i><b>确认后结算</b><span>确认收货后款项才给卖家</span>
        </div>
      </div>
    </section>

    <section class="nl-container nl-section" style="padding-bottom: 96px">
      <NlSectionHead eyebrow="PICKS · 推荐" title="为你推荐">
        <NuxtLink class="nl-link" to="/shop?sort=newest">查看全部 <i class="ph ph-arrow-right"></i></NuxtLink>
      </NlSectionHead>
      <div class="nl-grid">
        <template v-if="loading">
          <div v-for="i in 8" :key="i" class="nl-skel"></div>
        </template>
        <NlProductCard v-for="(p, i) in picks" v-else :key="p.id" :product="p" :index="i" />
      </div>
      <NlEmpty v-if="!loading && !picks.length" title="暂无商品" text="卖家正在上架中，稍后再来看看。" />
    </section>
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import { useSiteStore } from "~/stores/site";
import type { NlStats } from "~/components/nl/Hero.vue";
import type { NlChipItem } from "~/components/nl/Chips.vue";
import { categoryIcon, type NlCategory, type NlProduct } from "~/utils/nl";

useSeo({ title: "首页" });

const site = useSiteStore();
const stats = ref<NlStats | null>(null);
const flash = ref<NlProduct[]>([]);
const picks = ref<NlProduct[]>([]);
const flashMode = ref<"promo" | "sales">("promo");
const loading = ref(true);

type ListRes = { products: NlProduct[] };

const flashEnd = computed(() => {
  const now = Date.now();
  const ends = flash.value
    .map(p => (p.promo_end ? new Date(p.promo_end).getTime() : 0))
    .filter(t => t > now)
    .sort((a, b) => a - b);
  return ends.length ? new Date(ends[0]).toISOString() : null;
});

const bento = computed(() =>
  stats.value?.sold_keys
    ? { n: stats.value.sold_keys, label: "张卡密已自动发放" }
    : { n: stats.value?.products || 0, label: "件商品在售，自动发货" }
);

const chips = computed<NlChipItem[]>(() => {
  const cats = site.categories as NlCategory[];
  return [
    { key: "all", label: "全部", icon: "ph ph-squares-four", to: "/shop", count: stats.value?.products },
    ...cats.map(c => {
      const ic = categoryIcon(c);
      return {
        key: String(c.id),
        label: c.name,
        icon: ic.cls,
        img: ic.img,
        to: `/category/${c.id}`,
        count: c.product_count,
      };
    }),
  ];
});

onMounted(async () => {
  site.ensureLoaded();
  const get = <T,>(url: string) => $fetch<T>(url).catch(() => null);
  const [s, promo, newest] = await Promise.all([
    get<NlStats>("/api/stats"),
    get<ListRes>("/api/products?promo=1&sort=sales&page_size=8"),
    get<ListRes>("/api/products?sort=newest&page_size=8"),
  ]);
  stats.value = s;
  if (promo?.products?.length) {
    flash.value = promo.products;
  } else {
    flashMode.value = "sales";
    flash.value = (await get<ListRes>("/api/products?sort=sales&page_size=8"))?.products || [];
  }
  picks.value = newest?.products || [];
  loading.value = false;
});
</script>
