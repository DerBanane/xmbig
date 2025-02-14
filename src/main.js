import { createApp } from 'vue';
import App from './App.vue';
import router from './router'; // Importiere den Router
import store from './store'; // Importiere den Store
import $ from 'jquery';

// Stelle jQuery global zur Verfügung (optional, nur wenn nötig)
window.$ = window.jQuery = $;

// Importiere Bootstrap JavaScript (optional, je nach Bedarf)
import 'bootstrap';

createApp(App)
  .use(router) // Verwende den Router
  .use(store) // Verwende den Store
  .mount('#app');