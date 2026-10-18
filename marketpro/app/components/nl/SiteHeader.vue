<template>
  <header class="nl-header" :class="{ 'is-scrolled': scrolled }">
    <div class="nl-header__bar">
      <button type="button" class="nl-iconbtn nl-header__burger" aria-label="打开菜单" @click="drawer = true">
        <i class="ph ph-list"></i>
      </button>
      <NlLogo />
      <form class="nl-header__search" role="search" @submit.prevent="submit">
        <div class="nl-search">
          <label class="nl-search__cat">
            <span>{{ catLabel }}</span>
            <i class="ph ph-caret-down"></i>
            <select v-model="cat" aria-label="搜索分类">
              <option value="">全部分类</option>
              <option v-for="c in site.categories" :key="c.id" :value="String(c.id)">{{ c.name }}</option>
            </select>
          </label>
          <input
            ref="inputEl"
            v-model="q"
            class="nl-search__input"
            type="search"
            placeholder="搜索商品、分类或卖家…"
            aria-label="搜索"
            @input="err = false"
          />
          <kbd v-if="!q" class="nl-search__kbd">/</kbd>
          <button class="nl-search__go" type="submit" aria-label="搜索">
            <i class="ph-bold ph-magnifying-glass"></i>
          </button>
        </div>
        <span v-if="err" class="nl-search__err">请输入关键词</span>
      </form>
      <nav class="nl-header__actions" aria-label="账户">
        <button type="button" class="nl-iconbtn nl-header__searchbtn" aria-label="搜索" @click="mobileSearch = !mobileSearch">
          <i class="ph ph-magnifying-glass"></i>
        </button>
        <button
          type="button"
          class="nl-iconbtn"
          :aria-label="theme === 'dark' ? '切换到浅色主题' : '切换到深色主题'"
          :title="theme === 'dark' ? '浅色主题' : '深色主题'"
          @click="toggle"
        >
          <i :class="theme === 'dark' ? 'ph ph-sun' : 'ph ph-moon-stars'"></i>
        </button>
        <NuxtLink v-if="auth.isAdmin" class="nl-iconlink nl-iconlink--hide-sm" to="/admin">
          <i class="ph ph-gauge"></i><span>后台</span>
        </NuxtLink>
        <NuxtLink v-if="auth.isAuthenticated" class="nl-iconlink" to="/account">
          <i class="ph ph-user-circle"></i><span>{{ auth.user?.username }}</span>
        </NuxtLink>
        <a v-else class="nl-iconlink" :href="loginHref"><i class="ph ph-user-circle"></i><span>登录</span></a>
        <NuxtLink class="nl-iconlink nl-iconlink--hide-sm" to="/wishlist">
          <i class="ph ph-heart"></i>
          <span v-if="wishlist.count" :key="`w${wishlist.count}`" class="nl-iconlink__count">{{ wishlist.count }}</span>
          <span>收藏</span>
        </NuxtLink>
        <NuxtLink class="nl-iconlink" to="/cart">
          <i class="ph ph-shopping-cart-simple"></i>
          <span v-if="cart.count" :key="`c${cart.count}`" class="nl-iconlink__count">{{ cart.count }}</span>
          <span>购物车</span>
        </NuxtLink>
      </nav>
    </div>
    <form v-if="mobileSearch" class="nl-header__msearch" role="search" @submit.prevent="submit">
      <div class="nl-search">
        <input v-model="q" class="nl-search__input" type="search" placeholder="搜索商品…" aria-label="搜索" autofocus />
        <button class="nl-search__go" type="submit" aria-label="搜索"><i class="ph-bold ph-magnifying-glass"></i></button>
      </div>
    </form>
    <nav class="nl-header__nav" aria-label="主导航">
      <NuxtLink v-for="l in links" :key="l.key" :to="l.to" class="nl-navlink" :class="l.cls" :aria-current="l.active ? 'page' : undefined">
        <i v-if="l.icon" :class="l.icon"></i>{{ l.label }}
      </NuxtLink>
    </nav>
  </header>

  <div v-if="drawer" class="nl-drawer" role="dialog" aria-modal="true" aria-label="菜单" @keydown.esc="drawer = false">
    <div class="nl-drawer__scrim" @click="drawer = false"></div>
    <div class="nl-drawer__panel">
      <div class="nl-drawer__top">
        <NlLogo />
        <button type="button" class="nl-iconbtn" aria-label="关闭菜单" @click="drawer = false"><i class="ph ph-x"></i></button>
      </div>
      <nav class="nl-drawer__nav">
        <span class="nl-drawer__label">导航</span>
        <NuxtLink v-for="l in links" :key="l.key" :to="l.to" class="nl-navlink" :class="l.cls" :aria-current="l.active ? 'page' : undefined">
          <i v-if="l.icon" :class="l.icon"></i>{{ l.label }}
        </NuxtLink>
        <span class="nl-drawer__label">账户</span>
        <NuxtLink v-if="auth.isAuthenticated" class="nl-navlink" to="/account"><i class="ph ph-user-circle"></i>我的账户</NuxtLink>
        <a v-else class="nl-navlink" :href="loginHref"><i class="ph ph-sign-in"></i>登录</a>
        <NuxtLink class="nl-navlink" to="/account/orders"><i class="ph ph-receipt"></i>我的订单</NuxtLink>
        <NuxtLink class="nl-navlink" to="/wishlist"><i class="ph ph-heart"></i>收藏 ({{ wishlist.count }})</NuxtLink>
        <NuxtLink v-if="auth.isAdmin" class="nl-navlink" to="/admin"><i class="ph ph-gauge"></i>管理后台</NuxtLink>
      </nav>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "~/stores/auth";
import { useCartStore } from "~/stores/cart";
import { useSiteStore } from "~/stores/site";
import { useWishlistStore } from "~/stores/wishlist";

const auth = useAuthStore();
const cart = useCartStore();
const wishlist = useWishlistStore();
const site = useSiteStore();
const route = useRoute();
const { theme, toggle } = useNlTheme();

const q = ref(typeof route.query.keyword === "string" ? route.query.keyword : "");
const cat = ref(typeof route.query.category_id === "string" ? route.query.category_id : "");
const err = ref(false);
const drawer = ref(false);
const mobileSearch = ref(false);
const scrolled = ref(false);
const inputEl = ref<HTMLInputElement | null>(null);

const catLabel = computed(() => site.categories.find(c => String(c.id) === cat.value)?.name || "全部分类");
const loginHref = computed(() => `/auth/login?redirect=${encodeURIComponent(route.fullPath)}`);

const links = computed(() => {
  const p = route.path;
  const promo = p === "/shop" && route.query.promo === "1";
  const cats = site.categories.slice(0, 5).map(c => ({
    key: `c${c.id}`,
    label: c.name,
    to: `/category/${c.id}`,
    active: p === `/category/${c.id}`,
  }));
  return [
    { key: "all", label: "全部商品", to: "/shop", icon: "ph ph-squares-four", cls: "nl-navlink--all", active: p === "/shop" && !promo },
    { key: "home", label: "首页", to: "/", active: p === "/" },
    { key: "flash", label: "今日闪购", to: "/shop?promo=1", icon: "ph-fill ph-lightning", cls: "nl-navlink--hot", active: promo },
    ...cats,
    { key: "vendor", label: "店铺", to: "/vendor", active: p.startsWith("/vendor") },
    { key: "seller", label: "商家入驻", to: "/become-seller", active: p === "/become-seller" },
    { key: "blog", label: "博客", to: "/blog", active: p.startsWith("/blog") },
    { key: "contact", label: "联系我们", to: "/contact", active: p === "/contact" },
  ] as { key: string; label: string; to: string; icon?: string; cls?: string; active: boolean }[];
});

function submit() {
  const kw = q.value.trim();
  if (!kw) {
    err.value = true;
    return;
  }
  const params = new URLSearchParams({ keyword: kw });
  if (cat.value) params.set("category_id", cat.value);
  mobileSearch.value = false;
  navigateTo(`/shop?${params.toString()}`);
}

function onKey(e: KeyboardEvent) {
  if (e.key !== "/" || e.metaKey || e.ctrlKey || e.altKey) return;
  const t = e.target as HTMLElement | null;
  if (t && (t.isContentEditable || /^(INPUT|TEXTAREA|SELECT)$/.test(t.tagName))) return;
  e.preventDefault();
  inputEl.value?.focus();
}
function onScroll() {
  scrolled.value = window.scrollY > 4;
}

watch(() => route.fullPath, () => (drawer.value = false));

onMounted(() => {
  window.addEventListener("keydown", onKey);
  window.addEventListener("scroll", onScroll, { passive: true });
  onScroll();
});
onBeforeUnmount(() => {
  window.removeEventListener("keydown", onKey);
  window.removeEventListener("scroll", onScroll);
});
</script>
