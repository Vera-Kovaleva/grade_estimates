<template>
  <v-container max-width="800">
    <v-row>
      <!-- Форма добавления -->
      <v-col cols="12">
        <v-card class="pa-6" elevation="2">
          <v-card-title class="text-h6 mb-4">Добавить запись</v-card-title>
          <v-form @submit.prevent="handleAdd">
            <v-row>
              <v-col cols="12" sm="6">
                <v-text-field
                  v-model="form.subject"
                  label="Название предмета"
                  variant="outlined"
                />
              </v-col>
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.actual_grade"
                  :items="[2, 3, 4, 5]"
                  label="Полученная оценка"
                  variant="outlined"
                />
              </v-col>
              <v-col
                v-for="param in params"
                :key="param.key"
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model.number="form.parameters[param.key]"
                  :label="param.label"
                  type="number"
                  :hint="param.hint"
                  persistent-hint
                  variant="outlined"
                />
              </v-col>
            </v-row>

            <v-alert v-if="error" type="error" class="mt-4" density="compact">
              {{ error }}
            </v-alert>

            <v-btn
              type="submit"
              color="primary"
              class="mt-4"
              :loading="adding"
            >
              Добавить
            </v-btn>
          </v-form>
        </v-card>
      </v-col>

      <!-- Список семестров -->
      <v-col cols="12">
        <v-card class="pa-6" elevation="2">
          <v-card-title class="text-h6 mb-4">Журнал семестров</v-card-title>

          <v-data-table
            :headers="headers"
            :items="semesters"
            :loading="loadingSemesters"
            no-data-text="Записей пока нет"
            items-per-page="5"
          >
            <template #item.created_at="{ item }">
              {{ formatDate(item.created_at) }}
            </template>
            <template #item.actions="{ item }">
              <v-btn
                icon="mdi-delete"
                variant="text"
                color="error"
                size="small"
                @click="handleDelete(item.id)"
              />
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getSemesters, addSemester, deleteSemester } from '../api'

const params = [
  { key: 'lecture_hours',    label: 'Время на лекциях',     hint: 'часов в неделю (0–10)' },
  { key: 'attentiveness',    label: 'Внимательность',        hint: 'балл от 0 до 10' },
  { key: 'self_study_hours', label: 'Самостоятельная подготовка', hint: 'часов в неделю (0–10)' },
  { key: 'external_sources', label: 'Сторонние источники',   hint: 'балл от 0 до 10' },
  { key: 'sleep_hours',      label: 'Количество сна',        hint: 'часов в сутки (0–10)' },
]

const emptyForm = () => ({
  subject: '',
  actual_grade: null,
  parameters: Object.fromEntries(params.map((p) => [p.key, null])),
})

const form = ref(emptyForm())
const error = ref('')
const adding = ref(false)
const semesters = ref([])
const loadingSemesters = ref(false)

const headers = [
  { title: 'Предмет', key: 'subject' },
  { title: 'Оценка', key: 'actual_grade', align: 'center' },
  { title: 'Время на лекциях', key: 'parameters.lecture_hours' },
  { title: 'Внимательность', key: 'parameters.attentiveness' },
  { title: 'Сон', key: 'parameters.sleep_hours' },
  { title: 'Дата', key: 'created_at' },
  { title: '', key: 'actions', sortable: false, align: 'end' },
]

async function handleAdd() {
  error.value = ''
  adding.value = true
  try {
    await addSemester(form.value.subject, form.value.actual_grade, form.value.parameters)
    form.value = emptyForm()
    await loadSemesters()
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка при добавлении'
  } finally {
    adding.value = false
  }
}

async function handleDelete(id) {
  try {
    await deleteSemester(id)
    semesters.value = semesters.value.filter((s) => s.id !== id)
  } catch {
    // запись не найдена или уже удалена
  }
}

async function loadSemesters() {
  loadingSemesters.value = true
  try {
    const { data } = await getSemesters()
    semesters.value = data.semesters
  } finally {
    loadingSemesters.value = false
  }
}

function formatDate(iso) {
  return new Date(iso).toLocaleString('ru-RU')
}

onMounted(loadSemesters)
</script>
