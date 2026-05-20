<template>
  <v-container class="fill-height" max-width="420">
    <v-card class="pa-6 w-100" elevation="4">
      <v-card-title class="text-h5 text-center mb-4">
        {{ isRegister ? 'Регистрация' : 'Вход' }}
      </v-card-title>

      <v-form @submit.prevent="handleSubmit">
        <v-text-field
          v-model="login"
          label="Логин"
          prepend-inner-icon="mdi-account"
          variant="outlined"
          class="mb-3"
        />
        <v-text-field
          v-model="password"
          label="Пароль"
          type="password"
          prepend-inner-icon="mdi-lock"
          variant="outlined"
          class="mb-4"
        />

        <v-alert v-if="error" type="error" class="mb-4" density="compact">
          {{ error }}
        </v-alert>

        <v-btn
          type="submit"
          color="primary"
          block
          size="large"
          :loading="loading"
        >
          {{ isRegister ? 'Зарегистрироваться' : 'Войти' }}
        </v-btn>
      </v-form>

      <v-divider class="my-4" />

      <div class="text-center">
        <v-btn variant="text" @click="isRegister = !isRegister">
          {{ isRegister ? 'Уже есть аккаунт? Войти' : 'Нет аккаунта? Зарегистрироваться' }}
        </v-btn>
      </div>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login as apiLogin, register as apiRegister } from '../api'
import { authStore } from '../stores/auth'

const router = useRouter()

const isRegister = ref(false)
const login = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    if (isRegister.value) {
      await apiRegister(login.value, password.value)
      isRegister.value = false
    } else {
      const { data } = await apiLogin(login.value, password.value)
      authStore.setToken(data.token)
      router.push('/estimate')
    }
  } catch (e) {
    error.value = e.response?.data?.error || 'Произошла ошибка'
  } finally {
    loading.value = false
  }
}
</script>
