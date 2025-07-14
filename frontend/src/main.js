import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap";
import "quill/dist/quill.snow.css";
import "./index.css";
import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";

createApp(App).use(router).mount("#app");
