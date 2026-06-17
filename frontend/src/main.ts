import { createApp } from 'vue'
import App from './App.vue'

import Toast from 'vue-toastification'
import 'vue-toastification/dist/index.css'

import router from './router'
import './assets/main.css'

const app = createApp(App)

app.use(router)

app.use(Toast, {
  position: 'bottom-right',
  timeout: 4000,
})

app.mount('#app')
