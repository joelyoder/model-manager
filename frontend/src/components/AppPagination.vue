<template>
  <nav v-if="totalPages > 1" class="mb-4 d-flex justify-content-center align-items-center gap-2">
    <button 
        class="btn btn-outline-secondary btn-sm d-flex align-items-center justify-content-center" 
        :disabled="page === 1" 
        @click="$emit('changePage', 1)"
        title="First Page"
        style="width: 32px; height: 32px;"
    >
        <Icon icon="mdi:chevron-double-left" width="20" height="20" />
    </button>
    <button 
        class="btn btn-outline-secondary btn-sm d-flex align-items-center justify-content-center" 
        :disabled="page === 1" 
        @click="$emit('changePage', page - 1)"
        title="Previous Page"
        style="width: 32px; height: 32px;"
    >
        <Icon icon="mdi:chevron-left" width="20" height="20" />
    </button>
    
    <div class="d-flex align-items-center bg-dark-subtle rounded px-2 border border-secondary border-opacity-25" style="height: 32px;">
        <input
            type="number"
            min="1"
            :max="totalPages"
            :value="page"
            @input="pageInput = Number($event.target.value)"
            @keyup.enter="onEnter"
            class="form-control form-control-sm bg-transparent border-0 text-white shadow-none text-center p-0"
            style="width: 40px"
        />
        <span class="text-secondary small ms-1 user-select-none">/ {{ totalPages }}</span>
    </div>

    <button 
        class="btn btn-outline-secondary btn-sm d-flex align-items-center justify-content-center" 
        :disabled="page === totalPages" 
        @click="$emit('changePage', page + 1)"
        title="Next Page"
        style="width: 32px; height: 32px;"
    >
         <Icon icon="mdi:chevron-right" width="20" height="20" />
    </button>
    <button 
        class="btn btn-outline-secondary btn-sm d-flex align-items-center justify-content-center" 
        :disabled="page === totalPages" 
        @click="$emit('changePage', totalPages)"
        title="Last Page"
        style="width: 32px; height: 32px;"
    >
        <Icon icon="mdi:chevron-double-right" width="20" height="20" />
    </button>
  </nav>
</template>

<script setup>
import { ref, watch } from "vue";
import { Icon } from "@iconify/vue";

const props = defineProps({
  page: Number,
  totalPages: Number,
});

const emit = defineEmits(["changePage"]);

const pageInput = ref(props.page);

watch(
  () => props.page,
  (val) => {
    pageInput.value = val;
  }
);



const onEnter = () => {
  let p = pageInput.value;
  if (p < 1) p = 1;
  if (p > props.totalPages) p = props.totalPages;
  emit("changePage", p);
};
</script>
