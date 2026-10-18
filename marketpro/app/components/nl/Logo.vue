<template>
  <NuxtLink class="nl-logo" to="/" :aria-label="siteName">
    <img v-if="logo" :src="logo" :alt="siteName" style="max-height: 34px; width: auto; object-fit: contain" />
    <template v-else>
      <span class="nl-logo__mark"><i class="ph-fill ph-lightning"></i></span>
      <span class="nl-logo__word">NodeLoc<span>&nbsp;商店</span></span>
    </template>
  </NuxtLink>
</template>

<script setup lang="ts">
import { useSiteStore } from "~/stores/site";

const site = useSiteStore();
const logo = computed(() => (site.settings.site_logo as string) || "");
const siteName = computed(() => (site.settings.site_name as string) || "NodeLoc 商店");
</script>
