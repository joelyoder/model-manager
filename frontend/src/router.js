import { createRouter, createWebHashHistory } from "vue-router";
import ModelList from "./components/ModelList.vue";
import ModelDetail from "./components/ModelDetail.vue";
import AppSettings from "./components/AppSettings.vue";
import Utilities from "./components/UtilitiesPage.vue";

const routes = [
  { path: "/", component: ModelList },
  { path: "/settings", component: AppSettings },
  { path: "/utilities", component: Utilities },
  {
    path: "/model/:modelId/version/:versionId",
    component: ModelDetail,
    props: true,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
