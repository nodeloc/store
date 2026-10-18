<template>
  <div class="nl-empty">
    <i :class="icon || 'ph ph-package'" aria-hidden="true"></i>
    <b>{{ title }}</b>
    <span v-if="text">{{ text }}</span>
    <slot />
  </div>
</template>

<script setup lang="ts">
defineProps<{ title: string; text?: string; icon?: string }>();
</script>
