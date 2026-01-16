<template>
  <div>
    <div
      v-if="galleryImages.length"
      class="row row-cols-1 row-cols-md-2 row-cols-lg-3 g-4 mb-4"
    >
      <div v-for="img in galleryImages" :key="img.ID" class="col">
        <div class="card border-0 shadow-sm bg-dark-subtle h-100 overflow-hidden">
            <div class="position-relative">
                <img
                :src="img.url"
                :width="img.width"
                :height="img.height"
                class="card-img-top img-fluid"
                />
                 <span
                    v-if="img.path === currentImagePath"
                    class="position-absolute top-0 start-0 m-2 badge bg-success shadow-sm"
                >
                    Main Image
                </span>
            </div>
            
            <div class="card-body p-3">
                <div class="d-flex gap-2 mb-2">
                    <button
                        v-if="img.path !== currentImagePath"
                        @click="$emit('setMain', img)"
                        class="btn btn-outline-primary btn-sm flex-grow-1"
                    >
                        Set as Main
                    </button>
                    <button
                        @click="$emit('remove', img)"
                        class="btn btn-outline-danger btn-sm flex-grow-1"
                    >
                        Remove
                    </button>
                </div>

                <div v-if="Object.keys(img.parsedMeta || {}).length">
                    <button
                        class="btn btn-outline-secondary btn-sm w-100"
                        type="button"
                        data-bs-toggle="collapse"
                        :data-bs-target="'#meta-' + img.ID"
                        aria-expanded="false"
                        :aria-controls="'meta-' + img.ID"
                    >
                        Metadata
                    </button>
                    <div :id="'meta-' + img.ID" class="collapse mt-2">
                         <div class="bg-dark rounded p-2 small text-white">
                            <dl class="image-metadata mb-0">
                                <template v-for="(value, key) in img.parsedMeta" :key="key">
                                    <dt class="text-secondary fw-bold small text-uppercase mb-0">{{ key }}</dt>
                                    <dd class="mb-2">
                                        <div v-if="key.toLowerCase() === 'prompt'" class="d-flex gap-2 align-items-start justify-content-between">
                                            <span class="text-break">{{ value }}</span>
                                            <button 
                                                @click="copyToClipboard(value, 'Prompt copied', 'Failed to copy prompt')" 
                                                class="btn btn-outline-secondary btn-sm p-0 d-flex align-items-center justify-content-center border-0 flex-shrink-0" 
                                                title="Copy Prompt"
                                                style="width: 24px; height: 24px;"
                                            >
                                                <Icon icon="mdi:content-copy" width="14" height="14" />
                                            </button>
                                        </div>
                                        <span v-else class="text-break">{{ value }}</span>
                                    </dd>
                                </template>
                            </dl>
                        </div>
                    </div>
                </div>
            </div>
        </div>
      </div>
    </div>
    
    <div class="card border-0 shadow-sm bg-dark-subtle rounded-3 p-3">
        <label class="form-label text-secondary fw-bold small text-uppercase mb-2">Upload Image</label>
        <div class="input-group">
            <input 
                type="file" 
                @change="onFileChange" 
                class="form-control bg-dark border-0 text-white" 
            />
            <button @click="upload" class="btn btn-primary d-flex align-items-center gap-2">
                 <Icon icon="mdi:upload" width="20" height="20" />
                 Upload
            </button>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { Icon } from "@iconify/vue";
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

const copyToClipboard = async (text, successMessage, errorMessage) => {
  if (!text) return;
  try {
    await navigator.clipboard.writeText(text);
    showToast(successMessage, "success");
  } catch (err) {
    console.error("Failed to copy text", err);
    // Fallback
    try {
        const textarea = document.createElement("textarea");
        textarea.value = text;
        textarea.style.position = "fixed";
        textarea.style.left = "-9999px";
        document.body.appendChild(textarea);
        textarea.select();
        document.execCommand("copy");
        document.body.removeChild(textarea);
        showToast(successMessage, "success");
    } catch (fallbackErr) {
        console.error("Fallback copy failed", fallbackErr);
        showToast(errorMessage, "danger");
    }
  }
};
</script>
