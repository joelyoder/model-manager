<template>
  <div>
    <div class="row">
      <div class="col-md-4">
        <img v-if="imageUrl" :src="imageUrl" class="img-fluid rounded-3 shadow-sm mb-4" />
      </div>
      <div class="col-md-8">
        <h3 class="fw-bold text-white mb-1">{{ model.name }}</h3>
        <h5 v-if="version.name" class="text-white-50 mb-3">{{ version.name }}</h5>
        
        <div class="d-flex flex-wrap align-items-center gap-2 mb-4">
          <span v-if="version.type" class="badge rounded-pill" :class="getBadgeColor(version.type)">
            {{ version.type }}
          </span>
          <span
            v-if="version.baseModel"
            class="badge rounded-pill"
             :class="getBadgeColor(version.baseModel)"
          >
            {{ version.baseModel }}
          </span>
        </div>

        <div class="row g-3">
            <!-- Tags -->
            <div class="col-12" v-if="version.tags">
                <label class="form-label text-secondary fw-bold small text-uppercase mb-1">Tags</label>
                <div class="bg-dark rounded p-2 text-white small">
                    {{ version.tags.split(",").join(", ") }}
                </div>
            </div>

            <!-- Trained Words -->
            <div class="col-12" v-if="version.trainedWords">
                <label class="form-label text-secondary fw-bold small text-uppercase mb-1">Trained Words</label>
                <div class="bg-dark rounded p-2 text-white small d-flex align-items-center justify-content-between">
                    <span class="text-break me-2">{{ formattedTrainedWords }}</span>
                    <button
                        type="button"
                        class="btn btn-outline-secondary btn-sm p-0 d-flex align-items-center justify-content-center border-0"
                        @click="copyTrainedWords"
                        aria-label="Copy trained words"
                        style="width: 24px; height: 24px;"
                        title="Copy"
                    >
                        <Icon icon="mdi:content-copy" width="14" height="14" />
                    </button>
                </div>
            </div>

            <!-- Weight -->
            <div class="col-6 col-md-4">
                <label class="form-label text-secondary fw-bold small text-uppercase mb-1">Weight</label>
                 <div class="bg-dark rounded p-2 text-white small">
                    {{ weightDisplay }}
                </div>
            </div>

             <!-- File -->
            <div class="col-12" v-if="version.filePath">
                <label class="form-label text-secondary fw-bold small text-uppercase mb-1">File</label>
                <div class="bg-dark rounded p-2 text-white small d-flex align-items-center justify-content-between">
                    <span class="text-truncate me-2">{{ fileName }}</span>
                    <div class="d-flex gap-1">
                        <button
                        type="button"
                        class="btn btn-outline-secondary btn-sm p-0 d-flex align-items-center justify-content-center border-0"
                        @click="copyFileBaseName"
                        aria-label="Copy filename"
                        title="Copy filename"
                        style="width: 24px; height: 24px;"
                        >
                            <Icon icon="mdi:file-document-outline" width="14" height="14" />
                        </button>
                         <button
                        type="button"
                        class="btn btn-outline-secondary btn-sm p-0 d-flex align-items-center justify-content-center border-0"
                        @click="copyLoraTag"
                        aria-label="Copy LoRA tag"
                         title="Copy LoRA tag"
                        style="width: 24px; height: 24px;"
                        >
                            <Icon icon="mdi:tag-text-outline" width="14" height="14" />
                        </button>
                    </div>
                </div>
            </div>

             <!-- Model URL -->
            <div class="col-12" v-if="version.modelUrl">
                <label class="form-label text-secondary fw-bold small text-uppercase mb-1">Model URL</label>
                <div class="bg-dark rounded p-2 text-white small text-truncate">
                   <a :href="version.modelUrl" target="_blank" class="text-primary text-decoration-none">{{ version.modelUrl }}</a>
                </div>
            </div>
        </div>
      </div>
    </div>

    <!-- Description -->
    <div v-if="version.description" class="mt-4 mb-4">
         <label class="form-label text-secondary fw-bold small text-uppercase mb-1">Description</label>
        <div class="bg-dark rounded p-3 text-white small model-description" v-html="version.description"></div>
    </div>
    
    <!-- Stats Grid -->
    <div v-if="hasVersionStats" class="row row-cols-2 row-cols-md-4 g-3 mb-4">
      <div v-if="version.createdAt" class="col">
          <label class="form-label text-secondary fw-bold small text-uppercase mb-1">Created</label>
          <div class="bg-dark rounded p-2 text-white small text-center">
            {{ createdAtReadable }}
          </div>
      </div>
      <div v-if="version.updatedAt" class="col">
           <label class="form-label text-secondary fw-bold small text-uppercase mb-1">Updated</label>
          <div class="bg-dark rounded p-2 text-white small text-center">
            {{ updatedAtReadable }}
          </div>
      </div>
      <div v-if="version.sizeKB" class="col">
           <label class="form-label text-secondary fw-bold small text-uppercase mb-1">Size</label>
          <div class="bg-dark rounded p-2 text-white small text-center">
            {{ versionSizeMb }}
          </div>
      </div>
      <div v-if="version.sha256" class="col">
           <label class="form-label text-secondary fw-bold small text-uppercase mb-1">SHA256</label>
          <div class="bg-dark rounded p-2 text-white small text-center text-truncate" :title="version.sha256">
             {{ version.sha256.substring(0, 10) }}...
          </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { Icon } from "@iconify/vue";
import { showToast } from "../utils/ui";
import { getBadgeColor } from "../utils/colors";

const props = defineProps({
  model: Object,
  version: Object,
});

const imageUrl = computed(() => {
  const path = props.version.imagePath || props.model.imagePath;
  if (!path) return null;
  return path.replace(/^.*[\\/]backend[\\/]images/, "/images");
});

const normalizeWeight = (weight) => {
  const num = Number(weight);
  if (Number.isFinite(num) && num > 0) {
    return num;
  }
  return 1;
};

const normalizedWeight = computed(() => normalizeWeight(props.model.weight));

const weightDisplay = computed(() => {
  const weight = normalizedWeight.value;
  return Number(weight.toFixed(2));
});

const fileName = computed(() => {
  if (!props.version.filePath) return "";
  return props.version.filePath.split(/[/\\]/).pop();
});

const fileBaseName = computed(() => {
  const name = fileName.value;
  if (!name) return "";
  return name.replace(/\.[^./\\]+$/, "");
});

const loraTag = computed(() => {
  const base = fileBaseName.value;
  if (!base) return "";
  const weight = Number(normalizedWeight.value.toFixed(2));
  return `<lora:${base}:${weight}>`;
});

const formattedTrainedWords = computed(() => {
  if (!props.version.trainedWords) return "";
  return props.version.trainedWords
    .split(",")
    .map((word) => word.trim())
    .filter((word) => word.length)
    .join(", ");
});

const createdAtReadable = computed(() => {
  if (!props.version.createdAt) return "";
  return new Date(props.version.createdAt).toLocaleString();
});

const updatedAtReadable = computed(() => {
  if (!props.version.updatedAt) return "";
  return new Date(props.version.updatedAt).toLocaleString();
});

const versionSizeMb = computed(() => {
  if (!props.version.sizeKB) return "";
  return `${(props.version.sizeKB / 1024).toFixed(2)} MB`;
});

const hasVersionStats = computed(() => {
  return (
    !!props.version.sizeKB ||
    !!props.version.createdAt ||
    !!props.version.updatedAt ||
    !!props.version.sha256
  );
});

const copyToClipboard = async (
  text,
  successMessage,
  errorMessage,
  logLabel
) => {
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
  } catch {
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
    "trained words"
  );
};

const copyFileBaseName = async () => {
  await copyToClipboard(
    fileBaseName.value,
    "Filename copied",
    "Unable to copy filename",
    "filename"
  );
};

const copyLoraTag = async () => {
  await copyToClipboard(
    loraTag.value,
    "LoRA tag copied",
    "Unable to copy LoRA tag",
    "LoRA tag"
  );
};

// Local getBadgeColor removed in favor of utility
</script>
