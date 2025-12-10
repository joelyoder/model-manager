import { ref, computed, watch } from "vue";
import axios from "axios";
import debounce from "../utils/debounce";

export function useModels() {
  const models = ref([]);
  const search = ref("");
  const tagsSearch = ref("");
  const selectedCategory = ref("");
  const selectedBaseModel = ref("");
  const selectedModelType = ref("");
  const nsfwFilter = ref("no");
  const page = ref(1);
  const total = ref(0);
  const limit = 50;
  const initialized = ref(false);
  const baseModels = ref([]);

  const localStorageKey = "modelListState";

  const modelTypes = [
    "Checkpoint",
    "TextualInversion",
    "Hypernetwork",
    "AestheticGradient",
    "LORA",
    "LoCon",
    "DoRA",
    "Controlnet",
    "Upscaler",
    "MotionModule",
    "VAE",
    "Wildcards",
    "Poses",
    "Workflows",
    "Detection",
    "Other",
  ];

  const categories = [
    "character",
    "style",
    "concept",
    "clothing",
    "base model",
    "poses",
    "background",
    "tool",
    "vehicle",
    "buildings",
    "objects",
    "assets",
    "animal",
    "action",
  ];

  const saveState = () => {
    localStorage.setItem(
      localStorageKey,
      JSON.stringify({
        search: search.value,
        tagsSearch: tagsSearch.value,
        selectedCategory: selectedCategory.value,
        selectedBaseModel: selectedBaseModel.value,
        selectedModelType: selectedModelType.value,
        nsfwFilter: nsfwFilter.value,
        page: page.value,
      })
    );
  };

  const debouncedSave = debounce(saveState, 300);

  // Helper to normalize image paths for display
  // Handles old absolute paths (with /backend/images) and new relative paths
  const normalizeImagePath = (path) => {
    if (!path) return null;
    // Normalize backslashes to forward slashes
    let normalized = path.replace(/\\/g, "/");
    // If it's an old absolute path containing /backend/images, extract the relative part
    if (normalized.includes("/backend/images/")) {
      normalized = normalized.replace(/^.*\/backend\/images/, "/images");
    }
    // If path doesn't start with / or http, prepend /images/
    if (!normalized.startsWith("/") && !normalized.startsWith("http")) {
      normalized = "/images/" + normalized;
    }
    return normalized;
  };

  const mapModel = (model) => {
    const imageUrl = normalizeImagePath(model.imagePath);
    const versionsMap = new Map();
    (model.versions || []).forEach((v) => {
      if (!versionsMap.has(v.ID)) {
        const vImage = normalizeImagePath(v.imagePath);
        versionsMap.set(v.ID, { ...v, imageUrl: vImage });
      }
    });
    const versions = Array.from(versionsMap.values());
    return {
      ...model,
      versions,
      imageUrl,
    };
  };

  const fetchModels = async () => {
    const params = { page: page.value, limit, includeVersions: 1 };
    if (search.value) params.search = search.value;
    if (selectedBaseModel.value) params.baseModel = selectedBaseModel.value;
    if (selectedModelType.value) params.modelType = selectedModelType.value;
    if (nsfwFilter.value) params.nsfwFilter = nsfwFilter.value;
    const tagParts = [];
    if (selectedCategory.value) tagParts.push(selectedCategory.value);
    if (tagsSearch.value.trim()) tagParts.push(tagsSearch.value);
    if (tagParts.length) params.tags = tagParts.join(",");
    const res = await axios.get("/api/models", { params });
    models.value = res.data.map(mapModel);
  };

  const fetchTotal = async () => {
    const params = {};
    if (search.value) params.search = search.value;
    if (selectedBaseModel.value) params.baseModel = selectedBaseModel.value;
    if (selectedModelType.value) params.modelType = selectedModelType.value;
    if (nsfwFilter.value) params.nsfwFilter = nsfwFilter.value;
    const tagParts = [];
    if (selectedCategory.value) tagParts.push(selectedCategory.value);
    if (tagsSearch.value.trim()) tagParts.push(tagsSearch.value);
    if (tagParts.length) params.tags = tagParts.join(",");
    const res = await axios.get("/api/models/count", { params });
    total.value = res.data.count || 0;
  };

  const fetchBaseModels = async () => {
    try {
      const res = await axios.get("/api/base-models");
      baseModels.value = Array.isArray(res.data) ? res.data : [];
    } catch {
      baseModels.value = [];
    }
  };

  const debouncedUpdate = debounce(async () => {
    page.value = 1;
    await Promise.all([fetchTotal(), fetchModels()]);
  }, 300);

  const init = async () => {
    const saved = JSON.parse(localStorage.getItem(localStorageKey) || "{}");
    if (saved.search !== undefined) search.value = saved.search;
    if (saved.tagsSearch !== undefined) tagsSearch.value = saved.tagsSearch;
    if (saved.selectedCategory !== undefined)
      selectedCategory.value = saved.selectedCategory;
    if (saved.selectedBaseModel !== undefined)
      selectedBaseModel.value = saved.selectedBaseModel;
    if (saved.selectedModelType !== undefined)
      selectedModelType.value = saved.selectedModelType;
    if (saved.nsfwFilter !== undefined) nsfwFilter.value = saved.nsfwFilter;
    else if (saved.hideNsfw !== undefined)
      nsfwFilter.value = saved.hideNsfw ? "no" : "both";
    if (saved.page !== undefined) page.value = saved.page;

    await fetchBaseModels();
    await Promise.all([fetchTotal(), fetchModels()]);
    initialized.value = true;
  };

  watch([search, tagsSearch, selectedCategory, selectedBaseModel, selectedModelType, nsfwFilter], () => {
    if (initialized.value) {
      debouncedUpdate();
      debouncedSave();
    }
  });

  watch(page, () => {
    if (initialized.value) {
      fetchModels(); // Just fetch models on page change, total shouldn't change
      debouncedSave();
    }
  });

  const clearFilters = () => {
    search.value = "";
    tagsSearch.value = "";
    selectedCategory.value = "";
    selectedBaseModel.value = "";
    selectedModelType.value = "";
    nsfwFilter.value = "no";
    page.value = 1;
  };

  const totalPages = computed(() => Math.ceil(total.value / limit));

  return {
    models,
    search,
    tagsSearch,
    selectedCategory,
    selectedBaseModel,
    selectedModelType,
    nsfwFilter,
    page,
    total,
    totalPages,
    baseModels,
    modelTypes,
    categories,
    init,
    clearFilters,
    fetchModels, // Exporting if needed for manual refresh
  };
}
