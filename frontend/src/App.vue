<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-gray-900 dark:text-gray-100">
    <header class="p-6 flex items-center justify-between max-w-screen-lg mx-auto">
      <h1 class="text-2xl font-bold">📦 Local CivitAI Model Manager</h1>
      <button @click="toggleDark" class="p-2 rounded text-xl">
        <span v-if="isDark">🌞</span>
        <span v-else>🌙</span>
      </button>
    </header>
    <div class="max-w-screen-lg mx-auto px-4 pb-8">
      <router-view />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watchEffect } from "vue";

const isDark = ref(false);

const applyMode = (value) => {
  const root = document.documentElement;
  root.classList.toggle("dark", value);
};

onMounted(() => {
  const stored = localStorage.getItem("theme");
  if (stored === "dark" || stored === "light") {
    isDark.value = stored === "dark";
  } else {
    isDark.value = window.matchMedia("(prefers-color-scheme: dark)").matches;
  }
});

watchEffect(() => {
  applyMode(isDark.value);
  localStorage.setItem("theme", isDark.value ? "dark" : "light");
});

const toggleDark = () => {
  isDark.value = !isDark.value;
};
</script>
