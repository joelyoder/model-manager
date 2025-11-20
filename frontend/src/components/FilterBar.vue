<template>
  <div>
    <div class="row gap-2">
      <div class="col">
        <input
          :value="search"
          @input="$emit('update:search', $event.target.value)"
          placeholder="Search models..."
          class="form-control w-200 flex-grow-1"
          style="min-width: 200px"
        />
      </div>
      <div class="col">
        <input
          :value="tagsSearch"
          @input="$emit('update:tagsSearch', $event.target.value)"
          placeholder="Search tags (comma separated)"
          class="form-control"
          style="min-width: 200px"
        />
      </div>
    </div>
    <div class="row gap-2 my-2">
      <div class="col">
        <select
          :value="selectedCategory"
          @change="$emit('update:selectedCategory', $event.target.value)"
          class="form-select"
          style="min-width: 250px"
        >
          <option value="">All categories</option>
          <option v-for="cat in categories" :key="cat" :value="cat">
            {{ cat }}
          </option>
        </select>
      </div>
      <div class="col">
        <select
          :value="selectedBaseModel"
          @change="$emit('update:selectedBaseModel', $event.target.value)"
          class="form-select"
          style="min-width: 250px"
        >
          <option value="">All base models</option>
          <option v-for="bm in baseModels" :key="bm" :value="bm">
            {{ bm }}
          </option>
        </select>
      </div>
      <div class="col">
        <select
          :value="selectedModelType"
          @change="$emit('update:selectedModelType', $event.target.value)"
          class="form-select"
          style="min-width: 250px"
        >
          <option value="">All model types</option>
          <option v-for="t in modelTypes" :key="t" :value="t">
            {{ t }}
          </option>
        </select>
      </div>
      <div class="col">
        <select
          :value="nsfwFilter"
          @change="$emit('update:nsfwFilter', $event.target.value)"
          class="form-select"
          style="min-width: 200px"
        >
          <option value="no">No NSFW</option>
          <option value="only">Only NSFW</option>
          <option value="both">Both</option>
        </select>
      </div>
      <div class="col-auto d-flex align-items-center">
        <button
          type="button"
          @click="$emit('clear')"
          class="btn btn-outline-secondary d-inline-flex align-items-center justify-content-center"
          aria-label="Clear filters"
          title="Clear filters"
        >
          <Icon icon="mdi:filter-remove-outline" width="20" height="20" />
          <span class="visually-hidden">Clear Filters</span>
        </button>
      </div>
      <div class="col d-flex justify-content-end">
        <slot name="actions"></slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Icon } from "@iconify/vue";

defineProps({
  search: String,
  tagsSearch: String,
  selectedCategory: String,
  selectedBaseModel: String,
  selectedModelType: String,
  nsfwFilter: String,
  categories: Array,
  baseModels: Array,
  modelTypes: Array,
});

defineEmits([
  "update:search",
  "update:tagsSearch",
  "update:selectedCategory",
  "update:selectedBaseModel",
  "update:selectedModelType",
  "update:nsfwFilter",
  "clear",
]);
</script>
