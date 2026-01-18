<template>
  <div 
    class="card h-100 bg-dark-subtle border-0 shadow-sm transition-hover cursor-pointer"
    @click="$emit('click', collection)"
  >
    <div class="card-body d-flex flex-column">
      <div class="d-flex justify-content-between align-items-start mb-2">
        <h5 class="card-title text-truncate mb-0 fw-bold" :title="collection.name">{{ collection.name }}</h5>
        
        <div class="dropdown" @click.stop>
           <button 
                class="btn btn-link p-0 text-secondary z-2 position-relative" 
                type="button" 
                data-bs-toggle="dropdown"
                aria-expanded="false"
            >
             <Icon icon="mdi:dots-vertical" width="24" height="24" />
           </button>
           <ul class="dropdown-menu dropdown-menu-end shadow">
             <li><a class="dropdown-item" href="#" @click.prevent="$emit('rename', collection)">Rename</a></li>
             <li><hr class="dropdown-divider"></li>
             <li><a class="dropdown-item text-danger" href="#" @click.prevent="$emit('delete', collection)">Delete</a></li>
           </ul>
        </div>
      </div>
      
      <p class="card-text text-white-50 small flex-grow-1" style="min-height: 3em;">
        {{ collection.description || 'No description' }}
      </p>
      
      <div class="mt-auto pt-3 border-top border-secondary border-opacity-10 d-flex justify-content-between align-items-center">
         <span class="badge bg-secondary bg-opacity-25 text-white border border-secondary border-opacity-25 rounded-pill">
            {{ collection.versions ? collection.versions.length : 0 }} items
         </span>
         <span class="text-primary small d-flex align-items-center gap-1">
            View <Icon icon="mdi:arrow-right" />
         </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Icon } from "@iconify/vue";

defineProps({
  collection: {
    type: Object,
    required: true
  }
});

defineEmits(['click', 'rename', 'delete']);
</script>

<style scoped>
.transition-hover {
  transition: transform 0.2s, box-shadow 0.2s;
}
.transition-hover:hover {
  transform: translateY(-2px);
  box-shadow: 0 .5rem 1rem rgba(0,0,0,.15)!important;
}
.cursor-pointer {
    cursor: pointer;
}
</style>
