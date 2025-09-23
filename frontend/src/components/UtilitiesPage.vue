<template>
  <div class="container px-4">
    <div class="mb-2 d-flex gap-2 pb-2">
      <button @click="goBack" class="btn btn-secondary">Back</button>
    </div>
    <h2 class="my-3">Utilities</h2>
    <div class="card card-body mb-4" v-if="stats">
      <h3>Stats</h3>
      <div class="row g-2 mb-3">
        <div class="col-md-3">
          <label class="form-label mb-1">Category</label>
          <select v-model="selectedCategory" class="form-select">
            <option value="">All categories</option>
            <option v-for="cat in categories" :key="cat" :value="cat">
              {{ cat }}
            </option>
          </select>
        </div>
        <div class="col-md-3">
          <label class="form-label mb-1">Base Model</label>
          <select v-model="selectedBaseModel" class="form-select">
            <option value="">All base models</option>
            <option v-for="bm in baseModels" :key="bm" :value="bm">
              {{ bm }}
            </option>
          </select>
        </div>
        <div class="col-md-3">
          <label class="form-label mb-1">Model Type</label>
          <select v-model="selectedModelType" class="form-select">
            <option value="">All model types</option>
            <option v-for="type in modelTypes" :key="type" :value="type">
              {{ type }}
            </option>
          </select>
        </div>
        <div class="col-md-2 d-flex align-items-end">
          <div class="form-check mb-0">
            <input
              id="stats-hide-nsfw"
              class="form-check-input"
              type="checkbox"
              v-model="hideNsfw"
            />
            <label class="form-check-label" for="stats-hide-nsfw">
              Hide NSFW
            </label>
          </div>
        </div>
        <div class="col-md-1 d-flex align-items-end justify-content-end">
          <button @click="clearFilters" class="btn btn-outline-secondary">
            Clear
          </button>
        </div>
      </div>
      <p class="text-center h5 mb-3">
        Total Models: <strong>{{ stats.totalModels }}</strong>
        <br />
        Total Versions: <strong>{{ stats.totalVersions }}</strong>
      </p>
      <div class="row row-cols-1 row-cols-md-2 text-center g-3">
        <div class="col">
          <canvas id="typeChart"></canvas>
        </div>
        <div class="col">
          <canvas id="baseModelChart"></canvas>
        </div>
        <div class="col">
          <canvas id="categoryChart"></canvas>
        </div>
        <div class="col">
          <canvas id="nsfwChart"></canvas>
        </div>
      </div>
    </div>
    <div class="card card-body mb-4">
      <h3>Import & Export</h3>
      <h4 class="h5 my-3">Import JSON from Model Organizer</h4>
      <div class="input-group mb-3">
        <input
          type="file"
          accept=".json"
          @change="onFileChange"
          class="form-control"
        />
        <div class="input-group-append">
          <button
            @click="importJson"
            :disabled="!importFile"
            class="btn btn-primary"
          >
            Import
          </button>
        </div>
      </div>
      <div class="d-flex gap-2 mb-3">
        <span>Update:</span>
        <div class="form-check">
          <input
            class="form-check-input"
            type="checkbox"
            id="ie-pull-images"
            v-model="pullImages"
          />
          <label class="form-check-label" for="ie-pull-images">Images</label>
        </div>
        <div class="form-check">
          <input
            class="form-check-input"
            type="checkbox"
            id="ie-pull-meta"
            v-model="pullMeta"
          />
          <label class="form-check-label" for="ie-pull-meta">Metadata</label>
        </div>
        <div class="form-check">
          <input
            class="form-check-input"
            type="checkbox"
            id="ie-pull-desc"
            v-model="pullDesc"
          />
          <label class="form-check-label" for="ie-pull-desc">Description</label>
        </div>
      </div>
      <h4 class="h5 my-3">Export Database as JSON</h4>
      <div class="mb-3 d-flex gap-2">
        <button @click="exportJson" class="btn btn-primary">
          Export Models
        </button>
      </div>
      <h4 class="h5 my-3">Import Database from JSON</h4>
      <div class="input-group mb-3">
        <input
          type="file"
          accept=".json"
          @change="onDbFileChange"
          class="form-control"
        />
        <div class="input-group-append">
          <button
            @click="importDbJson"
            :disabled="!dbImportFile"
            class="btn btn-primary"
          >
            Import
          </button>
        </div>
      </div>
    </div>
    <div class="card card-body">
      <h3>Library Cleanup</h3>
      <h4 class="h5 my-3">Find Orphaned Model Files</h4>
      <div class="mb-3 flex gap-2">
        <button @click="findOrphanFiles" class="btn btn-primary">
          Search Library
        </button>
        <div v-if="orphanFiles.length">
          <ul class="list-group list-group-flush">
            <li v-for="file in orphanFiles" :key="file" class="list-group-item">
              {{ file }}
            </li>
          </ul>
          <button @click="exportOrphanFiles" class="btn btn-secondary mt-3">
            Export Results
          </button>
        </div>
        <p v-else-if="searchDone" class="mb-0">No orphaned files found</p>
      </div>
      <h4 class="h5 my-3">Find Duplicate File Paths</h4>
      <div class="mb-3 flex gap-2">
        <button @click="findDuplicatePaths" class="btn btn-primary">
          Search Duplicates
        </button>
        <div v-if="duplicatePaths.length">
          <ul class="list-group list-group-flush">
            <li
              v-for="dup in duplicatePaths"
              :key="dup.path"
              class="list-group-item"
            >
              <strong>{{ dup.path }}</strong>
              <ul class="mb-0">
                <li v-for="v in dup.versions" :key="v.versionId">
                  {{ v.modelName }} - {{ v.versionName }}
                </li>
              </ul>
            </li>
          </ul>
          <button @click="exportDuplicatePaths" class="btn btn-secondary mt-3">
            Export Results
          </button>
        </div>
        <p v-else-if="dupSearchDone" class="mb-0">
          No duplicate file paths found
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
/* global Chart */
import { ref, onMounted, nextTick, watch } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";
import { showToast } from "../utils/ui";

const stats = ref(null);
let typeChart = null;
let baseChart = null;
let nsfwChart = null;
let categoryChart = null;
let statsRequestId = 0;

const selectedCategory = ref("");
const selectedBaseModel = ref("");
const selectedModelType = ref("");
const hideNsfw = ref(false);
const baseModels = ref([]);

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

onMounted(async () => {
  await fetchBaseModels();
  await fetchStats();
});

watch([selectedCategory, selectedBaseModel, selectedModelType, hideNsfw], () => {
  fetchStats();
});

const fetchBaseModels = async () => {
  try {
    const res = await axios.get("/api/base-models");
    baseModels.value = Array.isArray(res.data) ? res.data : [];
  } catch (err) {
    console.error(err);
    baseModels.value = [];
  }
};

const fetchStats = async () => {
  const requestId = ++statsRequestId;
  try {
    const params = new URLSearchParams();
    if (selectedCategory.value) params.set("category", selectedCategory.value);
    if (selectedBaseModel.value)
      params.set("baseModel", selectedBaseModel.value);
    if (selectedModelType.value)
      params.set("modelType", selectedModelType.value);
    if (hideNsfw.value) params.set("hideNsfw", "1");

    const query = params.toString();
    const url = query ? `/api/stats?${query}` : "/api/stats";
    const res = await axios.get(url);
    if (requestId !== statsRequestId) {
      return;
    }
    stats.value = res.data;
    await nextTick();
    renderCharts();
  } catch (err) {
    console.error(err);
    if (requestId === statsRequestId) {
      stats.value = null;
      destroyCharts();
    }
  }
};

function renderCharts() {
  if (!stats.value) return;
  destroyCharts();

  const typeCounts = stats.value.typeCounts || [];
  const baseCounts = stats.value.baseModelCounts || [];
  const categoryCounts = stats.value.categoryCounts || [];

  const typeCtx = document.getElementById("typeChart");
  if (typeCtx) {
    typeChart = new Chart(typeCtx, {
      type: "pie",
      data: {
        labels: typeCounts.map((t) => t.Key || "Unknown"),
        datasets: [
          {
            data: typeCounts.map((t) => t.Count),
          },
        ],
      },
    });
  }

  const baseCtx = document.getElementById("baseModelChart");
  if (baseCtx) {
    baseChart = new Chart(baseCtx, {
      type: "pie",
      data: {
        labels: baseCounts.map((b) => b.Key || "Unknown"),
        datasets: [
          {
            data: baseCounts.map((b) => b.Count),
          },
        ],
      },
    });
  }

  const categoryCtx = document.getElementById("categoryChart");
  if (categoryCtx) {
    categoryChart = new Chart(categoryCtx, {
      type: "pie",
      data: {
        labels: categoryCounts.map((c) => c.Key || "Unknown"),
        datasets: [
          {
            data: categoryCounts.map((c) => c.Count),
          },
        ],
      },
    });
  }

  const nsfwCtx = document.getElementById("nsfwChart");
  if (nsfwCtx) {
    nsfwChart = new Chart(nsfwCtx, {
      type: "pie",
      data: {
        labels: ["Non-NSFW", "NSFW"],
        datasets: [
          {
            data: [stats.value.nonNsfwCount, stats.value.nsfwCount],
          },
        ],
      },
    });
  }
}

function destroyCharts() {
  if (typeChart) {
    typeChart.destroy();
    typeChart = null;
  }
  if (baseChart) {
    baseChart.destroy();
    baseChart = null;
  }
  if (categoryChart) {
    categoryChart.destroy();
    categoryChart = null;
  }
  if (nsfwChart) {
    nsfwChart.destroy();
    nsfwChart = null;
  }
}

const clearFilters = () => {
  if (
    !selectedCategory.value &&
    !selectedBaseModel.value &&
    !selectedModelType.value &&
    !hideNsfw.value
  ) {
    return;
  }
  selectedCategory.value = "";
  selectedBaseModel.value = "";
  selectedModelType.value = "";
  hideNsfw.value = false;
};

const importFile = ref(null);
const dbImportFile = ref(null);
const pullImages = ref(false);
const pullMeta = ref(false);
const pullDesc = ref(false);
const router = useRouter();
const orphanFiles = ref([]);
const searchDone = ref(false);
const duplicatePaths = ref([]);
const dupSearchDone = ref(false);

const onFileChange = (e) => {
  importFile.value = e.target.files[0] || null;
};

const onDbFileChange = (e) => {
  dbImportFile.value = e.target.files[0] || null;
};

const importJson = async () => {
  if (!importFile.value) return;
  const form = new FormData();
  form.append("file", importFile.value);
  try {
    const params = [];
    if (pullMeta.value) params.push("metadata");
    if (pullDesc.value) params.push("description");
    if (pullImages.value) params.push("images");
    const query = params.length ? `?fields=${params.join(",")}` : "";
    await axios.post(`/api/import${query}`, form);
    showToast("Import successful", "success");
  } catch (err) {
    console.error(err);
    showToast("Import failed", "danger");
  } finally {
    importFile.value = null;
  }
};

const importDbJson = async () => {
  if (!dbImportFile.value) return;
  const form = new FormData();
  form.append("file", dbImportFile.value);
  try {
    await axios.post("/api/import-db", form);
    showToast("Database import successful", "success");
  } catch (err) {
    console.error(err);
    showToast("Database import failed", "danger");
  } finally {
    dbImportFile.value = null;
  }
};

const exportJson = async () => {
  try {
    const res = await axios.get("/api/export", { responseType: "blob" });
    const url = window.URL.createObjectURL(
      new Blob([res.data], { type: "application/json" }),
    );
    const a = document.createElement("a");
    a.href = url;
    a.download = "model_export.json";
    a.click();
    window.URL.revokeObjectURL(url);
  } catch (err) {
    console.error(err);
    showToast("Export failed", "danger");
  }
};

const findOrphanFiles = async () => {
  try {
    const res = await axios.get("/api/orphaned-files");
    console.log("orphaned files response", res.data);
    orphanFiles.value = res.data.orphans || [];
  } catch (err) {
    console.error(err);
    showToast("Failed to fetch orphaned files", "danger");
    orphanFiles.value = [];
  } finally {
    searchDone.value = true;
  }
};

const findDuplicatePaths = async () => {
  try {
    const res = await axios.get("/api/duplicate-file-paths");
    duplicatePaths.value = res.data.duplicates || [];
  } catch (err) {
    console.error(err);
    showToast("Failed to fetch duplicate paths", "danger");
    duplicatePaths.value = [];
  } finally {
    dupSearchDone.value = true;
  }
};

const exportDuplicatePaths = () => {
  if (!duplicatePaths.value.length) return;
  const lines = duplicatePaths.value.map((d) =>
    [
      d.path,
      ...d.versions.map((v) => `  ${v.modelName} - ${v.versionName}`),
    ].join("\n"),
  );
  const blob = new Blob([lines.join("\n\n")], { type: "text/plain" });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "duplicate_file_paths.txt";
  a.click();
  window.URL.revokeObjectURL(url);
};

const exportOrphanFiles = () => {
  if (!orphanFiles.value.length) return;
  const blob = new Blob([orphanFiles.value.join("\n")], {
    type: "text/plain",
  });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "orphaned_files.txt";
  a.click();
  window.URL.revokeObjectURL(url);
};

const goBack = () => {
  router.push("/");
};
</script>
