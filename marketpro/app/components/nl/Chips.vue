<template>
  <div class="nl-chips-wrap">
    <button
      v-if="canLeft"
      type="button"
      class="nl-chips__arrow nl-chips__arrow--prev"
      aria-label="向左滚动"
      @click="scrollBy(-1)"
    >
      <i class="ph ph-caret-left"></i>
    </button>
    <div ref="scroller" class="nl-chips" :class="{ 'nl-chips--fade-l': canLeft, 'nl-chips--fade-r': canRight }" role="list" @scroll.passive="update">
      <template v-for="c in items" :key="c.key">
        <NuxtLink
          v-if="c.to"
          role="listitem"
          class="nl-chip"
          :class="c.className"
          :to="c.to"
          :aria-current="c.key === active ? 'true' : undefined"
        >
          <span class="nl-chip__icon">
            <img v-if="c.img" :src="c.img" alt="" style="width: 32px; height: 32px; border-radius: 50%; object-fit: cover" />
            <i v-else :class="c.icon"></i>
          </span>
          {{ c.label }}
          <span v-if="c.count !== undefined" class="nl-chip__count">{{ c.count.toLocaleString("en-US") }}</span>
        </NuxtLink>
        <button
          v-else
          type="button"
          role="listitem"
          class="nl-chip"
          :class="c.className"
          :aria-current="c.key === active ? 'true' : undefined"
          @click="$emit('select', c.key)"
        >
          <span class="nl-chip__icon">
            <img v-if="c.img" :src="c.img" alt="" style="width: 32px; height: 32px; border-radius: 50%; object-fit: cover" />
            <i v-else :class="c.icon"></i>
          </span>
          {{ c.label }}
          <span v-if="c.count !== undefined" class="nl-chip__count">{{ c.count.toLocaleString("en-US") }}</span>
        </button>
      </template>
    </div>
    <button
      v-if="canRight"
      type="button"
      class="nl-chips__arrow nl-chips__arrow--next"
      aria-label="向右滚动"
      @click="scrollBy(1)"
    >
      <i class="ph ph-caret-right"></i>
    </button>
  </div>
</template>

<script setup lang="ts">
export interface NlChipItem {
  key: string;
  label: string;
  icon?: string;
  img?: string;
  count?: number;
  to?: string;
  className?: string;
}

const props = defineProps<{ items: NlChipItem[]; active?: string }>();
defineEmits<{ select: [key: string] }>();

const scroller = ref<HTMLElement | null>(null);
const canLeft = ref(false);
const canRight = ref(false);
let ro: ResizeObserver | null = null;

function update() {
  const el = scroller.value;
  if (!el) return;
  canLeft.value = el.scrollLeft > 2;
  canRight.value = el.scrollLeft + el.clientWidth < el.scrollWidth - 2;
}

function scrollBy(dir: number) {
  const el = scroller.value;
  if (!el) return;
  el.scrollBy({ left: dir * el.clientWidth * 0.7, behavior: "smooth" });
}

watch(() => props.items.length, () => nextTick(update));

onMounted(() => {
  update();
  if (scroller.value) {
    ro = new ResizeObserver(update);
    ro.observe(scroller.value);
  }
  window.addEventListener("resize", update);
});

onBeforeUnmount(() => {
  ro?.disconnect();
  window.removeEventListener("resize", update);
});
</script>
