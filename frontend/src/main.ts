import { createApp } from 'vue'
import App from './App.vue'

import Toast, { POSITION, type PluginOptions } from 'vue-toastification'
import 'vue-toastification/dist/index.css'

import router from './router'
import './assets/main.css'

const app = createApp(App)

app.use(router)

const options: PluginOptions = {
  position: POSITION.TOP_CENTER,
  maxToasts: 1,
  newestOnTop: true,
  timeout: 3000,
  hideProgressBar: false,

  filterBeforeCreate: (toast, toasts) => {
    if (toast.type === 'error' && toasts.some((t) => t.type === 'error')) {
      return false
    }
    return toast
  },
}

app.use(Toast, options)

app.mount('#app')
