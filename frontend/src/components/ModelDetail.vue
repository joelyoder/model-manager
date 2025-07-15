<template>
  <div class="px-4 container">
    <div class="mb-2 d-flex gap-2 pb-2" v-if="!isEditing">
      <button @click="goBack" class="btn btn-secondary">Back</button>
      <button @click="startEdit" class="btn btn-primary">Edit</button>
      <button @click="deleteVersion" class="btn btn-danger ms-auto">
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
          <img
            v-if="imageUrl"
            :src="imageUrl"
            :width="model.imageWidth"
            :height="model.imageHeight"
            class="img-fluid mb-4"
          />
        </div>
        <div class="col-md-8">
          <h2 class="fw-bold">{{ model.name }}</h2>
          <h3 v-if="version.name" class="mb-2">{{ version.name }}</h3>
          <table class="table mt-4">
            <tbody>
              <tr v-if="version.tags">
                <th>Tags</th>
                <td>{{ version.tags.split(",").join(", ") }}</td>
              </tr>
              <tr>
                <th>Type</th>
                <td>{{ version.type }}</td>
              </tr>
              <tr>
                <th>NSFW</th>
                <td>{{ version.nsfw }}</td>
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
              <tr v-if="version.createdAt">
                <th>Created</th>
                <td>{{ createdAtReadable }}</td>
              </tr>
              <tr v-if="version.updatedAt">
                <th>Updated</th>
                <td>{{ updatedAtReadable }}</td>
              </tr>
              <tr v-if="version.sha256">
                <th>SHA256</th>
                <td>
                  <code>{{ version.sha256 }}</code>
                </td>
              </tr>
              <tr v-if="version.downloadUrl">
                <th>Download URL</th>
                <td>
                  <a :href="version.downloadUrl" target="_blank">{{
                    version.downloadUrl
                  }}</a>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <div
        v-if="version.description"
        v-html="version.description"
        class="mb-4"
      ></div>
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
          <button
            v-if="img.path !== version.imagePath"
            @click="setMainImage(img)"
            class="btn btn-secondary btn-sm mt-1 w-100"
          >
            Set as Main
          </button>
          <span v-else class="badge text-bg-success d-block text-center mt-1">
            Main Image
          </span>
          <table
            v-if="Object.keys(img.parsedMeta || {}).length"
            class="table table-sm bg-body-secondary rounded mb-0 mt-1"
          >
            <tbody>
              <tr v-for="(value, key) in img.parsedMeta" :key="key">
                <th class="fw-normal">{{ key }}</th>
                <td>{{ value }}</td>
              </tr>
            </tbody>
          </table>
        </div>
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
        <input v-model="version.tags" class="form-control" />
      </div>
      <div class="mb-3">
        <label class="form-label">Type</label>
        <input v-model="version.type" class="form-control" />
      </div>
      <div class="form-check mb-3">
        <input
          type="checkbox"
          class="form-check-input"
          id="nsfw"
          v-model="version.nsfw"
        />
        <label class="form-check-label" for="nsfw">NSFW</label>
      </div>
      <div class="mb-3">
        <label class="form-label">Description</label>
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

const fileName = computed(() => {
  if (!version.value.filePath) return "";
  return version.value.filePath.split("/").pop();
});

const createdAtReadable = computed(() => {
  if (!version.value.createdAt) return "";
  return new Date(version.value.createdAt).toLocaleString();
});

const updatedAtReadable = computed(() => {
  if (!version.value.updatedAt) return "";
  return new Date(version.value.updatedAt).toLocaleString();
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
    quill.clipboard.dangerouslyPasteHTML(version.value.description || "");
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
    version.value.description = quill.root.innerHTML;
  }
  await axios.put(`/api/models/${model.value.ID}`, model.value);
  await axios.put(`/api/versions/${version.value.ID}`, version.value);
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
  router.push("/");
};

const setMainImage = async (img) => {
  await axios.post(`/api/versions/${version.value.ID}/main-image/${img.ID}`);
  await fetchData();
  showToast("Main image updated", "success");
};
</script>
