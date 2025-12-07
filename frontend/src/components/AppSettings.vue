<template>
  <div class="container px-4">
    <div class="mb-2 d-flex gap-2 pb-2">
      <button @click="goBack" class="btn btn-secondary">Back</button>
      <button class="btn btn-primary" @click="save">Save</button>
    </div>
    <h2 class="mb-3">Settings</h2>
    <div class="mb-3 row">
      <label class="col-sm-3 col-form-label">Civitai API Key</label>
      <div class="col-sm-9">
        <input v-model="apiKey" type="text" class="form-control" />
      </div>
    </div>
    <div class="mb-3 row">
      <label class="col-sm-3 col-form-label">Model Path</label>
      <div class="col-sm-9">
        <input v-model="modelPath" type="text" class="form-control" placeholder="./backend/downloads" />
        <small class="text-white-50">Local directory where models are stored.</small>
      </div>
    </div>
    <div class="mb-3 row">
      <label class="col-sm-3 col-form-label">Image Path</label>
      <div class="col-sm-9">
        <input v-model="imagePath" type="text" class="form-control" placeholder="./backend/images" />
        <small class="text-white-50">Local directory where images are stored.</small>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import axios from "axios";
import { showToast } from "../utils/ui";

const apiKey = ref("");
const modelPath = ref("");
const imagePath = ref("");
const router = useRouter();

onMounted(async () => {
  const res = await axios.get("/api/settings");
  const keyItem = res.data.find((s) => s.key === "civitai_api_key");
  if (keyItem) apiKey.value = keyItem.value;
  
  const modelItem = res.data.find((s) => s.key === "model_path");
  if (modelItem) modelPath.value = modelItem.value;

  const imageItem = res.data.find((s) => s.key === "image_path");
  if (imageItem) imagePath.value = imageItem.value;
});

async function save() {
  await axios.post("/api/settings", {
    key: "civitai_api_key",
    value: apiKey.value,
  });
  if (modelPath.value) {
    await axios.post("/api/settings", {
        key: "model_path",
        value: modelPath.value,
    });
  }
  if (imagePath.value) {
    await axios.post("/api/settings", {
        key: "image_path",
        value: imagePath.value,
    });
  }
  showToast("Settings saved");
}

function goBack() {
  router.push("/");
}
</script>
