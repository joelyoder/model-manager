<template>
  <div class="min-vh-100">
    <div class="d-flex p-3 align-items-center sticky-top bg-dark">
      <!-- Left: Filter Button (Mobile Only) -->
      <div class="d-flex align-items-center d-md-none z-1">
        <button
          v-if="showFilterButton"
          class="btn btn-dark bg-opacity-25 btn-sm d-inline-flex align-items-center justify-content-center text-secondary-emphasis"
          @click="showSidebar = !showSidebar"
          aria-label="Toggle Filters"
          title="Toggle Filters"
          style="width: 32px; height: 32px;"
        >
          <Icon icon="mdi:filter" width="20" height="20" />
        </button>
      </div>

       <!-- Left: Add Model Button (Visible everywhere on ModelList) -->
      <div class="d-flex align-items-center z-1 ms-2">
        <button
          v-if="isModelList"
          class="btn btn-dark bg-opacity-25 btn-sm d-inline-flex align-items-center justify-content-center text-secondary-emphasis"
          @click="showAddPanel = !showAddPanel"
          aria-label="Add Model"
          title="Add Model"
          style="width: 32px; height: 32px;"
        >
           <Icon :icon="showAddPanel ? 'mdi:close' : 'mdi:plus'" width="20" height="20" />
        </button>
      </div>

      <!-- Center: Logo (Absolute Centered) -->
      <div class="position-absolute top-50 start-50 translate-middle pointer-events-none z-3">
        <router-link to="/" class="d-flex align-items-center text-decoration-none text-light" style="pointer-events: auto" aria-label="Home">
          <svg
            width="32px"
            height="32px"
            stroke-width="1.5"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            color="#fff"
          >
            <path
              d="M21 7.35304L21 16.647C21 16.8649 20.8819 17.0656 20.6914 17.1715L12.2914 21.8381C12.1102 21.9388 11.8898 21.9388 11.7086 21.8381L3.30861 17.1715C3.11814 17.0656 3 16.8649 3 16.647L2.99998 7.35304C2.99998 7.13514 3.11812 6.93437 3.3086 6.82855L11.7086 2.16188C11.8898 2.06121 12.1102 2.06121 12.2914 2.16188L20.6914 6.82855C20.8818 6.93437 21 7.13514 21 7.35304Z"
              stroke="#fff"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></path>
            <path
              d="M3.52844 7.29357L11.7086 11.8381C11.8898 11.9388 12.1102 11.9388 12.2914 11.8381L20.5 7.27777"
              stroke="#fff"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></path>
            <path
              d="M12 21L12 12"
              stroke="#fff"
              stroke-width="1.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            ></path>
          </svg>
        </router-link>
      </div>

      <!-- Right: Actions -->
      <div class="d-flex gap-2 flex-grow-1 justify-content-end ms-auto z-1">
        <router-link
          to="/collections"
          class="btn btn-dark bg-opacity-25 btn-sm d-inline-flex align-items-center justify-content-center text-secondary-emphasis"
          aria-label="Collections"
          title="Collections"
          style="width: 32px; height: 32px;"
        >
          <Icon icon="mdi:folder-multiple" width="20" height="20" />
        </router-link>
        <router-link
          to="/utilities"
          class="btn btn-dark bg-opacity-25 btn-sm d-inline-flex align-items-center justify-content-center text-secondary-emphasis"
          aria-label="Utilities"
          title="Utilities"
          style="width: 32px; height: 32px;"
        >
          <Icon icon="mdi:wrench" width="20" height="20" />
        </router-link>
        <router-link
          to="/settings"
          class="btn btn-dark bg-opacity-25 btn-sm d-inline-flex align-items-center justify-content-center text-secondary-emphasis"
          aria-label="Settings"
          title="Settings"
          style="width: 32px; height: 32px;"
        >
          <Icon icon="mdi:cog" width="20" height="20" />
        </router-link>
      </div>
    </div>
    <router-view />

     <!-- Add Model Slideout Panel -->
    <div 
        class="add-panel-slideout bg-dark shadow-lg" 
        :class="{ 'show': showAddPanel }"
    >
        <div class="d-flex justify-content-between align-items-center p-3 border-bottom border-dark-subtle">
            <h5 class="m-0">Add Model</h5>
            <button type="button" class="btn-close btn-close-white" aria-label="Close" @click="showAddPanel = false"></button>
        </div>
        <div class="p-3">
             <AddModelPanel
                v-if="showAddPanel" 
                @createManual="createManualModel"
                @added="onModelAdded"
            />
        </div>
    </div>
    <!-- Overlay for Add Panel -->
    <div v-if="showAddPanel" class="sidebar-overlay" @click="showAddPanel = false"></div>

    <BackToTop />
    <div
      id="toast-container"
      class="toast-container position-fixed bottom-0 start-0 p-3"
    ></div>
    <ConfirmModal />
  </div>
</template>

<script setup>
import BackToTop from "./components/BackToTop.vue";
import ConfirmModal from "./components/ConfirmModal.vue";
import AddModelPanel from "./components/AddModelPanel.vue";
import { Icon } from "@iconify/vue";
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useModels } from "./composables/useModels";
import axios from "axios";
import { showToast } from "./utils/ui";

const route = useRoute();
const router = useRouter();
const { showSidebar, showAddPanel, fetchModels } = useModels();

const isModelList = computed(() => route.name === "ModelList" || route.path === "/");
const showFilterButton = computed(() => isModelList.value || route.name === "CollectionDetail");

const onModelAdded = async () => {
    await fetchModels();
    showAddPanel.value = false;
};

const createManualModel = async () => {
  try {
    const res = await axios.post("/api/models", {
      name: "New Model",
      type: "Checkpoint",
    });
    await fetchModels();
    showAddPanel.value = false;
    router.push({
      name: "ModelDetail",
      params: { versionId: res.data.version.ID },
      query: { edit: "1" },
    });
  } catch (err) {
    console.error(err);
    showToast("Failed to create model", "danger");
  }
};
</script>

<style scoped>
.add-panel-slideout {
    position: fixed;
    top: 0;
    left: 0;
    height: 100vh;
    width: 100%;
    max-width: 1000px;
    z-index: 1050;
    transform: translateX(-100%);
    transition: transform 0.3s ease-in-out;
    overflow-y: auto;
}

.add-panel-slideout.show {
    transform: translateX(0);
}

.sidebar-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0,0,0,0.5);
    z-index: 1040;
}
</style>
