<template>
  <div>
    <Breadcrumb title="我的订单" :trail="[{ label: '我的账户', to: '/account' }]" />
    <div class="nl-container" style="padding-bottom: 96px">
      <div class="nl-account">
        <AccountSidebar />
        <div style="min-width: 0">
          <NlChips :items="tabs" :active="String(statusFilter)" @select="k => (statusFilter = Number(k))" />
          <div class="nl-toolbar">
            <span class="nl-toolbar__count">{{ loading ? "加载中…" : `共 ${filtered.length} 笔订单` }}</span>
            <div class="nl-toolbar__search">
              <input v-model="keyword" class="nl-input" type="search" placeholder="搜索订单号 / 商品名" aria-label="搜索订单" />
              <button type="button" aria-hidden="true" tabindex="-1"><i class="ph ph-magnifying-glass"></i></button>
            </div>
          </div>
          <div v-if="loading" class="nl-orders">
            <div v-for="i in 4" :key="i" class="nl-skel" style="min-height: 78px"></div>
          </div>
          <NlEmpty v-else-if="!filtered.length" icon="ph ph-receipt" title="暂无订单" :text="orders.length ? '换个筛选条件试试。' : '下单后会显示在这里。'">
            <NuxtLink v-if="!orders.length" class="nl-btn nl-btn--volt nl-btn--sm" to="/shop">去逛逛</NuxtLink>
          </NlEmpty>
          <div v-else class="nl-orders">
            <NlOrderRow v-for="o in filtered" :key="o.order_no" :order="o" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import AccountSidebar from "~/components/containers/account/AccountSidebar.vue";
import type { NlChipItem } from "~/components/nl/Chips.vue";
import type { NlOrderLite } from "~/components/nl/OrderRow.vue";
import { ORDER_STATUS } from "~/utils/nl";

definePageMeta({ layout: "layout-three" });
useHead({ title: "我的订单" });

const orders = ref<NlOrderLite[]>([]);
const loading = ref(true);
const statusFilter = ref(-1);
const keyword = ref("");

const tabs = computed<NlChipItem[]>(() => {
  const count = (s: number) => orders.value.filter(o => o.status === s).length;
  return [
    { key: "-1", label: "全部", icon: "ph ph-squares-four", count: orders.value.length },
    ...[0, 1, 4, 2, 3].map(s => ({ key: String(s), label: ORDER_STATUS[s]!.label, icon: ORDER_STATUS[s]!.icon, count: count(s) })),
  ];
});

const filtered = computed(() => {
  let list = orders.value;
  if (statusFilter.value >= 0) list = list.filter(o => o.status === statusFilter.value);
  const kw = keyword.value.trim().toLowerCase();
  if (kw) list = list.filter(o => o.order_no?.toLowerCase().includes(kw) || o.product?.name?.toLowerCase().includes(kw));
  return list;
});

onMounted(async () => {
  try {
    orders.value = await $fetch<NlOrderLite[]>("/api/orders", { credentials: "include" });
  } catch {
    orders.value = [];
  } finally {
    loading.value = false;
  }
});
</script>
