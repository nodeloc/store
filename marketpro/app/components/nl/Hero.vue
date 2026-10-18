<template>
  <section class="nl-hero">
    <div class="nl-hero__main">
      <span class="nl-hero__glow"></span><span class="nl-hero__glow nl-hero__glow--brand"></span>
      <svg class="nl-hero__art" viewBox="0 0 300 300" aria-hidden="true">
        <g class="nl-hero__ring">
          <circle cx="150" cy="150" r="146" fill="none" stroke="var(--volt)" stroke-width="2" stroke-dasharray="2 12" stroke-linecap="round" opacity=".8" />
        </g>
        <g class="nl-hero__ring nl-hero__ring--rev">
          <circle cx="150" cy="150" r="118" fill="none" stroke="var(--border-strong)" stroke-width="1" stroke-dasharray="40 14" />
        </g>
        <circle cx="150" cy="150" r="88" fill="var(--volt)" opacity=".1" />
        <path class="nl-hero__bolt" d="M164 58 100 168h48l-14 76 68-112h-50z" fill="var(--volt)" />
      </svg>
      <span class="nl-hero__eyebrow">
        <span class="nl-live"></span>
        <template v-if="stats">{{ stats.products.toLocaleString("en-US") }} 件商品在售 · 积分直接付款</template>
        <template v-else>NodeLoc 社区 · 积分直接付款</template>
      </span>
      <h1 class="nl-hero__title">
        <template v-for="(w, i) in line1" :key="`a${i}`"><span class="w" :style="{ '--i': i }">{{ w }}</span></template>
        <br />
        <em><template v-for="(w, i) in accent" :key="`b${i}`"><span class="w" :style="{ '--i': i }">{{ w }}</span></template></em>
        <template v-for="(w, i) in tail" :key="`c${i}`"><span class="w" :style="{ '--i': i }">{{ w }}</span></template>
      </h1>
      <p class="nl-hero__sub">{{ sub }}</p>
      <div class="nl-hero__cta">
        <NuxtLink class="nl-btn nl-btn--volt nl-btn--lg" to="/shop?promo=1">
          <i class="ph-fill ph-lightning"></i>逛逛闪购 <i class="ph ph-arrow-right"></i>
        </NuxtLink>
        <NuxtLink class="nl-btn nl-btn--outline-ink nl-btn--lg" to="/become-seller">成为卖家</NuxtLink>
      </div>
      <div class="nl-hero__stats">
        <div><b><NlCountUp :value="stats?.sold_keys || 0" /></b><span>累计发卡</span></div>
        <div><b><NlCountUp :value="stats?.paid_orders || 0" /></b><span>成交订单</span></div>
        <div><b><NlCountUp :value="stats?.shops || 0" /></b><span>入驻卖家</span></div>
      </div>
    </div>
    <div class="nl-hero__side">
      <NuxtLink class="nl-promo nl-promo--volt" to="/shop?promo=1" style="--i: 0">
        <span class="nl-promo__glyph">%</span>
        <div>
          <div class="nl-promo__meta">FLASH · 限时特惠</div>
          <h3 class="nl-promo__title">今日闪购</h3>
        </div>
        <span class="nl-promo__go"><span>去抢购</span><i class="ph-bold ph-arrow-right"></i></span>
      </NuxtLink>
      <NuxtLink class="nl-promo" :to="hot ? `/category/${hot.id}` : '/shop?sort=sales'" style="--i: 1">
        <span class="nl-promo__glyph" style="color: var(--brand)">{{ hotGlyph }}</span>
        <div>
          <div class="nl-promo__meta" style="color: var(--ink-muted)">HOT · 热门分类</div>
          <h3 class="nl-promo__title">{{ hot ? `${hot.name}专区` : "热销好物榜" }}</h3>
        </div>
        <span class="nl-promo__go" style="color: var(--ink)"><span>去看看</span><i class="ph-bold ph-arrow-right"></i></span>
      </NuxtLink>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useSiteStore } from "~/stores/site";
import type { NlCategory } from "~/utils/nl";

export interface NlStats {
  sold_keys: number;
  paid_orders: number;
  products: number;
  shops: number;
}

defineProps<{ stats: NlStats | null }>();

const site = useSiteStore();
const str = (k: string) => (typeof site.settings[k] === "string" ? (site.settings[k] as string).trim() : "");

const line1 = computed(() => Array.from(str("hero_title") || "社区好物，"));
const accent = computed(() => Array.from(str("hero_accent") || "秒发"));
const tail = computed(() => Array.from(str("hero_tail") || "到手"));
const sub = computed(
  () => str("hero_subtitle") || "VPS、AI 账号、点卡与兑换码。付款即自动发放卡密，订单页一键复制。"
);

const hot = computed(() => {
  const cats = site.categories as NlCategory[];
  return [...cats].sort((a, b) => (b.product_count || 0) - (a.product_count || 0))[0];
});
const hotGlyph = computed(() => {
  const n = hot.value?.name || "HOT";
  const latin = n.match(/[A-Za-z]{2,3}/);
  return latin ? latin[0].toUpperCase() : Array.from(n)[0];
});
</script>
