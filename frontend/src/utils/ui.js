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
    const okBtn = el.querySelector(".btn-primary");
    const cancelBtn = el.querySelector(".btn-secondary");
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
