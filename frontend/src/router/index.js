import { createRouter, createWebHistory } from 'vue-router'
import { authStore } from '../stores/auth'

import LoginView from '../views/LoginView.vue'
import EstimateView from '../views/EstimateView.vue'
import SemestersView from '../views/SemestersView.vue'

const routes = [
  { path: '/', redirect: '/estimate' },
  { path: '/login', component: LoginView },
  { path: '/estimate', component: EstimateView, meta: { requiresAuth: true } },
  { path: '/semesters', component: SemestersView, meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Защита маршрутов — незалогиненных отправляем на /login
router.beforeEach((to) => {
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    return '/login'
  }
})

export default router
