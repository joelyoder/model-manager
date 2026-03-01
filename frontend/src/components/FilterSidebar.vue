<template>
  <div class="h-100 p-3 bg-dark overflow-auto" style="min-width: 280px; width: 280px; transition: all 0.3s">
    <div class="d-flex justify-content-between align-items-center mb-4">
        <h5 class="m-0 text-white fw-bold">Filters</h5>
        <button type="button" class="btn-close btn-close-white d-md-none" aria-label="Close" @click="$emit('close')"></button>
    </div>

    <div class="d-flex flex-column gap-3">
        <!-- Search -->
        <div class="input-group">
            <span class="input-group-text bg-dark-subtle border-0 text-secondary">
                <Icon icon="mdi:magnify" width="18" height="18" />
            </span>
            <input
                :value="search"
                @input="$emit('update:search', $event.target.value)"
                placeholder="Search models..."
                class="form-control bg-dark-subtle border-0 text-white shadow-none"
            />
        </div>

        <!-- Tags Search -->
        <div class="input-group">
            <span class="input-group-text bg-dark-subtle border-0 text-secondary">
                <Icon icon="mdi:tag-outline" width="18" height="18" />
            </span>
            <input
              :value="tagsSearch"
              @input="$emit('update:tagsSearch', $event.target.value)"
              placeholder="Search tags..."
              class="form-control bg-dark-subtle border-0 text-white shadow-none"
            />
        </div>

        <!-- Category -->
         <div>
            <select
              :value="selectedCategory"
              @change="$emit('update:selectedCategory', $event.target.value)"
              class="form-select bg-dark-subtle border-0 text-white shadow-none"
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
              class="form-select bg-dark-subtle border-0 text-white shadow-none"
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
              class="form-select bg-dark-subtle border-0 text-white shadow-none"
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
              class="form-select bg-dark-subtle border-0 text-white shadow-none"
            >
              <option value="both">NSFW & Safe</option>
              <option value="no">Safe Only</option>
              <option value="only">NSFW Only</option>
            </select>
        </div>

        <!-- Synced Toggle -->
        <div class="form-check form-switch ps-0 mt-2 pt-2 border-top border-dark-subtle">
            <div class="d-flex justify-content-between align-items-center w-100">
                 <label class="form-check-label text-secondary" for="syncedFilter">Installed Only</label>
                 <input 
                    class="form-check-input ms-2 mt-0" 
                    type="checkbox" 
                    role="switch" 
                    id="syncedFilter" 
                    :checked="syncedFilter"
                    @change="$emit('update:syncedFilter', $event.target.checked)"
                    style="cursor: pointer;"
                >
            </div>
        </div>

        <!-- Uncollected Toggle -->
        <div class="form-check form-switch ps-0">
            <div class="d-flex justify-content-between align-items-center w-100">
                 <label class="form-check-label text-secondary" for="uncollectedFilter">Not in a Collection</label>
                 <input 
                    class="form-check-input ms-2 mt-0" 
                    type="checkbox" 
                    role="switch" 
                    id="uncollectedFilter" 
                    :checked="uncollectedFilter"
                    @change="$emit('update:uncollectedFilter', $event.target.checked)"
                    style="cursor: pointer;"
                >
            </div>
        </div>

        <!-- Clear Button -->
         <button
          type="button"
          @click="$emit('clear')"
          class="btn btn-dark bg-opacity-25 text-secondary-emphasis w-100 mt-2 border-0"
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
  uncollectedFilter: Boolean,
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
  "update:uncollectedFilter",
  "clear",
  "close",
]);
</script>
