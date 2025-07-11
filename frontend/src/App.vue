<template>
  <div class="app">
    <h1>📦 Local CivitAI Model Manager</h1>
    <ModelList />
  </div>
</template>

<script setup>
import ModelList from './components/ModelList.vue'
</script>

<style>
body {
  font-family: sans-serif;
  background: #f6f6f6;
  margin: 0;
}
h1 {
  padding: 1rem;
  text-align: center;
}
</style>
