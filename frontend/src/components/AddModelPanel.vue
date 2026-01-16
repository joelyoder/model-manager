<template>
  <div>
    <!-- Manual Add Button -->
    <div class="mb-4">
      <button @click="$emit('createManual')" class="btn btn-outline-primary w-auto px-4 py-2 border-2 fw-bold">
        <Icon icon="mdi:plus" class="me-2" />
        Create New Model Manually
      </button>
    </div>

    <div class="text-secondary mb-2 small fw-bold text-uppercase tracking-wide">Import from Civitai</div>

    <!-- URL Input -->
    <div class="input-group mb-3">
      <input
        v-model="modelUrl"
        placeholder="Paste Civitai model URL"
        class="form-control bg-dark-subtle border-0 text-white shadow-none"
        @keyup.enter="loadVersions"
      />
      <button
        @click="loadVersions"
        :disabled="loading || !modelUrl"
        class="btn btn-primary border-0"
      >
        Load Versions
      </button>
    </div>

    <!-- Version Selection -->
    <div>
      <div class="input-group mb-2">
        <select
          v-if="versions.length"
          v-model="selectedVersionId"
          class="form-select bg-dark-subtle border-0 text-white shadow-none"
          style="min-width: 200px"
        >
          <option disabled value="">Select version</option>
          <option v-for="v in versions" :value="v.id" :key="v.id">
            {{ v.name }} | {{ v.baseModel }} |
            {{ ((v.sizeKB || 0) / 1024).toFixed(2) }} MB
          </option>
        </select>
        
        <!-- Actions -->
        <template v-if="versions.length">
             <button
              v-if="selectedVersionId"
              @click="addSelectedVersion"
              :disabled="loading"
              class="btn btn-secondary border-0"
            >
              <span
                v-if="loading && adding"
                class="spinner-border spinner-border-sm"
                aria-hidden="true"
              ></span>
              <span v-if="loading && adding" role="status" class="ps-2"
                >Adding...</span
              >
              <span v-else>Add</span>
            </button>

            <button
              v-if="selectedVersionId"
              @click="downloadSelectedVersion"
              :disabled="loading"
              class="btn btn-primary border-0"
            >
              <span
                v-if="loading && !adding"
                class="spinner-border spinner-border-sm"
                aria-hidden="true"
              ></span>
              <span v-if="loading && !adding" role="status" class="ps-2"
                >Downloading...</span
              >
              <span v-else>Download</span>
            </button>
        </template>
      </div>

       <!-- Downloading Progress -->
      <div
        v-if="downloading"
        class="d-flex align-items-center gap-2 w-100 mb-2 mt-3"
      >
        <div class="progress flex-grow-1 bg-dark-subtle" style="height: 10px;">
          <div
            class="progress-bar progress-bar-striped progress-bar-animated bg-primary"
            :style="{ width: downloadProgress + '%' }"
          >
          </div>
        </div>
        <small class="text-white">{{ downloadProgress }}%</small>
        <button
          class="btn btn-outline-danger btn-sm border-0"
          type="button"
          @click="cancelDownload"
          :disabled="canceling"
        >
          <Icon icon="mdi:close" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import axios from "axios";
import { showToast } from "../utils/ui";
import { useDownloads } from "../composables/useDownloads";

const emit = defineEmits(["createManual", "added"]);

const modelUrl = ref("");
const versions = ref([]);
const selectedVersionId = ref("");
const loading = ref(false);
const adding = ref(false);

const {
  downloading,
  downloadProgress,
  canceling,
  startDownload,
  cancelDownload,
} = useDownloads();

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

const addSelectedVersion = async () => {
  if (!selectedVersionId.value) return;
  loading.value = true;
  adding.value = true;
  try {
    const buildSyncVersionUrl = (versionId, { download, modelId } = {}) => {
      const params = new URLSearchParams();
      if (modelId) params.set("modelId", String(modelId));
      if (download !== undefined) params.set("download", String(download));
      const query = params.toString();
      return `/api/sync/version/${versionId}${query ? `?${query}` : ""}`;
    };

    await axios.post(buildSyncVersionUrl(selectedVersionId.value, { download: 0 }));
    showToast("Version added", "success");
    emit("added");
  } catch (err) {
    console.error(err);
    showToast("Failed to add version", "danger");
  } finally {
    loading.value = false;
    adding.value = false;
  }
};

const downloadSelectedVersion = async () => {
  if (!selectedVersionId.value) return;
  loading.value = true;
  try {
    await startDownload(selectedVersionId.value);
    emit("added");
    modelUrl.value = "";
    versions.value = [];
    selectedVersionId.value = "";
  } finally {
    loading.value = false;
  }
};
</script>
