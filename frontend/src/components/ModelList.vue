<template>
  <div class="controls">
    <input v-model="search" placeholder="Search models..." />

    <button @click="fetchModels">🔄 Refresh</button>

    <!-- Paste URL and fetch versions -->
    <input v-model="modelUrl" placeholder="Paste CivitAI model URL" />
    <button @click="loadVersions" :disabled="loading || !modelUrl">
      🔍 Load Versions
    </button>

    <!-- Version selector -->
    <select v-if="versions.length" v-model="selectedVersionId">
      <option disabled value="">Select version</option>
      <option v-for="v in versions" :value="v.id" :key="v.id">
        {{ v.name }} | {{ v.baseModel }} |
        {{ ((v.sizeKB || 0) / 1024).toFixed(2) }} MB
      </option>
    </select>

    <!-- Download version -->
    <button
      v-if="selectedVersionId"
      @click="downloadSelectedVersion"
      :disabled="loading"
    >
      <span v-if="loading">⏳ Downloading...</span>
      <span v-else>📥 Download Selected</span>
    </button>
  </div>

  <div v-if="models.length === 0">No models found.</div>

  <div class="model-grid">
    <div v-for="card in versionCards" :key="card.version.id" class="card">
      <h3>{{ card.model.name }} - {{ card.version.name }}</h3>
      <img
        v-if="card.imageUrl"
        :src="card.imageUrl"
        :width="card.model.imageWidth"
        :height="card.model.imageHeight"
      />
      <p v-if="card.model.tags">
        Tags: {{ card.model.tags.split(",").join(", ") }}
      </p>
      <p v-if="card.version.filePath">
        File: {{ card.version.filePath.split("/").pop() }}
      </p>
      <p>Type: {{ card.model.type }}</p>
      <p>Base Model: {{ card.version.baseModel }}</p>
      <p
        v-if="
          card.version.trainedWordsArr && card.version.trainedWordsArr.length
        "
      >
        Trained Words: {{ card.version.trainedWordsArr.join(", ") }}
      </p>
      <p v-if="card.version.sizeKB">
        Size: {{ (card.version.sizeKB / 1024).toFixed(2) }} MB
      </p>
      <button @click="deleteModel(card.model.ID)">🗑 Delete</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import axios from "axios";

const models = ref([]);
const search = ref("");
const modelUrl = ref("");
const versions = ref([]);
const selectedVersionId = ref("");
const loading = ref(false);

const fetchModels = async () => {
  const res = await axios.get("/api/models");
  models.value = res.data.map((model) => {
    const imageUrl = model.imagePath
      ? model.imagePath.replace(/^.*\/backend\/images/, "/images")
      : null;
    const versions = (model.versions || []).map((v) => {
      const vImage = v.imagePath
        ? v.imagePath.replace(/^.*\/backend\/images/, "/images")
        : null;
      return { ...v, imageUrl: vImage };
    });
    return {
      ...model,
      versions,
      imageUrl,
    };
  });
};

onMounted(fetchModels);

const filteredModels = computed(() => {
  if (!search.value) return models.value;
  return models.value.filter((m) =>
    m.name.toLowerCase().includes(search.value.toLowerCase()),
  );
});

const versionCards = computed(() => {
  return filteredModels.value.flatMap((model) =>
    model.versions.map((v) => {
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
  );
});

const extractModelId = (url) => {
  const match = url.match(/models\/(\d+)/);
  return match ? match[1] : null;
};

const loadVersions = async () => {
  const id = extractModelId(modelUrl.value);
  if (!id) {
    alert("Invalid CivitAI model URL");
    return;
  }

  loading.value = true;
  try {
    const res = await axios.get(`/api/model/${id}/versions`);
    versions.value = res.data;
    selectedVersionId.value = "";
  } catch (err) {
    console.error(err);
    alert("Failed to load versions");
  } finally {
    loading.value = false;
  }
};

const downloadSelectedVersion = async () => {
  if (!selectedVersionId.value) return;

  loading.value = true;
  try {
    await axios.post(`/api/sync/version/${selectedVersionId.value}`);
    await fetchModels();
    alert("Version downloaded successfully");
  } catch (err) {
    console.error(err);
    alert("Download failed");
  } finally {
    modelUrl.value = "";
    versions.value = [];
    selectedVersionId.value = "";
    loading.value = false;
  }
};

const deleteModel = async (id) => {
  if (!confirm('Delete this model and all files?')) return;
  await axios.delete(`/api/models/${id}`);
  await fetchModels();
};
</script>

<style scoped>
.controls {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  padding: 1rem;
  align-items: center;
}
input,
select {
  padding: 0.5rem;
  flex: 1;
  min-width: 200px;
}
button {
  padding: 0.5rem 1rem;
}
.model-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1rem;
  padding: 1rem;
}
.card {
  background: #fff;
  padding: 1rem;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}
img {
  max-width: 100%;
  height: auto;
  object-fit: cover;
  border-radius: 4px;
  margin-bottom: 0.5rem;
}
</style>
