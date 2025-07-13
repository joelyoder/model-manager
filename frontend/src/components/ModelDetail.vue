<template>
  <div class="detail">
    <button @click="goBack">⬅ Back</button>
    <h2>{{ model.name }}</h2>
    <img
      v-if="imageUrl"
      :src="imageUrl"
      :width="model.imageWidth"
      :height="model.imageHeight"
    />
    <div v-if="model.description" v-html="model.description"></div>
    <p v-if="model.tags">Tags: {{ model.tags.split(",").join(", ") }}</p>
    <p>Type: {{ model.type }}</p>
    <p>NSFW: {{ model.nsfw }}</p>
    <p>Created: {{ model.createdAt }}</p>
    <p>Updated: {{ model.updatedAt }}</p>
    <h3>Versions</h3>
    <ul>
      <li v-for="v in model.versions" :key="v.ID">
        {{ v.name }} - {{ v.baseModel }} -
        <span v-if="v.sizeKB">{{ (v.sizeKB / 1024).toFixed(2) }} MB</span>
      </li>
    </ul>
    <button @click="deleteModel">🗑 Delete</button>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import axios from "axios";

const router = useRouter();
const route = useRoute();
const model = ref({});

const imageUrl = computed(() => {
  if (!model.value.imagePath) return null;
  return model.value.imagePath.replace(/^.*\/backend\/images/, "/images");
});

const fetchModel = async () => {
  const { id } = route.params;
  const res = await axios.get(`/api/models/${id}`);
  model.value = res.data;
};

onMounted(fetchModel);

const deleteModel = async () => {
  if (!confirm("Delete this model and all files?")) return;
  await axios.delete(`/api/models/${route.params.id}`);
  router.push("/");
};

const goBack = () => {
  router.push("/");
};
</script>

<style scoped>
.detail {
  padding: 1rem;
}
img {
  max-width: 100%;
  height: auto;
  margin-bottom: 1rem;
}
</style>
