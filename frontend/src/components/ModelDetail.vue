<template>
  <div class="px-4">
    <div
      class="mb-2 d-flex gap-2 justify-content-center px-4 pb-4"
      v-if="!isEditing"
    >
      <button @click="goBack" class="btn btn-secondary">⬅ Back</button>
      <button @click="startEdit" class="btn btn-primary">✏️ Edit</button>
      <button @click="deleteVersion" class="btn btn-danger">🗑 Delete</button>
    </div>
    <div v-else class="d-flex gap-2 justify-content-center">
      <button @click="saveEdit" class="btn btn-primary">💾 Save</button>
      <button @click="cancelEdit" class="btn btn-secondary">Cancel</button>
    </div>
    <div v-if="!isEditing" class="container">
      <div class="row">
        <div class="col-4">
          <img
            v-if="imageUrl"
            :src="imageUrl"
            :width="model.imageWidth"
            :height="model.imageHeight"
            class="img-fluid mb-4"
          />
        </div>
        <div class="col-8">
          <h2 class="fw-bold">{{ model.name }}</h2>
          <h3 v-if="version.name" class="mb-2">{{ version.name }}</h3>
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
                  <a :href="version.modelUrl" target="_blank">{{
                    version.modelUrl
                  }}</a>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <div
        v-if="model.description"
        v-html="model.description"
        class="mb-4"
      ></div>
    </div>
    <div v-else>
      <div class="mb-3">
        <label class="form-label">Name</label>
        <input v-model="model.name" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Version Name</label>
        <input v-model="version.name" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Tags</label>
        <input v-model="model.tags" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Type</label>
        <input v-model="model.type" class="form-control" />
      </div>
      <div class="form-check mb-3">
        <input
          type="checkbox"
          class="form-check-input"
          id="nsfw"
          v-model="model.nsfw"
        />
        <label class="form-check-label" for="nsfw">NSFW</label>
      </div>
      <div class="mb-3">
        <label class="form-label">Description</label>
        <div ref="editor" style="height: 200px"></div>
      </div>
      <div class="mb-3">
        <label class="form-label">Created</label>
        <input v-model="model.createdAt" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Updated</label>
        <input v-model="model.updatedAt" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Base Model</label>
        <input v-model="version.baseModel" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Trained Words</label>
        <input v-model="version.trainedWords" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">File Path</label>
        <input v-model="version.filePath" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Size (KB)</label>
        <input
          v-model.number="version.sizeKB"
          type="number"
          class="form-control"
        />
      </div>
      <div class="mb-3">
        <label class="form-label">Model URL</label>
        <input v-model="version.modelUrl" class="form-control" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, nextTick, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";
import Quill from "quill";
import { showToast, showConfirm } from "../utils/ui";

const router = useRouter();
const route = useRoute();
const model = ref({});
const version = ref({});
const isEditing = ref(false);
const editor = ref(null);
let quill;

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

watch(isEditing, async (val) => {
  if (val) {
    await nextTick();
    quill = new Quill(editor.value, { theme: "snow" });
    quill.root.innerHTML = model.value.description || "";
  }
});

const deleteVersion = async () => {
  if (!(await showConfirm("Delete this version and all files?"))) return;
  await axios.delete(`/api/versions/${route.params.versionId}`);
  router.push("/");
};

const startEdit = () => {
  isEditing.value = true;
};

const cancelEdit = async () => {
  isEditing.value = false;
  await fetchData();
};

const saveEdit = async () => {
  if (quill) {
    model.value.description = quill.root.innerHTML;
  }
  await axios.put(`/api/models/${model.value.ID}`, model.value);
  await axios.put(`/api/versions/${version.value.ID}`, version.value);
  isEditing.value = false;
  await fetchData();
  showToast("Saved", "success");
};

const goBack = () => {
  router.push("/");
};
</script>
