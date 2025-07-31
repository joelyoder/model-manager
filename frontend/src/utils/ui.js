import { Toast, Modal } from "bootstrap";

export function showToast(message, variant = "primary") {
  const container = document.getElementById("toast-container");
  const wrapper = document.createElement("div");
  wrapper.className = `toast align-items-center text-bg-${variant}`;
  wrapper.setAttribute("role", "alert");
  wrapper.setAttribute("aria-live", "assertive");
  wrapper.setAttribute("aria-atomic", "true");
  wrapper.innerHTML = `<div class="d-flex"><div class="toast-body">${message}</div><button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button></div>`;
  container.appendChild(wrapper);
  const toast = new Toast(wrapper, { delay: 5000 });
  toast.show();
}

export function showConfirm(message) {
  return new Promise((resolve) => {
    const el = document.getElementById("confirm-modal");
    el.querySelector(".modal-body").textContent = message;
    const okBtn = el.querySelector("#confirm-ok");
    const cancelBtn = el.querySelector(".btn-secondary");
    const delBtn = el.querySelector("#confirm-delete");
    const delFilesBtn = el.querySelector("#confirm-delete-files");
    delBtn.classList.add("d-none");
    delFilesBtn.classList.add("d-none");
    okBtn.classList.remove("d-none");
    const modal = new Modal(el);
    const cleanup = (result) => {
      okBtn.removeEventListener("click", ok);
      cancelBtn.removeEventListener("click", cancel);
      resolve(result);
    };
    const ok = () => {
      modal.hide();
      cleanup(true);
    };
    const cancel = () => {
      modal.hide();
      cleanup(false);
    };
    okBtn.addEventListener("click", ok);
    cancelBtn.addEventListener("click", cancel);
    modal.show();
  });
}

export function showDeleteConfirm(message) {
  return new Promise((resolve) => {
    const el = document.getElementById("confirm-modal");
    el.querySelector(".modal-body").textContent = message;
    const okBtn = el.querySelector("#confirm-ok");
    const delBtn = el.querySelector("#confirm-delete");
    const delFilesBtn = el.querySelector("#confirm-delete-files");
    const cancelBtn = el.querySelector(".btn-secondary");
    okBtn.classList.add("d-none");
    delBtn.classList.remove("d-none");
    delFilesBtn.classList.remove("d-none");
    const modal = new Modal(el);
    const cleanup = (result) => {
      delBtn.removeEventListener("click", doDelete);
      delFilesBtn.removeEventListener("click", doDeleteFiles);
      cancelBtn.removeEventListener("click", cancel);
      resolve(result);
    };
    const doDelete = () => {
      modal.hide();
      cleanup("delete");
    };
    const doDeleteFiles = () => {
      modal.hide();
      cleanup("deleteFiles");
    };
    const cancel = () => {
      modal.hide();
      cleanup(null);
    };
    delBtn.addEventListener("click", doDelete);
    delFilesBtn.addEventListener("click", doDeleteFiles);
    cancelBtn.addEventListener("click", cancel);
    modal.show();
  });
}
