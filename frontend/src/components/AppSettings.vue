<template>
  <div class="container px-2 px-md-4 max-w-xl mx-auto">
    <div class="mb-3 d-flex gap-2 align-items-center px-2 px-md-0">
      <button 
        @click="goBack" 
        class="btn btn-outline-secondary btn-sm d-flex align-items-center justify-content-center border-0"
        aria-label="Back"
        title="Back"
        style="width: 40px; height: 40px;"
      >
        <Icon icon="mdi:arrow-left" width="24" height="24" />
      </button>
      <h2 class="h5 mb-0 fw-bold ms-2">Settings</h2>
      <button 
        class="btn btn-outline-primary btn-sm d-flex align-items-center justify-content-center border-0 ms-auto" 
        @click="save"
        aria-label="Save"
        title="Save"
        style="width: 40px; height: 40px;"
      >
        <Icon icon="mdi:content-save" width="24" height="24" />
      </button>
    </div>

    <div class="card border-0 shadow-sm bg-dark-subtle rounded-3 p-4">
        <div class="mb-3">
            <label class="form-label text-secondary fw-bold small text-uppercase">Civitai API Key</label>
            <input v-model="apiKey" type="text" class="form-control bg-dark border-0 text-white shadow-none" placeholder="Enter API Key" />
        </div>
        
        <div class="mb-3">
            <label class="form-label text-secondary fw-bold small text-uppercase">Model Path</label>
            <input v-model="modelPath" type="text" class="form-control bg-dark border-0 text-white shadow-none" placeholder="./backend/downloads" />
            <small class="text-secondary opacity-75 d-block mt-1">Local directory where models are stored.</small>
        </div>

        <div class="mb-3">
            <label class="form-label text-secondary fw-bold small text-uppercase">Image Path</label>
            <input v-model="imagePath" type="text" class="form-control bg-dark border-0 text-white shadow-none" placeholder="./backend/images" />
            <small class="text-secondary opacity-75 d-block mt-1">Local directory where images are stored.</small>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { Icon } from "@iconify/vue";
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
