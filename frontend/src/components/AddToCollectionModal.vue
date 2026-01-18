<template>
  <div class="modal d-block" tabindex="-1" style="background: rgba(0,0,0,0.5); backdrop-filter: blur(2px);" @click.self="$emit('close')">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content bg-dark text-white border-0 shadow-lg">
        <div class="modal-header border-0">
          <h5 class="modal-title fw-bold">Add to Collection</h5>
          <button type="button" class="btn-close btn-close-white" @click="$emit('close')"></button>
        </div>
        <div class="modal-body pt-0">
            <!-- New Collection Input -->
            <div class="input-group mb-4">
                <input 
                    type="text" 
                    class="form-control bg-dark-subtle text-light border-0 shadow-none" 
                    placeholder="New collection name..." 
                    v-model="newCollectionName"
                    @keyup.enter="createCollection"
                >
                <button class="btn btn-primary shadow-sm" type="button" @click="createCollection" :disabled="!newCollectionName">
                    <Icon icon="mdi:plus" /> Create
                </button>
            </div>

            <div class="d-flex align-items-center mb-2">
                 <span class="text-secondary small fw-bold text-uppercase">Your Collections</span>
                 <div class="flex-grow-1 ms-3 border-top border-secondary opacity-25"></div>
            </div>

            <!-- Collection List -->
             <div v-if="loading" class="text-center py-4">
                <div class="spinner-border spinner-border-sm text-primary" role="status"></div>
            </div>
            <div v-else class="list-group list-group-flush" style="max-height: 300px; overflow-y: auto;">
                <div v-if="collections.length === 0" class="text-center text-muted py-4">
                    <Icon icon="mdi:folder-open-outline" width="32" height="32" class="mb-2 opacity-50"/>
                    <div>No collections found</div>
                </div>
                <label 
                    v-for="col in collections" 
                    :key="col.ID" 
                    class="list-group-item bg-transparent text-light border-0 px-2 py-2 d-flex justify-content-between align-items-center cursor-pointer rounded-2 hover-bg-subtle transition-all"
                >
                    <div class="d-flex align-items-center gap-2 overflow-hidden">
                        <Icon icon="mdi:folder-outline" class="text-secondary flex-shrink-0" />
                        <span class="text-truncate">{{ col.name }}</span>
                    </div>
                    <div class="form-check m-0">
                        <input 
                            class="form-check-input bg-dark-subtle border-secondary shadow-none cursor-pointer" 
                            type="checkbox" 
                            :checked="isInCollection(col.ID)" 
                            @change="toggleCollection(col, $event.target.checked)"
                        >
                    </div>
                </label>
            </div>
        </div>
        <div class="modal-footer border-0 pt-0">
            <button type="button" class="btn btn-outline-secondary border-0" @click="$emit('close')">Done</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { showToast } from '../utils/ui';

const props = defineProps({
    versionId: {
        type: Number,
        required: true
    }
});

const emit = defineEmits(['close']);

const collections = ref([]);
const versionCollections = ref([]); // IDs of collections this version is in
const loading = ref(true);
const newCollectionName = ref("");

const fetchCollections = async () => {
    loading.value = true;
    try {
        // Fetch all collections
        const [allRes, verRes] = await Promise.all([
            axios.get('/api/collections'),
            axios.get(`/api/versions/${props.versionId}/collections`)
        ]);
        
        collections.value = Array.isArray(allRes.data) ? allRes.data : [];
        // Map version collections to IDs for easy checking
        versionCollections.value = (Array.isArray(verRes.data) ? verRes.data : []).map(c => c.ID);
    } catch (err) {
        console.error(err);
        showToast("Failed to load collections", "danger");
    } finally {
        loading.value = false;
    }
};

const isInCollection = (collectionId) => {
    return versionCollections.value.includes(collectionId);
};

const toggleCollection = async (collection, isChecked) => {
    try {
        if (isChecked) {
            await axios.post(`/api/collections/${collection.ID}/versions`, { versionId: props.versionId });
            versionCollections.value.push(collection.ID);
            showToast(`Added to "${collection.name}"`, "success");
        } else {
            await axios.delete(`/api/collections/${collection.ID}/versions/${props.versionId}`);
            versionCollections.value = versionCollections.value.filter(id => id !== collection.ID);
            showToast(`Removed from "${collection.name}"`, "info");
        }
    } catch (err) {
        console.error(err);
        showToast("Failed to update collection", "danger");
        // Revert check state visually if possible, but simpler to just reload or let user try again.
        // For simplicity, we assume generic error handling is enough.
    }
};

const createCollection = async () => {
    if (!newCollectionName.value) return;
    try {
        const res = await axios.post('/api/collections', { name: newCollectionName.value });
        const newCol = res.data;
        collections.value.push(newCol);
        newCollectionName.value = "";
        
        // Auto-add to the new collection
        await toggleCollection(newCol, true);
    } catch (err) {
        console.error(err);
        showToast("Failed to create collection", "danger");
    }
};

onMounted(fetchCollections);
</script>

<style scoped>
.cursor-pointer {
    cursor: pointer;
}
/* Custom scrollbar for dark mode if needed, standard usually fine */
.hover-bg-subtle:hover {
    background-color: var(--bs-dark-bg-subtle) !important;
}
.transition-all {
    transition: all 0.2s;
}
</style>
