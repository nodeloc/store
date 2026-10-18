<template>
  <div class="nl-container nl-pagehead">
    <ol class="nl-crumbs" aria-label="面包屑">
      <li><NuxtLink to="/">首页</NuxtLink></li>
      <template v-for="c in trail" :key="c.label">
        <li aria-hidden="true">/</li>
        <li><NuxtLink v-if="c.to" :to="c.to">{{ c.label }}</NuxtLink><span v-else>{{ c.label }}</span></li>
      </template>
      <li aria-hidden="true">/</li>
      <li aria-current="page">{{ title }}</li>
    </ol>
    <template v-if="!compact">
      <h1 class="nl-pagehead__title">{{ title }}</h1>
      <p v-if="sub" class="nl-pagehead__sub">{{ sub }}</p>
    </template>
    <slot />
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ title: string; sub?: string; compact?: boolean; trail?: { label: string; to?: string }[] }>(), {
  sub: "",
  compact: false,
  trail: () => [],
});
</script>
