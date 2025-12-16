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
            return true;
        } catch (e) {
            console.error(e);
            alert("Failed to dispatch: " + e.message);
            return false;
        } finally {
            isDispatching.value = false;
        }
    };

    return { dispatchAction, isDispatching };
}
