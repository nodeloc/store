<template>
  <div v-reveal class="nl-sechead" data-reveal>
    <div>
      <span class="nl-sechead__eyebrow"><span v-if="live" class="nl-live"></span>{{ eyebrow }}</span>
      <h2 class="nl-sechead__title">{{ title }}</h2>
    </div>
    <div v-if="$slots.default" class="nl-sechead__aside"><slot /></div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ eyebrow: string; title: string; live?: boolean }>();
</script>
