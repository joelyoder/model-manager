<template>
  <div class="mx-4">
    <FilterBar
      v-model:search="search"
      v-model:tagsSearch="tagsSearch"
      v-model:selectedCategory="selectedCategory"
      v-model:selectedBaseModel="selectedBaseModel"
      v-model:selectedModelType="selectedModelType"
      v-model:nsfwFilter="nsfwFilter"
      :categories="categories"
      :baseModels="baseModels"
      :modelTypes="modelTypes"
      @clear="clearFilters"
    >
      <template #actions>
        <button
          type="button"
          class="btn btn-outline-primary d-inline-flex align-items-center justify-content-center"
          @click="showAddPanel = !showAddPanel"
          :aria-label="showAddPanel ? 'Close panel' : 'Add models'"
          :title="showAddPanel ? 'Close panel' : 'Add models'"
        >
          <Icon
            :icon="showAddPanel ? 'mdi:close' : 'mdi:plus'"
            width="20"
            height="20"
          />
          <span class="visually-hidden">
            {{ showAddPanel ? "Close Panel" : "Add Models" }}
          </span>
        </button>
      </template>
    </FilterBar>

    <AddModelPanel
      v-show="showAddPanel"
      @createManual="createManualModel"
      @added="fetchModels"
    />

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
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";
import { Icon } from "@iconify/vue";
import { showToast, showDeleteConfirm } from "../utils/ui";
import { useModels } from "../composables/useModels";
import FilterBar from "./FilterBar.vue";
import AppPagination from "./AppPagination.vue";
import ModelCard from "./ModelCard.vue";
import AddModelPanel from "./AddModelPanel.vue";

const router = useRouter();
const route = useRoute();
const showAddPanel = ref(false);

const {
  models,
  search,
  tagsSearch,
  selectedCategory,
  selectedBaseModel,
  selectedModelType,
  nsfwFilter,
  page,
  totalPages,
  baseModels,
  modelTypes,
  categories,
  init,
  clearFilters,
  fetchModels,
} = useModels();

const changePage = (p) => {
  page.value = p;
};

const createManualModel = async () => {
  try {
    const res = await axios.post("/api/models", {
      name: "New Model",
      type: "Checkpoint",
    });
    await fetchModels();
    router.push({
      name: "ModelDetail",
      params: { versionId: res.data.version.ID },
      query: { edit: "1" },
    });
  } catch (err) {
    console.error(err);
    showToast("Failed to create model", "danger");
  }
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

// Computed properties for filtering logic that was in ModelList
// Wait, useModels handles fetching, but does it handle client-side filtering?
// The original code had `filteredModels` and `versionCards` computed properties.
// `useModels` fetches from API with params.
// Let's check `useModels` again. It fetches from API.
// But `ModelList.vue` original code had `filteredModels` which filtered `models.value` client-side as well?
// Original `fetchModels` fetched with params.
// Original `filteredModels` filtered `models.value` again?
// Yes, lines 673-689 in original `ModelList.vue` filtered `models.value`.
// But `fetchModels` (line 471) already sent params to server.
// It seems the server filtering might be incomplete or the client filtering is for immediate feedback?
// Actually, `debouncedUpdate` calls `fetchModels`.
// The `filteredModels` computed property seems to be doing additional filtering or maybe the server returns more than needed?
// Wait, `models.value` is populated by `fetchModels`.
// If `fetchModels` uses the search params, then `models.value` should already be filtered.
// The original code had:
// `const res = await axios.get("/api/models", { params });`
// `models.value = res.data.map(mapModel);`
// And `filteredModels` used `models.value`.
// If the API handles filtering, `filteredModels` might be redundant or refining.
// However, `versionCards` (line 691) iterates over `filteredModels` and filters VERSIONS.
// The API returns models, which contain versions.
// The filters (like `selectedBaseModel`) might apply to versions, not just models.
// If the API filters models that HAVE a matching version, we still need to filter the versions to show only the matching ones in the card list.
// Yes, `versionCards` logic is crucial for flattening the list of versions.
// I need to move `versionCards` logic to `ModelList.vue` or `useModels`.
// Since it depends on the `models` data and filter states, it fits in `useModels` or `ModelList`.
// I'll put it in `ModelList` for now as it's view-specific (flattening for display), or `useModels` if I want to expose "cards".
// `useModels` has the filter states.
// I'll add `versionCards` to `ModelList.vue` using the data from `useModels`.

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
        // trainedWords parsing removed as ModelCard doesn't use it, 
        // and we need to preserve the reactive reference of 'v'.
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
