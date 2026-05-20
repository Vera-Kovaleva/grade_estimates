<template>
  <v-container max-width="800">
    <v-row>
      <!-- Форма расчёта -->
      <v-col cols="12">
        <v-card class="pa-6" elevation="2">
          <v-card-title class="text-h6 mb-4">Прогноз оценки</v-card-title>
          <v-form @submit.prevent="handleEstimate">
            <v-row>

              <!-- 1. Учебное время (S) — select 0..4 -->
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.study_hours"
                  :items="studyHoursOptions"
                  item-title="label"
                  item-value="value"
                  label="Сколько часов в день ты учишься?"
                  variant="outlined"
                />
              </v-col>

              <!-- 2. Прогулы (A) — select 0..4 -->
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.absences"
                  :items="absencesOptions"
                  item-title="label"
                  item-value="value"
                  label="Как часто ты пропускаешь занятия?"
                  variant="outlined"
                />
              </v-col>

              <!-- 3. Репетитор (T) — да/нет -->
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.tutor"
                  :items="yesNoOptions"
                  item-title="label"
                  item-value="value"
                  label="Занимаешься ли ты с репетитором?"
                  variant="outlined"
                />
              </v-col>

              <!-- 4. Поддержка родителей (P) — select 0..4 -->
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.parental"
                  :items="parentalOptions"
                  item-title="label"
                  item-value="value"
                  label="Как родители поддерживают твою учёбу?"
                  variant="outlined"
                />
              </v-col>

              <!-- 5. Внеклассная деятельность (E) — да/нет -->
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.extracurricular"
                  :items="yesNoOptions"
                  item-title="label"
                  item-value="value"
                  label="Участвуешь ли ты во внеклассных мероприятиях?"
                  variant="outlined"
                />
              </v-col>

              <!-- 6. Спорт (Sp) — да/нет -->
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.sports"
                  :items="yesNoOptions"
                  item-title="label"
                  item-value="value"
                  label="Занимаешься ли ты спортом?"
                  variant="outlined"
                />
              </v-col>

              <!-- 7. Музыка (M) — да/нет -->
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.music"
                  :items="yesNoOptions"
                  item-title="label"
                  item-value="value"
                  label="Занимаешься ли ты музыкой?"
                  variant="outlined"
                />
              </v-col>

              <!-- 8. Волонтерство (V) — да/нет -->
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.volunteering"
                  :items="yesNoOptions"
                  item-title="label"
                  item-value="value"
                  label="Занимаешься ли ты волонтерством?"
                  variant="outlined"
                />
              </v-col>

              <!-- 9. Образование родителей (Ed) — select 0..4 -->
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.parent_education"
                  :items="parentEducationOptions"
                  item-title="label"
                  item-value="value"
                  label="Какое образование у твоих родителей?"
                  variant="outlined"
                />
              </v-col>

              <!-- 10. Домашняя работа (H) — select 0..4 -->
              <v-col cols="12" sm="6">
                <v-select
                  v-model="form.homework"
                  :items="homeworkOptions"
                  item-title="label"
                  item-value="value"
                  label="Как часто ты делаешь домашнее задание?"
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
              size="large"
              :loading="loading"
            >
              Рассчитать
            </v-btn>
          </v-form>

          <!-- Результат -->
          <v-alert
            v-if="result !== null"
            :type="gradeAlertType"
            class="mt-6"
            prominent
          >
            <div class="text-h5 mb-2">Прогнозируемая оценка: <strong>{{ result }}</strong></div>
            <div class="text-body-1">{{ motivationalMessage }}</div>
          </v-alert>
        </v-card>
      </v-col>

      <!-- История расчётов -->
      <v-col cols="12">
        <v-card class="pa-6" elevation="2">
          <div class="d-flex align-center justify-space-between mb-4">
            <v-card-title class="text-h6 pa-0">История расчётов</v-card-title>
            <v-btn
              v-if="predictions.length"
              color="error"
              variant="outlined"
              size="small"
              :loading="clearing"
              @click="handleClear"
            >
              Очистить всё
            </v-btn>
          </div>

          <v-data-table
            :headers="historyHeaders"
            :items="predictions"
            :loading="loadingHistory"
            no-data-text="История пуста"
            items-per-page="5"
          >
            <template #item.created_at="{ item }">
              {{ formatDate(item.created_at) }}
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { estimate, getPredictions, clearPredictions } from '../api'

// ── Варианты ответов ────────────────────────────────────────────────────────

const yesNoOptions = [
  { label: 'Нет', value: 0 },
  { label: 'Да',  value: 1 },
]

const studyHoursOptions = [
  { label: 'Меньше 1 часа',  value: 0 },
  { label: '1–2 часа',       value: 1 },
  { label: '2–3 часа',       value: 2 },
  { label: '3–4 часа',       value: 3 },
  { label: 'Более 4 часов',  value: 4 },
]

const absencesOptions = [
  { label: 'Никогда',          value: 0 },
  { label: 'Редко (1–2 раза)', value: 1 },
  { label: 'Иногда (3–5)',     value: 2 },
  { label: 'Часто (6–9)',      value: 3 },
  { label: 'Очень часто (10+)',value: 4 },
]

const parentalOptions = [
  { label: 'Никакой',           value: 0 },
  { label: 'Минимальная',       value: 1 },
  { label: 'Умеренная',         value: 2 },
  { label: 'Хорошая',           value: 3 },
  { label: 'Очень активная',    value: 4 },
]

const parentEducationOptions = [
  { label: 'Нет образования',   value: 0 },
  { label: 'Среднее',           value: 1 },
  { label: 'Среднее специальное', value: 2 },
  { label: 'Неоконченное высшее', value: 3 },
  { label: 'Высшее',            value: 4 },
]

const homeworkOptions = [
  { label: 'Никогда',          value: 0 },
  { label: 'Редко',            value: 1 },
  { label: 'Иногда',           value: 2 },
  { label: 'Часто',            value: 3 },
  { label: 'Всегда',           value: 4 },
]

// ── Состояние формы ─────────────────────────────────────────────────────────

const form = ref({
  study_hours:      null,
  absences:         null,
  tutor:            null,
  parental:         null,
  extracurricular:  null,
  sports:           null,
  music:            null,
  volunteering:     null,
  parent_education: null,
  homework:         null,
})

const result = ref(null)
const error = ref('')
const loading = ref(false)
const predictions = ref([])
const loadingHistory = ref(false)
const clearing = ref(false)

const historyHeaders = [
  { title: 'Оценка',          key: 'predicted_grade',            align: 'center' },
  { title: 'Учёба (ч/день)',  key: 'parameters.study_hours' },
  { title: 'Пропуски',        key: 'parameters.absences' },
  { title: 'Репетитор',       key: 'parameters.tutor' },
  { title: 'Дата',            key: 'created_at' },
]

// ── Подбадривания ───────────────────────────────────────────────────────────

const messages = {
  2: [
    'Каждое начало требует смелости. Ты уже сделал первый шаг.',
    'Продолжай — даже маленькие усилия приносят результат.',
    'Ты справляешься, просто дай себе время.',
    'Не бойся ошибаться. Ошибки ведут к знаниям.',
    'Ты на правильном пути, даже если пока не видно.',
    'Просто не останавливайся — дальше будет легче.',
  ],
  3: [
    'Ты движешься вперёд, и это главное.',
    'Каждый новый ответ делает тебя опытнее.',
    'Немного терпения — и ты заметишь прогресс.',
    'Даже небольшой шаг — это движение к цели.',
    'Ты способен на большее, просто продолжай.',
    'У тебя всё получится, если не сворачивать.',
  ],
  4: [
    'Отличный результат! Ты проявляешь настойчивость.',
    'Ты близок к вершине — ещё чуть-чуть.',
    'Прекрасная работа, продолжай радовать себя.',
    'Твой прогресс впечатляет. Не сбавляй оборотов.',
    'Ты знаешь и умеешь больше, чем кажется.',
    'Горжусь твоими стараниями. Так держать!',
  ],
  5: [
    'Ты показываешь блестящий результат.',
    'Великолепно! Твои усилия приносят плоды.',
    'Ты настоящий пример для подражания.',
    'Идеально! Продолжай в том же духе.',
    'Ты на высоте — так держать всегда.',
    'Превосходно! Это заслуженный успех.',
  ],
}

const motivationalMessage = computed(() => {
  if (result.value === null) return ''
  const pool = messages[result.value] ?? []
  return pool[Math.floor(Math.random() * pool.length)] ?? ''
})

const gradeAlertType = computed(() => {
  if (result.value >= 5) return 'success'
  if (result.value >= 4) return 'info'
  if (result.value >= 3) return 'warning'
  return 'error'
})



async function handleEstimate() {
  error.value = ''
  loading.value = true
  try {
    const { data } = await estimate(form.value)
    result.value = data.grade
    await loadHistory()
  } catch (e) {
    error.value = e.response?.data?.error || 'Ошибка расчёта'
  } finally {
    loading.value = false
  }
}

async function loadHistory() {
  loadingHistory.value = true
  try {
    const { data } = await getPredictions()
    predictions.value = data.predictions
  } finally {
    loadingHistory.value = false
  }
}

async function handleClear() {
  clearing.value = true
  try {
    await clearPredictions()
    predictions.value = []
    result.value = null
  } finally {
    clearing.value = false
  }
}

function formatDate(iso) {
  return new Date(iso).toLocaleString('ru-RU')
}

onMounted(loadHistory)
</script>
