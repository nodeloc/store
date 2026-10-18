<template>
  <NuxtLink class="nl-orderrow" :to="`/order/${order.order_no}`">
    <span class="nl-order__thumb">
      <img v-if="order.product?.image" :src="order.product.image" alt="" />
      <i v-else :class="categoryIcon(order.product?.category).cls"></i>
    </span>
    <div style="min-width: 0">
      <div class="nl-orderrow__name">{{ order.product?.name || "-" }}</div>
      <div class="nl-orderrow__no">{{ order.order_no }} · × {{ order.quantity || 1 }}</div>
    </div>
    <span class="nl-orderrow__time">{{ fmtDateTime(order.created_at) }}</span>
    <NlPrice :now="order.total_amount || 0" />
    <span class="nl-badge" :class="st.tone && `nl-badge--${st.tone}`"><i :class="st.icon"></i>{{ st.label }}</span>
  </NuxtLink>
</template>

<script setup lang="ts">
import { ORDER_STATUS, categoryIcon, fmtDateTime, type NlProduct } from "~/utils/nl";

export interface NlOrderLite {
  order_no: string;
  status: number;
  quantity?: number;
  total_amount?: number;
  created_at?: string;
  product?: NlProduct;
}

const props = defineProps<{ order: NlOrderLite }>();
const st = computed(() => ORDER_STATUS[props.order.status] || { label: "未知", tone: "", icon: "ph ph-question" });
</script>
