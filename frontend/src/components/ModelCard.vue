<template>
  <div
    :id="`model-${version.ID}`"
    class="model-card card h-100 border-0 shadow-sm clickable-card transition-hover"
    @click="$emit('click', model.ID, version.ID)"
  >
    <!-- Image & Overlays -->
    <div class="position-relative overflow-hidden rounded-top-2">
      <img
        v-if="imageUrl"
        :src="imageUrl"
        class="card-img-top object-fit-cover transition-zoom"
        style="height: 450px; width: 100%"
        loading="lazy"
      />
      
      <!-- Top Left Badges -->
      <div class="position-absolute top-0 start-0 p-2 d-flex flex-row gap-1 z-2 flex-wrap">
        <span class="badge rounded-pill fw-normal shadow-sm" :class="getBadgeColor(version.type)">{{ version.type }}</span>
        <span class="badge rounded-pill fw-normal shadow-sm" :class="getBadgeColor(version.baseModel)">{{ version.baseModel }}</span>
      </div>

      <!-- Top Right Actions -->
      <div class="position-absolute top-0 end-0 p-2 z-2">
        <button
          @click.stop="$emit('toggleNsfw', version)"
          class="btn btn-sm rounded-circle d-flex align-items-center justify-content-center p-1"
          :class="version.nsfw ? 'btn-danger' : 'btn-dark bg-opacity-50 text-white'"
          style="width: 32px; height: 32px; backdrop-filter: blur(4px);"
          title="Toggle NSFW"
        >
          <Icon :icon="version.nsfw ? 'mdi:eye-off' : 'mdi:eye'" width="18" height="18" />
        </button>
      </div>
      
      <!-- Hover Gradient (Optional Polish) -->
      <div class="model-card-overlay position-absolute bottom-0 start-0 w-100 h-25 bg-gradient-to-t-transparent-black opacity-0"></div>
    </div>

    <!-- Minimal Footer -->
    <div class="card-body p-3 z-3 bg-dark-subtle rounded-bottom-2">
      <div class="d-flex justify-content-between align-items-start gap-2">
        <div class="flex-grow-1 min-width-0">
          <h6 class="card-title fw-bold mb-1 lh-sm text-body-emphasis">
            {{ model.name }}
          </h6>
          <div class="small text-muted mb-0">
            {{ version.name }}
          </div>
        </div>

        <!-- Actions Group -->
        <div class="d-flex align-items-center gap-1 flex-shrink-0">
             <!-- Smart Remote Button -->
            <button
                v-if="version.clientStatus === 'installed'"
                @click.stop="dispatch('delete')"
                @mouseenter="isHovering = true"
                @mouseleave="isHovering = false"
                :disabled="isDispatching"
                class="btn btn-sm rounded-circle d-flex align-items-center justify-content-center"
                :class="isHovering ? 'btn-danger' : 'btn-success'"
                :title="isHovering ? 'Remove from Client' : 'Installed on Client'"
                style="width: 32px; height: 32px;"
            >
                <Icon 
                :icon="isHovering ? 'mdi:trash-can' : 'mdi:check'" 
                width="16" 
                height="16" 
                />
            </button>
            <button
                v-else-if="version.clientStatus === 'pending'"
                disabled
                class="btn btn-warning btn-sm rounded-circle d-flex align-items-center justify-content-center"
                style="width: 32px; height: 32px;"
                title="Syncing..."
            >
                <span class="spinner-border spinner-border-sm" style="width:1rem;height:1rem;" aria-hidden="true"></span>
            </button>
            <button
                v-else
                @click.stop="dispatch('download')"
                :disabled="isDispatching"
                class="btn btn-outline-secondary btn-sm rounded-circle d-flex align-items-center justify-content-center download-btn"
                style="width: 32px; height: 32px;"
                title="Push to Client"
            >
                <Icon icon="mdi:cloud-download" width="16" height="16" />
            </button>

            <!-- Kebab Menu for Delete -->
            <div class="dropdown">
                <button 
                    class="btn btn-link p-0 ms-1 kebab-btn" 
                    type="button" 
                    data-bs-toggle="dropdown" 
                    aria-expanded="false"
                    @click.stop
                >
                    <Icon icon="mdi:dots-vertical" width="20" height="20" />
                </button>
                <ul class="dropdown-menu dropdown-menu-end shadow">
                    <li v-if="showCollectionRemove">
                        <button class="dropdown-item text-danger d-flex align-items-center" @click.stop="$emit('removeFromCollection', version.ID)">
                             <Icon icon="mdi:folder-remove-outline" class="me-2" /> Remove from Collection
                        </button>
                    </li>
                    <li>
                        <button class="dropdown-item d-flex align-items-center" @click.stop="$emit('addToCollection', version.ID)">
                             <Icon icon="mdi:folder-plus" class="me-2" /> Add to Collection
                        </button>
                    </li>
                    <li>
                        <button class="dropdown-item text-danger d-flex align-items-center" @click.stop="$emit('delete', version.ID)">
                             <Icon icon="mdi:trash-can-outline" class="me-2" /> Delete Version
                        </button>
                    </li>
                </ul>
            </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRemote } from '../composables/useRemote';
import { ref } from 'vue';
import { Icon } from "@iconify/vue";
import { getBadgeColor } from "../utils/colors";

const isHovering = ref(false);

const props = defineProps({
  model: Object,
  version: Object,
  imageUrl: String,
  showCollectionRemove: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(["click", "delete", "toggleNsfw", "addToCollection", "removeFromCollection"]);

const { dispatchAction, isDispatching } = useRemote();
const dispatch = (action) => {
  dispatchAction(action, props.model, props.version);
};

</script>

<style scoped>
.transition-zoom {
  transition: transform 0.4s ease;
}

.model-card:hover .transition-zoom {
  transform: scale(1.05); /* Slight zoom on hover */
}

/* Kebab button hover brighter */
.kebab-btn {
  color: var(--bs-secondary);
  transition: color 0.2s;
}
.kebab-btn:hover {
  color: var(--bs-white); /* Brighter on hover */
}

/* Download button fixes */
.download-btn {
  border-color: var(--bs-border-color-translucent);
  color: var(--bs-secondary);
}
.download-btn:hover {
  background-color: var(--bs-secondary);
  color: white;
  border-color: var(--bs-secondary);
}
</style>
