import { ref } from "vue";
import axios from "axios";
import { showToast } from "../utils/ui";

export function useModelDetail() {
    const model = ref({});
    const version = ref({});
    const isEditing = ref(false);
    const togglingNsfw = ref(false);

    const normalizeWeight = (weight) => {
        const num = Number(weight);
        if (Number.isFinite(num) && num > 0) {
            return num;
        }
        return 1;
    };

    const fetchData = async (versionId) => {
        const res = await axios.get(`/api/versions/${versionId}`);
        model.value = res.data.model;
        model.value.weight = normalizeWeight(model.value.weight);
        version.value = res.data.version;
    };

    const toggleNsfw = async () => {
        if (!version.value?.ID) return;
        const updated = { ...version.value, nsfw: !version.value.nsfw };
        togglingNsfw.value = true;
        try {
            await axios.put(`/api/versions/${version.value.ID}`, updated);
            version.value.nsfw = updated.nsfw;
            showToast("NSFW status updated", "success");
        } catch (err) {
            console.error(err);
            showToast("Failed to update NSFW status", "danger");
        } finally {
            togglingNsfw.value = false;
        }
    };

    const deleteVersion = async (versionId, deleteFiles = false) => {
        const files = deleteFiles ? 1 : 0;
        await axios.delete(`/api/versions/${versionId}?files=${files}`);
    };

    const saveEdit = async () => {
        model.value.weight = normalizeWeight(model.value.weight);
        await axios.put(`/api/models/${model.value.ID}`, model.value);
        try {
            await axios.put(`/api/versions/${version.value.ID}`, version.value);
        } catch (err) {
            if (err.response && err.response.status === 409) {
                showToast("Version ID already exists", "danger");
                throw err;
            }
            throw err;
        }
        isEditing.value = false;
        // Caller should refresh data
    };

    const refreshVersion = async (fields) => {
        await axios.post(`/api/versions/${version.value.ID}/refresh`, null, {
            params: { fields },
        });
    };

    return {
        model,
        version,
        isEditing,
        togglingNsfw,
        fetchData,
        toggleNsfw,
        deleteVersion,
        saveEdit,
        refreshVersion,
    };
}
