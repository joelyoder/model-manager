<template>
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
            class="btn btn-warning d-none"
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
</template>

<script setup>
// This component provides the DOM structure required by utils/ui.js
// It does not handle logic itself, but allows the imperative UI utils to work.
</script>
