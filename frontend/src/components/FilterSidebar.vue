<template>
  <div class="h-100 p-3 bg-dark overflow-auto" style="min-width: 280px; width: 280px; transition: all 0.3s">
    <div class="d-flex justify-content-between align-items-center mb-3">
        <h5 class="m-0">Filters</h5>
        <button type="button" class="btn-close btn-close-white d-md-none" aria-label="Close" @click="$emit('close')"></button>
    </div>

    <div class="d-flex flex-column gap-3">
        <!-- Search -->
        <div>
            <input
                :value="search"
                @input="$emit('update:search', $event.target.value)"
                placeholder="Search models..."
                class="form-control form-control-sm"
            />
        </div>

        <!-- Tags Search -->
        <div>
            <input
              :value="tagsSearch"
              @input="$emit('update:tagsSearch', $event.target.value)"
              placeholder="Search tags..."
              class="form-control form-control-sm"
            />
        </div>

        <!-- Category -->
         <div>
            <select
              :value="selectedCategory"
              @change="$emit('update:selectedCategory', $event.target.value)"
              class="form-select form-select-sm"
            >
              <option value="">All categories</option>
              <option v-for="cat in categories" :key="cat" :value="cat">
                {{ cat }}
              </option>
            </select>
        </div>
        
        <!-- Base Model -->
        <div>
             <select
              :value="selectedBaseModel"
              @change="$emit('update:selectedBaseModel', $event.target.value)"
              class="form-select form-select-sm"
            >
              <option value="">All base models</option>
              <option v-for="bm in baseModels" :key="bm" :value="bm">
                {{ bm }}
              </option>
            </select>
        </div>

        <!-- Model Type -->
        <div>
            <select
              :value="selectedModelType"
              @change="$emit('update:selectedModelType', $event.target.value)"
              class="form-select form-select-sm"
            >
              <option value="">All model types</option>
              <option v-for="t in modelTypes" :key="t" :value="t">
                {{ t }}
              </option>
            </select>
        </div>

        <!-- NSFW -->
        <div>
             <select
              :value="nsfwFilter"
              @change="$emit('update:nsfwFilter', $event.target.value)"
              class="form-select form-select-sm"
            >
              <option value="both">Both</option>
              <option value="no">No NSFW</option>
              <option value="only">Only NSFW</option>
            </select>
        </div>

        <!-- Synced Toggle -->
        <div class="form-check form-switch ps-0">
            <div class="d-flex justify-content-between align-items-center w-100">
                 <label class="form-check-label" for="syncedFilter">Synced</label>
                 <input 
                    class="form-check-input ms-2 mt-0" 
                    type="checkbox" 
                    role="switch" 
                    id="syncedFilter" 
                    :checked="syncedFilter"
                    @change="$emit('update:syncedFilter', $event.target.checked)"
                >
            </div>
        </div>

        <!-- Clear Button -->
         <button
          type="button"
          @click="$emit('clear')"
          class="btn btn-outline-secondary btn-sm w-100 mt-2"
        >
          <Icon icon="mdi:filter-remove-outline" class="me-1" />
          Clear Filters
        </button>

         <slot name="actions"></slot>
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
  syncedFilter: Boolean,
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
  "update:syncedFilter",
  "clear",
  "close",
]);
</script>
