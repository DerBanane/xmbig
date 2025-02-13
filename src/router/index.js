import { createRouter, createWebHistory } from 'vue-router';
import App from '../App.vue';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: App, // Hier sollte deine Hauptkomponente sein
  },
  // Weitere Routen hier...
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;