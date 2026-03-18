<template>
  <div>
    <div class="page-header">
      <h1>Мои оценки</h1>
      <n-button v-if="auth.isTeacher" type="primary" @click="$router.push('/grades/create')">
        <template #icon>
          <n-icon><AddOutline /></n-icon>
        </template>
        Выставить оценку
      </n-button>
    </div>

    <!-- Filters -->
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
    <div v-else-if="!store.list.length" class="empty">
      <n-empty description="Оценок нет" />
      <n-button v-if="auth.isTeacher" type="primary" style="margin-top: 16px" @click="$router.push('/grades/create')"
        >Выставить оценку</n-button
      >
    </div>
    <template v-else>
      <div class="grade-list">
        <GradeCard
          v-for="g in store.list"
          :key="g.id"
          :grade="g"
          @click="$router.push(`/grades/${g.id}`)"
        />
      </div>
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
import { NButton, NIcon, NSelect, NSpin, NEmpty, NPagination } from 'naive-ui'
import { AddOutline } from '@vicons/ionicons5'
import { disciplinesApi } from '~/api/disciplines'
import { gradeTypesApi } from '~/api/grade-types'
import type { Discipline, GradeType } from '~/types'

definePageMeta({ middleware: 'auth' })

const store = useGradesStore()
const auth = useAuthStore()
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

const statusOptions = [
  { label: 'Ожидает', value: 'pending' },
  { label: 'Оценено', value: 'graded' },
  { label: 'Не аттестован', value: 'not_certified' },
  { label: 'Аттестован', value: 'certified' },
  { label: 'Завершено', value: 'closed' },
]
const gradeTypeOptions = computed(() => gradeTypes.value.map((g) => ({ label: g.name, value: g.id })))
const disciplineOptions = computed(() => disciplines.value.map((d) => ({ label: d.name, value: d.id })))

function applyFilters() {
  store.setFilters({
    status: filters.status || '',
    grade_type_id: filters.grade_type_id,
    discipline_id: filters.discipline_id,
  })
}
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
.center {
  @include flex(row, center, center);
  padding: 60px;
}
</style>
