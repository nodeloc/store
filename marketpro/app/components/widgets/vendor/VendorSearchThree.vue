<template>
  <div class="d-flex align-items-center justify-content-between flex-wrap mt-40 mb-32 gap-16">
    <form class="input-group w-100 max-w-418" novalidate @submit="onSubmit">
      <input
        v-model="searchQuery"
        type="text"
        class="form-control common-input rounded-start-3"
        placeholder="Searching..."
        :class="{ 'is-invalid': errors.searchQuery }"
        :aria-invalid="errors.searchQuery ? 'true' : 'false'"
        aria-describedby="searchError"
      />
      <button
        type="submit"
        class="input-group-text border-0 bg-main-two-600 rounded-end-3 text-white text-2xl hover-bg-main-two-700 px-24"
      >
        <i class="ph ph-magnifying-glass"></i>
      </button>
      <div v-if="errors.searchQuery" id="searchError" class="invalid-feedback">
        {{ errors.searchQuery }}
      </div>
    </form>

    <div
      class="d-flex align-items-center justify-content-between justify-content-sm-end gap-16 flex-grow-1"
    >
      <div class="text-gray-600 text-md flex-shrink-0">
        <span class="text-neutral-900 fw-semibold">52</span> Results Found
      </div>

      <div class="d-flex align-items-center gap-8 d-sm-flex d-none">
        <button
          type="button"
          :class="[
            'grid-btn text-2xl d-flex w-48 h-48 border rounded-8 justify-content-center align-items-center',
            isListView ? 'border-gray-100' : 'border-main-600 text-white bg-main-600',
          ]"
          :aria-pressed="!isListView"
          aria-label="Grid view"
          @click="toggleGridView"
        >
          <i class="ph ph-squares-four"></i>
        </button>
        <button
          type="button"
          :class="[
            'list-btn text-2xl d-flex w-48 h-48 border rounded-8 justify-content-center align-items-center',
            isListView ? 'border-main-600 text-white bg-main-600' : 'border-gray-100',
          ]"
          :aria-pressed="isListView"
          aria-label="List view"
          @click="toggleListView"
        >
          <i class="ph ph-list-bullets"></i>
        </button>
      </div>

      <button
        type="button"
        class="w-48 h-48 d-lg-none d-flex flex-center border border-gray-100 rounded-6 text-2xl sidebar-btn"
        aria-label="Toggle filters"
        @click="toggleSidebar"
      >
        <i class="ph-bold ph-funnel"></i>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useSidebar } from "~/composables/useSidebar";
import { useViewToggle } from "~/composables/useViewToggle";

const searchQuery = ref("");
const errors = ref<{ searchQuery?: string }>({});

const { isListView, toggleListView, toggleGridView } = useViewToggle();

function validateForm() {
  errors.value = {};

  if (!searchQuery.value.trim()) {
    errors.value.searchQuery = "Please enter a search term.";
    return false;
  }

  return true;
}

function onSubmit(event: Event) {
  event.preventDefault();

  if (validateForm()) {
    // No API - you can implement search/filter logic here
  }
}

const { toggleSidebar } = useSidebar();
</script>
