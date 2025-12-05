<template>
  <div>
    <h5 class="mt-2 mb-3">Model Details</h5>
    <div class="mb-3">
      <label class="form-label">Name</label>
      <input
        :value="model.name"
        @input="updateModel('name', $event.target.value)"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Civit ID</label>
      <input
        :value="model.civitId"
        @input="updateModel('civitId', Number($event.target.value))"
        type="number"
        class="form-control"
      />
    </div>
    <hr />
    <h5 class="mt-2 mb-3">Version Details</h5>
    <div class="mb-3">
      <label class="form-label">Version Name</label>
      <input
        :value="version.name"
        @input="updateVersion('name', $event.target.value)"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Version ID</label>
      <input
        :value="version.versionId"
        @input="updateVersion('versionId', Number($event.target.value))"
        type="number"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Version Tags</label>
      <input
        :value="version.tags"
        @input="updateVersion('tags', $event.target.value)"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Version Type</label>
      <select
        :value="version.type"
        @change="updateVersion('type', $event.target.value)"
        class="form-select"
      >
        <option v-for="t in modelTypes" :key="t" :value="t">{{ t }}</option>
      </select>
    </div>
    <div class="form-check mb-3">
      <input
        type="checkbox"
        class="form-check-input"
        id="nsfw"
        :checked="version.nsfw"
        @change="updateVersion('nsfw', $event.target.checked)"
      />
      <label class="form-check-label" for="nsfw">Version NSFW</label>
    </div>
    <div class="mb-3">
      <label class="form-label">Version Description</label>
      <div ref="editor" style="height: 200px"></div>
    </div>
    <div class="mb-3">
      <label class="form-label">Base Model</label>
      <input
        :value="version.baseModel"
        @input="updateVersion('baseModel', $event.target.value)"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Trained Words</label>
      <input
        :value="version.trainedWords"
        @input="updateVersion('trainedWords', $event.target.value)"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Weight</label>
      <input
        :value="model.weight"
        @input="updateModel('weight', Number($event.target.value))"
        type="number"
        min="0"
        step="0.05"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Upload Image File</label>
      <div class="input-group">
        <input type="file" @change="onImageFileChange" class="form-control" />
        <button @click="uploadImage" class="btn btn-secondary">Upload</button>
      </div>
    </div>
    <div class="mb-3">
      <label class="form-label">Version Image Path</label>
      <input
        :value="version.imagePath"
        @input="updateVersion('imagePath', $event.target.value)"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Upload Model File</label>
      <div class="input-group">
        <input type="file" @change="onModelFileChange" class="form-control" />
        <button @click="uploadModel" class="btn btn-secondary">Upload</button>
      </div>
    </div>
    <div class="mb-3">
      <label class="form-label">Version File Path</label>
      <input
        :value="version.filePath"
        @input="updateVersion('filePath', $event.target.value)"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Size (KB)</label>
      <input
        :value="version.sizeKB"
        @input="updateVersion('sizeKB', Number($event.target.value))"
        type="number"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Model URL</label>
      <input
        :value="version.modelUrl"
        @input="updateVersion('modelUrl', $event.target.value)"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Early Access Time Frame</label>
      <input
        :value="version.earlyAccessTimeFrame"
        @input="updateVersion('earlyAccessTimeFrame', Number($event.target.value))"
        type="number"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Mode</label>
      <input
        :value="version.mode"
        @input="updateVersion('mode', $event.target.value)"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Created At</label>
      <input
        :value="version.createdAt"
        @input="updateVersion('createdAt', $event.target.value)"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Updated At</label>
      <input
        :value="version.updatedAt"
        @input="updateVersion('updatedAt', $event.target.value)"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">SHA256</label>
      <input
        :value="version.sha256"
        @input="updateVersion('sha256', $event.target.value)"
        class="form-control"
      />
    </div>
    <div class="mb-3">
      <label class="form-label">Download URL</label>
      <input
        :value="version.downloadUrl"
        @input="updateVersion('downloadUrl', $event.target.value)"
        class="form-control"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import Quill from "quill";
import axios from "axios";
import { showToast } from "../utils/ui";

const props = defineProps({
  model: Object,
  version: Object,
  modelTypes: Array,
});

const emit = defineEmits(["update:model", "update:version"]);

const editor = ref(null);
let quill;
const imageFile = ref(null);
const modelFile = ref(null);

const updateModel = (key, value) => {
  emit("update:model", { ...props.model, [key]: value });
};

const updateVersion = (key, value) => {
  emit("update:version", { ...props.version, [key]: value });
};

onMounted(() => {
  quill = new Quill(editor.value, { theme: "snow" });
  quill.clipboard.dangerouslyPasteHTML(props.version.description || "");
  quill.on("text-change", () => {
    updateVersion("description", quill.root.innerHTML);
  });
});

watch(
  () => props.version.description,
  (val) => {
    if (quill && val !== quill.root.innerHTML) {
      quill.clipboard.dangerouslyPasteHTML(val || "");
    }
  }
);

const onImageFileChange = (e) => {
  imageFile.value = e.target.files[0] || null;
};

const uploadImage = async () => {
  if (!imageFile.value) return;
  const fd = new FormData();
  fd.append("file", imageFile.value);
  try {
    const res = await axios.post(
      `/api/versions/${props.version.ID}/image`,
      fd,
      { headers: { "Content-Type": "multipart/form-data" } }
    );
    updateVersion("imagePath", res.data.path);
    imageFile.value = null;
    showToast("Image uploaded", "success");
  } catch (err) {
    console.error(err);
    showToast("Failed to upload image", "danger");
  }
};

const onModelFileChange = (e) => {
  modelFile.value = e.target.files[0] || null;
};

const uploadModel = async () => {
  if (!modelFile.value) return;
  const fd = new FormData();
  fd.append("file", modelFile.value);
  try {
    const res = await axios.post(
      `/api/versions/${props.version.ID}/file`,
      fd,
      { headers: { "Content-Type": "multipart/form-data" } }
    );
    updateVersion("filePath", res.data.path);
    updateVersion("sizeKB", res.data.sizeKB);
    updateVersion("sha256", res.data.sha256);
    modelFile.value = null;
    showToast("Model file uploaded", "success");
  } catch (err) {
    console.error(err);
    showToast("Failed to upload model file", "danger");
  }
};
</script>
