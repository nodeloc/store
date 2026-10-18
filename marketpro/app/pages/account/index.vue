<template>
  <div>
    <Breadcrumb title="我的账户" />
    <div class="nl-container" style="padding-bottom: 96px">
      <div class="nl-account">
        <AccountSidebar />
        <div class="nl-stack" style="gap: 16px">
          <div class="nl-welcome">
            <h2>欢迎回来，<em>{{ auth.user?.username }}</em></h2>
            <p>在这里查看订单、管理余额，或者开一家自己的店铺。</p>
            <div class="nl-row">
              <NuxtLink class="nl-btn nl-btn--volt" to="/account/orders"><i class="ph ph-receipt"></i>我的订单</NuxtLink>
              <NuxtLink class="nl-btn nl-btn--outline-ink" to="/shop">继续购物 <i class="ph ph-arrow-right"></i></NuxtLink>
            </div>
          </div>
          <div class="nl-bento">
            <div class="nl-bento__cell">
              <i class="ph-fill ph-wallet"></i><b class="nl-mono">{{ money(auth.user?.balance) }}</b><span>账户余额</span>
            </div>
            <div class="nl-bento__cell">
              <i class="ph-fill ph-receipt"></i><b>{{ orders === null ? "—" : orders.length }}</b><span>全部订单</span>
            </div>
            <div class="nl-bento__cell">
              <i class="ph-fill ph-clock"></i><b>{{ orders === null ? "—" : pending }}</b><span>待处理订单</span>
            </div>
            <NuxtLink class="nl-bento__cell" to="/account/shop" style="text-decoration: none; color: inherit">
              <i class="ph-fill ph-storefront"></i><b>我的店铺</b><span>上架商品，开始售卖</span>
            </NuxtLink>
          </div>
          <div v-if="recent.length" class="nl-panel">
            <h2 class="nl-panel__title">最近订单</h2>
            <div class="nl-orders">
              <NlOrderRow v-for="o in recent" :key="o.order_no" :order="o" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import AccountSidebar from "~/components/containers/account/AccountSidebar.vue";
import type { NlOrderLite } from "~/components/nl/OrderRow.vue";
import { useAuthStore } from "~/stores/auth";
import { money } from "~/utils/nl";

definePageMeta({ layout: "layout-three" });
useHead({ title: "我的账户" });

const auth = useAuthStore();
const orders = ref<NlOrderLite[] | null>(null);
const pending = computed(() => (orders.value || []).filter(o => o.status === 0 || o.status === 1 || o.status === 4).length);
const recent = computed(() => (orders.value || []).slice(0, 5));

onMounted(async () => {
  try {
    orders.value = await $fetch<NlOrderLite[]>("/api/orders", { credentials: "include" });
  } catch {
    orders.value = [];
  }
});
</script>
