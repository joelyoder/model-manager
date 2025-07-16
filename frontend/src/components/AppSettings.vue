<template>
  <div class="container px-4">
    <h2 class="mb-3">Settings</h2>
    <div class="mb-3 row">
      <label class="col-sm-3 col-form-label">Civitai API Key</label>
      <div class="col-sm-9">
        <input v-model="apiKey" type="text" class="form-control" />
      </div>
    </div>
    <button class="btn btn-primary" @click="save">Save</button>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import axios from "axios";
import { showToast } from "../utils/ui";

const apiKey = ref("");

onMounted(async () => {
  const res = await axios.get("/api/settings");
  const item = res.data.find((s) => s.key === "civitai_api_key");
  if (item) apiKey.value = item.value;
});

async function save() {
  await axios.post("/api/settings", {
    key: "civitai_api_key",
    value: apiKey.value,
  });
  showToast("Settings saved");
}
</script>
