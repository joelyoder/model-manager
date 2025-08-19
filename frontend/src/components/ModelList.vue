<template>
  <div class="mx-4">
    <div class="row gap-2">
      <div class="col">
        <input
          v-model="search"
          placeholder="Search models..."
          class="form-control w-200 flex-grow-1"
          style="min-width: 200px"
        />
      </div>
      <div class="col">
        <input
          v-model="tagsSearch"
          placeholder="Search tags (comma separated)"
          class="form-control"
          style="min-width: 200px"
        />
      </div>
    </div>
    <div class="row gap-2 my-2">
      <div class="col">
        <select
          v-model="selectedCategory"
          class="form-select"
          style="min-width: 250px"
        >
          <option value="">All categories</option>
          <option v-for="cat in categories" :key="cat" :value="cat">
            {{ cat }}
          </option>
        </select>
      </div>
      <div class="col">
        <select
          v-model="selectedBaseModel"
          class="form-select"
          style="min-width: 250px"
        >
          <option value="">All base models</option>
          <option v-for="bm in baseModels" :key="bm" :value="bm">
            {{ bm }}
          </option>
        </select>
      </div>
      <div class="col">
        <select
          v-model="selectedModelType"
          class="form-select"
          style="min-width: 250px"
        >
          <option value="">All model types</option>
          <option v-for="t in modelTypes" :key="t" :value="t">
            {{ t }}
          </option>
        </select>
      </div>
      <div class="col-auto d-flex align-items-center">
        <button
          @click="hideNsfw = !hideNsfw"
          class="btn btn-outline-secondary btn-sm"
        >
          <svg
            v-if="hideNsfw"
            width="22px"
            height="22px"
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
            width="22px"
            height="22px"
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
      <div class="col d-flex justify-content-end">
        <button
          class="btn btn-outline-primary"
          @click="showAddPanel = !showAddPanel"
        >
          {{ showAddPanel ? "Close Panel" : "Add Models" }}
        </button>
      </div>
    </div>
    <div v-show="showAddPanel" class="mx-auto card card-body my-3" style="max-width: 1000px;">
      <div class="row g-3">
        <div class="col-md-2">
          <button
            @click="createManualModel"
            class="btn btn-primary w-100"
          >
            Add Model
          </button>
        </div>
        <div class="col">
          <div class="input-group mb-2">
            <!-- Paste URL and fetch versions -->
            <input
              v-model="modelUrl"
              placeholder="Paste CivitAI model URL"
              class="form-control"
              style="min-width: 200px"
              @keyup.enter="loadVersions"
            />
            <div class="input-group-append">
              <button
                @click="loadVersions"
                :disabled="loading || !modelUrl"
                class="btn btn-primary"
              >
                Load Versions
              </button>
            </div>
          </div>

          <div class="input-group mb-2">
            <!-- Version selector -->
            <select
              v-if="versions.length"
              v-model="selectedVersionId"
              class="form-select"
              style="min-width: 200px"
            >
              <option disabled value="">Select version</option>
              <option v-for="v in versions" :value="v.id" :key="v.id">
                {{ v.name }} | {{ v.baseModel }} |
                {{ ((v.sizeKB || 0) / 1024).toFixed(2) }} MB
              </option>
            </select>
            <div class="input-group-append">
              <!-- Download version -->
              <button
                v-if="selectedVersionId"
                @click="downloadSelectedVersion"
                :disabled="loading"
                class="btn btn-primary"
              >
                <span
                  v-if="loading"
                  class="spinner-border spinner-border-sm"
                  aria-hidden="true"
                ></span>
                <span v-if="loading" role="status" class="ps-2"
                  >Downloading...</span
                >
                <span v-else>Download</span>
              </button>
            </div>
          </div>
          <div v-if="downloading" class="progress w-100 mb-2">
            <div
              class="progress-bar progress-bar-striped"
              :style="{ width: downloadProgress + '%' }"
            >
              {{ downloadProgress }}%
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <nav v-if="totalPages > 1" class="mb-4">
    <ul class="pagination justify-content-center align-items-center gap-1">
      <li class="page-item" :class="{ disabled: page === 1 }">
        <a class="page-link" href="#" @click.prevent="changePage(1)">First</a>
      </li>
      <li class="page-item" :class="{ disabled: page === 1 }">
        <a class="page-link" href="#" @click.prevent="changePage(page - 1)"
          >Previous</a
        >
      </li>
      <li class="d-flex align-items-center">
        <input
          type="number"
          min="1"
          :max="totalPages"
          v-model.number="pageInput"
          @keyup.enter="goToPage"
          class="form-control"
          style="width: 80px"
        />
        <span class="ms-1">/ {{ totalPages }}</span>
      </li>
      <li class="page-item" :class="{ disabled: page === totalPages }">
        <a class="page-link" href="#" @click.prevent="changePage(page + 1)"
          >Next</a
        >
      </li>
      <li class="page-item" :class="{ disabled: page === totalPages }">
        <a class="page-link" href="#" @click.prevent="changePage(totalPages)"
          >Last</a
        >
      </li>
    </ul>
  </nav>

  <div class="m-4 text-center" v-if="models.length === 0">No models found.</div>

  <div class="model-grid p-4">
    <div
      v-for="card in versionCards"
      :key="card.version.ID"
      :id="`model-${card.version.ID}`"
      class="model-card card h-100"
    >
      <img
        v-if="card.imageUrl"
        :src="card.imageUrl"
        class="card-img-top"
        style="width: 100%; height: 450px; object-fit: cover"
      />
      <div class="card-img-overlay z-2">
        <span class="badge rounded-pill text-bg-primary">{{
          card.version.type
        }}</span>
        <span class="ms-1 badge rounded-pill text-bg-success">{{
          card.version.baseModel
        }}</span>
        <button
          @click.stop="toggleVersionNsfw(card.version)"
          class="btn btn-sm position-absolute top-0 end-0 m-2"
          :class="card.version.nsfw ? 'btn-danger' : 'btn-secondary'"
          style="--bs-btn-padding-y: 0.25rem; --bs-btn-padding-x: 0.25rem"
        >
          <svg
            v-if="card.version.nsfw"
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
      <div class="card-body z-3">
        <h3 class="card-title h5">
          {{ card.model.name }} - {{ card.version.name }}
        </h3>
      </div>
      <div class="d-flex gap-2 card-footer z-2">
        <button
          @click="goToModel(card.model.ID, card.version.ID)"
          class="btn btn-outline-primary"
        >
          More details
        </button>
        <button
          @click="deleteVersion(card.version.ID)"
          class="btn btn-outline-danger ms-auto"
        >
          Delete
        </button>
      </div>
    </div>
  </div>
  <nav v-if="totalPages > 1" class="mb-4">
    <ul class="pagination justify-content-center align-items-center gap-1">
      <li class="page-item" :class="{ disabled: page === 1 }">
        <a class="page-link" href="#" @click.prevent="changePage(1)">First</a>
      </li>
      <li class="page-item" :class="{ disabled: page === 1 }">
        <a class="page-link" href="#" @click.prevent="changePage(page - 1)"
          >Previous</a
        >
      </li>
      <li class="d-flex align-items-center">
        <input
          type="number"
          min="1"
          :max="totalPages"
          v-model.number="pageInput"
          @keyup.enter="goToPage"
          class="form-control"
          style="width: 80px"
        />
        <span class="ms-1">/ {{ totalPages }}</span>
      </li>
      <li class="page-item" :class="{ disabled: page === totalPages }">
        <a class="page-link" href="#" @click.prevent="changePage(page + 1)"
          >Next</a
        >
      </li>
      <li class="page-item" :class="{ disabled: page === totalPages }">
        <a class="page-link" href="#" @click.prevent="changePage(totalPages)"
          >Last</a
        >
      </li>
    </ul>
  </nav>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";
import { showToast, showDeleteConfirm } from "../utils/ui";
import debounce from "../utils/debounce";

const models = ref([]);
const search = ref("");
const tagsSearch = ref("");
const selectedCategory = ref("");
const selectedBaseModel = ref("");
const selectedModelType = ref("");
const hideNsfw = ref(false);
const showAddPanel = ref(false);
const modelUrl = ref("");
const versions = ref([]);
const selectedVersionId = ref("");
const loading = ref(false);
const downloading = ref(false);
const downloadProgress = ref(0);
let progressInterval = null;
const router = useRouter();
const route = useRoute();

const page = ref(1);
const pageInput = ref(1);
const limit = 50;
const total = ref(0);
const localStorageKey = "modelListState";

const mapModel = (model) => {
  const imageUrl = model.imagePath
    ? model.imagePath.replace(/^.*[\\/]backend[\\/]images/, "/images")
    : null;
  const versions = (model.versions || []).map((v) => {
    const vImage = v.imagePath
      ? v.imagePath.replace(/^.*[\\/]backend[\\/]images/, "/images")
      : null;
    return { ...v, imageUrl: vImage };
  });
  return {
    ...model,
    versions,
    imageUrl,
  };
};

const fetchModels = async () => {
  const params = { page: page.value, limit, includeVersions: 1 };
  if (search.value) params.search = search.value;
  if (selectedBaseModel.value) params.baseModel = selectedBaseModel.value;
  if (selectedModelType.value) params.modelType = selectedModelType.value;
  if (hideNsfw.value) params.hideNsfw = 1;
  const tagParts = [];
  if (selectedCategory.value) tagParts.push(selectedCategory.value);
  if (tagsSearch.value.trim()) tagParts.push(tagsSearch.value);
  if (tagParts.length) params.tags = tagParts.join(",");
  const res = await axios.get("/api/models", { params });
  models.value = res.data.map(mapModel);
};

const fetchTotal = async () => {
  const params = {};
  if (search.value) params.search = search.value;
  if (selectedBaseModel.value) params.baseModel = selectedBaseModel.value;
  if (selectedModelType.value) params.modelType = selectedModelType.value;
  if (hideNsfw.value) params.hideNsfw = 1;
  const tagParts = [];
  if (selectedCategory.value) tagParts.push(selectedCategory.value);
  if (tagsSearch.value.trim()) tagParts.push(tagsSearch.value);
  if (tagParts.length) params.tags = tagParts.join(",");
  const res = await axios.get("/api/models/count", { params });
  total.value = res.data.count || 0;
};

const debouncedUpdate = debounce(async () => {
  page.value = 1;
  await fetchTotal();
  await fetchModels();
}, 300);

const initialized = ref(false);

onMounted(async () => {
  const saved = JSON.parse(localStorage.getItem(localStorageKey) || "{}");
  if (saved.search !== undefined) search.value = saved.search;
  if (saved.tagsSearch !== undefined) tagsSearch.value = saved.tagsSearch;
  if (saved.selectedCategory !== undefined)
    selectedCategory.value = saved.selectedCategory;
  if (saved.selectedBaseModel !== undefined)
    selectedBaseModel.value = saved.selectedBaseModel;
  if (saved.selectedModelType !== undefined)
    selectedModelType.value = saved.selectedModelType;
  if (saved.hideNsfw !== undefined) hideNsfw.value = saved.hideNsfw;
  if (saved.page !== undefined) page.value = saved.page;

  await fetchBaseModels();
  await fetchTotal();
  await fetchModels();
  initialized.value = true;
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

onUnmounted(() => {
  localStorage.setItem(
    localStorageKey,
    JSON.stringify({
      search: search.value,
      tagsSearch: tagsSearch.value,
      selectedCategory: selectedCategory.value,
      selectedBaseModel: selectedBaseModel.value,
      selectedModelType: selectedModelType.value,
      hideNsfw: hideNsfw.value,
      page: page.value,
    }),
  );
});

watch(search, () => {
  if (initialized.value) debouncedUpdate();
});

watch(tagsSearch, () => {
  if (initialized.value) debouncedUpdate();
});

watch(selectedCategory, () => {
  if (initialized.value) debouncedUpdate();
});

watch(selectedBaseModel, () => {
  if (initialized.value) debouncedUpdate();
});

watch(selectedModelType, () => {
  if (initialized.value) debouncedUpdate();
});

watch(hideNsfw, async () => {
  if (!initialized.value) return;
  page.value = 1;
  await fetchTotal();
  await fetchModels();
});

watch(page, () => {
  pageInput.value = page.value;
});

const baseModels = ref([]);
const fetchBaseModels = async () => {
  try {
    const res = await axios.get("/api/base-models");
    baseModels.value = Array.isArray(res.data) ? res.data : [];
  } catch {
    baseModels.value = [];
  }
};

const modelTypes = [
  "Checkpoint",
  "TextualInversion",
  "Hypernetwork",
  "AestheticGradient",
  "LORA",
  "LoCon",
  "DoRA",
  "Controlnet",
  "Upscaler",
  "MotionModule",
  "VAE",
  "Wildcards",
  "Poses",
  "Workflows",
  "Detection",
  "Other",
];

const categories = [
  "character",
  "style",
  "concept",
  "clothing",
  "base model",
  "poses",
  "background",
  "tool",
  "vehicle",
  "buildings",
  "objects",
  "assets",
  "animal",
  "action",
];

const totalPages = computed(() => Math.ceil(total.value / limit));

const filteredModels = computed(() => {
  return models.value.filter((m) => {
    if (hideNsfw.value && m.nsfw) return false;
    if (search.value && !m.name.toLowerCase().includes(search.value.toLowerCase()))
      return false;
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
        if (hideNsfw.value && v.nsfw) return false;

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
        let trained = v.trainedWords;
        if (typeof trained === "string") {
          trained = trained ? trained.split(",") : [];
        }
        return {
          model,
          version: { ...v, trainedWordsArr: trained },
          imageUrl: v.imageUrl || model.imageUrl,
        };
      });
  });
});

const extractModelId = (url) => {
  const match = url.match(/models\/(\d+)/);
  return match ? match[1] : null;
};

const extractVersionId = (url) => {
  const match = url.match(/modelVersionId=(\d+)/);
  return match ? match[1] : null;
};

const loadVersions = async () => {
  const id = extractModelId(modelUrl.value);
  if (!id) {
    showToast("Invalid CivitAI model URL", "danger");
    return;
  }

  loading.value = true;
  try {
    const res = await axios.get(`/api/model/${id}/versions`);
    versions.value = res.data;
    const vid = extractVersionId(modelUrl.value);
    if (vid && res.data.some((v) => v.id === Number(vid))) {
      selectedVersionId.value = vid;
    } else if (res.data.length) {
      selectedVersionId.value = String(res.data[0].id);
    } else {
      selectedVersionId.value = "";
    }
  } catch (err) {
    console.error(err);
    showToast("Failed to load versions", "danger");
  } finally {
    loading.value = false;
  }
};

const downloadSelectedVersion = async () => {
  if (!selectedVersionId.value) return;

  const ver = versions.value.find(
    (v) => v.id === Number(selectedVersionId.value),
  );
  if (ver) {
    const duplicate = models.value.some((m) =>
      (m.versions || []).some(
        (existing) => existing.sha256 && existing.sha256 === ver.sha256,
      ),
    );
    if (duplicate) {
      showToast("Model with the same hash already exists", "warning");
    }
  }

  loading.value = true;
  downloading.value = true;
  downloadProgress.value = 0;
  progressInterval = setInterval(async () => {
    const res = await axios.get("/api/download/progress");
    downloadProgress.value = res.data.progress || 0;
  }, 500);
  try {
    await axios.post(`/api/sync/version/${selectedVersionId.value}`);
    page.value = 1;
    await fetchTotal();
    await fetchModels();
    await fetchBaseModels();
    showToast("Version downloaded successfully", "success");
  } catch (err) {
    console.error(err);
    showToast("Download failed", "danger");
  } finally {
    clearInterval(progressInterval);
    progressInterval = null;
    downloadProgress.value = 0;
    modelUrl.value = "";
    versions.value = [];
    selectedVersionId.value = "";
    loading.value = false;
    downloading.value = false;
  }
};

const deleteVersion = async (id) => {
  const choice = await showDeleteConfirm("Delete this version?");
  if (!choice) return;
  const files = choice === "deleteFiles" ? 1 : 0;
  await axios.delete(`/api/versions/${id}?files=${files}`);
  page.value = 1;
  await fetchTotal();
  await fetchModels();
  await fetchBaseModels();
};

const toggleVersionNsfw = async (version) => {
  const updated = { ...version, nsfw: !version.nsfw };
  try {
    await axios.put(`/api/versions/${version.ID}`, updated);
    // update the card version
    version.nsfw = updated.nsfw;
    // update underlying model data so computed cards react immediately
    for (const m of models.value) {
      const v = m.versions.find((v) => v.ID === version.ID);
      if (v) {
        v.nsfw = updated.nsfw;
        break;
      }
    }
    showToast("NSFW status updated", "success");
  } catch (err) {
    console.error(err);
    showToast("Failed to update NSFW status", "danger");
  }
};

const goToModel = (modelId, versionId) => {
  router.push(`/model/${modelId}/version/${versionId}`);
};

const changePage = async (p) => {
  if (p < 1 || p > totalPages.value) return;
  page.value = p;
  await fetchModels();
};

const goToPage = async () => {
  await changePage(pageInput.value);
};

const createManualModel = async () => {
  try {
    const res = await axios.post("/api/models");
    const { modelId, versionId } = res.data;
    router.push(`/model/${modelId}/version/${versionId}?edit=1`);
  } catch (err) {
    console.error(err);
    showToast("Failed to create model", "danger");
  }
};
</script>
