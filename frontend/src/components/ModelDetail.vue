<template>
  <div class="detail">
    <button @click="goBack">⬅ Back</button>
    <h2>{{ model.name }}</h2>
    <h3 v-if="version.name">{{ version.name }}</h3>
    <img
      v-if="imageUrl"
      :src="imageUrl"
      :width="model.imageWidth"
      :height="model.imageHeight"
    />
    <div v-if="model.description" v-html="model.description"></div>
    <h3>Meta</h3>
    <table class="meta">
      <tbody>
        <tr v-if="model.tags">
          <th>Tags</th>
          <td>{{ model.tags.split(",").join(", ") }}</td>
        </tr>
        <tr>
          <th>Type</th>
          <td>{{ model.type }}</td>
        </tr>
        <tr>
          <th>NSFW</th>
          <td>{{ model.nsfw }}</td>
        </tr>
        <tr>
          <th>Created</th>
          <td>{{ model.createdAt }}</td>
        </tr>
        <tr>
          <th>Updated</th>
          <td>{{ model.updatedAt }}</td>
        </tr>
        <tr>
          <th>Base Model</th>
          <td>{{ version.baseModel }}</td>
        </tr>
        <tr v-if="version.trainedWords">
          <th>Trained Words</th>
          <td>{{ version.trainedWords.split(",").join(", ") }}</td>
        </tr>
        <tr v-if="version.filePath">
          <th>File</th>
          <td>{{ fileName }}</td>
        </tr>
        <tr v-if="version.sizeKB">
          <th>Size</th>
          <td>{{ (version.sizeKB / 1024).toFixed(2) }} MB</td>
        </tr>
        <tr v-if="version.modelUrl">
          <th>Model URL</th>
          <td>
            <a :href="version.modelUrl" target="_blank">View on CivitAI</a>
          </td>
        </tr>
      </tbody>
    </table>
    <button @click="deleteVersion">🗑 Delete Version</button>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";

const router = useRouter();
const route = useRoute();
const model = ref({});
const version = ref({});

const imageUrl = computed(() => {
  const path = version.value.imagePath || model.value.imagePath;
  if (!path) return null;
  return path.replace(/^.*\/backend\/images/, "/images");
});

const fileName = computed(() => {
  if (!version.value.filePath) return "";
  return version.value.filePath.split("/").pop();
});

const fetchData = async () => {
  const { versionId } = route.params;
  const res = await axios.get(`/api/versions/${versionId}`);
  model.value = res.data.model;
  version.value = res.data.version;
};

onMounted(fetchData);

const deleteVersion = async () => {
  if (!confirm("Delete this version and all files?")) return;
  await axios.delete(`/api/versions/${route.params.versionId}`);
  router.push("/");
};

const goBack = () => {
  router.push("/");
};
</script>

<style scoped>
.detail {
  padding: 1rem;
}
img {
  max-width: 100%;
  height: auto;
  margin-bottom: 1rem;
}
.meta {
  margin-top: 1rem;
  border-collapse: collapse;
}
.meta th {
  text-align: left;
  padding-right: 0.5rem;
}
.meta td {
  padding-bottom: 0.25rem;
}
</style>
