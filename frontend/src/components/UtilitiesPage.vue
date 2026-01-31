<template>
  <div class="container px-2 px-md-4 max-w-4xl mx-auto">
    <!-- Header -->
    <div class="mb-4 d-flex gap-2 align-items-center px-2 px-md-0">
      <button 
        @click="goBack" 
        class="btn btn-outline-secondary btn-sm d-flex align-items-center justify-content-center border-0"
        aria-label="Back"
        title="Back"
        style="width: 40px; height: 40px;"
      >
        <Icon icon="mdi:arrow-left" width="24" height="24" />
      </button>
      <h2 class="h5 mb-0 fw-bold ms-2">Utilities</h2>
    </div>

    <!-- Stats Card -->
    <div class="card border-0 shadow-sm bg-dark-subtle rounded-3 overflow-hidden mb-4" v-if="stats">
      <div class="card-body p-4">
          <h3 class="h6 fw-bold text-uppercase text-secondary mb-3">Library Stats</h3>
          
          <div class="row g-2 mb-4">
            <div class="col-md-3">
              <label class="form-label small text-secondary fw-bold mb-1">Category</label>
              <select v-model="selectedCategory" class="form-select bg-dark border-0 text-white shadow-none">
                <option value="">All categories</option>
                <option v-for="cat in categories" :key="cat" :value="cat">
                  {{ cat }}
                </option>
              </select>
            </div>
            <div class="col-md-3">
              <label class="form-label small text-secondary fw-bold mb-1">Base Model</label>
              <select v-model="selectedBaseModel" class="form-select bg-dark border-0 text-white shadow-none">
                <option value="">All base models</option>
                <option v-for="bm in baseModels" :key="bm" :value="bm">
                  {{ bm }}
                </option>
              </select>
            </div>
            <div class="col-md-3">
              <label class="form-label small text-secondary fw-bold mb-1">Model Type</label>
              <select v-model="selectedModelType" class="form-select bg-dark border-0 text-white shadow-none">
                <option value="">All model types</option>
                <option v-for="type in modelTypes" :key="type" :value="type">
                  {{ type }}
                </option>
              </select>
            </div>
            <div class="col-md-3">
              <label class="form-label small text-secondary fw-bold mb-1" for="stats-nsfw-filter">NSFW</label>
              <select
                id="stats-nsfw-filter"
                v-model="nsfwFilter"
                class="form-select bg-dark border-0 text-white shadow-none"
              >
                <option value="">NSFW &amp; Non-NSFW</option>
                <option value="non">Non-NSFW only</option>
                <option value="nsfw">NSFW only</option>
              </select>
            </div>
            <div class="col-md-12 col-lg-3 d-flex align-items-end justify-content-end ms-auto mt-2">
              <button @click="clearFilters" class="btn btn-outline-secondary btn-sm">
                Clear Filters
              </button>
            </div>
          </div>

          <div class="text-center mb-4 p-3 bg-dark bg-opacity-25 rounded-3">
              <div class="row">
                  <div class="col-6 border-end border-secondary border-opacity-25">
                      <div class="h3 fw-bold mb-0">{{ stats.totalModels }}</div>
                      <div class="small text-secondary text-uppercase tracking-wide">Total Models</div>
                  </div>
                  <div class="col-6">
                       <div class="h3 fw-bold mb-0">{{ stats.totalVersions }}</div>
                      <div class="small text-secondary text-uppercase tracking-wide">Total Versions</div>
                  </div>
              </div>
          </div>

          <div class="row row-cols-1 row-cols-md-2 text-center g-4">
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
    </div>

    <!-- Import / Export Card -->
    <div class="card border-0 shadow-sm bg-dark-subtle rounded-3 overflow-hidden mb-4">
      <div class="card-body p-4">
        <h3 class="h6 fw-bold text-uppercase text-secondary mb-3">Import & Export</h3>
        
        <h4 class="h6 fw-bold mt-4 mb-2">Import JSON from Model Organizer</h4>
        <div class="input-group mb-3">
          <input
            type="file"
            accept=".json"
            @change="onFileChange"
            class="form-control bg-dark border-0 text-white shadow-none"
          />
          <button
            @click="importJson"
            :disabled="!importFile"
            class="btn btn-primary"
          >
            Import
          </button>
        </div>

        <div class="d-flex gap-4 mb-3 ps-1">
          <span class="text-secondary small fw-bold text-uppercase pt-1">Update:</span>
          <div class="form-check">
            <input
              class="form-check-input"
              type="checkbox"
              id="ie-pull-images"
              v-model="pullImages"
            />
            <label class="form-check-label small" for="ie-pull-images">Images</label>
          </div>
          <div class="form-check">
            <input
              class="form-check-input"
              type="checkbox"
              id="ie-pull-meta"
              v-model="pullMeta"
            />
            <label class="form-check-label small" for="ie-pull-meta">Metadata</label>
          </div>
          <div class="form-check">
            <input
              class="form-check-input"
              type="checkbox"
              id="ie-pull-desc"
              v-model="pullDesc"
            />
            <label class="form-check-label small" for="ie-pull-desc">Description</label>
          </div>
        </div>
        
        <hr class="border-secondary opacity-25 my-4" />

        <h4 class="h6 fw-bold mb-2">Database Backup</h4>
        <div class="row g-3">
             <div class="col-md-6">
                <div class="p-3 bg-dark bg-opacity-25 rounded-3 h-100">
                    <label class="form-label small text-secondary fw-bold text-uppercase">Export Database JSON</label>
                    <button @click="exportJson" class="btn btn-outline-primary w-100 mt-1">
                        Export Models
                    </button>
                </div>
             </div>
              <div class="col-md-6">
                <div class="p-3 bg-dark bg-opacity-25 rounded-3 h-100">
                     <label class="form-label small text-secondary fw-bold text-uppercase">Import Database JSON</label>
                    <div class="input-group mt-1">
                        <input
                        type="file"
                        accept=".json"
                        @change="onDbFileChange"
                        class="form-control bg-dark border-0 text-white shadow-none form-control-sm"
                        />
                        <button
                        @click="importDbJson"
                        :disabled="!dbImportFile"
                        class="btn btn-primary btn-sm"
                        >
                        Import
                        </button>
                    </div>
                </div>
             </div>
        </div>
      </div>
    </div>

    <!-- Maintenance Card -->
    <div class="card border-0 shadow-sm bg-dark-subtle rounded-3 overflow-hidden mb-4">
      <div class="card-body p-4">
        <h3 class="h6 fw-bold text-uppercase text-secondary mb-3">Database Maintenance</h3>
        
        <div class="row g-4">
            <div class="col-md-6">
                 <h4 class="h6 fw-bold">Migrate Paths</h4>
                <p class="text-secondary small mb-3">
                    Convert absolute paths in the database to relative paths based on the current
                    configured Model and Image paths. This is useful for making your library portable.
                </p>
                <button
                @click="migratePaths"
                class="btn btn-warning btn-sm"
                :disabled="migrating"
                >
                {{ migrating ? "Migrating..." : "Migrate Paths to Relative" }}
                </button>
            </div>

             <div class="col-md-6 border-start border-secondary border-opacity-10">
                <h4 class="h6 fw-bold">Archive Description Images</h4>
                <p class="text-secondary small mb-3">
                    Download external images found in model descriptions and store them locally
                    to prevent broken links. This operation modifies model descriptions.
                </p>
                <button
                @click="archiveImages"
                class="btn btn-warning btn-sm"
                :disabled="isArchiving"
                >
                {{ isArchiving ? "Archiving..." : "Archive Images" }}
                </button>
                 <p v-if="archiveResult" class="mt-2 mb-0 small text-success">
                    {{ archiveResult }}
                </p>
                
                <hr class="border-secondary opacity-25 my-4" />
                
                <h4 class="h6 fw-bold">Reset Stuck Models</h4>
                <p class="text-secondary small mb-3">
                    If models are stuck in a "synching" state due to connection interruptions,
                    use this to reset their status to try again.
                </p>
                <button
                @click="resetPendingStatus"
                class="btn btn-danger btn-sm"
                :disabled="isResetting"
                >
                {{ isResetting ? "Resetting..." : "Reset Stuck Models" }}
                </button>
                
                <hr class="border-secondary opacity-25 my-4" />

                <h4 class="h6 fw-bold">Generate Thumbnails</h4>
                <p class="text-secondary small mb-3">
                    Generate WebP thumbnails for all models to improve performance. 
                    This will process any models missing thumbnails.
                </p>
                 <button
                @click="generateThumbnails"
                class="btn btn-primary btn-sm"
                :disabled="isGeneratingThumbs"
                >
                {{ isGeneratingThumbs ? "Generating..." : "Generate Missing Thumbnails" }}
                </button>
            </div>
        </div>
      </div>
    </div>

    <!-- Cleanup Card -->
    <div class="card border-0 shadow-sm bg-dark-subtle rounded-3 overflow-hidden mb-4">
      <div class="card-body p-4">
        <h3 class="h6 fw-bold text-uppercase text-secondary mb-3">Library Cleanup</h3>
        
        <div class="mb-4">
             <h4 class="h6 fw-bold">Find Orphaned Model Files</h4>
             <div class="d-flex gap-2 align-items-center mb-2">
                 <button @click="findOrphanFiles" class="btn btn-outline-primary btn-sm">
                    Search Library
                </button>
                <span v-if="searchDone && !orphanFiles.length" class="text-success small">No orphaned files found</span>
             </div>
            
             <div v-if="orphanFiles.length" class="mt-3 p-3 bg-dark bg-opacity-25 rounded-3">
                <div class="d-flex justify-content-between align-items-center mb-2">
                     <span class="fw-bold text-danger">{{ orphanFiles.length }} orphaned files found</span>
                     <button @click="exportOrphanFiles" class="btn btn-secondary btn-sm">
                        Export Results
                    </button>
                </div>
                <div class="overflow-auto" style="max-height: 150px;">
                     <ul class="list-group list-group-flush small">
                        <li v-for="file in orphanFiles" :key="file" class="list-group-item bg-transparent text-white-50 py-1 px-0 border-secondary border-opacity-25">
                        {{ file }}
                        </li>
                    </ul>
                </div>
            </div>
        </div>

        <hr class="border-secondary opacity-25 my-4" />

        <div>
            <h4 class="h6 fw-bold">Find Duplicate File Paths</h4>
             <div class="d-flex gap-2 align-items-center mb-2">
                 <button @click="findDuplicatePaths" class="btn btn-outline-primary btn-sm">
                    Search Duplicates
                </button>
                 <span v-if="dupSearchDone && !duplicatePaths.length" class="text-success small">No duplicate file paths found</span>
             </div>

             <div v-if="duplicatePaths.length" class="mt-3 p-3 bg-dark bg-opacity-25 rounded-3">
                 <div class="d-flex justify-content-between align-items-center mb-2">
                     <span class="fw-bold text-warning">{{ duplicatePaths.length }} duplicate paths found</span>
                     <button @click="exportDuplicatePaths" class="btn btn-secondary btn-sm">
                        Export Results
                    </button>
                </div>
                <div class="overflow-auto" style="max-height: 200px;">
                    <ul class="list-group list-group-flush small">
                        <li
                        v-for="dup in duplicatePaths"
                        :key="dup.path"
                        class="list-group-item bg-transparent py-2 px-0 border-secondary border-opacity-25"
                        >
                        <strong class="text-white d-block mb-1">{{ dup.path }}</strong>
                        <ul class="mb-0 ps-3 text-white-50">
                            <li v-for="v in dup.versions" :key="v.versionId">
                            {{ v.modelName }} - {{ v.versionName }}
                            </li>
                        </ul>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
/* global Chart */
import { ref, onMounted, nextTick, watch } from "vue";
import { Icon } from "@iconify/vue";
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
const nsfwFilter = ref("");
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

watch(
  [selectedCategory, selectedBaseModel, selectedModelType, nsfwFilter],
  () => {
    fetchStats();
  },
);

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
    if (nsfwFilter.value) params.set("nsfw", nsfwFilter.value);

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
    !nsfwFilter.value
  ) {
    return;
  }
  selectedCategory.value = "";
  selectedBaseModel.value = "";
  selectedModelType.value = "";
  nsfwFilter.value = "";
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

const migrating = ref(false);
const isArchiving = ref(false);
const archiveResult = ref("");
const isResetting = ref(false);

const resetPendingStatus = async () => {
  if (!confirm("This will reset the status of all models currently stuck in 'synching'. Continue?")) {
    return;
  }
  isResetting.value = true;
  try {
    await axios.post("/api/tools/reset-pending");
    showToast("Status reset successful", "success");
  } catch (err) {
    console.error(err);
    showToast("Failed to reset status", "danger");
  } finally {
    isResetting.value = false;
  }
};

const isGeneratingThumbs = ref(false);

const generateThumbnails = async () => {
    isGeneratingThumbs.value = true;
    try {
        await axios.post("/api/tools/generate-thumbnails");
        showToast("Thumbnail generation started in background", "success");
    } catch (err) {
        console.error(err);
        showToast("Failed to start thumbnail generation", "danger");
    } finally {
        isGeneratingThumbs.value = false;
    }
};

const archiveImages = async () => {
  if (!confirm("This will download external images and modify model descriptions. Continue?")) {
    return;
  }
  isArchiving.value = true;
  archiveResult.value = "";
  try {
    const res = await axios.post("/api/tools/archive-images");
    const { processed, updated } = res.data;
    archiveResult.value = `Scanned ${processed} versions, updated ${updated}.`;
    showToast("Archive complete", "success");
  } catch (err) {
    console.error(err);
    showToast("Archive failed", "danger");
  } finally {
    isArchiving.value = false;
  }
};

const migratePaths = async () => {
    if (!confirm("Are you sure you want to migrate paths? This operation modifies the database.")) {
        return;
    }
    migrating.value = true;
    try {
        await axios.post("/api/tools/migrate-paths");
        showToast("Path migration complete", "success");
    } catch (err) {
        console.error(err);
        showToast("Path migration failed", "danger");
    } finally {
        migrating.value = false;
    }
};

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
