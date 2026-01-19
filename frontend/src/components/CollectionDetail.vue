<template>
  <div class="d-flex" style="min-height: calc(100vh - 80px)">
    <!-- Sidebar -->
    <aside 
      class="sidebar-wrapper bg-dark" 
      :class="{ 'show-mobile': showSidebar }"
    >
      <FilterSidebar
        v-model:search="search"
        v-model:tagsSearch="tagsSearch"
        v-model:selectedCategory="selectedCategory"
        v-model:selectedBaseModel="selectedBaseModel"
        v-model:selectedModelType="selectedModelType"
        v-model:nsfwFilter="nsfwFilter"
        v-model:syncedFilter="syncedFilter"
        :categories="categories"
        :baseModels="baseModels"
        :modelTypes="modelTypes"
        @clear="clearFilters"
        @close="showSidebar = false"
      />
    </aside>

    <!-- Overlay -->
    <div v-if="showSidebar" class="sidebar-overlay d-md-none" @click="showSidebar = false"></div>

    <!-- Main Content -->
    <main class="flex-grow-1 p-3 p-md-4" style="min-width: 0;">
        <!-- Header -->
        <div class="mb-4 d-flex justify-content-between align-items-start">
            <div>
                <button 
                  @click="router.push('/collections')" 
                  class="btn btn-dark btn-sm mb-3 d-flex align-items-center gap-2 border-0 ps-0 text-white"
                >
                  <Icon icon="mdi:arrow-left" width="20" height="20" />
                  Back to Collections
                </button>
                <h2 class="fw-bold mb-1">{{ collection.name }}</h2>
                <p class="text-white-50 mb-0">{{ collection.description }}</p>
            </div>
            <div class="d-flex gap-2">
                 <button class="btn btn-outline-primary d-flex align-items-center gap-2" @click="openBulkModal" title="Bulk Add by Tag">
                    <Icon icon="mdi:tag-plus-outline" width="20" height="20" />
                    <span class="d-none d-md-inline">Bulk Add</span>
                </button>
                <button 
                  class="btn btn-outline-secondary d-flex align-items-center gap-2" 
                  @click="openRenameModal" 
                  title="Rename Collection"
                >
                    <Icon icon="mdi:pencil" width="20" height="20" />
                    <span class="d-none d-md-inline">Edit</span>
                </button>
            </div>
        </div>

        <div class="m-4 text-center" v-if="loading">
             <div class="spinner-border text-primary" role="status"></div>
        </div>

        <AppPagination
            v-if="!loading && totalPages > 1"
            :page="page"
            :totalPages="totalPages"
            @changePage="changePage"
            class="px-4"
        />

        <div class="m-4 text-center" v-else-if="models.length === 0">
            No models in this collection found matching your filters.
        </div>

        <div class="model-grid p-4" v-else>
            <ModelCard
                v-for="version in models"
                :key="version.ID"
                :model="version.ParentModel"
                :version="version"
                :imageUrl="version.imageUrl"
                @click="goToModel(version.ParentModel.ID, version.ID)"
                @addToCollection="openCollectionModal"
                @delete="deleteVersionGlobal"
                :showCollectionRemove="true"
                @removeFromCollection="removeVersionFromCollection"
                @toggleNsfw="toggleVersionNsfw"
            />
        </div>

         <!-- Pagination (Basic for now, can implement server-side if needed) -->
         <!-- Since API supports pagination, we should use it -->
        <AppPagination
            v-if="totalPages > 1"
            :page="page"
            :totalPages="totalPages"
            @changePage="changePage"
        />

        <AddToCollectionModal 
            v-if="showCollectionModal" 
            :versionId="selectedVersionId" 
            @close="showCollectionModal = false" 
        />
    </main>

    <!-- Bulk Add Modal -->
    <div v-if="showBulkModal" class="modal d-block" tabindex="-1" style="background: rgba(0,0,0,0.5); backdrop-filter: blur(2px);" @click.self="showBulkModal = false">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content bg-dark text-white border-0 shadow-lg">
          <div class="modal-header border-0">
            <h5 class="modal-title fw-bold">Bulk Add by Tag</h5>
            <button type="button" class="btn-close btn-close-white" @click="showBulkModal = false"></button>
          </div>
          <div class="modal-body pt-0">
            <p class="text-white-50 small mb-3">
                Enter a tag to find all matching model versions and add them to this collection.
            </p>
            <div class="mb-3">
              <label class="form-label text-secondary small fw-bold text-uppercase">Tag</label>
              <input 
                type="text" 
                class="form-control bg-dark-subtle text-light border-0 shadow-none" 
                v-model="bulkTag" 
                ref="bulkInput"
                @keyup.enter="performBulkAdd"
                placeholder="e.g. anime, realistic..."
              >
            </div>
          </div>
          <div class="modal-footer border-0 pt-0">
            <button type="button" class="btn btn-outline-secondary border-0" @click="showBulkModal = false">Cancel</button>
            <button type="button" class="btn btn-primary shadow-sm" @click="performBulkAdd" :disabled="!bulkTag">
                <Icon icon="mdi:plus-circle-multiple-outline" class="me-1"/> Add Models
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Rename Modal -->
    <div v-if="showRenameModal" class="modal d-block" tabindex="-1" style="background: rgba(0,0,0,0.5); backdrop-filter: blur(2px);" @click.self="showRenameModal = false">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content bg-dark text-white border-0 shadow-lg">
          <div class="modal-header border-0">
            <h5 class="modal-title fw-bold">Edit Collection</h5>
            <button type="button" class="btn-close btn-close-white" @click="showRenameModal = false"></button>
          </div>
          <div class="modal-body pt-0">
            <div class="mb-3">
              <label class="form-label text-secondary small fw-bold text-uppercase">Name</label>
              <input 
                type="text" 
                class="form-control bg-dark-subtle text-light border-0 shadow-none" 
                v-model="renameForm.name" 
                ref="renameInput"
                @keyup.enter="updateCollection"
                placeholder="Collection name..."
              >
            </div>
            <div class="mb-3">
              <label class="form-label text-secondary small fw-bold text-uppercase">Description (Optional)</label>
              <textarea 
                class="form-control bg-dark-subtle text-light border-0 shadow-none" 
                v-model="renameForm.description"
                rows="3"
                placeholder="Add a description..."
              ></textarea>
            </div>
          </div>
          <div class="modal-footer border-0 pt-0">
            <button type="button" class="btn btn-outline-secondary border-0" @click="showRenameModal = false">Cancel</button>
            <button type="button" class="btn btn-primary shadow-sm" @click="updateCollection" :disabled="!renameForm.name">
                Save Changes
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import { Icon } from "@iconify/vue";
import { showToast, showDeleteConfirm, showConfirm } from '../utils/ui';
import ModelCard from './ModelCard.vue';
import FilterSidebar from './FilterSidebar.vue';
import AppPagination from './AppPagination.vue';
import AddToCollectionModal from './AddToCollectionModal.vue';
import debounce from '../utils/debounce';
import { useModels } from '../composables/useModels';

const route = useRoute();
const router = useRouter();

// Filters - Initialize from query params
const search = ref(route.query.search || "");
const tagsSearch = ref(route.query.tags || "");
const selectedCategory = ref(route.query.category || "");
const selectedBaseModel = ref(route.query.baseModel || "");
const selectedModelType = ref(route.query.modelType || "");
const nsfwFilter = ref(route.query.nsfw || "no");
const syncedFilter = ref(route.query.synced === "true");
const page = ref(parseInt(route.query.page) || 1);
const collection = ref({});
const models = ref([]); // These are actually Versions
const loading = ref(true);
const { showSidebar } = useModels();
const showCollectionModal = ref(false);
const selectedVersionId = ref(0);

// Bulk Add
const showBulkModal = ref(false);
const bulkTag = ref("");
const bulkInput = ref(null);

// Rename
const showRenameModal = ref(false);
const renameForm = ref({ name: "", description: "" });
const renameInput = ref(null);

const openRenameModal = () => {
    renameForm.value = { 
        name: collection.value.name, 
        description: collection.value.description 
    };
    showRenameModal.value = true;
    nextTick(() => renameInput.value?.focus());
};

const updateCollection = async () => {
    if (!renameForm.value.name) return;
    try {
        await axios.put(`/api/collections/${route.params.id}`, renameForm.value);
        showToast("Collection updated", "success");
        showRenameModal.value = false;
        fetchCollection(); // Refresh details
    } catch (err) {
        console.error(err);
        showToast("Failed to update collection", "danger");
    }
};

const deleteVersionGlobal = async (versionId) => {
      const choice = await showDeleteConfirm("Delete this version globally?");
      if (!choice) return;
      const files = choice === "deleteFiles" ? 1 : 0;
      try {
        await axios.delete(`/api/versions/${versionId}?files=${files}`);
        showToast("Version deleted", "success");
        fetchVersions();
      } catch (err) {
        showToast("Failed to delete version", "danger");
      }
};

const openBulkModal = () => {
    bulkTag.value = "";
    showBulkModal.value = true;
    // Focus next tick
    setTimeout(() => bulkInput.value?.focus(), 100);
};

const performBulkAdd = async () => {
    if (!bulkTag.value) return;
    try {
        const res = await axios.post(`/api/collections/${route.params.id}/add-by-tag`, {
            tag: bulkTag.value
        });
        showToast(`Added ${res.data.added} versions to collection`, "success");
        showBulkModal.value = false;
        fetchVersions();
    } catch (err) {
        console.error(err);
        showToast("Failed to bulk add versions", "danger");
    }
};

const total = ref(0);
const totalPages = ref(1);

// Sync filters to URL
const updateUrl = () => {
    const query = {
        ...route.query,
        page: page.value,
        search: search.value || undefined,
        tags: tagsSearch.value || undefined,
        category: selectedCategory.value || undefined,
        baseModel: selectedBaseModel.value || undefined,
        modelType: selectedModelType.value || undefined,
        nsfw: nsfwFilter.value !== 'no' ? nsfwFilter.value : undefined,
        synced: syncedFilter.value ? 'true' : undefined
    };
    
    // Remove undefined keys
    Object.keys(query).forEach(key => query[key] === undefined && delete query[key]);
    
    router.replace({ query });
};

// Static data (could be fetched or shared)
const baseModels = ref([]);
// ... (populate like useModels)
// For now, fetch baseModels
const fetchBaseModels = async () => {
    try {
        const res = await axios.get("/api/base-models");
        baseModels.value = res.data;
    } catch (e) {}
};

const modelTypes = [
  "Checkpoint", "TextualInversion", "Hypernetwork", "AestheticGradient",
  "LORA", "LoCon", "DoRA", "Controlnet", "Upscaler", "MotionModule",
  "VAE", "Wildcards", "Poses", "Workflows", "Detection", "Other",
];

const categories = [
  "character", "style", "concept", "clothing", "base model", "poses",
  "background", "tool", "vehicle", "buildings", "objects", "assets",
  "animal", "action",
];

const fetchCollection = async () => {
    try {
        const res = await axios.get(`/api/collections/${route.params.id}`);
        collection.value = res.data;
    } catch (err) {
        showToast("Collection not found", "danger");
        router.push('/collections');
    }
};

const fetchVersions = async () => {
    loading.value = true;
    try {
        const params = {
            page: page.value,
            limit: 50,
            search: search.value,
            tags: [selectedCategory.value, tagsSearch.value].filter(Boolean).join(','),
            baseModel: selectedBaseModel.value,
            modelType: selectedModelType.value,
            nsfwFilter: nsfwFilter.value,
            synced: syncedFilter.value ? "1" : undefined
        };
        
        const res = await axios.get(`/api/collections/${route.params.id}/versions`, { params });
        
        // Handle pagination from header
        if (res.headers['x-total-count']) {
            total.value = parseInt(res.headers['x-total-count']);
            totalPages.value = Math.ceil(total.value / 50);
        }

        if (Array.isArray(res.data)) {
            models.value = res.data.map(v => {
                // Determine model name: use v.model.name if available, otherwise "Unknown"
                // API JSON tag for ParentModel is "model"
                
                // Construct a model object compatible with ModelCard
                const parentModel = v.model || { name: 'Unknown Model', ID: v.modelId };
                
                return {
                    ...v,
                    ParentModel: parentModel, // Store consistently for template
                    imageUrl: normalizeImagePath(v.imagePath)
                };
            });
        } else {
             models.value = [];
             console.error("Invalid versions response", res.data);
        }
        // We assume backend returns ALL for pagination if we don't have separate count endpoint for filtered versions?
        // Wait, I implemented handlers.go GetModelsCount but not GetCollectionVersionsCount.
        // For MVP, if pagination is needed, we need count. 
        // My handlers.go implementation for GetCollectionVersions uses Limit/Offset.
        // But it doesn't return count.
        // For now, simplified pagination or just load more?
        // I'll assume infinite scroll or just standard pagination without knowing total (Prev/Next only)?
        // Or I can add a count header or simple total.
        // Let's assume one page for MVP or add count to response later. 
        // Actually, without total count, AppPagination breaks?
        // I'll set totalPages to 1 or implement load more if needed.
        // Let's check AppPagination.
    } catch (err) {
        console.error(err);
        showToast("Failed to load versions", "danger");
    } finally {
        loading.value = false;
        updateUrl();
        
        // Handle Scroll Restoration
        if (route.query.scrollTo) {
            nextTick(() => {
                const el = document.getElementById(`model-${route.query.scrollTo}`);
                if (el) {
                    el.scrollIntoView({ behavior: 'smooth', block: 'center' });
                    // Remove scrollTo from URL
                    const query = { ...route.query };
                    delete query.scrollTo;
                    router.replace({ query }); 
                }
            });
        }
    }
};

const normalizeImagePath = (path) => {
    if (!path) return null;
    let normalized = path.replace(/\\/g, "/");
    if (normalized.includes("/backend/images/")) {
        normalized = normalized.replace(/^.*\/backend\/images/, "/images");
    }
    if (!normalized.startsWith("/") && !normalized.startsWith("http")) {
        normalized = "/images/" + normalized;
    }
    return normalized;
};

const debouncedFetch = debounce(() => {
    page.value = 1;
    fetchVersions();
}, 300);

watch([search, tagsSearch, selectedCategory, selectedBaseModel, selectedModelType, nsfwFilter, syncedFilter], debouncedFetch);
watch(page, fetchVersions);

const clearFilters = () => {
  search.value = "";
  tagsSearch.value = "";
  selectedCategory.value = "";
  selectedBaseModel.value = "";
  selectedModelType.value = "";
  nsfwFilter.value = "no";
  syncedFilter.value = false;
  page.value = 1;
};

const goToModel = (modelId, versionId) => {
  router.push({
    name: "ModelDetail",
    params: { modelId, versionId },
    query: { returnPath: route.fullPath }
  });
};

const openCollectionModal = (versionId) => {
    selectedVersionId.value = versionId;
    showCollectionModal.value = true;
};



const removeVersionFromCollection = async (versionId) => {
    if(!(await showConfirm("Remove this version from this collection?"))) return;

    try {
        await axios.delete(`/api/collections/${route.params.id}/versions/${versionId}`);
        showToast("Removed from collection", "success");
        fetchVersions();
    } catch (err) {
        console.error(err);
        showToast("Failed to remove from collection", "danger");
    }
};

const toggleVersionNsfw = async (version) => {
  const updated = { ...version, nsfw: !version.nsfw };
  try {
    // Determine the ID to update. Collection versions list often has ID as the join ID?
    // Wait, in fetchVersions mapper:
    // models.value = res.data.map(v => return { ...v, ParentModel... })
    // So 'version' passed here is the version object from models array.
    // 'version.ID' should be the Version ID.
    
    await axios.put(`/api/versions/${version.ID}`, updated);
    
    // Update local state to reflect change immediately
    const v = models.value.find(v => v.ID === version.ID);
    if (v) {
        v.nsfw = updated.nsfw;
    }
    showToast("NSFW status updated", "success");
  } catch (err) {
    console.error(err);
    showToast("Failed to update NSFW status", "danger");
  }
};

// Renamed from deleteVersionFromCollection to avoid confusion with global delete

const changePage = (p) => page.value = p;

onMounted(() => {
    fetchCollection();
    fetchBaseModels();
    fetchVersions();
});
</script>

<style scoped>
.sidebar-wrapper {
  width: 280px;
  flex-shrink: 0;
  transition: transform 0.3s ease-in-out;
  position: sticky;
  top: 80px;
  height: calc(100vh - 80px);
  overflow-y: auto;
}
@media (max-width: 767.98px) {
  .sidebar-wrapper {
    position: fixed;
    top: 0;
    left: 0;
    height: 100vh;
    z-index: 1045;
    transform: translateX(-100%);
  }
  .sidebar-wrapper.show-mobile {
    transform: translateX(0);
  }
  .sidebar-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0,0,0,0.5);
    z-index: 1040;
  }
}
</style>
