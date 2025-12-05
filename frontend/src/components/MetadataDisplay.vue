<template>
  <div>
    <div class="row">
      <div class="col-md-4">
        <img v-if="imageUrl" :src="imageUrl" class="img-fluid mb-4" />
      </div>
      <div class="col-md-8">
        <h2 class="fw-bold">{{ model.name }}</h2>
        <h3 v-if="version.name" class="mb-2">{{ version.name }}</h3>
        <div class="d-flex flex-wrap align-items-center gap-2 mb-3">
          <span v-if="version.type" class="badge rounded-pill text-bg-primary">
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
    <div v-if="hasVersionStats" class="row row-cols-1 row-cols-md-2 g-3 mb-4">
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
          <dd class="metadata-list__value">
            <code>{{ version.sha256 }}</code>
          </dd>
        </dl>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { Icon } from "@iconify/vue";
import { showToast } from "../utils/ui";

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
</script>
