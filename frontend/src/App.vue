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
      <router-link to="/utilities" class="btn btn-outline-secondary btn-sm"
        ><svg
          width="22px"
          height="22px"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          color="#ffffff"
        >
          <path
            d="M10.0503 10.6066L2.97923 17.6777C2.19818 18.4587 2.19818 19.725 2.97923 20.5061V20.5061C3.76027 21.2871 5.0266 21.2871 5.80765 20.5061L12.8787 13.435"
            stroke="#ffffff"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          ></path>
          <path
            d="M10.0502 10.6066C9.20638 8.45358 9.37134 5.6286 11.1109 3.88909C12.8504 2.14957 16.0606 1.76777 17.8284 2.82843L14.7877 5.8691L14.5051 8.98014L17.6161 8.69753L20.6568 5.65685C21.7175 7.42462 21.3357 10.6349 19.5961 12.3744C17.8566 14.1139 15.0316 14.2789 12.8786 13.435"
            stroke="#ffffff"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          ></path></svg
      ></router-link>
      <router-link to="/settings" class="btn btn-outline-secondary btn-sm"
        ><svg
          width="22px"
          height="22px"
          stroke-width="1.5"
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          color="#ffffff"
        >
          <path
            d="M12 15C13.6569 15 15 13.6569 15 12C15 10.3431 13.6569 9 12 9C10.3431 9 9 10.3431 9 12C9 13.6569 10.3431 15 12 15Z"
            stroke="#ffffff"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          ></path>
          <path
            d="M19.6224 10.3954L18.5247 7.7448L20 6L18 4L16.2647 5.48295L13.5578 4.36974L12.9353 2H10.981L10.3491 4.40113L7.70441 5.51596L6 4L4 6L5.45337 7.78885L4.3725 10.4463L2 11V13L4.40111 13.6555L5.51575 16.2997L4 18L6 20L7.79116 18.5403L10.397 19.6123L11 22H13L13.6045 19.6132L16.2551 18.5155C16.6969 18.8313 18 20 18 20L20 18L18.5159 16.2494L19.6139 13.598L21.9999 12.9772L22 11L19.6224 10.3954Z"
            stroke="#ffffff"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          ></path></svg
      ></router-link>
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
            <button id="confirm-ok" type="button" class="btn btn-primary">
              OK
            </button>
            <button
              id="confirm-delete"
              type="button"
              class="btn btn-primary d-none"
            >
              Delete Version
            </button>
            <button
              id="confirm-delete-files"
              type="button"
              class="btn btn-danger d-none"
            >
              Delete Version &amp; Files
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup></script>
