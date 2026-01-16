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
      >
        <template #actions>
        </template>
      </FilterSidebar>
    </aside>

    <!-- Overlay -->
    <div v-if="showSidebar" class="sidebar-overlay d-md-none" @click="showSidebar = false"></div>

    <!-- Main Content -->
    <main class="flex-grow-1 p-3 p-md-4" style="min-width: 0;">
        <!-- Mobile Header Removed -->

        <AppPagination
            :page="page"
            :totalPages="totalPages"
            @changePage="changePage"
        />

        <div class="m-4 text-center" v-if="models.length === 0">
            No models found.
        </div>

        <div class="model-grid p-4">
            <ModelCard
                v-for="card in versionCards"
                :key="card.version.ID"
                :model="card.model"
                :version="card.version"
                :imageUrl="card.imageUrl"
                @click="goToModel"
                @delete="deleteVersion"
                @toggleNsfw="toggleVersionNsfw"
            />
        </div>

        <AppPagination
            :page="page"
            :totalPages="totalPages"
            @changePage="changePage"
        />
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";
import { Icon } from "@iconify/vue";
import { showToast, showDeleteConfirm } from "../utils/ui";
import { useModels } from "../composables/useModels";
import FilterSidebar from "./FilterSidebar.vue";
import AppPagination from "./AppPagination.vue";
import ModelCard from "./ModelCard.vue";

const router = useRouter();
const route = useRoute();
const showAddPanel = ref(false); // remove this

const {
  models,
  search,
  tagsSearch,
  selectedCategory,
  selectedBaseModel,
  selectedModelType,
  nsfwFilter,
  syncedFilter,
  page,
  totalPages,
  baseModels,
  modelTypes,
  categories,
  showSidebar, // Destructured from useModels
  init,
  clearFilters,
  fetchModels,
} = useModels();

const changePage = (p) => {
  page.value = p;
};



const goToModel = (modelId, versionId) => {
  router.push({
    name: "ModelDetail",
    params: { modelId, versionId },
  });
};

const deleteVersion = async (versionId) => {
  const choice = await showDeleteConfirm("Delete this version?");
  if (!choice) return;
  const files = choice === "deleteFiles" ? 1 : 0;
  try {
    await axios.delete(`/api/versions/${versionId}?files=${files}`);
    showToast("Version deleted", "success");
    fetchModels();
  } catch (err) {
    console.error(err);
    showToast("Failed to delete version", "danger");
  }
};

const toggleVersionNsfw = async (version) => {
  const updated = { ...version, nsfw: !version.nsfw };
  try {
    await axios.put(`/api/versions/${version.ID}`, updated);
    
    // Update the source of truth in models.value to trigger reactivity
    const model = models.value.find(m => m.ID === version.modelId);
    if (model && model.versions) {
      const v = model.versions.find(v => v.ID === version.ID);
      if (v) {
        v.nsfw = updated.nsfw;
      }
    }
    showToast("NSFW status updated", "success");
  } catch (err) {
    console.error(err);
    showToast("Failed to update NSFW status", "danger");
  }
};

const matchesNsfwFilter = (value) => {
  const isNsfw = Boolean(value);
  if (nsfwFilter.value === "no") return !isNsfw;
  if (nsfwFilter.value === "only") return isNsfw;
  return true;
};

const modelMatchesNsfwFilter = (model) => {
  const versions = model.versions || [];
  if (nsfwFilter.value === "no") {
    if (versions.length) return versions.some((v) => !v.nsfw);
    return !model.nsfw;
  }
  if (nsfwFilter.value === "only") {
    if (versions.length) return versions.some((v) => Boolean(v.nsfw));
    return Boolean(model.nsfw);
  }
  return true;
};

const filteredModels = computed(() => {
  return models.value.filter((m) => {
    if (!modelMatchesNsfwFilter(m)) return false;
    // The API should have handled search, but if we want to be safe:
    if (search.value) {
      const s = search.value.toLowerCase();
      const matchModel = m.name.toLowerCase().includes(s);
      const matchVersion = (m.versions || []).some((v) =>
        v.name.toLowerCase().includes(s)
      );
      const matchTrained = (m.versions || []).some((v) =>
        (v.trainedWords || "").toLowerCase().includes(s)
      );
      if (!(matchModel || matchVersion || matchTrained)) return false;
    }
    return true;
  });
});

const versionCards = computed(() => {
  const sortedModels = filteredModels.value.slice().sort((a, b) => b.ID - a.ID);

  return sortedModels.flatMap((model) => {
    const versionsSorted = (model.versions || [])
      .slice()
      .sort((a, b) => b.ID - a.ID);

    return versionsSorted
      .filter((v) => {
        if (selectedBaseModel.value && v.baseModel !== selectedBaseModel.value)
          return false;
        if (selectedModelType.value && v.type !== selectedModelType.value)
          return false;
        if (!matchesNsfwFilter(v.nsfw)) return false;
        if (syncedFilter.value && v.clientStatus !== 'installed') return false;

        if (search.value) {
          const s = search.value.toLowerCase();
          const matchModel = model.name.toLowerCase().includes(s);
          const matchVersion = v.name.toLowerCase().includes(s);
          const matchTrained = (v.trainedWords || "").toLowerCase().includes(s);
          if (!(matchModel || matchVersion || matchTrained)) return false;
        }

        if (selectedCategory.value) {
          const tags = (v.tags || "")
            .split(",")
            .map((t) => t.trim().toLowerCase());
          if (!tags.includes(selectedCategory.value.toLowerCase()))
            return false;
        }

        if (tagsSearch.value.trim()) {
          const tags = (v.tags || "")
            .split(",")
            .map((t) => t.trim().toLowerCase());
          const searchTags = tagsSearch.value
            .split(/[, ]+/)
            .map((t) => t.trim().toLowerCase())
            .filter(Boolean);
          if (!searchTags.every((t) => tags.includes(t))) return false;
        }
        return true;
      })
      .map((v) => {
        return {
          model,
          version: v, 
          imageUrl: v.imageUrl || model.imageUrl,
        };
      });
  });
});

onMounted(async () => {
  await init();
  if (route.query.scrollTo) {
    await nextTick();
    const el = document.getElementById(`model-${route.query.scrollTo}`);
    if (el) {
      el.scrollIntoView({ behavior: "smooth" });
    }
    const rest = { ...route.query };
    delete rest.scrollTo;
    router.replace({ path: route.path, query: rest });
  }
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
