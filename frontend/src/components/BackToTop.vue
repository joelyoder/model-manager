<template>
  <button
    class="btn btn-primary back-to-top d-flex align-items-center justify-content-center p-0 shadow-lg"
    v-if="visible"
    @click="scrollToTop"
    style="width: 40px; height: 40px; border-radius: 50%;"
  >
    <Icon icon="mdi:arrow-up" width="24" height="24" />
  </button>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import { Icon } from "@iconify/vue";

const visible = ref(false);

const handleScroll = () => {
  visible.value = window.scrollY > 200;
};

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: "smooth" });
};

onMounted(() => {
  window.addEventListener("scroll", handleScroll);
});

onUnmounted(() => {
  window.removeEventListener("scroll", handleScroll);
});
</script>

<style scoped>
.back-to-top {
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 1050;
}
</style>
