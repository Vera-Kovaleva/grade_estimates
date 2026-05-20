import axios from 'axios'
import { authStore } from '../stores/auth'
import router from '../router'

// Единый экземпляр axios — все запросы идут через него
const api = axios.create({
  baseURL: '/api/v1',
})

// Перед каждым запросом подставляем токен из localStorage
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// При получении 401 (истёкший или недействительный токен) — разлогиниваем
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      authStore.logout()
      router.push('/login')
    }
    return Promise.reject(error)
  },
)

// --- Авторизация ---

export const register = (login, password) =>
  api.post('/auth/register', { login, password })

export const login = (login, password) =>
  api.post('/auth/login', { login, password })

// --- Прогноз оценки (страница 1) ---

export const estimate = (parameters) =>
  api.post('/estimate', { parameters })

export const getPredictions = () =>
  api.get('/predictions')

export const clearPredictions = () =>
  api.delete('/predictions')

// --- Журнал семестров (страница 2) ---

export const getSemesters = () =>
  api.get('/semesters')

export const addSemester = (subject, actual_grade, parameters) =>
  api.post('/semesters', { subject, actual_grade, parameters })

export const deleteSemester = (id) =>
  api.delete(`/semesters/${id}`)
