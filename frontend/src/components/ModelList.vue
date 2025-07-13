<template>
  <div class="flex flex-wrap gap-2 p-4 items-center">
    <input
      v-model="search"
      placeholder="Search models..."
      class="p-2 flex-1 min-w-[200px] border rounded"
    />

    <button @click="fetchModels" class="px-4 py-2 bg-gray-200 rounded">
      🔄 Refresh
    </button>

    <!-- Paste URL and fetch versions -->
    <input
      v-model="modelUrl"
      placeholder="Paste CivitAI model URL"
      class="p-2 flex-1 min-w-[200px] border rounded"
    />
    <button
      @click="loadVersions"
      :disabled="loading || !modelUrl"
      class="px-4 py-2 bg-gray-200 rounded disabled:opacity-50"
    >
      🔍 Load Versions
    </button>

    <!-- Version selector -->
    <select
      v-if="versions.length"
      v-model="selectedVersionId"
      class="p-2 flex-1 min-w-[200px] border rounded"
    >
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
      class="px-4 py-2 bg-blue-500 text-white rounded disabled:opacity-50"
    >
      <span v-if="loading">⏳ Downloading...</span>
      <span v-else>📥 Download Selected</span>
    </button>
  </div>

  <div v-if="models.length === 0">No models found.</div>

  <div class="grid gap-4 p-4 grid-cols-[repeat(auto-fill,minmax(320px,_1fr))]">
    <div v-for="card in versionCards" :key="card.version.ID" class="bg-white p-4 rounded shadow">
      <h3>{{ card.model.name }} - {{ card.version.name }}</h3>
      <img
        v-if="card.imageUrl"
        :src="card.imageUrl"
        :width="card.model.imageWidth"
        :height="card.model.imageHeight"
        class="w-full h-auto object-cover rounded mb-2"
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
      <button @click="deleteVersion(card.version.ID)" class="px-2 py-1 bg-red-500 text-white rounded">
        🗑 Delete
      </button>
      <button
        v-if="card.version.filePath"
        @click="goToModel(card.model.ID, card.version.ID)"
        class="px-2 py-1 bg-blue-500 text-white rounded"
      >
        ℹ️ More details
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";

const models = ref([]);
const search = ref("");
const modelUrl = ref("");
const versions = ref([]);
const selectedVersionId = ref("");
const loading = ref(false);
const router = useRouter();

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

const deleteVersion = async (id) => {
  if (!confirm("Delete this version and all files?")) return;
  await axios.delete(`/api/versions/${id}`);
  await fetchModels();
};

const goToModel = (modelId, versionId) => {
  router.push(`/model/${modelId}/version/${versionId}`);
};
</script>

