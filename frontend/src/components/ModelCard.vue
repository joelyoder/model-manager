<template>
  <div :id="`model-${version.ID}`" class="model-card card h-100">
    <div class="position-relative">
      <img
        v-if="imageUrl"
        :src="imageUrl"
        class="card-img-top"
        style="width: 100%; height: 450px; object-fit: cover"
      />
      <div class="card-img-overlay z-2">
        <span class="badge rounded-pill text-bg-primary">{{ version.type }}</span>
        <span class="ms-1 badge rounded-pill text-bg-success">{{
          version.baseModel
        }}</span>
        <span
          v-if="version.clientStatus === 'installed'"
          class="ms-1 badge rounded-pill text-bg-success"
          >Client</span
        >
        <span
          v-else-if="version.clientStatus === 'pending'"
          class="ms-1 badge rounded-pill text-bg-warning"
          >Syncing...</span
        >
        <button
          @click.stop="$emit('toggleNsfw', version)"
          class="btn btn-sm position-absolute top-0 end-0 m-2"
          :class="version.nsfw ? 'btn-danger' : 'btn-secondary'"
          style="--bs-btn-padding-y: 0.25rem; --bs-btn-padding-x: 0.25rem"
        >
          <svg
            v-if="version.nsfw"
            width="18px"
            height="18px"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            color="#ffffff"
          >
            <path
              d="M10.733 5.076a10.744 10.744 0 0 1 11.205 6.575 1 1 0 0 1 0 .696 10.747 10.747 0 0 1-1.444 2.49"
              stroke="#ffffff"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></path>
            <path
              d="M14.084 14.158a3 3 0 0 1-4.242-4.242"
              stroke="#ffffff"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></path>
            <path
              d="M17.479 17.499a10.75 10.75 0 0 1-15.417-5.151 1 1 0 0 1 0-.696 10.75 10.75 0 0 1 4.446-5.143"
              stroke="#ffffff"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></path>
            <path
              d="m2 2 20 20"
              stroke="#ffffff"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></path>
          </svg>
          <svg
            v-else
            width="18px"
            height="18px"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            color="#ffffff"
          >
            <path
              d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0"
              stroke="#ffffff"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></path>
            <circle
              cx="12"
              cy="12"
              r="3"
              stroke="#ffffff"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></circle>
          </svg>
        </button>
      </div>
    </div>
    <div class="card-body z-3">
      <h3 class="card-title h5">
        {{ model.name }} - {{ version.name }}
      </h3>
    </div>
    <div class="d-flex gap-2 card-footer z-2">
      <button
        @click="$emit('click', model.ID, version.ID)"
        class="btn btn-outline-primary"
      >
        More details
      </button>

      <button
        v-if="version.clientStatus === 'installed'"
        @click="dispatch('delete')"
        :disabled="isDispatching"
        class="btn btn-outline-danger btn-sm"
        title="Remove from Client"
      >
        <i class="bi bi-pc-display"></i> Remove
      </button>
      <button
        v-else-if="!version.clientStatus"
        @click="dispatch('download')"
        :disabled="isDispatching"
        class="btn btn-outline-info btn-sm"
        title="Push to Client"
      >
        <i class="bi bi-pc-display"></i> Push
      </button>
      <button
        @click="$emit('delete', version.ID)"
        class="btn btn-outline-danger ms-auto"
      >
        Delete
      </button>
    </div>
  </div>
</template>

<script setup>
import { useRemote } from '../composables/useRemote';

const props = defineProps({
  model: Object,
  version: Object,
  imageUrl: String,
});

const emit = defineEmits(["click", "delete", "toggleNsfw"]);

const { dispatchAction, isDispatching } = useRemote();
const dispatch = (action) => {
  dispatchAction(action, props.model, props.version);
};
</script>
