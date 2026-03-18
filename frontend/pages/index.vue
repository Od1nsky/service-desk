<template>
  <div>
    <div class="page-header"><h1>Дашборд</h1></div>

    <div v-if="loading" class="center"><n-spin /></div>

    <template v-else-if="stats">
      <!-- Stat cards row -->
      <div class="stats-grid">
        <div v-for="s in statCards" :key="s.label" class="stat-card">
          <div class="stat-card__value">{{ s.value }}</div>
          <div class="stat-card__label">{{ s.label }}</div>
        </div>
      </div>

      <!-- Grades by status -->
      <div class="card" style="margin-top: 24px">
        <h3 style="margin-bottom: 16px">Оценки по статусам</h3>
        <div class="status-bars">
          <div
            v-for="(count, status) in stats.grades_by_status"
            :key="status"
            class="status-bar-row"
          >
            <span class="status-badge" :class="`status-badge--${status}`">{{
              statusLabel(status)
            }}</span>
            <div class="bar-wrap">
              <div
                class="bar"
                :style="{ width: barWidth(count as number) + '%', background: statusColor(status) }"
              />
            </div>
            <span class="bar-count">{{ count }}</span>
          </div>
        </div>
      </div>

      <!-- Top teachers -->
      <div v-if="stats.top_teachers?.length" class="card" style="margin-top: 24px">
        <h3 style="margin-bottom: 16px">Топ преподавателей</h3>
        <n-data-table :columns="teacherColumns" :data="stats.top_teachers" :pagination="false" />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { NSpin, NDataTable } from 'naive-ui'
import { statsApi } from '~/api/stats'
import type { DashboardStats, TeacherStat } from '~/types'

definePageMeta({ middleware: 'auth' })

const auth = useAuthStore()
const router = useRouter()

onMounted(() => {
  // Redirect students to their grades list
  if (auth.user?.role === 'student') {
    router.replace('/grades')
    return
  }
  loadStats()
})

const loading = ref(true)
const stats = ref<DashboardStats | null>(null)

async function loadStats() {
  try {
    const { data } = await statsApi.dashboard()
    stats.value = data
  } finally {
    loading.value = false
  }
}

const statCards = computed(() => [
  { label: 'Всего оценок', value: stats.value?.total_grades ?? 0 },
  { label: 'Ожидают оценки', value: stats.value?.pending_grades ?? 0 },
  {
    label: 'Средний балл',
    value: stats.value?.avg_score != null
      ? Math.round(stats.value.avg_score * 10) / 10
      : '—',
  },
])

const teacherColumns = [
  { title: '#', key: 'index', render: (_: TeacherStat, i: number) => i + 1, width: 50 },
  { title: 'Преподаватель', key: 'full_name' },
  { title: 'Оценено работ', key: 'count' },
]

const maxCount = computed(() =>
  Math.max(...(Object.values(stats.value?.grades_by_status ?? {}) as number[]), 1)
)
const barWidth = (count: number) => Math.max(4, Math.round((count / maxCount.value) * 100))

const statusLabel = (s: string) =>
  ({
    pending: 'Ожидает',
    graded: 'Оценено',
    not_certified: 'Не аттестован',
    certified: 'Аттестован',
    closed: 'Завершено',
  }[s] || s)
const statusColor = (s: string) =>
  ({
    pending: '#0d6efd',
    graded: '#fd7e14',
    not_certified: '#dc3545',
    certified: '#198754',
    closed: '#adb5bd',
  }[s] || '#ccc')
</script>

<style lang="scss" scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: $sp-4;
  @include mobile {
    grid-template-columns: 1fr;
  }
}
.stat-card {
  @include card;
  padding: $sp-6;
  text-align: center;
  &__value {
    font-size: 32px;
    font-weight: 700;
    color: $primary;
  }
  &__label {
    color: $text-muted;
    font-size: 13px;
    margin-top: $sp-1;
  }
}
.status-bars {
  display: flex;
  flex-direction: column;
  gap: $sp-3;
}
.status-bar-row {
  @include flex(row, center, flex-start, $sp-3);
}
.bar-wrap {
  flex: 1;
  background: $bg;
  border-radius: 4px;
  height: 12px;
  overflow: hidden;
}
.bar {
  height: 100%;
  border-radius: 4px;
  transition: width 0.4s;
}
.bar-count {
  font-size: 13px;
  font-weight: 600;
  min-width: 30px;
  text-align: right;
}
.center {
  @include flex(row, center, center);
  padding: 60px;
}
</style>
