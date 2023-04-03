import Vue from 'vue'
import App from './App.vue'
import { BootstrapVue, IconsPlugin, FormFilePlugin } from 'bootstrap-vue'
// Import Bootstrap an BootstrapVue CSS files (order is important)
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import VueResource from 'vue-resource'
import router from "./router.js"

Vue.use(VueResource);
Vue.config.productionTip = false
Vue.use(BootstrapVue)
Vue.use(IconsPlugin)
Vue.use(FormFilePlugin)

Vue.http.interceptors.push((request, next) => {
  request.credentials = true;
  next();
});

new Vue({
  el: '#app',
  router,
  render: h => h(App),
})
