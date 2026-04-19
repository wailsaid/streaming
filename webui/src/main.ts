import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import { TOKEN_KEY } from './composables/useAuth'
import './assets/main.css'

// Set axios Authorization header from stored token on startup
const token = localStorage.getItem(TOKEN_KEY)
if (token) {
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
}

// Point axios at the Go backend (8080) when running with Vite dev server
// In production, they're on the same origin, so no base URL needed.
// To enable proxying in dev, configure vite.config.ts proxy below.

const app = createApp(App)
app.use(router)
app.mount('#app')
