<template>
  <div>
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
              v-if="img.path !== currentImagePath"
              @click="$emit('setMain', img)"
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
              @click="$emit('remove', img)"
              class="btn btn-outline-danger btn-sm mt-1 w-100"
            >
              Remove
            </button>
          </div>
        </div>
        <div v-if="Object.keys(img.parsedMeta || {}).length" class="mt-1">
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
      <input type="file" @change="onFileChange" class="form-control" />
      <button @click="upload" class="btn btn-secondary">Add Image</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import axios from "axios";
import { showToast } from "../utils/ui";

const props = defineProps({
  images: Array,
  currentImagePath: String,
  versionId: Number,
  versionType: String,
  versionMode: String,
});

const emit = defineEmits(["setMain", "remove", "uploaded"]);

const galleryFile = ref(null);

const parseMeta = (meta) => {
  try {
    if (typeof meta === "string") return JSON.parse(meta);
    return meta || {};
  } catch {
    return {};
  }
};

const galleryImages = computed(() => {
  const imgs = props.images || [];
  return imgs.map((img) => {
    const meta = parseMeta(img.meta);
    if (props.versionMode) {
      meta.mode = props.versionMode;
    }
    return {
      ...img,
      url: img.path.replace(/^.*[\\/]backend[\\/]images/, "/images"),
      parsedMeta: meta,
    };
  });
});

const onFileChange = (e) => {
  galleryFile.value = e.target.files[0] || null;
};

const upload = async () => {
  if (!galleryFile.value) return;
  const fd = new FormData();
  fd.append("file", galleryFile.value);
  try {
    await axios.post(
      `/api/versions/${props.versionId}/images?type=${encodeURIComponent(
        props.versionType
      )}`,
      fd,
      { headers: { "Content-Type": "multipart/form-data" } }
    );
    galleryFile.value = null;
    showToast("Image uploaded", "success");
    emit("uploaded");
  } catch (err) {
    console.error(err);
    showToast("Failed to upload image", "danger");
  }
};
</script>
