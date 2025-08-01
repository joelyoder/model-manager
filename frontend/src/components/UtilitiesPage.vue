<template>
  <div class="container px-4">
    <div class="mb-2 d-flex gap-2 pb-2">
      <button @click="goBack" class="btn btn-secondary">Back</button>
    </div>
    <h2 class="my-3">Utilities</h2>
    <div class="card card-body mb-4" v-if="stats">
      <h3 class="h5">Stats</h3>
      <p>Total Models: {{ stats.totalModels }}</p>
      <ul class="mb-3">
        <li v-for="t in stats.typeCounts" :key="t.Key">
          {{ t.Key || "Unknown" }}: {{ t.Count }}
        </li>
      </ul>
      <div class="row">
        <div class="col-md-6 mb-3">
          <canvas id="baseModelChart"></canvas>
        </div>
        <div class="col-md-6 mb-3">
          <canvas id="nsfwChart"></canvas>
        </div>
      </div>
    </div>
    <h3 class="h5 mt-5">Import JSON from Model Organizer</h3>
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
    <h3 class="h5 mt-5">Export Database as JSON</h3>
    <div class="mb-3 d-flex gap-2">
      <button @click="exportJson" class="btn btn-primary">Export Models</button>
    </div>
    <h3 class="h5 mt-5">Import Database JSON</h3>
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
</template>

<script setup>
/* global Chart */
import { ref, onMounted, nextTick } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";
import { showToast } from "../utils/ui";

const stats = ref(null);
let baseChart = null;
let nsfwChart = null;

onMounted(async () => {
  try {
    const res = await axios.get("/api/stats");
    stats.value = res.data;
    await nextTick();
    renderCharts();
  } catch (err) {
    console.error(err);
  }
});

function renderCharts() {
  if (!stats.value) return;
  if (baseChart) baseChart.destroy();
  if (nsfwChart) nsfwChart.destroy();

  const baseCtx = document.getElementById("baseModelChart");
  if (baseCtx) {
    baseChart = new Chart(baseCtx, {
      type: "pie",
      data: {
        labels: stats.value.baseModelCounts.map((b) => b.Key || "Unknown"),
        datasets: [
          {
            data: stats.value.baseModelCounts.map((b) => b.Count),
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

const importFile = ref(null);
const dbImportFile = ref(null);
const pullImages = ref(false);
const pullMeta = ref(false);
const pullDesc = ref(false);
const router = useRouter();

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

const goBack = () => {
  router.push("/");
};
</script>
