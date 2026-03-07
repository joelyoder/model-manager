import { createRouter, createWebHashHistory } from "vue-router";

const routes = [
  { path: "/", component: () => import("./components/ModelList.vue") },
  { path: "/settings", component: () => import("./components/AppSettings.vue") },
  { path: "/utilities", component: () => import("./components/UtilitiesPage.vue") },
  { path: "/collections", component: () => import("./components/CollectionsPage.vue") },
  { path: "/collections/:id", name: "CollectionDetail", component: () => import("./components/CollectionDetail.vue") },
  {
    path: "/model/:modelId/version/:versionId",
    name: "ModelDetail",
    component: () => import("./components/ModelDetail.vue"),
    props: true,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
