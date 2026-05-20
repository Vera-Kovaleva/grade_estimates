import { reactive } from 'vue'

// Простое реактивное хранилище авторизации без Pinia/Vuex
export const authStore = reactive({
  token: localStorage.getItem('token') || null,

  setToken(token) {
    this.token = token
    localStorage.setItem('token', token)
  },

  logout() {
    this.token = null
    localStorage.removeItem('token')
  },

  get isLoggedIn() {
    return !!this.token
  },
})
