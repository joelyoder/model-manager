import { describe, it, expect, vi, beforeEach } from "vitest";

vi.mock("bootstrap", () => {
  return {
    Toast: vi.fn().mockImplementation(() => ({ show: vi.fn() })),
    Modal: vi.fn().mockImplementation(() => ({ show: vi.fn(), hide: vi.fn() })),
  };
});

import { showToast, showConfirm, showDeleteConfirm } from "./ui";

beforeEach(() => {
  document.body.innerHTML = `
    <div id="toast-container"></div>
    <div id="confirm-modal">
      <div class="modal-body"></div>
      <button id="confirm-ok" class="d-none"></button>
      <button class="btn-secondary"></button>
      <button id="confirm-delete" class="d-none"></button>
      <button id="confirm-delete-files" class="d-none"></button>
    </div>`;
});

describe("showToast", () => {
  it("adds toast to container", () => {
    showToast("Hello", "success");
    const container = document.getElementById("toast-container");
    const toast = container.querySelector(".toast");
    expect(toast).not.toBeNull();
    expect(toast.className).toContain("text-bg-success");
    expect(toast.textContent).toContain("Hello");
  });
});

describe("showConfirm", () => {
  it("resolves true on ok", async () => {
    const p = showConfirm("Are you sure?");
    const ok = document.getElementById("confirm-ok");
    const del = document.getElementById("confirm-delete");
    const delFiles = document.getElementById("confirm-delete-files");
    expect(ok.classList.contains("d-none")).toBe(false);
    expect(del.classList.contains("d-none")).toBe(true);
    expect(delFiles.classList.contains("d-none")).toBe(true);
    ok.click();
    await expect(p).resolves.toBe(true);
  });

  it("resolves false on cancel", async () => {
    const p = showConfirm("Cancel?");
    const cancel = document.querySelector(".btn-secondary");
    cancel.click();
    await expect(p).resolves.toBe(false);
  });
});

describe("showDeleteConfirm", () => {
  it("resolves 'delete' on delete click", async () => {
    const p = showDeleteConfirm("Delete?");
    const del = document.getElementById("confirm-delete");
    const ok = document.getElementById("confirm-ok");
    const delFiles = document.getElementById("confirm-delete-files");
    expect(ok.classList.contains("d-none")).toBe(true);
    expect(del.classList.contains("d-none")).toBe(false);
    expect(delFiles.classList.contains("d-none")).toBe(false);
    del.click();
    await expect(p).resolves.toBe("delete");
  });

  it("resolves 'deleteFiles' on delete files click", async () => {
    const p = showDeleteConfirm("Delete files?");
    const delFiles = document.getElementById("confirm-delete-files");
    delFiles.click();
    await expect(p).resolves.toBe("deleteFiles");
  });

  it("resolves null on cancel", async () => {
    const p = showDeleteConfirm("Cancel?");
    const cancel = document.querySelector(".btn-secondary");
    cancel.click();
    await expect(p).resolves.toBe(null);
  });
});
