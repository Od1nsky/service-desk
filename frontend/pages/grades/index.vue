<template>
  <div>
    <div class="page-header">
      <h1>{{ auth.isTeacher ? 'Все оценки' : 'Мои оценки' }}</h1>
      <n-button v-if="auth.isTeacher" type="primary" @click="$router.push('/grades/create')">
        <template #icon><n-icon><AddOutline /></n-icon></template>
        Выставить оценку
      </n-button>
    </div>

    <n-tabs v-model:value="activeTab" type="line" animated>
      <!-- ===== TAB 1: По работам ===== -->
      <n-tab-pane name="works" tab="По работам">
        <div class="card filters-bar">
          <n-input v-model:value="filters.task_name" placeholder="Поиск по названию работы"
            clearable style="width:240px" @update:value="onSearchWorks">
            <template #prefix><n-icon><SearchOutline /></n-icon></template>
          </n-input>
          <n-select v-model:value="filters.status" placeholder="Все статусы" clearable
            :options="statusOptions" style="width:180px" @update:value="applyFilters" />
          <n-select v-model:value="filters.grade_type_id" placeholder="Тип оценки" clearable
            :options="gradeTypeOptions" style="width:160px" @update:value="applyFilters" />
          <n-select v-model:value="filters.discipline_id" placeholder="Дисциплина" clearable
            :options="disciplineOptions" style="width:180px" @update:value="applyFilters" />
        </div>

        <div v-if="store.loading" class="center"><n-spin /></div>
        <div v-else-if="!store.list.length" class="empty">
          <n-empty description="Оценок нет" />
        </div>
        <template v-else>
          <div class="grade-list">
            <GradeCard v-for="g in store.list" :key="g.id" :grade="g"
              @click="$router.push(`/grades/${g.id}`)" />
          </div>
          <div class="pagination">
            <n-pagination v-model:page="store.page" :item-count="store.total"
              :page-size="store.limit" @update:page="store.setPage" />
          </div>
        </template>
      </n-tab-pane>

      <!-- ===== TAB 2: По дисциплинам ===== -->
      <n-tab-pane name="disciplines" tab="По дисциплинам">
        <div class="card filters-bar">
          <n-select v-model:value="discFilter.discipline_id" placeholder="Дисциплина"
            clearable :options="disciplineOptions" style="width:220px" />
          <n-select v-if="auth.isTeacher" v-model:value="discFilter.student_id" placeholder="Студент"
            clearable filterable :options="studentOptions" style="width:220px" />
        </div>

        <div v-if="discDetailLoading" class="center"><n-spin /></div>

        <template v-else-if="discFilter.discipline_id && discFilter.student_id">
          <!-- Works list for this student + discipline -->
          <div class="card" style="margin-bottom:16px">
            <h3 style="margin-bottom:12px">Работы по дисциплине</h3>
            <div v-if="!discWorks.length" class="empty-sm">Нет оценок по данной дисциплине</div>
            <n-data-table v-else :columns="worksColumns" :data="discWorks" :pagination="false" size="small" />
          </div>

          <!-- Discipline grade summary -->
          <div class="card disc-summary">
            <div class="summary-row">
              <span class="summary-label">Сумма баллов за работы</span>
              <span class="summary-val">{{ discSummaryRow?.works_total ?? 0 }}</span>
            </div>
            <div class="summary-row">
              <span class="summary-label">Дополнительный балл</span>
              <div class="summary-input">
                <n-input-number v-model:value="extraScoreInput" :min="0" :max="100"
                  size="small" style="width:100px" placeholder="0" />
                <n-button type="primary" size="small" :loading="savingExtra" @click="saveExtra">
                  Сохранить
                </n-button>
              </div>
            </div>
            <div class="summary-row summary-row--total">
              <span class="summary-label">Итоговая оценка</span>
              <span class="summary-total">{{ computedTotal }}</span>
            </div>
            <div class="summary-row">
              <span class="summary-label">Статус по дисциплине</span>
              <div v-if="auth.isTeacher" class="summary-input">
                <n-select v-model:value="discStatusInput" :options="discStatusOptions"
                  size="small" style="width:180px" />
                <n-button type="primary" size="small" :loading="savingStatus" @click="saveStatus">
                  Сохранить
                </n-button>
              </div>
              <n-tag v-else :type="discStatusColors[discSummaryRow?.status ?? 'pending']" size="small">
                {{ discStatusLabels[discSummaryRow?.status ?? 'pending'] }}
              </n-tag>
            </div>
          </div>
        </template>

        <!-- Student-only view: all disciplines at once -->
        <template v-else-if="auth.isStudent">
          <div v-if="discSummaryLoading" class="center"><n-spin /></div>
          <div v-else class="card">
            <n-data-table :columns="summaryColumns" :data="allDiscSummary" :pagination="false" />
          </div>
        </template>

        <div v-else class="empty">
          <n-empty description="Выберите дисциплину и студента" />
        </div>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import {
  NButton, NIcon, NInput, NSelect, NSpin, NEmpty, NPagination,
  NDataTable, NInputNumber, NTabs, NTabPane, NTag, useMessage,
} from 'naive-ui'
import { AddOutline, SearchOutline } from '@vicons/ionicons5'
import { h } from 'vue'
import { disciplinesApi } from '~/api/disciplines'
import { gradeTypesApi } from '~/api/grade-types'
import { disciplineGradesApi } from '~/api/discipline-grades'
import { gradesApi } from '~/api/grades'
import { usersApi } from '~/api/users'
import type { Discipline, GradeType, DisciplineGrade, Grade, User } from '~/types'

definePageMeta({ middleware: 'auth' })

const store = useGradesStore()
const auth = useAuthStore()
const message = useMessage()

const activeTab = ref('works')

// Shared data
const disciplines = ref<Discipline[]>([])
const gradeTypes = ref<GradeType[]>([])
const students = ref<User[]>([])

// Tab 1 filters
const filters = reactive({
  task_name: '',
  status: null as string | null,
  grade_type_id: null as number | null,
  discipline_id: null as number | null,
})

let searchWorksTimer: ReturnType<typeof setTimeout> | null = null
function onSearchWorks() {
  if (searchWorksTimer) clearTimeout(searchWorksTimer)
  searchWorksTimer = setTimeout(applyFilters, 350)
}

// Tab 2
const discFilter = reactive({
  discipline_id: null as number | null,
  student_id: null as string | null,
})
const discWorks = ref<Grade[]>([])
const discSummaryRow = ref<DisciplineGrade | null>(null)
const discDetailLoading = ref(false)
const extraScoreInput = ref<number | null>(null)
const savingExtra = ref(false)
const discStatusInput = ref<string>('pending')
const savingStatus = ref(false)

// Student summary (all disciplines)
const allDiscSummary = ref<DisciplineGrade[]>([])
const discSummaryLoading = ref(false)

onMounted(async () => {
  store.fetchList()
  const [d, g, u] = await Promise.all([
    disciplinesApi.list(),
    gradeTypesApi.list(),
    auth.isTeacher ? usersApi.list() : Promise.resolve({ data: [] }),
  ])
  disciplines.value = d.data
  gradeTypes.value = g.data
  students.value = (u.data as User[]).filter((u) => u.role === 'student')

  if (auth.isStudent) loadAllDiscSummary()
})

async function loadAllDiscSummary() {
  discSummaryLoading.value = true
  try {
    const { data } = await disciplineGradesApi.list()
    allDiscSummary.value = data
  } finally {
    discSummaryLoading.value = false
  }
}

async function loadDiscDetail() {
  if (!discFilter.discipline_id) return
  const studentId = auth.isStudent ? auth.user!.id : discFilter.student_id
  if (!studentId) return

  discDetailLoading.value = true
  try {
    // Works for this discipline
    const { data: worksData } = await gradesApi.list({
      discipline_id: discFilter.discipline_id,
      student_id: studentId,
      limit: 100,
    })
    discWorks.value = worksData.data

    // Discipline grade summary
    const { data: summaryData } = await disciplineGradesApi.list(studentId)
    const row = summaryData.find((r) => r.discipline_id === discFilter.discipline_id) ?? null
    discSummaryRow.value = row
    extraScoreInput.value = row?.extra_score ?? null
    discStatusInput.value = row?.status ?? 'pending'
  } finally {
    discDetailLoading.value = false
  }
}

// Auto-load when either filter changes
watch(
  () => [discFilter.discipline_id, discFilter.student_id],
  () => { loadDiscDetail() }
)

const computedTotal = computed(() => {
  const works = discSummaryRow.value?.works_total ?? 0
  const extra = extraScoreInput.value ?? 0
  return Math.min(works + extra, 100)
})

async function saveExtra() {
  if (extraScoreInput.value === null) return
  const studentId = auth.isStudent ? auth.user!.id : discFilter.student_id
  if (!studentId || !discFilter.discipline_id) return
  savingExtra.value = true
  try {
    await disciplineGradesApi.upsert({
      student_id: studentId,
      discipline_id: discFilter.discipline_id,
      extra_score: extraScoreInput.value,
      notes: '',
      status: discStatusInput.value,
    })
    await loadDiscDetail()
    message.success('Доп. балл сохранён')
  } catch {
    message.error('Ошибка сохранения')
  } finally {
    savingExtra.value = false
  }
}

async function saveStatus() {
  const studentId = auth.isStudent ? auth.user!.id : discFilter.student_id
  if (!studentId || !discFilter.discipline_id) return
  savingStatus.value = true
  try {
    await disciplineGradesApi.upsert({
      student_id: studentId,
      discipline_id: discFilter.discipline_id,
      extra_score: extraScoreInput.value ?? 0,
      notes: '',
      status: discStatusInput.value,
    })
    await loadDiscDetail()
    message.success('Статус сохранён')
  } catch {
    message.error('Ошибка сохранения')
  } finally {
    savingStatus.value = false
  }
}

function applyFilters() {
  store.setFilters({
    task_name: filters.task_name || '',
    status: filters.status || '',
    grade_type_id: filters.grade_type_id,
    discipline_id: filters.discipline_id,
  })
}

// Options
const statusOptions = [
  { label: 'Ожидает', value: 'pending' },
  { label: 'Оценено', value: 'graded' },
  { label: 'Не аттестован', value: 'not_certified' },
  { label: 'Аттестован', value: 'certified' },
  { label: 'Завершено', value: 'closed' },
]
const gradeTypeOptions = computed(() => gradeTypes.value.map((g) => ({ label: g.name, value: g.id })))
const disciplineOptions = computed(() => disciplines.value.map((d) => ({ label: d.name, value: d.id })))
const studentOptions = computed(() => students.value.map((s) => ({ label: `${s.full_name}`, value: s.id })))

const discStatusOptions = [
  { label: 'Ожидает', value: 'pending' },
  { label: 'Аттестован', value: 'certified' },
  { label: 'Не аттестован', value: 'not_certified' },
]
const discStatusColors: Record<string, any> = {
  pending: 'default', certified: 'success', not_certified: 'error',
}
const discStatusLabels: Record<string, string> = {
  pending: 'Ожидает', certified: 'Аттестован', not_certified: 'Не аттестован',
}

const statusColors: Record<string, any> = {
  pending: 'info', graded: 'warning', not_certified: 'error', certified: 'success', closed: 'default',
}
const statusLabels: Record<string, string> = {
  pending: 'Ожидает', graded: 'Оценено', not_certified: 'Не аттестован', certified: 'Аттестован', closed: 'Завершено',
}

const workScores = ref<Record<string, number | null>>({})
const workSaving = ref<Record<string, boolean>>({})

watch(discWorks, (works) => {
  works.forEach((w) => {
    if (!(w.id in workScores.value)) {
      workScores.value[w.id] = w.score ?? null
    }
  })
}, { immediate: true })

async function saveWorkScore(row: Grade) {
  const score = workScores.value[row.id]
  if (score === null || score === undefined) return
  workSaving.value[row.id] = true
  try {
    await gradesApi.update(row.id, { score })
    await loadDiscDetail()
    message.success('Балл сохранён')
  } catch {
    message.error('Ошибка сохранения')
  } finally {
    workSaving.value[row.id] = false
  }
}

const worksColumns = computed(() => {
  const cols: any[] = [
    { title: '#', key: 'number', width: 60 },
    { title: 'Работа', key: 'task_name' },
    { title: 'Тип', key: 'grade_type', width: 110, render: (r: Grade) => r.grade_type?.name || '—' },
    {
      title: 'Статус', key: 'status', width: 140,
      render: (r: Grade) => h(NTag, { type: statusColors[r.status], size: 'small' }, { default: () => statusLabels[r.status] }),
    },
  ]

  if (auth.isTeacher) {
    cols.push({
      title: 'Балл (макс 100)', key: 'score', width: 220,
      render: (r: Grade) => h('div', { style: 'display:flex;gap:6px;align-items:center' }, [
        h(NInputNumber, {
          value: workScores.value[r.id] ?? null,
          min: 0,
          max: 100,
          size: 'small',
          style: 'width:90px',
          'onUpdate:value': (v: number | null) => { workScores.value[r.id] = v },
        }),
        h(NButton, {
          size: 'small',
          type: 'primary',
          loading: !!workSaving.value[r.id],
          disabled: workScores.value[r.id] === null,
          onClick: () => saveWorkScore(r),
        }, { default: () => 'Сохранить' }),
      ]),
    })
  } else {
    cols.push({
      title: 'Балл', key: 'score', width: 80,
      render: (r: Grade) => r.score != null ? h('strong', { style: 'color:#0d6efd' }, String(r.score)) : '—',
    })
  }

  return cols
})

const summaryColumns = [
  { title: 'Дисциплина', key: 'discipline', render: (r: DisciplineGrade) => r.discipline?.name || '—' },
  { title: 'Баллы за работы', key: 'works_total', width: 160, render: (r: DisciplineGrade) => r.works_total ?? 0 },
  { title: 'Доп. балл', key: 'extra_score', width: 120, render: (r: DisciplineGrade) => r.extra_score ?? '—' },
  {
    title: 'Итого', key: 'total', width: 100,
    render: (r: DisciplineGrade) => h('strong', { style: 'color:#0d6efd' }, String(r.total ?? r.works_total ?? 0)),
  },
  {
    title: 'Статус', key: 'status', width: 150,
    render: (r: DisciplineGrade) => h(NTag, { type: discStatusColors[r.status ?? 'pending'], size: 'small' }, { default: () => discStatusLabels[r.status ?? 'pending'] }),
  },
]
</script>

<style lang="scss" scoped>
.filters-bar {
  @include flex(row, center, flex-start, $sp-3);
  flex-wrap: wrap;
  margin-bottom: $sp-4;
}
.grade-list {
  display: flex;
  flex-direction: column;
  gap: $sp-3;
}
.pagination {
  @include flex(row, center, center);
  margin-top: $sp-6;
}
.empty {
  @include flex(column, center, center);
  padding: 60px;
}
.empty-sm {
  color: $text-muted;
  font-size: 13px;
  padding: $sp-4 0;
}
.center {
  @include flex(row, center, center);
  padding: 60px;
}
.disc-summary {
  display: flex;
  flex-direction: column;
  gap: $sp-1;
}
.summary-row {
  @include flex(row, center, space-between, $sp-3);
  padding: $sp-3 0;
  border-bottom: 1px solid $bg;
  font-size: 14px;
  &:last-child { border-bottom: none; }
  &--total { margin-top: $sp-2; }
}
.summary-label { color: $text-muted; }
.summary-input { @include flex(row, center, flex-end, $sp-2); }
.summary-val { font-weight: 600; }
.summary-total {
  font-size: 22px;
  font-weight: 700;
  color: $primary;
}
</style>
