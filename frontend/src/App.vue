<template>
  <div class="min-vh-100">
    <h1 class="p-4 h2 text-center">Model Manager</h1>
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
