<template>
  <div class="p-4">
    <button @click="goBack" class="btn btn-secondary mb-2">⬅ Back</button>
    <h2 class="h4 fw-bold">{{ model.name }}</h2>
    <h3 v-if="version.name" class="fs-5 mb-2">{{ version.name }}</h3>
    <img
      v-if="imageUrl"
      :src="imageUrl"
      :width="model.imageWidth"
      :height="model.imageHeight"
      class="img-fluid mb-4"
    />
    <div v-if="model.description" v-html="model.description" class="mb-4"></div>
    <h3 class="fs-5 fw-semibold">Meta</h3>
    <table class="table mt-4">
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
    <button @click="deleteVersion" class="btn btn-danger mt-4">🗑 Delete Version</button>
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

