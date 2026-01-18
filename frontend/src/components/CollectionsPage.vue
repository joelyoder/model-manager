<template>
  <div class="container-fluid p-4">
    <div class="d-flex flex-column flex-md-row justify-content-between align-items-start align-items-md-center mb-4 gap-3">
      <div>
        <h2 class="fw-bold mb-1">Collections</h2>
      </div>
      <div>
        <button class="btn btn-outline-primary d-flex align-items-center gap-2" @click="openCreateModal">
          <Icon icon="mdi:folder-plus" width="20" height="20" />
          New Collection
        </button>
      </div>
    </div>

    <!-- Search Bar -->
    <div class="row mb-4">
      <div class="col-md-6 col-lg-4">
        <div class="input-group">
          <span class="input-group-text bg-dark-subtle border-end-0">
            <Icon icon="mdi:magnify" class="text-muted" />
          </span>
          <input 
            type="text" 
            class="form-control bg-dark-subtle border-start-0" 
            placeholder="Search collections..." 
            v-model="searchQuery"
            @input="debouncedSearch"
          >
        </div>
      </div>
    </div>

    <!-- Collections Grid -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <div v-else-if="collections.length === 0" class="text-center py-5">
      <div class="text-muted mb-3">
        <Icon icon="mdi:folder-open-outline" width="64" height="64" class="opacity-50" />
      </div>
      <h5>No collections found</h5>
      <p class="text-muted">Create a collection to get started or try a different search.</p>
      <button class="btn btn-outline-primary btn-sm mt-2" @click="openCreateModal" v-if="!searchQuery">
        Create Collection
      </button>
    </div>

    <div v-else class="row row-cols-1 row-cols-md-2 row-cols-lg-3 row-cols-xl-4 g-4">
      <div class="col" v-for="collection in collections" :key="collection.ID">
        <CollectionCard 
            :collection="collection" 
            @click="viewCollection" 
            @rename="openRenameModal"
            @delete="confirmDelete"
        />
      </div>
    </div>

    <!-- Create/Rename Modal -->
    <div v-if="showModal" class="modal d-block" tabindex="-1" style="background: rgba(0,0,0,0.5); backdrop-filter: blur(2px);" @click.self="closeModal">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content bg-dark text-white border-0 shadow-lg">
          <div class="modal-header border-0">
            <h5 class="modal-title fw-bold">{{ isEditing ? 'Rename Collection' : 'New Collection' }}</h5>
            <button type="button" class="btn-close btn-close-white" @click="closeModal"></button>
          </div>
          <div class="modal-body pt-0">
            <div class="mb-3">
              <label class="form-label text-secondary small fw-bold text-uppercase">Name</label>
              <input 
                type="text" 
                class="form-control bg-dark-subtle text-light border-0 shadow-none" 
                v-model="modalForm.name" 
                ref="nameInput"
                @keyup.enter="saveCollection"
                placeholder="Collection name..."
              >
            </div>
            <div class="mb-3">
              <label class="form-label text-secondary small fw-bold text-uppercase">Description (Optional)</label>
              <textarea 
                class="form-control bg-dark-subtle text-light border-0 shadow-none" 
                v-model="modalForm.description"
                rows="3"
                placeholder="Add a description..."
              ></textarea>
            </div>
          </div>
          <div class="modal-footer border-0 pt-0">
            <button type="button" class="btn btn-outline-secondary border-0" @click="closeModal">Cancel</button>
            <button type="button" class="btn btn-primary shadow-sm" @click="saveCollection" :disabled="!modalForm.name">
                {{ isEditing ? 'Save Changes' : 'Create Collection' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { Icon } from "@iconify/vue";
import CollectionCard from './CollectionCard.vue';
import axios from 'axios';
import { showToast, showConfirm } from '../utils/ui';

const router = useRouter();
const collections = ref([]);
const loading = ref(true);
const searchQuery = ref("");
const showModal = ref(false);
const isEditing = ref(false);
const modalForm = ref({ id: null, name: "", description: "" });
const nameInput = ref(null);

let searchTimeout;

const fetchCollections = async () => {
    loading.value = true;
    try {
        const res = await axios.get('/api/collections', {
            params: { search: searchQuery.value }
        });
        if (Array.isArray(res.data)) {
            collections.value = res.data;
        } else {
            // If the backend isn't updated, it might return HTML (SPA fallback)
            console.error("Invalid response format, expected JSON array", res.data);
            collections.value = [];
            showToast("Failed to load collections: Backend might be outdated", "danger");
        }
    } catch (err) {
        console.error(err);
        showToast("Failed to load collections", "danger");
    } finally {
        loading.value = false;
    }
};

const debouncedSearch = () => {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(fetchCollections, 300);
};

const viewCollection = (collection) => {
    router.push(`/collections/${collection.ID}`);
};

const openCreateModal = () => {
    isEditing.value = false;
    modalForm.value = { id: null, name: "", description: "" };
    showModal.value = true;
    nextTick(() => nameInput.value?.focus());
};

const openRenameModal = (collection) => {
    isEditing.value = true;
    modalForm.value = { id: collection.ID, name: collection.name, description: collection.description };
    showModal.value = true;
    nextTick(() => nameInput.value?.focus());
};

const closeModal = () => {
    showModal.value = false;
};

const saveCollection = async () => {
    if (!modalForm.value.name) return;
    
    try {
        if (isEditing.value) {
            await axios.put(`/api/collections/${modalForm.value.id}`, {
                name: modalForm.value.name,
                description: modalForm.value.description
            });
            showToast("Collection updated", "success");
        } else {
            await axios.post('/api/collections', {
                name: modalForm.value.name,
                description: modalForm.value.description
            });
            showToast("Collection created", "success");
        }
        closeModal();
        fetchCollections();
    } catch (err) {
        console.error(err);
        showToast("Failed to save collection", "danger");
    }
};

const confirmDelete = async (collection) => {
    const choice = await showConfirm(`Are you sure you want to delete the collection "${collection.name}"? The models inside will not be deleted.`);
    if (!choice) return;
    
    try {
        await axios.delete(`/api/collections/${collection.ID}`);
        showToast("Collection deleted", "success");
        fetchCollections();
    } catch (err) {
        console.error(err);
        showToast("Failed to delete collection", "danger");
    }
};

onMounted(fetchCollections);
</script>
