<template>
  <div class="p-4 dark:text-white">
    <div class="flex">
      <button @click="goBack" class="mb-2 px-2 py-1 bg-gray-200 dark:text-gray-400 rounded">⬅ Back</button>
      <button @click="deleteVersion" class="mt-4 px-2 py-1 bg-red-500 text-white rounded">🗑 Delete</button>
    </div>
    
    <h2 class="text-xl font-bold">{{ model.name }}</h2>
    <h3 v-if="version.name" class="text-lg mb-2">{{ version.name }}</h3>
    <img
      v-if="imageUrl"
      :src="imageUrl"
      :width="model.imageWidth"
      :height="model.imageHeight"
      class="max-w-full h-auto mb-4"
    />
    <div v-if="model.description" v-html="model.description" class="mb-4"></div>
    <h3 class="text-lg font-semibold">Meta</h3>
    <table class="table-auto">
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

