<template>
  <span
    v-if="left > 0"
    class="nl-countdown"
    :class="{ 'nl-countdown--inline': inline }"
    role="timer"
    :aria-label="`距结束 ${parts[0]} 时 ${parts[1]} 分 ${parts[2]} 秒`"
  >
    <template v-for="(v, i) in parts" :key="i">
      <span v-if="i > 0" class="nl-countdown__sep" aria-hidden="true">:</span>
      <span class="nl-countdown__unit" aria-hidden="true">
        <span class="nl-countdown__num"><span :key="v" :class="{ flip: animate }">{{ v }}</span></span>
        <span v-if="!inline" class="nl-countdown__lbl">{{ labels[i] }}</span>
      </span>
    </template>
  </span>
</template>

<script setup lang="ts">
import { prefersReducedMotion } from "~/utils/nl";

const props = defineProps<{ until?: string | null; seconds?: number; inline?: boolean }>();
const emit = defineEmits<{ end: [] }>();

const labels = ["时", "分", "秒"];
const animate = !prefersReducedMotion();

function compute() {
  if (props.until) return Math.max(0, Math.floor((new Date(props.until).getTime() - Date.now()) / 1000));
  return Math.max(0, Math.floor(props.seconds || 0));
}

const left = ref(compute());
const pad = (n: number) => String(n).padStart(2, "0");
const parts = computed(() => [
  pad(Math.floor(left.value / 3600)),
  pad(Math.floor((left.value % 3600) / 60)),
  pad(left.value % 60),
]);

let timer: ReturnType<typeof setInterval> | null = null;
let started = 0;
let base = 0;

function start() {
  if (timer) clearInterval(timer);
  base = compute();
  started = Date.now();
  left.value = base;
  timer = setInterval(() => {
    left.value = props.until ? compute() : Math.max(0, base - Math.floor((Date.now() - started) / 1000));
    if (left.value <= 0 && timer) {
      clearInterval(timer);
      timer = null;
      emit("end");
    }
  }, 1000);
}

onMounted(start);
watch(() => [props.until, props.seconds], start);
onBeforeUnmount(() => timer && clearInterval(timer));
</script>
