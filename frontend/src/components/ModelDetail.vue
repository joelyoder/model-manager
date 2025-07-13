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
    <p v-if="model.tags">Tags: {{ model.tags.split(",").join(", ") }}</p>
    <p>Type: {{ model.type }}</p>
    <p>NSFW: {{ model.nsfw }}</p>
    <p>Created: {{ model.createdAt }}</p>
    <p>Updated: {{ model.updatedAt }}</p>
    <p>Base Model: {{ version.baseModel }}</p>
    <p v-if="version.trainedWords">
      Trained Words: {{ version.trainedWords.split(",").join(", ") }}
    </p>
    <p v-if="version.filePath">File: {{ fileName }}</p>
    <p v-if="version.sizeKB">
      Size: {{ (version.sizeKB / 1024).toFixed(2) }} MB
    </p>
    <p v-if="version.modelUrl">
      <a :href="version.modelUrl" target="_blank">View on CivitAI</a>
    </p>
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
</style>
