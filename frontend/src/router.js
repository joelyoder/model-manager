import { createRouter, createWebHashHistory } from "vue-router";
import ModelList from "./components/ModelList.vue";
import ModelDetail from "./components/ModelDetail.vue";
import AppSettings from "./components/AppSettings.vue";
import Utilities from "./components/UtilitiesPage.vue";
import CollectionsPage from "./components/CollectionsPage.vue";
import CollectionDetail from "./components/CollectionDetail.vue";

const routes = [
  { path: "/", component: ModelList },
  { path: "/settings", component: AppSettings },
  { path: "/utilities", component: Utilities },
  { path: "/collections", component: CollectionsPage },
  { path: "/collections/:id", name: "CollectionDetail", component: CollectionDetail },
  {
    path: "/model/:modelId/version/:versionId",
    name: "ModelDetail",
    component: ModelDetail,
    props: true,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
