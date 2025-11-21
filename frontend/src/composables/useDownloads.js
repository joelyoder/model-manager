import { ref } from "vue";
import axios from "axios";
import { showToast } from "../utils/ui";

export function useDownloads() {
    const downloading = ref(false);
    const downloadProgress = ref(0);
    const canceling = ref(false);
    const cancelledByUser = ref(false);
    let progressInterval = null;

    const startDownload = async (versionId, modelId) => {
        downloading.value = true;
        downloadProgress.value = 0;
        cancelledByUser.value = false;
        canceling.value = false;

        // Start polling progress
        progressInterval = setInterval(async () => {
            try {
                const res = await axios.get("/api/download/progress");
                if (res.data && typeof res.data.progress === "number") {
                    downloadProgress.value = Math.round(res.data.progress);
                }
            } catch {
                // ignore
            }
        }, 1000);

        try {
            const buildSyncVersionUrl = (vid, { download, modelId } = {}) => {
                const params = new URLSearchParams();
                if (modelId) params.set("modelId", String(modelId));
                if (download !== undefined) params.set("download", String(download));
                const query = params.toString();
                return `/api/sync/version/${vid}${query ? `?${query}` : ""}`;
            };

            await axios.post(buildSyncVersionUrl(versionId, { download: 1, modelId }));

            if (!cancelledByUser.value) {
                showToast("Download completed", "success");
            }
        } catch (err) {
            if (!cancelledByUser.value) {
                console.error(err);
                showToast("Download failed", "danger");
            }
        } finally {
            clearInterval(progressInterval);
            downloading.value = false;
            downloadProgress.value = 0;
        }
    };

    const cancelDownload = async () => {
        canceling.value = true;
        try {
            await axios.post("/api/download/cancel");
            cancelledByUser.value = true;
            showToast("Download cancelled", "info");
        } catch (err) {
            console.error(err);
            showToast("Failed to cancel download", "danger");
        } finally {
            canceling.value = false;
        }
    };

    return {
        downloading,
        downloadProgress,
        canceling,
        startDownload,
        cancelDownload,
    };
}
