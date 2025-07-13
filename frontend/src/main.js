import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import App from './App.vue'
import ModelList from './components/ModelList.vue'
import ModelDetails from './components/ModelDetails.vue'

const routes = [
  { path: '/', component: ModelList },
  { path: '/model/:id', component: ModelDetails, props: true }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

createApp(App).use(router).mount('#app')
