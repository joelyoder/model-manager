<template>
  <div class="mx-4">
    <h2 class="mb-3">Settings</h2>
    <div class="mb-3">
      <label class="form-label">CivitAI API Key</label>
      <input v-model="apiKey" class="form-control" />
    </div>
    <div class="mb-3">
      <label class="form-label">Images Folder</label>
      <input v-model="imagesPath" class="form-control" />
    </div>
    <div class="mb-3">
      <label class="form-label">Models Folder</label>
      <input v-model="modelsPath" class="form-control" />
    </div>
    <button @click="save" class="btn btn-primary">Save</button>
    <router-link to="/" class="btn btn-secondary ms-2">Back</router-link>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import axios from "axios";
import { showToast } from "../utils/ui";

const apiKey = ref("");
const imagesPath = ref("");
const modelsPath = ref("");

onMounted(async () => {
  const res = await axios.get("/api/settings");
  apiKey.value = res.data.api_key || "";
  imagesPath.value = res.data.images_path || "";
  modelsPath.value = res.data.models_path || "";
});

const save = async () => {
  await axios.put("/api/settings", {
    api_key: apiKey.value,
    images_path: imagesPath.value,
    models_path: modelsPath.value,
  });
  showToast("Settings saved", "success");
};
</script>
