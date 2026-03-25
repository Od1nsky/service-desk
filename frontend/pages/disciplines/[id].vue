<template>
  <div>
    <!-- Header -->
    <div class="page-header">
      <div class="back-row">
        <n-button text @click="$router.push('/disciplines')">
          <template #icon><n-icon><ArrowBackOutline /></n-icon></template>
          Назад
        </n-button>
      </div>
      <h1 v-if="discipline">{{ discipline.name }}</h1>
      <p v-if="discipline?.description" class="disc-desc">{{ discipline.description }}</p>
    </div>

    <!-- Filters -->
    <div class="card filters-bar">
      <n-input
        v-model:value="filterTaskName"
        placeholder="Поиск по названию работы"
        clearable
        style="width: 260px"
        @update:value="onSearch"
      >
        <template #prefix><n-icon><SearchOutline /></n-icon></template>
      </n-input>
      <n-select
        v-model:value="filterStudentId"
        placeholder="Все студенты"
        clearable
        filterable
        :options="studentOptions"
        style="width: 220px"
        @update:value="load"
      />
      <n-select
        v-model:value="filterStatus"
        placeholder="Все статусы"
        clearable
        :options="statusOptions"
        style="width: 180px"
        @update:value="load"
      />
    </div>

    <!-- Stats row -->
    <div v-if="!loading && grades.length" class="stats-row">
      <div class="stat-chip">
        <span class="stat-chip__val">{{ grades.length }}</span>
        <span class="stat-chip__label">Всего работ</span>
      </div>
      <div class="stat-chip">
        <span class="stat-chip__val">{{ gradedCount }}</span>
        <span class="stat-chip__label">Оценено</span>
      </div>
      <div class="stat-chip">
        <span class="stat-chip__val">{{ avgScore ?? '—' }}</span>
        <span class="stat-chip__label">Средний балл</span>
      </div>
    </div>

    <!-- Table -->
    <div v-if="loading" class="center"><n-spin /></div>
    <div v-else-if="!grades.length" class="empty">
      <n-empty description="Нет оценок по данной дисциплине" />
    </div>
    <div v-else class="card">
      <n-data-table
        :columns="columns"
        :data="grades"
        :pagination="{ pageSize: 20 }"
        :row-props="rowProps"
        size="small"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  NButton, NIcon, NInput, NSelect, NSpin, NEmpty, NDataTable, NTag,
} from 'naive-ui'
import { ArrowBackOutline, SearchOutline } from '@vicons/ionicons5'
import { h } from 'vue'
import { disciplinesApi } from '~/api/disciplines'
import { gradesApi } from '~/api/grades'
import { usersApi } from '~/api/users'
import type { Discipline, Grade, User } from '~/types'

definePageMeta({ middleware: ['auth', 'teacher'] })

const route = useRoute()
const router = useRouter()
const disciplineId = Number(route.params.id)

const discipline = ref<Discipline | null>(null)
const grades = ref<Grade[]>([])
const students = ref<User[]>([])
const loading = ref(false)
const filterStudentId = ref<string | null>(null)
const filterStatus = ref<string | null>(null)
const filterTaskName = ref('')

let searchTimer: ReturnType<typeof setTimeout> | null = null
function onSearch() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(load, 350)
}

onMounted(async () => {
  const [d, u] = await Promise.all([
    disciplinesApi.getById(disciplineId),
    usersApi.list(),
  ])
  discipline.value = d.data
  students.value = (u.data as User[]).filter((u) => u.role === 'student')
  await load()
})

async function load() {
  loading.value = true
  try {
    const { data } = await gradesApi.list({
      discipline_id: disciplineId,
      student_id: filterStudentId.value || undefined,
      status: filterStatus.value || undefined,
      task_name: filterTaskName.value || undefined,
      limit: 200,
    })
    grades.value = data.data
  } finally {
    loading.value = false
  }
}

const gradedCount = computed(() =>
  grades.value.filter((g) => g.score != null).length
)

const avgScore = computed(() => {
  const scored = grades.value.filter((g) => g.score != null)
  if (!scored.length) return null
  return Math.round(scored.reduce((s, g) => s + g.score!, 0) / scored.length)
})

const studentOptions = computed(() =>
  students.value.map((s) => ({ label: s.full_name, value: s.id }))
)

const statusOptions = [
  { label: 'Ожидает', value: 'pending' },
  { label: 'Оценено', value: 'graded' },
  { label: 'Не аттестован', value: 'not_certified' },
  { label: 'Аттестован', value: 'certified' },
  { label: 'Завершено', value: 'closed' },
]

const statusColors: Record<string, any> = {
  pending: 'info', graded: 'warning', not_certified: 'error', certified: 'success', closed: 'default',
}
const statusLabels: Record<string, string> = {
  pending: 'Ожидает', graded: 'Оценено', not_certified: 'Не аттестован', certified: 'Аттестован', closed: 'Завершено',
}

const columns = [
  { title: '#', key: 'number', width: 60 },
  {
    title: 'Студент', key: 'student', minWidth: 160,
    render: (r: Grade) => r.student?.full_name || r.student_id,
  },
  { title: 'Работа', key: 'task_name', minWidth: 180 },
  {
    title: 'Тип', key: 'grade_type', width: 120,
    render: (r: Grade) => r.grade_type?.name || '—',
  },
  {
    title: 'Балл', key: 'score', width: 80,
    render: (r: Grade) =>
      r.score != null
        ? h('strong', { style: 'color:#0d6efd' }, String(r.score))
        : '—',
  },
  {
    title: 'Статус', key: 'status', width: 150,
    render: (r: Grade) =>
      h(NTag, { type: statusColors[r.status], size: 'small' }, { default: () => statusLabels[r.status] }),
  },
]

function rowProps(row: Grade) {
  return {
    style: 'cursor:pointer',
    onClick: () => router.push(`/grades/${row.id}`),
  }
}
</script>

<style lang="scss" scoped>
.back-row {
  margin-bottom: $sp-2;
}
.disc-desc {
  color: $text-muted;
  font-size: 14px;
  margin-top: $sp-1;
}
.filters-bar {
  @include flex(row, center, flex-start, $sp-3);
  flex-wrap: wrap;
  margin-bottom: $sp-4;
}
.stats-row {
  @include flex(row, center, flex-start, $sp-3);
  flex-wrap: wrap;
  margin-bottom: $sp-4;
}
.stat-chip {
  @include flex(column, flex-start, flex-start, 2px);
  background: $surface;
  border: 1px solid $border;
  border-radius: $radius;
  padding: $sp-3 $sp-5;
  min-width: 100px;
  &__val {
    font-size: 22px;
    font-weight: 700;
    color: $primary;
  }
  &__label {
    font-size: 12px;
    color: $text-muted;
  }
}
.center {
  @include flex(row, center, center);
  padding: 60px;
}
.empty {
  @include flex(column, center, center);
  padding: 60px;
}
</style>
