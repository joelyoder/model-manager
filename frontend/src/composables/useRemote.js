import axios from "axios";
import { ref } from "vue";

export function useRemote() {
    const isDispatching = ref(false);

    const dispatchAction = async (action, model, version, clientId) => {
        if (!clientId) {
            clientId = "My-Desktop-PC"; // Default
            // TODO: Fetch from settings
        }

        isDispatching.value = true;
        try {
            // Determine filename
            let filename = version.filePath ? version.filePath.split(/[/\\]/).pop() : null;
            if (!filename) {
                // Construct a filename
                const safeName = (model.name + "_" + version.name).replace(/[^a-zA-Z0-9_-]/g, "_");
                filename = safeName + ".safetensors"; // Default assumption
            }

            // Determine subdirectory
            let subdir = version.type || model.type || "Checkpoint";

            const payload = {
                action: action,
                url: version.downloadUrl,
                filename: filename,
                subdirectory: subdir,
                model_version_id: version.ID,
                client_id: clientId
            };

            await axios.post("/api/remote/dispatch", payload);

            // Optimistic Update
            if (action === 'download') {
                version.clientStatus = 'pending';
                pollStatus(version);
            } else if (action === 'delete') {
                version.clientStatus = null;
                // We could poll to confirm deletion, but null is fine for now
            }

            return true;
        } catch (e) {
            console.error(e);
            alert("Failed to dispatch: " + e.message);
            return false;
        } finally {
            isDispatching.value = false;
        }
    };

    const pollStatus = async (version) => {
        console.log(`[Remote] Starting poll for version ${version.ID}`);
        const pollInterval = 2000;
        const maxAttempts = 1500;
        let attempts = 0;

        const check = async () => {
            if (attempts >= maxAttempts) {
                console.log(`[Remote] Polling timed out for version ${version.ID}`);
                return;
            }
            attempts++;

            try {
                const res = await axios.get(`/api/versions/${version.ID}`);
                const currentStatus = res.data.version.clientStatus;

                console.log(`[Remote] Poll #${attempts} for ${version.ID}: status=${currentStatus}`);

                if (currentStatus !== 'pending') {
                    console.log(`[Remote] Status changed to ${currentStatus}. Updating UI.`);
                    version.clientStatus = currentStatus;
                    return;
                }

                setTimeout(check, pollInterval);
            } catch (e) {
                console.error("Polling validation failed", e);
            }
        };

        setTimeout(check, pollInterval);
    };


    return { dispatchAction, isDispatching };
}
