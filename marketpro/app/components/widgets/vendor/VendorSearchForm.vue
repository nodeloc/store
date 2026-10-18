<template>
  <div class="flex-between flex-wrap gap-8 mb-40">
    <form class="search-form__wrapper position-relative d-block" novalidate @submit="onSubmit">
      <input
        v-model="searchQuery"
        type="text"
        class="search-form__input common-input py-13 ps-16 pe-18 rounded-pill pe-44"
        placeholder="Search vendors by name or ID..."
        aria-label="Search vendors by name or ID"
        :aria-invalid="!!errorMessage"
        :aria-describedby="errorMessage ? 'search-error' : undefined"
      />
      <button
        type="submit"
        class="w-32 h-32 bg-main-600 rounded-circle flex-center text-xl text-white position-absolute top-50 translate-middle-y inset-inline-end-0 me-8"
        aria-label="Submit search"
      >
        <i class="ph ph-magnifying-glass"></i>
      </button>
    </form>
    <p v-if="errorMessage" id="search-error" class="text-danger mt-2 text-sm">
      {{ errorMessage }}
    </p>

    <div class="flex-align gap-16">
      <span
        class="text-neutral-600 fw-medium px-40 py-12 rounded-pill border border-neutral-100 d-md-block d-none"
      >
        Showing 1-20 of 85 results
      </span>

      <div class="flex-align gap-8">
        <label for="sort-select" class="text-gray-900 flex-shrink-0 fw-medium">Sort by:</label>
        <select
          id="sort-select"
          v-model="selectedSort"
          class="common-input form-select rounded-pill border border-gray-100 d-inline-block ps-20 pe-36 h-48 py-0 fw-medium"
          aria-label="Sort vendors"
        >
          <option v-for="option in sortOptions" :key="option.value" :value="option.value">
            {{ option.label }}
          </option>
        </select>
      </div>

      <button
        type="button"
        class="w-44 h-44 d-lg-none d-flex flex-center border border-gray-100 rounded-6 text-2xl sidebar-btn"
        aria-label="Toggle filter sidebar"
        @click="toggleSidebar"
      >
        <i class="ph-bold ph-funnel"></i>
      </button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { useSidebar } from "~/composables/useSidebar";

const searchQuery = ref("");
const errorMessage = ref("");

const sortOptions = [
  { value: "latest", label: "Latest" },
  { value: "old", label: "Old" },
];

const selectedSort = ref("latest");

function validateForm(): boolean {
  if (!searchQuery.value.trim()) {
    errorMessage.value = "Please enter a vendor name or ID to search.";
    return false;
  }
  errorMessage.value = "";
  return true;
}

function onSubmit(event: Event) {
  event.preventDefault();
  if (validateForm()) {
    // you could trigger an event or handle your logic here
  }
}

const { toggleSidebar } = useSidebar();
</script>
