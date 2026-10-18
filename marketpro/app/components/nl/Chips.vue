<template>
  <div class="nl-chips" role="list">
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

defineProps<{ items: NlChipItem[]; active?: string }>();
defineEmits<{ select: [key: string] }>();
</script>
