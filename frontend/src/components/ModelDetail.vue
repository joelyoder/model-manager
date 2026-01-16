<template>
  <div class="px-2 px-md-4 container max-w-xl mx-auto">
    <!-- Header / Actions -->
    <div class="mb-3 d-flex gap-2 align-items-center px-2 px-md-0" v-if="!isEditing">
      <button 
        @click="goBack" 
        class="btn btn-outline-secondary btn-sm d-flex align-items-center justify-content-center border-0"
        aria-label="Back"
        title="Back"
        style="width: 40px; height: 40px;"
      >
        <Icon icon="mdi:arrow-left" width="24" height="24" />
      </button>
      
      <div class="ms-auto d-flex gap-2">
        <button 
            @click="isEditing = true" 
            class="btn btn-outline-primary btn-sm d-flex align-items-center justify-content-center border-0"
            aria-label="Edit"
            title="Edit"
            style="width: 40px; height: 40px;"
        >
            <Icon icon="mdi:pencil" width="24" height="24" />
        </button>
        
        <button
            type="button"
            @click="toggleNsfw"
            :disabled="togglingNsfw"
            class="btn btn-sm d-flex align-items-center justify-content-center border-0"
            :class="version.nsfw ? 'btn-danger' : 'btn-outline-secondary'"
            style="width: 40px; height: 40px;"
            aria-label="Toggle NSFW"
            title="Toggle NSFW"
        >
            <Icon :icon="version.nsfw ? 'mdi:eye-off' : 'mdi:eye'" width="24" height="24" />
        </button>
        
        <button 
            @click="handleDelete" 
            class="btn btn-outline-danger btn-sm d-flex align-items-center justify-content-center border-0"
            aria-label="Delete"
            title="Delete"
            style="width: 40px; height: 40px;"
        >
            <Icon icon="mdi:delete" width="24" height="24" />
        </button>
      </div>
    </div>
    
    <div v-else class="mb-3 d-flex gap-2 align-items-center justify-content-end">
      <button 
        @click="cancelEdit" 
        class="btn btn-outline-secondary btn-sm d-flex align-items-center justify-content-center border-0"
        aria-label="Cancel"
        title="Cancel"
        style="width: 40px; height: 40px;"
      >
        <Icon icon="mdi:close" width="24" height="24" />
      </button>
      <div class="ms-auto d-flex gap-2">
        <button 
        @click="saveEdit" 
        class="btn btn-outline-primary btn-sm d-flex align-items-center justify-content-center border-0"
        aria-label="Save"
        title="Save"
        style="width: 40px; height: 40px;"
        >
        <Icon icon="mdi:content-save" width="24" height="24" />
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div v-if="!isEditing">
      <div class="card border-0 shadow-sm bg-dark-subtle rounded-3 overflow-hidden">
        <div class="card-body p-4">
            <MetadataDisplay :model="model" :version="version" />

            <div class="my-4 border-top border-secondary opacity-25"></div>

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
        </div>
        
         <div class="card-footer border-0 bg-dark-subtle p-3 d-flex justify-content-center gap-2">
            <button @click="updateMeta" class="btn btn-outline-secondary btn-sm">
            Update Metadata
            </button>
            <button @click="updateDesc" class="btn btn-outline-secondary btn-sm">
            Update Description
            </button>
            <button @click="updateImages" class="btn btn-outline-secondary btn-sm">
            Refresh Images
            </button>
            <button @click="updateAll" class="btn btn-outline-secondary btn-sm">
            Update All
            </button>
        </div>
      </div>
    </div>

    <div v-else>
      <div class="card border-0 shadow-sm bg-dark-subtle rounded-3 p-4">
        <MetadataEditor
            v-model:model="model"
            v-model:version="version"
            :modelTypes="modelTypes"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from "vue";
import { Icon } from "@iconify/vue";
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
