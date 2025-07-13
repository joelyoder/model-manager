import { createRouter, createWebHashHistory } from "vue-router";
import ModelList from "./components/ModelList.vue";
import ModelDetail from "./components/ModelDetail.vue";

const routes = [
  { path: "/", component: ModelList },
  { path: "/model/:id", component: ModelDetail, props: true },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
