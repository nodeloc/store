<template>
  <div class="nl-site nl-canvas nl-shell">
    <NlSiteHeader />
    <NlLiveTicker v-if="ticker" />
    <main>
      <slot />
    </main>
    <NlSiteFooter />
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ ticker?: boolean }>(), { ticker: false });

const { apply } = useNlTheme();
onMounted(() => apply(true));
</script>
