<template>
  <nav v-if="totalPages > 1" class="mb-4">
    <ul class="pagination justify-content-center align-items-center gap-1">
      <li class="page-item" :class="{ disabled: page === 1 }">
        <a class="page-link" href="#" @click.prevent="$emit('changePage', 1)"
          >First</a
        >
      </li>
      <li class="page-item" :class="{ disabled: page === 1 }">
        <a
          class="page-link"
          href="#"
          @click.prevent="$emit('changePage', page - 1)"
          >Previous</a
        >
      </li>
      <li class="d-flex align-items-center">
        <input
          type="number"
          min="1"
          :max="totalPages"
          :value="page"
          @input="pageInput = Number($event.target.value)"
          @keyup.enter="onEnter"
          class="form-control"
          style="width: 80px"
        />
        <span class="ms-1">/ {{ totalPages }}</span>
      </li>
      <li class="page-item" :class="{ disabled: page === totalPages }">
        <a
          class="page-link"
          href="#"
          @click.prevent="$emit('changePage', page + 1)"
          >Next</a
        >
      </li>
      <li class="page-item" :class="{ disabled: page === totalPages }">
        <a
          class="page-link"
          href="#"
          @click.prevent="$emit('changePage', totalPages)"
          >Last</a
        >
      </li>
    </ul>
  </nav>
</template>

<script setup>
import { ref, watch } from "vue";

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
