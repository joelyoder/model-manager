<template>
  <div class="details" v-if="model">
    <button @click="$router.back()">⬅ Back</button>
    <h2>{{ model.name }}</h2>
    <img v-if="imageUrl" :src="imageUrl" :width="model.imageWidth" :height="model.imageHeight" />
    <div class="description" v-html="model.description"></div>
    <p v-if="model.tags">Tags: {{ model.tags.split(',').join(', ') }}</p>
    <p>Type: {{ model.type }}</p>
    <p>NSFW: {{ model.nsfw ? 'Yes' : 'No' }}</p>
    <p>Created: {{ model.createdAt }}</p>
    <p>Updated: {{ model.updatedAt }}</p>

    <h3>Versions</h3>
    <ul>
      <li v-for="v in model.versions" :key="v.ID">
        {{ v.name }} - {{ v.baseModel }} - {{ (v.sizeKB / 1024).toFixed(2) }} MB
      </li>
    </ul>

    <button class="delete" @click="deleteModel">🗑 Delete</button>
  </div>
  <div v-else>Loading...</div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const model = ref(null)
const imageUrl = ref('')

const fetchModel = async () => {
  const id = route.params.id
  const res = await axios.get(`/api/models/${id}`)
  const data = res.data
  imageUrl.value = data.imagePath ? data.imagePath.replace(/^.*\/backend\/images/, '/images') : null
  model.value = data
}

const deleteModel = async () => {
  if (!confirm('Delete this model and all versions?')) return
  const id = route.params.id
  await axios.delete(`/api/models/${id}`)
  router.push('/')
}

onMounted(fetchModel)
</script>

<style scoped>
.details {
  padding: 1rem;
}
img {
  max-width: 100%;
  height: auto;
  margin-bottom: 1rem;
}
.delete {
  margin-top: 1rem;
}
.description {
  margin: 1rem 0;
}
</style>
