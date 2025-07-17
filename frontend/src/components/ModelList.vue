<template>
  <div class="mx-4">
    <div class="row">
      <div class="col-md-6 d-flex align-content-start flex-wrap gap-2">
        <input
          v-model="search"
          placeholder="Search models..."
          class="form-control w-200 flex-grow-1"
          style="min-width: 200px"
        />

        <input
          v-model="tagsSearch"
          placeholder="Search tags (comma separated)"
          class="form-control"
          style="min-width: 200px"
        />

        <div class="row">
          <div class="col-12 col-sm-6">
            <select
              v-model="selectedBaseModel"
              class="form-select"
              style="min-width: 200px"
            >
              <option value="">All base models</option>
              <option v-for="bm in baseModels" :key="bm" :value="bm">
                {{ bm }}
              </option>
            </select>
          </div>
          <div class="col-12 col-sm-6">
            <select
              v-model="selectedModelType"
              class="form-select"
              style="min-width: 200px"
            >
              <option value="">All model types</option>
              <option v-for="t in modelTypes" :key="t" :value="t">
                {{ t }}
              </option>
            </select>
          </div>
        </div>

        <div class="form-check form-switch d-flex gap-2 align-items-center m-2">
          <input 
            class="form-check-input"
            type="checkbox"
            role="switch"
            id="hide-nsfw"
            v-model="hideNsfw"
          />
          <label class="form-check-label" for="hide-nsfw"
            >Hide NSFW</label
          >
        </div>
      </div>
      <div class="col-md-6 d-flex align-content-start flex-wrap gap-2">
        <div class="input-group">
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

        <div class="input-group">
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
        <div v-if="downloading" class="progress w-100 mt-2">
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

  <div class="m-4 text-center" v-if="models.length === 0">No models found.</div>

  <div class="row row-cols-1 row-cols-md-3 row-cols-lg-5 g-4 p-4">
    <div v-for="card in versionCards" :key="card.version.ID" class="col">
      <div class="card h-100">
        <img
          v-if="card.imageUrl"
          :src="card.imageUrl"
          :width="card.model.imageWidth"
          :height="card.model.imageHeight"
          class="img-fluid card-img-top"
        />
        <div class="card-img-overlay z-2">
          <span class="badge rounded-pill text-bg-primary">{{
            card.version.type
          }}</span>
          <span class="ms-1 badge rounded-pill text-bg-success">{{
            card.version.baseModel
          }}</span>
        </div>
        <div class="card-body z-3">
          <h3 class="card-title h5">
            {{ card.model.name }} - {{ card.version.name }}
          </h3>
        </div>
        <div class="mb-2 d-flex gap-2 card-footer z-2">
          <button
            v-if="card.version.filePath"
            @click="goToModel(card.model.ID, card.version.ID)"
            class="btn btn-primary"
          >
            More details
          </button>
          <button
            @click="deleteVersion(card.version.ID)"
            class="btn btn-danger ms-auto"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
  <div v-if="hasMore" class="text-center mb-4">
    <button @click="loadMore" class="btn btn-secondary">Load More</button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";
import { showToast, showConfirm } from "../utils/ui";
import debounce from "../utils/debounce";

const models = ref([]);
const search = ref("");
const tagsSearch = ref("");
const selectedBaseModel = ref("");
const selectedModelType = ref("");
const hideNsfw = ref(false);
const modelUrl = ref("");
const versions = ref([]);
const selectedVersionId = ref("");
const loading = ref(false);
const downloading = ref(false);
const downloadProgress = ref(0);
let progressInterval = null;
const router = useRouter();

const page = ref(1);
const limit = 50;
const hasMore = ref(true);

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

const fetchModels = async (reset = false) => {
  const params = { page: page.value, limit };
  if (search.value) params.search = search.value;
  const res = await axios.get("/api/models", { params });
  const fetched = res.data.map(mapModel);
  if (reset) {
    models.value = fetched;
  } else {
    models.value = [...models.value, ...fetched];
  }
  hasMore.value = fetched.length === limit;
};

const debouncedFetchModels = debounce(() => fetchModels(true), 300);

onMounted(() => fetchModels(true));

watch(search, () => {
  page.value = 1;
  hasMore.value = true;
  debouncedFetchModels();
});

const baseModels = computed(() => {
  const set = new Set();
  models.value.forEach((m) => {
    (m.versions || []).forEach((v) => {
      if (v.baseModel) set.add(v.baseModel);
    });
  });
  return Array.from(set);
});

const modelTypes = computed(() => {
  const set = new Set();
  models.value.forEach((m) => {
    (m.versions || []).forEach((v) => {
      if (v.type) set.add(v.type);
    });
  });
  return Array.from(set);
});

const filteredModels = computed(() => {
  if (!search.value) return models.value;
  return models.value.filter((m) =>
    m.name.toLowerCase().includes(search.value.toLowerCase()),
  );
});

const versionCards = computed(() => {
  return filteredModels.value
    .flatMap((model) =>
      model.versions
        .filter((v) => {
          if (
            selectedBaseModel.value &&
            v.baseModel !== selectedBaseModel.value
          )
            return false;
          if (selectedModelType.value && v.type !== selectedModelType.value)
            return false;
          if (hideNsfw.value && v.nsfw) return false;

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
        }),
    )
    .sort((a, b) => b.version.ID - a.version.ID);
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
    await fetchModels(true);
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
  if (!(await showConfirm("Delete this version and all files?"))) return;
  await axios.delete(`/api/versions/${id}`);
  page.value = 1;
  await fetchModels(true);
};

const goToModel = (modelId, versionId) => {
  router.push(`/model/${modelId}/version/${versionId}`);
};

const loadMore = async () => {
  page.value += 1;
  await fetchModels();
};

</script>
