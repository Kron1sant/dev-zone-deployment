import Vue from 'vue'
import App from './App.vue'
import vuetify from '@/plugins/vuetify'
import router from '@/router'
import store from '@/store'
import axios from "./axios-preset"

axios.defaults.baseURL = 'http://localhost:8082' // ToDo

Vue.config.productionTip = false
new Vue({
  vuetify,
  store,
  router,
  render: h => h(App),
}).$mount('#app')
