<template>
  <span ref="el">{{ shown }}{{ suffix }}</span>
</template>

<script setup lang="ts">
import { prefersReducedMotion } from "~/utils/nl";

const props = withDefaults(defineProps<{ value: number; decimals?: number; suffix?: string }>(), {
  decimals: 0,
  suffix: "",
});

const el = ref<HTMLElement | null>(null);
const current = ref(0);
const fmt = (n: number) =>
  n.toLocaleString("en-US", { minimumFractionDigits: props.decimals, maximumFractionDigits: props.decimals });
const shown = computed(() => fmt(current.value));

let raf = 0;
function run(to: number) {
  cancelAnimationFrame(raf);
  if (prefersReducedMotion()) {
    current.value = to;
    return;
  }
  const from = current.value;
  let t0: number | null = null;
  const step = (ts: number) => {
    if (t0 === null) t0 = ts;
    const p = Math.min(1, (ts - t0) / 1400);
    const e = 1 - Math.pow(1 - p, 4);
    current.value = from + (to - from) * e;
    if (p < 1) raf = requestAnimationFrame(step);
  };
  raf = requestAnimationFrame(step);
}

let io: IntersectionObserver | null = null;
let seen = false;
onMounted(() => {
  if (!el.value || !("IntersectionObserver" in window)) {
    seen = true;
    run(props.value);
    return;
  }
  io = new IntersectionObserver(entries => {
    if (entries.some(e => e.isIntersecting)) {
      seen = true;
      run(props.value);
      io?.disconnect();
    }
  });
  io.observe(el.value);
});
watch(
  () => props.value,
  v => seen && run(v)
);
onBeforeUnmount(() => {
  io?.disconnect();
  cancelAnimationFrame(raf);
});
</script>
