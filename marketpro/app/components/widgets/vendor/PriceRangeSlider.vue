<template>
  <div class="border border-gray-50 rounded-8 p-24">
    <h6 class="text-xl border-bottom border-gray-100 pb-24 mb-24">Filter by Price</h6>
    <div class="custom--range">
      <Range v-model="values" :min="min" :max="max" :step="step" class="slider-range-container" />
      <div class="flex-between flex-wrap-reverse gap-8 mt-24">
        <button type="button" class="btn btn-main h-40 flex-align" @click="onFilterClick">
          Filter
        </button>
        <div class="custom--range__content flex-align gap-8">
          <span class="text-gray-500 text-md flex-shrink-0">Price:</span>
          <input
            id="amount"
            type="text"
            class="custom--range__prices text-neutral-600 text-start text-md fw-medium"
            :value="formattedValue"
            readonly
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from "vue";
import { Range } from "vue-range-multi";
import "vue-range-multi/style.css";

const props = defineProps({
  min: {
    type: Number,
    required: true,
  },
  max: {
    type: Number,
    required: true,
  },
  step: {
    type: Number,
    default: 1,
  },
  initialLow: {
    type: Number,
    default: undefined,
  },
  initialHigh: {
    type: Number,
    default: undefined,
  },
});

const emits = defineEmits<{
  (e: "update" | "filter", range: { low: number; high: number }): void;
}>();

const values = ref<[number, number]>([
  props.initialLow ?? props.min,
  props.initialHigh ?? props.max,
]);

watch(values, ([low, high]) => {
  emits("update", { low, high });
});

const formattedValue = computed(() => {
  return `$${values.value[0]} - $${values.value[1]}`;
});

function onFilterClick() {
  emits("filter", { low: values.value[0], high: values.value[1] });
}
</script>

<style scoped>
.slider-range-container {
  margin-bottom: 1rem;
}
</style>
