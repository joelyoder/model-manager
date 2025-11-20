<template>
  <div class="px-4 container">
    <div class="mb-2 d-flex gap-2 pb-2" v-if="!isEditing">
      <button @click="goBack" class="btn btn-secondary">Back</button>
      <button @click="isEditing = true" class="btn btn-primary">Edit</button>
      <button
        type="button"
        @click="toggleNsfw"
        :disabled="togglingNsfw"
        class="btn btn-sm px-2"
        :class="version.nsfw ? 'btn-danger' : 'btn-secondary'"
        style="--bs-btn-padding-y: 0.25rem; --bs-btn-padding-x: 0.25rem"
        aria-label="Toggle NSFW"
        title="Toggle NSFW"
      >
        <span class="visually-hidden">Toggle NSFW</span>
        <svg
          v-if="version.nsfw"
          width="18px"
          height="18px"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          color="#ffffff"
        >
          <path
            d="M10.733 5.076a10.744 10.744 0 0 1 11.205 6.575 1 1 0 0 1 0 .696 10.747 10.747 0 0 1-1.444 2.49"
            stroke="#ffffff"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          ></path>
          <path
            d="M14.084 14.158a3 3 0 0 1-4.242-4.242"
            stroke="#ffffff"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          ></path>
          <path
            d="M17.479 17.499a10.75 10.75 0 0 1-15.417-5.151 1 1 0 0 1 0-.696 10.75 10.75 0 0 1 4.446-5.143"
            stroke="#ffffff"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          ></path>
          <path
            d="m2 2 20 20"
            stroke="#ffffff"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          ></path>
        </svg>
        <svg
          v-else
          width="18px"
          height="18px"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          color="#ffffff"
        >
          <path
            d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0"
            stroke="#ffffff"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          ></path>
          <circle
            cx="12"
            cy="12"
            r="3"
            stroke="#ffffff"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          ></circle>
        </svg>
      </button>
      <button @click="handleDelete" class="btn btn-outline-danger ms-auto">
        Delete
      </button>
    </div>
    <div v-else class="mb-2 d-flex gap-2 pb-2">
      <button @click="cancelEdit" class="btn btn-secondary">Cancel</button>
      <button @click="saveEdit" class="btn btn-primary">Save</button>
    </div>

    <div v-if="!isEditing">
      <MetadataDisplay :model="model" :version="version" />

      <ImageGallery
        :images="version.images"
        :currentImagePath="version.imagePath"
        :versionId="version.ID"
        :versionType="version.type"
        :versionMode="version.mode"
        @setMain="setMainImage"
        @remove="removeImage"
        @uploaded="fetchData(route.params.versionId)"
      />

      <div class="mb-2 d-flex justify-content-center gap-2 pb-2">
        <button @click="updateMeta" class="btn btn-secondary btn-sm">
          Update Metadata
        </button>
        <button @click="updateDesc" class="btn btn-secondary btn-sm">
          Update Description
        </button>
        <button @click="updateImages" class="btn btn-secondary btn-sm">
          Refresh Images
        </button>
        <button @click="updateAll" class="btn btn-secondary btn-sm">
          Update All
        </button>
      </div>
    </div>

    <div v-else>
      <MetadataEditor
        v-model:model="model"
        v-model:version="version"
        :modelTypes="modelTypes"
      />
    </div>
  </div>
</template>

<script setup>
import { onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";
import { showToast, showConfirm, showDeleteConfirm } from "../utils/ui";
import { useModelDetail } from "../composables/useModelDetail";
import MetadataDisplay from "./MetadataDisplay.vue";
import MetadataEditor from "./MetadataEditor.vue";
import ImageGallery from "./ImageGallery.vue";

const router = useRouter();
const route = useRoute();

const {
  model,
  version,
  isEditing,
  togglingNsfw,
  fetchData,
  toggleNsfw,
  deleteVersion,
  saveEdit,
  refreshVersion,
} = useModelDetail();

const modelTypes = [
  "Checkpoint",
  "TextualInversion",
  "Hypernetwork",
  "AestheticGradient",
  "LORA",
  "LoCon",
  "DoRA",
  "Controlnet",
  "Upscaler",
  "MotionModule",
  "VAE",
  "Wildcards",
  "Poses",
  "Workflows",
  "Detection",
  "Other",
];

onMounted(async () => {
  await fetchData(route.params.versionId);
  if (route.query.edit === "1") {
    isEditing.value = true;
  }
});

const cancelEdit = async () => {
  isEditing.value = false;
  await fetchData(route.params.versionId);
};

const handleDelete = async () => {
  const choice = await showDeleteConfirm("Delete this version?");
  if (!choice) return;
  const files = choice === "deleteFiles";
  await deleteVersion(route.params.versionId, files);
  router.push("/");
};

const updateMeta = async () => {
  if (!(await showConfirm("Pull latest metadata from CivitAI?"))) return;
  await refreshVersion("metadata");
  await fetchData(route.params.versionId);
  showToast("Updated", "success");
};
const updateDesc = async () => {
  if (!(await showConfirm("Pull latest description from CivitAI?"))) return;
  await refreshVersion("description");
  await fetchData(route.params.versionId);
  showToast("Updated", "success");
};
const updateImages = async () => {
  if (!(await showConfirm("Replace all images with the latest from CivitAI?")))
    return;
  await refreshVersion("images");
  await fetchData(route.params.versionId);
  showToast("Updated", "success");
};
const updateAll = async () => {
  if (
    !(await showConfirm(
      "Update all data from CivitAI? This will replace images."
    ))
  )
    return;
  await refreshVersion("all");
  await fetchData(route.params.versionId);
  showToast("Updated", "success");
};

const goBack = () => {
  router.push({ path: "/", query: { scrollTo: version.value.ID } });
};

const setMainImage = async (img) => {
  await axios.post(`/api/versions/${version.value.ID}/main-image/${img.ID}`);
  await fetchData(route.params.versionId);
  showToast("Main image updated", "success");
};

const removeImage = async (img) => {
  if (!(await showConfirm("Remove this image?"))) return;
  await axios.delete(`/api/versions/${version.value.ID}/images/${img.ID}`);
  if (version.value.imagePath === img.path) {
    await fetchData(route.params.versionId);
  } else {
    version.value.images = (version.value.images || []).filter(
      (i) => i.ID !== img.ID
    );
  }
  showToast("Image removed", "success");
};
</script>
