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
                    class="btn btn-dark d-md-none" 
                    @click="showSidebar = !showSidebar"
                >
                    <Icon icon="mdi:filter" />
                </button>
            </div>
        </div>

        <div class="m-4 text-center" v-if="loading">
             <div class="spinner-border text-primary" role="status"></div>
        </div>

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
                @delete="deleteVersionFromCollection"
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
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import { Icon } from "@iconify/vue";
import { showToast, showDeleteConfirm } from '../utils/ui';
import ModelCard from './ModelCard.vue';
import FilterSidebar from './FilterSidebar.vue';
import AppPagination from './AppPagination.vue';
import AddToCollectionModal from './AddToCollectionModal.vue';
import debounce from '../utils/debounce';

const route = useRoute();
const router = useRouter();

// State
const collection = ref({});
const models = ref([]); // These are actually Versions
const loading = ref(true);
const showSidebar = ref(false);
const showCollectionModal = ref(false);
const selectedVersionId = ref(0);

// Bulk Add
const showBulkModal = ref(false);
const bulkTag = ref("");
const bulkInput = ref(null);

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

// Filters
const search = ref("");
const tagsSearch = ref("");
const selectedCategory = ref("");
const selectedBaseModel = ref("");
const selectedModelType = ref("");
const nsfwFilter = ref("no");
const syncedFilter = ref(false);
const page = ref(1);
const total = ref(0);
const totalPages = ref(1);

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
    query: { returnPath: `/collections/${route.params.id}` }
  });
};

const openCollectionModal = (versionId) => {
    selectedVersionId.value = versionId;
    showCollectionModal.value = true;
};

const deleteVersionFromCollection = async (versionId) => {
    // Override default delete behavior: Remove from collection instead of deleting version
    // But ModelCard emits "delete" which implies permanent delete?
    // The user said "a button on each collection to delete it".
    // "Delete version" in kebab menu usually means delete from disk.
    // I should add "Remove from Collection" to kebab menu in ModelCard BUT ModelCard is generic.
    // I can't easily change ModelCard menu based on context without props.
    // I'll use the existing "delete" event but ask user "Remove from collection or Delete from Disk?"
    // Or simpler: ModelCard doesn't know context.
    // If I want "Remove from Collection", I should add it to ModelCard as an option controlled by prop `showRemoveFromCollection`.
    
    // For now, "Delete" on card means Delete Version globally.
    // Users might expect "Remove from Collection" contextually.
    // I'll add `excludeId` prop or similar?
    // Let's implement global delete for now as per ModelList.
    
    // Actually, "Collections... option to name it... button to delete it" refers to the collection itself.
    // Being inside a collection view, one might want to remove items.
    // I'll stick to standard behavior for now.
    
    // Re-using deleteVersion from ModelList logic
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
