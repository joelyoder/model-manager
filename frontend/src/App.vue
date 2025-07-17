<template>
  <div class="min-vh-100">
    <div class="d-flex p-4 gap-3 align-items-center">
      <svg
        width="30px"
        height="30px"
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
      <router-link
        to="/"
        class="h5 text-primary-emphasis my-auto fw-bold text-decoration-none me-auto"
        >Model Manager</router-link
      >
      <router-link to="/data" class="btn btn-outline-secondary btn-sm"
        >Import/Export</router-link
      >
      <router-link to="/settings" class="btn btn-outline-secondary btn-sm"
        >Settings</router-link
      >
    </div>
    <router-view />
    <div
      id="toast-container"
      class="toast-container position-fixed bottom-0 end-0 p-3"
    ></div>
    <div class="modal fade" id="confirm-modal" tabindex="-1">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-body">Are you sure?</div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-secondary"
              data-bs-dismiss="modal"
            >
              Cancel
            </button>
            <button type="button" class="btn btn-primary">OK</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup></script>
