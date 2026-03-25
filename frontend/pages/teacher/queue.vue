<template>
  <div>
    <div class="page-header"><h1>Очередь на оценивание</h1></div>

    <div class="card filters-bar">
      <n-select
        v-model:value="filters.status"
        placeholder="Все статусы"
        clearable
        :options="statusOptions"
        style="width: 180px"
        @update:value="applyFilters"
      />
      <n-select
        v-model:value="filters.grade_type_id"
        placeholder="Тип оценки"
        clearable
        :options="gradeTypeOptions"
        style="width: 160px"
        @update:value="applyFilters"
      />
      <n-select
        v-model:value="filters.discipline_id"
        placeholder="Дисциплина"
        clearable
        :options="disciplineOptions"
        style="width: 180px"
        @update:value="applyFilters"
      />
    </div>

    <div v-if="store.loading" class="center"><n-spin /></div>
    <div v-else-if="!store.list.length" class="center">
      <n-empty description="Оценок нет" />
    </div>
    <template v-else>
      <n-data-table
        :columns="columns"
        :data="store.list"
        :pagination="false"
        style="cursor: pointer"
        @row-props="rowProps"
      />
      <div class="pagination">
        <n-pagination
          v-model:page="store.page"
          :item-count="store.total"
          :page-size="store.limit"
          @update:page="store.setPage"
        />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { NSelect, NSpin, NEmpty, NDataTable, NPagination, NTag, NButton } from 'naive-ui'
import type { TagProps } from 'naive-ui'
import { h } from 'vue'
import { disciplinesApi } from '~/api/disciplines'
import { gradeTypesApi } from '~/api/grade-types'
import type { Discipline, GradeType, Grade } from '~/types'

definePageMeta({ middleware: ['auth', 'teacher'] })

const store = useGradesStore()
const router = useRouter()
const filters = reactive({
  status: null as string | null,
  grade_type_id: null as number | null,
  discipline_id: null as number | null,
})
const disciplines = ref<Discipline[]>([])
const gradeTypes = ref<GradeType[]>([])

onMounted(async () => {
  store.fetchList()
  const [d, g] = await Promise.all([disciplinesApi.list(), gradeTypesApi.list()])
  disciplines.value = d.data
  gradeTypes.value = g.data
})

const statusLabels: Record<string, string> = {
  pending: 'Ожидает',
  graded: 'Оценено',
  not_certified: 'Не аттестован',
  certified: 'Аттестован',
  closed: 'Завершено',
}
const statusColors: Record<string, TagProps['type']> = {
  pending: 'info',
  graded: 'warning',
  not_certified: 'error',
  certified: 'success',
  closed: 'default',
}

const columns = [
  { title: '#', key: 'number', width: 70 },
  { title: 'Название работы', key: 'task_name', ellipsis: { tooltip: true } },
  {
    title: 'Статус',
    key: 'status',
    width: 160,
    render: (row: Grade) =>
      h(NTag, { type: statusColors[row.status], size: 'small' }, { default: () => statusLabels[row.status] }),
  },
  {
    title: 'Тип оценки',
    key: 'grade_type',
    width: 130,
    render: (row: Grade) => row.grade_type?.name || '—',
  },
  {
    title: 'Дисциплина',
    key: 'discipline',
    width: 150,
    render: (row: Grade) => row.discipline?.name || '—',
  },
  {
    title: 'Студент',
    key: 'student',
    width: 160,
    render: (row: Grade) => row.student?.full_name || '—',
  },
  {
    title: 'Преподаватель',
    key: 'teacher',
    width: 160,
    render: (row: Grade) => row.teacher?.full_name || '—',
  },
  {
    title: 'Дата',
    key: 'created_at',
    width: 130,
    render: (row: Grade) => new Date(row.created_at).toLocaleDateString('ru-RU'),
  },
  {
    title: '',
    key: 'actions',
    width: 100,
    render: (row: Grade) =>
      h(NButton, { size: 'small', onClick: (e: MouseEvent) => { e.stopPropagation(); router.push(`/grades/${row.id}`) } }, { default: () => 'Оценить' }),
  },
]

const statusOptions = Object.entries(statusLabels).map(([v, l]) => ({ value: v, label: l }))
const gradeTypeOptions = computed(() => gradeTypes.value.map((g) => ({ value: g.id, label: g.name })))
const disciplineOptions = computed(() => disciplines.value.map((d) => ({ value: d.id, label: d.name })))

function applyFilters() {
  store.setFilters({
    status: filters.status || '',
    grade_type_id: filters.grade_type_id,
    discipline_id: filters.discipline_id,
  })
}

function rowProps(row: any) {
  return {
    onClick: () => router.push(`/grades/${row.id}`),
  }
}
</script>

<style lang="scss" scoped>
.filters-bar {
  @include flex(row, center, flex-start, $sp-3);
  flex-wrap: wrap;
  margin-bottom: $sp-4;
}
.pagination {
  @include flex(row, center, center);
  margin-top: $sp-5;
}
.center {
  @include flex(row, center, center);
  padding: 60px;
}
</style>
