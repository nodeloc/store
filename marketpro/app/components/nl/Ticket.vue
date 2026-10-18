<template>
  <div class="nl-ticket" :style="{ '--i': index }">
    <div class="nl-ticket__main">
      <div v-for="(row, i) in rows" :key="i" class="nl-ticket__row">
        <span class="nl-ticket__label">{{ row.label }}</span>
        <code class="nl-ticket__value"><NlScramble :text="row.value" :delay="index * 180 + i * 90" /></code>
        <button
          type="button"
          class="nl-ticket__copy"
          :data-copied="copied === i ? '' : undefined"
          :aria-label="`复制${row.label}`"
          @click="copy(row.value, i)"
        >
          <template v-if="copied === i"><i class="ph ph-check"></i>已复制</template>
          <template v-else><i class="ph ph-copy"></i>复制</template>
        </button>
      </div>
    </div>
    <div class="nl-ticket__stub" aria-hidden="true">
      <b>{{ String(index + 1).padStart(2, "0") }}</b>
      <span>KEY</span>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ index: number; rows: { label: string; value: string }[] }>();
const emit = defineEmits<{ copied: [label: string] }>();

const copied = ref<number | null>(null);
let timer: ReturnType<typeof setTimeout> | null = null;

async function copy(text: string, i: number) {
  try {
    await navigator.clipboard.writeText(text);
  } catch {
    const ta = document.createElement("textarea");
    ta.value = text;
    ta.style.position = "fixed";
    ta.style.opacity = "0";
    document.body.appendChild(ta);
    ta.select();
    document.execCommand("copy");
    ta.remove();
  }
  copied.value = i;
  emit("copied", props.rows[i]?.label || "");
  if (timer) clearTimeout(timer);
  timer = setTimeout(() => (copied.value = null), 2000);
}

onBeforeUnmount(() => timer && clearTimeout(timer));
</script>
