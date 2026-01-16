<template>
  <div>
    <!-- Manual Add Button -->
    <div class="mb-4">
      <button @click="$emit('createManual')" class="btn btn-primary w-100 py-2">
        <Icon icon="mdi:plus" class="me-2" />
        Add New Model Manually
      </button>
    </div>

    <div class="text-light mb-2 small fw-bold">Import from Civitai</div>

    <!-- URL Input -->
    <div class="input-group mb-3">
      <input
        v-model="modelUrl"
        placeholder="Paste Civitai model URL"
        class="form-control"
        @keyup.enter="loadVersions"
      />
      <button
        @click="loadVersions"
        :disabled="loading || !modelUrl"
        class="btn btn-primary"
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
          class="form-select"
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
              class="btn btn-secondary"
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
              class="btn btn-primary"
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
        <div class="progress flex-grow-1">
          <div
            class="progress-bar progress-bar-striped"
            :style="{ width: downloadProgress + '%' }"
          >
            {{ downloadProgress }}%
          </div>
        </div>
        <button
          class="btn btn-outline-danger btn-sm"
          type="button"
          @click="cancelDownload"
          :disabled="canceling"
        >
          <span
            v-if="canceling"
            class="spinner-border spinner-border-sm"
            aria-hidden="true"
          ></span>
          <span v-if="canceling" role="status" class="ps-2"
            >Cancelling...</span
          >
          <span v-else>Cancel</span>
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
