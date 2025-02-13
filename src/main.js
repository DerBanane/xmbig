import { createApp } from 'vue';
import App from './App.vue';
import router from './router'; // Importiere den Router
import store from './store'; // Importiere den Store
import $ from 'jquery';

// Stelle jQuery global zur Verfügung (optional, nur wenn nötig)
window.$ = window.jQuery = $;

// Importiere Bootstrap JavaScript (optional, je nach Bedarf)
import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css';
import 'datatables.net-bs4';
import 'datatables.net-bs4/css/dataTables.bootstrap4.min.css';
import 'datatables.net-buttons-bs4';
import 'datatables.net-buttons-bs4/css/buttons.bootstrap4.min.css';
import 'datatables.net-select-bs4';
import 'datatables.net-select-bs4/css/select.bootstrap4.min.css';
import 'datatables.net-buttons/js/dataTables.buttons.js';
import 'datatables.net-buttons/js/buttons.colVis.js';
import 'datatables.net-buttons/js/buttons.html5.js';
import 'datatables.net-buttons/js/buttons.print.js';
import 'datatables.net-select/js/dataTables.select.js';

createApp(App)
  .use(router) // Verwende den Router
  .use(store) // Verwende den Store
  .mount('#app');