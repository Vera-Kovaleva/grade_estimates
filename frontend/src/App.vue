<template>
  <v-app>
    <!-- Навигационная панель (только для залогиненных) -->
    <v-app-bar v-if="auth.isLoggedIn" color="primary" flat>
      <v-app-bar-title>Оценщик успеваемости</v-app-bar-title>
      <v-btn to="/estimate" variant="text">Прогноз</v-btn>
      <v-btn to="/semesters" variant="text">Журнал семестров</v-btn>
      <v-btn variant="text" @click="handleLogout">Выйти</v-btn>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { authStore as auth } from './stores/auth'

const router = useRouter()

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>
