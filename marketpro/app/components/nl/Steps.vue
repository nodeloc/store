<template>
  <div class="nl-steps" role="list" aria-label="订单进度">
    <div
      v-for="(s, i) in steps"
      :key="s.label"
      role="listitem"
      class="nl-step"
      :class="{ 'nl-step--done': s.state === 'done', 'nl-step--current': s.state === 'current' }"
      :style="{ '--i': i }"
      :aria-current="s.state === 'current' ? 'step' : undefined"
    >
      <span v-if="s.state === 'done'" class="nl-step__bar"></span>
      {{ s.label }}
      <small>{{ s.time || (s.state === "current" ? "进行中" : "—") }}</small>
    </div>
  </div>
</template>

<script setup lang="ts">
import { fmtDateTime } from "~/utils/nl";

const props = defineProps<{
  status: number;
  createdAt?: string;
  paidAt?: string | null;
  updatedAt?: string;
  autoDelivered?: boolean;
}>();

type State = "done" | "current" | "todo";

const steps = computed(() => {
  // 0 待支付 · 1 已支付待发货 · 4 已发货待确认 · 2 已完成
  const reached = { 0: 1, 1: 2, 4: 3, 2: 4 }[props.status as 0 | 1 | 4 | 2] ?? 1;
  const shippedAt = props.autoDelivered ? props.paidAt : props.status === 4 || props.status === 2 ? props.updatedAt : null;
  const items = [
    { label: "下单", time: fmtDateTime(props.createdAt) },
    { label: "已支付", time: fmtDateTime(props.paidAt) },
    { label: "已发货", time: fmtDateTime(shippedAt) },
    { label: "已完成", time: props.status === 2 ? fmtDateTime(props.updatedAt) : "" },
  ];
  return items.map((it, i) => {
    let state: State = "todo";
    if (i < reached) state = "done";
    else if (i === reached) state = "current";
    return { ...it, state };
  });
});
</script>
