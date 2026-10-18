<template>
  <footer class="nl-footer">
    <div class="nl-footer__inner">
      <div class="nl-footer__grid" :style="{ gridTemplateColumns: cols }">
        <div>
          <NlLogo />
          <p class="nl-footer__blurb">{{ cfg.description }}</p>
          <ul class="nl-footer__contact">
            <li v-if="cfg.email"><a :href="`mailto:${cfg.email}`"><i class="ph ph-envelope-simple"></i> {{ cfg.email }}</a></li>
            <li v-if="cfg.tel"><a :href="`tel:${cfg.tel}`"><i class="ph ph-phone"></i> {{ cfg.tel }}</a></li>
            <li v-if="cfg.address"><span><i class="ph ph-map-pin"></i> {{ cfg.address }}</span></li>
          </ul>
        </div>
        <div v-for="(s, i) in cfg.sections" :key="i">
          <h4>{{ s.title }}</h4>
          <ul>
            <li v-for="(l, j) in s.links" :key="j"><NuxtLink :to="l.url">{{ l.label }}</NuxtLink></li>
          </ul>
        </div>
      </div>
      <p class="nl-footer__giant" aria-hidden="true">{{ giant }}</p>
      <div class="nl-footer__base">
        <span>© {{ year }} {{ siteName }}</span>
        <span>以 NodeLoc 账号登录 · 积分支付</span>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { useSiteStore } from "~/stores/site";

interface FooterLink { label: string; url: string }
interface FooterSection { title: string; links: FooterLink[] }

const DEFAULT_SECTIONS: FooterSection[] = [
  {
    title: "购物",
    links: [
      { label: "全部商品", url: "/shop" },
      { label: "今日闪购", url: "/shop?promo=1" },
      { label: "店铺", url: "/vendor" },
    ],
  },
  {
    title: "账户",
    links: [
      { label: "我的订单", url: "/account/orders" },
      { label: "余额与提现", url: "/account/balance" },
      { label: "收藏", url: "/wishlist" },
    ],
  },
  {
    title: "社区",
    links: [
      { label: "商家入驻", url: "/become-seller" },
      { label: "博客", url: "/blog" },
      { label: "联系我们", url: "/contact" },
    ],
  },
];

const site = useSiteStore();
const year = new Date().getFullYear();
const str = (k: string) => (typeof site.settings[k] === "string" ? (site.settings[k] as string) : "");

const parsed = computed<Record<string, unknown>>(() => {
  try {
    return JSON.parse(str("footer_config") || "{}");
  } catch {
    return {};
  }
});

const siteName = computed(() => (parsed.value.site_name as string) || str("site_name") || "NodeLoc 商店");
const giant = computed(() => siteName.value.match(/[A-Za-z][A-Za-z0-9.]*/)?.[0] || "NodeLoc");

const cfg = computed(() => {
  const p = parsed.value;
  const sections = Array.isArray(p.sections) && p.sections.length ? (p.sections as FooterSection[]) : DEFAULT_SECTIONS;
  return {
    description:
      (p.description as string) || str("site_description") || "NodeLoc 社区发卡平台。使用社区账号登录，积分付款，卡密自动发放。",
    email: (p.email as string) || str("contact_email"),
    tel: (p.tel as string) || str("contact_tel"),
    address: (p.address as string) || str("contact_address"),
    sections: sections.slice(0, 4),
  };
});

const cols = computed(() => `1.6fr repeat(${cfg.value.sections.length}, 1fr)`);
</script>
