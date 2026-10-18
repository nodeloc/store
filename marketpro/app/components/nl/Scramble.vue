<template>
  <span :class="{ 'is-scrambling': running }" :aria-label="text">{{ display }}</span>
</template>

<script setup lang="ts">
import { prefersReducedMotion } from "~/utils/nl";

const props = withDefaults(defineProps<{ text: string; delay?: number }>(), { delay: 0 });

const GLYPHS = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789#%&*";
const display = ref(props.text);
const running = ref(false);
let raf = 0;

function play() {
  cancelAnimationFrame(raf);
  const final = props.text;
  display.value = final;
  if (prefersReducedMotion() || !final) return;
  const dur = 900 + final.length * 22;
  let start: number | null = null;
  running.value = true;
  const frame = (ts: number) => {
    if (start === null) start = ts + props.delay;
    const p = Math.max(0, Math.min(1, (ts - start) / dur));
    let out = "";
    for (let i = 0; i < final.length; i++) {
      const c = final[i]!;
      out += c === "-" || c === " " || i < p * final.length ? c : GLYPHS[(Math.random() * GLYPHS.length) | 0];
    }
    display.value = out;
    if (p < 1) raf = requestAnimationFrame(frame);
    else {
      display.value = final;
      running.value = false;
    }
  };
  raf = requestAnimationFrame(frame);
}

onMounted(play);
watch(() => props.text, play);
onBeforeUnmount(() => cancelAnimationFrame(raf));
</script>
