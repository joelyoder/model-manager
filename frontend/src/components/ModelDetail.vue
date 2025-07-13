<template>
  <div class="p-4 max-w-screen-lg mx-auto bg-gray-50 text-gray-900 dark:bg-gray-900 dark:text-gray-100">
    <button @click="goBack" class="mb-4 px-2 py-1 bg-gray-200 rounded dark:bg-gray-700 dark:text-white">⬅ Back</button>
    <h2 class="text-2xl font-bold mb-1">{{ model.name }}</h2>
    <h3 v-if="version.name" class="text-lg mb-2">{{ version.name }}</h3>
    <img
      v-if="imageUrl"
      :src="imageUrl"
      :width="model.imageWidth"
      :height="model.imageHeight"
      class="w-full max-w-md h-auto object-cover rounded mb-4 mx-auto"
    />
    <div v-if="model.description" v-html="model.description" class="mb-4"></div>
    <h3 class="text-lg font-semibold">Meta</h3>
    <table class="mt-4 border-collapse w-full text-sm">
      <tbody>
        <tr v-if="model.tags" class="border-b border-gray-700">
          <th class="pr-4 py-1">Tags</th>
          <td class="py-1">{{ model.tags.split(",").join(", ") }}</td>
        </tr>
        <tr class="border-b border-gray-700">
          <th class="pr-4 py-1">Type</th>
          <td class="py-1">{{ model.type }}</td>
        </tr>
        <tr class="border-b border-gray-700">
          <th class="pr-4 py-1">NSFW</th>
          <td class="py-1">{{ model.nsfw }}</td>
        </tr>
        <tr class="border-b border-gray-700">
          <th class="pr-4 py-1">Created</th>
          <td class="py-1">{{ model.createdAt }}</td>
        </tr>
        <tr class="border-b border-gray-700">
          <th class="pr-4 py-1">Updated</th>
          <td class="py-1">{{ model.updatedAt }}</td>
        </tr>
        <tr class="border-b border-gray-700">
          <th class="pr-4 py-1">Base Model</th>
          <td class="py-1">{{ version.baseModel }}</td>
        </tr>
        <tr v-if="version.trainedWords" class="border-b border-gray-700">
          <th class="pr-4 py-1">Trained Words</th>
          <td class="py-1">{{ version.trainedWords.split(",").join(", ") }}</td>
        </tr>
        <tr v-if="version.filePath" class="border-b border-gray-700">
          <th class="pr-4 py-1">File</th>
          <td class="py-1">{{ fileName }}</td>
        </tr>
        <tr v-if="version.sizeKB" class="border-b border-gray-700">
          <th class="pr-4 py-1">Size</th>
          <td class="py-1">{{ (version.sizeKB / 1024).toFixed(2) }} MB</td>
        </tr>
        <tr v-if="version.modelUrl" class="border-b border-gray-700">
          <th class="pr-4 py-1">Model URL</th>
          <td class="py-1">
            <a :href="version.modelUrl" target="_blank">View on CivitAI</a>
          </td>
        </tr>
      </tbody>
    </table>
    <button @click="deleteVersion" class="mt-4 px-2 py-1 bg-red-500 text-white rounded dark:bg-red-600">🗑 Delete Version</button>
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

