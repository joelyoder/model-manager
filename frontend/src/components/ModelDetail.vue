<template>
  <div class="px-4 container">
    <div class="mb-2 d-flex gap-2 pb-2" v-if="!isEditing">
      <button @click="goBack" class="btn btn-secondary">Back</button>
      <button @click="startEdit" class="btn btn-primary">Edit</button>
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
      <button @click="deleteVersion" class="btn btn-outline-danger ms-auto">
        Delete
      </button>
    </div>
    <div v-else class="mb-2 d-flex gap-2 pb-2">
      <button @click="cancelEdit" class="btn btn-secondary">Cancel</button>
      <button @click="saveEdit" class="btn btn-primary">Save</button>
    </div>
    <div v-if="!isEditing">
      <div class="row">
        <div class="col-md-4">
          <img v-if="imageUrl" :src="imageUrl" class="img-fluid mb-4" />
        </div>
        <div class="col-md-8">
          <h2 class="fw-bold">{{ model.name }}</h2>
          <h3 v-if="version.name" class="mb-2">{{ version.name }}</h3>
          <div class="d-flex flex-wrap align-items-center gap-2 mb-3">
            <span
              v-if="version.type"
              class="badge rounded-pill text-bg-primary"
            >
              {{ version.type }}
            </span>
            <span
              v-if="version.baseModel"
              class="badge rounded-pill text-bg-success"
            >
              {{ version.baseModel }}
            </span>
          </div>
          <dl class="metadata-list my-4">
            <template v-if="version.tags">
              <dt class="metadata-list__label">Tags</dt>
              <dd class="metadata-list__value">
                {{ version.tags.split(",").join(", ") }}
              </dd>
            </template>
            <template v-if="version.trainedWords">
              <dt class="metadata-list__label">Trained Words</dt>
              <dd class="metadata-list__value">
                <div class="d-flex align-items-center gap-2 flex-wrap">
                  <span>{{ formattedTrainedWords }}</span>
                  <button
                    type="button"
                    class="btn btn-outline-secondary btn-sm d-flex align-items-center justify-content-center"
                    @click="copyTrainedWords"
                    aria-label="Copy trained words"
                  >
                    <Icon
                      icon="mdi:content-copy"
                      width="16"
                      height="16"
                      aria-hidden="true"
                    />
                  </button>
                </div>
              </dd>
            </template>
            <dt class="metadata-list__label">Weight</dt>
            <dd class="metadata-list__value">{{ weightDisplay }}</dd>
            <template v-if="version.filePath">
              <dt class="metadata-list__label">File</dt>
              <dd class="metadata-list__value">
                <div class="d-flex align-items-center gap-2 flex-wrap">
                  <span>{{ fileName }}</span>
                  <button
                    type="button"
                    class="btn btn-outline-secondary btn-sm d-flex align-items-center justify-content-center"
                    @click="copyFileBaseName"
                    aria-label="Copy filename without extension"
                  >
                    <Icon
                      icon="mdi:file-document-outline"
                      width="16"
                      height="16"
                      aria-hidden="true"
                    />
                  </button>
                  <button
                    type="button"
                    class="btn btn-outline-secondary btn-sm d-flex align-items-center justify-content-center"
                    @click="copyLoraTag"
                    aria-label="Copy LoRA tag"
                  >
                    <Icon
                      icon="mdi:tag-text-outline"
                      width="16"
                      height="16"
                      aria-hidden="true"
                    />
                  </button>
                </div>
              </dd>
            </template>
            <template v-if="version.modelUrl">
              <dt class="metadata-list__label">Model URL</dt>
              <dd class="metadata-list__value">
                <a :href="version.modelUrl" target="_blank">{{
                  version.modelUrl
                }}</a>
              </dd>
            </template>
          </dl>
        </div>
      </div>
      <div
        v-if="version.description"
        v-html="version.description"
        class="mb-4"
      ></div>
      <div
        v-if="hasVersionStats"
        class="row row-cols-1 row-cols-md-2 g-3 mb-4"
      >
        <div v-if="version.createdAt" class="col">
          <dl class="metadata-summary">
            <dt class="metadata-list__label">Created</dt>
            <dd class="metadata-list__value">{{ createdAtReadable }}</dd>
          </dl>
        </div>
        <div v-if="version.updatedAt" class="col">
          <dl class="metadata-summary">
            <dt class="metadata-list__label">Updated</dt>
            <dd class="metadata-list__value">{{ updatedAtReadable }}</dd>
          </dl>
        </div>
         <div v-if="version.sizeKB" class="col">
          <dl class="metadata-summary">
            <dt class="metadata-list__label">Size</dt>
            <dd class="metadata-list__value">{{ versionSizeMb }}</dd>
          </dl>
        </div>
        <div v-if="version.sha256" class="col">
          <dl class="metadata-summary">
            <dt class="metadata-list__label">SHA256</dt>
            <dd class="metadata-list__value"><code>{{ version.sha256 }}</code></dd>
          </dl>
        </div>
      </div>
      <div
        v-if="galleryImages.length"
        class="row row-cols-1 row-cols-md-2 row-cols-lg-3 g-4 mb-4"
      >
        <div v-for="img in galleryImages" :key="img.ID" class="col">
          <img
            :src="img.url"
            :width="img.width"
            :height="img.height"
            class="img-fluid"
          />
          <div class="row g-1">
            <div class="col">
              <button
                v-if="img.path !== version.imagePath"
                @click="setMainImage(img)"
                class="btn btn-outline-primary btn-sm mt-1 w-100"
              >
                Set as Main
              </button>
              <span
                v-else
                class="badge text-bg-success d-block text-center mt-1 py-2"
              >
                Main Image
              </span>
            </div>
            <div class="col">
              <button
                @click="removeImage(img)"
                class="btn btn-outline-danger btn-sm mt-1 w-100"
              >
                Remove
              </button>
            </div>
          </div>
          <div
            v-if="Object.keys(img.parsedMeta || {}).length"
            class="mt-1"
          >
            <button
              class="btn btn-outline-secondary btn-sm w-100"
              type="button"
              data-bs-toggle="collapse"
              :data-bs-target="'#meta-' + img.ID"
              aria-expanded="false"
              :aria-controls="'meta-' + img.ID"
            >
              Show Metadata
            </button>
            <div :id="'meta-' + img.ID" class="collapse mt-1">
              <dl class="image-metadata">
                <template v-for="(value, key) in img.parsedMeta" :key="key">
                  <dt class="image-metadata__label">{{ key }}</dt>
                  <dd class="image-metadata__value">{{ value }}</dd>
                </template>
              </dl>
            </div>
          </div>
        </div>
      </div>
      <div class="input-group mb-3">
        <input type="file" @change="onGalleryFileChange" class="form-control" />
        <button @click="uploadGallery" class="btn btn-secondary">
          Add Image
        </button>
      </div>
      <div
        class="mb-2 d-flex justify-content-center gap-2 pb-2"
        v-if="!isEditing"
      >
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
      <h5 class="mt-2 mb-3">Model Details</h5>
      <div class="mb-3">
        <label class="form-label">Name</label>
        <input v-model="model.name" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Civit ID</label>
        <input
          v-model.number="model.civitId"
          type="number"
          class="form-control"
        />
      </div>
      <hr />
      <h5 class="mt-2 mb-3">Version Details</h5>
      <div class="mb-3">
        <label class="form-label">Version Name</label>
        <input v-model="version.name" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Version ID</label>
        <input
          v-model.number="version.versionId"
          type="number"
          class="form-control"
        />
      </div>
      <div class="mb-3">
        <label class="form-label">Version Tags</label>
        <input v-model="version.tags" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Version Type</label>
        <select v-model="version.type" class="form-select">
          <option v-for="t in modelTypes" :key="t" :value="t">{{ t }}</option>
        </select>
      </div>
      <div class="form-check mb-3">
        <input
          type="checkbox"
          class="form-check-input"
          id="nsfw"
          v-model="version.nsfw"
        />
        <label class="form-check-label" for="nsfw">Version NSFW</label>
      </div>
      <div class="mb-3">
        <label class="form-label">Version Description</label>
        <div ref="editor" style="height: 200px"></div>
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
        <label class="form-label">Weight</label>
        <input
          v-model.number="model.weight"
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
        <input v-model="version.imagePath" class="form-control" />
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
      <div class="mb-3">
        <label class="form-label">Early Access Time Frame</label>
        <input
          v-model.number="version.earlyAccessTimeFrame"
          type="number"
          class="form-control"
        />
      </div>
      <div class="mb-3">
        <label class="form-label">Mode</label>
        <input v-model="version.mode" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Created At</label>
        <input v-model="version.createdAt" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Updated At</label>
        <input v-model="version.updatedAt" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">SHA256</label>
        <input v-model="version.sha256" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Download URL</label>
        <input v-model="version.downloadUrl" class="form-control" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, nextTick, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";
import Quill from "quill";
import { showToast, showConfirm, showDeleteConfirm } from "../utils/ui";
import { Icon } from "@iconify/vue";

const router = useRouter();
const route = useRoute();
const model = ref({});
const version = ref({});
const isEditing = ref(false);
const editor = ref(null);
let quill;
const togglingNsfw = ref(false);

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

const imageFile = ref(null);
const modelFile = ref(null);
const galleryFile = ref(null);

const normalizeWeight = (weight) => {
  const num = Number(weight);
  if (Number.isFinite(num) && num > 0) {
    return num;
  }
  return 1;
};

const imageUrl = computed(() => {
  const path = version.value.imagePath || model.value.imagePath;
  if (!path) return null;
  return path.replace(/^.*[\\/]backend[\\/]images/, "/images");
});

const parseMeta = (meta) => {
  try {
    if (typeof meta === "string") return JSON.parse(meta);
    return meta || {};
  } catch {
    return {};
  }
};

const galleryImages = computed(() => {
  const imgs = version.value.images || [];
  return imgs.map((img) => {
    const meta = parseMeta(img.meta);
    if (version.value.mode) {
      meta.mode = version.value.mode;
    }
    return {
      ...img,
      url: img.path.replace(/^.*[\\/]backend[\\/]images/, "/images"),
      parsedMeta: meta,
    };
  });
});

const weightDisplay = computed(() => {
  const weight = normalizeWeight(model.value.weight);
  return Number(weight.toFixed(2));
});

const fileName = computed(() => {
  if (!version.value.filePath) return "";
  return version.value.filePath.split(/[/\\]/).pop();
});

const fileBaseName = computed(() => {
  const name = fileName.value;
  if (!name) return "";
  return name.replace(/\.[^./\\]+$/, "");
});

const loraTag = computed(() => {
  const base = fileBaseName.value;
  if (!base) return "";
  return `<lora:${base}:1>`;
});

const formattedTrainedWords = computed(() => {
  if (!version.value.trainedWords) return "";
  return version.value.trainedWords
    .split(",")
    .map((word) => word.trim())
    .filter((word) => word.length)
    .join(", ");
});

const createdAtReadable = computed(() => {
  if (!version.value.createdAt) return "";
  return new Date(version.value.createdAt).toLocaleString();
});

const updatedAtReadable = computed(() => {
  if (!version.value.updatedAt) return "";
  return new Date(version.value.updatedAt).toLocaleString();
});

const versionSizeMb = computed(() => {
  if (!version.value.sizeKB) return "";
  return `${(version.value.sizeKB / 1024).toFixed(2)} MB`;
});

const hasVersionStats = computed(() => {
  return (
    !!version.value.sizeKB ||
    !!version.value.createdAt ||
    !!version.value.updatedAt ||
    !!version.value.sha256
  );
});

const fetchData = async () => {
  const { versionId } = route.params;
  const res = await axios.get(`/api/versions/${versionId}`);
  model.value = res.data.model;
  model.value.weight = normalizeWeight(model.value.weight);
  version.value = res.data.version;
};

onMounted(async () => {
  await fetchData();
  if (route.query.edit === "1") {
    startEdit();
  }
});

watch(isEditing, async (val) => {
  if (val) {
    await nextTick();
    quill = new Quill(editor.value, { theme: "snow" });
    quill.clipboard.dangerouslyPasteHTML(version.value.description || "");
  }
});

const deleteVersion = async () => {
  const choice = await showDeleteConfirm("Delete this version?");
  if (!choice) return;
  const files = choice === "deleteFiles" ? 1 : 0;
  await axios.delete(`/api/versions/${route.params.versionId}?files=${files}`);
  router.push("/");
};

const toggleNsfw = async () => {
  if (!version.value?.ID) return;
  const updated = { ...version.value, nsfw: !version.value.nsfw };
  togglingNsfw.value = true;
  try {
    await axios.put(`/api/versions/${version.value.ID}`, updated);
    version.value.nsfw = updated.nsfw;
    showToast("NSFW status updated", "success");
  } catch (err) {
    console.error(err);
    showToast("Failed to update NSFW status", "danger");
  } finally {
    togglingNsfw.value = false;
  }
};

const copyToClipboard = async (text, successMessage, errorMessage, logLabel) => {
  if (!text) return;

  const fallbackCopy = () => {
    const textarea = document.createElement("textarea");
    textarea.value = text;
    textarea.setAttribute("readonly", "");
    textarea.style.position = "absolute";
    textarea.style.left = "-9999px";
    textarea.style.opacity = "0";
    document.body.appendChild(textarea);
    textarea.select();
    textarea.setSelectionRange(0, textarea.value.length);
    const successful = document.execCommand("copy");
    document.body.removeChild(textarea);
    if (!successful) {
      throw new Error("Copy command failed");
    }
  };

  try {
    if (navigator?.clipboard?.writeText) {
      await navigator.clipboard.writeText(text);
    } else {
      fallbackCopy();
    }
  } catch (err) {
    try {
      fallbackCopy();
    } catch (fallbackErr) {
      console.error(`Failed to copy ${logLabel || "text"}`, fallbackErr);
      showToast(errorMessage, "danger");
      return;
    }
  }

  showToast(successMessage, "success");
};

const copyTrainedWords = async () => {
  await copyToClipboard(
    formattedTrainedWords.value,
    "Trained words copied",
    "Unable to copy trained words",
    "trained words",
  );
};

const copyFileBaseName = async () => {
  await copyToClipboard(
    fileBaseName.value,
    "Filename copied",
    "Unable to copy filename",
    "filename",
  );
};

const copyLoraTag = async () => {
  await copyToClipboard(
    loraTag.value,
    "LoRA tag copied",
    "Unable to copy LoRA tag",
    "LoRA tag",
  );
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
    version.value.description = quill.root.innerHTML;
  }
  model.value.weight = normalizeWeight(model.value.weight);
  await axios.put(`/api/models/${model.value.ID}`, model.value);
  try {
    await axios.put(`/api/versions/${version.value.ID}`, version.value);
  } catch (err) {
    if (err.response && err.response.status === 409) {
      showToast("Version ID already exists", "danger");
      return;
    }
    throw err;
  }
  isEditing.value = false;
  await fetchData();
  showToast("Saved", "success");
};

const refreshVersion = async (fields) => {
  await axios.post(`/api/versions/${version.value.ID}/refresh`, null, {
    params: { fields },
  });
  await fetchData();
  showToast("Updated", "success");
};

const updateMeta = async () => {
  if (!(await showConfirm("Pull latest metadata from CivitAI?"))) return;
  await refreshVersion("metadata");
};
const updateDesc = async () => {
  if (!(await showConfirm("Pull latest description from CivitAI?"))) return;
  await refreshVersion("description");
};
const updateImages = async () => {
  if (!(await showConfirm("Replace all images with the latest from CivitAI?")))
    return;
  await refreshVersion("images");
};
const updateAll = async () => {
  if (
    !(await showConfirm(
      "Update all data from CivitAI? This will replace images.",
    ))
  )
    return;
  await refreshVersion("all");
};

const goBack = () => {
  router.push({ path: "/", query: { scrollTo: version.value.ID } });
};

const setMainImage = async (img) => {
  await axios.post(`/api/versions/${version.value.ID}/main-image/${img.ID}`);
  await fetchData();
  showToast("Main image updated", "success");
};

const onGalleryFileChange = (e) => {
  galleryFile.value = e.target.files[0] || null;
};

const uploadGallery = async () => {
  if (!galleryFile.value) return;
  const fd = new FormData();
  fd.append("file", galleryFile.value);
  const res = await axios.post(
    `/api/versions/${version.value.ID}/images?type=${encodeURIComponent(version.value.type)}`,
    fd,
    { headers: { "Content-Type": "multipart/form-data" } },
  );
  version.value.images = version.value.images || [];
  version.value.images.push(res.data);
  galleryFile.value = null;
  showToast("Image uploaded", "success");
};

const removeImage = async (img) => {
  if (!(await showConfirm("Remove this image?"))) return;
  await axios.delete(`/api/versions/${version.value.ID}/images/${img.ID}`);
  if (version.value.imagePath === img.path) {
    await fetchData();
  } else {
    version.value.images = (version.value.images || []).filter(
      (i) => i.ID !== img.ID,
    );
  }
  showToast("Image removed", "success");
};

watch(
  () => version.value.type,
  (val) => {
    if (model.value) model.value.type = val;
  },
);

const onImageFileChange = (e) => {
  imageFile.value = e.target.files[0] || null;
};

const onModelFileChange = (e) => {
  modelFile.value = e.target.files[0] || null;
};

const uploadImage = async () => {
  if (!imageFile.value) return;
  const fd = new FormData();
  fd.append("file", imageFile.value);
  const res = await axios.post(
    `/api/versions/${version.value.ID}/upload?kind=image&type=${encodeURIComponent(version.value.type)}`,
    fd,
    { headers: { "Content-Type": "multipart/form-data" } },
  );
  version.value.imagePath = res.data.path;
  imageFile.value = null;
  showToast("Image uploaded", "success");
};

const uploadModel = async () => {
  if (!modelFile.value) return;
  const fd = new FormData();
  fd.append("file", modelFile.value);
  const res = await axios.post(
    `/api/versions/${version.value.ID}/upload?kind=file&type=${encodeURIComponent(version.value.type)}`,
    fd,
    { headers: { "Content-Type": "multipart/form-data" } },
  );
  version.value.filePath = res.data.path;
  modelFile.value = null;
  showToast("File uploaded", "success");
};
</script>
