<template>
  <aside class="nl-side">
    <div class="nl-side__user">
      <span class="nl-avatar">
        <img v-if="auth.user?.avatar_url" :src="auth.user.avatar_url" :alt="auth.user?.username" />
        <i v-else class="ph ph-user"></i>
      </span>
      <div style="min-width: 0">
        <b>{{ auth.user?.username || "用户" }}</b>
        <span>余额 {{ money(auth.user?.balance) }}</span>
      </div>
    </div>
    <nav aria-label="账户导航">
      <NuxtLink v-for="item in navItems" :key="item.to" :to="item.to">
        <i :class="item.icon"></i><span>{{ item.label }}</span>
      </NuxtLink>
      <NuxtLink to="/wishlist"><i class="ph ph-heart"></i><span>我的收藏</span></NuxtLink>
      <button type="button" @click="auth.logout()"><i class="ph ph-sign-out"></i><span>退出登录</span></button>
    </nav>
  </aside>
</template>

<script setup lang="ts">
import { useAuthStore } from "~/stores/auth";
import { money } from "~/utils/nl";

const auth = useAuthStore();

const navItems = [
  { to: "/account", label: "账户概览", icon: "ph ph-squares-four" },
  { to: "/account/orders", label: "我的订单", icon: "ph ph-receipt" },
  { to: "/account/balance", label: "余额管理", icon: "ph ph-wallet" },
  { to: "/account/shop", label: "我的店铺", icon: "ph ph-storefront" },
  { to: "/account/withdrawals", label: "提现管理", icon: "ph ph-hand-coins" },
];
</script>
